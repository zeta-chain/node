// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

import "interfaces/IUniswapV2Router02.sol";

struct Context {
    bytes origin;
    address sender;
    uint256 chainID;
}

interface zContract {
    function onCrossChainCall(
        Context calldata context,
        address zrc20,
        uint256 amount,
        bytes calldata message
    ) external;
}
interface IZRC20 {
    function totalSupply() external view returns (uint256);

    function balanceOf(address account) external view returns (uint256);

    function transfer(address recipient, uint256 amount) external returns (bool);

    function allowance(address owner, address spender) external view returns (uint256);

    function approve(address spender, uint256 amount) external returns (bool);

    function transferFrom(
        address sender,
        address recipient,
        uint256 amount
    ) external returns (bool);

    function deposit(address to, uint256 amount) external returns (bool);

    function burn(address account, uint256 amount) external returns (bool);

    function withdraw(bytes memory to, uint256 amount) external returns (bool);

    function withdrawGasFee() external view returns (address, uint256);

    function PROTOCOL_FEE() external view returns (uint256);

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);
    event Deposit(bytes from, address indexed to, uint256 value);
    event Withdrawal(address indexed from, bytes to, uint256 value, uint256 gasFee, uint256 protocolFlatFee);
}




contract ZEVMSwapApp is zContract {
    error InvalidSender();
    error LowAmount();

    uint256 constant private _DEADLINE = 1 << 64;
    address immutable public router02;
    address immutable public systemContract;
    
    constructor(address router02_, address systemContract_) {
        router02 = router02_;
        systemContract = systemContract_;
    }

    function encodeMemo(address targetZRC20,  bytes calldata recipient) pure external returns (bytes memory) {
//        return abi.encode(targetZRC20, recipient, minAmountOut);
        return abi.encodePacked(targetZRC20, recipient);
    }

    // data
    function decodeMemo(bytes calldata data) pure public returns (address, bytes memory) {
        bytes memory decodedBytes;
        uint256 size;
        size = data.length;
        address addr;
        addr = address(uint160(bytes20(data[0:20])));
        decodedBytes = data[20:];

        return (addr, decodedBytes);
    }

    
    // Call this function to perform a cross-chain swap
    function onCrossChainCall(Context calldata, address zrc20, uint256 amount, bytes calldata message) external override {
        if (msg.sender != systemContract) {
            revert InvalidSender();
        }
        address targetZRC20;
        bytes memory recipient;
        (targetZRC20, recipient) = decodeMemo(message);
        address[] memory path;
        path = new address[](2);
        path[0] = zrc20;
        path[1] = targetZRC20;
        // Approve the usage of this token by router02
        IZRC20(zrc20).approve(address(router02), amount);
        // Swap for your target token
        uint256[] memory amounts = IUniswapV2Router02(router02).swapExactTokensForTokens(amount, 0, path, address(this), _DEADLINE);

        // this contract subsides withdraw gas fee
        (address gasZRC20Addr,uint256 gasFee) = IZRC20(targetZRC20).withdrawGasFee();
        IZRC20(gasZRC20Addr).approve(address(targetZRC20), gasFee);
        IZRC20(targetZRC20).approve(address(targetZRC20), amounts[1]); // this does not seem to be necessary
        IZRC20(targetZRC20).withdraw(recipient, amounts[1]-gasFee);
    }
}