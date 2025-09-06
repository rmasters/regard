FROM golang:1-alpine AS builder

# Install ca-certificates for SSL verification
RUN apk add --no-cache ca-certificates git

WORKDIR /src

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o regard cmd/regard/main.go

# Final stage - minimal image
FROM scratch

# Copy CA certificates for HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY --from=builder /src/regard /regard

# Set the binary as entrypoint
ENTRYPOINT ["/regard"]