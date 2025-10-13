# Staking Precompile

The Staking precompile provides an EVM interface to the Cosmos SDK staking module.
This enables smart contracts to perform staking operations including validator management, delegation, undelegation, redelegation.

## Address

The precompile is available at the fixed address: `0x0000000000000000000000000000000000000800`

## Interface

### Data Structures

```solidity
// Validator description
struct Description {
    string moniker;
    string identity;
    string website;
    string securityContact;
    string details;
}

// Commission rates for validators
struct CommissionRates {
    uint256 rate;           // Current commission rate (as integer, e.g., 100 = 0.01 = 1%)
    uint256 maxRate;        // Maximum commission rate
    uint256 maxChangeRate;  // Maximum daily increase
}

// Validator information
struct Validator {
    string operatorAddress;     // Validator operator address (bech32)
    string consensusPubkey;     // Consensus public key
    bool jailed;                // Whether validator is jailed
    BondStatus status;          // Bonding status
    uint256 tokens;             // Total tokens
    uint256 delegatorShares;    // Total delegator shares
    Description description;    // Description struct
    int64 unbondingHeight;      // Height when unbonding started
    int64 unbondingTime;        // Time when unbonding completes
    uint256 commission;         // Current commission rate
    uint256 minSelfDelegation;  // Minimum self delegation
}

// Validator bonding status
enum BondStatus {
    Unspecified,
    Unbonded,
    Unbonding,
    Bonded
}

// Unbonding delegation entry
struct UnbondingDelegationEntry {
    int64 creationHeight;
    int64 completionTime;
    uint256 initialBalance;
    uint256 balance;
    uint64 unbondingId;
    int64 unbondingOnHoldRefCount;
}
```

### Transaction Methods

```solidity
// Create a new validator
function createValidator(
    Description calldata description,
    CommissionRates calldata commissionRates,
    uint256 minSelfDelegation,
    address validatorAddress,
    string memory pubkey,
    uint256 value
) external returns (bool success);

// Edit validator parameters
function editValidator(
    Description calldata description,
    address validatorAddress,
    int256 commissionRate,      // Use -1 to keep current value
    int256 minSelfDelegation    // Use -1 to keep current value
) external returns (bool success);

// Delegate tokens to a validator
function delegate(
    address delegatorAddress,
    string memory validatorAddress,
    uint256 amount
) external returns (bool success);

// Undelegate tokens from a validator
function undelegate(
    address delegatorAddress,
    string memory validatorAddress,
    uint256 amount
) external returns (int64 completionTime);

// Redelegate tokens between validators
function redelegate(
    address delegatorAddress,
    string memory validatorSrcAddress,
    string memory validatorDstAddress,
    uint256 amount
) external returns (int64 completionTime);

// Cancel an unbonding delegation
function cancelUnbondingDelegation(
    address delegatorAddress,
    string memory validatorAddress,
    uint256 amount,
    uint256 creationHeight
) external returns (bool success);
```

### Query Methods

```solidity
// Query delegation info
function delegation(
    address delegatorAddress,
    string memory validatorAddress
) external view returns (uint256 shares, Coin calldata balance);

// Query unbonding delegation
function unbondingDelegation(
    address delegatorAddress,
    string memory validatorAddress
) external view returns (UnbondingDelegationOutput calldata unbondingDelegation);

// Query validator info
function validator(
    address validatorAddress
) external view returns (Validator calldata validator);

// Query validators by status
function validators(
    string memory status,
    PageRequest calldata pageRequest
) external view returns (
    Validator[] calldata validators,
    PageResponse calldata pageResponse
);

// Query redelegation info
function redelegation(
    address delegatorAddress,
    string memory srcValidatorAddress,
    string memory dstValidatorAddress
) external view returns (RedelegationOutput calldata redelegation);
```

## Gas Costs

Gas costs are calculated dynamically based on:

- Base gas for the method
- Complexity of the staking operation
- Storage operations for state changes

The precompile uses standard gas configuration for storage operations.

## Implementation Details

### Validator Creation

1. **Self-delegation**: Initial stake must meet minimum self-delegation requirement
2. **Commission rates**: Must be within valid ranges (0-100%)
3. **Public key**: Must be a valid ed25519 consensus public key
4. **Description**: All fields are optional except moniker

### Delegation Operations

- **Delegate**: Stakes tokens with a validator, receiving shares in return
- **Undelegate**: Initiates unbonding process (subject to unbonding period)
- **Redelegate**: Moves stake between validators without unbonding period
- **Cancel Unbonding**: Reverses an unbonding delegation before completion

### Address Formats

- **Validator addresses**: Can be either Ethereum hex or Cosmos bech32 format
- **Delegator addresses**: Ethereum hex addresses
- **Consensus pubkey**: Base64 encoded ed25519 public key

### Commission Updates

- Validators can update commission rates within constraints
- Cannot exceed `maxRate` or increase by more than `maxChangeRate` per day
- Use special constant `-1` to keep current values unchanged

## Events

```solidity
event CreateValidator(
    address indexed validatorAddress,
    uint256 value
);

event EditValidator(
    address indexed validatorAddress,
    int256 commissionRate,
    int256 minSelfDelegation
);

event Delegate(
    address indexed delegatorAddress,
    string indexed validatorAddress,
    uint256 amount,
    uint256 shares
);

event Unbond(
    address indexed delegatorAddress,
    string indexed validatorAddress,
    uint256 amount,
    int64 completionTime
);

event Redelegate(
    address indexed delegatorAddress,
    string indexed validatorSrcAddress,
    string indexed validatorDstAddress,
    uint256 amount,
    int64 completionTime
);

event CancelUnbondingDelegation(
    address indexed delegatorAddress,
    string indexed validatorAddress,
    uint256 amount,
    uint256 creationHeight
);
```

## Security Considerations

1. **Sender Verification**: All operations verify the transaction sender matches the specified address
2. **Balance Handling**: Uses the balance handler for proper native token management
3. **Unbonding Period**: Enforces chain-wide unbonding period for security
4. **Slashing Risk**: Delegated tokens are subject to slashing for validator misbehavior

## Usage Examples

### Creating a Validator

```solidity
StakingI staking = StakingI(STAKING_PRECOMPILE_ADDRESS);

Description memory desc = Description({
    moniker: "My Validator",
    identity: "keybase-identity",
    website: "https://validator.example.com",
    securityContact: "security@example.com",
    details: "Professional validator service"
});

CommissionRates memory rates = CommissionRates({
    rate: 100,           // 1% (100 / 10000)
    maxRate: 2000,       // 20% max
    maxChangeRate: 100   // 1% max daily change
});

// Create validator with 1000 tokens self-delegation
bool success = staking.createValidator(
    desc,
    rates,
    1000e18,             // Min self delegation
    msg.sender,          // Validator address
    "validator_pubkey",  // Consensus public key
    1000e18              // Initial self delegation
);
```

### Delegating to a Validator

```solidity
StakingI staking = StakingI(STAKING_PRECOMPILE_ADDRESS);

// Delegate 100 tokens to a validator
string memory validatorAddr = "evmosvaloper1..."; // Bech32 validator address
bool success = staking.delegate(msg.sender, validatorAddr, 100e18);

// Query delegation
(uint256 shares, Coin memory balance) = staking.delegation(msg.sender, validatorAddr);
```

### Managing Delegations

```solidity
// Undelegate 50 tokens (starts unbonding period)
int64 completionTime = staking.undelegate(msg.sender, validatorAddr, 50e18);

// Redelegate to another validator (no unbonding period)
string memory newValidator = "evmosvaloper2...";
int64 redelegationTime = staking.redelegate(
    msg.sender,
    validatorAddr,
    newValidator,
    25e18
);

// Cancel unbonding (must specify the creation height)
uint256 creationHeight = 12345;
staking.cancelUnbondingDelegation(
    msg.sender,
    validatorAddr,
    50e18,
    creationHeight
);
```

## Integration Notes

- The precompile integrates directly with the Cosmos SDK staking module
- All staking parameters and rules from the chain apply
- Amounts use the bond denomination precision (typically 18 decimals)
- Validator addresses can be provided in either hex or bech32 format
- Commission rates are integers where 10000 = 100%

