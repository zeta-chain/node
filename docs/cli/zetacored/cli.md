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

* [zetacored comet](#zetacored-comet)	 - CometBFT subcommands
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
* [zetacored status](#zetacored-status)	 - Query remote node for status
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

## zetacored comet

CometBFT subcommands

### Options

```
  -h, --help   help for comet
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
* [zetacored comet bootstrap-state](#zetacored-comet-bootstrap-state)	 - Bootstrap CometBFT state at an arbitrary block height using a light client
* [zetacored comet reset-state](#zetacored-comet-reset-state)	 - Remove all the data and WAL
* [zetacored comet show-address](#zetacored-comet-show-address)	 - Shows this node's CometBFT validator consensus address
* [zetacored comet show-node-id](#zetacored-comet-show-node-id)	 - Show this node's ID
* [zetacored comet show-validator](#zetacored-comet-show-validator)	 - Show this node's CometBFT validator info
* [zetacored comet unsafe-reset-all](#zetacored-comet-unsafe-reset-all)	 - (unsafe) Remove all the data and WAL, reset this node's validator to genesis state
* [zetacored comet version](#zetacored-comet-version)	 - Print CometBFT libraries' version

## zetacored comet bootstrap-state

Bootstrap CometBFT state at an arbitrary block height using a light client

```
zetacored comet bootstrap-state [flags]
```

### Options

```
      --height int   Block height to bootstrap state at, if not provided it uses the latest block height in app state
  -h, --help         help for bootstrap-state
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

* [zetacored comet](#zetacored-comet)	 - CometBFT subcommands

## zetacored comet reset-state

Remove all the data and WAL

```
zetacored comet reset-state [flags]
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

* [zetacored comet](#zetacored-comet)	 - CometBFT subcommands

## zetacored comet show-address

Shows this node's CometBFT validator consensus address

```
zetacored comet show-address [flags]
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

* [zetacored comet](#zetacored-comet)	 - CometBFT subcommands

## zetacored comet show-node-id

Show this node's ID

```
zetacored comet show-node-id [flags]
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

* [zetacored comet](#zetacored-comet)	 - CometBFT subcommands

## zetacored comet show-validator

Show this node's CometBFT validator info

```
zetacored comet show-validator [flags]
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

* [zetacored comet](#zetacored-comet)	 - CometBFT subcommands

## zetacored comet unsafe-reset-all

(unsafe) Remove all the data and WAL, reset this node's validator to genesis state

```
zetacored comet unsafe-reset-all [flags]
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

* [zetacored comet](#zetacored-comet)	 - CometBFT subcommands

## zetacored comet version

Print CometBFT libraries' version

### Synopsis

Print protocols' and libraries' version numbers against which this app has been compiled.

```
zetacored comet version [flags]
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

* [zetacored comet](#zetacored-comet)	 - CometBFT subcommands

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
      --timeout-duration duration           TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint                 DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                          Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                           Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
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
      --chain-id string             genesis file chain-id, if left blank will be randomly created
      --consensus-key-algo string   algorithm to use for the consensus key 
      --default-denom string        genesis file default denomination, if left blank default value is 'stake'
  -h, --help                        help for init
      --home string                 node's home directory 
      --initial-height int          specify the initial block height at genesis (default 1)
  -o, --overwrite                   overwrite the genesis.json file
      --recover                     provide seed phrase to recover existing key instead of creating
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
CometBFT crypto library and can be used by light-clients, full nodes, or any other application
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
* [zetacored keys list-key-types](#zetacored-keys-list-key-types)	 - List all key types
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
  -y, --yes             Skip confirmation prompt when export unarmored hex privkey
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

## zetacored keys list-key-types

List all key types

### Synopsis

Return a list of all supported key types (also known as algos)

```
zetacored keys list-key-types [flags]
```

### Options

```
  -h, --help   help for list-key-types
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
  -y, --yes              Skip confirmation prompt when check input entropy length
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
      --qrcode                   Display key address QR code (will be ignored if -a or --address is false)
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
* [zetacored query auth](#zetacored-query-auth)	 - Querying commands for the auth module
* [zetacored query authority](#zetacored-query-authority)	 - Querying commands for the authority module
* [zetacored query authz](#zetacored-query-authz)	 - Querying commands for the authz module
* [zetacored query bank](#zetacored-query-bank)	 - Querying commands for the bank module
* [zetacored query block](#zetacored-query-block)	 - Query for a committed block by height, hash, or event(s)
* [zetacored query block-results](#zetacored-query-block-results)	 - Query for a committed block's results by height
* [zetacored query blocks](#zetacored-query-blocks)	 - Query for paginated blocks that match a set of events
* [zetacored query comet-validator-set](#zetacored-query-comet-validator-set)	 - Get the full CometBFT validator set at given height
* [zetacored query consensus](#zetacored-query-consensus)	 - Querying commands for the consensus module
* [zetacored query crosschain](#zetacored-query-crosschain)	 - Querying commands for the crosschain module
* [zetacored query distribution](#zetacored-query-distribution)	 - Querying commands for the distribution module
* [zetacored query emissions](#zetacored-query-emissions)	 - Querying commands for the emissions module
* [zetacored query evidence](#zetacored-query-evidence)	 - Querying commands for the evidence module
* [zetacored query evm](#zetacored-query-evm)	 - Querying commands for the evm module
* [zetacored query feemarket](#zetacored-query-feemarket)	 - Querying commands for the fee market module
* [zetacored query fungible](#zetacored-query-fungible)	 - Querying commands for the fungible module
* [zetacored query gov](#zetacored-query-gov)	 - Querying commands for the gov module
* [zetacored query group](#zetacored-query-group)	 - Querying commands for the group module
* [zetacored query lightclient](#zetacored-query-lightclient)	 - Querying commands for the lightclient module
* [zetacored query observer](#zetacored-query-observer)	 - Querying commands for the observer module
* [zetacored query params](#zetacored-query-params)	 - Querying commands for the params module
* [zetacored query slashing](#zetacored-query-slashing)	 - Querying commands for the slashing module
* [zetacored query staking](#zetacored-query-staking)	 - Querying commands for the staking module
* [zetacored query tx](#zetacored-query-tx)	 - Query for a transaction by hash, "[addr]/[seq]" combination or comma-separated signatures in a committed block
* [zetacored query txs](#zetacored-query-txs)	 - Query for paginated transactions that match a set of events
* [zetacored query upgrade](#zetacored-query-upgrade)	 - Querying commands for the upgrade module

## zetacored query auth

Querying commands for the auth module

```
zetacored query auth [flags]
```

### Options

```
  -h, --help   help for auth
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
* [zetacored query auth account](#zetacored-query-auth-account)	 - Query account by address
* [zetacored query auth account-info](#zetacored-query-auth-account-info)	 - Query account info which is common to all account types.
* [zetacored query auth accounts](#zetacored-query-auth-accounts)	 - Query all the accounts
* [zetacored query auth address-by-acc-num](#zetacored-query-auth-address-by-acc-num)	 - Query account address by account number
* [zetacored query auth address-bytes-to-string](#zetacored-query-auth-address-bytes-to-string)	 - Transform an address bytes to string
* [zetacored query auth address-string-to-bytes](#zetacored-query-auth-address-string-to-bytes)	 - Transform an address string to bytes
* [zetacored query auth bech32-prefix](#zetacored-query-auth-bech32-prefix)	 - Query the chain bech32 prefix (if applicable)
* [zetacored query auth module-account](#zetacored-query-auth-module-account)	 - Query module account info by module name
* [zetacored query auth module-accounts](#zetacored-query-auth-module-accounts)	 - Query all module accounts
* [zetacored query auth params](#zetacored-query-auth-params)	 - Query the current auth parameters

## zetacored query auth account

Query account by address

```
zetacored query auth account [address] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for account
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query auth](#zetacored-query-auth)	 - Querying commands for the auth module

## zetacored query auth account-info

Query account info which is common to all account types.

```
zetacored query auth account-info [address] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for account-info
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query auth](#zetacored-query-auth)	 - Querying commands for the auth module

## zetacored query auth accounts

Query all the accounts

```
zetacored query auth accounts [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for accounts
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query auth](#zetacored-query-auth)	 - Querying commands for the auth module

## zetacored query auth address-by-acc-num

Query account address by account number

```
zetacored query auth address-by-acc-num [acc-num] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for address-by-acc-num
      --id int                   
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query auth](#zetacored-query-auth)	 - Querying commands for the auth module

## zetacored query auth address-bytes-to-string

Transform an address bytes to string

```
zetacored query auth address-bytes-to-string [address-bytes] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for address-bytes-to-string
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query auth](#zetacored-query-auth)	 - Querying commands for the auth module

## zetacored query auth address-string-to-bytes

Transform an address string to bytes

```
zetacored query auth address-string-to-bytes [address-string] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for address-string-to-bytes
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query auth](#zetacored-query-auth)	 - Querying commands for the auth module

## zetacored query auth bech32-prefix

Query the chain bech32 prefix (if applicable)

```
zetacored query auth bech32-prefix [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for bech32-prefix
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query auth](#zetacored-query-auth)	 - Querying commands for the auth module

## zetacored query auth module-account

Query module account info by module name

```
zetacored query auth module-account [module-name] [flags]
```

### Examples

```
zetacored q auth module-account gov
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for module-account
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query auth](#zetacored-query-auth)	 - Querying commands for the auth module

## zetacored query auth module-accounts

Query all module accounts

```
zetacored query auth module-accounts [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for module-accounts
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query auth](#zetacored-query-auth)	 - Querying commands for the auth module

## zetacored query auth params

Query the current auth parameters

```
zetacored query auth params [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for params
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query auth](#zetacored-query-auth)	 - Querying commands for the auth module

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

## zetacored query authz

Querying commands for the authz module

```
zetacored query authz [flags]
```

### Options

```
  -h, --help   help for authz
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
* [zetacored query authz grants](#zetacored-query-authz-grants)	 - Query grants for a granter-grantee pair and optionally a msg-type-url
* [zetacored query authz grants-by-grantee](#zetacored-query-authz-grants-by-grantee)	 - Query authorization grants granted to a grantee
* [zetacored query authz grants-by-granter](#zetacored-query-authz-grants-by-granter)	 - Query authorization grants granted by granter

## zetacored query authz grants

Query grants for a granter-grantee pair and optionally a msg-type-url

### Synopsis

Query authorization grants for a granter-grantee pair. If msg-type-url is set, it will select grants only for that msg type.

```
zetacored query authz grants [granter-addr] [grantee-addr] [msg-type-url] [flags]
```

### Examples

```
zetacored query authz grants cosmos1skj.. cosmos1skjwj.. /cosmos.bank.v1beta1.MsgSend
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for grants
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query authz](#zetacored-query-authz)	 - Querying commands for the authz module

## zetacored query authz grants-by-grantee

Query authorization grants granted to a grantee

```
zetacored query authz grants-by-grantee [grantee-addr] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for grants-by-grantee
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query authz](#zetacored-query-authz)	 - Querying commands for the authz module

## zetacored query authz grants-by-granter

Query authorization grants granted by granter

```
zetacored query authz grants-by-granter [granter-addr] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for grants-by-granter
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query authz](#zetacored-query-authz)	 - Querying commands for the authz module

## zetacored query bank

Querying commands for the bank module

```
zetacored query bank [flags]
```

### Options

```
  -h, --help   help for bank
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
* [zetacored query bank balance](#zetacored-query-bank-balance)	 - Query an account balance by address and denom
* [zetacored query bank balances](#zetacored-query-bank-balances)	 - Query for account balances by address
* [zetacored query bank denom-metadata](#zetacored-query-bank-denom-metadata)	 - Query the client metadata of a given coin denomination
* [zetacored query bank denom-metadata-by-query-string](#zetacored-query-bank-denom-metadata-by-query-string)	 - Execute the DenomMetadataByQueryString RPC method
* [zetacored query bank denom-owners](#zetacored-query-bank-denom-owners)	 - Query for all account addresses that own a particular token denomination.
* [zetacored query bank denom-owners-by-query](#zetacored-query-bank-denom-owners-by-query)	 - Execute the DenomOwnersByQuery RPC method
* [zetacored query bank denoms-metadata](#zetacored-query-bank-denoms-metadata)	 - Query the client metadata for all registered coin denominations
* [zetacored query bank params](#zetacored-query-bank-params)	 - Query the current bank parameters
* [zetacored query bank send-enabled](#zetacored-query-bank-send-enabled)	 - Query for send enabled entries
* [zetacored query bank spendable-balance](#zetacored-query-bank-spendable-balance)	 - Query the spendable balance of a single denom for a single account.
* [zetacored query bank spendable-balances](#zetacored-query-bank-spendable-balances)	 - Query for account spendable balances by address
* [zetacored query bank total-supply](#zetacored-query-bank-total-supply)	 - Query the total supply of coins of the chain
* [zetacored query bank total-supply-of](#zetacored-query-bank-total-supply-of)	 - Query the supply of a single coin denom

## zetacored query bank balance

Query an account balance by address and denom

```
zetacored query bank balance [address] [denom] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for balance
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query bank](#zetacored-query-bank)	 - Querying commands for the bank module

## zetacored query bank balances

Query for account balances by address

### Synopsis

Query the total balance of an account or of a specific denomination.

```
zetacored query bank balances [address] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for balances
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
      --resolve-denom            
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

* [zetacored query bank](#zetacored-query-bank)	 - Querying commands for the bank module

## zetacored query bank denom-metadata

Query the client metadata of a given coin denomination

```
zetacored query bank denom-metadata [denom] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for denom-metadata
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query bank](#zetacored-query-bank)	 - Querying commands for the bank module

## zetacored query bank denom-metadata-by-query-string

Execute the DenomMetadataByQueryString RPC method

```
zetacored query bank denom-metadata-by-query-string [flags]
```

### Options

```
      --denom string             
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for denom-metadata-by-query-string
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query bank](#zetacored-query-bank)	 - Querying commands for the bank module

## zetacored query bank denom-owners

Query for all account addresses that own a particular token denomination.

```
zetacored query bank denom-owners [denom] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for denom-owners
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query bank](#zetacored-query-bank)	 - Querying commands for the bank module

## zetacored query bank denom-owners-by-query

Execute the DenomOwnersByQuery RPC method

```
zetacored query bank denom-owners-by-query [flags]
```

### Options

```
      --denom string             
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for denom-owners-by-query
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query bank](#zetacored-query-bank)	 - Querying commands for the bank module

## zetacored query bank denoms-metadata

Query the client metadata for all registered coin denominations

```
zetacored query bank denoms-metadata [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for denoms-metadata
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query bank](#zetacored-query-bank)	 - Querying commands for the bank module

## zetacored query bank params

Query the current bank parameters

```
zetacored query bank params [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for params
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query bank](#zetacored-query-bank)	 - Querying commands for the bank module

## zetacored query bank send-enabled

Query for send enabled entries

### Synopsis

Query for send enabled entries that have been specifically set.
			
To look up one or more specific denoms, supply them as arguments to this command.
To look up all denoms, do not provide any arguments.

```
zetacored query bank send-enabled [denom1 ...] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for send-enabled
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query bank](#zetacored-query-bank)	 - Querying commands for the bank module

## zetacored query bank spendable-balance

Query the spendable balance of a single denom for a single account.

```
zetacored query bank spendable-balance [address] [denom] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for spendable-balance
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query bank](#zetacored-query-bank)	 - Querying commands for the bank module

## zetacored query bank spendable-balances

Query for account spendable balances by address

```
zetacored query bank spendable-balances [address] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for spendable-balances
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query bank](#zetacored-query-bank)	 - Querying commands for the bank module

## zetacored query bank total-supply

Query the total supply of coins of the chain

### Synopsis

Query total supply of coins that are held by accounts in the chain. To query for the total supply of a specific coin denomination use --denom flag.

```
zetacored query bank total-supply [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for total-supply
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query bank](#zetacored-query-bank)	 - Querying commands for the bank module

## zetacored query bank total-supply-of

Query the supply of a single coin denom

```
zetacored query bank total-supply-of [denom] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for total-supply-of
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query bank](#zetacored-query-bank)	 - Querying commands for the bank module

## zetacored query block

Query for a committed block by height, hash, or event(s)

### Synopsis

Query for a specific committed block using the CometBFT RPC `block` and `block_by_hash` method

```
zetacored query block --type=[height|hash] [height|hash] [flags]
```

### Examples

```
$ zetacored query block --type=height [height]
$ zetacored query block --type=hash [hash]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for block
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string      Output format (text|json) 
      --type string        The type to be used when querying tx, can be one of "height", "hash" 
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

## zetacored query block-results

Query for a committed block's results by height

### Synopsis

Query for a specific committed block's results using the CometBFT RPC `block_results` method

```
zetacored query block-results [height] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for block-results
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

* [zetacored query](#zetacored-query)	 - Querying subcommands

## zetacored query blocks

Query for paginated blocks that match a set of events

### Synopsis

Search for blocks that match the exact given events where results are paginated.
The events query is directly passed to CometBFT's RPC BlockSearch method and must
conform to CometBFT's query syntax.
Please refer to each module's documentation for the full set of events to query
for. Each module documents its respective events under 'xx_events.md'.


```
zetacored query blocks [flags]
```

### Examples

```
$ zetacored query blocks --query "message.sender='cosmos1...' AND block.height > 7" --page 1 --limit 30 --order_by asc
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for blocks
      --limit int          Query number of transactions results per page returned (default 100)
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
      --order_by string    The ordering semantics (asc|dsc)
  -o, --output string      Output format (text|json) 
      --page int           Query a specific page of paginated results (default 1)
      --query string       The blocks events query per CometBFT's query semantics
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

## zetacored query consensus

Querying commands for the consensus module

```
zetacored query consensus [flags]
```

### Options

```
  -h, --help   help for consensus
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
* [zetacored query consensus comet](#zetacored-query-consensus-comet)	 - Querying commands for the cosmos.base.tendermint.v1beta1.Service service
* [zetacored query consensus params](#zetacored-query-consensus-params)	 - Query the current consensus parameters

## zetacored query consensus comet

Querying commands for the cosmos.base.tendermint.v1beta1.Service service

```
zetacored query consensus comet [flags]
```

### Options

```
  -h, --help   help for comet
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

* [zetacored query consensus](#zetacored-query-consensus)	 - Querying commands for the consensus module
* [zetacored query consensus comet block-by-height](#zetacored-query-consensus-comet-block-by-height)	 - Query for a committed block by height
* [zetacored query consensus comet block-latest](#zetacored-query-consensus-comet-block-latest)	 - Query for the latest committed block
* [zetacored query consensus comet node-info](#zetacored-query-consensus-comet-node-info)	 - Query the current node info
* [zetacored query consensus comet syncing](#zetacored-query-consensus-comet-syncing)	 - Query node syncing status
* [zetacored query consensus comet validator-set](#zetacored-query-consensus-comet-validator-set)	 - Query for the latest validator set
* [zetacored query consensus comet validator-set-by-height](#zetacored-query-consensus-comet-validator-set-by-height)	 - Query for a validator set by height

## zetacored query consensus comet block-by-height

Query for a committed block by height

### Synopsis

Query for a specific committed block using the CometBFT RPC `block_by_height` method

```
zetacored query consensus comet block-by-height [height] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for block-by-height
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query consensus comet](#zetacored-query-consensus-comet)	 - Querying commands for the cosmos.base.tendermint.v1beta1.Service service

## zetacored query consensus comet block-latest

Query for the latest committed block

```
zetacored query consensus comet block-latest [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for block-latest
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query consensus comet](#zetacored-query-consensus-comet)	 - Querying commands for the cosmos.base.tendermint.v1beta1.Service service

## zetacored query consensus comet node-info

Query the current node info

```
zetacored query consensus comet node-info [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for node-info
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query consensus comet](#zetacored-query-consensus-comet)	 - Querying commands for the cosmos.base.tendermint.v1beta1.Service service

## zetacored query consensus comet syncing

Query node syncing status

```
zetacored query consensus comet syncing [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for syncing
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query consensus comet](#zetacored-query-consensus-comet)	 - Querying commands for the cosmos.base.tendermint.v1beta1.Service service

## zetacored query consensus comet validator-set

Query for the latest validator set

```
zetacored query consensus comet validator-set [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for validator-set
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query consensus comet](#zetacored-query-consensus-comet)	 - Querying commands for the cosmos.base.tendermint.v1beta1.Service service

## zetacored query consensus comet validator-set-by-height

Query for a validator set by height

```
zetacored query consensus comet validator-set-by-height [height] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for validator-set-by-height
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query consensus comet](#zetacored-query-consensus-comet)	 - Querying commands for the cosmos.base.tendermint.v1beta1.Service service

## zetacored query consensus params

Query the current consensus parameters

```
zetacored query consensus params [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for params
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query consensus](#zetacored-query-consensus)	 - Querying commands for the consensus module

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

## zetacored query distribution

Querying commands for the distribution module

```
zetacored query distribution [flags]
```

### Options

```
  -h, --help   help for distribution
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
* [zetacored query distribution commission](#zetacored-query-distribution-commission)	 - Query distribution validator commission
* [zetacored query distribution community-pool](#zetacored-query-distribution-community-pool)	 - Query the amount of coins in the community pool
* [zetacored query distribution delegator-validators](#zetacored-query-distribution-delegator-validators)	 - Execute the DelegatorValidators RPC method
* [zetacored query distribution delegator-withdraw-address](#zetacored-query-distribution-delegator-withdraw-address)	 - Execute the DelegatorWithdrawAddress RPC method
* [zetacored query distribution params](#zetacored-query-distribution-params)	 - Query the current distribution parameters.
* [zetacored query distribution rewards](#zetacored-query-distribution-rewards)	 - Query all distribution delegator rewards
* [zetacored query distribution rewards-by-validator](#zetacored-query-distribution-rewards-by-validator)	 - Query all distribution delegator from a particular validator
* [zetacored query distribution slashes](#zetacored-query-distribution-slashes)	 - Query distribution validator slashes
* [zetacored query distribution validator-distribution-info](#zetacored-query-distribution-validator-distribution-info)	 - Query validator distribution info
* [zetacored query distribution validator-outstanding-rewards](#zetacored-query-distribution-validator-outstanding-rewards)	 - Query distribution outstanding (un-withdrawn) rewards for a validator and all their delegations

## zetacored query distribution commission

Query distribution validator commission

```
zetacored query distribution commission [validator] [flags]
```

### Examples

```
$ zetacored query distribution commission [validator-address]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for commission
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query distribution](#zetacored-query-distribution)	 - Querying commands for the distribution module

## zetacored query distribution community-pool

Query the amount of coins in the community pool

```
zetacored query distribution community-pool [flags]
```

### Examples

```
$ zetacored query distribution community-pool
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for community-pool
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query distribution](#zetacored-query-distribution)	 - Querying commands for the distribution module

## zetacored query distribution delegator-validators

Execute the DelegatorValidators RPC method

```
zetacored query distribution delegator-validators [flags]
```

### Options

```
      --delegator-address account address or key name   
      --grpc-addr string                                the gRPC endpoint to use for this chain
      --grpc-insecure                                   allow gRPC over insecure channels, if not the server must use TLS
      --height int                                      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                                            help for delegator-validators
      --keyring-backend string                          Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string                              The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                                       Do not indent JSON output
      --node string                                     [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string                                   Output format (text|json) 
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

* [zetacored query distribution](#zetacored-query-distribution)	 - Querying commands for the distribution module

## zetacored query distribution delegator-withdraw-address

Execute the DelegatorWithdrawAddress RPC method

```
zetacored query distribution delegator-withdraw-address [flags]
```

### Options

```
      --delegator-address account address or key name   
      --grpc-addr string                                the gRPC endpoint to use for this chain
      --grpc-insecure                                   allow gRPC over insecure channels, if not the server must use TLS
      --height int                                      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                                            help for delegator-withdraw-address
      --keyring-backend string                          Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string                              The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                                       Do not indent JSON output
      --node string                                     [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string                                   Output format (text|json) 
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

* [zetacored query distribution](#zetacored-query-distribution)	 - Querying commands for the distribution module

## zetacored query distribution params

Query the current distribution parameters.

```
zetacored query distribution params [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for params
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query distribution](#zetacored-query-distribution)	 - Querying commands for the distribution module

## zetacored query distribution rewards

Query all distribution delegator rewards

### Synopsis

Query all rewards earned by a delegator

```
zetacored query distribution rewards [delegator-addr] [flags]
```

### Examples

```
$ zetacored query distribution rewards [delegator-address]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for rewards
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query distribution](#zetacored-query-distribution)	 - Querying commands for the distribution module

## zetacored query distribution rewards-by-validator

Query all distribution delegator from a particular validator

```
zetacored query distribution rewards-by-validator [delegator-addr] [validator-addr] [flags]
```

### Examples

```
$ zetacored query distribution rewards [delegator-address] [validator-address]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for rewards-by-validator
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query distribution](#zetacored-query-distribution)	 - Querying commands for the distribution module

## zetacored query distribution slashes

Query distribution validator slashes

```
zetacored query distribution slashes [validator] [start-height] [end-height] [flags]
```

### Examples

```
$ zetacored query distribution slashes [validator-address] 0 100
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for slashes
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query distribution](#zetacored-query-distribution)	 - Querying commands for the distribution module

## zetacored query distribution validator-distribution-info

Query validator distribution info

```
zetacored query distribution validator-distribution-info [validator] [flags]
```

### Examples

```
Example: $ zetacored query distribution validator-distribution-info [validator-address]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for validator-distribution-info
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query distribution](#zetacored-query-distribution)	 - Querying commands for the distribution module

## zetacored query distribution validator-outstanding-rewards

Query distribution outstanding (un-withdrawn) rewards for a validator and all their delegations

```
zetacored query distribution validator-outstanding-rewards [validator] [flags]
```

### Examples

```
$ zetacored query distribution validator-outstanding-rewards [validator-address]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for validator-outstanding-rewards
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query distribution](#zetacored-query-distribution)	 - Querying commands for the distribution module

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

## zetacored query evidence

Querying commands for the evidence module

```
zetacored query evidence [flags]
```

### Options

```
  -h, --help   help for evidence
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
* [zetacored query evidence evidence](#zetacored-query-evidence-evidence)	 - Query for evidence by hash
* [zetacored query evidence list](#zetacored-query-evidence-list)	 - Query all (paginated) submitted evidence

## zetacored query evidence evidence

Query for evidence by hash

```
zetacored query evidence evidence [hash] [flags]
```

### Examples

```
zetacored query evidence evidence DF0C23E8634E480F84B9D5674A7CDC9816466DEC28A3358F73260F68D28D7660
```

### Options

```
      --evidence-hash binary     
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for evidence
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query evidence](#zetacored-query-evidence)	 - Querying commands for the evidence module

## zetacored query evidence list

Query all (paginated) submitted evidence

```
zetacored query evidence list [flags]
```

### Examples

```
zetacored query evidence list --page=2 --page-limit=50
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for list
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query evidence](#zetacored-query-evidence)	 - Querying commands for the evidence module

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
* [zetacored query evm 0x-to-bech32](#zetacored-query-evm-0x-to-bech32)	 - Get the bech32 address for a given 0x address
* [zetacored query evm account](#zetacored-query-evm-account)	 - Gets account info from an address
* [zetacored query evm balance-bank](#zetacored-query-evm-balance-bank)	 - Get the bank balance for a given 0x address and bank denom
* [zetacored query evm balance-erc20](#zetacored-query-evm-balance-erc20)	 - Get the ERC20 balance for a given 0x address and erc20 address
* [zetacored query evm bech32-to-0x](#zetacored-query-evm-bech32-to-0x)	 - Get the 0x address for a given bech32 address
* [zetacored query evm code](#zetacored-query-evm-code)	 - Gets code from an account
* [zetacored query evm config](#zetacored-query-evm-config)	 - Get the evm config
* [zetacored query evm params](#zetacored-query-evm-params)	 - Get the evm params
* [zetacored query evm storage](#zetacored-query-evm-storage)	 - Gets storage for an account with a given key and height

## zetacored query evm 0x-to-bech32

Get the bech32 address for a given 0x address

### Synopsis

Get the bech32 address for a given 0x address.

```
zetacored query evm 0x-to-bech32 [flags]
```

### Examples

```
evmd query evm 0x-to-bech32 0x7cB61D4117AE31a12E393a1Cfa3BaC666481D02E
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for 0x-to-bech32
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

## zetacored query evm account

Gets account info from an address

### Synopsis

Gets account info from an address. If the height is not provided, it will use the latest height from context.

```
zetacored query evm account ADDRESS [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for account
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

## zetacored query evm balance-bank

Get the bank balance for a given 0x address and bank denom

### Synopsis

Get the bank balance for a given 0x address and bank denom.

```
zetacored query evm balance-bank [address] [denom] [flags]
```

### Examples

```
evmd query evm balance-bank 0xA2A8B87390F8F2D188242656BFb6852914073D06 atoken
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for balance-bank
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

## zetacored query evm balance-erc20

Get the ERC20 balance for a given 0x address and erc20 address

### Synopsis

Get the ERC20 balance for a given 0x address and erc20 address.

```
zetacored query evm balance-erc20 [address] [erc20-address] [flags]
```

### Examples

```
evmd query evm balance-erc20 0xA2A8B87390F8F2D188242656BFb6852914073D06 0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for balance-erc20
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

## zetacored query evm bech32-to-0x

Get the 0x address for a given bech32 address

### Synopsis

Get the 0x address for a given bech32 address.

```
zetacored query evm bech32-to-0x [flags]
```

### Examples

```
evmd query evm bech32-to-0x cosmos10jmp6sgh4cc6zt3e8gw05wavvejgr5pwsjskvv
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for bech32-to-0x
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

## zetacored query evm config

Get the evm config

### Synopsis

Get the evm configuration values.

```
zetacored query evm config [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for config
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

## zetacored query gov

Querying commands for the gov module

```
zetacored query gov [flags]
```

### Options

```
  -h, --help   help for gov
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
* [zetacored query gov constitution](#zetacored-query-gov-constitution)	 - Query the current chain constitution
* [zetacored query gov deposit](#zetacored-query-gov-deposit)	 - Query details of a deposit
* [zetacored query gov deposits](#zetacored-query-gov-deposits)	 - Query deposits on a proposal
* [zetacored query gov params](#zetacored-query-gov-params)	 - Query the parameters of the governance process
* [zetacored query gov proposal](#zetacored-query-gov-proposal)	 - Query details of a single proposal
* [zetacored query gov proposals](#zetacored-query-gov-proposals)	 - Query proposals with optional filters
* [zetacored query gov tally](#zetacored-query-gov-tally)	 - Query the tally of a proposal vote
* [zetacored query gov vote](#zetacored-query-gov-vote)	 - Query details of a single vote
* [zetacored query gov votes](#zetacored-query-gov-votes)	 - Query votes of a single proposal

## zetacored query gov constitution

Query the current chain constitution

```
zetacored query gov constitution [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for constitution
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query gov](#zetacored-query-gov)	 - Querying commands for the gov module

## zetacored query gov deposit

Query details of a deposit

```
zetacored query gov deposit [proposal-id] [depositer-addr] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for deposit
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query gov](#zetacored-query-gov)	 - Querying commands for the gov module

## zetacored query gov deposits

Query deposits on a proposal

```
zetacored query gov deposits [proposal-id] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for deposits
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query gov](#zetacored-query-gov)	 - Querying commands for the gov module

## zetacored query gov params

Query the parameters of the governance process

### Synopsis

Query the parameters of the governance process. Specify specific param types (voting|tallying|deposit) to filter results.

```
zetacored query gov params [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for params
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query gov](#zetacored-query-gov)	 - Querying commands for the gov module

## zetacored query gov proposal

Query details of a single proposal

```
zetacored query gov proposal [proposal-id] [flags]
```

### Examples

```
zetacored query gov proposal 1
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for proposal
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query gov](#zetacored-query-gov)	 - Querying commands for the gov module

## zetacored query gov proposals

Query proposals with optional filters

```
zetacored query gov proposals [flags]
```

### Examples

```
zetacored query gov proposals --depositor cosmos1...
zetacored query gov proposals --voter cosmos1...
zetacored query gov proposals --proposal-status (unspecified | deposit-period | voting-period | passed | rejected | failed)
```

### Options

```
      --depositor account address or key name                                                                        
      --grpc-addr string                                                                                             the gRPC endpoint to use for this chain
      --grpc-insecure                                                                                                allow gRPC over insecure channels, if not the server must use TLS
      --height int                                                                                                   Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                                                                                                         help for proposals
      --keyring-backend string                                                                                       Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string                                                                                           The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                                                                                                    Do not indent JSON output
      --node string                                                                                                  [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string                                                                                                Output format (text|json) 
      --page-count-total                                                                                             
      --page-key binary                                                                                              
      --page-limit uint                                                                                              
      --page-offset uint                                                                                             
      --page-reverse                                                                                                 
      --proposal-status ProposalStatus (unspecified | deposit-period | voting-period | passed | rejected | failed)    (default unspecified)
      --voter account address or key name                                                                            
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

* [zetacored query gov](#zetacored-query-gov)	 - Querying commands for the gov module

## zetacored query gov tally

Query the tally of a proposal vote

```
zetacored query gov tally [proposal-id] [flags]
```

### Examples

```
zetacored query gov tally 1
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for tally
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query gov](#zetacored-query-gov)	 - Querying commands for the gov module

## zetacored query gov vote

Query details of a single vote

```
zetacored query gov vote [proposal-id] [voter-addr] [flags]
```

### Examples

```
zetacored query gov vote 1 cosmos1...
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for vote
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query gov](#zetacored-query-gov)	 - Querying commands for the gov module

## zetacored query gov votes

Query votes of a single proposal

```
zetacored query gov votes [proposal-id] [flags]
```

### Examples

```
zetacored query gov votes 1
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for votes
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query gov](#zetacored-query-gov)	 - Querying commands for the gov module

## zetacored query group

Querying commands for the group module

```
zetacored query group [flags]
```

### Options

```
  -h, --help   help for group
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
* [zetacored query group group-info](#zetacored-query-group-group-info)	 - Query for group info by group id
* [zetacored query group group-members](#zetacored-query-group-group-members)	 - Query for group members by group id
* [zetacored query group group-policies-by-admin](#zetacored-query-group-group-policies-by-admin)	 - Query for group policies by admin account address
* [zetacored query group group-policies-by-group](#zetacored-query-group-group-policies-by-group)	 - Query for group policies by group id
* [zetacored query group group-policy-info](#zetacored-query-group-group-policy-info)	 - Query for group policy info by account address of group policy
* [zetacored query group groups](#zetacored-query-group-groups)	 - Query for all groups on chain
* [zetacored query group groups-by-admin](#zetacored-query-group-groups-by-admin)	 - Query for groups by admin account address
* [zetacored query group groups-by-member](#zetacored-query-group-groups-by-member)	 - Query for groups by member address
* [zetacored query group proposal](#zetacored-query-group-proposal)	 - Query for proposal by id
* [zetacored query group proposals-by-group-policy](#zetacored-query-group-proposals-by-group-policy)	 - Query for proposals by account address of group policy
* [zetacored query group tally-result](#zetacored-query-group-tally-result)	 - Query tally result of proposal
* [zetacored query group vote](#zetacored-query-group-vote)	 - Query for vote by proposal id and voter account address
* [zetacored query group votes-by-proposal](#zetacored-query-group-votes-by-proposal)	 - Query for votes by proposal id
* [zetacored query group votes-by-voter](#zetacored-query-group-votes-by-voter)	 - Query for votes by voter account address

## zetacored query group group-info

Query for group info by group id

```
zetacored query group group-info [group-id] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for group-info
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query group](#zetacored-query-group)	 - Querying commands for the group module

## zetacored query group group-members

Query for group members by group id

```
zetacored query group group-members [group-id] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for group-members
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query group](#zetacored-query-group)	 - Querying commands for the group module

## zetacored query group group-policies-by-admin

Query for group policies by admin account address

```
zetacored query group group-policies-by-admin [admin] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for group-policies-by-admin
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query group](#zetacored-query-group)	 - Querying commands for the group module

## zetacored query group group-policies-by-group

Query for group policies by group id

```
zetacored query group group-policies-by-group [group-id] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for group-policies-by-group
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query group](#zetacored-query-group)	 - Querying commands for the group module

## zetacored query group group-policy-info

Query for group policy info by account address of group policy

```
zetacored query group group-policy-info [group-policy-account] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for group-policy-info
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query group](#zetacored-query-group)	 - Querying commands for the group module

## zetacored query group groups

Query for all groups on chain

```
zetacored query group groups [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for groups
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query group](#zetacored-query-group)	 - Querying commands for the group module

## zetacored query group groups-by-admin

Query for groups by admin account address

```
zetacored query group groups-by-admin [admin] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for groups-by-admin
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query group](#zetacored-query-group)	 - Querying commands for the group module

## zetacored query group groups-by-member

Query for groups by member address

```
zetacored query group groups-by-member [address] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for groups-by-member
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query group](#zetacored-query-group)	 - Querying commands for the group module

## zetacored query group proposal

Query for proposal by id

```
zetacored query group proposal [proposal-id] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for proposal
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query group](#zetacored-query-group)	 - Querying commands for the group module

## zetacored query group proposals-by-group-policy

Query for proposals by account address of group policy

```
zetacored query group proposals-by-group-policy [group-policy-account] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for proposals-by-group-policy
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query group](#zetacored-query-group)	 - Querying commands for the group module

## zetacored query group tally-result

Query tally result of proposal

```
zetacored query group tally-result [proposal-id] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for tally-result
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query group](#zetacored-query-group)	 - Querying commands for the group module

## zetacored query group vote

Query for vote by proposal id and voter account address

```
zetacored query group vote [proposal-id] [voter] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for vote
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query group](#zetacored-query-group)	 - Querying commands for the group module

## zetacored query group votes-by-proposal

Query for votes by proposal id

```
zetacored query group votes-by-proposal [proposal-id] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for votes-by-proposal
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query group](#zetacored-query-group)	 - Querying commands for the group module

## zetacored query group votes-by-voter

Query for votes by voter account address

```
zetacored query group votes-by-voter [voter] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for votes-by-voter
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query group](#zetacored-query-group)	 - Querying commands for the group module

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
* [zetacored query observer list-ballots](#zetacored-query-observer-list-ballots)	 - Query all ballots
* [zetacored query observer list-ballots-for-height](#zetacored-query-observer-list-ballots-for-height)	 - Query BallotListForHeight
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

## zetacored query observer list-ballots

Query all ballots

```
zetacored query observer list-ballots [flags]
```

### Options

```
      --count-total        count total number of records in list-ballots to query for
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-ballots
      --limit uint         pagination limit of list-ballots to query for (default 100)
      --node string        [host]:[port] to CometBFT RPC interface for this chain 
      --offset uint        pagination offset of list-ballots to query for
  -o, --output string      Output format (text|json) 
      --page uint          pagination page of list-ballots to query for. This sets offset to a multiple of limit (default 1)
      --page-key string    pagination page-key of list-ballots to query for
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

## zetacored query observer list-ballots-for-height

Query BallotListForHeight

```
zetacored query observer list-ballots-for-height [height] [flags]
```

### Options

```
      --grpc-addr string   the gRPC endpoint to use for this chain
      --grpc-insecure      allow gRPC over insecure channels, if not the server must use TLS
      --height int         Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help               help for list-ballots-for-height
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

## zetacored query params

Querying commands for the params module

```
zetacored query params [flags]
```

### Options

```
  -h, --help   help for params
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
* [zetacored query params subspace](#zetacored-query-params-subspace)	 - Query for raw parameters by subspace and key
* [zetacored query params subspaces](#zetacored-query-params-subspaces)	 - Query for all registered subspaces and all keys for a subspace

## zetacored query params subspace

Query for raw parameters by subspace and key

```
zetacored query params subspace [subspace] [key] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for subspace
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query params](#zetacored-query-params)	 - Querying commands for the params module

## zetacored query params subspaces

Query for all registered subspaces and all keys for a subspace

```
zetacored query params subspaces [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for subspaces
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query params](#zetacored-query-params)	 - Querying commands for the params module

## zetacored query slashing

Querying commands for the slashing module

```
zetacored query slashing [flags]
```

### Options

```
  -h, --help   help for slashing
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
* [zetacored query slashing params](#zetacored-query-slashing-params)	 - Query the current slashing parameters
* [zetacored query slashing signing-info](#zetacored-query-slashing-signing-info)	 - Query a validator's signing information
* [zetacored query slashing signing-infos](#zetacored-query-slashing-signing-infos)	 - Query signing information of all validators

## zetacored query slashing params

Query the current slashing parameters

```
zetacored query slashing params [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for params
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query slashing](#zetacored-query-slashing)	 - Querying commands for the slashing module

## zetacored query slashing signing-info

Query a validator's signing information

### Synopsis

Query a validator's signing information, with a pubkey ('zetacored comet show-validator') or a validator consensus address

```
zetacored query slashing signing-info [validator-conspub/address] [flags]
```

### Examples

```
zetacored query slashing signing-info '{"@type":"/cosmos.crypto.ed25519.PubKey","key":"OauFcTKbN5Lx3fJL689cikXBqe+hcp6Y+x0rYUdR9Jk="}'
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for signing-info
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query slashing](#zetacored-query-slashing)	 - Querying commands for the slashing module

## zetacored query slashing signing-infos

Query signing information of all validators

```
zetacored query slashing signing-infos [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for signing-infos
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query slashing](#zetacored-query-slashing)	 - Querying commands for the slashing module

## zetacored query staking

Querying commands for the staking module

```
zetacored query staking [flags]
```

### Options

```
  -h, --help   help for staking
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
* [zetacored query staking delegation](#zetacored-query-staking-delegation)	 - Query a delegation based on address and validator address
* [zetacored query staking delegations](#zetacored-query-staking-delegations)	 - Query all delegations made by one delegator
* [zetacored query staking delegations-to](#zetacored-query-staking-delegations-to)	 - Query all delegations made to one validator
* [zetacored query staking delegator-validator](#zetacored-query-staking-delegator-validator)	 - Query validator info for given delegator validator pair
* [zetacored query staking delegator-validators](#zetacored-query-staking-delegator-validators)	 - Query all validators info for given delegator address
* [zetacored query staking historical-info](#zetacored-query-staking-historical-info)	 - Query historical info at given height
* [zetacored query staking params](#zetacored-query-staking-params)	 - Query the current staking parameters information
* [zetacored query staking pool](#zetacored-query-staking-pool)	 - Query the current staking pool values
* [zetacored query staking redelegation](#zetacored-query-staking-redelegation)	 - Query a redelegation record based on delegator and a source and destination validator address
* [zetacored query staking unbonding-delegation](#zetacored-query-staking-unbonding-delegation)	 - Query an unbonding-delegation record based on delegator and validator address
* [zetacored query staking unbonding-delegations](#zetacored-query-staking-unbonding-delegations)	 - Query all unbonding-delegations records for one delegator
* [zetacored query staking unbonding-delegations-from](#zetacored-query-staking-unbonding-delegations-from)	 - Query all unbonding delegatations from a validator
* [zetacored query staking validator](#zetacored-query-staking-validator)	 - Query a validator
* [zetacored query staking validators](#zetacored-query-staking-validators)	 - Query for all validators

## zetacored query staking delegation

Query a delegation based on address and validator address

### Synopsis

Query delegations for an individual delegator on an individual validator

```
zetacored query staking delegation [delegator-addr] [validator-addr] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for delegation
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query staking](#zetacored-query-staking)	 - Querying commands for the staking module

## zetacored query staking delegations

Query all delegations made by one delegator

### Synopsis

Query delegations for an individual delegator on all validators.

```
zetacored query staking delegations [delegator-addr] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for delegations
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query staking](#zetacored-query-staking)	 - Querying commands for the staking module

## zetacored query staking delegations-to

Query all delegations made to one validator

### Synopsis

Query delegations on an individual validator.

```
zetacored query staking delegations-to [validator-addr] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for delegations-to
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query staking](#zetacored-query-staking)	 - Querying commands for the staking module

## zetacored query staking delegator-validator

Query validator info for given delegator validator pair

```
zetacored query staking delegator-validator [delegator-addr] [validator-addr] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for delegator-validator
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query staking](#zetacored-query-staking)	 - Querying commands for the staking module

## zetacored query staking delegator-validators

Query all validators info for given delegator address

```
zetacored query staking delegator-validators [delegator-addr] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for delegator-validators
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query staking](#zetacored-query-staking)	 - Querying commands for the staking module

## zetacored query staking historical-info

Query historical info at given height

```
zetacored query staking historical-info [height] [flags]
```

### Examples

```
$ zetacored query staking historical-info 5
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for historical-info
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query staking](#zetacored-query-staking)	 - Querying commands for the staking module

## zetacored query staking params

Query the current staking parameters information

### Synopsis

Query values set as staking parameters.

```
zetacored query staking params [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for params
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query staking](#zetacored-query-staking)	 - Querying commands for the staking module

## zetacored query staking pool

Query the current staking pool values

### Synopsis

Query values for amounts stored in the staking pool.

```
zetacored query staking pool [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for pool
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query staking](#zetacored-query-staking)	 - Querying commands for the staking module

## zetacored query staking redelegation

Query a redelegation record based on delegator and a source and destination validator address

### Synopsis

Query a redelegation record for an individual delegator between a source and destination validator.

```
zetacored query staking redelegation [delegator-addr] [src-validator-addr] [dst-validator-addr] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for redelegation
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query staking](#zetacored-query-staking)	 - Querying commands for the staking module

## zetacored query staking unbonding-delegation

Query an unbonding-delegation record based on delegator and validator address

### Synopsis

Query unbonding delegations for an individual delegator on an individual validator.

```
zetacored query staking unbonding-delegation [delegator-addr] [validator-addr] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for unbonding-delegation
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query staking](#zetacored-query-staking)	 - Querying commands for the staking module

## zetacored query staking unbonding-delegations

Query all unbonding-delegations records for one delegator

### Synopsis

Query unbonding delegations for an individual delegator.

```
zetacored query staking unbonding-delegations [delegator-addr] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for unbonding-delegations
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query staking](#zetacored-query-staking)	 - Querying commands for the staking module

## zetacored query staking unbonding-delegations-from

Query all unbonding delegatations from a validator

### Synopsis

Query delegations that are unbonding _from_ a validator.

```
zetacored query staking unbonding-delegations-from [validator-addr] [flags]
```

### Examples

```
$ zetacored query staking unbonding-delegations-from [val-addr]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for unbonding-delegations-from
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
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

* [zetacored query staking](#zetacored-query-staking)	 - Querying commands for the staking module

## zetacored query staking validator

Query a validator

### Synopsis

Query details about an individual validator.

```
zetacored query staking validator [validator-addr] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for validator
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query staking](#zetacored-query-staking)	 - Querying commands for the staking module

## zetacored query staking validators

Query for all validators

### Synopsis

Query details about all validators on a network.

```
zetacored query staking validators [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for validators
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
      --page-count-total         
      --page-key binary          
      --page-limit uint          
      --page-offset uint         
      --page-reverse             
      --status string            
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

* [zetacored query staking](#zetacored-query-staking)	 - Querying commands for the staking module

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

## zetacored query upgrade

Querying commands for the upgrade module

```
zetacored query upgrade [flags]
```

### Options

```
  -h, --help   help for upgrade
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
* [zetacored query upgrade applied](#zetacored-query-upgrade-applied)	 - Query the block header for height at which a completed upgrade was applied
* [zetacored query upgrade authority](#zetacored-query-upgrade-authority)	 - Get the upgrade authority address
* [zetacored query upgrade module-versions](#zetacored-query-upgrade-module-versions)	 - Query the list of module versions
* [zetacored query upgrade plan](#zetacored-query-upgrade-plan)	 - Query the upgrade plan (if one exists)

## zetacored query upgrade applied

Query the block header for height at which a completed upgrade was applied

### Synopsis

If upgrade-name was previously executed on the chain, this returns the header for the block at which it was applied. This helps a client determine which binary was valid over a given range of blocks, as well as more context to understand past migrations.

```
zetacored query upgrade applied [upgrade-name] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for applied
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query upgrade](#zetacored-query-upgrade)	 - Querying commands for the upgrade module

## zetacored query upgrade authority

Get the upgrade authority address

```
zetacored query upgrade authority [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for authority
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query upgrade](#zetacored-query-upgrade)	 - Querying commands for the upgrade module

## zetacored query upgrade module-versions

Query the list of module versions

### Synopsis

Gets a list of module names and their respective consensus versions. Following the command with a specific module name will return only that module's information.

```
zetacored query upgrade module-versions [optional module_name] [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for module-versions
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query upgrade](#zetacored-query-upgrade)	 - Querying commands for the upgrade module

## zetacored query upgrade plan

Query the upgrade plan (if one exists)

### Synopsis

Gets the currently scheduled upgrade plan, if one exists

```
zetacored query upgrade plan [flags]
```

### Options

```
      --grpc-addr string         the gRPC endpoint to use for this chain
      --grpc-insecure            allow gRPC over insecure channels, if not the server must use TLS
      --height int               Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help                     help for plan
      --keyring-backend string   Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --no-indent                Do not indent JSON output
      --node string              [host]:[port] to CometBFT RPC interface for this chain 
  -o, --output string            Output format (text|json) 
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

* [zetacored query upgrade](#zetacored-query-upgrade)	 - Querying commands for the upgrade module

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

Run the full node application with CometBFT in or out of process. By
default, the application will run with CometBFT in process.

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
      --evm.cache-preimage                              Enables tracking of SHA3 preimages in the EVM (not implemented yet)
      --evm.evm-chain-id uint                           the EIP-155 compatible replay protection chain ID (default 262144)
      --evm.max-tx-gas-wanted uint                      the gas wanted for each eth tx returned in ante handler in check tx mode
      --evm.tracer string                               the EVM tracer type to collect execution traces from the EVM transaction execution (json|struct|access_list|markdown)
      --genesis_hash bytesHex                           optional SHA-256 hash of the genesis file
      --grpc-only                                       Start the node in gRPC query only mode without CometBFT process
      --grpc-web.address string                         The gRPC-Web server address to listen on 
      --grpc-web.enable                                 Define if the gRPC-Web server should be enabled. (Note: gRPC must also be enabled.)
      --grpc.address string                             the gRPC server address to listen on 
      --grpc.enable                                     Define if the gRPC server should be enabled
      --halt-height uint                                Block height at which to gracefully halt the chain and shutdown the node
      --halt-time uint                                  Minimum block time (in Unix seconds) at which to gracefully halt the chain and shutdown the node
  -h, --help                                            help for start
      --home string                                     The application home directory 
      --inter-block-cache                               Enable inter-block caching (default true)
      --inv-check-period uint                           Assert registered invariants every N blocks
      --json-rpc.address string                         the JSON-RPC server address to listen on 
      --json-rpc.allow-insecure-unlock                  Allow insecure account unlocking when account-related RPCs are exposed by http (default true)
      --json-rpc.allow-unprotected-txs                  Allow for unprotected (non EIP155 signed) transactions to be submitted via the node's RPC when the global parameter is disabled
      --json-rpc.api strings                            Defines a list of JSON-RPC namespaces that should be enabled (default [eth,net,web3])
      --json-rpc.block-range-cap eth_getLogs            Sets the max block range allowed for eth_getLogs query (default 10000)
      --json-rpc.enable                                 Define if the JSON-RPC server should be enabled
      --json-rpc.enable-indexer                         Enable the custom tx indexer for json-rpc
      --json-rpc.evm-timeout duration                   Sets a timeout used for eth_call (0=infinite) (default 5s)
      --json-rpc.filter-cap int32                       Sets the global cap for total number of filters that can be created (default 200)
      --json-rpc.gas-cap uint                           Sets a cap on gas that can be used in eth_call/estimateGas unit is aatom (0=infinite) (default 25000000)
      --json-rpc.http-idle-timeout duration             Sets a idle timeout for json-rpc http server (0=infinite) (default 2m0s)
      --json-rpc.http-timeout duration                  Sets a read/write timeout for json-rpc http server (0=infinite) (default 30s)
      --json-rpc.logs-cap eth_getLogs                   Sets the max number of results can be returned from single eth_getLogs query (default 10000)
      --json-rpc.max-open-connections int               Sets the maximum number of simultaneous connections for the server listener
      --json-rpc.txfee-cap float                        Sets a cap on transaction fee that can be sent via the RPC APIs (1 = default 1 evmos) (default 1)
      --json-rpc.ws-address string                      the JSON-RPC WS server address to listen on 
      --metrics                                         Define if EVM rpc metrics server should be enabled
      --min-retain-blocks uint                          Minimum block height offset during ABCI commit to prune CometBFT blocks
      --minimum-gas-prices string                       Minimum gas prices to accept for transactions; Any fee in a tx must meet this minimum (e.g. 20000000000azeta)
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
      --skip-config-overwrite                           Skip running the config configuration overwrite handler.This is used for testing purposes only and skips using the default timeouts hardcoded and uses the config file instead
      --state-sync.snapshot-interval uint               State sync snapshot interval
      --state-sync.snapshot-keep-recent uint32          State sync snapshot to keep (default 2)
      --tls.certificate-path string                     the cert.pem file path for the server TLS configuration
      --tls.key-path string                             the key.pem file path for the server TLS configuration
      --trace                                           Provide full stack traces for errors in ABCI Log
      --trace-store string                              Enable KVStore tracing to an output file
      --transport string                                Transport protocol: socket, grpc 
      --unsafe-skip-upgrades ints                       Skip a set of upgrade heights to continue the old binary
      --with-cometbft                                   Run abci app embedded in-process with CometBFT (default true)
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

## zetacored status

Query remote node for status

```
zetacored status [flags]
```

### Options

```
  -h, --help            help for status
  -n, --node string     Node to connect to 
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
* [zetacored tx auth](#zetacored-tx-auth)	 - Transactions commands for the auth module
* [zetacored tx authority](#zetacored-tx-authority)	 - authority transactions subcommands
* [zetacored tx authz](#zetacored-tx-authz)	 - Authorization transactions subcommands
* [zetacored tx bank](#zetacored-tx-bank)	 - Bank transaction subcommands
* [zetacored tx broadcast](#zetacored-tx-broadcast)	 - Broadcast transactions generated offline
* [zetacored tx consensus](#zetacored-tx-consensus)	 - Transactions commands for the consensus module
* [zetacored tx crisis](#zetacored-tx-crisis)	 - Transactions commands for the crisis module
* [zetacored tx crosschain](#zetacored-tx-crosschain)	 - crosschain transactions subcommands
* [zetacored tx decode](#zetacored-tx-decode)	 - Decode a binary encoded transaction string
* [zetacored tx distribution](#zetacored-tx-distribution)	 - Distribution transactions subcommands
* [zetacored tx emissions](#zetacored-tx-emissions)	 - emissions transactions subcommands
* [zetacored tx encode](#zetacored-tx-encode)	 - Encode transactions generated offline
* [zetacored tx evidence](#zetacored-tx-evidence)	 - Evidence transaction subcommands
* [zetacored tx evm](#zetacored-tx-evm)	 - evm subcommands
* [zetacored tx feemarket](#zetacored-tx-feemarket)	 - Transactions commands for the feemarket module
* [zetacored tx fungible](#zetacored-tx-fungible)	 - fungible transactions subcommands
* [zetacored tx gov](#zetacored-tx-gov)	 - Governance transactions subcommands
* [zetacored tx group](#zetacored-tx-group)	 - Group transaction subcommands
* [zetacored tx lightclient](#zetacored-tx-lightclient)	 - lightclient transactions subcommands
* [zetacored tx multi-sign](#zetacored-tx-multi-sign)	 - Generate multisig signatures for transactions generated offline
* [zetacored tx multisign-batch](#zetacored-tx-multisign-batch)	 - Assemble multisig transactions in batch from batch signatures
* [zetacored tx observer](#zetacored-tx-observer)	 - observer transactions subcommands
* [zetacored tx sign](#zetacored-tx-sign)	 - Sign a transaction generated offline
* [zetacored tx sign-batch](#zetacored-tx-sign-batch)	 - Sign transaction batch files
* [zetacored tx slashing](#zetacored-tx-slashing)	 - Transactions commands for the slashing module
* [zetacored tx staking](#zetacored-tx-staking)	 - Staking transaction subcommands
* [zetacored tx upgrade](#zetacored-tx-upgrade)	 - Upgrade transaction subcommands
* [zetacored tx validate-signatures](#zetacored-tx-validate-signatures)	 - validate transactions signatures
* [zetacored tx vesting](#zetacored-tx-vesting)	 - Vesting transaction subcommands

## zetacored tx auth

Transactions commands for the auth module

```
zetacored tx auth [flags]
```

### Options

```
  -h, --help   help for auth
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

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands
* [zetacored tx auth update-params-proposal](#zetacored-tx-auth-update-params-proposal)	 - Submit a proposal to update auth module params. Note: the entire params must be provided.

## zetacored tx auth update-params-proposal

Submit a proposal to update auth module params. Note: the entire params must be provided.

```
zetacored tx auth update-params-proposal [params] [flags]
```

### Examples

```
zetacored tx auth update-params-proposal '{ "max_memo_characters": 0, "tx_sig_limit": 0, "tx_size_cost_per_byte": 0, "sig_verify_cost_ed25519": 0, "sig_verify_cost_secp256k1": 0 }'
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for update-params-proposal
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx auth](#zetacored-tx-auth)	 - Transactions commands for the auth module

## zetacored tx authority

authority transactions subcommands

```
zetacored tx authority [flags]
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

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands
* [zetacored tx authority add-authorization](#zetacored-tx-authority-add-authorization)	 - Add a new authorization or update the policy of an existing authorization. Policy type can be 0 for groupEmergency, 1 for groupOperational, 2 for groupAdmin.
* [zetacored tx authority remove-authorization](#zetacored-tx-authority-remove-authorization)	 - removes an existing authorization
* [zetacored tx authority remove-chain-info](#zetacored-tx-authority-remove-chain-info)	 - Remove the chain info for the specified chain id
* [zetacored tx authority update-chain-info](#zetacored-tx-authority-update-chain-info)	 - Update the chain info
* [zetacored tx authority update-policies](#zetacored-tx-authority-update-policies)	 - Update policies to values provided in the JSON file.

## zetacored tx authority add-authorization

Add a new authorization or update the policy of an existing authorization. Policy type can be 0 for groupEmergency, 1 for groupOperational, 2 for groupAdmin.

```
zetacored tx authority add-authorization [msg-url] [authorized-policy] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for add-authorization
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx authority](#zetacored-tx-authority)	 - authority transactions subcommands

## zetacored tx authority remove-authorization

removes an existing authorization

```
zetacored tx authority remove-authorization [msg-url] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for remove-authorization
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx authority](#zetacored-tx-authority)	 - authority transactions subcommands

## zetacored tx authority remove-chain-info

Remove the chain info for the specified chain id

```
zetacored tx authority remove-chain-info [chain-id] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for remove-chain-info
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx authority](#zetacored-tx-authority)	 - authority transactions subcommands

## zetacored tx authority update-chain-info

Update the chain info

```
zetacored tx authority update-chain-info [chain-info-json-file] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for update-chain-info
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx authority](#zetacored-tx-authority)	 - authority transactions subcommands

## zetacored tx authority update-policies

Update policies to values provided in the JSON file.

```
zetacored tx authority update-policies [policies-json-file] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for update-policies
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx authority](#zetacored-tx-authority)	 - authority transactions subcommands

## zetacored tx authz

Authorization transactions subcommands

### Synopsis

Authorize and revoke access to execute transactions on behalf of your address

```
zetacored tx authz [flags]
```

### Options

```
  -h, --help   help for authz
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

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands
* [zetacored tx authz exec](#zetacored-tx-authz-exec)	 - execute tx on behalf of granter account
* [zetacored tx authz grant](#zetacored-tx-authz-grant)	 - Grant authorization to an address
* [zetacored tx authz revoke](#zetacored-tx-authz-revoke)	 - revoke authorization

## zetacored tx authz exec

execute tx on behalf of granter account

### Synopsis

execute tx on behalf of granter account:
Example:
 $ zetacored tx authz exec tx.json --from grantee
 $ zetacored tx bank send [granter] [recipient] --from [granter] --chain-id [chain-id] --generate-only > tx.json && zetacored tx authz exec tx.json --from grantee

```
zetacored tx authz exec [tx-json-file] --from [grantee] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for exec
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx authz](#zetacored-tx-authz)	 - Authorization transactions subcommands

## zetacored tx authz grant

Grant authorization to an address

### Synopsis

create a new grant authorization to an address to execute a transaction on your behalf:

Examples:
 $ zetacored tx authz grant cosmos1skjw.. send --spend-limit=1000stake --from=cosmos1skl..
 $ zetacored tx authz grant cosmos1skjw.. generic --msg-type=/cosmos.gov.v1.MsgVote --from=cosmos1sk..

```
zetacored tx authz grant [grantee] [authorization_type="send"|"generic"|"delegate"|"unbond"|"redelegate"] --from [granter] [flags]
```

### Options

```
  -a, --account-number uint          The account number of the signing account (offline mode only)
      --allow-list strings           Allowed addresses grantee is allowed to send funds separated by ,
      --allowed-validators strings   Allowed validators addresses separated by ,
      --aux                          Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string        Transaction broadcasting mode (sync|async) 
      --chain-id string              The network chain ID
      --deny-validators strings      Deny validators addresses separated by ,
      --dry-run                      ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --expiration int               Expire time as Unix timestamp. Set zero (0) for no expiry. Default is 0.
      --fee-granter string           Fee granter grants fees for the transaction
      --fee-payer string             Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                  Fees to pay along with transaction; eg: 10uatom
      --from string                  Name or address of private key with which to sign
      --gas string                   gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float         adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string            Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only                Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                         help for grant
      --keyring-backend string       Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string           The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                       Use a connected Ledger device
      --msg-type string              The Msg method name for which we are creating a GenericAuthorization
      --node string                  [host]:[port] to CometBFT rpc interface for this chain 
      --note string                  Note to add a description to the transaction (previously --memo)
      --offline                      Offline mode (does not allow any online functionality)
  -o, --output string                Output format (text|json) 
  -s, --sequence uint                The sequence number of the signing account (offline mode only)
      --sign-mode string             Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --spend-limit string           SpendLimit for Send Authorization, an array of Coins allowed spend
      --timeout-duration duration    TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint          DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                   Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                    Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                          Skip tx broadcasting prompt confirmation
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

* [zetacored tx authz](#zetacored-tx-authz)	 - Authorization transactions subcommands

## zetacored tx authz revoke

revoke authorization

### Synopsis

revoke authorization from a granter to a grantee:
Example:
 $ zetacored tx authz revoke cosmos1skj.. /cosmos.bank.v1beta1.MsgSend --from=cosmos1skj..

```
zetacored tx authz revoke [grantee] [msg-type-url] --from=[granter] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for revoke
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx authz](#zetacored-tx-authz)	 - Authorization transactions subcommands

## zetacored tx bank

Bank transaction subcommands

```
zetacored tx bank [flags]
```

### Options

```
  -h, --help   help for bank
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

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands
* [zetacored tx bank multi-send](#zetacored-tx-bank-multi-send)	 - Send funds from one account to two or more accounts.
* [zetacored tx bank send](#zetacored-tx-bank-send)	 - Send funds from one account to another.
* [zetacored tx bank set-send-enabled-proposal](#zetacored-tx-bank-set-send-enabled-proposal)	 - Submit a proposal to set/update/delete send enabled entries
* [zetacored tx bank update-params-proposal](#zetacored-tx-bank-update-params-proposal)	 - Submit a proposal to update bank module params. Note: the entire params must be provided.

## zetacored tx bank multi-send

Send funds from one account to two or more accounts.

### Synopsis

Send funds from one account to two or more accounts.
By default, sends the [amount] to each address of the list.
Using the '--split' flag, the [amount] is split equally between the addresses.
Note, the '--from' flag is ignored as it is implied from [from_key_or_address] and 
separate addresses with space.
When using '--dry-run' a key name cannot be used, only a bech32 address.

```
zetacored tx bank multi-send [from_key_or_address] [to_address_1 to_address_2 ...] [amount] [flags]
```

### Examples

```
zetacored tx bank multi-send cosmos1... cosmos1... cosmos1... cosmos1... 10stake
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for multi-send
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --split                       Send the equally split token amount to each address
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx bank](#zetacored-tx-bank)	 - Bank transaction subcommands

## zetacored tx bank send

Send funds from one account to another.

### Synopsis

Send funds from one account to another.
Note, the '--from' flag is ignored as it is implied from [from_key_or_address].
When using '--dry-run' a key name cannot be used, only a bech32 address.


```
zetacored tx bank send [from_key_or_address] [to_address] [amount] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for send
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx bank](#zetacored-tx-bank)	 - Bank transaction subcommands

## zetacored tx bank set-send-enabled-proposal

Submit a proposal to set/update/delete send enabled entries

```
zetacored tx bank set-send-enabled-proposal [send_enabled] [flags]
```

### Examples

```
zetacored tx bank set-send-enabled-proposal '{"denom":"stake","enabled":true}'
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for set-send-enabled-proposal
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
      --use-default-for strings     Use default for the given denom (delete a send enabled entry)
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx bank](#zetacored-tx-bank)	 - Bank transaction subcommands

## zetacored tx bank update-params-proposal

Submit a proposal to update bank module params. Note: the entire params must be provided.

```
zetacored tx bank update-params-proposal [params] [flags]
```

### Examples

```
zetacored tx bank update-params-proposal '{ "default_send_enabled": true }'
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for update-params-proposal
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx bank](#zetacored-tx-bank)	 - Bank transaction subcommands

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
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for broadcast
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

## zetacored tx consensus

Transactions commands for the consensus module

```
zetacored tx consensus [flags]
```

### Options

```
  -h, --help   help for consensus
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

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands
* [zetacored tx consensus update-params-proposal](#zetacored-tx-consensus-update-params-proposal)	 - Submit a proposal to update consensus module params. Note: the entire params must be provided.

## zetacored tx consensus update-params-proposal

Submit a proposal to update consensus module params. Note: the entire params must be provided.

```
zetacored tx consensus update-params-proposal [params] [flags]
```

### Examples

```
zetacored tx consensus update-params-proposal '{ params }'
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for update-params-proposal
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx consensus](#zetacored-tx-consensus)	 - Transactions commands for the consensus module

## zetacored tx crisis

Transactions commands for the crisis module

```
zetacored tx crisis [flags]
```

### Options

```
  -h, --help   help for crisis
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

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands
* [zetacored tx crisis invariant-broken](#zetacored-tx-crisis-invariant-broken)	 - Submit proof that an invariant broken

## zetacored tx crisis invariant-broken

Submit proof that an invariant broken

```
zetacored tx crisis invariant-broken [module-name] [invariant-route] --from mykey [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for invariant-broken
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx crisis](#zetacored-tx-crisis)	 - Transactions commands for the crisis module

## zetacored tx crosschain

crosschain transactions subcommands

```
zetacored tx crosschain [flags]
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

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands
* [zetacored tx crosschain abort-stuck-cctx](#zetacored-tx-crosschain-abort-stuck-cctx)	 - abort a stuck CCTX
* [zetacored tx crosschain add-inbound-tracker](#zetacored-tx-crosschain-add-inbound-tracker)	 - Add an inbound tracker 
				Use 0:Zeta,1:Gas,2:ERC20
* [zetacored tx crosschain add-outbound-tracker](#zetacored-tx-crosschain-add-outbound-tracker)	 - Add an outbound tracker
* [zetacored tx crosschain migrate-tss-funds](#zetacored-tx-crosschain-migrate-tss-funds)	 - Migrate TSS funds to the latest TSS address
* [zetacored tx crosschain refund-aborted](#zetacored-tx-crosschain-refund-aborted)	 - Refund an aborted tx , the refund address is optional, if not provided, the refund will be sent to the sender/tx origin of the cctx.
* [zetacored tx crosschain remove-inbound-tracker](#zetacored-tx-crosschain-remove-inbound-tracker)	 - Remove an inbound tracker
* [zetacored tx crosschain remove-outbound-tracker](#zetacored-tx-crosschain-remove-outbound-tracker)	 - Remove an outbound tracker
* [zetacored tx crosschain update-tss-address](#zetacored-tx-crosschain-update-tss-address)	 - Create a new TSSVoter
* [zetacored tx crosschain vote-gas-price](#zetacored-tx-crosschain-vote-gas-price)	 - Broadcast message to vote gas price
* [zetacored tx crosschain vote-inbound](#zetacored-tx-crosschain-vote-inbound)	 - Broadcast message to vote an inbound
* [zetacored tx crosschain vote-outbound](#zetacored-tx-crosschain-vote-outbound)	 - Broadcast message to vote an outbound
* [zetacored tx crosschain whitelist-erc20](#zetacored-tx-crosschain-whitelist-erc20)	 - Add a new erc20 token to whitelist

## zetacored tx crosschain abort-stuck-cctx

abort a stuck CCTX

```
zetacored tx crosschain abort-stuck-cctx [index] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for abort-stuck-cctx
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx crosschain](#zetacored-tx-crosschain)	 - crosschain transactions subcommands

## zetacored tx crosschain add-inbound-tracker

Add an inbound tracker 
				Use 0:Zeta,1:Gas,2:ERC20

```
zetacored tx crosschain add-inbound-tracker [chain-id] [tx-hash] [coin-type] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for add-inbound-tracker
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx crosschain](#zetacored-tx-crosschain)	 - crosschain transactions subcommands

## zetacored tx crosschain add-outbound-tracker

Add an outbound tracker

```
zetacored tx crosschain add-outbound-tracker [chain] [nonce] [tx-hash] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for add-outbound-tracker
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx crosschain](#zetacored-tx-crosschain)	 - crosschain transactions subcommands

## zetacored tx crosschain migrate-tss-funds

Migrate TSS funds to the latest TSS address

```
zetacored tx crosschain migrate-tss-funds [chainID] [amount] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for migrate-tss-funds
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx crosschain](#zetacored-tx-crosschain)	 - crosschain transactions subcommands

## zetacored tx crosschain refund-aborted

Refund an aborted tx , the refund address is optional, if not provided, the refund will be sent to the sender/tx origin of the cctx.

```
zetacored tx crosschain refund-aborted [cctx-index] [refund-address] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for refund-aborted
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx crosschain](#zetacored-tx-crosschain)	 - crosschain transactions subcommands

## zetacored tx crosschain remove-inbound-tracker

Remove an inbound tracker

```
zetacored tx crosschain remove-inbound-tracker [chain-id] [tx-hash] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for remove-inbound-tracker
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx crosschain](#zetacored-tx-crosschain)	 - crosschain transactions subcommands

## zetacored tx crosschain remove-outbound-tracker

Remove an outbound tracker

```
zetacored tx crosschain remove-outbound-tracker [chain] [nonce] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for remove-outbound-tracker
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx crosschain](#zetacored-tx-crosschain)	 - crosschain transactions subcommands

## zetacored tx crosschain update-tss-address

Create a new TSSVoter

```
zetacored tx crosschain update-tss-address [pubkey] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for update-tss-address
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx crosschain](#zetacored-tx-crosschain)	 - crosschain transactions subcommands

## zetacored tx crosschain vote-gas-price

Broadcast message to vote gas price

```
zetacored tx crosschain vote-gas-price [chain] [price] [priorityFee] [blockNumber] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for vote-gas-price
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx crosschain](#zetacored-tx-crosschain)	 - crosschain transactions subcommands

## zetacored tx crosschain vote-inbound

Broadcast message to vote an inbound

```
zetacored tx crosschain vote-inbound [sender] [senderChainID] [txOrigin] [receiver] [receiverChainID] [amount] [message] [inboundHash] [inBlockHeight] [coinType] [asset] [eventIndex] [protocolContractVersion] [isArbitraryCall] [confirmationMode] [inboundStatus] [flags]
```

### Examples

```
zetacored tx crosschain vote-inbound 0xfa233D806C8EB69548F3c4bC0ABb46FaD4e2EB26 8453 "" 0xfa233D806C8EB69548F3c4bC0ABb46FaD4e2EB26 7000 1000000 "" 0x66b59ad844404e91faa9587a3061e2f7af36f7a7a1a0afaca3a2efd811bc9463 26170791 Gas 0x0000000000000000000000000000000000000000 587 V2 FALSE SAFE SUCCESS
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for vote-inbound
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx crosschain](#zetacored-tx-crosschain)	 - crosschain transactions subcommands

## zetacored tx crosschain vote-outbound

Broadcast message to vote an outbound

```
zetacored tx crosschain vote-outbound [sendHash] [outboundHash] [outBlockHeight] [outGasUsed] [outEffectiveGasPrice] [outEffectiveGasLimit] [valueReceived] [Status] [chain] [outTXNonce] [coinType] [confirmationMode] [flags]
```

### Examples

```
zetacored tx crosschain vote-outbound 0x12044bec3b050fb28996630e9f2e9cc8d6cf9ef0e911e73348ade46c7ba3417a 0x4f29f9199b10189c8d02b83568aba4cb23984f11adf23e7e5d2eb037ca309497 67773716 65646 30011221226 100000 297254 0 137 13812 ERC20 SAFE
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for vote-outbound
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx crosschain](#zetacored-tx-crosschain)	 - crosschain transactions subcommands

## zetacored tx crosschain whitelist-erc20

Add a new erc20 token to whitelist

```
zetacored tx crosschain whitelist-erc20 [erc20Address] [chainID] [name] [symbol] [decimals] [gasLimit] [liquidityCap] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for whitelist-erc20
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx crosschain](#zetacored-tx-crosschain)	 - crosschain transactions subcommands

## zetacored tx decode

Decode a binary encoded transaction string

```
zetacored tx decode [protobuf-byte-string] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for decode
  -x, --hex                         Treat input as hexadecimal instead of base64
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

## zetacored tx distribution

Distribution transactions subcommands

```
zetacored tx distribution [flags]
```

### Options

```
  -h, --help   help for distribution
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

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands
* [zetacored tx distribution community-pool-spend-proposal](#zetacored-tx-distribution-community-pool-spend-proposal)	 - Submit a proposal to spend from the community pool
* [zetacored tx distribution fund-community-pool](#zetacored-tx-distribution-fund-community-pool)	 - Funds the community pool with the specified amount
* [zetacored tx distribution fund-validator-rewards-pool](#zetacored-tx-distribution-fund-validator-rewards-pool)	 - Fund the validator rewards pool with the specified amount
* [zetacored tx distribution set-withdraw-addr](#zetacored-tx-distribution-set-withdraw-addr)	 - change the default withdraw address for rewards associated with an address
* [zetacored tx distribution update-params-proposal](#zetacored-tx-distribution-update-params-proposal)	 - Submit a proposal to update distribution module params. Note: the entire params must be provided.
* [zetacored tx distribution withdraw-all-rewards](#zetacored-tx-distribution-withdraw-all-rewards)	 - withdraw all delegations rewards for a delegator
* [zetacored tx distribution withdraw-rewards](#zetacored-tx-distribution-withdraw-rewards)	 - Withdraw rewards from a given delegation address, and optionally withdraw validator commission if the delegation address given is a validator operator
* [zetacored tx distribution withdraw-validator-commission](#zetacored-tx-distribution-withdraw-validator-commission)	 - Withdraw commissions from a validator address (must be a validator operator)

## zetacored tx distribution community-pool-spend-proposal

Submit a proposal to spend from the community pool

```
zetacored tx distribution community-pool-spend-proposal [recipient] [amount] [flags]
```

### Examples

```
$ zetacored tx distribution community-pool-spend-proposal [recipient] 100uatom
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for community-pool-spend-proposal
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx distribution](#zetacored-tx-distribution)	 - Distribution transactions subcommands

## zetacored tx distribution fund-community-pool

Funds the community pool with the specified amount

### Synopsis

Funds the community pool with the specified amount

Example:
$ zetacored tx distribution fund-community-pool 100uatom --from mykey

```
zetacored tx distribution fund-community-pool [amount] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for fund-community-pool
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx distribution](#zetacored-tx-distribution)	 - Distribution transactions subcommands

## zetacored tx distribution fund-validator-rewards-pool

Fund the validator rewards pool with the specified amount

```
zetacored tx distribution fund-validator-rewards-pool [val_addr] [amount] [flags]
```

### Examples

```
zetacored tx distribution fund-validator-rewards-pool cosmosvaloper1x20lytyf6zkcrv5edpkfkn8sz578qg5sqfyqnp 100uatom --from mykey
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for fund-validator-rewards-pool
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx distribution](#zetacored-tx-distribution)	 - Distribution transactions subcommands

## zetacored tx distribution set-withdraw-addr

change the default withdraw address for rewards associated with an address

### Synopsis

Set the withdraw address for rewards associated with a delegator address.

Example:
$ zetacored tx distribution set-withdraw-addr zeta1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p --from mykey

```
zetacored tx distribution set-withdraw-addr [withdraw-addr] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for set-withdraw-addr
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx distribution](#zetacored-tx-distribution)	 - Distribution transactions subcommands

## zetacored tx distribution update-params-proposal

Submit a proposal to update distribution module params. Note: the entire params must be provided.

```
zetacored tx distribution update-params-proposal [params] [flags]
```

### Examples

```
zetacored tx distribution update-params-proposal '{ "community_tax": "20000", "base_proposer_reward": "0", "bonus_proposer_reward": "0", "withdraw_addr_enabled": true }'
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for update-params-proposal
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx distribution](#zetacored-tx-distribution)	 - Distribution transactions subcommands

## zetacored tx distribution withdraw-all-rewards

withdraw all delegations rewards for a delegator

### Synopsis

Withdraw all rewards for a single delegator.
Note that if you use this command with --broadcast-mode=sync or --broadcast-mode=async, the max-msgs flag will automatically be set to 0.

Example:
$ zetacored tx distribution withdraw-all-rewards --from mykey

```
zetacored tx distribution withdraw-all-rewards [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for withdraw-all-rewards
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --max-msgs int                Limit the number of messages per tx (0 for unlimited)
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx distribution](#zetacored-tx-distribution)	 - Distribution transactions subcommands

## zetacored tx distribution withdraw-rewards

Withdraw rewards from a given delegation address, and optionally withdraw validator commission if the delegation address given is a validator operator

### Synopsis

Withdraw rewards from a given delegation address,
and optionally withdraw validator commission if the delegation address given is a validator operator.

Example:
$ zetacored tx distribution withdraw-rewards zetavaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj --from mykey
$ zetacored tx distribution withdraw-rewards zetavaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj --from mykey --commission

```
zetacored tx distribution withdraw-rewards [validator-addr] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --commission                  Withdraw the validator's commission in addition to the rewards
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for withdraw-rewards
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx distribution](#zetacored-tx-distribution)	 - Distribution transactions subcommands

## zetacored tx distribution withdraw-validator-commission

Withdraw commissions from a validator address (must be a validator operator)

```
zetacored tx distribution withdraw-validator-commission [validator-addr] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for withdraw-validator-commission
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx distribution](#zetacored-tx-distribution)	 - Distribution transactions subcommands

## zetacored tx emissions

emissions transactions subcommands

```
zetacored tx emissions [flags]
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

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands
* [zetacored tx emissions withdraw-emission](#zetacored-tx-emissions-withdraw-emission)	 - create a new withdrawEmission

## zetacored tx emissions withdraw-emission

create a new withdrawEmission

```
zetacored tx emissions withdraw-emission [amount] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for withdraw-emission
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx emissions](#zetacored-tx-emissions)	 - emissions transactions subcommands

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
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for encode
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

## zetacored tx evidence

Evidence transaction subcommands

```
zetacored tx evidence [flags]
```

### Options

```
  -h, --help   help for evidence
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

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands

## zetacored tx evm

evm subcommands

```
zetacored tx evm [flags]
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

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands
* [zetacored tx evm raw](#zetacored-tx-evm-raw)	 - Build cosmos transaction from raw ethereum transaction
* [zetacored tx evm send](#zetacored-tx-evm-send)	 - Send funds from one account to another.

## zetacored tx evm raw

Build cosmos transaction from raw ethereum transaction

```
zetacored tx evm raw TX_HEX [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for raw
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx evm](#zetacored-tx-evm)	 - evm subcommands

## zetacored tx evm send

Send funds from one account to another.

### Synopsis

Send funds from one account to another. Both 0x and bech32 addresses
may be used.
Note, the '--from' flag is ignored as it is implied from [from_key_or_address].
When using '--dry-run' a key name cannot be used, only an 0x or bech32 address.


```
zetacored tx evm send [from_key_or_address] [to_address] [amount] [flags]
```

### Examples

```
evmd tx evm send 0x7cB61D4117AE31a12E393a1Cfa3BaC666481D02E 0xA2A8B87390F8F2D188242656BFb6852914073D06 10utoken
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for send
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx evm](#zetacored-tx-evm)	 - evm subcommands

## zetacored tx feemarket

Transactions commands for the feemarket module

```
zetacored tx feemarket [flags]
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

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands
* [zetacored tx feemarket update-params](#zetacored-tx-feemarket-update-params)	 - Execute the UpdateParams RPC method

## zetacored tx feemarket update-params

Execute the UpdateParams RPC method

```
zetacored tx feemarket update-params [flags]
```

### Options

```
  -a, --account-number uint                            The account number of the signing account (offline mode only)
      --aux                                            Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string                          Transaction broadcasting mode (sync|async) 
      --chain-id string                                The network chain ID
      --dry-run                                        ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string                             Fee granter grants fees for the transaction
      --fee-payer string                               Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                                    Fees to pay along with transaction; eg: 10uatom
      --from string                                    Name or address of private key with which to sign
      --gas string                                     gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float                           adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string                              Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only                                  Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                                           help for update-params
      --keyring-backend string                         Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string                             The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                                         Use a connected Ledger device
      --node string                                    [host]:[port] to CometBFT rpc interface for this chain 
      --note string                                    Note to add a description to the transaction (previously --memo)
      --offline                                        Offline mode (does not allow any online functionality)
  -o, --output string                                  Output format (text|json) 
      --params cosmos.evm.feemarket.v1.Params (json)   
  -s, --sequence uint                                  The sequence number of the signing account (offline mode only)
      --sign-mode string                               Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration                      TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint                            DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                                     Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                                      Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                                            Skip tx broadcasting prompt confirmation
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

* [zetacored tx feemarket](#zetacored-tx-feemarket)	 - Transactions commands for the feemarket module

## zetacored tx fungible

fungible transactions subcommands

```
zetacored tx fungible [flags]
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

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands
* [zetacored tx fungible deploy-fungible-coin-zrc-4](#zetacored-tx-fungible-deploy-fungible-coin-zrc-4)	 - Broadcast message DeployFungibleCoinZRC20
* [zetacored tx fungible deploy-system-contracts](#zetacored-tx-fungible-deploy-system-contracts)	 - Broadcast message SystemContracts
* [zetacored tx fungible pause-zrc20](#zetacored-tx-fungible-pause-zrc20)	 - Broadcast message PauseZRC20
* [zetacored tx fungible remove-foreign-coin](#zetacored-tx-fungible-remove-foreign-coin)	 - Broadcast message RemoveForeignCoin
* [zetacored tx fungible unpause-zrc20](#zetacored-tx-fungible-unpause-zrc20)	 - Broadcast message UnpauseZRC20
* [zetacored tx fungible update-contract-bytecode](#zetacored-tx-fungible-update-contract-bytecode)	 - Broadcast message UpdateContractBytecode
* [zetacored tx fungible update-gateway-contract](#zetacored-tx-fungible-update-gateway-contract)	 - Broadcast message UpdateGatewayContract to update the gateway contract address
* [zetacored tx fungible update-system-contract](#zetacored-tx-fungible-update-system-contract)	 - Broadcast message UpdateSystemContract
* [zetacored tx fungible update-zrc20-liquidity-cap](#zetacored-tx-fungible-update-zrc20-liquidity-cap)	 - Broadcast message UpdateZRC20LiquidityCap
* [zetacored tx fungible update-zrc20-withdraw-fee](#zetacored-tx-fungible-update-zrc20-withdraw-fee)	 - Broadcast message UpdateZRC20WithdrawFee

## zetacored tx fungible deploy-fungible-coin-zrc-4

Broadcast message DeployFungibleCoinZRC20

```
zetacored tx fungible deploy-fungible-coin-zrc-4 [erc-20] [foreign-chain] [decimals] [name] [symbol] [coin-type] [gas-limit] [liquidity-cap] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for deploy-fungible-coin-zrc-4
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx fungible](#zetacored-tx-fungible)	 - fungible transactions subcommands

## zetacored tx fungible deploy-system-contracts

Broadcast message SystemContracts

```
zetacored tx fungible deploy-system-contracts [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for deploy-system-contracts
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx fungible](#zetacored-tx-fungible)	 - fungible transactions subcommands

## zetacored tx fungible pause-zrc20

Broadcast message PauseZRC20

```
zetacored tx fungible pause-zrc20 [contractAddress1, contractAddress2, ...] [flags]
```

### Examples

```
zetacored tx fungible pause-zrc20 "0xece40cbB54d65282c4623f141c4a8a0bE7D6AdEc, 0xece40cbB54d65282c4623f141c4a8a0bEjgksncf" 
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for pause-zrc20
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx fungible](#zetacored-tx-fungible)	 - fungible transactions subcommands

## zetacored tx fungible remove-foreign-coin

Broadcast message RemoveForeignCoin

```
zetacored tx fungible remove-foreign-coin [name] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for remove-foreign-coin
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx fungible](#zetacored-tx-fungible)	 - fungible transactions subcommands

## zetacored tx fungible unpause-zrc20

Broadcast message UnpauseZRC20

```
zetacored tx fungible unpause-zrc20 [contractAddress1, contractAddress2, ...] [flags]
```

### Examples

```
zetacored tx fungible unpause-zrc20 "0xece40cbB54d65282c4623f141c4a8a0bE7D6AdEc, 0xece40cbB54d65282c4623f141c4a8a0bEjgksncf" 
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for unpause-zrc20
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx fungible](#zetacored-tx-fungible)	 - fungible transactions subcommands

## zetacored tx fungible update-contract-bytecode

Broadcast message UpdateContractBytecode

```
zetacored tx fungible update-contract-bytecode [contract-address] [new-code-hash] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for update-contract-bytecode
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx fungible](#zetacored-tx-fungible)	 - fungible transactions subcommands

## zetacored tx fungible update-gateway-contract

Broadcast message UpdateGatewayContract to update the gateway contract address

```
zetacored tx fungible update-gateway-contract [contract-address] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for update-gateway-contract
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx fungible](#zetacored-tx-fungible)	 - fungible transactions subcommands

## zetacored tx fungible update-system-contract

Broadcast message UpdateSystemContract

```
zetacored tx fungible update-system-contract [contract-address]  [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for update-system-contract
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx fungible](#zetacored-tx-fungible)	 - fungible transactions subcommands

## zetacored tx fungible update-zrc20-liquidity-cap

Broadcast message UpdateZRC20LiquidityCap

```
zetacored tx fungible update-zrc20-liquidity-cap [zrc20] [liquidity-cap] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for update-zrc20-liquidity-cap
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx fungible](#zetacored-tx-fungible)	 - fungible transactions subcommands

## zetacored tx fungible update-zrc20-withdraw-fee

Broadcast message UpdateZRC20WithdrawFee

```
zetacored tx fungible update-zrc20-withdraw-fee [contractAddress] [newWithdrawFee] [newGasLimit] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for update-zrc20-withdraw-fee
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx fungible](#zetacored-tx-fungible)	 - fungible transactions subcommands

## zetacored tx gov

Governance transactions subcommands

```
zetacored tx gov [flags]
```

### Options

```
  -h, --help   help for gov
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

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands
* [zetacored tx gov cancel-proposal](#zetacored-tx-gov-cancel-proposal)	 - Cancel governance proposal before the voting period ends. Must be signed by the proposal creator.
* [zetacored tx gov deposit](#zetacored-tx-gov-deposit)	 - Deposit tokens for an active proposal
* [zetacored tx gov draft-proposal](#zetacored-tx-gov-draft-proposal)	 - Generate a draft proposal json file. The generated proposal json contains only one message (skeleton).
* [zetacored tx gov submit-legacy-proposal](#zetacored-tx-gov-submit-legacy-proposal)	 - Submit a legacy proposal along with an initial deposit
* [zetacored tx gov submit-proposal](#zetacored-tx-gov-submit-proposal)	 - Submit a proposal along with some messages, metadata and deposit
* [zetacored tx gov vote](#zetacored-tx-gov-vote)	 - Vote for an active proposal, options: yes/no/no_with_veto/abstain
* [zetacored tx gov weighted-vote](#zetacored-tx-gov-weighted-vote)	 - Vote for an active proposal, options: yes/no/no_with_veto/abstain

## zetacored tx gov cancel-proposal

Cancel governance proposal before the voting period ends. Must be signed by the proposal creator.

```
zetacored tx gov cancel-proposal [proposal-id] [flags]
```

### Examples

```
$ zetacored tx gov cancel-proposal 1 --from mykey
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for cancel-proposal
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx gov](#zetacored-tx-gov)	 - Governance transactions subcommands

## zetacored tx gov deposit

Deposit tokens for an active proposal

### Synopsis

Submit a deposit for an active proposal. You can
find the proposal-id by running "zetacored query gov proposals".

Example:
$ zetacored tx gov deposit 1 10stake --from mykey

```
zetacored tx gov deposit [proposal-id] [deposit] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for deposit
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx gov](#zetacored-tx-gov)	 - Governance transactions subcommands

## zetacored tx gov draft-proposal

Generate a draft proposal json file. The generated proposal json contains only one message (skeleton).

```
zetacored tx gov draft-proposal [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for draft-proposal
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --skip-metadata               skip metadata prompt
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx gov](#zetacored-tx-gov)	 - Governance transactions subcommands

## zetacored tx gov submit-legacy-proposal

Submit a legacy proposal along with an initial deposit

### Synopsis

Submit a legacy proposal along with an initial deposit.
Proposal title, description, type and deposit can be given directly or through a proposal JSON file.

Example:
$ zetacored tx gov submit-legacy-proposal --proposal="path/to/proposal.json" --from mykey

Where proposal.json contains:

{
  "title": "Test Proposal",
  "description": "My awesome proposal",
  "type": "Text",
  "deposit": "10test"
}

Which is equivalent to:

$ zetacored tx gov submit-legacy-proposal --title="Test Proposal" --description="My awesome proposal" --type="Text" --deposit="10test" --from mykey

```
zetacored tx gov submit-legacy-proposal [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --deposit string              The proposal deposit
      --description string          The proposal description
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for submit-legacy-proposal
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
      --proposal string             Proposal file path (if this path is given, other proposal flags are ignored)
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --title string                The proposal title
      --type string                 The proposal Type
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx gov](#zetacored-tx-gov)	 - Governance transactions subcommands

## zetacored tx gov submit-proposal

Submit a proposal along with some messages, metadata and deposit

### Synopsis

Submit a proposal along with some messages, metadata and deposit.
They should be defined in a JSON file.

Example:
$ zetacored tx gov submit-proposal path/to/proposal.json

Where proposal.json contains:

{
  // array of proto-JSON-encoded sdk.Msgs
  "messages": [
    {
      "@type": "/cosmos.bank.v1beta1.MsgSend",
      "from_address": "cosmos1...",
      "to_address": "cosmos1...",
      "amount":[{"denom": "stake","amount": "10"}]
    }
  ],
  // metadata can be any of base64 encoded, raw text, stringified json, IPFS link to json
  // see below for example metadata
  "metadata": "4pIMOgIGx1vZGU=",
  "deposit": "10stake",
  "title": "My proposal",
  "summary": "A short summary of my proposal",
  "expedited": false
}

metadata example: 
{
	"title": "",
	"authors": [""],
	"summary": "",
	"details": "", 
	"proposal_forum_url": "",
	"vote_option_context": "",
}

```
zetacored tx gov submit-proposal [path/to/proposal.json] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for submit-proposal
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx gov](#zetacored-tx-gov)	 - Governance transactions subcommands

## zetacored tx gov vote

Vote for an active proposal, options: yes/no/no_with_veto/abstain

### Synopsis

Submit a vote for an active proposal. You can
find the proposal-id by running "zetacored query gov proposals".

Example:
$ zetacored tx gov vote 1 yes --from mykey

```
zetacored tx gov vote [proposal-id] [option] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for vote
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --metadata string             Specify metadata of the vote
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx gov](#zetacored-tx-gov)	 - Governance transactions subcommands

## zetacored tx gov weighted-vote

Vote for an active proposal, options: yes/no/no_with_veto/abstain

### Synopsis

Submit a vote for an active proposal. You can
find the proposal-id by running "zetacored query gov proposals".

Example:
$ zetacored tx gov weighted-vote 1 yes=0.6,no=0.3,abstain=0.05,no_with_veto=0.05 --from mykey

```
zetacored tx gov weighted-vote [proposal-id] [weighted-options] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for weighted-vote
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --metadata string             Specify metadata of the weighted vote
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx gov](#zetacored-tx-gov)	 - Governance transactions subcommands

## zetacored tx group

Group transaction subcommands

```
zetacored tx group [flags]
```

### Options

```
  -h, --help   help for group
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

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands
* [zetacored tx group create-group](#zetacored-tx-group-create-group)	 - Create a group which is an aggregation of member accounts with associated weights and an administrator account.
* [zetacored tx group create-group-policy](#zetacored-tx-group-create-group-policy)	 - Create a group policy which is an account associated with a group and a decision policy. Note, the '--from' flag is ignored as it is implied from [admin].
* [zetacored tx group create-group-with-policy](#zetacored-tx-group-create-group-with-policy)	 - Create a group with policy which is an aggregation of member accounts with associated weights, an administrator account and decision policy.
* [zetacored tx group draft-proposal](#zetacored-tx-group-draft-proposal)	 - Generate a draft proposal json file. The generated proposal json contains only one message (skeleton).
* [zetacored tx group exec](#zetacored-tx-group-exec)	 - Execute a proposal
* [zetacored tx group leave-group](#zetacored-tx-group-leave-group)	 - Remove member from the group
* [zetacored tx group submit-proposal](#zetacored-tx-group-submit-proposal)	 - Submit a new proposal
* [zetacored tx group update-group-admin](#zetacored-tx-group-update-group-admin)	 - Update a group's admin
* [zetacored tx group update-group-members](#zetacored-tx-group-update-group-members)	 - Update a group's members. Set a member's weight to "0" to delete it.
* [zetacored tx group update-group-metadata](#zetacored-tx-group-update-group-metadata)	 - Update a group's metadata
* [zetacored tx group update-group-policy-admin](#zetacored-tx-group-update-group-policy-admin)	 - Update a group policy admin
* [zetacored tx group update-group-policy-decision-policy](#zetacored-tx-group-update-group-policy-decision-policy)	 - Update a group policy's decision policy
* [zetacored tx group update-group-policy-metadata](#zetacored-tx-group-update-group-policy-metadata)	 - Update a group policy metadata
* [zetacored tx group vote](#zetacored-tx-group-vote)	 - Vote on a proposal
* [zetacored tx group withdraw-proposal](#zetacored-tx-group-withdraw-proposal)	 - Withdraw a submitted proposal

## zetacored tx group create-group

Create a group which is an aggregation of member accounts with associated weights and an administrator account.

### Synopsis

Create a group which is an aggregation of member accounts with associated weights and an administrator account.
Note, the '--from' flag is ignored as it is implied from [admin]. Members accounts can be given through a members JSON file that contains an array of members.

```
zetacored tx group create-group [admin] [metadata] [members-json-file] [flags]
```

### Examples

```

zetacored tx group create-group [admin] [metadata] [members-json-file]

Where members.json contains:

{
	"members": [
		{
			"address": "addr1",
			"weight": "1",
			"metadata": "some metadata"
		},
		{
			"address": "addr2",
			"weight": "1",
			"metadata": "some metadata"
		}
	]
}
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for create-group
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx group](#zetacored-tx-group)	 - Group transaction subcommands

## zetacored tx group create-group-policy

Create a group policy which is an account associated with a group and a decision policy. Note, the '--from' flag is ignored as it is implied from [admin].

```
zetacored tx group create-group-policy [admin] [group-id] [metadata] [decision-policy-json-file] [flags]
```

### Examples

```

zetacored tx group create-group-policy [admin] [group-id] [metadata] policy.json

where policy.json contains:

{
    "@type": "/cosmos.group.v1.ThresholdDecisionPolicy",
    "threshold": "1",
    "windows": {
        "voting_period": "120h",
        "min_execution_period": "0s"
    }
}

Here, we can use percentage decision policy when needed, where 0 < percentage <= 1:

{
    "@type": "/cosmos.group.v1.PercentageDecisionPolicy",
    "percentage": "0.5",
    "windows": {
        "voting_period": "120h",
        "min_execution_period": "0s"
    }
}
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for create-group-policy
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx group](#zetacored-tx-group)	 - Group transaction subcommands

## zetacored tx group create-group-with-policy

Create a group with policy which is an aggregation of member accounts with associated weights, an administrator account and decision policy.

### Synopsis

Create a group with policy which is an aggregation of member accounts with associated weights,
an administrator account and decision policy. Note, the '--from' flag is ignored as it is implied from [admin].
Members accounts can be given through a members JSON file that contains an array of members.
If group-policy-as-admin flag is set to true, the admin of the newly created group and group policy is set with the group policy address itself.

```
zetacored tx group create-group-with-policy [admin] [group-metadata] [group-policy-metadata] [members-json-file] [decision-policy-json-file] [flags]
```

### Examples

```

zetacored tx group create-group-with-policy [admin] [group-metadata] [group-policy-metadata] members.json policy.json

where members.json contains:

{
	"members": [
		{
			"address": "addr1",
			"weight": "1",
			"metadata": "some metadata"
		},
		{
			"address": "addr2",
			"weight": "1",
			"metadata": "some metadata"
		}
	]
}

and policy.json contains:

{
    "@type": "/cosmos.group.v1.ThresholdDecisionPolicy",
    "threshold": "1",
    "windows": {
        "voting_period": "120h",
        "min_execution_period": "0s"
    }
}

```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
      --group-policy-as-admin       Sets admin of the newly created group and group policy with group policy address itself when true
  -h, --help                        help for create-group-with-policy
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx group](#zetacored-tx-group)	 - Group transaction subcommands

## zetacored tx group draft-proposal

Generate a draft proposal json file. The generated proposal json contains only one message (skeleton).

```
zetacored tx group draft-proposal [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for draft-proposal
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --skip-metadata               skip metadata prompt
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx group](#zetacored-tx-group)	 - Group transaction subcommands

## zetacored tx group exec

Execute a proposal

```
zetacored tx group exec [proposal-id] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for exec
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx group](#zetacored-tx-group)	 - Group transaction subcommands

## zetacored tx group leave-group

Remove member from the group

### Synopsis

Remove member from the group

Parameters:
		   group-id: unique id of the group
		   member-address: account address of the group member
		   Note, the '--from' flag is ignored as it is implied from [member-address]
		

```
zetacored tx group leave-group [member-address] [group-id] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for leave-group
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx group](#zetacored-tx-group)	 - Group transaction subcommands

## zetacored tx group submit-proposal

Submit a new proposal

### Synopsis

Submit a new proposal.
Parameters:
			msg_tx_json_file: path to json file with messages that will be executed if the proposal is accepted.

```
zetacored tx group submit-proposal [proposal_json_file] [flags]
```

### Examples

```

zetacored tx group submit-proposal path/to/proposal.json
	
	Where proposal.json contains:

{
	"group_policy_address": "cosmos1...",
	// array of proto-JSON-encoded sdk.Msgs
	"messages": [
	{
		"@type": "/cosmos.bank.v1beta1.MsgSend",
		"from_address": "cosmos1...",
		"to_address": "cosmos1...",
		"amount":[{"denom": "stake","amount": "10"}]
	}
	],
	// metadata can be any of base64 encoded, raw text, stringified json, IPFS link to json
	// see below for example metadata
	"metadata": "4pIMOgIGx1vZGU=", // base64-encoded metadata
	"title": "My proposal",
	"summary": "This is a proposal to send 10 stake to cosmos1...",
	"proposers": ["cosmos1...", "cosmos1..."],
}

metadata example: 
{
	"title": "",
	"authors": [""],
	"summary": "",
	"details": "", 
	"proposal_forum_url": "",
	"vote_option_context": "",
} 

```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --exec string                 Set to 1 or 'try' to try to execute proposal immediately after creation (proposers signatures are considered as Yes votes)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for submit-proposal
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx group](#zetacored-tx-group)	 - Group transaction subcommands

## zetacored tx group update-group-admin

Update a group's admin

```
zetacored tx group update-group-admin [admin] [group-id] [new-admin] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for update-group-admin
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx group](#zetacored-tx-group)	 - Group transaction subcommands

## zetacored tx group update-group-members

Update a group's members. Set a member's weight to "0" to delete it.

```
zetacored tx group update-group-members [admin] [group-id] [members-json-file] [flags]
```

### Examples

```

zetacored tx group update-group-members [admin] [group-id] [members-json-file]

Where members.json contains:

{
	"members": [
		{
			"address": "addr1",
			"weight": "1",
			"metadata": "some new metadata"
		},
		{
			"address": "addr2",
			"weight": "0",
			"metadata": "some metadata"
		}
	]
}

Set a member's weight to "0" to delete it.

```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for update-group-members
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx group](#zetacored-tx-group)	 - Group transaction subcommands

## zetacored tx group update-group-metadata

Update a group's metadata

```
zetacored tx group update-group-metadata [admin] [group-id] [metadata] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for update-group-metadata
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx group](#zetacored-tx-group)	 - Group transaction subcommands

## zetacored tx group update-group-policy-admin

Update a group policy admin

```
zetacored tx group update-group-policy-admin [admin] [group-policy-account] [new-admin] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for update-group-policy-admin
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx group](#zetacored-tx-group)	 - Group transaction subcommands

## zetacored tx group update-group-policy-decision-policy

Update a group policy's decision policy

```
zetacored tx group update-group-policy-decision-policy [admin] [group-policy-account] [decision-policy-json-file] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for update-group-policy-decision-policy
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx group](#zetacored-tx-group)	 - Group transaction subcommands

## zetacored tx group update-group-policy-metadata

Update a group policy metadata

```
zetacored tx group update-group-policy-metadata [admin] [group-policy-account] [new-metadata] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for update-group-policy-metadata
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx group](#zetacored-tx-group)	 - Group transaction subcommands

## zetacored tx group vote

Vote on a proposal

### Synopsis

Vote on a proposal.

Parameters:
			proposal-id: unique ID of the proposal
			voter: voter account addresses.
			vote-option: choice of the voter(s)
				VOTE_OPTION_UNSPECIFIED: no-op
				VOTE_OPTION_NO: no
				VOTE_OPTION_YES: yes
				VOTE_OPTION_ABSTAIN: abstain
				VOTE_OPTION_NO_WITH_VETO: no-with-veto
			Metadata: metadata for the vote


```
zetacored tx group vote [proposal-id] [voter] [vote-option] [metadata] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --exec string                 Set to 1 to try to execute proposal immediately after voting
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for vote
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx group](#zetacored-tx-group)	 - Group transaction subcommands

## zetacored tx group withdraw-proposal

Withdraw a submitted proposal

### Synopsis

Withdraw a submitted proposal.

Parameters:
			proposal-id: unique ID of the proposal.
			group-policy-admin-or-proposer: either admin of the group policy or one the proposer of the proposal.
			Note: --from flag will be ignored here.


```
zetacored tx group withdraw-proposal [proposal-id] [group-policy-admin-or-proposer] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for withdraw-proposal
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx group](#zetacored-tx-group)	 - Group transaction subcommands

## zetacored tx lightclient

lightclient transactions subcommands

```
zetacored tx lightclient [flags]
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

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands
* [zetacored tx lightclient disable-header-verification](#zetacored-tx-lightclient-disable-header-verification)	 - Disable header verification for the list of chains separated by comma
* [zetacored tx lightclient enable-header-verification](#zetacored-tx-lightclient-enable-header-verification)	 - Enable verification for the list of chains separated by comma

## zetacored tx lightclient disable-header-verification

Disable header verification for the list of chains separated by comma

### Synopsis

Provide a list of chain ids separated by comma to disable block header verification for the specified chain ids.

  				Example:
                    To disable verification flags for chain ids 1 and 56
					zetacored tx lightclient disable-header-verification "1,56"
				

```
zetacored tx lightclient disable-header-verification [list of chain-id] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for disable-header-verification
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx lightclient](#zetacored-tx-lightclient)	 - lightclient transactions subcommands

## zetacored tx lightclient enable-header-verification

Enable verification for the list of chains separated by comma

### Synopsis

Provide a list of chain ids separated by comma to enable block header verification for the specified chain ids.

  				Example:
                    To enable verification flags for chain ids 1 and 56
					zetacored tx lightclient enable-header-verification "1,56"
				

```
zetacored tx lightclient enable-header-verification [list of chain-id] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for enable-header-verification
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx lightclient](#zetacored-tx-lightclient)	 - lightclient transactions subcommands

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
      --timeout-duration duration     TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint           DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                    Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                     Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
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
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for multisign-batch
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --multisig string             Address of the multisig account that the transaction signs on behalf of
      --no-auto-increment           disable sequence auto increment
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
      --output-document string      The document is written to the given file instead of STDOUT
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

## zetacored tx observer

observer transactions subcommands

```
zetacored tx observer [flags]
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

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands
* [zetacored tx observer add-observer](#zetacored-tx-observer-add-observer)	 - Broadcast message add-observer
* [zetacored tx observer disable-cctx](#zetacored-tx-observer-disable-cctx)	 - Disable inbound and outbound for CCTX
* [zetacored tx observer disable-fast-confirmation](#zetacored-tx-observer-disable-fast-confirmation)	 - Disable fast confirmation for the given chain ID
* [zetacored tx observer enable-cctx](#zetacored-tx-observer-enable-cctx)	 - Enable inbound and outbound for CCTX
* [zetacored tx observer encode](#zetacored-tx-observer-encode)	 - Encode a json string into hex
* [zetacored tx observer remove-chain-params](#zetacored-tx-observer-remove-chain-params)	 - Broadcast message to remove chain params
* [zetacored tx observer reset-chain-nonces](#zetacored-tx-observer-reset-chain-nonces)	 - Broadcast message to reset chain nonces
* [zetacored tx observer update-chain-params](#zetacored-tx-observer-update-chain-params)	 - Broadcast message updateChainParams
* [zetacored tx observer update-gas-price-increase-flags](#zetacored-tx-observer-update-gas-price-increase-flags)	 - Update the gas price increase flags
* [zetacored tx observer update-keygen](#zetacored-tx-observer-update-keygen)	 - command to update the keygen block via a group proposal
* [zetacored tx observer update-observer](#zetacored-tx-observer-update-observer)	 - Broadcast message add-observer
* [zetacored tx observer update-operational-flags](#zetacored-tx-observer-update-operational-flags)	 - Broadcast message UpdateOperationalFlags
* [zetacored tx observer vote-blame](#zetacored-tx-observer-vote-blame)	 - Broadcast message vote-blame
* [zetacored tx observer vote-tss](#zetacored-tx-observer-vote-tss)	 - Vote for a new TSS creation

## zetacored tx observer add-observer

Broadcast message add-observer

```
zetacored tx observer add-observer [observer-address] [zetaclient-grantee-pubkey] [add_node_account_only] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for add-observer
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx observer](#zetacored-tx-observer)	 - observer transactions subcommands

## zetacored tx observer disable-cctx

Disable inbound and outbound for CCTX

```
zetacored tx observer disable-cctx [disable-inbound] [disable-outbound] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for disable-cctx
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx observer](#zetacored-tx-observer)	 - observer transactions subcommands

## zetacored tx observer disable-fast-confirmation

Disable fast confirmation for the given chain ID

```
zetacored tx observer disable-fast-confirmation [chain-id] [flags]
```

### Examples

```
zetacored tx observer disable-fast-confirmation 1
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for disable-fast-confirmation
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx observer](#zetacored-tx-observer)	 - observer transactions subcommands

## zetacored tx observer enable-cctx

Enable inbound and outbound for CCTX

```
zetacored tx observer enable-cctx [enable-inbound] [enable-outbound] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for enable-cctx
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx observer](#zetacored-tx-observer)	 - observer transactions subcommands

## zetacored tx observer encode

Encode a json string into hex

```
zetacored tx observer encode [file.json] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for encode
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx observer](#zetacored-tx-observer)	 - observer transactions subcommands

## zetacored tx observer remove-chain-params

Broadcast message to remove chain params

```
zetacored tx observer remove-chain-params [chain-id] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for remove-chain-params
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx observer](#zetacored-tx-observer)	 - observer transactions subcommands

## zetacored tx observer reset-chain-nonces

Broadcast message to reset chain nonces

```
zetacored tx observer reset-chain-nonces [chain-id] [chain-nonce-low] [chain-nonce-high] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for reset-chain-nonces
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx observer](#zetacored-tx-observer)	 - observer transactions subcommands

## zetacored tx observer update-chain-params

Broadcast message updateChainParams

```
zetacored tx observer update-chain-params [chain-id] [client-params.json] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for update-chain-params
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx observer](#zetacored-tx-observer)	 - observer transactions subcommands

## zetacored tx observer update-gas-price-increase-flags

Update the gas price increase flags

```
zetacored tx observer update-gas-price-increase-flags [epochLength] [retryInterval] [gasPriceIncreasePercent] [gasPriceIncreaseMax] [maxPendingCctxs] [retryIntervalBTC] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for update-gas-price-increase-flags
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx observer](#zetacored-tx-observer)	 - observer transactions subcommands

## zetacored tx observer update-keygen

command to update the keygen block via a group proposal

```
zetacored tx observer update-keygen [block] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for update-keygen
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx observer](#zetacored-tx-observer)	 - observer transactions subcommands

## zetacored tx observer update-observer

Broadcast message add-observer

```
zetacored tx observer update-observer [old-observer-address] [new-observer-address] [update-reason] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for update-observer
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx observer](#zetacored-tx-observer)	 - observer transactions subcommands

## zetacored tx observer update-operational-flags

Broadcast message UpdateOperationalFlags

```
zetacored tx observer update-operational-flags [flags]
```

### Options

```
  -a, --account-number uint                 The account number of the signing account (offline mode only)
      --aux                                 Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string               Transaction broadcasting mode (sync|async) 
      --chain-id string                     The network chain ID
      --dry-run                             ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string                  Fee granter grants fees for the transaction
      --fee-payer string                    Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                         Fees to pay along with transaction; eg: 10uatom
      --file string                         Path to a JSON file containing OperationalFlags
      --from string                         Name or address of private key with which to sign
      --gas string                          gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float                adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string                   Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only                       Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                                help for update-operational-flags
      --keyring-backend string              Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string                  The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                              Use a connected Ledger device
      --node string                         [host]:[port] to CometBFT rpc interface for this chain 
      --note string                         Note to add a description to the transaction (previously --memo)
      --offline                             Offline mode (does not allow any online functionality)
  -o, --output string                       Output format (text|json) 
      --restart-height int                  Height for a coordinated zetaclient restart
  -s, --sequence uint                       The sequence number of the signing account (offline mode only)
      --sign-mode string                    Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --signer-block-time-offset duration   Offset from the zetacore block time to initiate signing
      --timeout-duration duration           TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint                 DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                          Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                           Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                                 Skip tx broadcasting prompt confirmation
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

* [zetacored tx observer](#zetacored-tx-observer)	 - observer transactions subcommands

## zetacored tx observer vote-blame

Broadcast message vote-blame

```
zetacored tx observer vote-blame [chain-id] [index] [failure-reason] [node-list] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for vote-blame
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx observer](#zetacored-tx-observer)	 - observer transactions subcommands

## zetacored tx observer vote-tss

Vote for a new TSS creation

```
zetacored tx observer vote-tss [pubkey] [keygen-block] [status] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for vote-tss
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx observer](#zetacored-tx-observer)	 - observer transactions subcommands

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
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for sign
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --multisig string             Address or key name of the multisig account on behalf of which the transaction shall be signed
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
      --output-document string      The document will be written to the given file instead of STDOUT
      --overwrite                   Overwrite existing signatures with a new one. If disabled, new signature will be appended
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --signature-only              Print only the signatures
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --append                      Combine all message and generate single signed transaction for broadcast.
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for sign-batch
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --multisig string             Address or key name of the multisig account on behalf of which the transaction shall be signed
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
      --output-document string      The document will be written to the given file instead of STDOUT
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --signature-only              Print only the generated signature, then exit
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

## zetacored tx slashing

Transactions commands for the slashing module

```
zetacored tx slashing [flags]
```

### Options

```
  -h, --help   help for slashing
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

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands
* [zetacored tx slashing unjail](#zetacored-tx-slashing-unjail)	 - Unjail a jailed validator
* [zetacored tx slashing update-params-proposal](#zetacored-tx-slashing-update-params-proposal)	 - Submit a proposal to update slashing module params. Note: the entire params must be provided.

## zetacored tx slashing unjail

Unjail a jailed validator

```
zetacored tx slashing unjail [flags]
```

### Examples

```
zetacored tx slashing unjail --from [validator]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for unjail
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx slashing](#zetacored-tx-slashing)	 - Transactions commands for the slashing module

## zetacored tx slashing update-params-proposal

Submit a proposal to update slashing module params. Note: the entire params must be provided.

### Synopsis

Submit a proposal to update slashing module params. Note: the entire params must be provided.
 See the fields to fill in by running `zetacored query slashing params --output json`

```
zetacored tx slashing update-params-proposal [params] [flags]
```

### Examples

```
zetacored tx slashing update-params-proposal '{ "signed_blocks_window": "100", ... }'
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for update-params-proposal
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx slashing](#zetacored-tx-slashing)	 - Transactions commands for the slashing module

## zetacored tx staking

Staking transaction subcommands

```
zetacored tx staking [flags]
```

### Options

```
  -h, --help   help for staking
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

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands
* [zetacored tx staking cancel-unbond](#zetacored-tx-staking-cancel-unbond)	 - Cancel unbonding delegation and delegate back to the validator
* [zetacored tx staking create-validator](#zetacored-tx-staking-create-validator)	 - create new validator initialized with a self-delegation to it
* [zetacored tx staking delegate](#zetacored-tx-staking-delegate)	 - Delegate liquid tokens to a validator
* [zetacored tx staking edit-validator](#zetacored-tx-staking-edit-validator)	 - edit an existing validator account
* [zetacored tx staking redelegate](#zetacored-tx-staking-redelegate)	 - Redelegate illiquid tokens from one validator to another
* [zetacored tx staking unbond](#zetacored-tx-staking-unbond)	 - Unbond shares from a validator

## zetacored tx staking cancel-unbond

Cancel unbonding delegation and delegate back to the validator

### Synopsis

Cancel Unbonding Delegation and delegate back to the validator.

Example:
$ zetacored tx staking cancel-unbond zetavaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj 100stake 2 --from mykey

```
zetacored tx staking cancel-unbond [validator-addr] [amount] [creation-height] [flags]
```

### Examples

```
$ zetacored tx staking cancel-unbond zetavaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj 100stake 2 --from mykey
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for cancel-unbond
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx staking](#zetacored-tx-staking)	 - Staking transaction subcommands

## zetacored tx staking create-validator

create new validator initialized with a self-delegation to it

### Synopsis

Create a new validator initialized with a self-delegation by submitting a JSON file with the new validator details.

```
zetacored tx staking create-validator [path/to/validator.json] [flags]
```

### Examples

```
$ zetacored tx staking create-validator path/to/validator.json --from keyname

Where validator.json contains:

{
	"pubkey": {"@type":"/cosmos.crypto.ed25519.PubKey","key":"oWg2ISpLF405Jcm2vXV+2v4fnjodh6aafuIdeoW+rUw="},
	"amount": "1000000stake",
	"moniker": "myvalidator",
	"identity": "optional identity signature (ex. UPort or Keybase)",
	"website": "validator's (optional) website",
	"security": "validator's (optional) security contact email",
	"details": "validator's (optional) details",
	"commission-rate": "0.1",
	"commission-max-rate": "0.2",
	"commission-max-change-rate": "0.01",
	"min-self-delegation": "1"
}

where we can get the pubkey using "zetacored tendermint show-validator"
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for create-validator
      --ip string                   The node's public IP. It takes effect only when used in combination with --generate-only
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --node-id string              The node's ID
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx staking](#zetacored-tx-staking)	 - Staking transaction subcommands

## zetacored tx staking delegate

Delegate liquid tokens to a validator

### Synopsis

Delegate an amount of liquid coins to a validator from your wallet.

Example:
$ zetacored tx staking delegate cosmosvalopers1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm 1000stake --from mykey

```
zetacored tx staking delegate [validator-addr] [amount] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for delegate
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx staking](#zetacored-tx-staking)	 - Staking transaction subcommands

## zetacored tx staking edit-validator

edit an existing validator account

```
zetacored tx staking edit-validator [flags]
```

### Options

```
  -a, --account-number uint          The account number of the signing account (offline mode only)
      --aux                          Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string        Transaction broadcasting mode (sync|async) 
      --chain-id string              The network chain ID
      --commission-rate string       The new commission rate percentage
      --details string               The validator's (optional) details 
      --dry-run                      ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string           Fee granter grants fees for the transaction
      --fee-payer string             Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                  Fees to pay along with transaction; eg: 10uatom
      --from string                  Name or address of private key with which to sign
      --gas string                   gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float         adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string            Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only                Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                         help for edit-validator
      --identity string              The (optional) identity signature (ex. UPort or Keybase) 
      --keyring-backend string       Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string           The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                       Use a connected Ledger device
      --min-self-delegation string   The minimum self delegation required on the validator
      --new-moniker string           The validator's name 
      --node string                  [host]:[port] to CometBFT rpc interface for this chain 
      --note string                  Note to add a description to the transaction (previously --memo)
      --offline                      Offline mode (does not allow any online functionality)
  -o, --output string                Output format (text|json) 
      --security-contact string      The validator's (optional) security contact email 
  -s, --sequence uint                The sequence number of the signing account (offline mode only)
      --sign-mode string             Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration    TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint          DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                   Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                    Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
      --website string               The validator's (optional) website 
  -y, --yes                          Skip tx broadcasting prompt confirmation
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

* [zetacored tx staking](#zetacored-tx-staking)	 - Staking transaction subcommands

## zetacored tx staking redelegate

Redelegate illiquid tokens from one validator to another

### Synopsis

Redelegate an amount of illiquid staking tokens from one validator to another.

Example:
$ zetacored tx staking redelegate cosmosvalopers1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj cosmosvalopers1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm 100stake --from mykey

```
zetacored tx staking redelegate [src-validator-addr] [dst-validator-addr] [amount] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for redelegate
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx staking](#zetacored-tx-staking)	 - Staking transaction subcommands

## zetacored tx staking unbond

Unbond shares from a validator

### Synopsis

Unbond an amount of bonded shares from a validator.

Example:
$ zetacored tx staking unbond zetavaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj 100stake --from mykey

```
zetacored tx staking unbond [validator-addr] [amount] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for unbond
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx staking](#zetacored-tx-staking)	 - Staking transaction subcommands

## zetacored tx upgrade

Upgrade transaction subcommands

### Options

```
  -h, --help   help for upgrade
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

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands
* [zetacored tx upgrade cancel-software-upgrade](#zetacored-tx-upgrade-cancel-software-upgrade)	 - Cancel the current software upgrade proposal
* [zetacored tx upgrade cancel-upgrade-proposal](#zetacored-tx-upgrade-cancel-upgrade-proposal)	 - Submit a proposal to cancel a planned chain upgrade.
* [zetacored tx upgrade software-upgrade](#zetacored-tx-upgrade-software-upgrade)	 - Submit a software upgrade proposal

## zetacored tx upgrade cancel-software-upgrade

Cancel the current software upgrade proposal

### Synopsis

Cancel a software upgrade along with an initial deposit.

```
zetacored tx upgrade cancel-software-upgrade [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --authority string            The address of the upgrade module authority (defaults to gov)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --deposit string              The deposit to include with the governance proposal
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for cancel-software-upgrade
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --metadata string             The metadata to include with the governance proposal
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --summary string              The summary to include with the governance proposal
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --title string                The title to put on the governance proposal
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx upgrade](#zetacored-tx-upgrade)	 - Upgrade transaction subcommands

## zetacored tx upgrade cancel-upgrade-proposal

Submit a proposal to cancel a planned chain upgrade.

```
zetacored tx upgrade cancel-upgrade-proposal [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for cancel-upgrade-proposal
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx upgrade](#zetacored-tx-upgrade)	 - Upgrade transaction subcommands

## zetacored tx upgrade software-upgrade

Submit a software upgrade proposal

### Synopsis

Submit a software upgrade along with an initial deposit.
Please specify a unique name and height for the upgrade to take effect.
You may include info to reference a binary download link, in a format compatible with: https://docs.cosmos.network/main/tooling/cosmovisor

```
zetacored tx upgrade software-upgrade [name] (--upgrade-height [height]) (--upgrade-info [info]) [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --authority string            The address of the upgrade module authority (defaults to gov)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --daemon-name string          The name of the executable being upgraded (for upgrade-info validation). Default is the DAEMON_NAME env var if set, or else this executable 
      --deposit string              The deposit to include with the governance proposal
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for software-upgrade
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --metadata string             The metadata to include with the governance proposal
      --no-checksum-required        Skip requirement of checksums for binaries in the upgrade info
      --no-validate                 Skip validation of the upgrade info (dangerous!)
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --summary string              The summary to include with the governance proposal
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --title string                The title to put on the governance proposal
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
      --upgrade-height int          The height at which the upgrade must happen
      --upgrade-info string         Info for the upgrade plan such as new version download urls, etc.
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx upgrade](#zetacored-tx-upgrade)	 - Upgrade transaction subcommands

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
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for validate-signatures
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

## zetacored tx vesting

Vesting transaction subcommands

```
zetacored tx vesting [flags]
```

### Options

```
  -h, --help   help for vesting
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

* [zetacored tx](#zetacored-tx)	 - Transactions subcommands
* [zetacored tx vesting create-periodic-vesting-account](#zetacored-tx-vesting-create-periodic-vesting-account)	 - Create a new vesting account funded with an allocation of tokens.
* [zetacored tx vesting create-permanent-locked-account](#zetacored-tx-vesting-create-permanent-locked-account)	 - Create a new permanently locked account funded with an allocation of tokens.
* [zetacored tx vesting create-vesting-account](#zetacored-tx-vesting-create-vesting-account)	 - Create a new vesting account funded with an allocation of tokens.

## zetacored tx vesting create-periodic-vesting-account

Create a new vesting account funded with an allocation of tokens.

### Synopsis

A sequence of coins and period length in seconds. Periods are sequential, in that the duration of of a period only starts at the end of the previous period. The duration of the first period starts upon account creation. For instance, the following periods.json file shows 20 "test" coins vesting 30 days apart from each other.
		Where periods.json contains:

		An array of coin strings and unix epoch times for coins to vest
{ "start_time": 1625204910,
"periods":[
 {
  "coins": "10test",
  "length_seconds":2592000 //30 days
 },
 {
	"coins": "10test",
	"length_seconds":2592000 //30 days
 },
]
	}
		

```
zetacored tx vesting create-periodic-vesting-account [to_address] [periods_json_file] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for create-periodic-vesting-account
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx vesting](#zetacored-tx-vesting)	 - Vesting transaction subcommands

## zetacored tx vesting create-permanent-locked-account

Create a new permanently locked account funded with an allocation of tokens.

### Synopsis

Create a new account funded with an allocation of permanently locked tokens. These
tokens may be used for staking but are non-transferable. Staking rewards will acrue as liquid and transferable
tokens.

```
zetacored tx vesting create-permanent-locked-account [to_address] [amount] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for create-permanent-locked-account
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx vesting](#zetacored-tx-vesting)	 - Vesting transaction subcommands

## zetacored tx vesting create-vesting-account

Create a new vesting account funded with an allocation of tokens.

### Synopsis

Create a new vesting account funded with an allocation of tokens. The
account can either be a delayed or continuous vesting account, which is determined
by the '--delayed' flag. All vesting accounts created will have their start time
set by the committed block's time. The end_time must be provided as a UNIX epoch
timestamp.

```
zetacored tx vesting create-vesting-account [to_address] [amount] [end_time] [flags]
```

### Options

```
  -a, --account-number uint         The account number of the signing account (offline mode only)
      --aux                         Generate aux signer data instead of sending a tx
  -b, --broadcast-mode string       Transaction broadcasting mode (sync|async) 
      --chain-id string             The network chain ID
      --delayed                     Create a delayed vesting account if true
      --dry-run                     ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)
      --fee-granter string          Fee granter grants fees for the transaction
      --fee-payer string            Fee payer pays fees for the transaction instead of deducting from the signer
      --fees string                 Fees to pay along with transaction; eg: 10uatom
      --from string                 Name or address of private key with which to sign
      --gas string                  gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically. Note: "auto" option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of "fees". (default 200000)
      --gas-adjustment float        adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)
      --gas-prices string           Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only               Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase only accessed when providing a key name)
  -h, --help                        help for create-vesting-account
      --keyring-backend string      Select keyring's backend (os|file|kwallet|pass|test|memory) 
      --keyring-dir string          The client Keyring directory; if omitted, the default 'home' directory will be used
      --ledger                      Use a connected Ledger device
      --node string                 [host]:[port] to CometBFT rpc interface for this chain 
      --note string                 Note to add a description to the transaction (previously --memo)
      --offline                     Offline mode (does not allow any online functionality)
  -o, --output string               Output format (text|json) 
  -s, --sequence uint               The sequence number of the signing account (offline mode only)
      --sign-mode string            Choose sign mode (direct|amino-json|direct-aux|textual), this is an advanced feature
      --timeout-duration duration   TimeoutDuration is the duration the transaction will be considered valid in the mempool. The transaction's unordered nonce will be set to the time of transaction creation + the duration value passed. If the transaction is still in the mempool, and the block time has passed the time of submission + TimeoutTimestamp, the transaction will be rejected.
      --timeout-height uint         DEPRECATED: Please use --timeout-duration instead. Set a block timeout height to prevent the tx from being committed past a certain height
      --tip string                  Tip is the amount that is going to be transferred to the fee payer on the target chain. This flag is only valid when used with --aux, and is ignored if the target chain didn't enable the TipDecorator
      --unordered                   Enable unordered transaction delivery; must be used in conjunction with --timeout-duration
  -y, --yes                         Skip tx broadcasting prompt confirmation
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

* [zetacored tx vesting](#zetacored-tx-vesting)	 - Vesting transaction subcommands

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

