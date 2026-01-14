FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
COPY pkg/ ./pkg/
COPY services/discovery-service/ ./services/discovery-service/

# Download dependencies
WORKDIR /app/services/discovery-service
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /discovery-service ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /discovery-service .

EXPOSE 8083

CMD ["./discovery-service"]