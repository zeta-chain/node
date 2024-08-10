// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

interface IERC20 {
    function transferFrom(address sender, address recipient, uint256 amount) external returns (bool);
}

contract TestDAppV2 {
    struct zContext {
        bytes origin;
        address sender;
        uint256 chainID;
    }

    zContext public lastContext;
    address public lastZRC20;
    uint256 public lastAmount;
    string public lastMessage;

    function onCrossChainCall(
        zContext calldata context,
        address zrc20,
        uint256 amount,
        bytes calldata message
    )
    external
    {
        require(!isRevertMessage(message));

        // Store the context and parameters
        lastContext = context;
        lastZRC20 = zrc20;
        lastAmount = amount;
        lastMessage = string(message);
    }

    function gasCall(string memory message) external payable {
        // Revert if the message is "revert"
        require(!isRevertMessage(bytes(message)));

        lastMessage = message;
        lastAmount = msg.value;
    }

    function erc20Call(IERC20 erc20, uint256 amount, string memory message) external {
        require(!isRevertMessage(bytes(message)));
        require(erc20.transferFrom(msg.sender, address(this), amount));

        lastMessage = message;
        lastAmount = amount;
    }

    function simpleCall(string memory message) external {
        require(!isRevertMessage(bytes(message)));

        lastMessage = message;
        lastAmount = 0;
    }

    function isRevertMessage(bytes memory message) internal pure returns (bool) {
        return keccak256(message) == keccak256(abi.encodePacked("revert"));
    }
}