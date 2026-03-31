FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd


FROM alpine:latest

WORKDIR /app

# security
RUN adduser -D appuser

COPY --from=builder /app/main .

RUN chown appuser:appuser main

USER appuser

EXPOSE 8080

CMD ["./main"]