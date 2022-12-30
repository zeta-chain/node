// SPDX-License-Identifier: GPL-3.0

pragma solidity 0.8.7;

contract ZETABridge {
    address public constant FUNGIBLE_MODULE_ADDRESS = payable(0x735b14BB79463307AAcBED86DAf3322B1e6226aB);

    event ZetaSent(bytes to, uint256 toChainID, uint256 value);

    function sendZeta(bytes calldata to, uint256 toChainID) external payable {
        (bool success, ) = FUNGIBLE_MODULE_ADDRESS.call{value: msg.value}("");
        require(success, "ZETABridge: failed to transfer");
        emit ZetaSent(to, toChainID, msg.value);
    }
}