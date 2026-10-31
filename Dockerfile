#base image as builder
FROM golang:1.27-alpine AS builder


WORKDIR /app

# Copy dependency manifest
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy application source
COPY . .

# Compile the Go application
RUN CGO_ENABLED=0 GOOS=linux \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /bulk-parser .


FROM alpine:3.22 AS runtime

# Install runtime certificates and create application user
RUN apk add --no-cache ca-certificates && \
    addgroup -S appgroup && \
    adduser -S -G appgroup appuser  #CA certificates?They help applications verify the identities of servers during TLS connections.
#Create a non-root user

WORKDIR /app

# Copy only the compiled executable
COPY --from=builder --chown=appuser:appgroup \
    /bulk-parser /app/bulk-parser

# Run as non-root
USER appuser:appgroup

# Document application port
EXPOSE 8081

# Graceful termination signal
STOPSIGNAL SIGTERM

# Start application
ENTRYPOINT ["/app/bulk-parser"]







