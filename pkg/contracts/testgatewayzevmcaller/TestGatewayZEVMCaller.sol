// SPDX-License-Identifier: MIT
pragma solidity 0.8.10;

struct CallOptions {
    uint256 gasLimit;
    bool isArbitraryCall;
}

struct RevertOptions {
    address revertAddress;
    bool callOnRevert;
    address abortAddress;
    bytes revertMessage;
    uint256 onRevertGasLimit;
}

interface IGatewayZEVM {
    function call(
        bytes memory receiver,
        address zrc20,
        bytes calldata message,
        CallOptions calldata callOptions,
        RevertOptions calldata revertOptions
    )
        external;
}

interface IZRC20 {
    function approve(address spender, uint256 amount) external returns (bool);
}

contract TestGatewayZEVMCaller {
    IGatewayZEVM private gatewayZEVM;
    constructor(address gatewayZEVMAddress) {
        gatewayZEVM = IGatewayZEVM(gatewayZEVMAddress);
    }

    function callGatewayZEVM(
        bytes memory receiver,
        address zrc20,
        bytes calldata message,
        CallOptions calldata callOptions,
        RevertOptions calldata revertOptions
    ) external {
        IZRC20(zrc20).approve(address(gatewayZEVM), 100000000000000000);
        gatewayZEVM.call(receiver, zrc20, message, callOptions, revertOptions);
    }
}