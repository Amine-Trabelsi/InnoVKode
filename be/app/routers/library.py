from fastapi import APIRouter, Depends, HTTPException, Query
from pydantic import BaseModel
from sqlalchemy import insert, or_, select
from sqlalchemy.ext.asyncio import AsyncSession

from ..db import get_session
from ..tables import (
    library_books,
    library_digital_assets,
    library_loans,
    library_reservations,
)

router = APIRouter(prefix="/api/v1/library", tags=["Library"])


@router.get("/books/search")
async def search_books(
    q: str = Query(""),
    session: AsyncSession = Depends(get_session),
) -> list[dict]:
    query = select(library_books)
    if q:
        like_value = f"%{q}%"
        query = query.where(
            or_(
                library_books.c.title.ilike(like_value),
                library_books.c.author.ilike(like_value),
            )
        )
    result = await session.execute(query.order_by(library_books.c.title))
    return [dict(row) for row in result.mappings().all()]


class ReservationPayload(BaseModel):
    book_id: int
    student_id: int


@router.post("/books/reserve")
async def reserve_book(payload: ReservationPayload, session: AsyncSession = Depends(get_session)) -> dict:
    stmt = (
        insert(library_reservations)
        .values(book_id=payload.book_id, student_id=payload.student_id)
        .returning(library_reservations.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"reservation_id": result.scalar_one()}


@router.get("/borrowed/{student_id}")
async def borrowed_books(student_id: int, session: AsyncSession = Depends(get_session)) -> list[dict]:
    query = (
        select(
            library_loans.c.id.label("loan_id"),
            library_loans.c.borrowed_at,
            library_loans.c.due_at,
            library_loans.c.status,
            library_books.c.id.label("book_id"),
            library_books.c.title,
            library_books.c.author,
        )
        .join(library_books, library_books.c.id == library_loans.c.book_id)
        .where(library_loans.c.student_id == student_id)
        .order_by(library_loans.c.borrowed_at.desc())
    )
    result = await session.execute(query)
    return [dict(row) for row in result.mappings().all()]


@router.get("/digital/{book_id}")
async def digital_asset(book_id: int, session: AsyncSession = Depends(get_session)) -> dict:
    query = select(library_digital_assets).where(library_digital_assets.c.book_id == book_id)
    row = (await session.execute(query)).mappings().first()
    if not row:
        raise HTTPException(status_code=404, detail="Digital copy not available")
    return dict(row)

