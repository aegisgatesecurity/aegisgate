#!/bin/bash
# Multi-platform Docker build script for AegisGate
# Builds for: linux/amd64, linux/arm64, linux/arm/v7

set -e

VERSION=${VERSION:-v0.39.0}
REGISTRY=${REGISTRY:-ghcr.io/aegisgatesecurity}
IMAGE_NAME=aegisgate

echo "Building AegisGate ${VERSION} for multiple platforms..."

# Build multi-platform images
docker buildx build \
    --platform linux/amd64,linux/arm64,linux/arm/v7 \
    -t ${REGISTRY}/${IMAGE_NAME}:${VERSION} \
    -t ${REGISTRY}/${IMAGE_NAME}:latest \
    -f deploy/docker/Dockerfile.multiplatform \
    --push \
    .

echo "✅ Multi-platform build complete!"
echo "   - linux/amd64"
echo "   - linux/arm64" 
echo "   - linux/arm/v7"
echo ""
echo "Images pushed to:"
echo "   - ${REGISTRY}/${IMAGE_NAME}:${VERSION}"
echo "   - ${REGISTRY}/${IMAGE_NAME}:latest"
