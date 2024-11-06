// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

/// @title IGatewayZEVM
/// @notice Interface for the GatewayZEVM contract.
/// @dev Defines functions for cross-chain interactions and token handling.
interface IGatewayZEVM {
    /// @notice Withdraw ZRC20 tokens to an external chain.
    /// @param receiver The receiver address on the external chain.
    /// @param amount The amount of tokens to withdraw.
    /// @param zrc20 The address of the ZRC20 token.
    /// @param revertOptions Revert options.
    function withdraw(
        bytes memory receiver,
        uint256 amount,
        address zrc20,
        RevertOptions calldata revertOptions
    )
    external;
}

/// @title IZRC20
/// @notice Interface for the ZRC20 token contract.
interface IZRC20 {
    function approve(address spender, uint256 amount) external returns (bool);
    function withdrawGasFee() external view returns (address, uint256);
}

interface IERC20 {
    function transferFrom(address sender, address recipient, uint256 amount) external returns (bool);
}

/// @notice Struct containing revert options
/// @param revertAddress Address to receive revert.
/// @param callOnRevert Flag if onRevert hook should be called.
/// @param abortAddress Address to receive funds if aborted.
/// @param revertMessage Arbitrary data sent back in onRevert.
/// @param onRevertGasLimit Gas limit for revert tx, unused on GatewayZEVM methods
struct RevertOptions {
    address revertAddress;
    bool callOnRevert;
    address abortAddress;
    bytes revertMessage;
    uint256 onRevertGasLimit;
}

struct zContext {
    bytes origin;
    address sender;
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

contract TestDAppV2 {
    // used to simulate gas consumption
    uint256[] private storageArray;

    string public constant NO_MESSAGE_CALL = "called with no message";
    string public constant WITHDRAW = "withdraw";

    // these structures allow to assess contract calls
    mapping(bytes32 => bool) public calledWithMessage;
    mapping(bytes => address) public senderWithMessage;
    mapping(bytes32 => uint256) public amountWithMessage;

    // return the index used for the "WithMessage" mapping when the message for calls is empty
    // this allows testing the message with empty message
    // this function includes the sender of the message to avoid collisions when running parallel tests with different senders
    function getNoMessageIndex(address sender) pure public returns (string memory) {
        return string(abi.encodePacked(NO_MESSAGE_CALL, sender));
    }

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

    // Universal contract interface on ZEVM
    function onCall(
        zContext calldata context,
        address zrc20,
        uint256 amount,
        bytes calldata message
    )
    external
    {
        require(!isRevertMessage(string(message)));

        // if withdraw is specified in the call, withdraw back the funds
        if (isWithdrawMessage(string(message))) {
            // get the fees for the withdraw, note: zrc20 should be the gas token for now
            (address feeToken, uint256 feeAmount) = IZRC20(zrc20).withdrawGasFee();
            require(feeToken == zrc20, "zrc20 is not gas token");
            require(feeAmount <= amount, "fee amount is higher than the amount");
            uint256 withdrawAmount = amount - feeAmount;

            // in the E2E test doing a withdraw, we want to test the gas usage for processing the outbound
            // this function allow to represent more accurate use case where gas would be consumed in the onCall hook
            consumeGas();

            IZRC20(zrc20).approve(msg.sender, amount);

            // caller is the gateway
            IGatewayZEVM(msg.sender).withdraw(
                abi.encode(context.sender),
                withdrawAmount,
                zrc20,
                RevertOptions(address(0), false, address(0), "", 0)
            );
        } else {
            string memory messageStr = message.length == 0 ? getNoMessageIndex(context.sender) : string(message);

            setCalledWithMessage(messageStr);
            setAmountWithMessage(messageStr, amount);
        }
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

    // used to trigger withdraw back
    function isWithdrawMessage(string memory message) internal pure returns (bool) {
        return keccak256(abi.encodePacked(message)) == keccak256(abi.encodePacked(WITHDRAW));
    }

    // Revertable interface
    function onRevert(RevertContext calldata revertContext) external {
        setCalledWithMessage(string(revertContext.revertMessage));
        setAmountWithMessage(string(revertContext.revertMessage), 0);
        senderWithMessage[revertContext.revertMessage] = revertContext.sender;
    }

    // Callable interface on connected EVM chains
    function onCall(MessageContext calldata messageContext, bytes calldata message) external payable returns (bytes memory) {
        string memory messageStr = message.length == 0 ? getNoMessageIndex(messageContext.sender) : string(message);

        setCalledWithMessage(messageStr);
        setAmountWithMessage(messageStr, msg.value);
        senderWithMessage[bytes(messageStr)] = messageContext.sender;

        return "";
    }

    function consumeGas() internal {
        // Approximate target gas consumption
        uint256 targetGas = 500000;
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

    receive() external payable {}
}