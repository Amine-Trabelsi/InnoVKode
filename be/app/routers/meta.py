from datetime import datetime

from fastapi import APIRouter

router = APIRouter(prefix="/api/v1", tags=["Meta"])


@router.get("/healthz")
async def healthz() -> dict:
    """Simple readiness probe."""
    return {"status": "ok", "timestamp": datetime.utcnow().isoformat()}


@router.get("/version")
async def api_version() -> dict:
    """Return static API metadata."""
    return {
        "version": "v1",
        "name": "MAX Bot API",
        "build": "2025.11.0",
    }

