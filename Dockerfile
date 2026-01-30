# Etapa de compilación
FROM golang:1.25.5-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o badgerStatus ./cmd/badgerStatus/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/badgerStatus .

ENV VOLUME="/"
ENV SERIAL_PORT="/dev/ttyACM0"
ENV BAUD_RATE="115200"
ENV INTERVAL="10s"

ENTRYPOINT ["./badgerStatus"]