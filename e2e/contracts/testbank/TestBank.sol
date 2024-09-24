// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

// @dev Interface for IBank contract
interface IBank {
    function deposit(
        address zrc20,
        uint256 amount
    ) external returns (bool success);

    function withdraw(
        address zrc20,
        uint256 amount
    ) external returns (bool success);

    function balanceOf(
        address zrc20,
        address user
    ) external view returns (uint256 balance);
}

// @dev Call IBank contract functions
contract TestBank {
    IBank bank = IBank(0x0000000000000000000000000000000000000067);
    address immutable owner;

    constructor() {
        owner = msg.sender;
    }

    modifier onlyOwner() {
        require(msg.sender == owner);
        _;
    }

    function deposit(
        address zrc20,
        uint256 amount
    ) external onlyOwner returns (bool) {
        bool success = bank.deposit(zrc20, amount);
        return success;
    }

    function withdraw(
        address zrc20,
        uint256 amount
    ) external onlyOwner returns (bool) {
        bool success = bank.withdraw(zrc20, amount);
        return success;
    }

    function balanceOf(
        address zrc20,
        address user
    ) external view onlyOwner returns (uint256) {
        uint256 balance = bank.balanceOf(zrc20, user);
        return balance;
    }

    fallback() external payable {}

    receive() external payable {}
}
