# Multi-stage build for minimal image
FROM golang:1.25.8-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod ./
COPY go.sum ./

# Verify go mod files were copied
RUN ls -la

# Copy source code
COPY cmd/ ./cmd/
COPY pkg/ ./pkg/
COPY config/ ./config/

# Verify source code was copied
RUN ls -la ./cmd/aegisgate/ && ls -la ./config/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags="-s -w" \
    -o /build/aegisgate \
    ./cmd/aegisgate

# Verify the binary was built
RUN ls -la /build/aegisgate

# Final stage - minimal runtime image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 aegisgate && \
    adduser -u 1001 -G aegisgate -s /bin/false -D aegisgate

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/aegisgate /app/aegisgate

# Verify the binary was copied
RUN ls -la /app/aegisgate

# Copy configuration file
COPY config/aegisgate.yml /app/aegisgate.yml

# Verify the config file was copied
RUN ls -la /app/aegisgate.yml

# Use non-root user
USER aegisgate

# Expose ports (IPv4 only)
EXPOSE 8080/tcp
EXPOSE 8443/tcp
EXPOSE 8444/tcp

# Set environment variables to force IPv4 binding
ENV AEGISGATE_BIND="0.0.0.0:8080" \
    AEGISGATE_LOG_LEVEL="debug" \
    AEGISGATE_COMPLIANCE_ENABLED="true" \
    AEGISGATE_TIER="community"

# Run the AegisGate binary with IPv4 binding
ENTRYPOINT ["/app/aegisgate"]
CMD ["-bind", "0.0.0.0:8080", "-target", "https://api.openai.com", "-tier", "community"]
