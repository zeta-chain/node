// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

// Sample contract for evm tests
contract Example {
    error Foo();

    struct zContext {
        bytes sender;
        address senderEVM;
        uint256 chainID;
    }

    uint256 public bar;
    bytes public lastMessage;
    bytes public lastSender;

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
        bar = amount;
        lastMessage = message;
        lastSender = context.sender;
    }
}