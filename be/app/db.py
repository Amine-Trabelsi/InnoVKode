import asyncio
import logging
import os
from typing import AsyncGenerator

from sqlalchemy import MetaData, select
from sqlalchemy.exc import SQLAlchemyError
from sqlalchemy.ext.asyncio import AsyncEngine, AsyncSession, async_sessionmaker, create_async_engine

DATABASE_URL = os.getenv(
    "DATABASE_URL",
    "postgresql+asyncpg://app_user:app_password@db:5432/app_db",
)

logger = logging.getLogger("server-be")

metadata = MetaData()
engine: AsyncEngine = create_async_engine(DATABASE_URL, echo=False, future=True)
AsyncSessionLocal = async_sessionmaker(engine, expire_on_commit=False, class_=AsyncSession)


async def wait_for_db(max_attempts: int = 10, delay_seconds: float = 1.5) -> None:
    """Retry database connectivity before continuing app startup."""
    for attempt in range(1, max_attempts + 1):
        try:
            async with engine.connect() as conn:
                await conn.execute(select(1))
            logger.info("Database connection established.")
            return
        except (SQLAlchemyError, OSError) as exc:
            if attempt == max_attempts:
                logger.error("Failed to connect to DB after %s attempts: %s", max_attempts, exc)
                raise
            logger.warning("DB not ready (attempt %s/%s): %s", attempt, max_attempts, exc)
            await asyncio.sleep(delay_seconds)


async def get_session() -> AsyncGenerator[AsyncSession, None]:
    """Provide a scoped AsyncSession for FastAPI dependencies."""
    async with AsyncSessionLocal() as session:
        yield session

