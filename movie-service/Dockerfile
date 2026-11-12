FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o /bin/movie-service \
    ./cmd/movie-service

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /bin/movie-service /bin/movie-service
COPY --from=builder /app/migrations ./migrations

EXPOSE 50051

CMD ["/bin/movie-service"]