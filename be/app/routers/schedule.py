from fastapi import APIRouter, Depends
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from ..db import get_session
from ..tables import course_enrollments, course_sessions, courses_table

router = APIRouter(prefix="/api/v1", tags=["Schedule & Courses"])


@router.get("/schedule/{student_id}")
async def get_schedule(student_id: int, session: AsyncSession = Depends(get_session)) -> list[dict]:
    """Return chronologically ordered sessions for the student."""
    query = (
        select(
            course_sessions.c.id.label("session_id"),
            course_sessions.c.session_type,
            course_sessions.c.start_time,
            course_sessions.c.end_time,
            course_sessions.c.location,
            course_sessions.c.week_label,
            courses_table.c.id.label("course_id"),
            courses_table.c.code,
            courses_table.c.title,
        )
        .join(course_enrollments, course_enrollments.c.course_id == course_sessions.c.course_id)
        .join(courses_table, courses_table.c.id == course_sessions.c.course_id)
        .where(course_enrollments.c.student_id == student_id)
        .order_by(course_sessions.c.start_time)
    )
    result = await session.execute(query)
    return [dict(row) for row in result.mappings().all()]


@router.get("/courses/{student_id}")
async def get_courses(student_id: int, session: AsyncSession = Depends(get_session)) -> list[dict]:
    """List all courses the student is enrolled in."""
    query = (
        select(courses_table)
        .join(course_enrollments, course_enrollments.c.course_id == courses_table.c.id)
        .where(course_enrollments.c.student_id == student_id)
        .order_by(courses_table.c.title)
    )
    result = await session.execute(query)
    return [dict(row) for row in result.mappings().all()]

