// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

/// @dev The IPrototype contract's address.
address constant IPROTOTYPE_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000065; // 101

/// @dev The IPrototype contract's instance.
IPrototype constant IPROTOTYPE_CONTRACT = IPrototype(
    IPROTOTYPE_PRECOMPILE_ADDRESS
);

interface IPrototype {
    /// @dev converting a bech32 address to hexadecimal address.
    /// @param bech32 The bech32 address.
    /// @return addr The hexadecimal address.
    function bech32ToHexAddr(
        string memory bech32
    ) external view returns (address addr);

    /// @dev converting a hex address to bech32 address.
    /// @param prefix of the bech32, e.g. zeta.
    /// @param addr The hex address
    /// @return bech32 The bech32 address.
    function bech32ify(
        string memory prefix,
        address addr
    ) external view returns (string memory bech32);

    /// @dev returns the balance of the gas stability pool
    /// @param chainID to query gas.
    /// @return result of the call.
    function getGasStabilityPoolBalance(
        int64 chainID
    ) external view returns (uint256 result);
}
