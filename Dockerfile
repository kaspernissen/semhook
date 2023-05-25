# Start from a minimal Golang image
FROM golang:1.17-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Install the StarHook binary
RUN go install github.com/fatih/starhook@latest

# Copy the source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 go build -o myapp .

# Start a new stage for the final minimal image
FROM alpine:latest

# Install semgrep
RUN apk --no-cache add curl && \
    curl -L https://github.com/returntocorp/semgrep/releases/latest/download/semgrep.tar.gz | tar zxvf - -C /usr/local/bin && \
    chmod +x /usr/local/bin/semgrep

# Copy the binary from the previous stage
COPY --from=build /go/bin/starhook /usr/local/bin/starhook
COPY --from=build /cmd/main /usr/local/bin/myapp

# Expose port 8080
EXPOSE 8080

# Set the entrypoint for the container
ENTRYPOINT ["main"]
