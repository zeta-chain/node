pragma solidity ^0.8.26;

/// @dev The IStaking contract's address.
address constant ISTAKING_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000066; // 102

/// @dev The IStaking contract's instance.
IStaking constant ISTAKING_CONTRACT = IStaking(
    ISTAKING_PRECOMPILE_ADDRESS
);

interface IStaking {
    /// @dev Delegate coin to validator
    /// @param delegator Delegator address
    /// @param validator Validator address
    /// @param amount Coins amount
    /// @return success Delegation success
    function delegate(
        address delegator,
        string memory validator,
        uint256 amount
    ) external returns (bool success);
}
