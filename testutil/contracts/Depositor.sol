// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

interface IERC20 {
    function approve(address spender, uint256 amount) external;
}


interface ERC20Custody {
    function deposit(
        bytes calldata recipient,
        IERC20 asset,
        uint256 amount,
        bytes calldata message
    ) external;
}

// Sample contract for running deposit on EVM
contract Depositor {
    ERC20Custody immutable private _custody;

    constructor(address custody_) {
        _custody = ERC20Custody(custody_);
    }

    // Run n deposits of amount on asset to custody
    function runDeposits(
        bytes calldata recipient,
        IERC20 asset,
        uint256 amount,
        bytes calldata message,
        uint256 count
    ) external {
        asset.approve(address(_custody), amount * count);
        for (uint256 i = 0; i < count; i++) {
            _custody.deposit(recipient, asset, amount, message);
        }
    }
}