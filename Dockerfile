# Use the official Go image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go server source code into the container
COPY . .

# Build the Go server binary
RUN go build -o app .

# Expose the port your Go server is listening on
EXPOSE 8000

# Define the command to run your Go server when the container starts
CMD ["./app"]
