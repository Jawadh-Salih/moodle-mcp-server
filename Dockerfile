FROM golang:1.22-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o moodle-mcp ./cmd/moodle-mcp/

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/moodle-mcp .

# Port for REST API
ENV PORT=8080
EXPOSE 8080

# Default to REST mode
CMD ["./moodle-mcp", "-mode", "rest", "-port", "8080"]
