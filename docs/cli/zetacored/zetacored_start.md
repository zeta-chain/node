# start

Run the full node

### Synopsis

Run the full node application with Tendermint in or out of process. By
default, the application will run with Tendermint in process.

Pruning options can be provided via the '--pruning' flag or alternatively with '--pruning-keep-recent',
'pruning-keep-every', and 'pruning-interval' together.

For '--pruning' the options are as follows:

default: the last 100 states are kept in addition to every 500th state; pruning at 10 block intervals
nothing: all historic states will be saved, nothing will be deleted (i.e. archiving node)
everything: all saved states will be deleted, storing only the current state; pruning at 10 block intervals
custom: allow pruning options to be manually specified through 'pruning-keep-recent', 'pruning-keep-every', and 'pruning-interval'

Node halting configurations exist in the form of two flags: '--halt-height' and '--halt-time'. During
the ABCI Commit phase, the node will check if the current block height is greater than or equal to
the halt-height or if the current block time is greater than or equal to the halt-time. If so, the
node will attempt to gracefully shutdown and the block will not be committed. In addition, the node
will not be able to commit subsequent blocks.

For profiling and benchmarking purposes, CPU profiling can be enabled via the '--cpu-profile' flag
which accepts a path for the resulting pprof file.


```
zetacored start [flags]
```

### Options

```
      --abci string                                     specify abci transport (socket | grpc) 
      --address string                                  Listen address 
      --api.enable                                      Defines if Cosmos-sdk REST server should be enabled
      --api.enabled-unsafe-cors                         Defines if CORS should be enabled (unsafe - use it at your own risk)
      --app-db-backend string                           The type of database for application and snapshots databases
      --consensus.create_empty_blocks                   set this to false to only produce blocks when there are txs or when the AppHash changes (default true)
      --consensus.create_empty_blocks_interval string   the possible interval between empty blocks 
      --consensus.double_sign_check_height int          how many blocks to look back to check existence of the node's consensus votes before joining consensus
      --cpu-profile string                              Enable CPU profiling and write to the provided file
      --db_backend string                               database backend: goleveldb | cleveldb | boltdb | rocksdb | badgerdb 
      --db_dir string                                   database directory 
      --evm.max-tx-gas-wanted uint                      the gas wanted for each eth tx returned in ante handler in check tx mode
      --evm.tracer string                               the EVM tracer type to collect execution traces from the EVM transaction execution (json|struct|access_list|markdown)
      --fast_sync                                       fast blockchain syncing (default true)
      --genesis_hash bytesHex                           optional SHA-256 hash of the genesis file
      --grpc-only                                       Start the node in gRPC query only mode without Tendermint process
      --grpc-web.address string                         The gRPC-Web server address to listen on 
      --grpc-web.enable                                 Define if the gRPC-Web server should be enabled. (Note: gRPC must also be enabled.) (default true)
      --grpc.address string                             the gRPC server address to listen on 
      --grpc.enable                                     Define if the gRPC server should be enabled (default true)
      --halt-height uint                                Block height at which to gracefully halt the chain and shutdown the node
      --halt-time uint                                  Minimum block time (in Unix seconds) at which to gracefully halt the chain and shutdown the node
  -h, --help                                            help for start
      --home string                                     The application home directory 
      --inter-block-cache                               Enable inter-block caching (default true)
      --inv-check-period uint                           Assert registered invariants every N blocks
      --json-rpc.address string                         the JSON-RPC server address to listen on 
      --json-rpc.allow-unprotected-txs                  Allow for unprotected (non EIP155 signed) transactions to be submitted via the node's RPC when the global parameter is disabled
      --json-rpc.api strings                            Defines a list of JSON-RPC namespaces that should be enabled (default [eth,net,web3])
      --json-rpc.block-range-cap eth_getLogs            Sets the max block range allowed for eth_getLogs query (default 10000)
      --json-rpc.enable                                 Define if the JSON-RPC server should be enabled (default true)
      --json-rpc.enable-indexer                         Enable the custom tx indexer for json-rpc
      --json-rpc.evm-timeout duration                   Sets a timeout used for eth_call (0=infinite) (default 5s)
      --json-rpc.filter-cap int32                       Sets the global cap for total number of filters that can be created (default 200)
      --json-rpc.gas-cap uint                           Sets a cap on gas that can be used in eth_call/estimateGas unit is aphoton (0=infinite) (default 25000000)
      --json-rpc.http-idle-timeout duration             Sets a idle timeout for json-rpc http server (0=infinite) (default 2m0s)
      --json-rpc.http-timeout duration                  Sets a read/write timeout for json-rpc http server (0=infinite) (default 30s)
      --json-rpc.logs-cap eth_getLogs                   Sets the max number of results can be returned from single eth_getLogs query (default 10000)
      --json-rpc.max-open-connections int               Sets the maximum number of simultaneous connections for the server listener
      --json-rpc.txfee-cap float                        Sets a cap on transaction fee that can be sent via the RPC APIs (1 = default 1 photon) (default 1)
      --json-rpc.ws-address string                      the JSON-RPC WS server address to listen on 
      --metrics                                         Define if EVM rpc metrics server should be enabled
      --min-retain-blocks uint                          Minimum block height offset during ABCI commit to prune Tendermint blocks
      --minimum-gas-prices string                       Minimum gas prices to accept for transactions; Any fee in a tx must meet this minimum (e.g. 0.01photon;0.0001stake)
      --moniker string                                  node name 
      --p2p.external-address string                     ip:port address to advertise to peers for them to dial
      --p2p.laddr string                                node listen address. (0.0.0.0:0 means any interface, any port) 
      --p2p.persistent_peers string                     comma-delimited ID@host:port persistent peers
      --p2p.pex                                         enable/disable Peer-Exchange (default true)
      --p2p.private_peer_ids string                     comma-delimited private peer IDs
      --p2p.seed_mode                                   enable/disable seed mode
      --p2p.seeds string                                comma-delimited ID@host:port seed nodes
      --p2p.unconditional_peer_ids string               comma-delimited IDs of unconditional peers
      --p2p.upnp                                        enable/disable UPNP port forwarding
      --priv_validator_laddr string                     socket address to listen on for connections from external priv_validator process
      --proxy_app string                                proxy app address, or one of: 'kvstore', 'persistent_kvstore', 'counter', 'e2e' or 'noop' for local testing. 
      --pruning string                                  Pruning strategy (default|nothing|everything|custom) 
      --pruning-interval uint                           Height interval at which pruned heights are removed from disk (ignored if pruning is not 'custom')
      --pruning-keep-recent uint                        Number of recent heights to keep on disk (ignored if pruning is not 'custom')
      --rpc.grpc_laddr string                           GRPC listen address (BroadcastTx only). Port required
      --rpc.laddr string                                RPC listen address. Port required 
      --rpc.pprof_laddr string                          pprof listen address (https://golang.org/pkg/net/http/pprof)
      --rpc.unsafe                                      enabled unsafe rpc methods
      --state-sync.snapshot-interval uint               State sync snapshot interval
      --state-sync.snapshot-keep-recent uint32          State sync snapshot to keep (default 2)
      --tls.certificate-path string                     the cert.pem file path for the server TLS configuration
      --tls.key-path string                             the key.pem file path for the server TLS configuration
      --trace                                           Provide full stack traces for errors in ABCI Log
      --trace-store string                              Enable KVStore tracing to an output file
      --transport string                                Transport protocol: socket, grpc 
      --unsafe-skip-upgrades ints                       Skip a set of upgrade heights to continue the old binary
      --with-tendermint                                 Run abci app embedded in-process with tendermint (default true)
      --x-crisis-skip-assert-invariants                 Skip x/crisis invariants check on startup
```

### Options inherited from parent commands

```
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic) 
```

### SEE ALSO

* [zetacored](zetacored.md)	 - Zetacore Daemon (server)

###### Auto generated by spf13/cobra on 22-Aug-2023
