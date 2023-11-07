// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

interface IZRC20 {
    function approve(address spender, uint256 amount) external;
    function transferFrom(address sender, address recipient, uint256 amount) external;
    function withdraw(bytes memory to, uint256 amount) external returns (bool);
}

// Sample contract for running withdraw on zEVM
contract Withdrawer {
    // Run n withdraws of amount on asset to custody
    function runWithdraws(
        bytes calldata recipient,
        IZRC20 asset,
        uint256 amount,
        bytes calldata message,
        uint256 count
    ) external {
        asset.transferFrom(msg.sender, address(this), amount * count);
        for (uint256 i = 0; i < count; i++) {
            asset.withdraw(recipient, amount);
        }
    }
}