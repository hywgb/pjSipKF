from __future__ import annotations

import hashlib
import os
from concurrent.futures import ThreadPoolExecutor, as_completed
from pathlib import Path
from typing import Iterable
import mimetypes

from pathspec import PathSpec
from pathspec.patterns import GitWildMatchPattern
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select

from ..logging_config import get_logger
from ..models import FileRecord, ScanRun

logger = get_logger("scanner")


def load_gitignore(root: Path) -> PathSpec:
    patterns: list[str] = []
    for name in [".scanignore", ".gitignore"]:
        p = root / name
        if p.exists():
            try:
                patterns.extend(p.read_text(encoding="utf-8", errors="ignore").splitlines())
            except Exception:
                continue
    return PathSpec.from_lines(GitWildMatchPattern, patterns)


def is_binary_file(sample: bytes) -> bool:
    if b"\0" in sample:
        return True
    try:
        sample.decode("utf-8")
        return False
    except UnicodeDecodeError:
        return True


def compute_sha256(path: Path) -> str:
    h = hashlib.sha256()
    with path.open("rb") as f:
        for chunk in iter(lambda: f.read(1024 * 1024), b""):
            h.update(chunk)
    return h.hexdigest()


def iter_files(root: Path, follow_symlinks: bool, ignore: PathSpec) -> Iterable[Path]:
    for dirpath, dirnames, filenames in os.walk(root, followlinks=follow_symlinks):
        rel_dir = Path(dirpath).relative_to(root)
        # Filter ignored directories
        pruned = []
        for d in list(dirnames):
            rel = (rel_dir / d).as_posix()
            if ignore.match_file(rel):
                pruned.append(d)
        for d in pruned:
            dirnames.remove(d)

        for fn in filenames:
            rel = (rel_dir / fn).as_posix()
            if ignore.match_file(rel):
                continue
            yield Path(dirpath) / fn


def summarize_file(path: Path, compute_hash_flag: bool) -> dict:
    try:
        stat = path.stat()
        size = stat.st_size
        mtime = int(stat.st_mtime)
        mime, _ = mimetypes.guess_type(path.as_posix())
        sha = None
        is_bin = False
        with path.open("rb") as f:
            head = f.read(4096)
            is_bin = is_binary_file(head)
        if compute_hash_flag and size <= 1024 * 1024 * 100:  # cap at 100MB
            sha = compute_sha256(path)
        return {
            "path": str(path),
            "size_bytes": size,
            "mtime_ts": mtime,
            "is_binary": is_bin,
            "mime_type": mime,
            "sha256": sha,
        }
    except Exception as e:
        logger.warning("summarize_file_failed", path=str(path), error=str(e))
        return {
            "path": str(path),
            "size_bytes": 0,
            "mtime_ts": 0,
            "is_binary": True,
            "mime_type": None,
            "sha256": None,
        }


async def run_scan(session: AsyncSession, root: Path, max_workers: int, follow_symlinks: bool, compute_hash: bool) -> ScanRun:
    scan = ScanRun(root_path=str(root))
    session.add(scan)
    await session.flush()

    ignore = load_gitignore(root)

    paths = list(iter_files(root, follow_symlinks=follow_symlinks, ignore=ignore))

    total_bytes = 0
    files_to_add: list[FileRecord] = []

    with ThreadPoolExecutor(max_workers=max_workers) as executor:
        futures = {executor.submit(summarize_file, p, compute_hash): p for p in paths}
        for fut in as_completed(futures):
            data = fut.result()
            total_bytes += data["size_bytes"]
            files_to_add.append(FileRecord(
                scan_run_id=scan.id,
                path=data["path"],
                size_bytes=data["size_bytes"],
                mtime_ts=data["mtime_ts"],
                is_binary=data["is_binary"],
                mime_type=data["mime_type"],
                sha256=data["sha256"],
            ))

    scan.total_files = len(files_to_add)
    scan.total_bytes = int(total_bytes)
    scan.success = True

    session.add_all(files_to_add)
    await session.commit()

    # Refresh to get finished_at if set elsewhere later
    await session.refresh(scan)
    return scan