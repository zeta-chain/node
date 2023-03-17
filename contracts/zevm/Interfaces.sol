// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

/**
 * @dev Interfaces of SystemContract and ZRC20 to make easier to import.
 */
interface ISystem {
    function FUNGIBLE_MODULE_ADDRESS() external view returns (address);
    function wZetaContractAddress() external view returns (address);
    function uniswapv2FactoryAddress() external view returns (address);
    function gasPriceByChainId(uint256 chainID) external view returns (uint256);
    function gasCoinZRC20ByChainId(uint256 chainID) external view returns (address);
    function gasZetaPoolByChainId(uint256 chainID) external view returns (address);
}
interface IZRC20 {
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
    event Withdrawal(address indexed from, bytes to, uint256 value, uint256 gasfee, uint256 protocolFlatFee);
}

abstract contract Context {
    function _msgSender() internal view virtual returns (address) {
        return msg.sender;
    }

    function _msgData() internal view virtual returns (bytes calldata) {
        return msg.data;
    }
}

interface IZRC20Metadata is IZRC20 {
    function name() external view returns (string memory);
    function symbol() external view returns (string memory);
    function decimals() external view returns (uint8);
}


enum CoinType {
    Zeta, // this should not be used
    Gas,
    ERC20
}

/**
 * @dev Any Zeta Contract must implement this interface to allow SystemContract to interact with. 
 * It's only require if the contract wants to interact with other chains.
 */
interface zContract {
    function onCrossChainCall(address zrc20, uint256 amount, bytes calldata message) external;
}