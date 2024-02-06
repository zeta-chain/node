// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

// Sample contract for evm tests
contract Example {
    error Foo();

    struct zContext {
        bytes origin;
        address sender;
        uint256 chainID;
    }

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
        assert(false, "foo");
    }

    function doSucceed() external {
        bar = 1;
    }

    function onCrossChainCall(
        zContext calldata context,
        address zrc20,
        uint256 amount,
        bytes calldata message
    ) external {
        bar = amount;
    }
}