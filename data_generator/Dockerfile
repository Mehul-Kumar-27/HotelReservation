### The builder stage where we build the go binary of the data generator
FROM golang:1.21.3 AS builder

WORKDIR /app

COPY . .
COPY ./go.mod .
COPY ./go.sum .

ENV GOOS=linux
ENV CGO_ENABLED=1
RUN go build -a -o dataGenerator .

### Start a new stage from scratch
FROM alpine:latest
FROM golang:1.21.3

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/dataGenerator .

# Ensure the binary is executable
RUN chmod +x dataGenerator

# Command to run the binary when the container starts
CMD ["./dataGenerator"]
