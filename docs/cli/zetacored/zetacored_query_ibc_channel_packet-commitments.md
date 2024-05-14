# query ibc channel packet-commitments

Query all packet commitments associated with a channel

### Synopsis

Query all packet commitments associated with a channel

```
zetacored query ibc channel packet-commitments [port-id] [channel-id] [flags]
```

### Examples

```
zetacored query ibc channel packet-commitments [port-id] [channel-id]
```

### Options

```
      --count-total        count total number of records in packet commitments associated with a channel to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for packet-commitments
      --limit uint         pagination limit of packet commitments associated with a channel to query for (default 100)
      --node string        [host]:[port] to Tendermint RPC interface for this chain 
      --offset uint        pagination offset of packet commitments associated with a channel to query for
  -o, --output string      Output format (text|json) 
      --page uint          pagination page of packet commitments associated with a channel to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of packet commitments associated with a channel to query for
      --reverse            results are sorted in descending order
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

