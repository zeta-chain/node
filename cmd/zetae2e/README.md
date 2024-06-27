# `zetae2e`

### Basics

`zetae2e` is a CLI tool allowing to quickly test ZetaChain functionality.

`zetae2e` can communicate with ZetaChain, EVM and Bitcoin network to test and track the full workflow of a cross-chain transaction.

### Getting Started

`zetae2e` can be installed with the command:

```go
make install-zetae2e

zetae2e -h
```

### Config

The command takes a config file as input containing RPC nodes to connect to, hotkey wallet information for interaction with networks, and addresses of the deployed contracts to be used.

This is an example of config for interaction with Athens3:

```go
zeta_chain_id: "athens_7001-1"
default_account:
  evm_address: "<your evm address>"
  private_key: "<your evm private key>"
rpcs:
  zevm: "<zevm (ZetaChain) url>, generally using port 8545"
  evm: "<evm url>, generally using port 8545"
  bitcoin:
    host: "<bitcoin rpc url>"
    user: "<bitcoin user>"
    pass: "<bitcoin pass>"
    http_post_mode: true
    disable_tls: true
    params: "<mainnet|testnet3|regnet>"
  zetacore_grpc: "<ZetaChain grpc url>, generally using port 9090"
  zetacore_rpc: "<ZetaChain grpc url>, generally using port 26657"
contracts:
  zevm:
    system_contract: "0xEdf1c3275d13489aCdC6cD6eD246E72458B8795B"
    eth_zrc20: "0x13A0c5930C028511Dc02665E7285134B6d11A5f4"
    erc20_zrc20: "0x0cbe0dF132a6c6B4a2974Fa1b7Fb953CF0Cc798a"
    btc_zrc20: "0x65a45c57636f9BcCeD4fe193A602008578BcA90b"
    uniswap_factory: "0x9fd96203f7b22bCF72d9DCb40ff98302376cE09c"
    uniswap_router: "0x2ca7d64A7EFE2D62A725E2B35Cf7230D6677FfEe"
  evm:
    zeta_eth: "0x0000c304d2934c00db1d51995b9f6996affd17c0"
    connector_eth: "0x00005e3125aba53c5652f9f0ce1a4cf91d8b15ea"
    custody: "0x000047f11c6e42293f433c82473532e869ce4ec5"
    erc20: "0x07865c6e87b9f70255377e024ace6630c1eaa37f"
```

### Bitcoin setup
Interaction with the Bitcoin node will require setting up a specific node tracking the address. It can be set locally following the guide [Using Bitcoin Docker Image for Local Development](https://www.notion.so/Using-Bitcoin-Docker-Image-for-Local-Development-bf7e84c58f22431fb52f17a471997e1f?pvs=21) 

If an error occurs mention that wallets are not loaded, this command can be run to set it up:

```go
zetae2e setup-bitcoin [config]
```

Note: This command is intended for local testing only and won't work on live networks, where a node with an enabled wallet should already be set up.

### Commands

Show the balances of the accounts used on the different networks:

```go
zetae2e balances [config]
```

Show the Bitcoin address (the address is derived from the Ethereum private key, this address must therefore be found to perform the Bitcoin test):

```go
zetae2e bitcoin-address [config]
```

The list of tests to be run can be found by running following command:

```go
zetae2e list-tests
```

Run tests:

Tests can be run using `run` command. Arguments of this command should be in `testName:testArgs` format, where `testArgs` is optional comma-separated list of test arguments. If test arguments are not provided, default arguments will be used. Example of `run` command:

```go
zetae2e run zeta_deposit:2000000000000000000 eth_deposit:2000000000000000000 erc20_deposit:200000 --config cmd/zetae2e/config/local.yml
```

Since cctxs might take a longer time to be processed on live networks, it is highly recommended to use `--verbose` flag to see the current status of the cctx workflow.

### Testing a gas ZRC20 from an EVM chain

Testing a gas token requires the following values to be defined in the config:

```go
zeta_chain_id
default_account:
  evm_address
  private_key
rpcs:
  zevm
  evm
  zetacore_grpc
  zetacore_rpc
contracts:
  zevm:
    system_contract
    eth_zrc20
    uniswap_factory
    uniswap_router
  evm:
    zeta_eth
    connector_eth
    custody: "0x000047f11c6e42293f433c82473532e869ce4ec5"
```

One of the tests can be commented out in case only a deposit or a withdrawal is to be tested.
Testing an ERC20 ZRC20 from an EVM chain

Testing ZRC20 requires the same config as for the gas tokens, but must include the `erc20` field that contains the address of the ERC20 on the evm chain and `erc20_zrc20` on ZetaChain.

```go
zeta_chain_id
accounts:
  evm_address
  evm_priv_key
rpcs:
  zevm
  evm
  zetacore_grpc
  zetacore_rpc
contracts:
  zevm:
    erc20_zrc20
  evm:
		erc20
```

### Testing a ZRC20 from a Bitcoin chain

The following values must be set in the config in order to test Bitcoin functionality

```go
zeta_chain_id
accounts:
  evm_address
  evm_priv_key
rpcs:
  zevm
  bitcoin:
    host
    user
    pass
    http_post_mode
    disable_tls
    params
  zetacore_grpc
  zetacore_rpc
contracts:
  zevm:
    system_contract
    btc_zrc20
    uniswap_factory
    uniswap_router
```

### TODO: message passing