# ZetaChain

ZetaChain is an EVM-compatible L1 blockchain that enables omnichain, generic
smart contracts and messaging between any blockchain.

## Prerequisites

- [Go](https://golang.org/doc/install) 1.19
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
For athens3 testnet, clone this repository and checkout the latest release tag, and type
```
make install-testnet
```
to build. 

This command will install the `zetacoded` and `zetaclientd` binaries in your
`$GOPATH/bin` directory.

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
Cosmos SDK messages. It then search the source code for the corresponding
message handler functions and retrieves the documentation for those functions.
Finally, it creates a `messages.md` file for each module, which contains the
documentation for all the messages in that module.

## Running tests

To check that the source code is working as expected, refer to the manual on how
to [run the smoke test](./contrib/localnet/README.md).

## Community

[Twitter](https://twitter.com/zetablockchain) |
[Discord](https://discord.com/invite/zetachain) |
[Telegram](https://t.me/zetachainofficial) | [Website](https://zetachain.com)
