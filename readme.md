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


Here is the formatted documentation in Markdown:

---

### Starting Full Zetacored Nodes

#### Step 1: Choose the Network

To start a node, use the `make` command with the `DOCKER_TAG` of the image you wish to use from Docker Hub.

- **For Mainnet:**

  ```shell
  # Use this command to start a mainnet node with a specific Docker tag
  make mainnet-zetarpc-node DOCKER_TAG={THE_DOCKER_TAG_FROM_DOCKER_HUB_YOU_WANT_TO_USE}
  # Example:
  make mainnet-zetarpc-node DOCKER_TAG=ubuntu-v12.3.0-docker-test
  ```

- **For Athens3:**

  ```shell
  # The command is the same for Athens3, just ensure you're specifying the correct Docker tag
  make mainnet-zetarpc-node DOCKER_TAG={THE_DOCKER_TAG_FROM_DOCKER_HUB_YOU_WANT_TO_USE}
  # Example:
  make mainnet-zetarpc-node DOCKER_TAG=ubuntu-v12.3.0-docker-test
  ```

**Note:** The default configuration is to restore from state sync. This process will download the necessary configurations and information from [Zeta-Chain Network Config](https://github.com/zeta-chain/network-config) and configure the node for state sync restore.

#### Changing the Sync Type

If you wish to change the sync type, you will need to modify the `docker-compose.yml` file located in `contrib/{NETWORK}/zetacored/`.

Change the following values according to your needs:

```yaml
# Possible values for RESTORE_TYPE are "snapshot", "snapshot-archive", or "statesync"
RESTORE_TYPE: "statesync"
MONIKER: "local-test"
RE_DO_START_SEQUENCE: "false"
```

To perform a snapshot restore from the latest snapshot, simply change the `RESTORE_TYPE` to either `snapshot` or `snapshot-archive`.

---

Here's the formatted documentation in Markdown for starting a full Bitcoind Mainnet node:

---

### Starting Full Bitcoind Mainnet Node

#### Step 1: Restore a Mainnet BTC Watcher Node

To restore a mainnet BTC watcher node from a BTC snapshot, run the following `make` command and specify the `DOCKER_TAG` with the image you want to use from Docker Hub.

```commandline
make mainnet-bitcoind-node DOCKER_TAG={DOCKER_TAG_FROM_DOCKER_HUB_TO_USE}
# Example:
make mainnet-bitcoind-node DOCKER_TAG=36-mainnet
```

#### Updating the TSS Address

If you need to update the TSS (Threshold Signature Scheme) address being watched, please edit the `docker-compose.yml` file located at `contrib/mainnet/bitcoind/docker-compose.yml`.

To update, simply change the user and password you wish to use, and the TSS address to watch. Then, run the command provided above to apply your changes.

---