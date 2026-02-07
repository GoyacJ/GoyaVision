# Stage 1: Build Frontend
FROM node:20-alpine AS builder-web

WORKDIR /web

# Enable pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

# Copy package files
COPY web/package.json web/pnpm-lock.yaml ./

# Install dependencies
RUN pnpm install --frozen-lockfile

# Copy source code
COPY web/ ./

# Build frontend
RUN pnpm run build

# Stage 2: Build Backend
FROM golang:1.24-alpine AS builder-backend

WORKDIR /app

# Install build dependencies if needed (none for pure Go)
# RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Copy frontend build artifacts to expected embed location
COPY --from=builder-web /web/dist ./web/dist

# Build binary
# CGO_ENABLED=0 for static binary
# -ldflags="-w -s" to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o goyavision ./cmd/server

# Stage 3: Final Image
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    ffmpeg \
    ca-certificates \
    tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder-backend /app/goyavision .

# Expose port
EXPOSE 8080

# Run application
CMD ["./goyavision"]
