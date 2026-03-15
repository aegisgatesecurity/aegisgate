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
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o aegisgate ./cmd/aegisgate

# =============================================================================
# Runtime Stage
# =============================================================================
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata curl

# Create non-root user
RUN adduser -D -u 1000 aegisgate

# Set working directory
WORKDIR /home/aegisgate

# Copy binary from builder
COPY --from=builder /build/aegisgate .
COPY --from=builder /build/config ./config

# Create data directories
RUN mkdir -p /data /logs && chown -R aegisgate:aegisgate /home/aegisgate

# Switch to non-root user
USER aegisgate

# Expose ports
EXPOSE 8080 8443 9090 3000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Default command
ENTRYPOINT ["./aegisgate"]
CMD ["serve", "--config", "./config/community.env"]
