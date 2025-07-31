// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// Opcode is a smart contract used to test if specific opcodes or other EVM functionalities are supported by the ZEVM
contract Opcode {
    function testPUSH0() public returns (uint256) {
        assembly {
            let result := 0  // Compiler uses PUSH0 if supported
            mstore(0, result)
            return(0, 32)
        }
    }
}