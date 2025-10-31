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

ZetaChain involves fees in several distinct contexts, which can be grouped into fees related to local transactions on the ZetaChain L1 and fees related to cross-chain interactions with connected chains.

### Local L1 Fees (ZetaChain)

These fees cover operations executed directly on the ZetaChain Layer 1 (L1) blockchain.
Most users will primarily interact with ZetaChain EVM Transaction Fees.

* **ZetaChain EVM Transaction Fees**: ZetaChain integrates the [EVM module](https://github.com/cosmos/evm) that makes it an EVM-compatible Layer 1 blockchain. Sending EVM transactions to ZetaChain involves using the standard EVM transaction fee mechanism.

* **Cosmos Transaction Fees**: As a Layer 1 blockchain built with Cosmos SDK, sending direct Cosmos transactions on ZetaChain involves using the Cosmos gas mechanism for fee calculation and payment.

### Cross-Chain Fees

These fees are incurred when a user initiates a transaction that bridges assets or data between ZetaChain and a connected chains.

* **Connected Chains Withdraw Fees**: The calculation and payment of fees when a user initiates an outgoing cross-chain transaction from ZetaChain to a connected chain.

* **Connected Chains Deposit Fees**: The calculation and payment of fees when a user initiates an incoming cross-chain transaction from a connected chain to ZetaChain.

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

Or using the [Cast CLI](https://getfoundry.sh):

```shell
cast block latest --rpc-url https://zetachain-evm.blockpi.network/v1/rpc/public
```

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
There is currently no distinctive message between voting and executing inbound, the vote message reaching sufficient quorum will execute the inbound automatically.
However, the execution of the inbound might trigger a smart contract call that requires an increased amount of gas.
Since there is no way to predict which vote will trigger execution, and higher gas limits result in higher fees in Cosmos transactions, the following mechanism is used to avoid requiring high gas limits for each vote message:

- Each vote message is sent with a base gas limit (currently 500,000)
- When the vote message fails due to out of gas during execution, the ZetaClient process re-sends the vote message with an increased gas limit (currently 7,000,000)

This mechanism will be removed in the future once a distinct message for executing inbound is implemented.

### Administrative Operations

The administrative operations can be found at [this page](https://www.zetachain.com/docs/developers/architecture/privileged).

They must be executed using the [group module](https://tutorials.cosmos.network/tutorials/8-understand-sdk-modules/3-group.html).

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

// for withdraws that include a call on the connected chain, gas limit is specified by the user depending on the call complexity
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
These gas price and gas limit values are specific to the connected chain and are completely independent from the gas parameters used for the transaction on the ZetaChain EVM.

**Important Consideration**: The gas limit is fully consumed regardless of whether the cross-chain call uses all the allocated gas.
In the current architecture, any remaining unused gas is sent to a gas stability pool.
This pool is used to compensate for additional fees for transactions during periods of gas price surges,
helping to ensure transaction execution during network congestion.

> **Note**: In a future upgrade, a portion of the unused gas will be refunded to users to optimize fee efficiency.

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

For Solana transactions, ZetaChain uses Solana's compute unit model for fee calculation.
More information about Solana fees can be found in the [Solana Fee Documentation](https://solana.com/docs/core/fees).

In Solana, fees are represented using compute units, which measure the computational resources required to execute a transaction.

In ZetaChain's implementation, there is no gas price parameter for Solana.
The ZetaChain observer set does not vote on gas prices for Solana chains. Instead, only the gas limit is used, which specifies the number of compute units allocated for the transaction.

The withdraw fee for Solana is determined solely by the gas limit (compute units) specified for the transaction.

### Sui

For Sui transactions, ZetaChain adapts to Sui's unique gas model.
More information about Sui gas can be found in the [Sui Gas Documentation](https://docs.sui.io/concepts/tokenomics/gas-pricing).

ZetaChain maintains a gas price parameter for Sui that represents the computation price on the Sui network.

**Important**: Unlike other chains, the gas limit parameter for Sui does not represent a gas limit in computation units.
Instead, the gas limit represents the maximum number of units used to calculate the overall gas budget.

The withdraw fee for Sui is calculated as the gas budget:
```
Withdraw Fee (Gas Budget in MIST) = Gas Price × Gas Limit
```

### TON

For TON transactions, ZetaChain uses a simplified fee model adapted to TON's architecture.
More information about TON fees can be found in the [TON Fee Documentation](https://docs.ton.org/v3/documentation/smart-contracts/transaction-fees/fees).

The gas price for TON is fixed at **400** as part of the base chain configuration.

TON's native fee structure includes multiple components that influence the total transaction cost, such as storage fees,
compute fees, and action fees.
However, for a simpler model—and because ZetaChain currently only supports asset transfers on TON where fees are relatively inexpensive—a capped gas unit value is used.

The gas limit is set to **21,000** to ensure sufficient coverage for all transaction fees. In practice, the withdraw fee is hardcoded.

The withdraw fee for TON is calculated as:
```
Withdraw Fee = Gas Price × Gas Limit = 400 × 21,000 = 8,400,000
```

### Reverting Withdraw

When a withdraw transaction reverts on the connected chain, there are two scenarios:

**Without `callOnRevert`**: If no `callOnRevert` option was specified during the withdraw, the entire withdraw amount is returned to the user on ZetaChain.

**With `callOnRevert`**: If a `callOnRevert` option was specified during the withdraw, an `onRevert` hook is called on ZetaChain.

In both scenarios, the fees for the revert transaction are covered by the protocol. This means the user receives the full withdraw amount back.

For `onRevert` hook calls, a fixed gas limit is used by the protocol.
This gas limit can be retrieved with the [following query](https://zetachain.blockpi.network/lcd/v1/public/zeta-chain/fungible/system_contract) under the `gateway_gas_limit` field.

### Aborting Withdraw

Aborts occur when a revert transaction itself reverts or cannot be executed.

When the user specifies an abort address, the protocol will attempt to call the `onAbort` hook on this address on ZetaChain.

As with revert transactions, the fees for the abort transaction are covered by the protocol and use the fixed gas limit defined in the system contract.

## Connected Chains Deposit Fees

When making deposits from a connected chain to ZetaChain,
users typically pay the inherent fees of the connected chain for the deposit transaction.
The deposit call on ZetaChain is typically covered by the protocol.
There are some exceptions per chain covered below.

For deposits on ZetaChain that include a smart contract call, a fixed gas limit is used by the protocol for the call.
This gas limit can be retrieved with the [following query](https://zetachain.blockpi.network/lcd/v1/public/zeta-chain/fungible/system_contract) under the `gateway_gas_limit` field.

### EVM

[Ethereum Gas Documentation](https://ethereum.org/developers/docs/gas/).

When executing a transaction that includes a single deposit to ZetaChain, the user only pays the inherent fees of the connected EVM chain for the deposit transaction.

As a protection against spam deposits, if the transaction contains multiple deposits to ZetaChain, an additional fee per deposit is charged.

The additional fee is the `additionalActionFeeWei` parameter defined in the connected chain gateway contract.

Therefore, the total additional fee for multiple deposits is calculated in wei (or other base unit) as:
```
Total Additional Fee = (Number of Deposits - 1) × additionalActionFeeWei
```

> **Note**: As of this writing, the multiple deposits mechanism is not yet implemented on ZetaChain mainnet. Subsequent deposits in a single transaction are currently ignored.

### Bitcoin

For Bitcoin transactions, refer to the [Bitcoin Fee Documentation](https://developer.bitcoin.org/devguide/transactions.html).

In addition to the Bitcoin network fees, a depositor fee is deducted from the deposited amount.
The depositor fee is charged to cover the cost of spending the deposited UTXO in the future.

The logic to calculate this depositor fee can be found in the [following page](https://www.zetachain.com/docs/developers/chains/bitcoin#fees).

### Solana

[Solana Fee Documentation](https://solana.com/docs/core/fees).

An additional deposit fee of 2,000,000 Lamports (`0.002` SOL) is charged for deposits from the Solana network to ZetaChain, on top of the standard Solana network transaction fees.

### Sui

[Sui Gas Documentation](https://docs.sui.io/concepts/tokenomics/gas-pricing).

No additional fees are charged for deposits from the Sui network to ZetaChain.

### TON

[TON Fee Documentation](https://docs.ton.org/v3/documentation/smart-contracts/transaction-fees/fees).

TON differs from other chains in that the provided amount is not separate from the fee paid for the transaction. The gas fees are deducted directly from the deposited amount.

For `deposit` and `depositAndCall` operations, ZetaChain currently uses the [ordinary mode for messages](https://docs.ton.org/v3/documentation/smart-contracts/message-management/sending-messages#message-modes).

The fee calculation formula is defined in the [gateway contract](https://github.com/zeta-chain/protocol-contracts-ton/blob/d30343520a7e4167658fcacd786b6adb38c20959/contracts/common/gas.fc#L57C12-L57C78):
```
Fee = flat_gas_price + (gas_amount - flat_gas_limit) × (gas_price >> 16)
```

Where:
- `gas_amount` depends on the operation type ([defined here](https://github.com/zeta-chain/protocol-contracts-ton/blob/d30343520a7e4167658fcacd786b6adb38c20959/contracts/gateway.fc#L37)):
    - **10,000** for `deposit`
    - **13,000** for `depositAndCall`
- `flat_gas_price`, `flat_gas_limit`, and `gas_price` are determined by the [TON chain configuration](https://tonviewer.com/config#21) and remain consistent across testnet and mainnet.

### Reverting Deposit

When a deposit reverts, ZetaChain will initiate a revert transaction on the connected chain.

The fees for the revert transaction follow the same rules as the withdraw fees section, with the formula:
```
Revert Fee = gasLimit × gasPrice
```

The `gasLimit` is determined by the ZRC20 default gas limit value, or by a custom `onRevertGasLimit` specified during the deposit.

The fees are directly deducted from the amount transferred during the deposit:
- When the asset is the native gas token of the connected chain, the fees are deducted from the deposited amount.
- When the asset is a token, the tokens are swapped to the native gas token to cover the fees using internal liquidity pools.

### Aborting Deposit

Aborts occur when a revert transaction itself reverts or cannot be executed.

Aborting deposits have the same behavior as aborting withdraws. When the user specifies an abort address, the protocol will attempt to call the `onAbort` hook on this address on ZetaChain.

As with revert transactions, the fees for the abort transaction are covered by the protocol and use the fixed gas limit defined in the system contract.

## Resources

You can find below resources and references for the fee mechanisms of the different technologies that ZetaChain infrastructure relies on.

* [Cosmos SDK Fees Documentation](https://docs.cosmos.network/v0.53/learn/beginner/gas-fees)
* [EVM Fees Documentation](https://ethereum.org/developers/docs/gas/)
* [Bitcoin Transaction documentation](https://developer.bitcoin.org/devguide/transactions.html)
* [Solana Fees Documentation](https://solana.com/docs/core/fees)
* [Sui Fees Documentation](https://docs.sui.io/concepts/tokenomics/gas-pricing)
* [TON Fees Documentation](https://docs.ton.org/v3/documentation/smart-contracts/transaction-fees/fees)

---

*Last updated: 2025-10-29*