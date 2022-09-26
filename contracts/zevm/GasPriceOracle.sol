// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;
import "./ifaces.sol";

contract GasPriceOracle {
    mapping (uint256 => uint256) public gasPrice; // chainid => gas price
    mapping (uint256 => address) public gasCoinERC4;  // chainid => gas coin erc4
    address public constant FUNGIBLE_MODULE_ADDRESS = 0x735b14BB79463307AAcBED86DAf3322B1e6226aB;

    event Deployed();
    event SetGasPrice(uint256, uint256);
    event SetGasCoin(uint256, address);

    constructor() {
        require(msg.sender == FUNGIBLE_MODULE_ADDRESS, "only fungible module can deploy");
        emit Deployed();
    }

    // fungible module updates the gas price oracle periodically
    function setGasPrice(uint256 chainID, uint256 price) external {
        require(msg.sender == FUNGIBLE_MODULE_ADDRESS, "Only fungible module can set gas price");
        gasPrice[chainID] = price;
        emit SetGasPrice(chainID, price);
    }

    function setGasCoinERC4(uint256 chainID, address erc4) external {
        require(msg.sender == FUNGIBLE_MODULE_ADDRESS, "Only fungible module can set gas coin erc4");
        gasCoinERC4[chainID] = erc4;
        emit SetGasCoin(chainID, erc4);
    }
}