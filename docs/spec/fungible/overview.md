# Overview

The `fungible` facilitates the deployment of fungible tokens of connected
blockchains (called "foreign coins") on ZetaChain.

Foreign coins are represented as ZRC20 tokens on ZetaChain.

When a foreign coin is deployed on ZetaChain, a ZRC20 contract is deployed, a
pool is created, liquidity is added to the pool, and the foreign coin is added
to the list of foreign coins in the module's state.

The module contains the logic for:

- Deploying a foreign coin on ZetaChain
- Deploying a system contract, Uniswap and wrapped ZETA
- Depositing to and calling omnichain smart contracts on ZetaChain from
  connected chains (`DepositZRC20AndCallContract` and `DepositZRC20`)

the module depends heavily on the
[protocol contracts](https://github.com/zeta-chain/protocol-contracts).

## State

The `fungible` module keeps track of the following state:

- System contract address
- A list of foreign coins
