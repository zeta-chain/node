// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

// Sample contract for evm tests
contract Reverter {
    error Foo();

    struct zContext {
        bytes origin;
        address sender;
        uint256 chainID;
    }

    function onCrossChainCall(
        zContext calldata context,
        address zrc20,
        uint256 amount,
        bytes calldata message
    ) external {
        onCall(context, zrc20, amount, message);
    }

    function onCall(
        zContext calldata context,
        address zrc20,
        uint256 amount,
        bytes calldata message
    ) public {
        revert Foo();
    }
}