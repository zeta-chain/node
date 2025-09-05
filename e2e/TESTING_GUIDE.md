# Regular E2E tests

This page lists the regular E2E tests to run when testing the network, in case of upgrade, etc..
These snippets are aimed to be copy-pasted in the input in the E2E CI tool.

## Inbounds and outbounds observation

When we only want to verify the network correctly observe cross-chain transactions, simple deposits and withdraws are sufficient.

The amount provided represent `0.0001` unit for coin with 18 decimals.

```
eth_deposit:100000000000000 eth_withdraw:100000000000000
```

## ERC20 observation

When we want to verify the network correctly observe cross-chain transactions for ERC20 tokens.

The amount is set to a small value so it can be used for most ERC20s regardless of the decimals.

```
erc20_deposit:1000 erc20_withdraw:1000
```

## Gateway basic workflow

When we want to verify the gateway basic workflow, the happy path where cross-chain calls succeed.

The amount is arbitrarily set to a small value, currently the tokens sent to the test contracts are lost.

```
eth_deposit_and_call:1000 eth_withdraw_and_call:1000 erc20_deposit_and_call:1000 erc20_withdraw_and_call:1000 zevm_to_evm_call evm_to_zevm_call
```

## Solana

When it is necessary to test the Solana workflows, SOL and SPL tokens.

```
solana_deposit:10000000 solana_withdraw:10000000 solana_deposit_and_call:1000 spl_deposit:1000 spl_withdraw:1000 spl_deposit_and_call:1000 solana_deposit_and_call_revert:10000000
```

## Gateway revert workflow

When we want to verify the gateway revert workflow, the unhappy path where cross-chain calls fail

### WithdrawAndCall

The `withdrawAndCall` tests doesn't depend on the provided amount, this list can be used across all networks

```
eth_withdraw_and_call_revert:1000 eth_withdraw_and_call_revert_with_call:1000 erc20_withdraw_and_call_revert:1000 erc20_withdraw_and_call_revert_with_call:1000
```

### DepositAndCall

The amount for reverting `depositAndCall` must depend on the chain as the value in the CCTX is used to pay for the revert fee.

Note: these are estimated required values for mainnet based on the current gas price, the actual value might be different and fine-tuned. The values for ERC20 tests are set for USDC token.

Ethereum: `0.0007ETH` and `3USDC`

```
eth_deposit_and_call_revert:700000000000000 eth_deposit_and_call_revert_with_call:700000000000000 erc20_deposit_and_call_revert:3000000 erc20_deposit_and_call_revert_with_call:3000000
```

BSC: `0.0008BNB` and `0.5USDC`

```
eth_deposit_and_call_revert:800000000000000 eth_deposit_and_call_revert_with_call:800000000000000 erc20_deposit_and_call_revert:500000 erc20_deposit_and_call_revert_with_call:500000
```

Polygon: `0.008POL` and `0.01USDC`

```
eth_deposit_and_call_revert:8000000000000000 eth_deposit_and_call_revert_with_call:8000000000000000 erc20_deposit_and_call_revert:10000 erc20_deposit_and_call_revert_with_call:10000
```

Base: `0.000005ETH` and `0.02USDC`

```
eth_deposit_and_call_revert:5000000000000 eth_deposit_and_call_revert_with_call:5000000000000 erc20_deposit_and_call_revert:20000 erc20_deposit_and_call_revert_with_call:20000
```

## Gateway arbitrary calls

Arbitrary calls feature is an experimental and niche use case for now, these tests are not necessary for regular testing.

```
eth_withdraw_and_arbitrary_call:1000 erc20_withdraw_and_arbitrary_call:1000
```

## Gov proposal

Any governace proposals can be executed in the E2E tests. There are two options for the sequence of the proposal execution.
1. Start of E2E tests: The proposal json needs to be placed in the `contrib/localnet/orchestrator/proposal_e2e_start` directory.All proposals in this directory will be executed before running the e2e tests
2. End of E2E tests: The proposal json needs to be placed in the `contrib/localnet/orchestrator/proposal_e2e_end` directory.All proposals in this directory will be executed after running the e2e tests
