# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN go install github.com/air-verse/air@latest

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]
