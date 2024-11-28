# Use the official Golang image as the base image
FROM golang:1.20-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first to cache dependencies
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the rest of the application code
COPY . .

# Ensure SQLite3 and build tools are installed (Alpine package)
RUN apk add --no-cache build-base sqlite-dev

# Build the Go app
RUN go build -o froum ./main.go

# Expose the correct port (1703)
EXPOSE 1703

# Health check to ensure the service is running
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --spider http://localhost:1703 || exit 1

# Command to run the executable
CMD ["./froum"]