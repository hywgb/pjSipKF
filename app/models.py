from __future__ import annotations

from datetime import datetime
from sqlalchemy.orm import Mapped, mapped_column
from sqlalchemy import Integer, String, DateTime, BigInteger, Boolean, ForeignKey, UniqueConstraint
from sqlalchemy.orm import relationship

from .db import Base


class ScanRun(Base):
    __tablename__ = "scan_runs"

    id: Mapped[int] = mapped_column(Integer, primary_key=True, autoincrement=True)
    root_path: Mapped[str] = mapped_column(String(1024), nullable=False)
    started_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), default=datetime.utcnow)
    finished_at: Mapped[datetime | None] = mapped_column(DateTime(timezone=True))
    success: Mapped[bool] = mapped_column(Boolean, default=False, nullable=False)
    total_files: Mapped[int] = mapped_column(Integer, default=0, nullable=False)
    total_bytes: Mapped[int] = mapped_column(BigInteger, default=0, nullable=False)

    files: Mapped[list[FileRecord]] = relationship("FileRecord", back_populates="scan_run")


class FileRecord(Base):
    __tablename__ = "file_records"
    __table_args__ = (
        UniqueConstraint("scan_run_id", "path", name="uq_scan_file_path"),
    )

    id: Mapped[int] = mapped_column(Integer, primary_key=True, autoincrement=True)
    scan_run_id: Mapped[int] = mapped_column(ForeignKey("scan_runs.id"), nullable=False, index=True)
    path: Mapped[str] = mapped_column(String(4096), nullable=False)
    size_bytes: Mapped[int] = mapped_column(BigInteger, default=0, nullable=False)
    mtime_ts: Mapped[int] = mapped_column(BigInteger, default=0, nullable=False)
    is_binary: Mapped[bool] = mapped_column(Boolean, default=False, nullable=False)
    mime_type: Mapped[str | None] = mapped_column(String(256))
    sha256: Mapped[str | None] = mapped_column(String(64))

    scan_run: Mapped[ScanRun] = relationship("ScanRun", back_populates="files")