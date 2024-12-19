// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

// TestGasConsumer is a contract used to simulate high gas consumption
contract TestGasConsumer {
    // used to simulate gas consumption
    uint256[] private storageArray;

    struct zContext {
        bytes origin;
        address sender;
        uint256 chainID;
    }

    // Universal contract interface on ZEVM
    function onCall(
        zContext calldata _context,
        address _zrc20,
        uint256 _amount,
        bytes calldata _message
    )
    external
    {
        consumeGas();
    }

    function consumeGas() internal {
        // Approximate target gas consumption
        uint256 targetGas = 5000000;
        // Approximate gas cost for a single storage write
        uint256 storageWriteGasCost = 20000;
        uint256 iterations = targetGas / storageWriteGasCost;

        // Perform the storage writes
        for (uint256 i = 0; i < iterations; i++) {
            storageArray.push(i);
        }

        // Reset the storage array to avoid accumulation of storage cost
        delete storageArray;
    }
}