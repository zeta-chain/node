// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

struct Context {
    bytes origin;
    address sender;
    uint256 chainID;
}

interface IZRC20 {
    function approve(address spender, uint256 amount) external returns (bool);
    function withdraw(bytes memory to, uint256 amount) external returns (bool);
}

// Withdrawer is a simple contract performing a withdraw of deposited ZRC20
// The amount to withdraw can be set during the contract deployment, it also to tests some edge cases like withdrawing BTC dust amount
contract Withdrawer {
    uint256 immutable public withdrawAmount;
    
    constructor(uint256 _withdrawAmount) {
        withdrawAmount = _withdrawAmount;
    }

    // perform a withdraw on cross chain call
    function onCall(Context calldata context, address zrc20, uint256, bytes calldata) external {
        // perform withdrawal with the target token
        IZRC20(zrc20).approve(address(zrc20), type(uint256).max);
        IZRC20(zrc20).withdraw(context.origin, withdrawAmount);
    }
}