from __future__ import annotations

import asyncio
from typing import AsyncGenerator
from sqlalchemy.ext.asyncio import AsyncSession, create_async_engine, async_sessionmaker
from sqlalchemy.orm import DeclarativeBase
from sqlalchemy import MetaData

from .config import settings


class Base(DeclarativeBase):
    metadata = MetaData()


engine = create_async_engine(settings.database_url, echo=False, pool_pre_ping=True)
SessionLocal = async_sessionmaker(engine, class_=AsyncSession, expire_on_commit=False)

_tables_ready: bool = False
_tables_lock = asyncio.Lock()


async def create_all() -> None:
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)


async def ensure_tables_ready() -> None:
    global _tables_ready
    if _tables_ready:
        return
    async with _tables_lock:
        if _tables_ready:
            return
        await create_all()
        _tables_ready = True


async def get_db_session() -> AsyncGenerator[AsyncSession, None]:
    await ensure_tables_ready()
    async with SessionLocal() as session:
        yield session