from fastapi import APIRouter, Depends, HTTPException
from pydantic import BaseModel
from sqlalchemy import insert, select, update
from sqlalchemy.ext.asyncio import AsyncSession

from ..db import get_session
from ..tables import visa_applications, visa_documents

router = APIRouter(prefix="/api/v1/visa", tags=["Visa Services"])


@router.get("/applications/{user_id}")
async def get_visa_applications(user_id: int, session: AsyncSession = Depends(get_session)) -> list[dict]:
    query = (
        select(visa_applications)
        .where(visa_applications.c.user_id == user_id)
        .order_by(visa_applications.c.created_at.desc())
    )
    result = await session.execute(query)
    return [dict(row) for row in result.mappings().all()]


class CreateApplicationPayload(BaseModel):
    user_id: int
    application_type: str  # visa_renewal or registration_renewal


@router.post("/applications")
async def create_visa_application(payload: CreateApplicationPayload, session: AsyncSession = Depends(get_session)) -> dict:
    stmt = (
        insert(visa_applications)
        .values(
            user_id=payload.user_id,
            application_type=payload.application_type,
        )
        .returning(visa_applications.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"application_id": result.scalar_one()}


@router.post("/applications/{application_id}/withdraw")
async def withdraw_visa_application(application_id: int, session: AsyncSession = Depends(get_session)) -> dict:
    stmt = (
        update(visa_applications)
        .where(visa_applications.c.id == application_id)
        .values(status="withdrawn")
    )
    result = await session.execute(stmt)
    if result.rowcount == 0:
        raise HTTPException(status_code=404, detail="Application not found")
    await session.commit()
    return {"status": "withdrawn"}


@router.get("/applications/{application_id}/documents")
async def get_visa_documents(application_id: int, session: AsyncSession = Depends(get_session)) -> list[dict]:
    query = (
        select(visa_documents)
        .where(visa_documents.c.application_id == application_id)
        .order_by(visa_documents.c.uploaded_at.desc())
    )
    result = await session.execute(query)
    return [dict(row) for row in result.mappings().all()]


class UploadDocumentPayload(BaseModel):
    file_name: str
    file_url: str


@router.post("/applications/{application_id}/documents")
async def upload_visa_document(application_id: int, payload: UploadDocumentPayload, session: AsyncSession = Depends(get_session)) -> dict:
    stmt = (
        insert(visa_documents)
        .values(
            application_id=application_id,
            file_name=payload.file_name,
            file_url=payload.file_url,
        )
        .returning(visa_documents.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"document_id": result.scalar_one()}