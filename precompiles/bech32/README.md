# Bech32 Precompile

## Address

`0x0000000000000000000000000000000000000400`

## Description

The Bech32 precompile provides address format conversion between Ethereum hex addresses and Cosmos bech32 addresses.
This enables smart contracts to interact with Cosmos SDK modules that require bech32-formatted addresses.

## Interface

### Methods

#### hexToBech32

```solidity
function hexToBech32(
    address addr,
    string memory prefix
) external returns (string memory bech32Address);
```

Converts an Ethereum hex address to bech32 format with the specified human-readable prefix (HRP).

**Parameters:**

- `addr`: The Ethereum address to convert
- `prefix`: The bech32 human-readable prefix (e.g., "cosmos", "evmos")

**Returns:**

- Bech32-formatted address string

**Validation:**

- Prefix must be non-empty and properly formatted
- Address must be a valid 20-byte Ethereum address
- Reverts if bech32 encoding fails

#### bech32ToHex

```solidity
function bech32ToHex(
    string memory bech32Address
) external returns (address addr);
```

Converts a bech32-formatted address to Ethereum hex format.

**Parameters:**

- `bech32Address`: The bech32 address string to convert

**Returns:**

- Ethereum address in hex format

**Validation:**

- Input must be a valid bech32 address with proper formatting
- Address must contain the separator character "1"
- The decoded address must be 20 bytes
- Reverts if bech32 decoding fails

## Implementation Details

### Gas Usage

The precompile uses a configurable base gas amount for all operations.
The gas cost is fixed regardless of string length *within reasonable bounds*.

### Address Validation

Both methods perform validation on the address format:

- Hex addresses must be exactly 20 bytes
- Bech32 addresses must conform to the bech32 specification
- Invalid addresses result in execution reversion

### Prefix Handling

For `hexToBech32`:

- The prefix parameter determines the human-readable part of the bech32 address
- Common prefixes include account addresses, validator addresses, and consensus addresses
- Empty or whitespace-only prefixes are rejected

For `bech32ToHex`:

- The prefix is automatically extracted from the bech32 address
- No prefix parameter is required as it's embedded in the address

### State Mutability

Both methods are marked as `nonpayable` in the ABI but function as read-only operations.
They do not modify blockchain state and could technically be seen as `view` functions.

