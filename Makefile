.PHONY: backend-dev backend-test backend-build

backend-dev:
	cd backend && go run ./cmd/server

backend-test:
	cd backend && go test ./...

backend-build:
	cd backend && go build -o bin/cnzamnt-api ./cmd/server
