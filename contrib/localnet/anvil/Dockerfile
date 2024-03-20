# Start from the latest Rust image as Anvil is built with Rust
FROM ghcr.io/foundry-rs/foundry:latest

# Expose the default Anvil port
EXPOSE 8545

# Run Anvil with specified chain ID and a prefunded account when the container starts
ENTRYPOINT ["anvil", "--host", "0.0.0.0", "--chain-id", "11155111"]