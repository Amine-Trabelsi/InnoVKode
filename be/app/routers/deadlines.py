from datetime import datetime

from fastapi import APIRouter, Depends
from pydantic import BaseModel, Field
from sqlalchemy import insert, select
from sqlalchemy.ext.asyncio import AsyncSession

from ..db import get_session
from ..tables import deadlines_table, notifications_table

router = APIRouter(prefix="/api/v1", tags=["Deadlines & Notifications"])


@router.get("/deadlines/{student_id}")
async def list_deadlines(student_id: int, session: AsyncSession = Depends(get_session)) -> list[dict]:
    """Return all current deadlines affecting the student."""
    query = (
        select(deadlines_table)
        .where(deadlines_table.c.student_id == student_id)
        .order_by(deadlines_table.c.due_date)
    )
    result = await session.execute(query)
    return [dict(row) for row in result.mappings().all()]


class NotificationRequest(BaseModel):
    recipient_id: int | None = Field(default=None, description="Optional specific user target")
    subject: str
    body: str
    channel: str = "in_app"


@router.post("/notifications/send")
async def send_notification(payload: NotificationRequest, session: AsyncSession = Depends(get_session)) -> dict:
    """Persist a notification entry for follow-up delivery."""
    stmt = (
        insert(notifications_table)
        .values(
            recipient_id=payload.recipient_id,
            subject=payload.subject,
            body=payload.body,
            channel=payload.channel,
            created_at=datetime.utcnow(),
        )
        .returning(notifications_table.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    new_id = result.scalar_one()
    return {"notification_id": new_id, "status": "queued"}

