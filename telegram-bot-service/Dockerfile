FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o telegram-bot ./cmd/bot/main.go

FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/telegram-bot .
COPY ./migrations ./migrations
ENTRYPOINT ["./telegram-bot"]