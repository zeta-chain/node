# This Dockerfile is used to build a Docker image for Anvil, a localnet for testing purposes.
# Currently we directly set the chain ID to 11155111 and expose the default Anvil port specifically for the chain migration test.

# Start from the latest Rust image as Anvil is built with Rust
FROM ghcr.io/foundry-rs/foundry:latest

# Expose the default Anvil port
EXPOSE 8545

# Run Anvil with specified chain ID and a prefunded account when the container starts
ENTRYPOINT ["anvil", "--host", "0.0.0.0", "--chain-id", "11155111"]