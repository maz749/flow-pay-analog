# Build stage
FROM golang:1.21-alpine AS builder

# Force cache bust - increment this to force rebuild
ARG CACHEBUST=2

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata postgresql-client bash

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy web files
COPY --from=builder /app/web ./web

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Copy start script
COPY start.sh .
RUN chmod +x start.sh

# Expose port
EXPOSE 8080

# Run the application with start script
CMD ["./start.sh"]
