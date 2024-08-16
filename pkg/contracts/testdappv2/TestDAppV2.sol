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

    /// @notice Struct containing revert context passed to onRevert.
    /// @param asset Address of asset, empty if it's gas token.
    /// @param amount Amount specified with the transaction.
    /// @param revertMessage Arbitrary data sent back in onRevert.
    struct RevertContext {
        address asset;
        uint64 amount;
        bytes revertMessage;
    }

    // these structures allow to assess contract calls
    mapping(bytes32 => bool) public calledWithMessage;
    mapping(bytes32 => uint256) public amountWithMessage;

    function setCalledWithMessage(string memory message) internal {
        calledWithMessage[keccak256(abi.encodePacked(message))] = true;
    }
    function setAmountWithMessage(string memory message, uint256 amount) internal {
        amountWithMessage[keccak256(abi.encodePacked(message))] = amount;
    }

    function getCalledWithMessage(string memory message) public view returns (bool) {
        return calledWithMessage[keccak256(abi.encodePacked(message))];
    }

    function getAmountWithMessage(string memory message) public view returns (uint256) {
        return amountWithMessage[keccak256(abi.encodePacked(message))];
    }

    // Universal contract interface
    function onCrossChainCall(
        zContext calldata _context,
        address _zrc20,
        uint256 amount,
        bytes calldata message
    )
    external
    {
        require(!isRevertMessage(string(message)));

        setCalledWithMessage(string(message));
        setAmountWithMessage(string(message), amount);
    }

    // called with gas token
    function gasCall(string memory message) external payable {
        // Revert if the message is "revert"
        require(!isRevertMessage(message));

        setCalledWithMessage(message);
        setAmountWithMessage(message, msg.value);
    }

    // called with ERC20 token
    function erc20Call(IERC20 erc20, uint256 amount, string memory message) external {
        require(!isRevertMessage(message));
        require(erc20.transferFrom(msg.sender, address(this), amount));

        setCalledWithMessage(message);
        setAmountWithMessage(message, amount);
    }

    // called without token
    function simpleCall(string memory message) external {
        require(!isRevertMessage(message));

        setCalledWithMessage(message);
        setAmountWithMessage(message, 0);
    }

    // used to make functions revert
    function isRevertMessage(string memory message) internal pure returns (bool) {
        return keccak256(abi.encodePacked(message)) == keccak256(abi.encodePacked("revert"));
    }

    // Revertable interface
    function onRevert(RevertContext calldata revertContext) external {
        setCalledWithMessage(string(revertContext.revertMessage));
        setAmountWithMessage(string(revertContext.revertMessage), 0);
    }

    receive() external payable {}
}