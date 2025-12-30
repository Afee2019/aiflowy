# AIFlowy Go Backend - Multi-stage Dockerfile
# Target image size: < 50MB

# ============================================
# Stage 1: Build
# ============================================
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build arguments
ARG VERSION=dev
ARG BUILD_TIME=unknown
ARG COMMIT_SHA=unknown

# Build the application
# CGO_ENABLED=0: Static binary, no C dependencies
# -ldflags: Embed version info and strip debug symbols
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-X 'main.version=${VERSION}' -X 'main.buildTime=${BUILD_TIME}' -X 'main.commitSHA=${COMMIT_SHA}' -s -w" \
    -o /build/aiflowy-go \
    ./cmd/server

# ============================================
# Stage 2: Runtime
# ============================================
FROM alpine:3.19

# Labels
LABEL maintainer="AIFlowy Team"
LABEL org.opencontainers.image.title="AIFlowy Go Backend"
LABEL org.opencontainers.image.description="AIFlowy AI Application Platform - Go Backend"
LABEL org.opencontainers.image.version="${VERSION}"

# Install runtime dependencies
# - ca-certificates: For HTTPS connections to external APIs
# - tzdata: Timezone data for time operations
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -g 1000 -S aiflowy \
    && adduser -u 1000 -S aiflowy -G aiflowy

# Set timezone
ENV TZ=Asia/Shanghai

# Create necessary directories
RUN mkdir -p /app/configs /app/logs /app/uploads \
    && chown -R aiflowy:aiflowy /app

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/aiflowy-go /app/aiflowy-go

# Copy default config (can be overridden by ConfigMap)
COPY --from=builder /build/configs/config.prod.yaml /app/configs/config.yaml

# Change ownership
RUN chown -R aiflowy:aiflowy /app

# Switch to non-root user
USER aiflowy

# Expose port
EXPOSE 8213

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8213/health || exit 1

# Run the application
ENTRYPOINT ["/app/aiflowy-go"]
CMD ["-config", "/app/configs/config.yaml"]
