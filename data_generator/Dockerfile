# Start with the golang base image
FROM golang:1.21.3

# Set the working directory
WORKDIR /app

# Copy the datagenerator binary into the container
COPY datagenerator .

# Ensure the binary is executable
RUN chmod +x datagenerator

# Command to run the binary when the container starts
CMD ["./datagenerator"]
