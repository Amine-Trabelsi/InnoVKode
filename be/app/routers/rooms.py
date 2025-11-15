from datetime import datetime

from fastapi import APIRouter, Depends, HTTPException, Query
from pydantic import BaseModel
from sqlalchemy import insert, select, update
from sqlalchemy.ext.asyncio import AsyncSession

from ..db import get_session
from ..tables import room_bookings, rooms_table

router = APIRouter(prefix="/api/v1", tags=["Rooms"])


@router.get("/rooms/available")
async def available_rooms(
    session: AsyncSession = Depends(get_session),
    start_time: datetime | None = Query(default=None),
    end_time: datetime | None = Query(default=None),
) -> list[dict]:
    """Return rooms, optionally filtered by availability in the provided slot."""
    query = select(rooms_table)

    if start_time and end_time:
        conflicting = (
            select(room_bookings.c.room_id)
            .where(room_bookings.c.status == "confirmed")
            .where(room_bookings.c.start_time < end_time)
            .where(room_bookings.c.end_time > start_time)
        ).scalar_subquery()
        query = query.where(~rooms_table.c.id.in_(conflicting))

    result = await session.execute(query.order_by(rooms_table.c.capacity.desc()))
    return [dict(row) for row in result.mappings().all()]


class RoomBookingRequest(BaseModel):
    room_id: int
    user_id: int
    start_time: datetime
    end_time: datetime
    purpose: str | None = None


@router.post("/rooms/book")
async def book_room(payload: RoomBookingRequest, session: AsyncSession = Depends(get_session)) -> dict:
    """Attempt to reserve a room slot."""
    conflict_query = (
        select(room_bookings.c.id)
        .where(room_bookings.c.room_id == payload.room_id)
        .where(room_bookings.c.status == "confirmed")
        .where(room_bookings.c.start_time < payload.end_time)
        .where(room_bookings.c.end_time > payload.start_time)
    )
    conflict = await session.execute(conflict_query)
    if conflict.first():
        raise HTTPException(status_code=409, detail="Room already booked for that slot")

    stmt = (
        insert(room_bookings)
        .values(
            room_id=payload.room_id,
            user_id=payload.user_id,
            start_time=payload.start_time,
            end_time=payload.end_time,
            purpose=payload.purpose,
        )
        .returning(room_bookings.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"booking_id": result.scalar_one(), "status": "confirmed"}


@router.delete("/rooms/bookings/{booking_id}")
async def cancel_booking(booking_id: int, session: AsyncSession = Depends(get_session)) -> dict:
    """Cancel an existing booking."""
    stmt = (
        update(room_bookings)
        .where(room_bookings.c.id == booking_id)
        .values(status="cancelled")
        .returning(room_bookings.c.id)
    )
    result = await session.execute(stmt)
    row = result.first()
    if not row:
        raise HTTPException(status_code=404, detail="Booking not found")
    await session.commit()
    return {"booking_id": booking_id, "status": "cancelled"}
