FROM golang:1.25.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o /bin/recommendation-service \
    ./cmd/recommendation-service

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

COPY --from=builder /bin/recommendation-service /bin/recommendation-service

USER nobody

EXPOSE 50051

ENTRYPOINT ["/bin/recommendation-service"]