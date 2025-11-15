from fastapi import APIRouter, Depends, HTTPException
from pydantic import BaseModel
from sqlalchemy import insert, select, update
from sqlalchemy.ext.asyncio import AsyncSession

from ..db import get_session
from ..tables import (
    admission_applications,
    admission_documents,
    admission_event_bookings,
    admission_events,
    admission_faq_queries,
    admission_programs,
)

router = APIRouter(prefix="/api/v1/admissions", tags=["Admissions"])


@router.get("/programs")
async def list_programs(session: AsyncSession = Depends(get_session)) -> list[dict]:
    result = await session.execute(select(admission_programs).order_by(admission_programs.c.title))
    return [dict(row) for row in result.mappings().all()]


@router.get("/programs/{program_id}")
async def get_program(program_id: int, session: AsyncSession = Depends(get_session)) -> dict:
    query = select(admission_programs).where(admission_programs.c.id == program_id)
    row = (await session.execute(query)).mappings().first()
    if not row:
        raise HTTPException(status_code=404, detail="Program not found")
    return dict(row)


@router.get("/events")
async def list_admission_events(session: AsyncSession = Depends(get_session)) -> list[dict]:
    query = select(admission_events).order_by(admission_events.c.date_time)
    result = await session.execute(query)
    return [dict(row) for row in result.mappings().all()]


class EventBookingPayload(BaseModel):
    applicant_name: str
    email: str
    phone: str | None = None
    note: str | None = None


@router.post("/events/{event_id}/book")
async def book_event_seat(event_id: int, payload: EventBookingPayload, session: AsyncSession = Depends(get_session)) -> dict:
    """Book a seat for an admission event (e.g., Open Doors)."""
    # Check if event exists
    event_query = select(admission_events).where(admission_events.c.id == event_id)
    event_result = await session.execute(event_query)
    event_row = event_result.mappings().first()
    if not event_row:
        raise HTTPException(status_code=404, detail="Event not found")

    # Check capacity
    current_attendees = event_row.get("current_attendees") or 0
    max_attendees = event_row.get("max_attendees")
    if max_attendees and current_attendees >= max_attendees:
        raise HTTPException(status_code=409, detail="Event is fully booked")

    # Check if already booked
    existing_booking_query = select(admission_event_bookings).where(
        admission_event_bookings.c.event_id == event_id,
        admission_event_bookings.c.email == payload.email,
    )
    existing_result = await session.execute(existing_booking_query)
    if existing_result.first():
        raise HTTPException(status_code=409, detail="You have already booked this event")

    # Create booking
    stmt = (
        insert(admission_event_bookings)
        .values(
            event_id=event_id,
            applicant_name=payload.applicant_name,
            email=payload.email,
            phone=payload.phone,
            note=payload.note,
        )
        .returning(admission_event_bookings.c.id)
    )
    insert_result = await session.execute(stmt)

    # Update event attendee count
    update_stmt = (
        update(admission_events)
        .where(admission_events.c.id == event_id)
        .values(current_attendees=current_attendees + 1)
    )
    await session.execute(update_stmt)
    await session.commit()

    return {"booking_id": insert_result.scalar_one(), "status": "confirmed"}


@router.get("/events/{event_id}/bookings")
async def list_event_bookings(event_id: int, session: AsyncSession = Depends(get_session)) -> list[dict]:
    """List all bookings for a specific admission event."""
    # Check if event exists
    event_query = select(admission_events).where(admission_events.c.id == event_id)
    event_result = await session.execute(event_query)
    if not event_result.first():
        raise HTTPException(status_code=404, detail="Event not found")

    # Get bookings
    query = select(admission_event_bookings).where(admission_event_bookings.c.event_id == event_id).order_by(admission_event_bookings.c.created_at)
    result = await session.execute(query)
    return [dict(row) for row in result.mappings().all()]


class ApplicationPayload(BaseModel):
    applicant_name: str
    email: str
    program_id: int | None = None
    details: dict | None = None


@router.post("/applications")
async def submit_application(payload: ApplicationPayload, session: AsyncSession = Depends(get_session)) -> dict:
    stmt = (
        insert(admission_applications)
        .values(
            applicant_name=payload.applicant_name,
            email=payload.email,
            program_id=payload.program_id,
            details=payload.details,
        )
        .returning(admission_applications.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"application_id": result.scalar_one(), "status": "received"}


@router.get("/status/{application_id}")
async def application_status(application_id: int, session: AsyncSession = Depends(get_session)) -> dict:
    query = select(admission_applications).where(admission_applications.c.id == application_id)
    row = (await session.execute(query)).mappings().first()
    if not row:
        raise HTTPException(status_code=404, detail="Application not found")
    return dict(row)


class FAQPayload(BaseModel):
    question: str


@router.post("/faq/query")
async def ask_faq(payload: FAQPayload, session: AsyncSession = Depends(get_session)) -> dict:
    response = f"Our team will reach out about: {payload.question}"
    stmt = (
        insert(admission_faq_queries)
        .values(question=payload.question, response=response)
        .returning(admission_faq_queries.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"faq_query_id": result.scalar_one(), "answer": response}


class DocumentPayload(BaseModel):
    application_id: int
    file_name: str
    file_type: str
    storage_url: str


@router.post("/upload")
async def upload_document(payload: DocumentPayload, session: AsyncSession = Depends(get_session)) -> dict:
    stmt = (
        insert(admission_documents)
        .values(
            application_id=payload.application_id,
            file_name=payload.file_name,
            file_type=payload.file_type,
            storage_url=payload.storage_url,
        )
        .returning(admission_documents.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"document_id": result.scalar_one(), "status": "uploaded"}
