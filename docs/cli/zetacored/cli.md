## zetacored

Zetacore Daemon (server)

### Options

```
  -h, --help                help for zetacored
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored add-genesis-account](#zetacored-add-genesis-account)	 - Add a genesis account to genesis.json
* [zetacored add-observer-list](#zetacored-add-observer-list)	 - Add a list of observers to the observer mapper ,default path is ~/.zetacored/os_info/observer_info.json
* [zetacored addr-conversion](#zetacored-addr-conversion)	 - convert a zeta1xxx address to validator operator address zetavaloper1xxx
* [zetacored collect-gentxs](#zetacored-collect-gentxs)	 - Collect genesis txs and output a genesis.json file
* [zetacored collect-observer-info](#zetacored-collect-observer-info)	 - collect observer info into the genesis from a folder , default path is ~/.zetacored/os_info/ 

* [zetacored config](#zetacored-config)	 - Utilities for managing application configuration
* [zetacored debug](#zetacored-debug)	 - Tool for helping with debugging your application
* [zetacored docs](#zetacored-docs)	 - Generate markdown documentation for zetacored
* [zetacored export](#zetacored-export)	 - Export state to JSON
* [zetacored gentx](#zetacored-gentx)	 - Generate a genesis tx carrying a self delegation
* [zetacored get-pubkey](#zetacored-get-pubkey)	 - Get the node account public key
* [zetacored index-eth-tx](#zetacored-index-eth-tx)	 - Index historical eth txs
* [zetacored init](#zetacored-init)	 - Initialize private validator, p2p, genesis, and application configuration files
* [zetacored keys](#zetacored-keys)	 - Manage your application's keys
* [zetacored parse-genesis-file](#zetacored-parse-genesis-file)	 - Parse the provided genesis file and import the required data into the optionally provided genesis file
* [zetacored query](#zetacored-query)	 - Querying subcommands
* [zetacored rollback](#zetacored-rollback)	 - rollback Cosmos SDK and CometBFT state by one height
* [zetacored snapshots](#zetacored-snapshots)	 - Manage local snapshots
* [zetacored start](#zetacored-start)	 - Run the full node
* [zetacored tendermint](#zetacored-tendermint)	 - Tendermint subcommands
* [zetacored testnet](#zetacored-testnet)	 - subcommands for starting or configuring local testnets
* [zetacored tx](#zetacored-tx)	 - Transactions subcommands
* [zetacored upgrade-handler-version](#zetacored-upgrade-handler-version)	 - Print the default upgrade handler version
* [zetacored validate](#zetacored-validate)	 - Validates the genesis file at the default location or at the location passed as an arg
* [zetacored version](#zetacored-version)	 - Print the application binary version information

## zetacored add-genesis-account

Add a genesis account to genesis.json

### Synopsis

Add a genesis account to genesis.json. The provided account must specify
the account address or key name and a list of initial coins. If a key name is given,
the address will be looked up in the local Keybase. The list of initial tokens must
contain valid denominations. Accounts may optionally be supplied with vesting parameters.


```
zetacored add-genesis-account [address_or_key_name] [coin][,[coin]] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for add-genesis-account
      --home string              The application home directory 
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test) 
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --vesting-amount string    amount of coins for vesting accounts
      --vesting-end-time int     schedule end time (unix epoch) for vesting accounts
      --vesting-start-time int   schedule start time (unix epoch) for vesting accounts
```

### Options inherited from parent commands

```
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)

## zetacored add-observer-list

Add a list of observers to the observer mapper ,default path is ~/.zetacored/os_info/observer_info.json

```
zetacored add-observer-list [observer-list.json]  [flags]
```

### Options

```
  -h, --help                help for add-observer-list
      --keygen-block int    set keygen block , default is 20 (default 20)
      --tss-pubkey string   set TSS pubkey if using older keygen
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)

## zetacored addr-conversion

convert a zeta1xxx address to validator operator address zetavaloper1xxx

### Synopsis


read a zeta1xxx or zetavaloper1xxx address and convert it to the other type;
it always outputs three lines; the first line is the zeta1xxx address, the second line is the zetavaloper1xxx address
and the third line is the ethereum address.
			

```
zetacored addr-conversion [zeta address] [flags]
```

### Options

```
  -h, --help   help for addr-conversion
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)

## zetacored collect-gentxs

Collect genesis txs and output a genesis.json file

```
zetacored collect-gentxs [flags]
```

### Options

```
      --gentx-dir string   override default "gentx" directory from which collect and execute genesis transactions; default [--home]/config/gentx/
  -h, --help               help for collect-gentxs
      --home string        The application home directory 
```

### Options inherited from parent commands

```
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)

## zetacored collect-observer-info

collect observer info into the genesis from a folder , default path is ~/.zetacored/os_info/ 


```
zetacored collect-observer-info [folder] [flags]
```

### Options

```
  -h, --help   help for collect-observer-info
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)

## zetacored config

Utilities for managing application configuration

### Options

```
  -h, --help   help for config
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)
* [zetacored config diff](#zetacored-config-diff)	 - Outputs all config values that are different from the app.toml defaults.
* [zetacored config get](#zetacored-config-get)	 - Get an application config value
* [zetacored config home](#zetacored-config-home)	 - Outputs the folder used as the binary home. No home directory is set when using the `confix` tool standalone.
* [zetacored config migrate](#zetacored-config-migrate)	 - Migrate Cosmos SDK app configuration file to the specified version
* [zetacored config set](#zetacored-config-set)	 - Set an application config value
* [zetacored config view](#zetacored-config-view)	 - View the config file

## zetacored config diff

Outputs all config values that are different from the app.toml defaults.

```
zetacored config diff [target-version] [app-toml-path] [flags]
```

### Options

```
  -h, --help   help for diff
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored config](#zetacored-config)	 - Utilities for managing application configuration

## zetacored config get

Get an application config value

### Synopsis

Get an application config value. The [config] argument must be the path of the file when using the `confix` tool standalone, otherwise it must be the name of the config file without the .toml extension.

```
zetacored config get [config] [key] [flags]
```

### Options

```
  -h, --help   help for get
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored config](#zetacored-config)	 - Utilities for managing application configuration

## zetacored config home

Outputs the folder used as the binary home. No home directory is set when using the `confix` tool standalone.

### Synopsis

Outputs the folder used as the binary home. In order to change the home directory path, set the $APPD_HOME environment variable, or use the "--home" flag.

```
zetacored config home [flags]
```

### Options

```
  -h, --help   help for home
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored config](#zetacored-config)	 - Utilities for managing application configuration

## zetacored config migrate

Migrate Cosmos SDK app configuration file to the specified version

### Synopsis

Migrate the contents of the Cosmos SDK app configuration (app.toml) to the specified version.
The output is written in-place unless --stdout is provided.
In case of any error in updating the file, no output is written.

```
zetacored config migrate [target-version] [app-toml-path] (options) [flags]
```

### Options

```
  -h, --help            help for migrate
      --skip-validate   skip configuration validation (allows to migrate unknown configurations)
      --stdout          print the updated config to stdout
      --verbose         log changes to stderr
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored config](#zetacored-config)	 - Utilities for managing application configuration

## zetacored config set

Set an application config value

### Synopsis

Set an application config value. The [config] argument must be the path of the file when using the `confix` tool standalone, otherwise it must be the name of the config file without the .toml extension.

```
zetacored config set [config] [key] [value] [flags]
```

### Options

```
  -h, --help            help for set
  -s, --skip-validate   skip configuration validation (allows to mutate unknown configurations)
      --stdout          print the updated config to stdout
  -v, --verbose         log changes to stderr
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored config](#zetacored-config)	 - Utilities for managing application configuration

## zetacored config view

View the config file

### Synopsis

View the config file. The [config] argument must be the path of the file when using the `confix` tool standalone, otherwise it must be the name of the config file without the .toml extension.

```
zetacored config view [config] [flags]
```

### Options

```
  -h, --help                   help for view
      --output-format string   Output format (json|toml) 
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored config](#zetacored-config)	 - Utilities for managing application configuration

## zetacored debug

Tool for helping with debugging your application

```
zetacored debug [flags]
```

### Options

```
  -h, --help   help for debug
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)
* [zetacored debug addr](#zetacored-debug-addr)	 - Convert an address between hex and bech32
* [zetacored debug codec](#zetacored-debug-codec)	 - Tool for helping with debugging your application codec
* [zetacored debug prefixes](#zetacored-debug-prefixes)	 - List prefixes used for Human-Readable Part (HRP) in Bech32
* [zetacored debug pubkey](#zetacored-debug-pubkey)	 - Decode a pubkey from proto JSON
* [zetacored debug pubkey-raw](#zetacored-debug-pubkey-raw)	 - Decode a ED25519 or secp256k1 pubkey from hex, base64, or bech32
* [zetacored debug raw-bytes](#zetacored-debug-raw-bytes)	 - Convert raw bytes output (eg. [10 21 13 255]) to hex

## zetacored debug addr

Convert an address between hex and bech32

### Synopsis

Convert an address between hex encoding and bech32.

Example:
$ zetacored debug addr cosmos1e0jnq2sun3dzjh8p2xq95kk0expwmd7shwjpfg
			

```
zetacored debug addr [address] [flags]
```

### Options

```
  -h, --help   help for addr
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored debug](#zetacored-debug)	 - Tool for helping with debugging your application

## zetacored debug codec

Tool for helping with debugging your application codec

```
zetacored debug codec [flags]
```

### Options

```
  -h, --help   help for codec
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored debug](#zetacored-debug)	 - Tool for helping with debugging your application
* [zetacored debug codec list-implementations](#zetacored-debug-codec-list-implementations)	 - List the registered type URLs for the provided interface
* [zetacored debug codec list-interfaces](#zetacored-debug-codec-list-interfaces)	 - List all registered interface type URLs

## zetacored debug codec list-implementations

List the registered type URLs for the provided interface

### Synopsis

List the registered type URLs that can be used for the provided interface name using the application codec

```
zetacored debug codec list-implementations [interface] [flags]
```

### Options

```
  -h, --help   help for list-implementations
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored debug codec](#zetacored-debug-codec)	 - Tool for helping with debugging your application codec

## zetacored debug codec list-interfaces

List all registered interface type URLs

### Synopsis

List all registered interface type URLs using the application codec

```
zetacored debug codec list-interfaces [flags]
```

### Options

```
  -h, --help   help for list-interfaces
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored debug codec](#zetacored-debug-codec)	 - Tool for helping with debugging your application codec

## zetacored debug prefixes

List prefixes used for Human-Readable Part (HRP) in Bech32

### Synopsis

List prefixes used in Bech32 addresses.

```
zetacored debug prefixes [flags]
```

### Examples

```
$ zetacored debug prefixes
```

### Options

```
  -h, --help   help for prefixes
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored debug](#zetacored-debug)	 - Tool for helping with debugging your application

## zetacored debug pubkey

Decode a pubkey from proto JSON

### Synopsis

Decode a pubkey from proto JSON and display it's address.

Example:
$ zetacored debug pubkey '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"AurroA7jvfPd1AadmmOvWM2rJSwipXfRf8yD6pLbA2DJ"}'
			

```
zetacored debug pubkey [pubkey] [flags]
```

### Options

```
  -h, --help   help for pubkey
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored debug](#zetacored-debug)	 - Tool for helping with debugging your application

## zetacored debug pubkey-raw

Decode a ED25519 or secp256k1 pubkey from hex, base64, or bech32

### Synopsis

Decode a pubkey from hex, base64, or bech32.

```
zetacored debug pubkey-raw [pubkey] -t [{ed25519, secp256k1}] [flags]
```

### Examples

```

zetacored debug pubkey-raw 8FCA9D6D1F80947FD5E9A05309259746F5F72541121766D5F921339DD061174A
zetacored debug pubkey-raw j8qdbR+AlH/V6aBTCSWXRvX3JUESF2bV+SEzndBhF0o=
zetacored debug pubkey-raw cosmospub1zcjduepq3l9f6mglsz28l40f5pfsjfvhgm6lwf2pzgtkd40eyyeem5rpza9q47axrz
			
```

### Options

```
  -h, --help          help for pubkey-raw
  -t, --type string   Pubkey type to decode (oneof secp256k1, ed25519) 
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored debug](#zetacored-debug)	 - Tool for helping with debugging your application

## zetacored debug raw-bytes

Convert raw bytes output (eg. [10 21 13 255]) to hex

### Synopsis

Convert raw-bytes to hex.

```
zetacored debug raw-bytes [raw-bytes] [flags]
```

### Examples

```
zetacored debug raw-bytes '[72 101 108 108 111 44 32 112 108 97 121 103 114 111 117 110 100]'
```

### Options

```
  -h, --help   help for raw-bytes
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored debug](#zetacored-debug)	 - Tool for helping with debugging your application

## zetacored docs

Generate markdown documentation for zetacored

```
zetacored docs [path] [flags]
```

### Options

```
  -h, --help          help for docs
      --path string   Path where the docs will be generated 
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)

## zetacored export

Export state to JSON

```
zetacored export [flags]
```

### Options

```
      --for-zero-height              Export state to start at height zero (perform preproccessing)
      --height int                   Export state from a particular height (-1 means latest height) (default -1)
  -h, --help                         help for export
      --home string                  The application home directory 
      --jail-allowed-addrs strings   Comma-separated list of operator addresses of jailed validators to unjail
      --modules-to-export strings    Comma-separated list of modules to export. If empty, will export all modules
      --output-document string       Exported state is written to the given file instead of STDOUT
```

### Options inherited from parent commands

```
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)

## zetacored gentx

Generate a genesis tx carrying a self delegation

### Synopsis

Generate a genesis transaction that creates a validator with a self-delegation,
that is signed by the key in the Keyring referenced by a given name. A node ID and consensus
pubkey may optionally be provided. If they are omitted, they will be retrieved from the priv_validator.json
file. The following default parameters are included:
    
	delegation amount:           100000000stake
	commission rate:             0.1
	commission max rate:         0.2
	commission max change rate:  0.01
	minimum self delegation:     1


Example:
$ zetacored gentx my-key-name 1000000stake --home=/path/to/home/dir --keyring-backend=os --chain-id=test-chain-1 \
    --moniker="myvalidator" \
    --commission-max-change-rate=0.01 \
    --commission-max-rate=1.0 \
    --commission-rate=0.07 \
    --details="..." \
    --security-contact="..." \
    --website="..."


```
zetacored gentx [key_name] [amount] [flags]
```

### Options

```
  -a, --account-number uint                 The account number of the signing account (offline mode only)
      --amount string                       Amount of coins to bond
      --aux                                 Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string               Transaction broadcasting mode (sync|async) 
      --chain-id string                     The network chain ID
      --commission-max-change-rate string   The maximum commission change rate percentage (per day)
      --commission-max-rate string          The maximum commission rate percentage
      --commission-rate string              The initial commission rate percentage
      --details string                      The validator's (optional) details
      --dry-run                             ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string                  Fee granter grants fees for the transaction
      --fee-payer string                    Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                         Fees to pay along with transaction; eg: 10uatom
      --from string                         Name or address of private key with which to sign
      --gas string                          gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float                adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string                   Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only                       Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                                help for gentx
      --home string                         The application home directory 
      --identity string                     The (optional) identity signature (ex. UPort or Keybase)
      --ip string                           The node's public P2P IP 
      --keyring-backend string              Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string                  The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                              Use a connected Ledger device
      --min-self-delegation string          The minimum self delegation required on the validator
      --moniker string                      The validator's (optional) moniker
      --node string                         [host]:[port] to CometBFT rpc interface for this chain 
      --node-id string                      The node's NodeID
      --note string                         Note to add a description to the transaction (previously --memo)
      --offline                             Offline mode (does not allow any online functionality)
      --output-document string              Write the genesis transaction JSON document to the given file instead of the default location
      --p2p-port uint                       The node's public P2P port (default 26656)
      --pubkey string                       The validator's Protobuf JSON encoded public key
      --security-contact string             The validator's (optional) security contact email
  -s, --sequence uint                       The sequence number of the signing account (offline mode only)
      --sign-mode string                    Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-height uint                 Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                          Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --website string                      The validator's (optional) website
  -y, --yes                                 Skip tx broadcasting prompt confirmation
```

### Options inherited from parent commands

```
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)

## zetacored get-pubkey

Get the node account public key

```
zetacored get-pubkey [tssKeyName] [password] [flags]
```

### Options

```
  -h, --help   help for get-pubkey
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)

## zetacored index-eth-tx

Index historical eth txs

### Synopsis

Index historical eth txs, it only support two traverse direction to avoid creating gaps in the indexer db if using arbitrary block ranges:
		- backward: index the blocks from the first indexed block to the earliest block in the chain, if indexer db is empty, start from the latest block.
		- forward: index the blocks from the latest indexed block to latest block in the chain.

		When start the node, the indexer start from the latest indexed block to avoid creating gap.
        Backward mode should be used most of the time, so the latest indexed block is always up-to-date.
		

```
zetacored index-eth-tx [backward|forward] [flags]
```

### Options

```
  -h, --help   help for index-eth-tx
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)

## zetacored init

Initialize private validator, p2p, genesis, and application configuration files

### Synopsis

Initialize validators's and node's configuration files.

```
zetacored init [moniker] [flags]
```

### Options

```
      --chain-id string        genesis file chain-id, if left blank will be randomly created
      --default-denom string   genesis file default denomination, if left blank default value is 'stake'
  -h, --help                   help for init
      --home string            node's home directory 
      --initial-height int     specify the initial block height at genesis (default 1)
  -o, --overwrite              overwrite the genesis.json file
      --recover                provide seed phrase to recover existing key instead of creating
```

### Options inherited from parent commands

```
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)

## zetacored keys

Manage your application's keys

### Synopsis

Keyring management commands. These keys may be in any format supported by the
Tendermint crypto library and can be used by light-clients, full nodes, or any other application
that needs to sign with a private key.

The keyring supports the following backends:

    os          Uses the operating system's default credentials store.
    file        Uses encrypted file-based keystore within the app's configuration directory.
                This keyring will request a password each time it is accessed, which may occur
                multiple times in a single command resulting in repeated password prompts.
    kwallet     Uses KDE Wallet Manager as a credentials management application.
    pass        Uses the pass command line utility to store and retrieve keys.
    test        Stores keys insecurely to disk. It does not prompt for a password to be unlocked
                and it should be use only for testing purposes.

kwallet and pass backends depend on external tools. Refer to their respective documentation for more
information:
    KWallet     https://github.com/KDE/kwallet
    pass        https://www.passwordstore.org/

The pass backend requires GnuPG: https://gnupg.org/


### Options

```
  -h, --help                     help for keys
      --home string              The application home directory 
      --keyring-backend string   Select keyring's backend (os|file|test) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --output string            Output format (text|json) 
```

### Options inherited from parent commands

```
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)
* [zetacored keys ](#zetacored-keys-)	 - 
* [zetacored keys add](#zetacored-keys-add)	 - Add an encrypted private key (either newly generated or recovered), encrypt it, and save to [name] file
* [zetacored keys delete](#zetacored-keys-delete)	 - Delete the given keys
* [zetacored keys export](#zetacored-keys-export)	 - Export private keys
* [zetacored keys import](#zetacored-keys-import)	 - Import private keys into the local keybase
* [zetacored keys list](#zetacored-keys-list)	 - List all keys
* [zetacored keys migrate](#zetacored-keys-migrate)	 - Migrate keys from amino to proto serialization format
* [zetacored keys mnemonic](#zetacored-keys-mnemonic)	 - Compute the bip39 mnemonic for some input entropy
* [zetacored keys parse](#zetacored-keys-parse)	 - Parse address from hex to bech32 and vice versa
* [zetacored keys rename](#zetacored-keys-rename)	 - Rename an existing key
* [zetacored keys show](#zetacored-keys-show)	 - Retrieve key information by name or address
* [zetacored keys unsafe-export-eth-key](#zetacored-keys-unsafe-export-eth-key)	 - **UNSAFE** Export an Ethereum private key
* [zetacored keys unsafe-import-eth-key](#zetacored-keys-unsafe-import-eth-key)	 - **UNSAFE** Import Ethereum private keys into the local keybase

## zetacored keys 



```
zetacored keys  [flags]
```

### Options

```
  -h, --help   help for this command
```

### Options inherited from parent commands

```
      --home string              The application home directory 
      --keyring-backend string   Select keyring's backend (os|file|test) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --log_format string        The logging format (json|plain) 
      --log_level string         The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color             Disable colored logs
      --output string            Output format (text|json) 
      --trace                    print out full stack trace on errors
```

### SEE ALSO

* [zetacored keys](#zetacored-keys)	 - Manage your application's keys

## zetacored keys add

Add an encrypted private key (either newly generated or recovered), encrypt it, and save to [name] file

### Synopsis

Derive a new private key and encrypt to disk.
Optionally specify a BIP39 mnemonic, a BIP39 passphrase to further secure the mnemonic,
and a bip32 HD path to derive a specific account. The key will be stored under the given name
and encrypted with the given password. The only input that is required is the encryption password.

If run with -i, it will prompt the user for BIP44 path, BIP39 mnemonic, and passphrase.
The flag --recover allows one to recover a key from a seed passphrase.
If run with --dry-run, a key would be generated (or recovered) but not stored to the
local keystore.
Use the --pubkey flag to add arbitrary public keys to the keystore for constructing
multisig transactions.

Use the --source flag to import mnemonic from a file in recover or interactive mode. 
Example:

	keys add testing --recover --source ./mnemonic.txt

You can create and store a multisig key by passing the list of key names stored in a keyring
and the minimum number of signatures required through --multisig-threshold. The keys are
sorted by address, unless the flag --nosort is set.
Example:

    keys add mymultisig --multisig "keyname1,keyname2,keyname3" --multisig-threshold 2


```
zetacored keys add [name] [flags]
```

### Options

```
      --account uint32           Account number for HD derivation (less than equal 2147483647)
      --coin-type uint32         coin type number for HD derivation (default 118)
      --dry-run                  Perform action, but don't add key to local keystore
      --hd-path string           Manual HD Path derivation (overrides BIP44 config) 
  -h, --help                     help for add
      --index uint32             Address index number for HD derivation (less than equal 2147483647)
  -i, --interactive              Interactively prompt user for BIP39 passphrase and mnemonic
      --key-type string          Key signing algorithm to generate keys for 
      --ledger                   Store a local reference to a private key on a Ledger device
      --multisig strings         List of key names stored in keyring to construct a public legacy multisig key
      --multisig-threshold int   K out of N required signatures. For use in conjunction with --multisig (default 1)
      --no-backup                Don't print out seed phrase (if others are watching the terminal)
      --nosort                   Keys passed to --multisig are taken in the order they're supplied
      --pubkey string            Parse a public key in JSON format and saves key info to [name] file.
      --pubkey-base64 string     Parse a public key in base64 format and saves key info.
      --recover                  Provide seed phrase to recover existing key instead of creating
      --source string            Import mnemonic from a file (only usable when recover or interactive is passed)
```

### Options inherited from parent commands

```
      --home string              The application home directory 
      --keyring-backend string   Select keyring's backend (os|file|test) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --log_format string        The logging format (json|plain) 
      --log_level string         The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color             Disable colored logs
      --output string            Output format (text|json) 
      --trace                    print out full stack trace on errors
```

### SEE ALSO

* [zetacored keys](#zetacored-keys)	 - Manage your application's keys

## zetacored keys delete

Delete the given keys

### Synopsis

Delete keys from the Keybase backend.

Note that removing offline or ledger keys will remove
only the public key references stored locally, i.e.
private keys stored in a ledger device cannot be deleted with the CLI.


```
zetacored keys delete [name]... [flags]
```

### Options

```
  -f, --force   Remove the key unconditionally without asking for the passphrase. Deprecated.
  -h, --help    help for delete
  -y, --yes     Skip confirmation prompt when deleting offline or ledger key references
```

### Options inherited from parent commands

```
      --home string              The application home directory 
      --keyring-backend string   Select keyring's backend (os|file|test) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --log_format string        The logging format (json|plain) 
      --log_level string         The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color             Disable colored logs
      --output string            Output format (text|json) 
      --trace                    print out full stack trace on errors
```

### SEE ALSO

* [zetacored keys](#zetacored-keys)	 - Manage your application's keys

## zetacored keys export

Export private keys

### Synopsis

Export a private key from the local keyring in ASCII-armored encrypted format.

When both the --unarmored-hex and --unsafe flags are selected, cryptographic
private key material is exported in an INSECURE fashion that is designed to
allow users to import their keys in hot wallets. This feature is for advanced
users only that are confident about how to handle private keys work and are
FULLY AWARE OF THE RISKS. If you are unsure, you may want to do some research
and export your keys in ASCII-armored encrypted format.

```
zetacored keys export [name] [flags]
```

### Options

```
  -h, --help            help for export
      --unarmored-hex   Export unarmored hex privkey. Requires --unsafe.
      --unsafe          Enable unsafe operations. This flag must be switched on along with all unsafe operation-specific options.
```

### Options inherited from parent commands

```
      --home string              The application home directory 
      --keyring-backend string   Select keyring's backend (os|file|test) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --log_format string        The logging format (json|plain) 
      --log_level string         The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color             Disable colored logs
      --output string            Output format (text|json) 
      --trace                    print out full stack trace on errors
```

### SEE ALSO

* [zetacored keys](#zetacored-keys)	 - Manage your application's keys

## zetacored keys import

Import private keys into the local keybase

### Synopsis

Import a ASCII armored private key into the local keybase.

```
zetacored keys import [name] [keyfile] [flags]
```

### Options

```
  -h, --help   help for import
```

### Options inherited from parent commands

```
      --home string              The application home directory 
      --keyring-backend string   Select keyring's backend (os|file|test) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --log_format string        The logging format (json|plain) 
      --log_level string         The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color             Disable colored logs
      --output string            Output format (text|json) 
      --trace                    print out full stack trace on errors
```

### SEE ALSO

* [zetacored keys](#zetacored-keys)	 - Manage your application's keys

## zetacored keys list

List all keys

### Synopsis

Return a list of all public keys stored by this key manager
along with their associated name and address.

```
zetacored keys list [flags]
```

### Options

```
  -h, --help         help for list
  -n, --list-names   List names only
```

### Options inherited from parent commands

```
      --home string              The application home directory 
      --keyring-backend string   Select keyring's backend (os|file|test) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --log_format string        The logging format (json|plain) 
      --log_level string         The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color             Disable colored logs
      --output string            Output format (text|json) 
      --trace                    print out full stack trace on errors
```

### SEE ALSO

* [zetacored keys](#zetacored-keys)	 - Manage your application's keys

## zetacored keys migrate

Migrate keys from amino to proto serialization format

### Synopsis

Migrate keys from Amino to Protocol Buffers records.
For each key material entry, the command will check if the key can be deserialized using proto.
If this is the case, the key is already migrated. Therefore, we skip it and continue with a next one. 
Otherwise, we try to deserialize it using Amino into LegacyInfo. If this attempt is successful, we serialize 
LegacyInfo to Protobuf serialization format and overwrite the keyring entry. If any error occurred, it will be 
outputted in CLI and migration will be continued until all keys in the keyring DB are exhausted.
See https://github.com/cosmos/cosmos-sdk/pull/9695 for more details.


```
zetacored keys migrate [flags]
```

### Options

```
  -h, --help   help for migrate
```

### Options inherited from parent commands

```
      --home string              The application home directory 
      --keyring-backend string   Select keyring's backend (os|file|test) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --log_format string        The logging format (json|plain) 
      --log_level string         The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color             Disable colored logs
      --output string            Output format (text|json) 
      --trace                    print out full stack trace on errors
```

### SEE ALSO

* [zetacored keys](#zetacored-keys)	 - Manage your application's keys

## zetacored keys mnemonic

Compute the bip39 mnemonic for some input entropy

### Synopsis

Create a bip39 mnemonic, sometimes called a seed phrase, by reading from the system entropy. To pass your own entropy, use --unsafe-entropy

```
zetacored keys mnemonic [flags]
```

### Options

```
  -h, --help             help for mnemonic
      --unsafe-entropy   Prompt the user to supply their own entropy, instead of relying on the system
```

### Options inherited from parent commands

```
      --home string              The application home directory 
      --keyring-backend string   Select keyring's backend (os|file|test) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --log_format string        The logging format (json|plain) 
      --log_level string         The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color             Disable colored logs
      --output string            Output format (text|json) 
      --trace                    print out full stack trace on errors
```

### SEE ALSO

* [zetacored keys](#zetacored-keys)	 - Manage your application's keys

## zetacored keys parse

Parse address from hex to bech32 and vice versa

### Synopsis

Convert and print to stdout key addresses and fingerprints from
hexadecimal into bech32 cosmos prefixed format and vice versa.


```
zetacored keys parse [hex-or-bech32-address] [flags]
```

### Options

```
  -h, --help   help for parse
```

### Options inherited from parent commands

```
      --home string              The application home directory 
      --keyring-backend string   Select keyring's backend (os|file|test) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --log_format string        The logging format (json|plain) 
      --log_level string         The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color             Disable colored logs
      --output string            Output format (text|json) 
      --trace                    print out full stack trace on errors
```

### SEE ALSO

* [zetacored keys](#zetacored-keys)	 - Manage your application's keys

## zetacored keys rename

Rename an existing key

### Synopsis

Rename a key from the Keybase backend.

Note that renaming offline or ledger keys will rename
only the public key references stored locally, i.e.
private keys stored in a ledger device cannot be renamed with the CLI.


```
zetacored keys rename [old_name] [new_name] [flags]
```

### Options

```
  -h, --help   help for rename
  -y, --yes    Skip confirmation prompt when renaming offline or ledger key references
```

### Options inherited from parent commands

```
      --home string              The application home directory 
      --keyring-backend string   Select keyring's backend (os|file|test) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --log_format string        The logging format (json|plain) 
      --log_level string         The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color             Disable colored logs
      --output string            Output format (text|json) 
      --trace                    print out full stack trace on errors
```

### SEE ALSO

* [zetacored keys](#zetacored-keys)	 - Manage your application's keys

## zetacored keys show

Retrieve key information by name or address

### Synopsis

Display keys details. If multiple names or addresses are provided,
then an ephemeral multisig key will be created under the name "multi"
consisting of all the keys provided by name and multisig threshold.

```
zetacored keys show [name_or_address [name_or_address...]] [flags]
```

### Options

```
  -a, --address                  Output the address only (cannot be used with --output)
      --bech string              The Bech32 prefix encoding for a key (acc|val|cons) 
  -d, --device                   Output the address in a ledger device (cannot be used with --pubkey)
  -h, --help                     help for show
      --multisig-threshold int   K out of N required signatures (default 1)
  -p, --pubkey                   Output the public key only (cannot be used with --output)
```

### Options inherited from parent commands

```
      --home string              The application home directory 
      --keyring-backend string   Select keyring's backend (os|file|test) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --log_format string        The logging format (json|plain) 
      --log_level string         The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color             Disable colored logs
      --output string            Output format (text|json) 
      --trace                    print out full stack trace on errors
```

### SEE ALSO

* [zetacored keys](#zetacored-keys)	 - Manage your application's keys

## zetacored keys unsafe-export-eth-key

**UNSAFE** Export an Ethereum private key

### Synopsis

**UNSAFE** Export an Ethereum private key unencrypted to use in dev tooling

```
zetacored keys unsafe-export-eth-key [name] [flags]
```

### Options

```
  -h, --help   help for unsafe-export-eth-key
```

### Options inherited from parent commands

```
      --home string              The application home directory 
      --keyring-backend string   Select keyring's backend (os|file|test) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --log_format string        The logging format (json|plain) 
      --log_level string         The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color             Disable colored logs
      --output string            Output format (text|json) 
      --trace                    print out full stack trace on errors
```

### SEE ALSO

* [zetacored keys](#zetacored-keys)	 - Manage your application's keys

## zetacored keys unsafe-import-eth-key

**UNSAFE** Import Ethereum private keys into the local keybase

### Synopsis

**UNSAFE** Import a hex-encoded Ethereum private key into the local keybase.

```
zetacored keys unsafe-import-eth-key [name] [pk] [flags]
```

### Options

```
  -h, --help   help for unsafe-import-eth-key
```

### Options inherited from parent commands

```
      --home string              The application home directory 
      --keyring-backend string   Select keyring's backend (os|file|test) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --log_format string        The logging format (json|plain) 
      --log_level string         The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color             Disable colored logs
      --output string            Output format (text|json) 
      --trace                    print out full stack trace on errors
```

### SEE ALSO

* [zetacored keys](#zetacored-keys)	 - Manage your application's keys

## zetacored parse-genesis-file

Parse the provided genesis file and import the required data into the optionally provided genesis file

```
zetacored parse-genesis-file [import-genesis-file] [optional-genesis-file] [flags]
```

### Options

```
  -h, --help     help for parse-genesis-file
      --modify   modify the genesis file before importing
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)

## zetacored query

Querying subcommands

```
zetacored query [flags]
```

### Options

```
      --chain-id string   The network chain ID
  -h, --help              help for query
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)
* [zetacored query authority](#zetacored-query-authority)	 - Querying commands for the authority module
* [zetacored query comet-validator-set](#zetacored-query-comet-validator-set)	 - Get the full CometBFT validator set at given height
* [zetacored query crosschain](#zetacored-query-crosschain)	 - Querying commands for the crosschain module
* [zetacored query emissions](#zetacored-query-emissions)	 - Querying commands for the emissions module
* [zetacored query evm](#zetacored-query-evm)	 - Querying commands for the evm module
* [zetacored query feemarket](#zetacored-query-feemarket)	 - Querying commands for the fee market module
* [zetacored query fungible](#zetacored-query-fungible)	 - Querying commands for the fungible module
* [zetacored query lightclient](#zetacored-query-lightclient)	 - Querying commands for the lightclient module
* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module
* [zetacored query tx](#zetacored-query-tx)	 - Query for a transaction by hash, "[addr]/[seq]" combination or comma-separated signatures in a committed block
* [zetacored query txs](#zetacored-query-txs)	 - Query for paginated transactions that match a set of events

## zetacored query authority

Querying commands for the authority module

```
zetacored query authority [flags]
```

### Options

```
  -h, --help   help for authority
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query](#zetacored-query)	 - Querying subcommands
* [zetacored query authority list-authorizations](#zetacored-query-authority-list-authorizations)	 - lists all authorizations
* [zetacored query authority show-authorization](#zetacored-query-authority-show-authorization)	 - shows the authorization for a given message URL
* [zetacored query authority show-chain-info](#zetacored-query-authority-show-chain-info)	 - show the chain info
* [zetacored query authority show-policies](#zetacored-query-authority-show-policies)	 - show the policies

## zetacored query authority list-authorizations

lists all authorizations

```
zetacored query authority list-authorizations [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-authorizations
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query authority](#zetacored-query-authority)	 - Querying commands for the authority module

## zetacored query authority show-authorization

shows the authorization for a given message URL

```
zetacored query authority show-authorization [msg-url] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-authorization
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query authority](#zetacored-query-authority)	 - Querying commands for the authority module

## zetacored query authority show-chain-info

show the chain info

```
zetacored query authority show-chain-info [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-chain-info
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query authority](#zetacored-query-authority)	 - Querying commands for the authority module

## zetacored query authority show-policies

show the policies

```
zetacored query authority show-policies [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-policies
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query authority](#zetacored-query-authority)	 - Querying commands for the authority module

## zetacored query comet-validator-set

Get the full CometBFT validator set at given height

```
zetacored query comet-validator-set [height] [flags]
```

### Options

```
  -h, --help            help for comet-validator-set
      --limit int       Query number of results returned per page (default 100)
      --node string     [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string   Output format (text|json) 
      --page int        Query a specific page of paginated results (default 1)
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query](#zetacored-query)	 - Querying subcommands

## zetacored query crosschain

Querying commands for the crosschain module

```
zetacored query crosschain [flags]
```

### Options

```
  -h, --help   help for crosschain
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query](#zetacored-query)	 - Querying subcommands
* [zetacored query crosschain get-zeta-accounting](#zetacored-query-crosschain-get-zeta-accounting)	 - Query zeta accounting
* [zetacored query crosschain inbound-hash-to-cctx-data](#zetacored-query-crosschain-inbound-hash-to-cctx-data)	 - query a cctx data from a inbound hash
* [zetacored query crosschain last-zeta-height](#zetacored-query-crosschain-last-zeta-height)	 - Query last Zeta Height
* [zetacored query crosschain list-all-inbound-trackers](#zetacored-query-crosschain-list-all-inbound-trackers)	 - shows all inbound trackers
* [zetacored query crosschain list-cctx](#zetacored-query-crosschain-list-cctx)	 - list all CCTX
* [zetacored query crosschain list-gas-price](#zetacored-query-crosschain-list-gas-price)	 - list all gasPrice
* [zetacored query crosschain list-inbound-hash-to-cctx](#zetacored-query-crosschain-list-inbound-hash-to-cctx)	 - list all inboundHashToCctx
* [zetacored query crosschain list-inbound-tracker](#zetacored-query-crosschain-list-inbound-tracker)	 - shows a list of inbound trackers by chainId
* [zetacored query crosschain list-outbound-tracker](#zetacored-query-crosschain-list-outbound-tracker)	 - list all outbound trackers
* [zetacored query crosschain list-pending-cctx](#zetacored-query-crosschain-list-pending-cctx)	 - shows pending CCTX
* [zetacored query crosschain list_pending_cctx_within_rate_limit](#zetacored-query-crosschain-list-pending-cctx-within-rate-limit)	 - list all pending CCTX within rate limit
* [zetacored query crosschain show-cctx](#zetacored-query-crosschain-show-cctx)	 - shows a CCTX
* [zetacored query crosschain show-gas-price](#zetacored-query-crosschain-show-gas-price)	 - shows a gasPrice
* [zetacored query crosschain show-inbound-hash-to-cctx](#zetacored-query-crosschain-show-inbound-hash-to-cctx)	 - shows a inboundHashToCctx
* [zetacored query crosschain show-inbound-tracker](#zetacored-query-crosschain-show-inbound-tracker)	 - shows an inbound tracker by chainID and txHash
* [zetacored query crosschain show-outbound-tracker](#zetacored-query-crosschain-show-outbound-tracker)	 - shows an outbound tracker
* [zetacored query crosschain show-rate-limiter-flags](#zetacored-query-crosschain-show-rate-limiter-flags)	 - shows the rate limiter flags

## zetacored query crosschain get-zeta-accounting

Query zeta accounting

```
zetacored query crosschain get-zeta-accounting [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for get-zeta-accounting
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query crosschain](#zetacored-query-crosschain)	 - Querying commands for the crosschain module

## zetacored query crosschain inbound-hash-to-cctx-data

query a cctx data from a inbound hash

```
zetacored query crosschain inbound-hash-to-cctx-data [inbound-hash] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for inbound-hash-to-cctx-data
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query crosschain](#zetacored-query-crosschain)	 - Querying commands for the crosschain module

## zetacored query crosschain last-zeta-height

Query last Zeta Height

```
zetacored query crosschain last-zeta-height [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for last-zeta-height
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query crosschain](#zetacored-query-crosschain)	 - Querying commands for the crosschain module

## zetacored query crosschain list-all-inbound-trackers

shows all inbound trackers

```
zetacored query crosschain list-all-inbound-trackers [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-all-inbound-trackers
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query crosschain](#zetacored-query-crosschain)	 - Querying commands for the crosschain module

## zetacored query crosschain list-cctx

list all CCTX

```
zetacored query crosschain list-cctx [flags]
```

### Options

```
      --count-total        count total number of records in list-cctx to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-cctx
      --limit uint         pagination limit of list-cctx to query for (default 100)
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
      --offset uint        pagination offset of list-cctx to query for
  -o, --output string      Output format (text|json) 
      --page uint          pagination page of list-cctx to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of list-cctx to query for
      --reverse            results are sorted in descending order
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query crosschain](#zetacored-query-crosschain)	 - Querying commands for the crosschain module

## zetacored query crosschain list-gas-price

list all gasPrice

```
zetacored query crosschain list-gas-price [flags]
```

### Options

```
      --count-total        count total number of records in list-gas-price to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-gas-price
      --limit uint         pagination limit of list-gas-price to query for (default 100)
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
      --offset uint        pagination offset of list-gas-price to query for
  -o, --output string      Output format (text|json) 
      --page uint          pagination page of list-gas-price to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of list-gas-price to query for
      --reverse            results are sorted in descending order
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query crosschain](#zetacored-query-crosschain)	 - Querying commands for the crosschain module

## zetacored query crosschain list-inbound-hash-to-cctx

list all inboundHashToCctx

```
zetacored query crosschain list-inbound-hash-to-cctx [flags]
```

### Options

```
      --count-total        count total number of records in list-inbound-hash-to-cctx to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-inbound-hash-to-cctx
      --limit uint         pagination limit of list-inbound-hash-to-cctx to query for (default 100)
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
      --offset uint        pagination offset of list-inbound-hash-to-cctx to query for
  -o, --output string      Output format (text|json) 
      --page uint          pagination page of list-inbound-hash-to-cctx to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of list-inbound-hash-to-cctx to query for
      --reverse            results are sorted in descending order
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query crosschain](#zetacored-query-crosschain)	 - Querying commands for the crosschain module

## zetacored query crosschain list-inbound-tracker

shows a list of inbound trackers by chainId

```
zetacored query crosschain list-inbound-tracker [chainId] [flags]
```

### Options

```
      --count-total        count total number of records in list-inbound-tracker [chainId] to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-inbound-tracker
      --limit uint         pagination limit of list-inbound-tracker [chainId] to query for (default 100)
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
      --offset uint        pagination offset of list-inbound-tracker [chainId] to query for
  -o, --output string      Output format (text|json) 
      --page uint          pagination page of list-inbound-tracker [chainId] to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of list-inbound-tracker [chainId] to query for
      --reverse            results are sorted in descending order
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query crosschain](#zetacored-query-crosschain)	 - Querying commands for the crosschain module

## zetacored query crosschain list-outbound-tracker

list all outbound trackers

```
zetacored query crosschain list-outbound-tracker [flags]
```

### Options

```
      --count-total        count total number of records in list-outbound-tracker to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-outbound-tracker
      --limit uint         pagination limit of list-outbound-tracker to query for (default 100)
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
      --offset uint        pagination offset of list-outbound-tracker to query for
  -o, --output string      Output format (text|json) 
      --page uint          pagination page of list-outbound-tracker to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of list-outbound-tracker to query for
      --reverse            results are sorted in descending order
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query crosschain](#zetacored-query-crosschain)	 - Querying commands for the crosschain module

## zetacored query crosschain list-pending-cctx

shows pending CCTX

```
zetacored query crosschain list-pending-cctx [chain-id] [limit] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-pending-cctx
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query crosschain](#zetacored-query-crosschain)	 - Querying commands for the crosschain module

## zetacored query crosschain list_pending_cctx_within_rate_limit

list all pending CCTX within rate limit

```
zetacored query crosschain list_pending_cctx_within_rate_limit [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list_pending_cctx_within_rate_limit
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query crosschain](#zetacored-query-crosschain)	 - Querying commands for the crosschain module

## zetacored query crosschain show-cctx

shows a CCTX

```
zetacored query crosschain show-cctx [index] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-cctx
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query crosschain](#zetacored-query-crosschain)	 - Querying commands for the crosschain module

## zetacored query crosschain show-gas-price

shows a gasPrice

```
zetacored query crosschain show-gas-price [index] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-gas-price
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query crosschain](#zetacored-query-crosschain)	 - Querying commands for the crosschain module

## zetacored query crosschain show-inbound-hash-to-cctx

shows a inboundHashToCctx

```
zetacored query crosschain show-inbound-hash-to-cctx [inbound-hash] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-inbound-hash-to-cctx
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query crosschain](#zetacored-query-crosschain)	 - Querying commands for the crosschain module

## zetacored query crosschain show-inbound-tracker

shows an inbound tracker by chainID and txHash

```
zetacored query crosschain show-inbound-tracker [chainID] [txHash] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-inbound-tracker
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query crosschain](#zetacored-query-crosschain)	 - Querying commands for the crosschain module

## zetacored query crosschain show-outbound-tracker

shows an outbound tracker

```
zetacored query crosschain show-outbound-tracker [chainId] [nonce] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-outbound-tracker
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query crosschain](#zetacored-query-crosschain)	 - Querying commands for the crosschain module

## zetacored query crosschain show-rate-limiter-flags

shows the rate limiter flags

```
zetacored query crosschain show-rate-limiter-flags [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-rate-limiter-flags
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query crosschain](#zetacored-query-crosschain)	 - Querying commands for the crosschain module

## zetacored query emissions

Querying commands for the emissions module

```
zetacored query emissions [flags]
```

### Options

```
  -h, --help   help for emissions
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query](#zetacored-query)	 - Querying subcommands
* [zetacored query emissions list-pool-addresses](#zetacored-query-emissions-list-pool-addresses)	 - Query list-pool-addresses
* [zetacored query emissions params](#zetacored-query-emissions-params)	 - shows the parameters of the module
* [zetacored query emissions show-available-emissions](#zetacored-query-emissions-show-available-emissions)	 - Query show-available-emissions

## zetacored query emissions list-pool-addresses

Query list-pool-addresses

```
zetacored query emissions list-pool-addresses [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-pool-addresses
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query emissions](#zetacored-query-emissions)	 - Querying commands for the emissions module

## zetacored query emissions params

shows the parameters of the module

```
zetacored query emissions params [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for params
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query emissions](#zetacored-query-emissions)	 - Querying commands for the emissions module

## zetacored query emissions show-available-emissions

Query show-available-emissions

```
zetacored query emissions show-available-emissions [address] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-available-emissions
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query emissions](#zetacored-query-emissions)	 - Querying commands for the emissions module

## zetacored query evm

Querying commands for the evm module

```
zetacored query evm [flags]
```

### Options

```
  -h, --help   help for evm
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query](#zetacored-query)	 - Querying subcommands
* [zetacored query evm code](#zetacored-query-evm-code)	 - Gets code from an account
* [zetacored query evm params](#zetacored-query-evm-params)	 - Get the evm params
* [zetacored query evm storage](#zetacored-query-evm-storage)	 - Gets storage for an account with a given key and height

## zetacored query evm code

Gets code from an account

### Synopsis

Gets code from an account. If the height is not provided, it will use the latest height from context.

```
zetacored query evm code ADDRESS [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for code
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query evm](#zetacored-query-evm)	 - Querying commands for the evm module

## zetacored query evm params

Get the evm params

### Synopsis

Get the evm parameter values.

```
zetacored query evm params [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for params
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query evm](#zetacored-query-evm)	 - Querying commands for the evm module

## zetacored query evm storage

Gets storage for an account with a given key and height

### Synopsis

Gets storage for an account with a given key and height. If the height is not provided, it will use the latest height from context.

```
zetacored query evm storage ADDRESS KEY [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for storage
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query evm](#zetacored-query-evm)	 - Querying commands for the evm module

## zetacored query feemarket

Querying commands for the fee market module

```
zetacored query feemarket [flags]
```

### Options

```
  -h, --help   help for feemarket
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query](#zetacored-query)	 - Querying subcommands
* [zetacored query feemarket base-fee](#zetacored-query-feemarket-base-fee)	 - Get the base fee amount at a given block height
* [zetacored query feemarket block-gas](#zetacored-query-feemarket-block-gas)	 - Get the block gas used at a given block height
* [zetacored query feemarket params](#zetacored-query-feemarket-params)	 - Get the fee market params

## zetacored query feemarket base-fee

Get the base fee amount at a given block height

### Synopsis

Get the base fee amount at a given block height.
If the height is not provided, it will use the latest height from context.

```
zetacored query feemarket base-fee [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for base-fee
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query feemarket](#zetacored-query-feemarket)	 - Querying commands for the fee market module

## zetacored query feemarket block-gas

Get the block gas used at a given block height

### Synopsis

Get the block gas used at a given block height.
If the height is not provided, it will use the latest height from context

```
zetacored query feemarket block-gas [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for block-gas
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query feemarket](#zetacored-query-feemarket)	 - Querying commands for the fee market module

## zetacored query feemarket params

Get the fee market params

### Synopsis

Get the fee market parameter values.

```
zetacored query feemarket params [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for params
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query feemarket](#zetacored-query-feemarket)	 - Querying commands for the fee market module

## zetacored query fungible

Querying commands for the fungible module

```
zetacored query fungible [flags]
```

### Options

```
  -h, --help   help for fungible
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query](#zetacored-query)	 - Querying subcommands
* [zetacored query fungible code-hash](#zetacored-query-fungible-code-hash)	 - shows the code hash of an account
* [zetacored query fungible gas-stability-pool-address](#zetacored-query-fungible-gas-stability-pool-address)	 - query the address of a gas stability pool
* [zetacored query fungible gas-stability-pool-balance](#zetacored-query-fungible-gas-stability-pool-balance)	 - query the balance of a gas stability pool for a chain
* [zetacored query fungible gas-stability-pool-balances](#zetacored-query-fungible-gas-stability-pool-balances)	 - query all gas stability pool balances
* [zetacored query fungible list-foreign-coins](#zetacored-query-fungible-list-foreign-coins)	 - list all ForeignCoins
* [zetacored query fungible show-foreign-coins](#zetacored-query-fungible-show-foreign-coins)	 - shows a ForeignCoins
* [zetacored query fungible system-contract](#zetacored-query-fungible-system-contract)	 - query system contract

## zetacored query fungible code-hash

shows the code hash of an account

```
zetacored query fungible code-hash [address] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for code-hash
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query fungible](#zetacored-query-fungible)	 - Querying commands for the fungible module

## zetacored query fungible gas-stability-pool-address

query the address of a gas stability pool

```
zetacored query fungible gas-stability-pool-address [flags]
```

### Options

```
      --count-total        count total number of records in gas-stability-pool-address to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for gas-stability-pool-address
      --limit uint         pagination limit of gas-stability-pool-address to query for (default 100)
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
      --offset uint        pagination offset of gas-stability-pool-address to query for
  -o, --output string      Output format (text|json) 
      --page uint          pagination page of gas-stability-pool-address to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of gas-stability-pool-address to query for
      --reverse            results are sorted in descending order
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query fungible](#zetacored-query-fungible)	 - Querying commands for the fungible module

## zetacored query fungible gas-stability-pool-balance

query the balance of a gas stability pool for a chain

```
zetacored query fungible gas-stability-pool-balance [chain-id] [flags]
```

### Options

```
      --count-total        count total number of records in gas-stability-pool-balance [chain-id] to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for gas-stability-pool-balance
      --limit uint         pagination limit of gas-stability-pool-balance [chain-id] to query for (default 100)
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
      --offset uint        pagination offset of gas-stability-pool-balance [chain-id] to query for
  -o, --output string      Output format (text|json) 
      --page uint          pagination page of gas-stability-pool-balance [chain-id] to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of gas-stability-pool-balance [chain-id] to query for
      --reverse            results are sorted in descending order
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query fungible](#zetacored-query-fungible)	 - Querying commands for the fungible module

## zetacored query fungible gas-stability-pool-balances

query all gas stability pool balances

```
zetacored query fungible gas-stability-pool-balances [flags]
```

### Options

```
      --count-total        count total number of records in gas-stability-pool-balances to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for gas-stability-pool-balances
      --limit uint         pagination limit of gas-stability-pool-balances to query for (default 100)
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
      --offset uint        pagination offset of gas-stability-pool-balances to query for
  -o, --output string      Output format (text|json) 
      --page uint          pagination page of gas-stability-pool-balances to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of gas-stability-pool-balances to query for
      --reverse            results are sorted in descending order
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query fungible](#zetacored-query-fungible)	 - Querying commands for the fungible module

## zetacored query fungible list-foreign-coins

list all ForeignCoins

```
zetacored query fungible list-foreign-coins [flags]
```

### Options

```
      --count-total        count total number of records in list-foreign-coins to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-foreign-coins
      --limit uint         pagination limit of list-foreign-coins to query for (default 100)
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
      --offset uint        pagination offset of list-foreign-coins to query for
  -o, --output string      Output format (text|json) 
      --page uint          pagination page of list-foreign-coins to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of list-foreign-coins to query for
      --reverse            results are sorted in descending order
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query fungible](#zetacored-query-fungible)	 - Querying commands for the fungible module

## zetacored query fungible show-foreign-coins

shows a ForeignCoins

```
zetacored query fungible show-foreign-coins [index] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-foreign-coins
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query fungible](#zetacored-query-fungible)	 - Querying commands for the fungible module

## zetacored query fungible system-contract

query system contract

```
zetacored query fungible system-contract [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for system-contract
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query fungible](#zetacored-query-fungible)	 - Querying commands for the fungible module

## zetacored query lightclient

Querying commands for the lightclient module

```
zetacored query lightclient [flags]
```

### Options

```
  -h, --help   help for lightclient
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query](#zetacored-query)	 - Querying subcommands
* [zetacored query lightclient list-block-header](#zetacored-query-lightclient-list-block-header)	 - List all the block headers
* [zetacored query lightclient list-chain-state](#zetacored-query-lightclient-list-chain-state)	 - List all the chain states
* [zetacored query lightclient show-block-header](#zetacored-query-lightclient-show-block-header)	 - Show a block header from its hash
* [zetacored query lightclient show-chain-state](#zetacored-query-lightclient-show-chain-state)	 - Show a chain state from its chain id
* [zetacored query lightclient show-header-enabled-chains](#zetacored-query-lightclient-show-header-enabled-chains)	 - Show the verification flags

## zetacored query lightclient list-block-header

List all the block headers

```
zetacored query lightclient list-block-header [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-block-header
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query lightclient](#zetacored-query-lightclient)	 - Querying commands for the lightclient module

## zetacored query lightclient list-chain-state

List all the chain states

```
zetacored query lightclient list-chain-state [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-chain-state
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query lightclient](#zetacored-query-lightclient)	 - Querying commands for the lightclient module

## zetacored query lightclient show-block-header

Show a block header from its hash

```
zetacored query lightclient show-block-header [block-hash] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-block-header
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query lightclient](#zetacored-query-lightclient)	 - Querying commands for the lightclient module

## zetacored query lightclient show-chain-state

Show a chain state from its chain id

```
zetacored query lightclient show-chain-state [chain-id] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-chain-state
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query lightclient](#zetacored-query-lightclient)	 - Querying commands for the lightclient module

## zetacored query lightclient show-header-enabled-chains

Show the verification flags

```
zetacored query lightclient show-header-enabled-chains [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-header-enabled-chains
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query lightclient](#zetacored-query-lightclient)	 - Querying commands for the lightclient module

## zetacored query observer

Querying commands for the observer module

```
zetacored query observer [flags]
```

### Options

```
  -h, --help   help for observer
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query](#zetacored-query)	 - Querying subcommands
* [zetacored query observer get-historical-tss-address](#zetacored-query-observer-get-historical-tss-address)	 - Query tss address by finalized zeta height (for historical tss addresses)
* [zetacored query observer get-tss-address](#zetacored-query-observer-get-tss-address)	 - Query current tss address
* [zetacored query observer list-blame](#zetacored-query-observer-list-blame)	 - Query AllBlameRecords
* [zetacored query observer list-blame-by-msg](#zetacored-query-observer-list-blame-by-msg)	 - Query AllBlameRecords
* [zetacored query observer list-chain-nonces](#zetacored-query-observer-list-chain-nonces)	 - list all chainNonces
* [zetacored query observer list-chain-params](#zetacored-query-observer-list-chain-params)	 - Query GetChainParams
* [zetacored query observer list-chains](#zetacored-query-observer-list-chains)	 - list all SupportedChains
* [zetacored query observer list-node-account](#zetacored-query-observer-list-node-account)	 - list all NodeAccount
* [zetacored query observer list-observer-set](#zetacored-query-observer-list-observer-set)	 - Query observer set
* [zetacored query observer list-pending-nonces](#zetacored-query-observer-list-pending-nonces)	 - shows a chainNonces
* [zetacored query observer list-tss-funds-migrator](#zetacored-query-observer-list-tss-funds-migrator)	 - list all tss funds migrators
* [zetacored query observer list-tss-history](#zetacored-query-observer-list-tss-history)	 - show historical list of TSS
* [zetacored query observer show-ballot](#zetacored-query-observer-show-ballot)	 - Query BallotByIdentifier
* [zetacored query observer show-blame](#zetacored-query-observer-show-blame)	 - Query BlameByIdentifier
* [zetacored query observer show-chain-nonces](#zetacored-query-observer-show-chain-nonces)	 - shows a chainNonces
* [zetacored query observer show-chain-params](#zetacored-query-observer-show-chain-params)	 - Query GetChainParamsForChain
* [zetacored query observer show-crosschain-flags](#zetacored-query-observer-show-crosschain-flags)	 - shows the crosschain flags
* [zetacored query observer show-keygen](#zetacored-query-observer-show-keygen)	 - shows keygen
* [zetacored query observer show-node-account](#zetacored-query-observer-show-node-account)	 - shows a NodeAccount
* [zetacored query observer show-observer-count](#zetacored-query-observer-show-observer-count)	 - Query show-observer-count
* [zetacored query observer show-operational-flags](#zetacored-query-observer-show-operational-flags)	 - shows the operational flags
* [zetacored query observer show-tss](#zetacored-query-observer-show-tss)	 - shows a TSS
* [zetacored query observer show-tss-funds-migrator](#zetacored-query-observer-show-tss-funds-migrator)	 - show the tss funds migrator for a chain

## zetacored query observer get-historical-tss-address

Query tss address by finalized zeta height (for historical tss addresses)

```
zetacored query observer get-historical-tss-address [finalizedZetaHeight] [bitcoinChainId] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for get-historical-tss-address
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module

## zetacored query observer get-tss-address

Query current tss address

```
zetacored query observer get-tss-address [bitcoinChainId]] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for get-tss-address
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module

## zetacored query observer list-blame

Query AllBlameRecords

```
zetacored query observer list-blame [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-blame
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module

## zetacored query observer list-blame-by-msg

Query AllBlameRecords

```
zetacored query observer list-blame-by-msg [chainId] [nonce] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-blame-by-msg
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module

## zetacored query observer list-chain-nonces

list all chainNonces

```
zetacored query observer list-chain-nonces [flags]
```

### Options

```
      --count-total        count total number of records in list-chain-nonces to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-chain-nonces
      --limit uint         pagination limit of list-chain-nonces to query for (default 100)
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
      --offset uint        pagination offset of list-chain-nonces to query for
  -o, --output string      Output format (text|json) 
      --page uint          pagination page of list-chain-nonces to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of list-chain-nonces to query for
      --reverse            results are sorted in descending order
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module

## zetacored query observer list-chain-params

Query GetChainParams

```
zetacored query observer list-chain-params [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-chain-params
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module

## zetacored query observer list-chains

list all SupportedChains

```
zetacored query observer list-chains [flags]
```

### Options

```
      --count-total        count total number of records in list-chains to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-chains
      --limit uint         pagination limit of list-chains to query for (default 100)
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
      --offset uint        pagination offset of list-chains to query for
  -o, --output string      Output format (text|json) 
      --page uint          pagination page of list-chains to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of list-chains to query for
      --reverse            results are sorted in descending order
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module

## zetacored query observer list-node-account

list all NodeAccount

```
zetacored query observer list-node-account [flags]
```

### Options

```
      --count-total        count total number of records in list-node-account to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-node-account
      --limit uint         pagination limit of list-node-account to query for (default 100)
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
      --offset uint        pagination offset of list-node-account to query for
  -o, --output string      Output format (text|json) 
      --page uint          pagination page of list-node-account to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of list-node-account to query for
      --reverse            results are sorted in descending order
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module

## zetacored query observer list-observer-set

Query observer set

```
zetacored query observer list-observer-set [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-observer-set
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module

## zetacored query observer list-pending-nonces

shows a chainNonces

```
zetacored query observer list-pending-nonces [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-pending-nonces
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module

## zetacored query observer list-tss-funds-migrator

list all tss funds migrators

```
zetacored query observer list-tss-funds-migrator [flags]
```

### Options

```
      --count-total        count total number of records in list-tss-funds-migrator to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-tss-funds-migrator
      --limit uint         pagination limit of list-tss-funds-migrator to query for (default 100)
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
      --offset uint        pagination offset of list-tss-funds-migrator to query for
  -o, --output string      Output format (text|json) 
      --page uint          pagination page of list-tss-funds-migrator to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of list-tss-funds-migrator to query for
      --reverse            results are sorted in descending order
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module

## zetacored query observer list-tss-history

show historical list of TSS

```
zetacored query observer list-tss-history [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-tss-history
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module

## zetacored query observer show-ballot

Query BallotByIdentifier

```
zetacored query observer show-ballot [ballot-identifier] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-ballot
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module

## zetacored query observer show-blame

Query BlameByIdentifier

```
zetacored query observer show-blame [blame-identifier] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-blame
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module

## zetacored query observer show-chain-nonces

shows a chainNonces

```
zetacored query observer show-chain-nonces [chain-id] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-chain-nonces
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module

## zetacored query observer show-chain-params

Query GetChainParamsForChain

```
zetacored query observer show-chain-params [chain-id] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-chain-params
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module

## zetacored query observer show-crosschain-flags

shows the crosschain flags

```
zetacored query observer show-crosschain-flags [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-crosschain-flags
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module

## zetacored query observer show-keygen

shows keygen

```
zetacored query observer show-keygen [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-keygen
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module

## zetacored query observer show-node-account

shows a NodeAccount

```
zetacored query observer show-node-account [operator_address] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-node-account
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module

## zetacored query observer show-observer-count

Query show-observer-count

```
zetacored query observer show-observer-count [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-observer-count
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module

## zetacored query observer show-operational-flags

shows the operational flags

```
zetacored query observer show-operational-flags [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-operational-flags
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module

## zetacored query observer show-tss

shows a TSS

```
zetacored query observer show-tss [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-tss
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module

## zetacored query observer show-tss-funds-migrator

show the tss funds migrator for a chain

```
zetacored query observer show-tss-funds-migrator [chain-id] [flags]
```

### Options

```
      --count-total        count total number of records in show-tss-funds-migrator [chain-id] to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for show-tss-funds-migrator
      --limit uint         pagination limit of show-tss-funds-migrator [chain-id] to query for (default 100)
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
      --offset uint        pagination offset of show-tss-funds-migrator [chain-id] to query for
  -o, --output string      Output format (text|json) 
      --page uint          pagination page of show-tss-funds-migrator [chain-id] to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of show-tss-funds-migrator [chain-id] to query for
      --reverse            results are sorted in descending order
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module

## zetacored query tx

Query for a transaction by hash, "[addr]/[seq]" combination or comma-separated signatures in a committed block

### Synopsis

Example:
$ zetacored query tx [hash]
$ zetacored query tx --type=acc_seq [addr]/[sequence]
$ zetacored query tx --type=signature [sig1_base64],[sig2_base64...]

```
zetacored query tx --type=[hash|acc_seq|signature] [hash|acc_seq|signature] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for tx
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
      --type string        The type to be used when querying tx, can be one of "hash", "acc_seq", "signature" 
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query](#zetacored-query)	 - Querying subcommands

## zetacored query txs

Query for paginated transactions that match a set of events

### Synopsis

Search for transactions that match the exact given events where results are paginated.
The events query is directly passed to Tendermint's RPC TxSearch method and must
conform to Tendermint's query syntax.

Please refer to each module's documentation for the full set of events to query
for. Each module documents its respective events under 'xx_events.md'.


```
zetacored query txs [flags]
```

### Examples

```
$ zetacored query txs --query "message.sender='cosmos1...' AND message.action='withdraw_delegator_reward' AND tx.height > 7" --page 1 --limit 30
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for txs
      --limit int          Query number of transactions results per page returned (default 100)
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
      --order_by string    The ordering semantics (asc|dsc)
  -o, --output string      Output format (text|json) 
      --page int           Query a specific page of paginated results (default 1)
      --query string       The transactions events query per Tendermint's query semantics
```

### Options inherited from parent commands

```
      --chain-id string     The network chain ID
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored query](#zetacored-query)	 - Querying subcommands

## zetacored rollback

rollback Cosmos SDK and CometBFT state by one height

### Synopsis


A state rollback is performed to recover from an incorrect application state transition,
when CometBFT has persisted an incorrect app hash and is thus unable to make
progress. Rollback overwrites a state at height n with the state at height n - 1.
The application also rolls back to height n - 1. No blocks are removed, so upon
restarting CometBFT the transactions in block n will be re-executed against the
application.


```
zetacored rollback [flags]
```

### Options

```
      --hard          remove last block as well as state
  -h, --help          help for rollback
      --home string   The application home directory 
```

### Options inherited from parent commands

```
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)

## zetacored snapshots

Manage local snapshots

### Options

```
  -h, --help   help for snapshots
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)
* [zetacored snapshots delete](#zetacored-snapshots-delete)	 - Delete a local snapshot
* [zetacored snapshots dump](#zetacored-snapshots-dump)	 - Dump the snapshot as portable archive format
* [zetacored snapshots export](#zetacored-snapshots-export)	 - Export app state to snapshot store
* [zetacored snapshots list](#zetacored-snapshots-list)	 - List local snapshots
* [zetacored snapshots load](#zetacored-snapshots-load)	 - Load a snapshot archive file (.tar.gz) into snapshot store
* [zetacored snapshots restore](#zetacored-snapshots-restore)	 - Restore app state from local snapshot

## zetacored snapshots delete

Delete a local snapshot

```
zetacored snapshots delete [height] [format] [flags]
```

### Options

```
  -h, --help   help for delete
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored snapshots](#zetacored-snapshots)	 - Manage local snapshots

## zetacored snapshots dump

Dump the snapshot as portable archive format

```
zetacored snapshots dump [height] [format] [flags]
```

### Options

```
  -h, --help            help for dump
  -o, --output string   output file
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored snapshots](#zetacored-snapshots)	 - Manage local snapshots

## zetacored snapshots export

Export app state to snapshot store

```
zetacored snapshots export [flags]
```

### Options

```
      --height int   Height to export, default to latest state height
  -h, --help         help for export
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored snapshots](#zetacored-snapshots)	 - Manage local snapshots

## zetacored snapshots list

List local snapshots

```
zetacored snapshots list [flags]
```

### Options

```
  -h, --help   help for list
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored snapshots](#zetacored-snapshots)	 - Manage local snapshots

## zetacored snapshots load

Load a snapshot archive file (.tar.gz) into snapshot store

```
zetacored snapshots load [archive-file] [flags]
```

### Options

```
  -h, --help   help for load
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored snapshots](#zetacored-snapshots)	 - Manage local snapshots

## zetacored snapshots restore

Restore app state from local snapshot

### Synopsis

Restore app state from local snapshot

```
zetacored snapshots restore [height] [format] [flags]
```

### Options

```
  -h, --help   help for restore
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored snapshots](#zetacored-snapshots)	 - Manage local snapshots

## zetacored start

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
      --genesis_hash bytesHex                           optional SHA-256 hash of the genesis file
      --grpc-only                                       Start the node in gRPC query only mode without Tendermint process
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
      --priv_validator_laddr string                     socket address to listen on for connections from external priv_validator process
      --proxy_app string                                proxy app address, or one of: 'kvstore', 'persistent_kvstore' or 'noop' for local testing. 
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
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)

## zetacored tendermint

Tendermint subcommands

### Options

```
  -h, --help   help for tendermint
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)
* [zetacored tendermint reset-state](#zetacored-tendermint-reset-state)	 - Remove all the data and WAL
* [zetacored tendermint show-address](#zetacored-tendermint-show-address)	 - Shows this node's CometBFT validator consensus address
* [zetacored tendermint show-node-id](#zetacored-tendermint-show-node-id)	 - Show this node's ID
* [zetacored tendermint show-validator](#zetacored-tendermint-show-validator)	 - Show this node's CometBFT validator info
* [zetacored tendermint unsafe-reset-all](#zetacored-tendermint-unsafe-reset-all)	 - (unsafe) Remove all the data and WAL, reset this node's validator to genesis state
* [zetacored tendermint version](#zetacored-tendermint-version)	 - Print CometBFT libraries' version

## zetacored tendermint reset-state

Remove all the data and WAL

```
zetacored tendermint reset-state [flags]
```

### Options

```
  -h, --help   help for reset-state
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored tendermint](#zetacored-tendermint)	 - Tendermint subcommands

## zetacored tendermint show-address

Shows this node's CometBFT validator consensus address

```
zetacored tendermint show-address [flags]
```

### Options

```
  -h, --help   help for show-address
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored tendermint](#zetacored-tendermint)	 - Tendermint subcommands

## zetacored tendermint show-node-id

Show this node's ID

```
zetacored tendermint show-node-id [flags]
```

### Options

```
  -h, --help   help for show-node-id
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored tendermint](#zetacored-tendermint)	 - Tendermint subcommands

## zetacored tendermint show-validator

Show this node's CometBFT validator info

```
zetacored tendermint show-validator [flags]
```

### Options

```
  -h, --help   help for show-validator
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored tendermint](#zetacored-tendermint)	 - Tendermint subcommands

## zetacored tendermint unsafe-reset-all

(unsafe) Remove all the data and WAL, reset this node's validator to genesis state

```
zetacored tendermint unsafe-reset-all [flags]
```

### Options

```
  -h, --help             help for unsafe-reset-all
      --keep-addr-book   keep the address book intact
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored tendermint](#zetacored-tendermint)	 - Tendermint subcommands

## zetacored tendermint version

Print CometBFT libraries' version

### Synopsis

Print protocols' and libraries' version numbers against which this app has been compiled.

```
zetacored tendermint version [flags]
```

### Options

```
  -h, --help   help for version
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored tendermint](#zetacored-tendermint)	 - Tendermint subcommands

## zetacored testnet

subcommands for starting or configuring local testnets

```
zetacored testnet [flags]
```

### Options

```
  -h, --help   help for testnet
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)
* [zetacored testnet init-files](#zetacored-testnet-init-files)	 - Initialize config directories & files for a multi-validator testnet running locally via separate processes (e.g. Docker Compose or similar)
* [zetacored testnet start](#zetacored-testnet-start)	 - Launch an in-process multi-validator testnet

## zetacored testnet init-files

Initialize config directories & files for a multi-validator testnet running locally via separate processes (e.g. Docker Compose or similar)

### Synopsis

init-files will setup "v" number of directories and populate each with
necessary files (private validator, genesis, config, etc.) for running "v" validator nodes.

Booting up a network with these validator folders is intended to be used with Docker Compose,
or a similar setup where each node has a manually configurable IP address.

Note, strict routability for addresses is turned off in the config file.

Example:
	evmosd testnet init-files --v 4 --output-dir ./.testnets --starting-ip-address 192.168.10.2
	

```
zetacored testnet init-files [flags]
```

### Options

```
      --chain-id string              genesis file chain-id, if left blank will be randomly created
  -h, --help                         help for init-files
      --key-type string              Key signing algorithm to generate keys for 
      --keyring-backend string       Select keyring's backend (os|file|test) 
      --minimum-gas-prices string    Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum (e.g. 0.01photino,0.001stake) 
      --node-daemon-home string      Home directory of the node's daemon configuration 
      --node-dir-prefix string       Prefix the directory name for each node with (node results in node0, node1, ...) 
  -o, --output-dir string            Directory to store initialization data for the testnet 
      --starting-ip-address string   Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:46656, ID1@192.168.0.2:46656, ...) 
      --v int                        Number of validators to initialize the testnet with (default 4)
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored testnet](#zetacored-testnet)	 - subcommands for starting or configuring local testnets

## zetacored testnet start

Launch an in-process multi-validator testnet

### Synopsis

testnet will launch an in-process multi-validator testnet,
and generate "v" directories, populated with necessary validator configuration files
(private validator, genesis, config, etc.).

Example:
	evmosd testnet --v 4 --output-dir ./.testnets
	

```
zetacored testnet start [flags]
```

### Options

```
      --api.address string          the address to listen on for REST API 
      --chain-id string             genesis file chain-id, if left blank will be randomly created
      --enable-logging              Enable INFO logging of tendermint validator nodes
      --grpc.address string         the gRPC server address to listen on 
  -h, --help                        help for start
      --json-rpc.address string     the JSON-RPC server address to listen on 
      --key-type string             Key signing algorithm to generate keys for 
      --minimum-gas-prices string   Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum (e.g. 0.01photino,0.001stake) 
  -o, --output-dir string           Directory to store initialization data for the testnet 
      --print-mnemonic              print mnemonic of first validator to stdout for manual testing (default true)
      --rpc.address string          the RPC address to listen on 
      --v int                       Number of validators to initialize the testnet with (default 4)
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored testnet](#zetacored-testnet)	 - subcommands for starting or configuring local testnets

## zetacored tx

Transactions subcommands

```
zetacored tx [flags]
```

### Options

```
      --chain-id string   The network chain ID
  -h, --help              help for tx
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)
* [zetacored tx broadcast](#zetacored-tx-broadcast)	 - Broadcast transactions generated offline
* [zetacored tx decode](#zetacored-tx-decode)	 - Decode a binary encoded transaction string
* [zetacored tx encode](#zetacored-tx-encode)	 - Encode transactions generated offline
* [zetacored tx multi-sign](#zetacored-tx-multi-sign)	 - Generate multisig signatures for transactions generated offline
* [zetacored tx multisign-batch](#zetacored-tx-multisign-batch)	 - Assemble multisig transactions in batch from batch signatures
* [zetacored tx sign](#zetacored-tx-sign)	 - Sign a transaction generated offline
* [zetacored tx sign-batch](#zetacored-tx-sign-batch)	 - Sign transaction batch files
* [zetacored tx validate-signatures](#zetacored-tx-validate-signatures)	 - validate transactions signatures

## zetacored tx broadcast

Broadcast transactions generated offline

### Synopsis

Broadcast transactions created with the --generate-only
flag and signed with the sign command. Read a transaction from [file_path] and
broadcast it to a node. If you supply a dash (-) argument in place of an input
filename, the command reads from standard input.

$ zetacored tx broadcast ./mytxn.json

```
zetacored tx broadcast [file_path] [flags]
```

### Options

```
  -a, --account-number uint      The account number of the signing account (offline mode only)
      --aux                      Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string    Transaction broadcasting mode (sync|async) 
      --chain-id string          The network chain ID
      --dry-run                  ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string       Fee granter grants fees for the transaction
      --fee-payer string         Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string              Fees to pay along with transaction; eg: 10uatom
      --from string              Name or address of private key with which to sign
      --gas string               gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float     adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string        Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only            Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                     help for broadcast
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                   Use a connected Ledger device
      --node string              [host]:[port] to CometBFT rpc interface for this chain 
      --note string              Note to add a description to the transaction (previously --memo)
      --offline                  Offline mode (does not allow any online functionality)
  -o, --output string            Output format (text|json) 
  -s, --sequence uint            The sequence number of the signing account (offline mode only)
      --sign-mode string         Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-height uint      Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string               Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
  -y, --yes                      Skip tx broadcasting prompt confirmation
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands

## zetacored tx decode

Decode a binary encoded transaction string

```
zetacored tx decode [protobuf-byte-string] [flags]
```

### Options

```
  -a, --account-number uint      The account number of the signing account (offline mode only)
      --aux                      Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string    Transaction broadcasting mode (sync|async) 
      --chain-id string          The network chain ID
      --dry-run                  ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string       Fee granter grants fees for the transaction
      --fee-payer string         Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string              Fees to pay along with transaction; eg: 10uatom
      --from string              Name or address of private key with which to sign
      --gas string               gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float     adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string        Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only            Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                     help for decode
  -x, --hex                      Treat input as hexadecimal instead of base64
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                   Use a connected Ledger device
      --node string              [host]:[port] to CometBFT rpc interface for this chain 
      --note string              Note to add a description to the transaction (previously --memo)
      --offline                  Offline mode (does not allow any online functionality)
  -s, --sequence uint            The sequence number of the signing account (offline mode only)
      --sign-mode string         Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-height uint      Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string               Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
  -y, --yes                      Skip tx broadcasting prompt confirmation
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands

## zetacored tx encode

Encode transactions generated offline

### Synopsis

Encode transactions created with the --generate-only flag or signed with the sign command.
Read a transaction from [file], serialize it to the Protobuf wire protocol, and output it as base64.
If you supply a dash (-) argument in place of an input filename, the command reads from standard input.

```
zetacored tx encode [file] [flags]
```

### Options

```
  -a, --account-number uint      The account number of the signing account (offline mode only)
      --aux                      Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string    Transaction broadcasting mode (sync|async) 
      --chain-id string          The network chain ID
      --dry-run                  ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string       Fee granter grants fees for the transaction
      --fee-payer string         Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string              Fees to pay along with transaction; eg: 10uatom
      --from string              Name or address of private key with which to sign
      --gas string               gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float     adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string        Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only            Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                     help for encode
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                   Use a connected Ledger device
      --node string              [host]:[port] to CometBFT rpc interface for this chain 
      --note string              Note to add a description to the transaction (previously --memo)
      --offline                  Offline mode (does not allow any online functionality)
  -s, --sequence uint            The sequence number of the signing account (offline mode only)
      --sign-mode string         Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-height uint      Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string               Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
  -y, --yes                      Skip tx broadcasting prompt confirmation
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands

## zetacored tx multi-sign

Generate multisig signatures for transactions generated offline

### Synopsis

Sign transactions created with the --generate-only flag that require multisig signatures.

Read one or more signatures from one or more [signature] file, generate a multisig signature compliant to the
multisig key [name], and attach the key name to the transaction read from [file].

Example:
$ zetacored tx multisign transaction.json k1k2k3 k1sig.json k2sig.json k3sig.json

If --signature-only flag is on, output a JSON representation
of only the generated signature.

If the --offline flag is on, the client will not reach out to an external node.
Account number or sequence number lookups are not performed so you must
set these parameters manually.

If the --skip-signature-verification flag is on, the command will not verify the
signatures in the provided signature files. This is useful when the multisig
account is a signer in a nested multisig scenario.

The current multisig implementation defaults to amino-json sign mode.
The SIGN_MODE_DIRECT sign mode is not supported.'

```
zetacored tx multi-sign [file] [name] [[signature]...] [flags]
```

### Options

```
  -a, --account-number uint           The account number of the signing account (offline mode only)
      --aux                           Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string         Transaction broadcasting mode (sync|async) 
      --chain-id string               The network chain ID
      --dry-run                       ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string            Fee granter grants fees for the transaction
      --fee-payer string              Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                   Fees to pay along with transaction; eg: 10uatom
      --from string                   Name or address of private key with which to sign
      --gas string                    gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float          adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string             Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only                 Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                          help for multi-sign
      --keyring-backend string        Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string            The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                        Use a connected Ledger device
      --node string                   [host]:[port] to CometBFT rpc interface for this chain 
      --note string                   Note to add a description to the transaction (previously --memo)
      --offline                       Offline mode (does not allow any online functionality)
      --output-document string        The document is written to the given file instead of STDOUT
  -s, --sequence uint                 The sequence number of the signing account (offline mode only)
      --sign-mode string              Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --signature-only                Print only the generated signature, then exit
      --skip-signature-verification   Skip signature verification
      --timeout-height uint           Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                    Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
  -y, --yes                           Skip tx broadcasting prompt confirmation
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands

## zetacored tx multisign-batch

Assemble multisig transactions in batch from batch signatures

### Synopsis

Assemble a batch of multisig transactions generated by batch sign command.

Read one or more signatures from one or more [signature] file, generate a multisig signature compliant to the
multisig key [name], and attach the key name to the transaction read from [file].

Example:
$ zetacored tx multisign-batch transactions.json multisigk1k2k3 k1sigs.json k2sigs.json k3sig.json

The current multisig implementation defaults to amino-json sign mode.
The SIGN_MODE_DIRECT sign mode is not supported.'

```
zetacored tx multisign-batch [file] [name] [[signature-file]...] [flags]
```

### Options

```
  -a, --account-number uint      The account number of the signing account (offline mode only)
      --aux                      Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string    Transaction broadcasting mode (sync|async) 
      --chain-id string          The network chain ID
      --dry-run                  ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string       Fee granter grants fees for the transaction
      --fee-payer string         Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string              Fees to pay along with transaction; eg: 10uatom
      --from string              Name or address of private key with which to sign
      --gas string               gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float     adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string        Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only            Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                     help for multisign-batch
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                   Use a connected Ledger device
      --multisig string          Address of the multisig account that the transaction signs on behalf of
      --no-auto-increment        disable sequence auto increment
      --node string              [host]:[port] to CometBFT rpc interface for this chain 
      --note string              Note to add a description to the transaction (previously --memo)
      --offline                  Offline mode (does not allow any online functionality)
      --output-document string   The document is written to the given file instead of STDOUT
  -s, --sequence uint            The sequence number of the signing account (offline mode only)
      --sign-mode string         Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-height uint      Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string               Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
  -y, --yes                      Skip tx broadcasting prompt confirmation
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands

## zetacored tx sign

Sign a transaction generated offline

### Synopsis

Sign a transaction created with the --generate-only flag.
It will read a transaction from [file], sign it, and print its JSON encoding.

If the --signature-only flag is set, it will output the signature parts only.

The --offline flag makes sure that the client will not reach out to full node.
As a result, the account and sequence number queries will not be performed and
it is required to set such parameters manually. Note, invalid values will cause
the transaction to fail.

The --multisig=[multisig_key] flag generates a signature on behalf of a multisig account
key. It implies --signature-only. Full multisig signed transactions may eventually
be generated via the 'multisign' command.


```
zetacored tx sign [file] [flags]
```

### Options

```
  -a, --account-number uint      The account number of the signing account (offline mode only)
      --aux                      Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string    Transaction broadcasting mode (sync|async) 
      --chain-id string          The network chain ID
      --dry-run                  ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string       Fee granter grants fees for the transaction
      --fee-payer string         Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string              Fees to pay along with transaction; eg: 10uatom
      --from string              Name or address of private key with which to sign
      --gas string               gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float     adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string        Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only            Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                     help for sign
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                   Use a connected Ledger device
      --multisig string          Address or key name of the multisig account on behalf of which the transaction shall be signed
      --node string              [host]:[port] to CometBFT rpc interface for this chain 
      --note string              Note to add a description to the transaction (previously --memo)
      --offline                  Offline mode (does not allow any online functionality)
  -o, --output string            Output format (text|json) 
      --output-document string   The document will be written to the given file instead of STDOUT
      --overwrite                Overwrite existing signatures with a new one. If disabled, new signature will be appended
  -s, --sequence uint            The sequence number of the signing account (offline mode only)
      --sign-mode string         Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --signature-only           Print only the signatures
      --timeout-height uint      Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string               Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
  -y, --yes                      Skip tx broadcasting prompt confirmation
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands

## zetacored tx sign-batch

Sign transaction batch files

### Synopsis

Sign batch files of transactions generated with --generate-only.
The command processes list of transactions from a file (one StdTx each line), or multiple files.
Then generates signed transactions or signatures and print their JSON encoding, delimited by '\n'.
As the signatures are generated, the command updates the account and sequence number accordingly.

If the --signature-only flag is set, it will output the signature parts only.

The --offline flag makes sure that the client will not reach out to full node.
As a result, the account and the sequence number queries will not be performed and
it is required to set such parameters manually. Note, invalid values will cause
the transaction to fail. The sequence will be incremented automatically for each
transaction that is signed.

If --account-number or --sequence flag is used when offline=false, they are ignored and 
overwritten by the default flag values.

The --multisig=[multisig_key] flag generates a signature on behalf of a multisig
account key. It implies --signature-only.


```
zetacored tx sign-batch [file] ([file2]...) [flags]
```

### Options

```
  -a, --account-number uint      The account number of the signing account (offline mode only)
      --append                   Combine all message and generate single signed transaction for broadcast.
      --aux                      Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string    Transaction broadcasting mode (sync|async) 
      --chain-id string          The network chain ID
      --dry-run                  ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string       Fee granter grants fees for the transaction
      --fee-payer string         Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string              Fees to pay along with transaction; eg: 10uatom
      --from string              Name or address of private key with which to sign
      --gas string               gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float     adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string        Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only            Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                     help for sign-batch
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                   Use a connected Ledger device
      --multisig string          Address or key name of the multisig account on behalf of which the transaction shall be signed
      --node string              [host]:[port] to CometBFT rpc interface for this chain 
      --note string              Note to add a description to the transaction (previously --memo)
      --offline                  Offline mode (does not allow any online functionality)
  -o, --output string            Output format (text|json) 
      --output-document string   The document will be written to the given file instead of STDOUT
  -s, --sequence uint            The sequence number of the signing account (offline mode only)
      --sign-mode string         Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --signature-only           Print only the generated signature, then exit
      --timeout-height uint      Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string               Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
  -y, --yes                      Skip tx broadcasting prompt confirmation
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands

## zetacored tx validate-signatures

validate transactions signatures

### Synopsis

Print the addresses that must sign the transaction, those who have already
signed it, and make sure that signatures are in the correct order.

The command would check whether all required signers have signed the transactions, whether
the signatures were collected in the right order, and if the signature is valid over the
given transaction. If the --offline flag is also set, signature validation over the
transaction will be not be performed as that will require RPC communication with a full node.


```
zetacored tx validate-signatures [file] [flags]
```

### Options

```
  -a, --account-number uint      The account number of the signing account (offline mode only)
      --aux                      Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string    Transaction broadcasting mode (sync|async) 
      --chain-id string          The network chain ID
      --dry-run                  ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string       Fee granter grants fees for the transaction
      --fee-payer string         Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string              Fees to pay along with transaction; eg: 10uatom
      --from string              Name or address of private key with which to sign
      --gas string               gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float     adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string        Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only            Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                     help for validate-signatures
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                   Use a connected Ledger device
      --node string              [host]:[port] to CometBFT rpc interface for this chain 
      --note string              Note to add a description to the transaction (previously --memo)
      --offline                  Offline mode (does not allow any online functionality)
  -o, --output string            Output format (text|json) 
  -s, --sequence uint            The sequence number of the signing account (offline mode only)
      --sign-mode string         Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-height uint      Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string               Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
  -y, --yes                      Skip tx broadcasting prompt confirmation
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands

## zetacored upgrade-handler-version

Print the default upgrade handler version

```
zetacored upgrade-handler-version [flags]
```

### Options

```
  -h, --help   help for upgrade-handler-version
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)

## zetacored validate

Validates the genesis file at the default location or at the location passed as an arg

```
zetacored validate [file] [flags]
```

### Options

```
  -h, --help   help for validate
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)

## zetacored version

Print the application binary version information

```
zetacored version [flags]
```

### Options

```
  -h, --help            help for version
      --long            Print long version information
  -o, --output string   Output format (text|json) 
```

### Options inherited from parent commands

```
      --home string         directory for config and data 
      --log_format string   The logging format (json|plain) 
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:[level],[key]:[level]') 
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors
```

### SEE ALSO

* [zetacored](#zetacored)	 - Zetacore Daemon (server)

