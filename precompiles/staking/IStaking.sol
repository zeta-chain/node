// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

/// @dev The IStaking contract's address.
address constant ISTAKING_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000066; // 102

/// @dev The IStaking contract's instance.
IStaking constant ISTAKING_CONTRACT = IStaking(ISTAKING_PRECOMPILE_ADDRESS);

/// @notice Bond status for validator
enum BondStatus {
    Unspecified,
    Unbonded,
    Unbonding,
    Bonded
}

/// @notice Validator info
struct Validator {
    string operatorAddress;
    string consensusPubKey;
    bool jailed;
    BondStatus bondStatus;
}

/// @notice Cosmos coin representation.
/// ref: https://github.com/cosmos/cosmos-sdk/blob/470e0859462b28a53adb411843539561d11d7bf5/x/distribution/README.md?plain=1#L139
struct DecCoin {
    string denom;
    uint256 amount;
}

interface IStaking {
    /// @notice Stake event is emitted when stake function is called
    /// @param staker Staker address
    /// @param validator Validator address
    /// @param amount Coins amount
    event Stake(
        address indexed staker,
        address indexed validator,
        uint256 amount
    );

    /// @notice Unstake event is emitted when unstake function is called
    /// @param staker Staker address
    /// @param validator Validator address
    /// @param amount Coins amount
    event Unstake(
        address indexed staker,
        address indexed validator,
        uint256 amount
    );

    /// @notice MoveStake event is emitted when moveStake function is called
    /// @param staker Staker address
    /// @param validatorSrc Validator from address
    /// @param validatorDst Validator to address
    /// @param amount Coins amount
    event MoveStake(
        address indexed staker,
        address indexed validatorSrc,
        address indexed validatorDst,
        uint256 amount
    );

    /// @notice Distributed event is emitted when distribute function is called successfully.
    /// @param zrc20_distributor Distributor address.
    /// @param zrc20_token ZRC20 token address.
    /// @param amount Distributed amount.
    event Distributed(
        address indexed zrc20_distributor,
        address indexed zrc20_token,
        uint256 amount
    );

    /// @notice ClaimedRewards is emitted when a delegator claims ZRC20.
    /// @param claim_address Delegator address where the funds were withdrawed.
    /// @param zrc20_token ZRC20 token address.
    /// @param amount Claimed amount.
    event ClaimedRewards(
        address indexed claim_address,
        address indexed zrc20_token,
        uint256 amount
    );

    /// @notice Stake coins to validator
    /// @param staker Staker address
    /// @param validator Validator address
    /// @param amount Coins amount
    /// @return success Staking success
    function stake(
        address staker,
        string memory validator,
        uint256 amount
    ) external returns (bool success);

    /// @notice Unstake coins from validator
    /// @param staker Staker address
    /// @param validator Validator address
    /// @param amount Coins amount
    /// @return completionTime Time when unstaking is done
    function unstake(
        address staker,
        string memory validator,
        uint256 amount
    ) external returns (int64 completionTime);

    /// @notice Move coins from validatorSrc to validatorDst
    /// @param staker Staker address
    /// @param validatorSrc Validator from address
    /// @param validatorDst Validator to address
    /// @param amount Coins amount
    /// @return completionTime Time when stake move is done
    function moveStake(
        address staker,
        string memory validatorSrc,
        string memory validatorDst,
        uint256 amount
    ) external returns (int64 completionTime);

    /// @notice Get all validators
    /// @return validators All validators
    function getAllValidators()
        external
        view
        returns (Validator[] calldata validators);

    /// @notice Get shares for staker in validator
    /// @return shares Staker shares in validator
    function getShares(
        address staker,
        string memory validator
    ) external view returns (uint256 shares);

    /// @notice Distribute a ZRC20 token as staking rewards.
    /// @param zrc20 The ZRC20 token address to be distributed.
    /// @param amount The amount of ZRC20 tokens to distribute.
    /// @return success Boolean indicating whether the distribution was successful.
    function distribute(
        address zrc20,
        uint256 amount
    ) external returns (bool success);

    /// @notice Claim ZRC20 staking rewards.
    /// @param validator The validator address to claim rewards from.
    /// @return success Boolean indicating whether the claim was successful.
    function claimRewards(
        address delegator,
        string memory validator
    ) external returns (bool success);

    /// @dev Queries all validators the delegator has delegated to.
    /// @param delegator The delegator address to query rewards from.
    /// @return validators List of the validators the caller has delegated to.
    function getDelegatorValidators(
        address delegator
    ) external view returns (string[] calldata validators);

    /// @notice Query ZRC20 outstanding staking rewards.
    /// @param delegator The delegator address to query rewards from.
    /// @param validator The validator address to query rewards from.
    /// @return rewards The list of coins rewarded on the validator.
    function getRewards(
        address delegator,
        string memory validator
    ) external view returns (DecCoin[] calldata rewards);
}
