FROM golang:1.22.3

WORKDIR /app

COPY /migrations ./migrations

RUN go install github.com/rubenv/sql-migrate/...@latest
