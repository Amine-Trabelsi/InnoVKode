from fastapi import APIRouter, Depends, HTTPException, Query
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from ..db import get_session
from ..tables import users_table

router = APIRouter(prefix="/api/v1", tags=["Users"])


@router.get("/users/by-email")
async def get_user_by_email(
    email: str = Query(..., min_length=3),
    session: AsyncSession = Depends(get_session),
) -> dict:
    """Fetch profile data by email address."""
    query = select(users_table).where(users_table.c.email == email)
    result = await session.execute(query)
    row = result.mappings().first()
    if not row:
        raise HTTPException(status_code=404, detail="User not found")
    return dict(row)


@router.get("/users/{student_id}")
async def get_user(student_id: int, session: AsyncSession = Depends(get_session)) -> dict:
    """Fetch the core profile data for a student/user."""
    query = select(users_table).where(users_table.c.id == student_id)
    result = await session.execute(query)
    row = result.mappings().first()
    if not row:
        raise HTTPException(status_code=404, detail="User not found")
    return dict(row)
