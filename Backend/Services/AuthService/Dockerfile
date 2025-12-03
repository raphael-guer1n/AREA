# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY app/go.mod app/go.sum ./app/

# Download dependencies
WORKDIR /build/app
RUN go mod download

# Copy entire project structure
WORKDIR /build
COPY . .

# Build the application (entrypoint in ./app/cmd/main.go)
WORKDIR /build/app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /build/app/main .

# Copy all configuration files from root directory
# This includes: *.json, *.yaml, *.yml, .env* files
COPY --from=builder /build/*.json ./
COPY --from=builder /build/*.yaml ./
COPY --from=builder /build/*.yml ./

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
