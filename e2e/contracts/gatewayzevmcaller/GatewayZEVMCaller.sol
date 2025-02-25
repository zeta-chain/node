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

    function withdrawAndCall(
        bytes memory receiver,
        uint256 amount,
        uint256 chainId,
        bytes calldata message,
        CallOptions calldata callOptions,
        RevertOptions calldata revertOptions
    )
        external;

    function withdrawAndCall(
        bytes memory receiver,
        uint256 amount,
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

interface WZETA {
    function deposit() external payable;
    function approve(address guy, uint256 wad) external returns (bool);
}

contract GatewayZEVMCaller {
    IGatewayZEVM private gatewayZEVM;
    WZETA wzeta;
    constructor(address gatewayZEVMAddress, address wzetaAddress) {
        gatewayZEVM = IGatewayZEVM(gatewayZEVMAddress);
        wzeta = WZETA(wzetaAddress);
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

    function withdrawAndCallGatewayZEVM(
        bytes memory receiver,
        uint256 amount,
        uint256 chainId,
        bytes calldata message,
        CallOptions calldata callOptions,
        RevertOptions calldata revertOptions
    ) external {
        wzeta.approve(address(gatewayZEVM), amount);
        gatewayZEVM.withdrawAndCall(receiver, amount, chainId, message, callOptions, revertOptions);
    }

    function withdrawAndCallGatewayZEVM(
        bytes memory receiver,
        uint256 amount,
        address zrc20,
        bytes calldata message,
        CallOptions calldata callOptions,
        RevertOptions calldata revertOptions
    ) external {
        IZRC20(zrc20).approve(address(gatewayZEVM), 100000000000000000);
        gatewayZEVM.withdrawAndCall(receiver, amount, zrc20, message, callOptions, revertOptions);
    }

    function depositWZETA() external payable {
        wzeta.deposit{value: msg.value}();
    }
}