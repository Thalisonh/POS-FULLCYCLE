# Start from the official Go image
FROM golang:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Set environment variable for the views directory
ENV DB_DRIVER=mysql
ENV DB_HOST=host.docker.internal
ENV DB_PORT=3306
ENV DB_USER=root
ENV DB_PASSWORD=root
ENV DB_NAME=orders
ENV WEB_SERVER_PORT=8080
ENV GRPC_SERVER_PORT=8081
ENV GRAPHQL_SERVER_PORT=8082
ENV CGO_ENABLED=1

# Build the Go app
RUN go build -o /app/cmd/main /app/cmd/main.go /app/cmd/wire_gen.go 
# Expose port 8080 to the outside world

from scratch
COPY --from=builder /app/cmd/main .

EXPOSE 8080
EXPOSE 8081
EXPOSE 8082

# Run the executable
CMD ["/app/cmd/main"]