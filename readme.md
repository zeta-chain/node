# ZetaChain

ZetaChain is an EVM-compatible L1 blockchain that enables omnichain, generic
smart contracts and messaging between any blockchain.

## Prerequisites

- [Go](https://golang.org/doc/install) 1.20
- [Docker](https://docs.docker.com/install/) and
  [Docker Compose](https://docs.docker.com/compose/install/) (optional, for
  running tests locally)
- [buf](https://buf.build/) (optional, for processing protocol buffer files)
- [jq](https://stedolan.github.io/jq/download/) (optional, for running scripts)

## Components of ZetaChain

ZetaChain is built with [Cosmos SDK](https://github.com/cosmos/cosmos-sdk), a
modular framework for building blockchain and
[Ethermint](https://github.com/evmos/ethermint), a module that implements
EVM-compatibility.

- [zeta-node](https://github.com/zeta-chain/zeta-node) (this repository)
  contains the source code for the ZetaChain node (`zetacored`) and the
  ZetaChain client (`zetaclientd`).
- [protocol-contracts](https://github.com/zeta-chain/protocol-contracts)
  contains the source code for the Solidity smart contracts that implement the
  core functionality of ZetaChain.

## Building the zetacored/zetaclientd binaries
For the Athens 3 testnet, clone this repository, checkout the latest release tag, and type the following command to build the binaries:
```
make install
```
to build. 

This command will install the `zetacored` and `zetaclientd` binaries in your
`$GOPATH/bin` directory.

Verify that the version of the binaries match the release tag.  
```
zetacored version
zetaclientd version
```

## Making changes to the source code

After making changes to any of the protocol buffer files, run the following
command to generate the Go files:

```
make proto
```

This command will use `buf` to generate the Go files from the protocol buffer
files and move them to the correct directories inside `x/`. It will also
generate an OpenAPI spec.

### Generate documentation

To generate the documentation, run the following command:

```
make specs
```

This command will run a script to update the modules' documentation. The script
uses static code analysis to read the protocol buffer files and identify all
Cosmos SDK messages. It then searches the source code for the corresponding
message handler functions and retrieves the documentation for those functions.
Finally, it creates a `messages.md` file for each module, which contains the
documentation for all the messages in that module.

## Running tests

To check that the source code is working as expected, refer to the manual on how
to [run the E2E test](./LOCAL_TESTING.md).

## Community

[Twitter](https://twitter.com/zetablockchain) |
[Discord](https://discord.com/invite/zetachain) |
[Telegram](https://t.me/zetachainofficial) | [Website](https://zetachain.com)


## Creating a Release Candidate
Creating a release candidate for testing is a straightforward process. Here are the steps to follow:

### Steps
 - Step 1. Create the release candidate tag with the following format (e.g., vx.x.x-rc) ex. v11.0.0-rc.
 - Step 2. Once a RC branch is created the automation will kickoff to build and upload the release and its binaries.

By following these steps, you can efficiently create a release candidate for QA and validation. In the future we will make this automatically deploy to a testnet when a -rc branch is created. 
Currently, raising the proposal to deploy to testnet is a manual process via GovOps repo. 

## Creating a Release / Hotfix Release

To create a release simply execute the publish-release workflow and follow the steps below.

### Steps
 - Step 1. Go to this pipeline: https://github.com/zeta-chain/node/actions/workflows/publish-release.yml
 - Step 2. Select the dropdown branch / tag you want to create the release with.
 - Step 3. In the version input, include the version of your release. Note. The major version must match what is in the upgrade handler.
 - Step 4. Select if you want to skip the tests by checking the checkbox for skip tests.
 - Step 5. Once the testing steps pass it will create a Github Issue. This Github Issue needes to be approved by one of the approvers: kingpinXD,lumtis,brewmaster012

Once the release is approved the pipeline will continue and will publish the releases with the title / version you specified in the user input.

---
## Full Deployment Guide for Zetacored and Bitcoin Nodes

This guide details deploying Zetacored nodes on both ZetaChain mainnet and Athens3 (testnet), alongside setting up a Bitcoin node for mainnet. The setup utilizes Docker Compose with environment variables for a streamlined deployment process.

### Deploying Zetacored Nodes

#### Launching a Node

**For Mainnet:**
- Use the `make` command with a specified Docker tag to initiate a mainnet Zetacored node.
  ```shell
  make mainnet-zetarpc-node DOCKER_TAG=ubuntu-v14.0.1
  ```

**For Athens3 (Testnet):**
- Similar command structure for Athens3, ensuring the correct Docker tag is used.
  ```shell
  make testnet-zetarpc-node DOCKER_TAG=ubuntu-v14.0.1
  ```
  
#### Modifying the Sync Type

**To change the sync type for your node:**
- Edit docker-compose.yml in contrib/{NETWORK}/zetacored/.
- Set RESTORE_TYPE to your desired method (snapshot, snapshot-archive, statesync).

#### Zetacored Environment Variables

| Variable | Description |
|----------|-------------|
| `DAEMON_HOME` | Daemon's home directory (`/root/.zetacored`). |
| `NETWORK` | Network identifier: `mainnet` or `athens3` (for testnet). |
| `RESTORE_TYPE` | Node restoration method: `snapshot`, `snapshot-archive`, `statesync`. |
| `SNAPSHOT_API` | API URL for fetching snapshots. |
| `TRUST_HEIGHT_DIFFERENCE_STATE_SYNC` | Trust height difference for state synchronization. |
| `CHAIN_ID` | Chain ID for the network. |
| `VISOR_NAME` | Visor software name, typically `cosmovisor`. |
| `DAEMON_NAME` | Daemon software name, `zetacored`. |
| `DAEMON_ALLOW_DOWNLOAD_BINARIES` | Enable daemon to download binaries. |
| `DAEMON_RESTART_AFTER_UPGRADE` | Restart daemon after software upgrade. |
| `UNSAFE_SKIP_BACKUP` | Skip backup during potentially unsafe operations. |
| `CLIENT_DAEMON_NAME` | Client daemon name, such as `zetaclientd`. |
| `CLIENT_DAEMON_ARGS` | Extra arguments for the client daemon. |
| `CLIENT_SKIP_UPGRADE` | Skip client software upgrade. |
| `CLIENT_START_PROCESS` | Begin client process start-up. |
| `MONIKER` | Node's moniker or nickname. |
| `RE_DO_START_SEQUENCE` | Restart node setup from scratch if necessary. |

### Bitcoin Node Setup for Mainnet

**Restoring a BTC Watcher Node:**
- To deploy a Bitcoin mainnet node, specify the `DOCKER_TAG` for your Docker image.
  ```shell
  make mainnet-bitcoind-node DOCKER_TAG=36-mainnet
  ```

#### Bitcoin Node Environment Variables

| Variable | Description |
|----------|-------------|
| `bitcoin_username` | Username for Bitcoin RPC. |
| `bitcoin_password` | Password for Bitcoin RPC. |
| `NETWORK_HEIGHT_URL` | URL to fetch the latest block height. |
| `WALLET_NAME` | Name of the Bitcoin wallet. |
| `WALLET_ADDRESS` | Bitcoin wallet address for transactions. |
| `SNAPSHOT_URL` | URL for downloading the blockchain snapshot. |
| `SNAPSHOT_RESTORE` | Enable restoration from snapshot. |
| `CLEAN_SNAPSHOT` | Clean existing data before restoring snapshot. |
| `DOWNLOAD_SNAPSHOT` | Download the snapshot if not present. |

### Docker Compose Configurations

#### Zetacored Mainnet

```yaml
version: '3.8'
services:
  zetachain_mainnet_rpc:
    platform: linux/amd64
    image: zetachain/zetacored:${DOCKER_TAG:-ubuntu-v14.0.1}
    environment:
      DAEMON_HOME: "/root/.zetacored"
      NETWORK: mainnet
      RESTORE_TYPE: "snapshot"
      SNAPSHOT_API: https://snapshots.zetachain.com
      TRUST_HEIGHT_DIFFERENCE_STATE_SYNC: 40000
      CHAIN_ID: "zetachain_7000-1"
      VISOR_NAME: "cosmovisor"
      DAEMON_NAME: "zetacored"
      DAEMON_ALLOW_DOWNLOAD_BINARIES: "false"
      DAEMON_RESTART_AFTER_UPGRADE: "true"
      UNSAFE_SKIP_BACKUP: "true"
      CLIENT_DAEMON_NAME: "zetaclientd"
      CLIENT_DAEMON_ARGS: ""
      CLIENT_SKIP_UPGRADE: "true"
      CLIENT_START_PROCESS: "false"
      MONIKER: local-test
      RE_DO_START_SEQUENCE: "false"
    ports:
      - "26656:26656"
      - "1317:1317"
      - "8545:8545"
      - "8546:8546"
      - "26657:26657"
      - "9090:9090"
      - "9091:9091"
    volumes:


      - zetacored_data_mainnet:/root/.zetacored/
    entrypoint: bash /scripts/start.sh
volumes:
  zetacored_data_mainnet:
```

#### Zetacored Athens3/Testnet

```yaml
version: '3.8'
services:
  zetachain_testnet_rpc:
    platform: linux/amd64
    image: zetachain/zetacored:${DOCKER_TAG:-ubuntu-v14-testnet}
    environment:
      DAEMON_HOME: "/root/.zetacored"
      NETWORK: athens3
      RESTORE_TYPE: "snapshot"
      SNAPSHOT_API: https://snapshots.zetachain.com
      TRUST_HEIGHT_DIFFERENCE_STATE_SYNC: 40000
      CHAIN_ID: "athens_7001-1"
      VISOR_NAME: "cosmovisor"
      DAEMON_NAME: "zetacored"
      DAEMON_ALLOW_DOWNLOAD_BINARIES: "false"
      DAEMON_RESTART_AFTER_UPGRADE: "true"
      UNSAFE_SKIP_BACKUP: "true"
      CLIENT_DAEMON_NAME: "zetaclientd"
      CLIENT_DAEMON_ARGS: ""
      CLIENT_SKIP_UPGRADE: "true"
      CLIENT_START_PROCESS: "false"
      MONIKER: local-test
      RE_DO_START_SEQUENCE: "false"
    ports:
      - "26656:26656"
      - "1317:1317"
      - "8545:8545"
      - "8546:8546"
      - "26657:26657"
      - "9090:9090"
      - "9091:9091"
    volumes:
      - zetacored_data_athens3:/root/.zetacored/
    entrypoint: bash /scripts/start.sh
volumes:
  zetacored_data_athens3:
```

#### Bitcoin Mainnet Node

```yaml
version: '3'
services:
  bitcoin:
    image: zetachain/bitcoin:${DOCKER_TAG:-36-mainnet}
    platform: linux/amd64
    environment:
      - bitcoin_username=test
      - bitcoin_password=test
      - NETWORK_HEIGHT_URL=https://blockstream.info/api/blocks/tip/height
      - WALLET_NAME=tssMainnet
      - WALLET_ADDRESS=bc1qm24wp577nk8aacckv8np465z3dvmu7ry45el6y
      - SNAPSHOT_URL=https://storage.googleapis.com/bitcoin-rpc-snapshots-prod/bitcoind-mainnet-2024-02-20-00-22-06.tar.gz
      - SNAPSHOT_RESTORE=true
      - CLEAN_SNAPSHOT=true
      - DOWNLOAD_SNAPSHOT=true
    volumes:
      - bitcoin_data:/root/
    ports:
      - 8332:8332
volumes:
  bitcoin_data:
```

Replace placeholders in Docker Compose files and `make` commands with actual values appropriate for your deployment scenario. This complete setup guide is designed to facilitate the deployment and management of Zetacored and Bitcoin nodes in various environments.