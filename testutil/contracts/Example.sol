// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

// Sample contract for evm tests
contract Example {
    error Foo();

    // always reverts
    function doRevert() external {
        revert Foo();
    }
}