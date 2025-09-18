# Management for Unused Gas Fee

## Overview

This document explains how gas fees are managed when processing outbound transactions on connected chains, particularly focusing on how differences between estimated and actual gas usage are handled.
The Gas fee is always paid in the Gas ZRC20 token of the connected chain ,irrespective of coin type.

## Gas Fee Scenarios

When processing outbound transactions, there are two possible scenarios regarding gas fees:

1. **User Overpays** (Common)
    - Users often overpay gas fees to ensure transactions are processed
    - EVM chains automatically refund unused gas to the caller (TSS address)
    - This creates an opportunity to return a portion of these funds to users

2. **User Underpays**
    - Gas prices may increase after transaction submission
    - In these cases, the stability pool covers the difference
    - No refund mechanism applies in underpayment scenarios

## Refund Mechanism for Overpayments

For overpayment scenarios, we implement the following refund logic:

### Fee Tracking

- We track the initial gas fee paid by the user when initiating the transaction. This does not include the Protocol Fee.
- When an outbound transaction completes (regardless of status), we calculate the actual fee used:
  ```
  actualFee = receipt.GasUsed * transaction.GasPrice()
  ```

### Difference in Fee Calculation

- The difference in fee is calculated as:
  ```
  totalRemainingFees = userGasFeePaid - actualFee
  ```
- We then take 95% of this amount for further calculations:
  ```
  remainingFees = 95% of totalRemainingFees
  ```
- We intentionally use only 95% of the unused fee to avoid any potential over-minting:
    - The 5% of tokens that are not minted back are retained by the TSS address
    - This creates a safety buffer since we burn the total amount of tokens when initiating the transaction, which creates a deficit

### Fee Distribution

The remaining fees are distributed as follows:

1. **Stability Pool Allocation**
    - For non-EVM chains: 100% of this amount goes to the stability pool
    - For EVM chains: A configurable percentage (defined in the chain parameters) goes to the stability pool

2. **User Refund**
    - For EVM chains: The remaining amount after stability pool allocation is refunded to the user
    - Refunds are provided as ZRC20 gas tokens of the connected chain

3. **Fallback**
    - If refunding fails for any reason, all remaining fees are allocated to the stability pool

## Implementation Details

- The refund is always in the form of gas tokens from the connected chain
- The stability pool funded is always associated with the gas token of the outbound chain
- The system requires chain parameters to be properly configured
- The mechanism handles various edge cases (invalid addresses, zero amounts, etc.).
- We refund the amount to the sender irrespective of whether it is a user or a contract. The contract should implement a mechanism to handle the refund.