# `e2e`

`e2e` is a comprehensive suite of E2E tests designed to validate the integration and functionality of the ZetaChain network, particularly its interactions with Bitcoin and EVM (Ethereum Virtual Machine) networks. This tool is essential for ensuring the robustness and reliability of ZetaChain's cross-chain functionalities.

## Packages
The E2E testing project is organized into several packages, each with a specific role:

- `config`: Provides general configuration for E2E tests, including RPC addresses for connected networks, addresses of deployed smart contracts, and account details for test transactions.
- `contracts`: Includes sample Solidity smart contracts used in testing scenarios.
- `runner`: Responsible for executing E2E tests, handling interactions with various network clients.
- `e2etests`: Houses a collection of E2E tests that can be run against the ZetaChain network. Each test is implemented as a separate Go file prefixed with `test_`.
- `txserver`: A minimalistic client for interacting with the ZetaChain RPC interface.
- `utils`: Offers utility functions to facilitate interactions with the different blockchain networks involved in testing.

## Config

The E2E testing suite utilizes a flexible and comprehensive configuration system defined in the config package, which is central to setting up and customizing your test environments. The configuration is structured as follows:

A config YAML file can be provided to the E2E test tool via the `--config` flag. If no config file is provided, the tool will use default values for all configuration parameters.

### Config Structure
- `RPCs`: Defines the RPC endpoints for various networks involved in the testing.
- `Contracts`: Specifies the addresses of pre-deployed smart contracts relevant to the tests.
- `ZetaChainID`: The specific chain ID of the ZetaChain network being tested.

### RPCs Configuration

- `Zevm`: RPC endpoint for the ZetaChain EVM.
- `EVM`: RPC endpoint for the Ethereum network.
- `Bitcoin`: RPC endpoint for the Bitcoin network.
- `ZetaCoreGRPC`: GRPC endpoint for zetacore.
- `ZetaCoreRPC`: RPC endpoint for zetacore.

### Contracts Configuration:

**EVM Contracts**
- `ZetaEthAddress`: Address of Zeta token contract on EVM chain.
- `ConnectorEthAddr`: Address of a connector contract on EVM chain.
- `ERC20`: Address of the ERC20 token contract on EVM chain.

### Config Example

```yaml
rpcs:
  zevm: "http://localhost:8545"
  evm: "http://localhost:8546"
  bitcoin: "http://localhost:18332"
  zetacore_grpc: "localhost:9090"
  zetacore_rpc: "http://localhost:26657"
contracts:
  evm:
    zeta_eth: "0x..."
    connector_eth: "0x..."
    erc20: "0x..."
zeta_chain_id: "zetachain-1"
```

NOTE: config is in progress, contracts on the zEVM must be added

## Debugging

It's possible to debug a single test using Delve debugger.

1. Make sure delve is installed. `go install github.com/go-delve/delve/cmd/dlv@latest`
2. Configure your IDE to use Delve as the debugger. For Goland, you can do the following:
    - Go to "Run" > "Edit Run Configurations"
    - Hit "+" > "Go Remote". Keep port as default (`2345`). Toggle "On Disconnect" > "Stop Delve process"
3. Make sure that localnet is running. For a quick start, you can use `make start-localnet-skip-build`.
   Networks need some time to generate blocks.
4. Run test as following: `./e2e/scripts/debug.sh my_test_name arg1 arg2 arg_n`.
   Example: `./e2e/scripts/debug.sh bitcoin_withdraw_restricted 0.001`
5. Place a breakpoint in the code.
6. Go to the editor's debug panel and hit "Debug" button.

You can also run an alias of `zetae2e run` like so:
```shell
  `./e2e/scripts/run.sh bitcoin_withdraw_restricted 0.001`
```
