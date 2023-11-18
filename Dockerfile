# Stage 1: Build stage
FROM golang:1.21 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download and install Go dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -o server -a -tags netgo -ldflags '-w -extldflags "-static"' ./cmd/server/main.go

# Stage 2: Final stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage to the final stage
COPY --from=builder /app .

# Install necessary libraries for SQLite3
RUN apk --no-cache add sqlite

# Expose any necessary ports
EXPOSE 50051

# Command to run the executable
CMD ["/app/server"]
