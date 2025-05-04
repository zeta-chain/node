## Full Deployment Guide for Zetacored

This guide details deploying Zetacored nodes on both ZetaChain mainnet and Athens3 (testnet). The setup utilizes Docker Compose with environment variables for a streamlined deployment process.

Here's a comprehensive documentation using markdown tables to cover all the `make` commands for managing Zetacored, including where to modify the environment variables in Docker Compose configurations.

### Zetacored / BTC Node Deployment and Management

#### Commands Overview for Zetacored

| Environment                         | Action                      | Command                                                       | Docker Compose Location                  |
|-------------------------------------|-----------------------------|---------------------------------------------------------------|------------------------------------------|
| **Mainnet**                         | Start Ethereum Node         | `make start-eth-node-mainnet`                                 | `contrib/rpc/ethereum`                   |
| **Mainnet**                         | Stop Ethereum Node          | `make stop-eth-node-mainnet`                                  | `contrib/rpc/ethereum`                   |
| **Mainnet**                         | Clean Ethereum Node Data    | `make clean-eth-node-mainnet`                                 | `contrib/rpc/ethereum`                   |
| **Mainnet**                         | Start Zetacored Node        | `make start-mainnet-zetarpc-node DOCKER_TAG=ubuntu-v14.0.1`   | `contrib/rpc/zetacored`                  |
| **Mainnet**                         | Stop Zetacored Node         | `make stop-mainnet-zetarpc-node`                              | `contrib/rpc/zetacored`                  |
| **Mainnet**                         | Clean Zetacored Node Data   | `make clean-mainnet-zetarpc-node`                             | `contrib/rpc/zetacored`                  |
| **Testnet (Athens3)**               | Start Zetacored Node        | `make start-testnet-zetarpc-node DOCKER_TAG=ubuntu-v14.0.1`   | `contrib/rpc/zetacored`                  |
| **Testnet (Athens3)**               | Stop Zetacored Node         | `make stop-testnet-zetarpc-node`                              | `contrib/rpc/zetacored`                  |
| **Testnet (Athens3)**               | Clean Zetacored Node Data   | `make clean-testnet-zetarpc-node`                             | `contrib/rpc/zetacored`                  |
| **Mainnet Local Build**             | Start Zetacored Node        | `make start-zetacored-rpc-mainnet-localbuild`                 | `contrib/rpc/zetacored`                  |
| **Mainnet Local Build**             | Stop Zetacored Node         | `make stop-zetacored-rpc-mainnet-localbuild`                  | `contrib/rpc/zetacored`                  |
| **Mainnet Local Build**             | Clean Zetacored Node Data   | `make clean-zetacored-rpc-mainnet-localbuild`                 | `contrib/rpc/zetacored`                  |
| **Testnet Local Build (Athens3)**   | Start Zetacored Node        | `make start-zetacored-rpc-testnet-localbuild`                 | `contrib/rpc/zetacored`                  |
| **Testnet Local Build (Athens3)**   | Stop Zetacored Node         | `make stop-zetacored-rpc-testnet-localbuild`                  | `contrib/rpc/zetacored`                  |
| **Testnet Local Build (Athens3)**   | Clean Zetacored Node Data   | `make clean-zetacored-rpc-testnet-localbuild`                 | `contrib/rpc/zetacored`                  |
