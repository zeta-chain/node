package signer

import (
	"context"
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/zeta-chain/protocol-contracts/v1/pkg/contracts/evm/erc20custody.sol"
	connectorevm "github.com/zeta-chain/protocol-contracts/v1/pkg/contracts/evm/zetaconnector.base.sol"
	erc20custodyv2 "github.com/zeta-chain/protocol-contracts/v2/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayevm.sol"

	"github.com/zeta-chain/zetacore/zetaclient/chains/evm"
)

// SignConnectorOnReceive
// function onReceive(
//
//	bytes calldata originSenderAddress,
//	uint256 originChainId,
//	address destinationAddress,
//	uint zetaAmount,
//	bytes calldata message,
//	bytes32 internalSendHash
//
// ) external virtual {}
func (signer *Signer) SignConnectorOnReceive(ctx context.Context, txData *OutboundData) (*ethtypes.Transaction, error) {
	var data []byte
	var err error

	zetaConnectorABI, err := connectorevm.ZetaConnectorBaseMetaData.GetAbi()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get ZetaConnectorZEVMMetaData ABI")
	}

	data, err = zetaConnectorABI.Pack("onReceive",
		txData.sender.Bytes(),
		txData.srcChainID,
		txData.to,
		txData.amount,
		txData.message,
		txData.cctxIndex)
	if err != nil {
		return nil, errors.Wrap(err, "onReceive pack error")
	}

	tx, _, _, err := signer.Sign(
		ctx,
		data,
		signer.zetaConnectorAddress,
		zeroValue,
		txData.gasLimit,
		txData.gasPrice,
		txData.nonce,
		txData.height)
	if err != nil {
		return nil, errors.Wrap(err, "sign onReceive error")
	}

	return tx, nil
}

// SignConnectorOnRevert
// function onRevert(
// address originSenderAddress,
// uint256 originChainId,
// bytes calldata destinationAddress,
// uint256 destinationChainId,
// uint256 zetaAmount,
// bytes calldata message,
// bytes32 internalSendHash
// ) external override whenNotPaused onlyTssAddress
func (signer *Signer) SignConnectorOnRevert(ctx context.Context, txData *OutboundData) (*ethtypes.Transaction, error) {
	var data []byte
	var err error

	zetaConnectorABI, err := connectorevm.ZetaConnectorBaseMetaData.GetAbi()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get ZetaConnectorZEVMMetaData ABI")
	}

	data, err = zetaConnectorABI.Pack("onRevert",
		txData.sender,
		txData.srcChainID,
		txData.to.Bytes(),
		txData.toChainID,
		txData.amount,
		txData.message,
		txData.cctxIndex)
	if err != nil {
		return nil, errors.Wrap(err, "onRevert pack error")
	}

	tx, _, _, err := signer.Sign(
		ctx,
		data,
		signer.zetaConnectorAddress,
		zeroValue,
		txData.gasLimit,
		txData.gasPrice,
		txData.nonce,
		txData.height)
	if err != nil {
		return nil, errors.Wrap(err, "sign onRevert error")
	}

	return tx, nil
}

// SignCancel signs a transaction from TSS address to itself with a zero amount in order to increment the nonce
func (signer *Signer) SignCancel(ctx context.Context, txData *OutboundData) (*ethtypes.Transaction, error) {
	tx, _, _, err := signer.Sign(
		ctx,
		nil,
		signer.TSS().EVMAddress(),
		zeroValue, // zero out the amount to cancel the tx
		evm.EthTransferGasLimit,
		txData.gasPrice,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "SignCancel error")
	}

	return tx, nil
}

// SignGasWithdraw signs a withdrawal transaction sent from the TSS address to the destination
func (signer *Signer) SignGasWithdraw(ctx context.Context, txData *OutboundData) (*ethtypes.Transaction, error) {
	tx, _, _, err := signer.Sign(
		ctx,
		nil,
		txData.to,
		txData.amount,
		evm.EthTransferGasLimit,
		txData.gasPrice,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "SignGasWithdraw error")
	}

	return tx, nil
}

// SignERC20Withdraw
// function withdraw(
// address recipient,
// address asset,
// uint256 amount,
// ) external onlyTssAddress
func (signer *Signer) SignERC20Withdraw(ctx context.Context, txData *OutboundData) (*ethtypes.Transaction, error) {
	var data []byte
	var err error

	erc20CustodyV1ABI, err := erc20custody.ERC20CustodyMetaData.GetAbi()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get ERC20CustodyMetaData ABI")
	}

	data, err = erc20CustodyV1ABI.Pack("withdraw", txData.to, txData.asset, txData.amount)
	if err != nil {
		return nil, fmt.Errorf("withdraw pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(
		ctx,
		data,
		signer.er20CustodyAddress,
		zeroValue,
		txData.gasLimit,
		txData.gasPrice,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("sign withdraw error: %w", err)
	}

	return tx, nil
}

// SignWhitelistERC20Cmd signs a whitelist command for ERC20 token
func (signer *Signer) SignWhitelistERC20Cmd(
	ctx context.Context,
	txData *OutboundData,
	params string,
) (*ethtypes.Transaction, error) {
	outboundParams := txData.outboundParams
	erc20 := ethcommon.HexToAddress(params)
	if erc20 == (ethcommon.Address{}) {
		return nil, fmt.Errorf("SignCommandTx: invalid erc20 address %s", params)
	}
	custodyAbi, err := erc20custody.ERC20CustodyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	data, err := custodyAbi.Pack("whitelist", erc20)
	if err != nil {
		return nil, fmt.Errorf("whitelist pack error: %w", err)
	}
	tx, _, _, err := signer.Sign(
		ctx,
		data,
		txData.to,
		zeroValue,
		txData.gasLimit,
		txData.gasPrice,
		outboundParams.TssNonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("sign whitelist error: %w", err)
	}
	return tx, nil
}

// SignMigrateTssFundsCmd signs a migrate TSS funds command
func (signer *Signer) SignMigrateTssFundsCmd(ctx context.Context, txData *OutboundData) (*ethtypes.Transaction, error) {
	tx, _, _, err := signer.Sign(
		ctx,
		nil,
		txData.to,
		txData.amount,
		txData.gasLimit,
		txData.gasPrice,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("SignMigrateTssFundsCmd error: %w", err)
	}
	return tx, nil
}

// SignGatewayExecute signs a gateway execute
// used for gas withdrawal and call transaction
// function execute
// address destination,
// bytes calldata data
func (signer *Signer) SignGatewayExecute(ctx context.Context, txData *OutboundData) (*ethtypes.Transaction, error) {
	var data []byte
	var err error

	gatewayABI, err := gatewayevm.GatewayEVMMetaData.GetAbi()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get GatewayEVMMetaData ABI")
	}

	data, err = gatewayABI.Pack("execute", txData.to, txData.message)
	if err != nil {
		return nil, fmt.Errorf("execute pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(
		ctx,
		data,
		signer.gatewayAddress,
		txData.amount,
		txData.gasLimit,
		txData.gasPrice,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("sign execute error: %w", err)
	}

	return tx, nil
}

// SignERC20CustodyWithdraw signs a erc20 withdrawal transaction
// function withdrawAndCall
// address token,
// address to,
// uint256 amount,
func (signer *Signer) SignERC20CustodyWithdraw(
	ctx context.Context,
	txData *OutboundData,
) (*ethtypes.Transaction, error) {
	var data []byte
	var err error

	erc20CustodyV2ABI, err := erc20custodyv2.ERC20CustodyMetaData.GetAbi()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get ERC20CustodyMetaData ABI")
	}

	data, err = erc20CustodyV2ABI.Pack("withdraw", txData.asset, txData.to, txData.amount)
	if err != nil {
		return nil, fmt.Errorf("withdraw pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(
		ctx,
		data,
		signer.er20CustodyAddress,
		zeroValue,
		txData.gasLimit,
		txData.gasPrice,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("sign withdraw error: %w", err)
	}

	return tx, nil
}

// SignERC20CustodyWithdrawAndCall signs a erc20 withdrawal and call transaction
// function withdrawAndCall
// address token,
// address to,
// uint256 amount,
// bytes calldata data
func (signer *Signer) SignERC20CustodyWithdrawAndCall(
	ctx context.Context,
	txData *OutboundData,
) (*ethtypes.Transaction, error) {
	var data []byte
	var err error

	erc20CustodyV2ABI, err := erc20custodyv2.ERC20CustodyMetaData.GetAbi()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get ERC20CustodyMetaData ABI")
	}

	data, err = erc20CustodyV2ABI.Pack("withdrawAndCall", txData.asset, txData.to, txData.amount, txData.message)
	if err != nil {
		return nil, fmt.Errorf("withdraw pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(
		ctx,
		data,
		signer.er20CustodyAddress,
		zeroValue,
		txData.gasLimit,
		txData.gasPrice,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("sign withdrawAndCall error: %w", err)
	}

	return tx, nil
}
