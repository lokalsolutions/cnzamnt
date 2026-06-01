FROM golang:1.22-bookworm AS build

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN go build -o /cnzamnt-api ./cmd/server

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=build /cnzamnt-api /usr/local/bin/cnzamnt-api

ENV CNZAMNT_ADDR=:8080
ENV CNZAMNT_DB_PATH=/data/cnzamnt.db

EXPOSE 8080

CMD ["cnzamnt-api"]
