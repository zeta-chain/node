// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.7;

import "./IRegular.sol";

//IRegular constant IREGULAR_CONTRACT = IRegular(IREGULAR_PRECOMPILE_ADDRESS);
contract RegularCaller {
    function testBech32ToHexAddr() public view returns (bool) {
        // Test input and expected output
        string
            memory testBech32 = "zeta1h8duy2dltz9xz0qqhm5wvcnj02upy887fyn43u";
        address expectedHexAddr = 0xB9Dbc229Bf588A613C00BEE8e662727AB8121cfE;

        // Call the precompiled contract
        address result = IREGULAR_CONTRACT.bech32ToHexAddr(testBech32);

        // Check if the result matches the expected output
        return result == expectedHexAddr;
    }

    function testBech32ify() public view returns (bool) {
        // Test input and expected output
        string memory testPrefix = "zeta";
        address testHexAddr = 0xB9Dbc229Bf588A613C00BEE8e662727AB8121cfE;
        string
            memory expectedBech32 = "zeta1h8duy2dltz9xz0qqhm5wvcnj02upy887fyn43u";

        // Call the precompiled contract
        string memory result = IREGULAR_CONTRACT.bech32ify(
            testPrefix,
            testHexAddr
        );

        // Check if the result matches the expected output
        return
            keccak256(abi.encodePacked(result)) ==
            keccak256(abi.encodePacked(expectedBech32));
    }

    function testRegularCall(
        string memory method,
        address addr
    ) public returns (uint256) {
        // Call the precompiled contract with the given method and address
        uint256 result = IREGULAR_CONTRACT.regularCall(method, addr);

        // Return the result
        return result;
    }
}
