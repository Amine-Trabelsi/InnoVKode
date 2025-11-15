from fastapi import APIRouter, Depends, HTTPException
from pydantic import BaseModel
from sqlalchemy import insert, select
from sqlalchemy.ext.asyncio import AsyncSession

from ..db import get_session
from ..tables import dorm_payments, dorm_requests, dorm_rooms

router = APIRouter(prefix="/api/v1/dorms", tags=["Dormitories"])


@router.get("/rooms/{student_id}")
async def dorm_room(student_id: int, session: AsyncSession = Depends(get_session)) -> dict:
    query = select(dorm_rooms).where(dorm_rooms.c.student_id == student_id)
    row = (await session.execute(query)).mappings().first()
    if not row:
        raise HTTPException(status_code=404, detail="Dorm assignment not found")
    return dict(row)


class MaintenancePayload(BaseModel):
    student_id: int
    request_type: str
    description: str | None = None


@router.post("/maintenance")
async def create_maintenance(payload: MaintenancePayload, session: AsyncSession = Depends(get_session)) -> dict:
    stmt = (
        insert(dorm_requests)
        .values(
            student_id=payload.student_id,
            request_type=payload.request_type,
            description=payload.description,
        )
        .returning(dorm_requests.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"request_id": result.scalar_one(), "status": "open"}


@router.get("/requests/{request_id}")
async def get_request(request_id: int, session: AsyncSession = Depends(get_session)) -> dict:
    query = select(dorm_requests).where(dorm_requests.c.id == request_id)
    row = (await session.execute(query)).mappings().first()
    if not row:
        raise HTTPException(status_code=404, detail="Request not found")
    return dict(row)


class PaymentPayload(BaseModel):
    student_id: int
    amount: float
    reference: str | None = None


@router.post("/payments")
async def submit_payment(payload: PaymentPayload, session: AsyncSession = Depends(get_session)) -> dict:
    stmt = (
        insert(dorm_payments)
        .values(student_id=payload.student_id, amount=payload.amount, reference=payload.reference)
        .returning(dorm_payments.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"payment_id": result.scalar_one(), "status": "recorded"}

