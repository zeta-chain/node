# ZetaTool

ZetaTool is a utility CLI for Zetachain.It currently provides a command to fetch the ballot/cctx identifier from the inbound hash

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

### Command
```shell
zetatool get-ballot [inboundHash] [chainID] --config <filename.json>
```
### Example
```shell
zetatool get-ballot 0x61008d7f79b2955a15e3cb95154a80e19c7385993fd0e083ff0cbe0b0f56cb9a 1
{"level":"info","time":"2025-01-20T11:30:47-05:00","message":"ballot identifier: 0xae189ab5cd884af784835297ac43eb55deb8a7800023534c580f44ee2b3eb5ed"}
```

- `inboundHash`: The inbound hash of the transaction for which the ballot identifier is to be fetched
- `chainID`: The chain ID of the chain to which the transaction belongs
- `config`: [Optional] The path to the configuration file. When not provided, the configuration in the file is user. A sample config is provided at `cmd/zetatool/config/sample_config.json`

The Config contains the rpcs needed for the tool to function,
if not provided the tool automatically uses the default rpcs.It is able to fetch the rpc needed using the chain ID

The command returns a ballot identifier for the given inbound hash.

