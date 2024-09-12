// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

/// @dev The IStaking contract's address.
address constant ISTAKING_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000066; // 102

/// @dev The IStaking contract's instance.
IStaking constant ISTAKING_CONTRACT = IStaking(
    ISTAKING_PRECOMPILE_ADDRESS
);

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
    function getAllValidators() external view returns (Validator[] calldata validators);

    /// @notice Get shares for staker in validator
    /// @return shares Staker shares in validator
    function getShares(address staker, string memory validator) external view returns (uint256 shares);
}
