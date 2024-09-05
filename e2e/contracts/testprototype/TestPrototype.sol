// SPDX-License-Identifier: MIT
pragma solidity 0.8.10;

interface IPrototype {
    function bech32ToHexAddr(
        string memory bech32
    ) external view returns (address addr);

    function bech32ify(
        string memory prefix,
        address addr
    ) external view returns (string memory bech32);

    function getGasStabilityPoolBalance(
        int64 chainID
    ) external view returns (uint256 result);
}

// @dev Purpose of this contract is to test calling prototype precompile through contract
// every function calling precompiles must have return so solidity doesn't check for precompile code size and revert because it's 0
// version of solidity used must be >= 0.8.10 to support this
contract TestPrototype {
    IPrototype prototype = IPrototype(0x0000000000000000000000000000000000000065);

    function bech32ToHexAddr(string memory bech32) external view returns (address addr) {
        return prototype.bech32ToHexAddr(bech32);
    }

    function bech32ify(string memory prefix, address addr) external view returns (string memory bech32) {
        return prototype.bech32ify(prefix, addr);
    }

     function getGasStabilityPoolBalance(int64 chainID) external view returns (uint256 result) {
        return prototype.getGasStabilityPoolBalance(chainID);
    }
}