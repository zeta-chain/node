// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

// Opcode is a smart contract used to test if specific opcodes or other EVM functionalities are supported by the ZEVM
contract Opcode {
    function testPUSH0() public returns (uint256) {
        assembly {
            let result := 0  // Compiler uses PUSH0 if supported
            mstore(0, result)
            return(0, 32)
        }
    }

    function testTLOAD() public returns (uint256) {
        uint256 result;
        assembly {
            // Store value in transient storage
            tstore(0x00, 0x1234)

            // Load it back with TLOAD
            result := tload(0x00)
        }
        return result; // Should return 0x1234 if TLOAD works
    }
}