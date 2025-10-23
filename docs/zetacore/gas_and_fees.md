# Gas and Fees in ZetaChain

Gas and fees are fundamental mechanisms in blockchain networks.
Gas represents the computational effort required to execute operations on a blockchain,
while fees are the cost users pay for these resources,
typically calculated by multiplying gas consumed by a gas price.
Without fees, malicious actors could flood the network with spam transactions,
overwhelming validators and halting the network.
Fees make such attacks economically prohibitive,
create efficient resource allocation through market mechanisms,
and compensate validators for their computational work,
ensuring the sustainability of the decentralized infrastructure.

ZetaChain is a generalist interoperability protocol that enables execution of contracts on connected chains,
including EVM and other virtual machine architectures.
This creates multiple distinct contexts and situations where the notions of gas and fees apply differently.
For example, the fee mechanism when directly interacting with the ZetaChain Layer 1 for a smart contract call differs from initiating a cross-chain transaction.
This page summarizes all gas and fee mechanisms across ZetaChain's various components and connected chains.

For simplicity, we will mostly refer to gas and fees as "fees" in the following sections.

## Usages of Fees in ZetaChain

ZetaChain involves fees in several distinct contexts:

* **ZetaChain EVM Transaction Fees**: ZetaChain integrates the EVM module that makes it an EVM-compatible Layer 1 blockchain. Sending EVM transactions to ZetaChain involves using the standard EVM transaction fee mechanism.

* **Connected Chains Withdraw Fees**: The calculation and payment of fees when a user initiates an outgoing cross-chain transaction from ZetaChain to a connected chain.

* **Connected Chains Deposit Fees**: The calculation and payment of fees when a user initiates an incoming cross-chain transaction from a connected chain to ZetaChain.

* **Cosmos Transaction Fees**: As a Layer 1 blockchain built with Cosmos SDK, sending direct Cosmos transactions on ZetaChain involves using the Cosmos gas mechanism for fee calculation and payment.

## ZetaChain EVM Transaction Fees

When users interact with the ZetaChain EVM, for any operations, they will pay fees using the standard EVM fee mechanism.
More information can be found at the [EVM Fees Documentation](https://ethereum.org/developers/docs/gas/).

ZetaChain currently uses Go-Ethereum `v1.15.11` and supports both legacy and EIP-1559 transaction types.

ZetaChain EVM currently doesn't support `EIP-7702` for sponsored transactions.

The base fee parameter can be queried using the ETH JSON-RPC:

```shell
curl -X POST https://zetachain-evm.blockpi.network/v1/rpc/public \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getBlockByNumber",
    "params": ["latest", false],
    "id": 1
  }'
```

Or using the Cast CLI:

```shell
cast block latest --rpc-url https://zetachain-evm.blockpi.network/v1/rpc/public
```

## Connected Chains Withdraw Fees

Users create outgoing cross-chain transactions to connected chains using the withdraw interface of the ZetaChain gateway:
- `withdraw`
- `withdrawAndCall`
- `call`

As functions called on the ZetaChain EVM, the standard EVM transaction fee mechanism applies for executing these functions on ZetaChain.

In addition to the EVM transaction fees, users must also pay additional fees to cover the cost of executing the transaction on the connected chain.

The additional withdraw fees can be queried on the ZRC20 asset contract with the following functions:
```solidity
// for simple withdraws, default gas limit set per zrc20
function withdrawGasFee() external view returns (address, uint256);

// for withdraws that include a call on the connected chain
function withdrawGasFeeWithGasLimit(uint256 gasLimit) external view returns (address, uint256);
```

Users must approve the gateway contract to spend the required fee amount in addition to the withdraw amount using the ERC20 standard.

For each connected chain, the ZetaChain observer set votes on and maintains an updated "gas price".
The gas price and gas limit are two parameters used to calculate the withdraw fees on connected chains.
We use this terminology for all chains for simplicity and consistency, although not all connected chains use the same concepts, as we support non-EVM chains.
The meaning of gas price and gas limit depends on the connected chain.

These differences are summarized below.

### EVM

For EVM-compatible chains, ZetaChain relies on the standard Ethereum gas mechanism.
More information about Ethereum gas can be found in the [Ethereum Gas Documentation](https://ethereum.org/developers/docs/gas/).

The withdraw fee for EVM chains is calculated using the gas price and gas limit parameters:
```
Withdraw Fee = Gas Price × Gas Limit
```

The gas price represents the cost per unit of gas on the connected EVM chain,
while the gas limit represents the maximum amount of gas allocated for executing the transaction on that chain.

**Important Consideration**: The gas limit is fully consumed regardless of whether the cross-chain call uses all the allocated gas.
In the current architecture, any remaining unused gas is sent to a gas stability pool.
This pool is used to increase the gas price for transactions during periods of gas price surges,
helping to ensure transaction execution during network congestion.

### Bitcoin

For Bitcoin transactions, ZetaChain uses a simplified fee model adapted to Bitcoin's UTXO architecture.
More information about Bitcoin transaction fees can be found in the [Bitcoin Fee Documentation](https://developer.bitcoin.org/devguide/transactions.html).

Cross-chain calls are not supported for Bitcoin. All Bitcoin withdraw transactions use a simple transfer mechanism.

A fixed "gas limit" value of **254** is used for all Bitcoin transactions.
The gas price represents the fee rate (satoshis per byte) on the Bitcoin network.

The withdraw fee for Bitcoin is calculated as:
```
Withdraw Fee = Fee Rate × 254
```

ZetaChain supports Replace-By-Fee (RBF) to increase the fee of a transaction in case of a surge in network fees.
However, the gas stability pool must be manually funded to enable fee increases during periods of high congestion.

### Solana

### Sui

### TON

### Reverting Withdraw 

### Aborting Withdraw

## Connected Chains Deposit Fees

### EVM

### Bitcoin

### Solana

### Sui

### TON

### Reverting Deposit

### Aborting Deposit

## Cosmos Transaction Fees

ZetaChain, being a Cosmos SDK-based blockchain, employs the Cosmos gas mechanism for fee calculation and payment for direct Cosmos transactions.
Documentation for Cosmos SDK fees can be found [here](https://docs.cosmos.network/v0.53/learn/beginner/gas-fees).

Direct interactions are used for operations such as staking, ZETA transfers, or governance voting, although these operations can now be performed via the [EVM precompiles](https://evm.cosmos.network/docs/evm/v0.4.x/documentation/smart-contracts/precompiles). Therefore, it is generally rare for users to perform direct Cosmos transactions on ZetaChain. All interactions for cross-chain transactions are done via the EVM interface.

Cosmos transactions are still used for system transactions sent by the observer signers and for administrative operations.

### System Messages

System messages represent the following operations:

- Voting inbound observations
- Voting outbound observations
- Voting gas prices for connected chains
- Adding and removing inbound and outbound trackers

For these system messages, the ZetaClient process uses fixed gas limit values defined in [this file](https://github.com/zeta-chain/node/blob/develop/zetaclient/zetacore/constant.go).

### Gas Limit for Inbound Voting

The gas limit used for inbound voting has some specificities.
There is currently no distinctive message between voting and executing inbound—the vote message reaching sufficient quorum will execute the inbound automatically.
However, the execution of the inbound might trigger a smart contract call that requires an increased amount of gas.
Since higher gas limits result in higher fees in Cosmos transactions, the following mechanism is used to avoid requiring high gas limits for each vote message:

- Each vote message is sent with a base gas limit (currently 500,000)
- When the vote message fails due to out of gas during execution, the ZetaClient process re-sends the vote message with an increased gas limit (currently 7,000,000)

This mechanism will be removed in the future once a distinct message for executing inbound is implemented.

### Administrative Operations

The admnistrative operations can be found at [this page](https://www.zetachain.com/docs/developers/architecture/privileged).

They must be executed using the [group module](https://tutorials.cosmos.network/tutorials/8-understand-sdk-modules/3-group.html).

## Resources

You can find below resources and references for the fee mechanisms of the different technologies that ZetaChain infrastructure relies on.

* [Cosmos SDK Fees Documentation](https://docs.cosmos.network/v0.53/learn/beginner/gas-fees)
* [EVM Fees Documentation](https://ethereum.org/developers/docs/gas/)
* [Bitcoin Transaction documentation](https://developer.bitcoin.org/devguide/transactions.html)
* [Solana Fees Documentation](https://solana.com/docs/core/fees)
* [Sui Fees Documentation](https://docs.sui.io/concepts/tokenomics/gas-pricing)
* [TON Fees Documentation](https://docs.ton.org/v3/documentation/smart-contracts/transaction-fees/fees)
