# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api cmd/api/main.go

# Run stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/bin/api .

# Expose port
EXPOSE 8080

# Run application
CMD ["./api"]