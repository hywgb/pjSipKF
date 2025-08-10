import os
import pytest
from httpx import AsyncClient
from fastapi import status

from app.main import app


@pytest.mark.asyncio
async def test_scan_workspace(tmp_path):
    # Create a small temp project
    d = tmp_path / "proj"
    d.mkdir()
    f1 = d / "a.txt"
    f1.write_text("hello", encoding="utf-8")
    f2 = d / "bin.bin"
    f2.write_bytes(b"\x00\x01\x02")

    async with AsyncClient(app=app, base_url="http://test") as ac:
        resp = await ac.post("/api/scan", json={"root_path": str(d), "compute_hash": True})
        assert resp.status_code == status.HTTP_200_OK
        run = resp.json()
        assert run["success"] is True
        assert run["total_files"] >= 2

        # Query files
        resp2 = await ac.get("/api/files", params={"scan_run_id": run["id"]})
        assert resp2.status_code == 200
        items = resp2.json()
        paths = {os.path.basename(x["path"]) for x in items}
        assert {"a.txt", "bin.bin"}.issubset(paths)