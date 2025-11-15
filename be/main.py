import logging
import os

from fastapi import FastAPI

from app import tables  # noqa: F401  # ensure table metadata is registered
from app.db import engine, metadata, wait_for_db
from app.routers import ROUTERS
from app.seed_data import seed_initial_data

logger = logging.getLogger("server-be")

app = FastAPI(title="MAX Bot API", version="1.0.0")

for router in ROUTERS:
    app.include_router(router)


RESET_DB_ON_STARTUP = os.getenv("RESET_DB_ON_STARTUP", "true").lower() in {"1", "true", "yes"}


@app.on_event("startup")
async def startup_event() -> None:
    """Ensure all database tables exist before serving requests."""
    await wait_for_db()
    async with engine.begin() as conn:
        if RESET_DB_ON_STARTUP:
            await conn.run_sync(metadata.drop_all)
        await conn.run_sync(metadata.create_all)
    await seed_initial_data()
