# query ibc channel unreceived-acks

Query all the unreceived acks associated with a channel

### Synopsis

Given a list of acknowledgement sequences from counterparty, determine if an ack on the counterparty chain has been received on the executing chain.

The return value represents:
- Unreceived packet acknowledgement: packet commitment exists on original sending (executing) chain and ack exists on receiving chain.


```
zetacored query ibc channel unreceived-acks [port-id] [channel-id] [flags]
```

### Examples

```
zetacored query ibc channel unreceived-acks [port-id] [channel-id] --sequences=1,2,3
```

### Options

```
      --grpc-addr string       the gRPC endpoint to use for this chain
      --grpc-insecure          allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int             Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                   help for unreceived-acks
      --node string            [host]:[port] to Tendermint RPC interface for this chain 
  -o, --output string          Output format (text|json) 
      --sequences int64Slice   comma separated list of packet sequence numbers (default [])
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic) 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query ibc channel](zetacored_query_ibc_channel.md)	 - IBC channel query subcommands

