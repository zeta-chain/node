// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

interface ZetaInterfaces {
    struct SendInput {
        uint256 destinationChainId;
        bytes destinationAddress;
        uint256 gasLimit;
        bytes message;
        uint256 zetaAmount;
        bytes zetaParams;
    }

    struct ZetaMessage {
        bytes originSenderAddress;
        uint256 originChainId;
        address destinationAddress;
        uint256 zetaAmount;
        bytes message;
    }

    struct ZetaRevert {
        address originSenderAddress;
        uint256 originChainId;
        bytes destinationAddress;
        uint256 destinationChainId;
        uint256 zetaAmount;
        bytes message;
    }
}
