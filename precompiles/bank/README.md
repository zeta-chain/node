# Bank Precompile

## Address

`0x0000000000000000000000000000000000000804`

## Description

The Bank precompile provides read-only access to the Cosmos SDK `x/bank` module state through an EVM-compatible interface.
This enables smart contracts to query native token balances and supply information
for accounts and tokens registered with corresponding ERC-20 representations.

## Interface

### Methods

#### balances

```solidity
function balances(address account) external view returns (Balance[] memory)
```

Retrieves all native token balances for the specified account.
Each balance includes the ERC-20 contract address and amount in the token's original precision.

**Parameters:**

- `account`: The account address to query

**Returns:**

- Array of `Balance` structs containing:
    - `contractAddress`: ERC-20 contract address representing the native token
    - `amount`: Token balance in smallest denomination

**Gas Cost:** 2,851 + (2,851 × (n-1)) where n = number of tokens returned

#### totalSupply

```solidity
function totalSupply() external view returns (Balance[] memory)
```

Retrieves the total supply of all native tokens in the system.

**Parameters:** None

**Returns:**

- Array of `Balance` structs containing:
    - `contractAddress`: ERC-20 contract address representing the native token
    - `amount`: Total supply in smallest denomination

**Gas Cost:** 2,477 + (2,477 × (n-1)) where n = number of tokens returned

#### supplyOf

```solidity
function supplyOf(address erc20Address) external view returns (uint256)
```

Retrieves the total supply of a specific token by its ERC-20 contract address.

**Parameters:**

- `erc20Address`: The ERC-20 contract address of the token

**Returns:**

- Total supply as `uint256`. Returns 0 if the token is not registered.

**Gas Cost:** 2,477

### Data Structures

```solidity
struct Balance {
    address contractAddress;  // ERC-20 contract address
    uint256 amount;          // Amount in smallest denomination
}
```

## Implementation Details

### Token Resolution

The precompile resolves native Cosmos SDK denominations to their corresponding ERC-20
contract addresses through the `x/erc20` module's token pair registry.
Only tokens with registered token pairs are returned in query results.

### Decimal Precision

All amounts returned preserve the original decimal precision stored in the `x/bank` module.
No decimal conversion is performed by the precompile.

### Gas Metering

The precompile implements efficient gas metering by:

- Charging base gas for the first result
- Incrementally charging for each additional result in batch queries
- Consuming gas before returning results to prevent DoS vectors

### Error Handling

- Invalid token addresses in `supplyOf` return 0 rather than reverting
- Queries for accounts with no balances return empty arrays
- All methods are read-only and cannot modify state

