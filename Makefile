PYTHON = python3
VENV_BIN = .venv/bin
PIP = $(VENV_BIN)/pip
PYTEST = $(VENV_BIN)/pytest
RUFF = $(VENV_BIN)/ruff
BLACK = $(VENV_BIN)/black
ISORT = $(VENV_BIN)/isort
UVICORN = $(VENV_BIN)/uvicorn

.PHONY: venv install fmt lint test run dev alembic-up alembic-rev docker-up docker-down

venv:
	$(PYTHON) -m venv .venv

install: venv
	$(PIP) install -r requirements.txt

fmt:
	$(BLACK) .
	$(ISORT) .

lint:
	$(RUFF) .

alembic-up:
	$(VENV_BIN)/alembic upgrade head

alembic-rev:
	$(VENV_BIN)/alembic revision --autogenerate -m "auto"

test:
	$(PYTEST) -q --maxfail=1

run:
	$(UVICORN) app.main:app --host 0.0.0.0 --port 8000 --reload

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down -v