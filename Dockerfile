# Build stage
FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binaries
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/node ./cmd/node
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/client ./cmd/client

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binaries from builder
COPY --from=builder /app/bin/ .

# Copy configuration
COPY config.example.yaml config.yaml

# Create directories
RUN mkdir -p data/storage data/files logs

# Expose API port
EXPOSE 8080

# Run API server by default
CMD ["./api"]
