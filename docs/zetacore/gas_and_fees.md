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

## Connected Chains Withdraw Fees

### EVM

### Solana

### Sui

### TON

### Withdraw reverting

### Withdraw aborting

## Connected Chains Deposit Fees

### EVM

### Solana

### Sui

### TON

### Deposit reverting

### Deposit aborting

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
