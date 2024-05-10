# query ibc client consensus-states

Query all the consensus states of a client.

### Synopsis

Query all the consensus states from a given client state.

```
zetacored query ibc client consensus-states [client-id] [flags]
```

### Examples

```
zetacored query ibc client consensus-states [client-id]
```

### Options

```
      --count-total        count total number of records in consensus states to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for consensus-states
      --limit uint         pagination limit of consensus states to query for (default 100)
      --node string        [host]:[port] to Tendermint RPC interface for this chain 
      --offset uint        pagination offset of consensus states to query for
  -o, --output string      Output format (text|json) 
      --page uint          pagination page of consensus states to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of consensus states to query for
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

* [zetacored query ibc client](zetacored_query_ibc_client.md)	 - IBC client query subcommands

