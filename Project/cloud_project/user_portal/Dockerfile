# Use a minimal base image for efficiency
FROM golang:alpine AS builder

# Set working directory for the build stage
WORKDIR /app

# Copy your source code to the working directory
COPY . .

# Install dependencies (if any)
RUN go mod download

# Build the application 
RUN go build -o main ./cmd/main.go  # Assuming main.go is within the cmd/main.go path

# Use a smaller image for the final container
FROM alpine:latest

# Copy the built binary 
COPY main ./

# Set the entrypoint to run your application directly (no need for ./main)
ENTRYPOINT ["main"]
