from __future__ import annotations

from pathlib import Path
from typing import List

from fastapi import APIRouter, Depends, BackgroundTasks, HTTPException, Query
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select, func

from ..config import settings
from ..db import get_db_session
from ..models import FileRecord, ScanRun
from ..schemas import HealthResponse, ScanRequest, ScanRunOut, FileRecordOut
from ..services.scanner import run_scan

router = APIRouter()


@router.get("/health", response_model=HealthResponse)
async def health() -> HealthResponse:
    return HealthResponse(status="ok", service=settings.app_name)


@router.post("/scan", response_model=ScanRunOut)
async def scan_repo(
    req: ScanRequest,
    background_tasks: BackgroundTasks,
    session: AsyncSession = Depends(get_db_session),
) -> ScanRunOut:
    root = Path(req.root_path)
    if not root.exists() or not root.is_dir():
        raise HTTPException(status_code=400, detail="Invalid root_path")

    follow_symlinks = req.follow_symlinks if req.follow_symlinks is not None else settings.scan_follow_symlinks
    max_workers = req.max_workers if req.max_workers is not None else settings.scan_max_workers

    # Create a ScanRun now and run in background
    scan = await run_scan(
        session=session,
        root=root,
        max_workers=max_workers,
        follow_symlinks=follow_symlinks,
        compute_hash=req.compute_hash,
    )
    return ScanRunOut.model_validate(scan)


@router.get("/files", response_model=List[FileRecordOut])
async def list_files(
    session: AsyncSession = Depends(get_db_session),
    scan_run_id: int | None = Query(default=None),
    q: str | None = Query(default=None),
    limit: int = Query(default=100, ge=1, le=1000),
    offset: int = Query(default=0, ge=0),
) -> list[FileRecordOut]:
    stmt = select(FileRecord)
    if scan_run_id is not None:
        stmt = stmt.where(FileRecord.scan_run_id == scan_run_id)
    if q:
        like = f"%{q}%"
        stmt = stmt.where(FileRecord.path.ilike(like))
    stmt = stmt.order_by(FileRecord.id).limit(limit).offset(offset)

    result = await session.execute(stmt)
    rows = result.scalars().all()
    return [FileRecordOut.model_validate(r) for r in rows]


@router.get("/stats")
async def stats(
    session: AsyncSession = Depends(get_db_session),
) -> dict:
    total_files = await session.scalar(select(func.count()).select_from(FileRecord))
    total_bytes = await session.scalar(select(func.coalesce(func.sum(FileRecord.size_bytes), 0)))
    scans = await session.scalar(select(func.count()).select_from(ScanRun))
    return {
        "total_files": int(total_files or 0),
        "total_bytes": int(total_bytes or 0),
        "scan_runs": int(scans or 0),
    }