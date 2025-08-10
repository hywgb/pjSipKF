from alembic import op
import sqlalchemy as sa

revision = "0001_init"
down_revision = None
branch_labels = None
depends_on = None

def upgrade() -> None:
    op.create_table(
        "scan_runs",
        sa.Column("id", sa.Integer(), primary_key=True, autoincrement=True),
        sa.Column("root_path", sa.String(length=1024), nullable=False),
        sa.Column("started_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("finished_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("success", sa.Boolean(), nullable=False, server_default=sa.text("0")),
        sa.Column("total_files", sa.Integer(), nullable=False, server_default=sa.text("0")),
        sa.Column("total_bytes", sa.BigInteger(), nullable=False, server_default=sa.text("0")),
    )

    op.create_table(
        "file_records",
        sa.Column("id", sa.Integer(), primary_key=True, autoincrement=True),
        sa.Column("scan_run_id", sa.Integer(), sa.ForeignKey("scan_runs.id"), nullable=False, index=True),
        sa.Column("path", sa.String(length=4096), nullable=False),
        sa.Column("size_bytes", sa.BigInteger(), nullable=False, server_default=sa.text("0")),
        sa.Column("mtime_ts", sa.BigInteger(), nullable=False, server_default=sa.text("0")),
        sa.Column("is_binary", sa.Boolean(), nullable=False, server_default=sa.text("0")),
        sa.Column("mime_type", sa.String(length=256), nullable=True),
        sa.Column("sha256", sa.String(length=64), nullable=True),
        sa.UniqueConstraint("scan_run_id", "path", name="uq_scan_file_path"),
    )


def downgrade() -> None:
    op.drop_table("file_records")
    op.drop_table("scan_runs")