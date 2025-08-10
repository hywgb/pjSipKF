from pydantic_settings import BaseSettings, SettingsConfigDict
from pydantic import Field
from pathlib import Path


class Settings(BaseSettings):
    model_config = SettingsConfigDict(env_file=(".env", ".env.local"), env_prefix="APP_", case_sensitive=False)

    environment: str = Field(default="development")
    app_name: str = Field(default="repo-scanner")
    log_level: str = Field(default="INFO")
    http_host: str = Field(default="0.0.0.0")
    http_port: int = Field(default=8000)

    database_url: str = Field(default="sqlite+aiosqlite:///./data.db")
    run_migrations: bool = Field(default=False)

    scan_default_root: Path = Field(default=Path("/workspace"))
    scan_max_workers: int = Field(default=8)
    scan_follow_symlinks: bool = Field(default=False)



settings = Settings()