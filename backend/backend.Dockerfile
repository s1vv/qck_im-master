# Используем официальный образ Go
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . . 

RUN go build -o main ./cmd/main.go

# Финальный образ
FROM debian:latest

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
