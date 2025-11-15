from datetime import datetime

from fastapi import APIRouter, Depends
from pydantic import BaseModel
from sqlalchemy import func, insert, select
from sqlalchemy.ext.asyncio import AsyncSession

from ..db import get_session
from ..tables import (
    course_sessions,
    courses_table,
    teaching_announcements,
    teaching_attendance,
    teaching_feedback,
    teaching_grade_uploads,
    teaching_submissions,
)

router = APIRouter(prefix="/api/v1/teaching", tags=["Professors"])


@router.get("/schedule/{professor_id}")
async def professor_schedule(professor_id: int, session: AsyncSession = Depends(get_session)) -> list[dict]:
    query = (
        select(
            course_sessions.c.id.label("session_id"),
            course_sessions.c.session_type,
            course_sessions.c.start_time,
            course_sessions.c.end_time,
            course_sessions.c.location,
            courses_table.c.id.label("course_id"),
            courses_table.c.title,
        )
        .join(courses_table, courses_table.c.id == course_sessions.c.course_id)
        .where(courses_table.c.teacher_id == professor_id)
        .order_by(course_sessions.c.start_time)
    )
    result = await session.execute(query)
    return [dict(row) for row in result.mappings().all()]


@router.get("/courses/{professor_id}")
async def professor_courses(professor_id: int, session: AsyncSession = Depends(get_session)) -> list[dict]:
    query = select(courses_table).where(courses_table.c.teacher_id == professor_id).order_by(courses_table.c.code)
    result = await session.execute(query)
    return [dict(row) for row in result.mappings().all()]


class AttendancePayload(BaseModel):
    course_id: int
    professor_id: int
    session_date: datetime
    attendance: list[dict]


@router.post("/attendance")
async def submit_attendance(payload: AttendancePayload, session: AsyncSession = Depends(get_session)) -> dict:
    stmt = (
        insert(teaching_attendance)
        .values(
            course_id=payload.course_id,
            professor_id=payload.professor_id,
            session_date=payload.session_date,
            attendance=payload.attendance,
        )
        .returning(teaching_attendance.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"attendance_id": result.scalar_one()}


class GradeUploadPayload(BaseModel):
    course_id: int
    professor_id: int
    grades: list[dict]


@router.post("/grades/upload")
async def upload_grades(payload: GradeUploadPayload, session: AsyncSession = Depends(get_session)) -> dict:
    stmt = (
        insert(teaching_grade_uploads)
        .values(course_id=payload.course_id, professor_id=payload.professor_id, payload=payload.grades)
        .returning(teaching_grade_uploads.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"upload_id": result.scalar_one()}


@router.get("/submissions/{course_id}")
async def list_submissions(course_id: int, session: AsyncSession = Depends(get_session)) -> list[dict]:
    query = (
        select(teaching_submissions)
        .where(teaching_submissions.c.course_id == course_id)
        .order_by(teaching_submissions.c.submitted_at.desc())
    )
    result = await session.execute(query)
    return [dict(row) for row in result.mappings().all()]


class AnnouncementPayload(BaseModel):
    course_id: int
    professor_id: int
    message: str


@router.post("/announcements")
async def create_announcement(payload: AnnouncementPayload, session: AsyncSession = Depends(get_session)) -> dict:
    stmt = (
        insert(teaching_announcements)
        .values(course_id=payload.course_id, professor_id=payload.professor_id, message=payload.message)
        .returning(teaching_announcements.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"announcement_id": result.scalar_one()}


@router.get("/feedback/{course_id}")
async def course_feedback(course_id: int, session: AsyncSession = Depends(get_session)) -> dict:
    """Aggregate feedback for a course."""
    stats = (
        select(
            func.count(teaching_feedback.c.id).label("responses"),
            func.avg(teaching_feedback.c.rating).label("avg_rating"),
        )
        .where(teaching_feedback.c.course_id == course_id)
    )
    stats_row = (await session.execute(stats)).mappings().first() or {"responses": 0, "avg_rating": None}
    feedback_rows = await session.execute(
        select(teaching_feedback)
        .where(teaching_feedback.c.course_id == course_id)
        .order_by(teaching_feedback.c.submitted_at.desc())
    )
    return {
        "course_id": course_id,
        "responses": stats_row["responses"] or 0,
        "avg_rating": float(stats_row["avg_rating"]) if stats_row["avg_rating"] is not None else None,
        "items": [dict(row) for row in feedback_rows.mappings().all()],
    }

