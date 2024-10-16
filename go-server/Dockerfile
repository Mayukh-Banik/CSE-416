# Use Go image to build the Go server
FROM golang:1.23.2-alpine

# Set the working directory inside the container
WORKDIR /go-server

# Copy go.mod and go.sum to download dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the Go server code to the container
COPY . .

# Build the Go server
RUN go build -o go-server .

# Expose the Go server port
EXPOSE 8080

# Run the Go server
CMD ["./go-server"]
