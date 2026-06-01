.PHONY: backend-dev backend-test backend-build frontend-dev frontend-build

BACKEND_GOCACHE ?= $(CURDIR)/backend/.gocache

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
