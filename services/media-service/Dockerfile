FROM golang:1.23-alpine AS builder

WORKDIR /workspace

# Install build dependencies
RUN apk add --no-cache git

# Copy shared pkg first
COPY pkg/ /workspace/pkg/

# Copy service files
COPY services/media-service/go.mod services/media-service/go.sum /workspace/services/media-service/
WORKDIR /workspace/services/media-service

# Download dependencies
RUN go mod download

# Copy service source
COPY services/media-service/ /workspace/services/media-service/

# Build the service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/media-service ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/media-service .

EXPOSE 8082

CMD ["./media-service"]