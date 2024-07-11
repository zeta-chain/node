# Stage 1: Set permissions using a base image that allows it
FROM alpine:latest as builder

# Install necessary tools
RUN apk add --no-cache bash

# Set the working directory
WORKDIR /root/.hermes

# Copy entrypoint script
COPY entrypoint.sh .

# Set permissions
RUN chmod +x entrypoint.sh

# Stage 2: Use the Hermes image
FROM informalsystems/hermes:1.10.0

# Set the working directory
WORKDIR /root/.hermes

# Copy the Hermes config file
COPY config.toml /root/.hermes/config.toml

# Copy the entrypoint script from the builder stage
COPY --from=builder /root/.hermes/entrypoint.sh /root/.hermes/entrypoint.sh

# Set the entrypoint
ENTRYPOINT ["/root/.hermes/entrypoint.sh"]