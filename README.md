# Repo Scanner API

一个高性能的并发文件扫描与索引服务，提供 REST API：健康检查、触发扫描、查询文件、统计信息。

## 快速开始（本地）

```
make install
make run
```

打开 `http://localhost:8000/docs` 查看交互式文档。

## 运行测试

```
make test
```

## 使用 Docker + Postgres（生产）

```
make docker-up
```

环境变量（以 `APP_` 前缀）：
- `APP_DATABASE_URL`：数据库连接（默认 `sqlite+aiosqlite:///./data.db`）
- `APP_SCAN_MAX_WORKERS`：扫描线程数
- `APP_LOG_LEVEL`：日志级别

## 主要特性
- FastAPI + SQLAlchemy(Async) + Alembic
- 结构化 JSON 日志（structlog）
- 多线程扫描器（忽略 .gitignore/.scanignore）
- Dockerfile 与 docker-compose 生产部署
- CI/Lint/Format/Test 工具链
