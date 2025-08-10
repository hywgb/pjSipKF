from __future__ import annotations

import asyncio
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from .config import settings
from .logging_config import configure_logging, get_logger
from .db import create_all
from .api.routes import router as api_router


def create_app() -> FastAPI:
    configure_logging()
    logger = get_logger("startup")

    app = FastAPI(title=settings.app_name)

    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )

    @app.on_event("startup")
    async def on_startup() -> None:
        logger.info("app_starting", environment=settings.environment)
        # In dev/test, auto create tables for SQLite
        await create_all()

    app.include_router(api_router, prefix="/api")
    return app


app = create_app()