# query ibc-transfer total-escrow

Query the total amount of tokens in escrow for a denom

### Synopsis

Query the total amount of tokens in escrow for a denom

```
zetacored query ibc-transfer total-escrow [denom] [flags]
```

### Examples

```
zetacored query ibc-transfer total-escrow uosmo
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for total-escrow
      --node string        [host]:[port] to Tendermint RPC interface for this chain 
  -o, --output string      Output format (text|json) 
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

* [zetacored query ibc-transfer](zetacored_query_ibc-transfer.md)	 - IBC fungible token transfer query subcommands

