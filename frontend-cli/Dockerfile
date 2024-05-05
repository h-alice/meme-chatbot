FROM golang:1.22.2-bullseye

# Set the Current Working Directory inside the container
WORKDIR /app

COPY ./src .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main

# Command to run the executable
CMD ["./main"]