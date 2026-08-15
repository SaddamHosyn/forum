# Stage 1: Builder Stage
FROM golang:alpine AS builder

WORKDIR /app

# Install build dependencies including sqlite
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copy go.mod and go.sum first for efficient caching of dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code and assets
COPY . .

# Build with CGO enabled (required for SQLite)
RUN CGO_ENABLED=1 GOOS=linux go build -o reel-movie-talk .

# Stage 2: Runtime Stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies for SQLite
RUN apk add --no-cache sqlite-libs

# Copy the compiled binary from the builder stage
COPY --from=builder /app/reel-movie-talk .

# Copy static assets and templates
COPY --from=builder /app/assets ./assets
COPY --from=builder /app/templates ./templates

# Expose the port your application listens on
EXPOSE 8999

# Command to run your application
CMD ["./reel-movie-talk"]
