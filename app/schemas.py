from __future__ import annotations

from datetime import datetime
from pydantic import BaseModel, Field


class HealthResponse(BaseModel):
    status: str = Field(default="ok")
    service: str


class ScanRequest(BaseModel):
    root_path: str
    compute_hash: bool = False
    follow_symlinks: bool | None = None
    max_workers: int | None = None


class ScanRunOut(BaseModel):
    id: int
    root_path: str
    started_at: datetime
    finished_at: datetime | None
    success: bool
    total_files: int
    total_bytes: int

    class Config:
        from_attributes = True


class FileRecordOut(BaseModel):
    id: int
    path: str
    size_bytes: int
    mtime_ts: int
    is_binary: bool
    mime_type: str | None
    sha256: str | None

    class Config:
        from_attributes = True