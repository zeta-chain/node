// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

struct AbortContext {
    bytes sender;
    address asset;
    uint256 amount;
    bool outgoing;
    uint256 chainID;
    bytes revertMessage;
}

contract TestAbort {
    // allow to assess onAbort calls
    mapping(bytes32 => AbortContext) public abortedWithMessage;

    function setAbortedWithMessage(string memory message, AbortContext memory abortContext) internal {
        abortedWithMessage[keccak256(abi.encodePacked(message))] = abortContext;
    }

    function getAbortedWithMessage(string memory message) public view returns (AbortContext memory) {
        return abortedWithMessage[keccak256(abi.encodePacked(message))];
    }

    function onAbort(AbortContext calldata abortContext) external {
        setAbortedWithMessage(string(abortContext.revertMessage), abortContext);
    }

    fallback() external payable {}

    receive() external payable {}
}
