FROM golang:alpine AS builder

# Set working directory to match the Go module path
WORKDIR /go/src/github.com/Ahmed1monm/Axis-BE-assessment

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/server

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Add ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /go/src/github.com/Ahmed1monm/Axis-BE-assessment/app .

# Copy environment file
COPY .env* ./

# Expose the application port
EXPOSE 8080

CMD ["./app"]
