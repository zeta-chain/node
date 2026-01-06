# Slashing Precompile

The Slashing precompile provides an EVM interface to the Cosmos SDK slashing module, enabling smart contracts
to interact with validator slashing information and allowing jailed validators to unjail themselves.

## Address

The precompile is available at the fixed address: `0x0000000000000000000000000000000000000806`

## Interface

### Data Structures

```solidity
// Validator signing information for liveness monitoring
struct SigningInfo {
    address validatorAddress;       // Validator operator address
    int64 startHeight;             // Height at which validator was first a candidate or was unjailed
    int64 indexOffset;             // Index offset into signed block bit array
    int64 jailedUntil;             // Timestamp until which validator is jailed
    bool tombstoned;               // Whether validator has been permanently removed
    int64 missedBlocksCounter;     // Count of missed blocks
}

// Slashing module parameters
struct Params {
    int64 signedBlocksWindow;      // Number of blocks to track for signing
    Dec minSignedPerWindow;        // Minimum percentage of blocks to sign
    int64 downtimeJailDuration;    // Duration of jail time for downtime
    Dec slashFractionDoubleSign;   // Slash percentage for double signing
    Dec slashFractionDowntime;     // Slash percentage for downtime
}

// Decimal type representation
struct Dec {
    string value;  // Decimal string representation
}
```

### Transaction Methods

```solidity
// Unjail a validator after downtime slashing
function unjail(address validatorAddress) external returns (bool success);
```

### Query Methods

```solidity
// Get signing info for a specific validator
function getSigningInfo(
    address consAddress
) external view returns (SigningInfo memory signingInfo);

// Get signing info for all validators with pagination
function getSigningInfos(
    PageRequest calldata pagination
) external view returns (
    SigningInfo[] memory signingInfos,
    PageResponse memory pageResponse
);

// Get slashing module parameters
function getParams() external view returns (Params memory params);
```

## Gas Costs

Gas costs are calculated dynamically based on:

- Base gas for the method
- Storage operations for state changes
- Query complexity for read operations

The precompile uses standard gas configuration for storage operations.

## Implementation Details

### Unjail Mechanism

1. **Eligibility Check**: Validator must be jailed and jail period must have expired
2. **Sender Verification**: Only the validator themselves can request unjailing
3. **State Update**: Updates validator status from jailed to active
4. **Event Emission**: Emits ValidatorUnjailed event

### Signing Information

- **Consensus Address**: Uses the validator's consensus address (from Tendermint ed25519 public key)
- **Liveness Tracking**: Monitors block signing to detect downtime
- **Jail Status**: Tracks jail duration and tombstone status

### Parameter Management

The slashing parameters control:

- **Downtime Detection**: Window size and minimum signing percentage
- **Penalties**: Slash fractions for different infractions
- **Jail Duration**: Time validators must wait before unjailing

## Events

```solidity
event ValidatorUnjailed(address indexed validator);
```

## Security Considerations

1. **Authorization**: Only validators can unjail themselves - no third-party unjailing
2. **Jail Period Enforcement**: Cannot unjail before jail duration expires
3. **Tombstone Protection**: Tombstoned validators cannot be unjailed
4. **Balance Handler**: Proper integration with native token management

## Usage Example

```solidity
ISlashing slashing = ISlashing(SLASHING_PRECOMPILE_ADDRESS);

// Query validator signing info
address consAddress = 0x...; // Validator consensus address
SigningInfo memory info = slashing.getSigningInfo(consAddress);

// Check if validator is jailed
if (info.jailedUntil > int64(block.timestamp)) {
    // Validator is currently jailed

    // If jail period has expired and caller is the validator
    if (block.timestamp >= uint64(info.jailedUntil)) {
        // Unjail the validator
        bool success = slashing.unjail(msg.sender);
        require(success, "Failed to unjail");
    }
}

// Query slashing parameters
Params memory params = slashing.getParams();
// Access parameters like params.signedBlocksWindow
```

## Integration with Validator Operations

```solidity
contract ValidatorManager {
    ISlashing constant slashing = ISlashing(SLASHING_PRECOMPILE_ADDRESS);

    function checkValidatorStatus(address validatorAddr) public view returns (
        bool isJailed,
        int64 jailedUntil,
        bool canUnjail
    ) {
        // Convert to consensus address (implementation specific)
        address consAddr = getConsensusAddress(validatorAddr);

        SigningInfo memory info = slashing.getSigningInfo(consAddr);

        isJailed = info.jailedUntil > int64(block.timestamp);
        jailedUntil = info.jailedUntil;
        canUnjail = isJailed && block.timestamp >= uint64(info.jailedUntil) && !info.tombstoned;
    }

    function autoUnjail(address validatorAddr) external {
        (, , bool canUnjail) = checkValidatorStatus(validatorAddr);
        require(canUnjail, "Cannot unjail yet");
        require(msg.sender == validatorAddr, "Only validator can unjail");

        slashing.unjail(validatorAddr);
    }
}
```

## Address Conversion

The precompile uses different address types:

- **Validator Address**: Standard Ethereum address (operator address)
- **Consensus Address**: Derived from validator's CometBFT public key

Consensus addresses are typically found in:

- `$HOME/.evmd/config/priv_validator_key.json`
- Validator info queries

## Integration Notes

- The precompile integrates directly with the Cosmos SDK slashing module
- All slashing rules and parameters from the chain apply
- Validators must monitor their signing performance to avoid jailing
- Smart contracts can build automation around validator management

