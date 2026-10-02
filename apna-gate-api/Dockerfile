ARG GO_VERSION=1.26.3

FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/worker ./cmd/worker

FROM postgres:15-alpine

RUN apk add --no-cache ca-certificates tzdata wget \
    && pg_dump --version | grep -q '^pg_dump (PostgreSQL) 15\.' \
    && pg_restore --version | grep -q '^pg_restore (PostgreSQL) 15\.' \
    && pg_amcheck --version | grep -q '^pg_amcheck (PostgreSQL) 15\.' \
    && addgroup -S appgroup \
    && adduser -S -D -H -u 10001 -G appgroup appuser \
    && mkdir -p /app /var/log/society /var/lib/apna-gate/backups/tmp \
    && chown -R appuser:appgroup /app /var/log/society /var/lib/apna-gate/backups \
    && chmod 700 /var/lib/apna-gate/backups/tmp

WORKDIR /app

COPY --from=builder --chown=appuser:appgroup /out/server ./server
COPY --from=builder --chown=appuser:appgroup /out/worker ./worker
COPY --from=builder --chown=appuser:appgroup /src/migrations ./migrations

USER appuser

# This is the API container; do not run the PostgreSQL image's database entrypoint.
ENTRYPOINT []

EXPOSE 8080

CMD ["./server"]
