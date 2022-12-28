// SPDX-License-Identifier: GPL-3.0

pragma solidity 0.8.7;

contract ZETABridge {
    event ZetaSent(bytes to, uint256 toChainID, uint256 value);

    function sendZeta(bytes memory to, uint256 toChainID) payable public {
        emit ZetaSent(to, toChainID, msg.value);
    }
}