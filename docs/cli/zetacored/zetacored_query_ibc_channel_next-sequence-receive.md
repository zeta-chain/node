# query ibc channel next-sequence-receive

Query a next receive sequence

### Synopsis

Query the next receive sequence for a given channel

```
zetacored query ibc channel next-sequence-receive [port-id] [channel-id] [flags]
```

### Examples

```
zetacored query ibc channel next-sequence-receive [port-id] [channel-id]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for next-sequence-receive
      --node string        [host]:[port] to Tendermint RPC interface for this chain 
  -o, --output string      Output format (text|json) 
      --prove              show proofs for the query results (default true)
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

