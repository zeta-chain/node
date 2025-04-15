# Transaction Query Validation

This workflow queries and validates transactions from recent blocks in ZetaChain (Athens or mainnet).

The idea is to query a block, then query some of its transactions and check that RPC is not failing to return the result.

## Configuration

The `config/default.json` file contains the configuration for the transaction query limits:
- `max_blocks`: The maximum number of recent blocks to query
- `max_transactions`: The maximum number of transactions to process

## Scaling Guidance

Start with small limits (10 blocks, 10 transactions) to verify functionality.
After confirming everything works properly, gradually increase the limits:
- Increase to 100 blocks and 1,000 transactions
- Then to 500 blocks and 50,000 transactions
- Finally to 1,000 blocks and 100,000 transactions (or as needed)
