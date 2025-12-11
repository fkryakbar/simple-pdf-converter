# ============================================
# Stage 1: Builder
# ============================================
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod and sum files first (for better layer caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations for smaller binary
# CGO_ENABLED=0: Static linking (no C dependencies)
# -ldflags: Strip debug symbols and reduce binary size
# -trimpath: Remove file system paths from binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static'" \
    -trimpath \
    -o /app/server .

# ============================================
# Stage 2: Production (Minimal Runtime)
# ============================================
FROM scratch

# Import certificates for HTTPS calls
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Set timezone (optional, adjust as needed)
ENV TZ=Asia/Jakarta

# Copy the binary from builder
COPY --from=builder /app/server /server

# Expose port (configurable via PORT env var)
EXPOSE 8080

# Run the application
ENTRYPOINT ["/server"]
