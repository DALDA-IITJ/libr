# FROM golang:1.24.1-alpine AS builder

# WORKDIR /app

# # Set environment variables for build
# ENV GO111MODULE=on \
#     CGO_ENABLED=0 \
#     GOOS=linux \
#     GOARCH=amd64

# # Cache Go modules installation
# COPY ./go.mod ./go.sum ./
# RUN go mod download

# # Copy the source code and build the app
# COPY . .
# # RUN apk add --no-cache postgresql-client sqlite
# # RUN apk add --no-cache git && \
# #     go install github.com/air-verse/air@latest

# # RUN go mod tidy 
# RUN go build -o node .

# FROM alpine:latest

# WORKDIR /app

# # Install runtime dependencies
# # RUN apk add --no-cache postgresql-client sqlite

# # Copy only the binary from the builder image
# COPY --from=builder /app/node .

# # Ensure the binary has execute permissions
# # RUN chmod +x /app/node

# # Expose the node port
# EXPOSE 8080

# # Use absolute path for the executable
# CMD ["./node"]



# Builder stage
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Set environment variables for build
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Install Git (needed for go install)
RUN apk add --no-cache git

# Cache Go modules installation
COPY ./go.mod ./go.sum ./
RUN go mod download

# Copy source code and build the app binary
COPY . .
RUN go build -o node .

# Install Air for live reloading
RUN go install github.com/air-verse/air@latest

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the built application and Air binary from the builder stage
COPY --from=builder /app/node .

# Expose your application port
EXPOSE 8080

# Use Air as the command to start the development server
# CMD ["air", "-c",".air.toml"]

CMD ["air", "-c", ".air.toml"]

# Builder stage: Build the application and install Air
# FROM golang:1.24.1-alpine AS builder

# WORKDIR /app

# # Set environment variables for build
# ENV GO111MODULE=on \
#     CGO_ENABLED=0 \
#     GOOS=linux \
#     GOARCH=amd64

# # Install Git (needed for go install)
# RUN apk add --no-cache git

# # Cache Go modules installation
# COPY ./go.mod ./go.sum ./
# RUN go mod download

# # Copy the source code
# COPY . .

# # Install Air (hot-reload tool)
# RUN go install github.com/air-verse/air@latest

# # Optionally build the binary (Air will rebuild on file changes)
# RUN go build -o node .

# # Final stage: Use a base image that includes Go so Air can build on changes
# FROM golang:1.24.1-alpine

# WORKDIR /app

# # Copy the entire source code from the builder stage
# COPY --from=builder /app /app

# # Expose the port your application listens on
# EXPOSE 8080

# # Start Air for live reloading
# CMD ["air"]
