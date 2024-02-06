// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

interface IERC20 {
    function transferFrom(address sender, address recipient, uint256 amount) external returns (bool);
    function transfer(address recipient, uint256 amount) external returns (bool);
    function balanceOf(address account) external view returns (uint256);
}

// Sample contract that locks and unlocks tokens
contract Vault {
    mapping(address => mapping(address => uint256)) public balances;

    function deposit(address tokenAddress, uint256 amount) external {
        require(amount > 0, "Amount should be greater than 0");

        IERC20 token = IERC20(tokenAddress);
        require(token.transferFrom(msg.sender, address(this), amount), "Transfer failed");

        balances[msg.sender][tokenAddress] += amount;
    }

    function withdraw(address tokenAddress, uint256 amount) external {
        require(amount > 0, "Amount should be greater than 0");
        require(balances[msg.sender][tokenAddress] >= amount, "Insufficient balance");

        balances[msg.sender][tokenAddress] -= amount;

        IERC20 token = IERC20(tokenAddress);
        require(token.transfer(msg.sender, amount), "Transfer failed");
    }
}