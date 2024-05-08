# query bank send-enabled

Query for send enabled entries

### Synopsis

Query for send enabled entries that have been specifically set.

To look up one or more specific denoms, supply them as arguments to this command.
To look up all denoms, do not provide any arguments.

```
zetacored query bank send-enabled [denom1 ...] [flags]
```

### Examples

```
Getting one specific entry:
  $ zetacored query bank send-enabled foocoin

Getting two specific entries:
  $ zetacored query bank send-enabled foocoin barcoin

Getting all entries:
  $ zetacored query bank send-enabled
```

### Options

```
      --count-total        count total number of records in send enabled entries to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for send-enabled
      --limit uint         pagination limit of send enabled entries to query for (default 100)
      --node string        [host]:[port] to Tendermint RPC interface for this chain 
      --offset uint        pagination offset of send enabled entries to query for
  -o, --output string      Output format (text|json) 
      --page uint          pagination page of send enabled entries to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of send enabled entries to query for
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

* [zetacored query bank](zetacored_query_bank.md)	 - Querying commands for the bank module

