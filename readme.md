# ZetaChain

ZetaChain is an EVM-compatible L1 blockchain that enables omnichain, generic
smart contracts and messaging between any blockchain.

## Prerequisites

- [Go](https://golang.org/doc/install) 1.23
- [Docker](https://docs.docker.com/install/) and
  [Docker Compose](https://docs.docker.com/compose/install/) (optional, for
  running tests locally)
- [buf](https://buf.build/) (optional, for processing protocol buffer files)
- [jq](https://stedolan.github.io/jq/download/) (optional, for running scripts)

## Components of ZetaChain

ZetaChain is built with [Cosmos SDK](https://github.com/cosmos/cosmos-sdk), a
modular framework for building blockchain and
[Cosmos EVM](https://github.com/zeta-chain/evm), a module that implements
EVM-compatibility.

This repository contains the core components:

* [Blockchain Modules (ZetaCore)](x):
This section contains the core logic of the ZetaChain blockchain, built using Cosmos SDK modules. These modules are responsible for managing the state, state transitions, and overall functionality of the ZetaChain network.
* [ZetaClient](zetaclient):
The ZetaClient is a specialized client designed to act as an observer and signer for the ZetaChain network. It is responsible for communicating with the blockchain, relaying messages, and performing signature tasks to ensure the network operates cross-chain transactions.

### Protocol Contracts

In addition to the blockchain codebase, ZetaChain’s architecture includes a set of protocol contracts that serve as an interface for developers to interact with the blockchain. These smart contracts are deployed across various blockchain networks. The smart contract source code is maintained in separate repositories, depending on the network they are deployed on:

* [ZetaChain EVM and EVM connected chains](https://github.com/zeta-chain/protocol-contracts)
* [Solana connected chains (SVM)](https://github.com/zeta-chain/protocol-contracts-solana)
* [TON connected chains (TVM)](https://github.com/zeta-chain/protocol-contracts-ton)
* [Sui connected chains (Sui's MVM)](https://github.com/zeta-chain/protocol-contracts-sui)

These repositories contain the necessary code and tools to deploy, interact with, and extend the functionality of ZetaChain’s cross-chain protocol on each respective blockchain network.

### Versions

For a complete compatibility matrix showing which protocol contract versions are compatible with specific ZetaCore and ZetaClient versions, see [VERSIONS.md](VERSIONS.md).

## Building the `zetacored`/`zetaclientd` binaries

Clone this repository, checkout the latest release tag, and type the following command to build the binaries:

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
command to run generated files generation (ProtoBuf, OpenAPI and docs):

```
make generate
```

This command will use `buf` to generate the Go files from the protocol buffer
files and move them to the correct directories inside `x/`. It will also
generate an OpenAPI spec.

This command will run a script to update the modules' documentation. The script
uses static code analysis to read the protocol buffer files and identify all
Cosmos SDK messages. It then searches the source code for the corresponding
message handler functions and retrieves the documentation for those functions.
Finally, it creates a `messages.md` file for each module, which contains the
documentation for all the messages in that module.

## Further Reading

Find below further documentation for development and running your own ZetaChain node:

- [Get familiar with our release lifecycle](docs/development/RELEASE_LIFECYCLE.md)
- [Run the E2E tests and interact with the localnet](docs/development/LOCAL_TESTING.md)
- [Make a new ZetaChain release](docs/development/RELEASES.md)
- [Deploy your own ZetaChain or Bitcoin node](docs/development/DEPLOY_NODES.md)
- [Run the simulation tests](docs/development/SIMULATION_TESTING.md)

## Community

[X (formerly Twitter)](https://x.com/zetablockchain) |
[Discord](https://discord.com/invite/zetachain) |
[Telegram](https://t.me/zetachainofficial) | [Website](https://zetachain.com)
