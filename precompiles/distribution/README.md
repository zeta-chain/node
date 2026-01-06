# Distribution Precompile

## Address

`0x0000000000000000000000000000000000000801`

## Description

The Distribution precompile provides an EVM interface to the Cosmos SDK `x/distribution` module,
enabling smart contracts to interact with staking rewards, validator commissions,
and the community pool. It supports both reward queries and distribution operations including claiming,
withdrawing, and funding.

## Interface

### Transaction Methods

#### setWithdrawAddress

```solidity
function setWithdrawAddress(
    address delegator,
    string memory withdrawerAddress
) external returns (bool);
```

Sets the address authorized to withdraw rewards for a delegator.

**Parameters:**

- `delegator`: The delegator address setting the withdraw address
- `withdrawerAddress`: The address that will be authorized to withdraw rewards

**Authorization:** Caller must be the delegator

**Gas Cost:** 2000 + (30 × input data size in bytes)

#### withdrawDelegatorRewards

```solidity
function withdrawDelegatorRewards(
    address delegator,
    string memory validator
) external returns (Coin[] memory);
```

Withdraws pending rewards from a specific validator.

**Parameters:**

- `delegator`: The delegator withdrawing rewards
- `validator`: The validator address to withdraw from

**Returns:**

- Array of `Coin` structs representing withdrawn amounts

**Authorization:** Caller must be the delegator

**Gas Cost:** 2000 + (30 × input data size in bytes)

#### withdrawValidatorCommission

```solidity
function withdrawValidatorCommission(
    string memory validator
) external returns (Coin[] memory);
```

Withdraws accumulated commission for a validator.

**Parameters:**

- `validator`: The validator address withdrawing commission

**Returns:**

- Array of `Coin` structs representing withdrawn commission

**Authorization:** Caller must be the validator

**Gas Cost:** 2000 + (30 × input data size in bytes)

#### claimRewards

```solidity
function claimRewards(
    address delegator,
    uint32 maxRetrieve
) external returns (bool);
```

Claims rewards from all validators at once (custom batch operation).

**Parameters:**

- `delegator`: The delegator claiming rewards
- `maxRetrieve`: Maximum number of validators to claim from

**Authorization:** Caller must be the delegator

**Gas Cost:** 2000 + (30 × input data size in bytes)

#### fundCommunityPool

```solidity
function fundCommunityPool(
    address depositor,
    Coin[] memory coins
) external returns (bool);
```

Deposits tokens into the community pool.

**Parameters:**

- `depositor`: The address funding the pool
- `coins`: Array of coins to deposit

**Authorization:** Caller must be the depositor

**Gas Cost:** 2000 + (30 × input data size in bytes)

#### depositValidatorRewardsPool

```solidity
function depositValidatorRewardsPool(
    string memory validator,
    Coin[] memory coins
) external returns (bool);
```

Deposits tokens into a validator's rewards pool.

**Parameters:**

- `validator`: The validator whose pool receives the deposit
- `coins`: Array of coins to deposit

**Gas Cost:** 2000 + (30 × input data size in bytes)

### Query Methods

#### delegationTotalRewards

```solidity
function delegationTotalRewards(
    address delegator
) external view returns (
    DelegatorTotal[] memory,
    DecCoin[] memory
);
```

Returns total rewards across all validators for a delegator.

**Parameters:**

- `delegator`: The delegator address

**Returns:**

- Array of `DelegatorTotal` structs (per-validator rewards)
- Array of `DecCoin` structs (total rewards sum)

**Gas Cost:** 1000 + (3 × input data size in bytes)

#### delegationRewards

```solidity
function delegationRewards(
    address delegator,
    string memory validator
) external view returns (DecCoin[] memory);
```

Returns rewards for a specific delegator-validator pair.

**Parameters:**

- `delegator`: The delegator address
- `validator`: The validator address

**Returns:**

- Array of `DecCoin` structs representing rewards

**Gas Cost:** 1000 + (3 × input data size in bytes)

#### delegatorValidators

```solidity
function delegatorValidators(
    address delegator
) external view returns (string[] memory);
```

Lists all validators from which a delegator can claim rewards.

**Parameters:**

- `delegator`: The delegator address

**Returns:**

- Array of validator addresses

**Gas Cost:** 1000 + (3 × input data size in bytes)

#### delegatorWithdrawAddress

```solidity
function delegatorWithdrawAddress(
    address delegator
) external view returns (string memory);
```

Returns the configured withdraw address for a delegator.

**Parameters:**

- `delegator`: The delegator address

**Returns:**

- The withdraw address

**Gas Cost:** 1000 + (3 × input data size in bytes)

#### communityPool

```solidity
function communityPool() external view returns (DecCoin[] memory);
```

Returns the current balance of the community pool.

**Returns:**

- Array of `DecCoin` structs representing pool balance

**Gas Cost:** 1000 + (3 × input data size in bytes)

#### validatorCommission

```solidity
function validatorCommission(
    string memory validator
) external view returns (DecCoin[] memory);
```

Returns accumulated commission for a validator.

**Parameters:**

- `validator`: The validator address

**Returns:**

- Array of `DecCoin` structs representing commission

**Gas Cost:** 1000 + (3 × input data size in bytes)

#### validatorDistributionInfo

```solidity
function validatorDistributionInfo(
    string memory validator
) external view returns (DistInfo memory);
```

Returns comprehensive distribution information for a validator.

**Parameters:**

- `validator`: The validator address

**Returns:**

- `DistInfo` struct containing commission and self-delegation rewards

**Gas Cost:** 1000 + (3 × input data size in bytes)

#### validatorOutstandingRewards

```solidity
function validatorOutstandingRewards(
    string memory validator
) external view returns (DecCoin[] memory);
```

Returns outstanding (undistributed) rewards for a validator.

**Parameters:**

- `validator`: The validator address

**Returns:**

- Array of `DecCoin` structs representing outstanding rewards

**Gas Cost:** 1000 + (3 × input data size in bytes)

#### validatorSlashes

```solidity
function validatorSlashes(
    string memory validator,
    uint64 startingHeight,
    uint64 endingHeight,
    PageRequest memory pageRequest
) external view returns (
    ValidatorSlashEvent[] memory,
    PageResponse memory
);
```

Returns slashing events for a validator within a height range.

**Parameters:**

- `validator`: The validator address
- `startingHeight`: Start of the query range
- `endingHeight`: End of the query range
- `pageRequest`: Pagination parameters

**Returns:**

- Array of `ValidatorSlashEvent` structs
- `PageResponse` with pagination information

**Gas Cost:** 1000 + (3 × input data size in bytes)

### Data Structures

```solidity
struct Coin {
    string denom;
    uint256 amount;
}

struct DecCoin {
    string denom;
    uint256 amount;
    uint8 precision;
}

struct DelegatorTotal {
    string validatorAddress;
    DecCoin[] rewards;
}

struct DistInfo {
    string operatorAddress;
    DecCoin[] commission;
    DecCoin[] selfBondRewards;
}

struct ValidatorSlashEvent {
    uint64 validatorPeriod;
    Fraction fraction;
}

struct Fraction {
    uint256 numerator;
    uint256 denominator;
}
```

## Message Type Constants

```solidity
string constant MSG_SET_WITHDRAWER_ADDRESS = "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress"
string constant MSG_WITHDRAW_DELEGATOR_REWARD = "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward"
string constant MSG_WITHDRAW_VALIDATOR_COMMISSION = "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission"
```

## Implementation Details

### Authorization

All transaction methods enforce that the caller matches the relevant account (delegator or validator)
to prevent unauthorized operations.

### Balance Tracking

The precompile tracks native token balance changes during transaction execution to accurately return transfer amounts.

### Event Emission

Each transaction emits corresponding events for on-chain tracking and indexing.

### Address Format Support

The precompile accepts both hex and bech32 address formats, automatically converting as needed for Cosmos SDK compatibility.

