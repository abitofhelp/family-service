FROM golang:1.24-alpine3.21 AS builder

ENV CGO_ENABLED=1

# Install build dependencies for CGO
RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY dev.docker.env ./
COPY "cmd/" "./cmd/"
COPY config/ ./config/
COPY core/ ./core/
RUN mkdir -p /app/data/dev/sqlite
COPY data/dev/sqlite/family_service.db /app/data/dev/sqlite/family_service.db
COPY infrastructure/ ./infrastructure/
COPY interface/ ./interface/

RUN go build -o family_service "./cmd/server/graphql"

FROM alpine:3.19

LABEL maintainer="mjgardner@abitofhelp.com"
LABEL version="1.1.0"
LABEL description="Family Service application"

# Add necessary runtime dependencies
RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /app
COPY --from=builder /app/family_service .
COPY --from=builder /app/dev.docker.env .
COPY --from=builder /app/config ./config
RUN mkdir -p /app/data/dev/sqlite
COPY --from=builder /app/data/dev/sqlite/family_service.db /app/data/dev/sqlite/family_service.db
COPY entrypoint.sh .
COPY secrets ./secrets

RUN chmod +x /app/entrypoint.sh
RUN mkdir -p /app/secrets && chmod -R 755 /app/secrets
RUN chmod -R 755 /app/data/dev/sqlite && chmod 644 /app/data/dev/sqlite/family_service.db

RUN adduser -D appuser
USER appuser

EXPOSE 8089
ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["./family_service"]

HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget --quiet --tries=1 --spider http://localhost:8089/health || exit 1
