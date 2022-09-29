// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

interface ISystem {
    function FUNGIBLE_MODULE_ADDRESS() external view returns (address);
    function wzetaContractAddress() external view returns (address);
    function uniswapv2FactoryAddress() external view returns (address);
    function gasPrice(uint256) external view returns (uint256); // chainid => gas price
    function gasCoinERC4(uint256) external view returns (address);  // chainid => gas coin erc4
}

interface IZRC4 {
    function totalSupply() external view returns (uint256);
    function balanceOf(address account) external view returns (uint256);
    function transfer(address recipient, uint256 amount) external returns (bool);
    function allowance(address owner, address spender) external view returns (uint256);
    function approve(address spender, uint256 amount) external returns (bool);
    function transferFrom(address sender, address recipient, uint256 amount) external returns (bool);

    function deposit(address to, uint256 amount) external returns (bool);
    function withdraw(bytes memory to, uint256 amount) external returns (bool);

    function withdrawGasFee() external view returns (address,uint256);

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);
    event Deposit(bytes from, address indexed to, uint256 value);
    event Withdrawal(address indexed from, bytes to, uint256 value, uint256 gasfee);
}

abstract contract Context {
    function _msgSender() internal view virtual returns (address) {
        return msg.sender;
    }

    function _msgData() internal view virtual returns (bytes calldata) {
        return msg.data;
    }
}

interface IZRC4Metadata is IZRC4 {
    function name() external view returns (string memory);
    function symbol() external view returns (string memory);
    function decimals() external view returns (uint8);
}


enum CoinType {
    Zeta, // this should not be used
    Gas,
    ERC20
}

interface zContract {
    function onCrossChainCall(address zrc4, uint256 amount, bytes calldata message) external;
}