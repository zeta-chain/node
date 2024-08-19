// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.26;

/// @dev The IRegular contract's address.
address constant IPROTOTYPE_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000065; // 101

/// @dev The IRegular contract's instance.
IPrototype constant IREGULAR_CONTRACT = IPrototype(
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

    /// @dev Function to verify calling regular contact through precompiled contact
    /// @param method to call, e.g. bar.
    /// @param addr of deployed regular contract.
    /// @return result of the call.
    function regularCall(
        string memory method,
        address addr
    ) external returns (uint256 result);
}
