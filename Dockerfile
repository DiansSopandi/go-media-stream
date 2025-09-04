# Stage 1: Build
FROM golang:1.24.3 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o goride main.go

# Stage 2: Run
FROM alpine:3.19
WORKDIR /app

RUN apk add --no-cache ca-certificates && adduser -D appuser
USER appuser

COPY --from=builder /app/goride .
COPY env.conf env.conf

EXPOSE 8001
CMD ["./goride"]