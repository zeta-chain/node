## Usage 

### Track the status of a CCTX

### Command
```shell
zetatool track-cctx [inboundHash] [chainID] --config <filename.json>
```
### Example
```shell
zetatool track-cctx 0x61008d7f79b2955a15e3cb95154a80e19c7385993fd0e083ff0cbe0b0f56cb9a 1
{"level":"info","time":"2025-02-03T12:59:33-05:00","message":"CCTX Identifier: 0xae189ab5cd884af784835297ac43eb55deb8a7800023534c580f44ee2b3eb5ed Status: OutboundMined"}
```

- `inboundHash`: The inbound hash of the transaction for which the ballot identifier is to be fetched
- `chainID`: The chain ID of the chain to which the transaction belongs
- `config`: [Optional] The path to the configuration file. When not provided, the configuration in the file is user. A sample config is provided at `cmd/zetatool/config/sample_config.json`
- `debug`: [Optional] The debug flag is used to print additional debug information when set to true

The Config contains the rpcs needed for the tool to function,
if not provided the tool automatically uses the default rpcs.It is able to fetch the rpc needed using the chain ID

The command returns
- The CCTX identifier
- The status of the CCTX

