# Precompiled Contracts

This directory contains the precompiled contracts for the ZetaChain node. These contracts provide native functionality that can be called from smart contracts.

## Available Contracts

### 1. Prototype Contract (Address: 0x0000000000000000000000000000000000000065)
Enabled by default, this contract provides utility functions for address conversion and gas pool queries.

#### Methods:
- `bech32ToHexAddr(string bech32) view returns (address addr)`: Converts a bech32 address to hexadecimal address
- `bech32ify(string prefix, address addr) view returns (string bech32)`: Converts a hex address to bech32 address
- `getGasStabilityPoolBalance(int64 chainID) view returns (uint256 result)`: Returns the balance of the gas stability pool for a given chain ID

### 2. Bank Contract (Address: 0x0000000000000000000000000000000000000067)
Enabled by default, this contract handles cross-chain token deposits and withdrawals.

#### Methods:
- `deposit(address zrc20, uint256 amount) returns (bool success)`: Deposits ZRC20 tokens and mints corresponding Cosmos tokens
- `withdraw(address zrc20, uint256 amount) returns (bool success)`: Withdraws Cosmos tokens and converts them back to ZRC20 tokens
- `balanceOf(address zrc20, address user) view returns (uint256 balance)`: Retrieves the Cosmos token balance for a specific ZRC20 token and user

### 3. Staking Contract (Address: 0x0000000000000000000000000000000000000066)
Currently disabled by default, this contract provides staking functionality.

#### Methods:
- `stake(address staker, string validator, uint256 amount) returns (bool success)`: Stakes tokens with a validator
- `unstake(address staker, string validator, uint256 amount) returns (int64 completionTime)`: Unstakes tokens from a validator
- `moveStake(address staker, string validatorSrc, string validatorDst, uint256 amount) returns (int64 completionTime)`: Moves stake from one validator to another
- `getAllValidators() view returns (Validator[] validators)`: Returns all validators
- `getShares(address staker, string validator) view returns (uint256 shares)`: Returns staker's shares in a validator
- `distribute(address zrc20, uint256 amount) returns (bool success)`: Distributes ZRC20 tokens as staking rewards
- `claimRewards(address delegator, string validator) returns (bool success)`: Claims staking rewards for a delegator


## Usage

These precompiled contracts can be called from smart contracts using their respective interfaces, which are defined in the corresponding `.sol` files. Referer to `node/precompiles`

Example usage in Solidity:
```solidity
// Using the Prototype contract
IPrototype prototype = IPrototype(IPROTOTYPE_PRECOMPILE_ADDRESS);
address addr = prototype.bech32ToHexAddr(bech32Address);
``` 

```go
// Using Go bindings
iPrototype, err := prototype.NewIPrototype(IPROTOTYPE_PRECOMPILE_ADDRESS, ZEVMClient)
addr, err := iPrototype.Bech32ToHexAddr(opts, bech32Address)
```