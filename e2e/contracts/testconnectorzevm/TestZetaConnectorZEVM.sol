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

interface WZETA {
    function transferFrom(address src, address dst, uint wad) external returns (bool);

    function withdraw(uint wad) external;
}

contract TestZetaConnectorZEVM is ZetaInterfaces {
    /// @notice Contract custom errors.
    error OnlyWZETA();
    error WZETATransferFailed();
    error OnlyFungibleModule();
    error FailedZetaSent();

    /// @notice Fungible module address.
    address public constant FUNGIBLE_MODULE_ADDRESS = payable(0x735b14BB79463307AAcBED86DAf3322B1e6226aB);
    /// @notice WZETA token address.
    address public wzeta;

    event ZetaSent(
        address sourceTxOriginAddress,
        address indexed zetaTxSenderAddress,
        uint256 indexed destinationChainId,
        bytes destinationAddress,
        uint256 zetaValueAndGas,
        uint256 destinationGasLimit,
        bytes message,
        bytes zetaParams
    );
    event SetWZETA(address wzeta_);

    constructor(address wzeta_) {
        wzeta = wzeta_;
    }

    /// @dev Receive function to receive ZETA from WETH9.withdraw().
    receive() external payable {
        if (msg.sender != wzeta) revert OnlyWZETA();
    }

    /**
     * @dev Sends ZETA and bytes messages (to execute it) crosschain.
     * @param input, SendInput struct, checkout above.
     */
    function send(ZetaInterfaces.SendInput calldata input) external {
        // Transfer wzeta to "fungible" module, which will be burnt by the protocol post processing via hooks.
        if (!WZETA(wzeta).transferFrom(msg.sender, address(this), input.zetaValueAndGas)) revert WZETATransferFailed();
        WZETA(wzeta).withdraw(input.zetaValueAndGas);
        (bool sent, ) = FUNGIBLE_MODULE_ADDRESS.call{value: input.zetaValueAndGas}("");
        if (!sent) revert FailedZetaSent();
        emit ZetaSent(
            tx.origin,
            msg.sender,
            input.destinationChainId,
            input.destinationAddress,
            input.zetaValueAndGas,
            input.destinationGasLimit,
            input.message,
            input.zetaParams
        );
    }

    /**
     * @dev Sends ZETA and bytes messages (to execute it) crosschain.
     * @param wzeta_, new WZETA address.
     */
    function setWzetaAddress(address wzeta_) external {
        if (msg.sender != FUNGIBLE_MODULE_ADDRESS) revert OnlyFungibleModule();
        wzeta = wzeta_;
        emit SetWZETA(wzeta_);
    }

    function foo() external pure returns (string memory) {
        return "foo";
    }
}