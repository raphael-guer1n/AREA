# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY app/go.mod app/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY app/ .

# Build the application (entrypoint in ./cmd/main.go)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /build/main .

# Copy provider config files
COPY --from=builder /build/internal/config/providers ./internal/config/providers
COPY --from=builder /build/internal/config/services ./internal/config/services

# Run the application
CMD ["./main"]
