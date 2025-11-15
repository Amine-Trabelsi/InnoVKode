from fastapi import APIRouter, Depends
from pydantic import BaseModel, Field
from sqlalchemy import insert
from sqlalchemy.ext.asyncio import AsyncSession

from ..db import get_session
from ..tables import ai_queries, ai_quizzes, ai_sources, ai_summaries, ai_transcriptions

router = APIRouter(prefix="/api/v1/ai", tags=["AI Module"])


class SourceUpload(BaseModel):
    source_type: str = Field(..., description="youtube|pdf|doc|web")
    reference: str
    title: str | None = None
    metadata: dict | None = None


@router.post("/rag/upload")
async def register_source(payload: SourceUpload, session: AsyncSession = Depends(get_session)) -> dict:
    stmt = (
        insert(ai_sources)
        .values(
            source_type=payload.source_type,
            reference=payload.reference,
            title=payload.title,
            metadata=payload.metadata,
        )
        .returning(ai_sources.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"source_id": result.scalar_one()}


class QueryPayload(BaseModel):
    question: str
    filters: dict | None = None


@router.post("/rag/query")
async def query_sources(payload: QueryPayload, session: AsyncSession = Depends(get_session)) -> dict:
    """Store the query and return a canned response until RAG is wired up."""
    response_text = f"Knowledge base lookup placeholder for: {payload.question}"
    stmt = (
        insert(ai_queries)
        .values(query_text=payload.question, response_text=response_text)
        .returning(ai_queries.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"query_id": result.scalar_one(), "answer": response_text}


class QuizRequest(BaseModel):
    course_id: int | None = None
    prompt: str


@router.post("/quiz/generate")
async def generate_quiz(payload: QuizRequest, session: AsyncSession = Depends(get_session)) -> dict:
    """Persist quiz metadata and return simple stub questions."""
    questions = [
        {
            "question": f"{payload.prompt} - concept check",
            "options": ["A", "B", "C", "D"],
            "answer": "A",
        }
    ]
    stmt = (
        insert(ai_quizzes)
        .values(course_id=payload.course_id, prompt=payload.prompt, questions=questions)
        .returning(ai_quizzes.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"quiz_id": result.scalar_one(), "questions": questions}


class SummaryRequest(BaseModel):
    source_text: str


@router.post("/summary/create")
async def create_summary(payload: SummaryRequest, session: AsyncSession = Depends(get_session)) -> dict:
    summary = payload.source_text[:200] + ("..." if len(payload.source_text) > 200 else "")
    stmt = (
        insert(ai_summaries)
        .values(source=payload.source_text, summary=summary)
        .returning(ai_summaries.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"summary_id": result.scalar_one(), "summary": summary}


class TranscriptionRequest(BaseModel):
    audio_ref: str


@router.post("/audio/transcribe")
async def transcribe_audio(payload: TranscriptionRequest, session: AsyncSession = Depends(get_session)) -> dict:
    transcript = f"Transcription placeholder for {payload.audio_ref}"
    stmt = (
        insert(ai_transcriptions)
        .values(audio_ref=payload.audio_ref, transcript=transcript)
        .returning(ai_transcriptions.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"transcription_id": result.scalar_one(), "transcript": transcript}
