# Build stage
FROM golang:1.24-alpine AS builder

# Install git for fetching dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Install migrate tool
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -o vinylvault ./cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS connections
RUN apk add --no-cache ca-certificates tzdata

# Create app user
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/vinylvault .

# Copy config and migrations
COPY --from=builder /app/internal/config/config.yml ./internal/config/
COPY --from=builder /app/internal/db/migrations ./internal/db/migrations
COPY --from=builder /go/bin/migrate ./migrate

# Use non-root user
USER appuser

# Expose port
EXPOSE 8080

# Command to run
CMD ["./vinylvault"]