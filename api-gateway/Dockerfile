FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o api-gateway ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN adduser -D -g '' appuser
WORKDIR /home/appuser
COPY --from=builder /app/api-gateway .
COPY .env .env
USER appuser
EXPOSE 80
CMD ["./api-gateway"]
