// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

import "../interfaces/zContract.sol";
import "../interfaces/IZRC20.sol";
import "../interfaces/IUniswapV2Router02.sol";

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
    
    // Call this function to perform a cross-chain swap
    function onCrossChainCall(address zrc20, uint256 amount, bytes calldata message) external override {
        if (msg.sender != systemContract) {
            revert InvalidSender();
        }
        address targetZRC20;
        address recipient;
        uint256 minAmountOut; 
        (targetZRC20, recipient, minAmountOut) = abi.decode(message, (address,address,uint256));
        address[] memory path;
        path = new address[](2);
        path[0] = zrc20;
        path[1] = targetZRC20;
        // Approve the usage of this token by router02
        IZRC20(zrc20).approve(address(router02), amount);
        // Swap for your target token
        uint256[] memory amounts = IUniswapV2Router02(router02).swapExactTokensForTokens(amount, minAmountOut, path, address(this), _DEADLINE);
        // Withdraw amount to target recipient
        (, uint256 gasFee) = IZRC20(targetZRC20).withdrawGasFee();
        if (gasFee > amounts[1]) {
            revert LowAmount();
        }
        IZRC20(targetZRC20).approve(address(targetZRC20), gasFee);
        IZRC20(targetZRC20).withdraw(abi.encodePacked(recipient), amounts[1] - gasFee);
    }
}