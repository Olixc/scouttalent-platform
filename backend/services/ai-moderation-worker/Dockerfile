FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
COPY pkg/ ./pkg/
COPY services/ai-moderation-worker/ ./services/ai-moderation-worker/

# Download dependencies
WORKDIR /app/services/ai-moderation-worker
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /ai-moderation-worker ./cmd/worker

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /ai-moderation-worker .

CMD ["./ai-moderation-worker"]