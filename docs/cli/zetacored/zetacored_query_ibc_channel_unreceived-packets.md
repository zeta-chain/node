# query ibc channel unreceived-packets

Query all the unreceived packets associated with a channel

### Synopsis

Determine if a packet, given a list of packet commitment sequences, is unreceived.

The return value represents:
- Unreceived packet commitments: no acknowledgement exists on receiving chain for the given packet commitment sequence on sending chain.


```
zetacored query ibc channel unreceived-packets [port-id] [channel-id] [flags]
```

### Examples

```
zetacored query ibc channel unreceived-packets [port-id] [channel-id] --sequences=1,2,3
```

### Options

```
      --grpc-addr string       the gRPC endpoint to use for this chain
      --grpc-insecure          allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int             Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                   help for unreceived-packets
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

