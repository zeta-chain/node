// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

interface ZetaInterfaces {
    /**
     * @dev Use SendInput to interact with the Connector: connector.send(SendInput)
     */
    struct SendInput {
        /// @dev Chain id of the destination chain. More about chain ids https://docs.zetachain.com/learn/glossary#chain-id
        uint256 destinationChainId;
        /// @dev Address receiving the message on the destination chain (expressed in bytes since it can be non-EVM)
        bytes destinationAddress;
        /// @dev Gas limit for the destination chain's transaction
        uint256 destinationGasLimit;
        /// @dev An encoded, arbitrary message to be parsed by the destination contract
        bytes message;
        /// @dev ZETA to be sent cross-chain + ZetaChain gas fees + destination chain gas fees (expressed in ZETA)
        uint256 zetaValueAndGas;
        /// @dev Optional parameters for the ZetaChain protocol
        bytes zetaParams;
    }

    /**
     * @dev Our Connector calls onZetaMessage with this struct as argument
     */
    struct ZetaMessage {
        bytes zetaTxSenderAddress;
        uint256 sourceChainId;
        address destinationAddress;
        /// @dev Remaining ZETA from zetaValueAndGas after subtracting ZetaChain gas fees and destination gas fees
        uint256 zetaValue;
        bytes message;
    }

    /**
     * @dev Our Connector calls onZetaRevert with this struct as argument
     */
    struct ZetaRevert {
        address zetaTxSenderAddress;
        uint256 sourceChainId;
        bytes destinationAddress;
        uint256 destinationChainId;
        /// @dev Equals to: zetaValueAndGas - ZetaChain gas fees - destination chain gas fees - source chain revert tx gas fees
        uint256 remainingZetaValue;
        bytes message;
    }
}

// Dapp is a sample comtract that implements ZetaReceiver and is used for unit testing
// It sets the values of the ZetaMessage struct to its public variables which can then be queried to check if the function was called correctly
contract Dapp {
    bytes public zetaTxSenderAddress;
    uint256 public sourceChainId;
    address public destinationAddress;
    uint256 public destinationChainId;
    uint256 public zetaValue;
    bytes public  message;

    constructor() {
        zetaTxSenderAddress = "";
        sourceChainId = 0;
        destinationAddress = address(0);
        destinationChainId = 0;
        zetaValue = 0;
        message = "";
    }

    function onZetaMessage(ZetaInterfaces.ZetaMessage calldata zetaMessage) external{
        zetaTxSenderAddress = zetaMessage.zetaTxSenderAddress;
        sourceChainId = zetaMessage.sourceChainId;
        destinationAddress = zetaMessage.destinationAddress;
        zetaValue = zetaMessage.zetaValue;
        message = zetaMessage.message;
    }
    function onZetaRevert(ZetaInterfaces.ZetaRevert calldata zetaRevert) external {
        zetaTxSenderAddress = abi.encodePacked(zetaRevert.zetaTxSenderAddress);
        sourceChainId = zetaRevert.sourceChainId;
        destinationAddress = address(uint160(uint256(keccak256(zetaRevert.destinationAddress))));
        destinationChainId = zetaRevert.destinationChainId;
        zetaValue = zetaRevert.remainingZetaValue;
        message = zetaRevert.message;
    }
}