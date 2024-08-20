## Full Deployment Guide for Zetacored and Bitcoin Nodes

This guide details deploying Zetacored nodes on both ZetaChain mainnet and Athens3 (testnet), alongside setting up a Bitcoin node for mainnet. The setup utilizes Docker Compose with environment variables for a streamlined deployment process.

Here's a comprehensive documentation using markdown tables to cover all the `make` commands for managing Zetacored and Bitcoin nodes, including where to modify the environment variables in Docker Compose configurations.

### Zetacored / BTC Node Deployment and Management

#### Commands Overview for Zetacored

| Environment                         | Action                      | Command                                                       | Docker Compose Location                  |
|-------------------------------------|-----------------------------|---------------------------------------------------------------|------------------------------------------|
| **Mainnet**                         | Start Bitcoin Node          | `make start-bitcoin-node-mainnet`                             | `contrib/rpc/bitcoind-mainnet`           |
| **Mainnet**                         | Stop Bitcoin Node           | `make stop-bitcoin-node-mainnet`                              | `contrib/rpc/bitcoind-mainnet`           |
| **Mainnet**                         | Clean Bitcoin Node Data     | `make clean-bitcoin-node-mainnet`                             | `contrib/rpc/bitcoind-mainnet`           |
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


### Bitcoin Node Setup for Mainnet

#### Commands Overview for Bitcoin

| Action | Command | Docker Compose Location |
|--------|---------|-------------------------|
| Start Node | `make start-mainnet-bitcoind-node DOCKER_TAG=36-mainnet` | `contrib/mainnet/bitcoind` |
| Stop Node | `make stop-mainnet-bitcoind-node` | `contrib/mainnet/bitcoind` |
| Clean Node Data | `make clean-mainnet-bitcoind-node` | `contrib/mainnet/bitcoind` |

### Configuration Options

#### Where to Modify Environment Variables

The environment variables for both Zetacored and Bitcoin nodes are defined in the `docker-compose.yml` files located in the respective directories mentioned above. These variables control various operational aspects like the sync type, networking details, and client behavior.

#### Example Environment Variables for Zetacored

| Variable | Description | Example |
|----------|-------------|---------|
| `DAEMON_HOME` | Daemon's home directory | `/root/.zetacored` |
| `NETWORK` | Network identifier | `mainnet`, `athens3` |
| `CHAIN_ID` | Chain ID for the network | `zetachain_7000-1`, `athens_7001-1` |
| `RESTORE_TYPE` | Node restoration method | `snapshot`, `statesync` |
| `SNAPSHOT_API` | API URL for fetching snapshots | `https://snapshots.rpc.zetachain.com` |

#### Example Environment Variables for Bitcoin

| Variable | Description | Example |
|----------|-------------|---------|
| `bitcoin_username` | Username for Bitcoin RPC | `user` |
| `bitcoin_password` | Password for Bitcoin RPC | `pass` |
| `WALLET_NAME` | Name of the Bitcoin wallet | `tssMainnet` |
| `WALLET_ADDRESS` | Bitcoin wallet address for transactions | `bc1qm24wp577nk8aacckv8np465z3dvmu7ry45el6y` |

This detailed tabulation ensures all necessary commands and configurations are neatly organized, providing clarity on where to manage the settings and how to execute different operations for Zetacored and Bitcoin nodes across different environments.
