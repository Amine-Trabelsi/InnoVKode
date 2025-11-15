from fastapi import APIRouter, Depends
from sqlalchemy import func, select
from sqlalchemy.ext.asyncio import AsyncSession

from ..db import get_session
from ..tables import (
    ai_queries,
    ai_summaries,
    course_enrollments,
    courses_table,
    teaching_attendance,
    teaching_feedback,
    users_table,
)

router = APIRouter(prefix="/api/v1/dashboard", tags=["Administration"])


@router.get("/overview")
async def overview(session: AsyncSession = Depends(get_session)) -> dict:
    students = (
        select(func.count())
        .select_from(users_table)
        .where(users_table.c.role == "student")
        .scalar_subquery()
    )
    employees = (
        select(func.count())
        .select_from(users_table)
        .where(users_table.c.role == "employee")
        .scalar_subquery()
    )
    leadership = (
        select(func.count())
        .select_from(users_table)
        .where(users_table.c.role == "leadership")
        .scalar_subquery()
    )
    courses = select(func.count()).select_from(courses_table).scalar_subquery()
    enrollments = select(func.count()).select_from(course_enrollments).scalar_subquery()

    result = await session.execute(select(students, employees, leadership, courses, enrollments))
    row = result.first()
    return {
        "students": row[0],
        "employees": row[1],
        "leadership": row[2],
        "courses": row[3],
        "enrollments": row[4],
    }


@router.get("/attendance")
async def attendance(session: AsyncSession = Depends(get_session)) -> dict:
    total_sessions_query = select(func.count()).select_from(teaching_attendance)
    total_sessions = (await session.execute(total_sessions_query)).scalar_one()
    latest_entries = await session.execute(
        select(teaching_attendance)
        .order_by(teaching_attendance.c.session_date.desc())
        .limit(10)
    )
    return {
        "sessions_recorded": total_sessions,
        "recent": [dict(row) for row in latest_entries.mappings().all()],
    }


@router.get("/feedback")
async def feedback_summary(session: AsyncSession = Depends(get_session)) -> list[dict]:
    query = (
        select(
            teaching_feedback.c.course_id,
            func.avg(teaching_feedback.c.rating).label("avg_rating"),
            func.count(teaching_feedback.c.id).label("responses"),
        )
        .group_by(teaching_feedback.c.course_id)
    )
    result = await session.execute(query)
    return [dict(row) for row in result.mappings().all()]


@router.get("/ai/insights")
async def ai_insights(session: AsyncSession = Depends(get_session)) -> dict:
    query_count = select(func.count()).select_from(ai_queries)
    summary_count = select(func.count()).select_from(ai_summaries)
    row = (await session.execute(select(query_count.scalar_subquery(), summary_count.scalar_subquery()))).first()
    return {"queries": row[0], "summaries": row[1]}
