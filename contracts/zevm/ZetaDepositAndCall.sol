// SPDX-License-Identifier: MIT
pragma solidity 0.8.7;

interface IZRC4 {
    function totalSupply() external view returns (uint256);
    function balanceOf(address account) external view returns (uint256);
    function transfer(address recipient, uint256 amount) external returns (bool);
    function allowance(address owner, address spender) external view returns (uint256);
    function approve(address spender, uint256 amount) external returns (bool);
    function transferFrom(address sender, address recipient, uint256 amount) external returns (bool);

    function deposit(address to, uint256 amount) external returns (bool);
    function withdraw(bytes memory to, uint256 amount) external returns (bool);

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);
    event Deposit(bytes from, address indexed to, uint256 value);
    event Withdrawal(address indexed from, bytes to, uint256 value);
}

interface zContract {
    function onCCC(address zrc4, uint256 amount, bytes calldata message) external;
}

contract ZetaDepositAndCall {
    address public FUNGIBLE_MODULE_ADDRESS;

    constructor() {
        FUNGIBLE_MODULE_ADDRESS = msg.sender;
    }

    function DepositAndCall(address zrc4, uint256 amount, address target, bytes calldata message) external {
        require(msg.sender == FUNGIBLE_MODULE_ADDRESS);
        require(target != FUNGIBLE_MODULE_ADDRESS && target != address(this));

        IZRC4(zrc4).deposit(target, amount);
        zContract(target).onCCC(zrc4, amount, message);
    }

   
}
