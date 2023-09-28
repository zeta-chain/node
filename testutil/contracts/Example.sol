// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

// Sample contract for evm tests
contract Example {
    error Foo();

    uint256 public bar;

    constructor() {
        bar = 0;
    }

    function doRevert() external {
        revert Foo();
    }

    function doRevertWithMessage() external {
        revert("foo");
    }

    function doRevertWithRequire() external {
        require(false, "foo");
    }

    function doSucceed() external {
        bar = 1;
    }

    function setBar(uint256 _bar) external {
        bar = _bar;
    }
}