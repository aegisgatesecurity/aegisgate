# Multi-stage build for minimal image
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod ./
COPY go.sum ./

# Copy source code
COPY cmd/ ./cmd/
COPY pkg/ ./pkg/
COPY config/ ./config/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags="-s -w" \
    -o /build/aegisgate \
    ./cmd/aegisgate

# Final stage - minimal runtime image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 aegisgate && \
    adduser -u 1001 -G aegisgate -s /bin/false -D aegisgate

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/aegisgate /app/aegisgate
COPY config/aegisgate.yml.example /app/aegisgate.yml

# Use non-root user
USER aegisgate

# Expose ports
EXPOSE 8080 8443 8444

# Set entrypoint
ENTRYPOINT ["/app/aegisgate"]
CMD ["--config", "/app/aegisgate.yml"]
