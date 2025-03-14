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

struct RevertOptions {
    address revertAddress;
    bool callOnRevert;
    address abortAddress;
    bytes revertMessage;
    uint256 onRevertGasLimit;
}

interface IGatewayZEVM {
    function withdraw(
        bytes memory receiver,
        uint256 amount,
        address zrc20,
        RevertOptions calldata revertOptions
    ) external;
}

interface IZRC20 {
    function approve(address spender, uint256 amount) external returns (bool);
    function withdrawGasFee() external view returns (address, uint256);
}

contract TestAbort {
    // allow to assess onAbort calls
    mapping(bytes32 => AbortContext) public abortedWithMessage;
    bool public aborted;

    function setAbortedWithMessage(string memory message, AbortContext memory abortContext) internal {
        abortedWithMessage[keccak256(abi.encodePacked(message))] = abortContext;
        aborted = true;
    }

    function getAbortedWithMessage(string memory message) public view returns (AbortContext memory) {
        return abortedWithMessage[keccak256(abi.encodePacked(message))];
    }

    function isAborted() public view returns (bool) {
        return aborted;
    }

    function onAbort(AbortContext calldata abortContext) external {
        setAbortedWithMessage(string(abortContext.revertMessage), abortContext);

        // withdraw funds to the sender on connected chain
        if (isWithdrawMessage(string(abortContext.revertMessage))) {
            (address feeToken, uint256 feeAmount) = IZRC20(abortContext.asset).withdrawGasFee();
            require(feeToken == abortContext.asset, "zrc20 is not gas token");
            require(feeAmount <= abortContext.amount, "fee amount is higher than the amount");
            uint256 withdrawAmount = abortContext.amount - feeAmount;

            IZRC20(abortContext.asset).approve(msg.sender, abortContext.amount);

            // caller is the gateway
            IGatewayZEVM(msg.sender).withdraw(
                abi.encode(abortContext.sender),
                withdrawAmount,
                abortContext.asset,
                RevertOptions(address(0), false, address(0), "", 0)
            );
        }
    }

    function isWithdrawMessage(string memory message) internal pure returns (bool) {
        return keccak256(abi.encodePacked(message)) == keccak256(abi.encodePacked("withdraw"));
    }

    fallback() external payable {}

    receive() external payable {}
}
