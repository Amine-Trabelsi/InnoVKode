from fastapi import APIRouter, Depends, HTTPException
from pydantic import BaseModel
from sqlalchemy import delete, insert, select, update
from sqlalchemy.ext.asyncio import AsyncSession

from ..db import get_session
from ..tables import clubs_table, event_registrations, events_table, news_table

router = APIRouter(prefix="/api/v1", tags=["Events & News & Clubs"])


@router.get("/events")
async def list_events(session: AsyncSession = Depends(get_session)) -> list[dict]:
    """List upcoming campus events."""
    query = select(events_table).order_by(events_table.c.date_time)
    result = await session.execute(query)
    return [dict(row) for row in result.mappings().all()]


class RSVPRequest(BaseModel):
    user_id: int
    registration_type: str = "attendee"
    note: str | None = None


@router.post("/events/{event_id}/rsvp")
async def rsvp_event(event_id: int, payload: RSVPRequest, session: AsyncSession = Depends(get_session)) -> dict:
    """Create or update participation in an event."""
    event_query = select(events_table).where(events_table.c.id == event_id)
    event_result = await session.execute(event_query)
    event_row = event_result.mappings().first()
    if not event_row:
        raise HTTPException(status_code=404, detail="Event not found")

    # Check if already registered
    reg_query = select(event_registrations).where(
        event_registrations.c.event_id == event_id,
        event_registrations.c.user_id == payload.user_id
    )
    reg_result = await session.execute(reg_query)
    existing_reg = reg_result.mappings().first()

    current_attendees = event_row.get("current_attendees") or 0
    max_attendees = event_row.get("max_attendees")

    if existing_reg:
        # Already registered, update type
        if existing_reg["registration_type"] == payload.registration_type:
            return {"registration_id": existing_reg["id"], "status": "already_registered"}
        # Update registration type
        update_reg_stmt = (
            update(event_registrations)
            .where(
                event_registrations.c.event_id == event_id,
                event_registrations.c.user_id == payload.user_id
            )
            .values(registration_type=payload.registration_type, note=payload.note)
            .returning(event_registrations.c.id)
        )
        update_result = await session.execute(update_reg_stmt)
        await session.commit()
        return {"registration_id": update_result.scalar_one(), "status": "updated"}
    else:
        # New registration
        if max_attendees and current_attendees >= max_attendees:
            raise HTTPException(status_code=409, detail="Event is fully booked")

        stmt = (
            insert(event_registrations)
            .values(
                event_id=event_id,
                user_id=payload.user_id,
                registration_type=payload.registration_type,
                note=payload.note,
            )
            .returning(event_registrations.c.id)
        )
        insert_result = await session.execute(stmt)

        update_stmt = (
            update(events_table)
            .where(events_table.c.id == event_id)
            .values(current_attendees=current_attendees + 1)
        )
        await session.execute(update_stmt)
        await session.commit()
        return {"registration_id": insert_result.scalar_one(), "status": "registered"}


class CancelRequest(BaseModel):
    user_id: int


@router.post("/events/{event_id}/cancel")
async def cancel_rsvp(event_id: int, payload: CancelRequest, session: AsyncSession = Depends(get_session)) -> dict:
    """Cancel participation in an event."""
    # Check if registration exists
    reg_query = select(event_registrations).where(
        event_registrations.c.event_id == event_id,
        event_registrations.c.user_id == payload.user_id
    )
    reg_result = await session.execute(reg_query)
    reg_row = reg_result.mappings().first()
    if not reg_row:
        raise HTTPException(status_code=404, detail="Registration not found")

    # Delete registration
    delete_stmt = delete(event_registrations).where(
        event_registrations.c.event_id == event_id,
        event_registrations.c.user_id == payload.user_id
    )
    await session.execute(delete_stmt)

    # Decrement current_attendees
    update_stmt = (
        update(events_table)
        .where(events_table.c.id == event_id)
        .values(current_attendees=events_table.c.current_attendees - 1)
    )
    await session.execute(update_stmt)
    await session.commit()
    return {"status": "cancelled"}


@router.get("/events/user/{user_id}")
async def list_user_events(user_id: int, session: AsyncSession = Depends(get_session)) -> list[dict]:
    """List events the user is registered for."""
    query = (
        select(
            events_table,
            event_registrations.c.registration_type.label("user_registration_type")
        )
        .join(event_registrations, events_table.c.id == event_registrations.c.event_id)
        .where(event_registrations.c.user_id == user_id)
        .order_by(events_table.c.date_time)
    )
    result = await session.execute(query)
    return [dict(row) for row in result.mappings().all()]


@router.get("/news")
async def list_news(session: AsyncSession = Depends(get_session)) -> list[dict]:
    """Return news items ordered by publish date."""
    query = select(news_table).order_by(news_table.c.published_at.desc())
    result = await session.execute(query)
    return [dict(row) for row in result.mappings().all()]


@router.get("/clubs")
async def list_clubs(session: AsyncSession = Depends(get_session)) -> list[dict]:
    """Return available clubs and organizations."""
    result = await session.execute(select(clubs_table).order_by(clubs_table.c.name))
    return [dict(row) for row in result.mappings().all()]
