# ZetaTool

ZetaTool is a utility tool for the Zeta-Chain project. It provides a command to fetch the ballot/cctx identifier from the inbound hash

## Installation

To install ZetaTool, clone the repository and build the project:

```sh
git clone https://github.com/zeta-chain/node.git
cd node/cmd/zetatool
go build -o zetatool
```

Alternatively you can also use the target `make install-zetatool`

## Usage 

### Fetching the Ballot Identifier

```shell
get-ballot [inboundHash] [chainID] --config <filename.json>
```

- `inboundHash`: The inbound hash of the transaction for which the ballot identifier is to be fetched
- `chainID`: The chain ID of the chain to which the transaction belongs
- `--config`: [Optional]The path to the configuration file. When not provided, the default configuration is used 

The command returns a ballot identifier for the given inbound hash.

