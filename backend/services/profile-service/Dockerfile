# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /workspace

# Copy pkg directory to match local structure
COPY pkg ./pkg

# Copy service files to services/profile-service
COPY services/profile-service ./services/profile-service

# Set working directory to service
WORKDIR /workspace/services/profile-service

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/profile-service ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/profile-service .
COPY --from=builder /workspace/services/profile-service/migrations ./migrations

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 8081

CMD ["./profile-service"]