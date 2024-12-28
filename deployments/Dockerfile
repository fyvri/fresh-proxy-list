FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy compiled binary into the container
COPY fresh-proxy-list /app/fresh-proxy-list

# Set default command to run application
ENTRYPOINT ["/app/fresh-proxy-list"]
