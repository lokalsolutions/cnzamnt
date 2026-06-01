-include .env
export

.PHONY: dev prod up down logs logs-backend logs-frontend logs-caddy backend-dev backend-test backend-build frontend-dev frontend-build test-e2e require-prod-env

BACKEND_GOCACHE ?= $(CURDIR)/backend/.gocache

require-prod-env:
	@if [ -z "$(CNZAMNT_DOMAIN)" ] || [ -z "$(CNZAMNT_API_DOMAIN)" ]; then \
		echo "CNZAMNT_DOMAIN and CNZAMNT_API_DOMAIN are required. Copy .env.example to .env and edit it."; \
		exit 2; \
	fi

dev:
	docker compose up --build -d

prod: require-prod-env
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

logs-caddy:
	docker compose logs -f caddy

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

test-e2e:
	cd frontend && npm run test:e2e
