# Build stage for frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm ci --only=production

COPY frontend/ ./
RUN npm run build

# Build stage for backend
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o nabd main.go

# Final stage
FROM alpine:latest

# Install required packages
RUN apk --no-cache add ca-certificates sqlite curl

# Create app directory
WORKDIR /app

# Copy backend binary
COPY --from=backend-builder /app/backend/nabd .

# Copy frontend build
COPY --from=frontend-builder /app/frontend/build ./static

# Create data directory for SQLite
RUN mkdir -p /data

# Copy example config
COPY config.yaml.example /app/config.yaml

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Set environment variables
ENV GIN_MODE=release
ENV NABD_DB_PATH=/data/nabd.db

# Run the application
CMD ["./nabd"]