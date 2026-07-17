FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies for CGO (SQLite)
RUN apk add --no-cache gcc musl-dev

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o accountabel main.go

FROM alpine:latest

WORKDIR /app

# Copy the built binary
COPY --from=builder /app/accountabel /app/accountabel

# Copy static assets and configurations
COPY --from=builder /app/static /app/static
COPY --from=builder /app/conf /app/conf
COPY --from=builder /app/i18n /app/i18n

# Create data directory for SQLite
RUN mkdir -p /app/data

# Expose the default port
EXPOSE 36060

# Run the binary
CMD /app/accountabel -type=gemini -http_host=:${PORT:-36060}
