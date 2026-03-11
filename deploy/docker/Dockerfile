# Build Stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o padlock ./cmd/padlock

# =============================================================================
# Runtime Stage
# =============================================================================
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata curl

# Create non-root user
RUN adduser -D -u 1000 padlock

# Set working directory
WORKDIR /home/padlock

# Copy binary from builder
COPY --from=builder /build/padlock .
COPY --from=builder /build/config ./config

# Create data directories
RUN mkdir -p /data /logs && chown -R padlock:padlock /home/padlock

# Switch to non-root user
USER padlock

# Expose ports
EXPOSE 8080 8443 9090 3000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Default command
ENTRYPOINT ["./padlock"]
CMD ["serve", "--config", "./config/community.env"]
