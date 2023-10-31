# query staking delegation

Query a delegation based on address and validator address

### Synopsis

Query delegations for an individual delegator on an individual validator.

Example:
$ zetacored query staking delegation zeta1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p zetavaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj

```
zetacored query staking delegation [delegator-addr] [validator-addr] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not TLS the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for delegation
      --node string        [host]:[port] to Tendermint RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic) 
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query staking](zetacored_query_staking.md)	 - Querying commands for the staking module

