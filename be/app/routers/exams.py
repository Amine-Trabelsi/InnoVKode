import logging
from fastapi import APIRouter, Depends
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from ..db import get_session
from ..tables import course_enrollments, courses_table, exam_schedules, grade_records

logger = logging.getLogger("server-be")
router = APIRouter(prefix="/api/v1", tags=["Exams & Grades"])


@router.get("/exams/{student_id}")
async def get_exams(student_id: int, session: AsyncSession = Depends(get_session)) -> list[dict]:
    """Return upcoming exams for the student's enrolled courses."""
    query = (
        select(
            exam_schedules.c.id.label("exam_id"),
            exam_schedules.c.exam_date,
            exam_schedules.c.room,
            exam_schedules.c.exam_format,
            courses_table.c.id.label("course_id"),
            courses_table.c.code,
            courses_table.c.title,
        )
        .join(courses_table, courses_table.c.id == exam_schedules.c.course_id)
        .join(course_enrollments, course_enrollments.c.course_id == courses_table.c.id)
        .where(course_enrollments.c.student_id == student_id)
        .order_by(exam_schedules.c.exam_date)
    )
    result = await session.execute(query)
    return [dict(row) for row in result.mappings().all()]


@router.get("/grades/{student_id}")
async def get_grades(student_id: int, session: AsyncSession = Depends(get_session)) -> list[dict]:
    """Return recorded grades and GPA contributions."""
    query = (
        select(
            grade_records.c.id.label("grade_id"),
            grade_records.c.grade,
            grade_records.c.gpa_points,
            grade_records.c.graded_on,
            courses_table.c.id.label("course_id"),
            courses_table.c.code,
            courses_table.c.title,
        )
        .join(courses_table, courses_table.c.id == grade_records.c.course_id)
        .where(grade_records.c.student_id == student_id)
        .order_by(grade_records.c.graded_on.desc())
    )
    result = await session.execute(query)
    rows = [dict(row) for row in result.mappings().all()]
    for row in rows:
        row['gpa_points'] = float(row['gpa_points'])
        logger.info(f"Grade record: gpa_points type={type(row['gpa_points'])}, value={row['gpa_points']}")
    return rows

