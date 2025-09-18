# ZetaChain Gas Fee Documentation

## Overview

Gas fees in ZetaChain refer to the amount paid by users to execute transactions on connected chains. These fees **do not** include any gas payments that users make on the source chain to initiate transactions.

## Key Concepts

- **Deposit**: Transfer tokens from a connected chain to ZetaChain
- **Withdraw**: Transfer tokens from ZetaChain to a connected chain
- **Coin Types**: Different token types used for transactions (Gas, ERC20, ZETA, No Asset)
- **Revert**: When a transaction fails and needs to be reverted

## V2 Protocol Flows

### Deposit Transactions

Deposits transfer tokens from connected chains to ZetaChain.

- Gas Fee is always in GAS ZRC20. If the amount is in a different ERC20, it's swapped to buy Gas tokens

| Transaction Type | GAS FEE Deposit | GAS FEE Revert | Notes |
|------------------|-----------------|----------------|-------|
| Deposit | Free | User Pays based on cointype | Revert fees paid using the GAS ZRC20 on source chain |
| DepositAndCall | Free | User Pays based on cointype | Revert fees paid using the GAS ZRC20 on source chain |
| Call | Free | Revert Not supported | Revert fees paid using the GAS ZRC20 on source chain |

#### Gas Fee Payment by Coin Type

1. **CoinType ERC20**
   - User is using ERC20/ZRC20 tokens
   - `Protocol` calculates expected gas cost on connected chain, using `Current Gas Price` * `Revert Gas Limit` + Protocol Flat Fee
   - ERC20 tokens are swapped to buy gas ZRC20 tokens
   - Gas ZRC20 tokens are burned to pay for the gas fee

2. **CoinType Gas**
   - User paid in the native gas token on a connected chain
   - `Protocol` calculates expected gas cost on connected chain, using `Current Gas Price` * `Revert Gas Limit` + Protocol Flat Fee
   - Gas token can be directly burned to pay for the gas fee

3. **NoAssetCall**
   - No revert processing in this case, the cross-chain transaction is aborted

### Withdraw Transactions

Withdrawals transfer tokens from ZetaChain to connected chains.

| Transaction Type | GAS FEE Withdraw | GAS FEE Revert | Notes |
|------------------|-----------------|----------------|-------|
| Withdraw | User Pays based on withdraw type | Free call to OnRevert using fixed GasLimit | Initiated through Gateway smart contract |
| WithdrawAndCall | User Pays based on withdraw type | Free call to OnRevert using fixed GasLimit | Initiated through Gateway smart contract |
| Call | User Pays based on withdraw type | Free call to OnRevert using fixed GasLimit | Initiated through Gateway smart contract |

- Withdrawals are initiated by users through zEVM smart contract calls
- Fees are always paid in the GAS ZRC20 token for the target connected chain, this fees is paid separately by the user an is not deducted from the amount.However when the transaction reverts, the revert fee is deducted from the amount
- CCTX type is either GAS or ERC20 depending on the token

**Gas Fee Payment by withdraw type:**

1. **Withdraw**
   - `Gateway Smart contract` calculates expected gas cost on connected chain based on `gasPrice` * `ZRC20.GasLimit` + Protocol Flat Fee

2. **WithdrawAndCall**
   - `Gateway Smart contract` calculates expected gas cost on connected chain based on `gasPrice` * `CallOptions.GasLimit` + Protocol Flat Fee

3. **Call**
   - `Gateway Smart contract` calculates expected gas cost on connected chain based on `gasPrice` * `ZRC20.GasLimit` + Protocol Flat Fee

## V1 Protocol Flows

- V1 flows for ERC20/ZRC20 and GAS tokens have already been deprecated
- V1 flows for ZETA token deposits and withdrawals are still supported but are planned to be deprecated soon.

### Deposit Transactions

| Transaction Type | GAS FEE Deposit | GAS FEE Revert | Notes |
|------------------|-----------------|----------------|-------|
| Deposit zEVM | Free | User pays based on coin type | Gas fee for revert based on cointype |
| Deposit Connected Chain (Msg-Passing) | User pays based on coin type | User pays based on coin type | Fees in connected chain's gas token |

#### Gas Fee Payment by Coin Type

1. **CoinType Zeta**
   - User is using ERC20/ZRC20 "Zeta Token"
   - `Protocol` calculates expected gas cost on connected chain using `gasPrice * 2` * `CallOptions.GasLimit` + Protocol Flat Fee
   - Zeta tokens are swapped to buy gas ZRC20 tokens
   - Gas tokens are burned to pay for the gas fee
   - **Important**: Gas price is multiplied by 2× when calculating the fee

2. **CoinType ERC20**
   - User is using ERC20/ZRC20 tokens
   - `Protocol` calculates expected gas cost on connected chain using `gasPrice` * `GasZRC20.GasLimit` + Protocol Flat Fee
   - ERC20 tokens are swapped to buy gas ZRC20 tokens
   - Gas tokens are burned to pay for the gas fee

3. **CoinType Gas**
   - User paid in native gas token on connected chain
   - `Protocol` calculates expected gas cost on connected chain using `gasPrice` * `GasZRC20.GasLimit` + Protocol Flat Fee
   - Gas token can be directly burned to pay for fee

### Withdraw Transactions

| Transaction Type | GAS FEE Deposit | GAS FEE Revert | Notes |
|------------------|-----------------|----------------|-------|
| ZRC20 Withdraw | User pays | Not supported | Uses same mechanism as V2 withdrawals |
| ZETA Sent | User pays | Free call to OnRevert using fixed GasLimit | Total Zeta (value + gas) is burned |

#### Details by Withdraw Type

1. **ZRC20 Withdraw**
   - Uses the same mechanism as V2 withdrawals but initiated through ZRC20 contract

2. **ZETA Sent**
   - Total Zeta amount (value + gas) is burned
   - A portion is minted to swap and pay for gas in outbound chain
   - `Protocol` calculates expected gas cost on a connected chain using `gasPrice * 2` * `CallOptions.GasLimit` + Protocol Flat Fee
   - Value portion is unlocked in connected chain after successful outbound transaction