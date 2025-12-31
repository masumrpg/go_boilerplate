# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies needed for build
RUN apk add --no-cache git make

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go
# Build migration tool
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate-tool cmd/migrate/main.go

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for external HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Copy binaries from builder
COPY --from=builder /app/main .
COPY --from=builder /app/migrate-tool .

# Copy config and migration files
# Assumes .env is passed via volume or env vars
COPY --from=builder /app/db ./db

# Expose port
EXPOSE 3000

# Run the application
CMD ["./main"]
