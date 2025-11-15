from datetime import datetime

from fastapi import APIRouter, Depends
from pydantic import BaseModel, Field
from sqlalchemy import insert, select
from sqlalchemy.ext.asyncio import AsyncSession

from ..db import get_session
from ..tables import (
    business_trip_requests,
    hr_certificates,
    hr_notifications,
    vacation_requests,
)

router = APIRouter(prefix="/api/v1/hr", tags=["Employees"])


@router.get("/vacations/{employee_id}")
async def get_vacations(employee_id: int, session: AsyncSession = Depends(get_session)) -> list[dict]:
    query = (
        select(vacation_requests)
        .where(vacation_requests.c.employee_id == employee_id)
        .order_by(vacation_requests.c.created_at.desc())
    )
    result = await session.execute(query)
    return [dict(row) for row in result.mappings().all()]


class VacationPayload(BaseModel):
    employee_id: int
    start_date: datetime
    end_date: datetime
    vacation_type: str = "paid"


@router.post("/vacations/request")
async def request_vacation(payload: VacationPayload, session: AsyncSession = Depends(get_session)) -> dict:
    stmt = (
        insert(vacation_requests)
        .values(
            employee_id=payload.employee_id,
            start_date=payload.start_date,
            end_date=payload.end_date,
            vacation_type=payload.vacation_type,
        )
        .returning(vacation_requests.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"vacation_request_id": result.scalar_one()}


@router.get("/business_trips/{employee_id}")
async def get_business_trips(employee_id: int, session: AsyncSession = Depends(get_session)) -> list[dict]:
    query = (
        select(business_trip_requests)
        .where(business_trip_requests.c.employee_id == employee_id)
        .order_by(business_trip_requests.c.created_at.desc())
    )
    result = await session.execute(query)
    return [dict(row) for row in result.mappings().all()]


class BusinessTripPayload(BaseModel):
    employee_id: int
    destination: str
    start_date: datetime
    end_date: datetime
    purpose: str | None = None


@router.post("/business_trips/request")
async def request_business_trip(payload: BusinessTripPayload, session: AsyncSession = Depends(get_session)) -> dict:
    stmt = (
        insert(business_trip_requests)
        .values(
            employee_id=payload.employee_id,
            destination=payload.destination,
            start_date=payload.start_date,
            end_date=payload.end_date,
            purpose=payload.purpose,
        )
        .returning(business_trip_requests.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"business_trip_id": result.scalar_one()}


@router.get("/certificates/{employee_id}")
async def get_certificates(employee_id: int, session: AsyncSession = Depends(get_session)) -> list[dict]:
    query = (
        select(hr_certificates)
        .where(hr_certificates.c.employee_id == employee_id)
        .order_by(hr_certificates.c.requested_at.desc())
    )
    result = await session.execute(query)
    return [dict(row) for row in result.mappings().all()]


class CertificatePayload(BaseModel):
    employee_id: int
    certificate_type: str = Field(..., description="employment|income|custom")


@router.post("/certificates/request")
async def request_certificate(payload: CertificatePayload, session: AsyncSession = Depends(get_session)) -> dict:
    stmt = (
        insert(hr_certificates)
        .values(employee_id=payload.employee_id, certificate_type=payload.certificate_type)
        .returning(hr_certificates.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"certificate_request_id": result.scalar_one()}


@router.get("/notifications")
async def hr_notifications_feed(session: AsyncSession = Depends(get_session)) -> list[dict]:
    result = await session.execute(select(hr_notifications).order_by(hr_notifications.c.published_at.desc()))
    return [dict(row) for row in result.mappings().all()]

