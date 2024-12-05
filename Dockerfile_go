# Use an official Go runtime as a parent image
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Start a new stage from a minimal image
FROM debian:buster-slim

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/main /usr/local/bin/main

# Expose the port the app runs on
EXPOSE 8080

# Run the binary program produced by `go build`
ENTRYPOINT ["/usr/local/bin/main"]
