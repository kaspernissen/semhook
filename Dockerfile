# Start from a minimal Golang image
FROM golang:1.20.4-alpine3.18 AS build

# Set the working directory inside the container
WORKDIR /deps

# Install the StarHook binary
RUN go install github.com/fatih/starhook/cmd/starhook@latest

WORKDIR /app
# Copy the Go module files and download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o semhook cmd/main.go

# Start a new stage for the final minimal image
FROM returntocorp/semgrep:latest

# Copy the binary from the previous stage
COPY --from=build /go/bin/starhook /usr/local/bin/starhook
COPY --from=build /app/semhook /usr/local/bin/semhook

# Expose port 8080
EXPOSE 8080

# Set the entrypoint for the container
ENTRYPOINT ["semhook"]
