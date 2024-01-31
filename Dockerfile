# Use a minimal Alpine Linux image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /usr/src/app

# Copy the Go binary into the container
COPY * /app/

# Expose any necessary ports (if your Go binary listens on a specific port)
EXPOSE 3000

# Define the command to run your Go binary
CMD ["GOGC=4000","./BalkanLinGO_arm64"]
