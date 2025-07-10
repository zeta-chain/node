// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

contract TestDAppEmpty {
    struct zContext {
        bytes sender;
        address senderEVM;
        uint256 chainID;
    }

    /// @notice Struct containing revert context passed to onRevert.
    /// @param sender Address of account that initiated smart contract call.
    /// @param asset Address of asset, empty if it's gas token.
    /// @param amount Amount specified with the transaction.
    /// @param revertMessage Arbitrary data sent back in onRevert.
    struct RevertContext {
        address sender;
        address asset;
        uint256 amount;
        bytes revertMessage;
    }

    /// @notice Message context passed to execute function.
    /// @param sender Sender from omnichain contract.
    struct MessageContext {
        address sender;
    }

    // the constructor is used to determine if the chain is ZetaChain
    constructor() {
    }

    // Universal contract interface on ZEVM
    function onCall(
        zContext calldata context,
        address _zrc20,
        uint256 amount,
        bytes calldata message
    )
    external
    {
    }


    // Revertable interface
    function onRevert(RevertContext calldata revertContext) external {
    }

    // Callable interface on connected EVM chains
    function onCall(MessageContext calldata messageContext, bytes calldata message) external payable returns (bytes memory) {
        return "";
    }

    receive() external payable {}
}