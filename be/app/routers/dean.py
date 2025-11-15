from fastapi import APIRouter, Depends, HTTPException
from pydantic import BaseModel
from sqlalchemy import insert, select
from sqlalchemy.ext.asyncio import AsyncSession

from ..db import get_session
from ..tables import dean_requests

router = APIRouter(prefix="/api/v1/dean", tags=["Dean's Office"])


class DeanRequest(BaseModel):
    user_id: int
    request_type: str
    payload: dict


@router.post("/requests")
async def create_dean_request(payload: DeanRequest, session: AsyncSession = Depends(get_session)) -> dict:
    stmt = (
        insert(dean_requests)
        .values(user_id=payload.user_id, request_type=payload.request_type, payload=payload.payload)
        .returning(dean_requests.c.id)
    )
    result = await session.execute(stmt)
    await session.commit()
    return {"request_id": result.scalar_one(), "status": "submitted"}


@router.get("/requests/{request_id}")
async def get_dean_request(request_id: int, session: AsyncSession = Depends(get_session)) -> dict:
    query = select(dean_requests).where(dean_requests.c.id == request_id)
    result = await session.execute(query)
    row = result.mappings().first()
    if not row:
        raise HTTPException(status_code=404, detail="Request not found")
    return dict(row)

