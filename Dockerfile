# Use a minimal Alpine Linux image
FROM debian:latest

# Set the working directory inside the container
WORKDIR /

# Copy the current directory contents into the container at /app
COPY db /db
COPY BalkanLinGO /
COPY views /views
COPY .env /

# Expose any necessary ports (if your Go binary listens on a specific port)
EXPOSE 3000

# Define the command to run your Go binary
#CMD ["ls", "|", "cat"]
RUN ls /
#CMD /app/BalkanLinGO
CMD ["./BalkanLinGO"]