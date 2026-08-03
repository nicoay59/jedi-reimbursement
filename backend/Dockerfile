FROM golang:1.23-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/api \
    ./cmd/api

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/migrate \
    ./cmd/migrate

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/seed \
    ./cmd/seed

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/wait-db \
    ./cmd/wait-db

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -G app app

WORKDIR /app

COPY --from=builder /out/api /app/api
COPY --from=builder /out/migrate /app/migrate
COPY --from=builder /out/seed /app/seed
COPY --from=builder /out/wait-db /app/wait-db
COPY docker-entrypoint.sh /app/docker-entrypoint.sh

RUN mkdir -p /app/storage/uploads/parking-receipts \
    && chmod +x /app/docker-entrypoint.sh \
    && chown -R app:app /app

USER app

EXPOSE 8080

ENTRYPOINT ["/app/docker-entrypoint.sh"]
