pragma solidity ^0.8.26;

/// @dev The IStaking contract's address.
address constant ISTAKING_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000066; // 102

/// @dev The IStaking contract's instance.
IStaking constant ISTAKING_CONTRACT = IStaking(
    ISTAKING_PRECOMPILE_ADDRESS
);

interface IStaking {
    /// @dev Delegate coins to validator
    /// @param delegator Delegator address
    /// @param validator Validator address
    /// @param amount Coins amount
    /// @return success Delegation success
    function delegate(
        address delegator,
        string memory validator,
        uint256 amount
    ) external returns (bool success);

    /// @dev Undelegate coins from validator
    /// @param delegator Delegator address
    /// @param validator Validator address
    /// @param amount Coins amount
    /// @return completionTime Time when undelegation is done
    function undelegate(
        address delegator,
        string memory validator,
        uint256 amount
    ) external returns (int64 completionTime);

    /// @dev Redelegate coins from validatorSrd to validatorDst
    /// @param delegator Delegator address
    /// @param validatorSrc Validator from address
    /// @param validatorDst Validator to address
    /// @param amount Coins amount
    /// @return completionTime Time when redelegation is done
    function redelegate(
        address delegator,
        string memory validatorSrc,
        string memory validatorDst,
        uint256 amount
    ) external returns (int64 completionTime);
}
