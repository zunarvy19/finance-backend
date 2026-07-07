# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git tzdata

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# CGO_ENABLED=0 creates a statically linked binary which runs perfectly on alpine
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Install tzdata for timezone support
RUN apk add --no-cache tzdata

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/bin/api .

# Expose port (matches APP_PORT in .env)
EXPOSE 5050

# Command to run the executable
CMD ["./api"]
