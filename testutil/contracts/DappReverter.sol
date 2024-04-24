// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

// DappReverter is a contract that can be used to test the reversion of a cross-chain call.
// It implements the onZetaMessage and onZetaRevert functions, which are called the ZEVM connector
contract DappReverter {
    function onZetaMessage() external{}
    function onZetaRevert() external {}
}