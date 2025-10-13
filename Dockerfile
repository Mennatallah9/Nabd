# Build stage for frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm ci --only=production

COPY frontend/ ./
RUN npm run build

# Build stage for backend
FROM golang:1.21-alpine AS backend-builder

# Install build dependencies for CGO
RUN apk --no-cache add gcc musl-dev

# Set Go proxy and sum database
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org

WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./

# Download dependencies with retry mechanism
RUN go mod download || \
    (sleep 5 && go mod download) || \
    (sleep 10 && go mod download)

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

# Copy frontend build files
COPY --from=frontend-builder /app/frontend/build ./web

# Create static directory and copy static assets
RUN mkdir -p ./static && \
    cp -r ./web/static/* ./static/ && \
    cp ./web/index.html ./static/ && \
    cp ./web/*.png ./static/ 2>/dev/null || true && \
    cp ./web/*.ico ./static/ 2>/dev/null || true

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