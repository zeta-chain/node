pragma solidity ^0.8.26;

/// @dev The IStaking contract's address.
address constant ISTAKING_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000066; // 102

/// @dev The IStaking contract's instance.
IStaking constant ISTAKING_CONTRACT = IStaking(
    ISTAKING_PRECOMPILE_ADDRESS
);

interface IStaking {
    /// @dev Stake coins to validator
    /// @param staker Staker address
    /// @param validator Validator address
    /// @param amount Coins amount
    /// @return success Staking success
    function stake(
        address staker,
        string memory validator,
        uint256 amount
    ) external returns (bool success);

    /// @dev Unstake coins from validator
    /// @param staker Staker address
    /// @param validator Validator address
    /// @param amount Coins amount
    /// @return completionTime Time when unstaking is done
    function unstake(
        address staker,
        string memory validator,
        uint256 amount
    ) external returns (int64 completionTime);

    /// @dev Transfer coins from validatorSrc to validatorDst
    /// @param staker Staker address
    /// @param validatorSrc Validator from address
    /// @param validatorDst Validator to address
    /// @param amount Coins amount
    /// @return completionTime Time when staket transfer is done
    function transferStake(
        address staker,
        string memory validatorSrc,
        string memory validatorDst,
        uint256 amount
    ) external returns (int64 completionTime);
}
