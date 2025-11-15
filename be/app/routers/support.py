from fastapi import APIRouter, Depends, HTTPException
from pydantic import BaseModel
from sqlalchemy import insert, select
from sqlalchemy.ext.asyncio import AsyncSession

from ..db import get_session
from ..tables import ai_advisor_sessions, support_queries, support_tickets

router = APIRouter(prefix="/api/v1", tags=["Support & AI Assistants"])


class SupportQueryPayload(BaseModel):
    user_id: int | None = None
    question: str


@router.post("/support/query")
async def support_query(payload: SupportQueryPayload, session: AsyncSession = Depends(get_session)) -> dict:
    answer = f"Our support team will respond regarding: {payload.question}"
    stmt = (
        insert(support_queries)
        .values(user_id=payload.user_id, question=payload.question, answer=answer)
        .returning(support_queries.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"query_id": result.scalar_one(), "answer": answer}


class SupportTicketPayload(BaseModel):
    user_id: int | None = None
    category: str
    subject: str
    description: str | None = None


@router.post("/support/tickets")
async def create_ticket(payload: SupportTicketPayload, session: AsyncSession = Depends(get_session)) -> dict:
    stmt = (
        insert(support_tickets)
        .values(
            user_id=payload.user_id,
            category=payload.category,
            subject=payload.subject,
            description=payload.description,
        )
        .returning(support_tickets.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"ticket_id": result.scalar_one(), "status": "open"}


@router.get("/support/tickets/{ticket_id}")
async def get_ticket(ticket_id: int, session: AsyncSession = Depends(get_session)) -> dict:
    query = select(support_tickets).where(support_tickets.c.id == ticket_id)
    row = (await session.execute(query)).mappings().first()
    if not row:
        raise HTTPException(status_code=404, detail="Ticket not found")
    return dict(row)


class AdvisorPayload(BaseModel):
    user_id: int | None = None
    topic: str | None = None
    prompt: str


@router.post("/ai/chat/advisor")
async def advisor_chat(payload: AdvisorPayload, session: AsyncSession = Depends(get_session)) -> dict:
    response = f"Advisor tip for {payload.topic or 'general guidance'}"
    stmt = (
        insert(ai_advisor_sessions)
        .values(user_id=payload.user_id, topic=payload.topic, prompt=payload.prompt, response=response)
        .returning(ai_advisor_sessions.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"session_id": result.scalar_one(), "response": response}

