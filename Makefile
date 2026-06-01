.PHONY: dev up down logs logs-backend logs-frontend backend-dev backend-test backend-build frontend-dev frontend-build

BACKEND_GOCACHE ?= $(CURDIR)/backend/.gocache

dev:
	docker compose up --build -d

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

logs-backend:
	docker compose logs -f backend

logs-frontend:
	docker compose logs -f frontend

backend-dev:
	cd backend && GOCACHE="$(BACKEND_GOCACHE)" go run ./cmd/server

backend-test:
	cd backend && GOCACHE="$(BACKEND_GOCACHE)" go test ./...

backend-build:
	cd backend && GOCACHE="$(BACKEND_GOCACHE)" go build -o bin/cnzamnt-api ./cmd/server

frontend-dev:
	cd frontend && npm run dev -- --host 0.0.0.0

frontend-build:
	cd frontend && npm run build
