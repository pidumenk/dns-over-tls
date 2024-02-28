# Multi-stage build
FROM golang:alpine as builder

# Set working dir to copy files in there
WORKDIR /app

#  Download necessary dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into working dir
COPY main.go ./

# Build Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /godns

#------------------------------------------------------------#

# Use minimal image to optimize the final size (https://hub.docker.com/_/alpine)
FROM alpine

# Copy data from the first stage
COPY --from=builder /godns /

# Specify the port to be exposed
EXPOSE 53/tcp
EXPOSE 53/udp

# Set entrypoint and default TCP argument
ENTRYPOINT ["/godns"]
CMD ["tcp"]