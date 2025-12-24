package signer

import (
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	erc20custodyv2 "github.com/zeta-chain/protocol-contracts-evm/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayevm.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/revert.sol"
	connector "github.com/zeta-chain/protocol-contracts-evm/pkg/zetaconnectornative.sol"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(errors.Wrap(err, "must"))
	}
	return v
}

var connectorABI = must(connector.ZetaConnectorNativeMetaData.GetAbi())
var gatewayABI = must(gatewayevm.GatewayEVMMetaData.GetAbi())
var erc20CustodyV2ABI = must(erc20custodyv2.ERC20CustodyMetaData.GetAbi())

// signGatewayExecute signs a gateway execute
// used for gas withdrawal and call transaction
// function execute
// address destination,
// bytes calldata data
func (signer *Signer) signGatewayExecute(txData *OutboundData) (*ethtypes.Transaction, error) {
	messageContext, err := txData.MessageContext()
	if err != nil {
		return nil, err
	}

	var data []byte

	data, err = gatewayABI.Pack("execute", messageContext, txData.to, txData.message)
	if err != nil {
		return nil, errors.Wrap(err, "execute pack error")
	}

	tx, _, _, err := signer.Sign(
		data,
		signer.gatewayAddress,
		txData.amount,
		txData.gas,
		txData.nonce,
	)
	if err != nil {
		return nil, errors.Wrap(err, "sign execute error")
	}

	return tx, nil
}

// signGatewayExecuteRevert signs a gateway execute revert
// function executeRevert
// address destination,
// bytes calldata data
func (signer *Signer) signGatewayExecuteRevert(
	inboundSender string,
	txData *OutboundData,
) (*ethtypes.Transaction, error) {
	data, err := gatewayABI.Pack(
		"executeRevert",
		txData.to,
		txData.message,
		revert.RevertContext{
			Sender:        common.HexToAddress(inboundSender),
			Asset:         txData.asset,
			Amount:        txData.amount,
			RevertMessage: txData.revertOptions.RevertMessage,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "executeRevert pack error")
	}

	tx, _, _, err := signer.Sign(
		data,
		signer.gatewayAddress,
		txData.amount,
		txData.gas,
		txData.nonce,
	)
	if err != nil {
		return nil, errors.Wrap(err, "sign executeRevert error")
	}

	return tx, nil
}

// signERC20CustodyWithdraw signs a erc20 withdrawal transaction
// function withdrawAndCall
// address to,
// address token,
// uint256 amount,
func (signer *Signer) signERC20CustodyWithdraw(txData *OutboundData) (*ethtypes.Transaction, error) {
	data, err := erc20CustodyV2ABI.Pack("withdraw", txData.to, txData.asset, txData.amount)
	if err != nil {
		return nil, errors.Wrap(err, "withdraw pack error")
	}

	tx, _, _, err := signer.Sign(
		data,
		signer.er20CustodyAddress,
		zeroValue,
		txData.gas,
		txData.nonce,
	)
	if err != nil {
		return nil, errors.Wrap(err, "sign withdraw error")
	}

	return tx, nil
}

func (signer *Signer) signZetaConnectorWithdraw(txData *OutboundData) (*ethtypes.Transaction, error) {
	data, err := connectorABI.Pack("withdraw", txData.to, txData.amount)
	if err != nil {
		return nil, errors.Wrap(err, "withdraw pack error")
	}

	tx, _, _, err := signer.Sign(
		data,
		signer.zetaConnectorAddress,
		zeroValue,
		txData.gas,
		txData.nonce,
	)
	if err != nil {
		return nil, errors.Wrap(err, "sign withdraw error")
	}

	return tx, nil
}

func (signer *Signer) signZetaConnectorWithdrawAndCall(
	ctx context.Context,
	txData *OutboundData,
) (*ethtypes.Transaction, error) {
	connectorABI, err := connector.ZetaConnectorNativeMetaData.GetAbi()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get ZetaConnectorNativeMetaData ABI")
	}

	messageContext, err := txData.MessageContext()
	if err != nil {
		return nil, err
	}

	data, err := connectorABI.Pack("withdrawAndCall", messageContext, txData.to, txData.amount, txData.message)
	if err != nil {
		return nil, errors.Wrap(err, "withdraw and call pack error")
	}

	tx, _, _, err := signer.Sign(
		ctx,
		data,
		signer.zetaConnectorAddress,
		zeroValue,
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "sign withdrawAndCall error")
	}
	return tx, nil
}

// signERC20CustodyWithdrawAndCall signs a erc20 withdrawal and call transaction
// function withdrawAndCall
// address token,
// address to,
// uint256 amount,
// bytes calldata data
func (signer *Signer) signERC20CustodyWithdrawAndCall(txData *OutboundData) (*ethtypes.Transaction, error) {
	messageContext, err := txData.MessageContext()
	if err != nil {
		return nil, err
	}

	data, err := erc20CustodyV2ABI.Pack(
		"withdrawAndCall",
		messageContext,
		txData.to,
		txData.asset,
		txData.amount,
		txData.message,
	)
	if err != nil {
		return nil, errors.Wrap(err, "withdraw pack error")
	}

	tx, _, _, err := signer.Sign(
		data,
		signer.er20CustodyAddress,
		zeroValue,
		txData.gas,
		txData.nonce,
	)
	if err != nil {
		return nil, errors.Wrap(err, "sign withdrawAndCall error")
	}

	return tx, nil
}

// signERC20CustodyWithdrawRevert signs a erc20 withdrawal revert transaction
// function withdrawAndRevert
// address token,
// address to,
// uint256 amount,
// bytes calldata data
func (signer *Signer) signERC20CustodyWithdrawRevert(
	inboundSender string,
	txData *OutboundData,
) (*ethtypes.Transaction, error) {
	data, err := erc20CustodyV2ABI.Pack(
		"withdrawAndRevert",
		txData.to,
		txData.asset,
		txData.amount,
		txData.message,
		revert.RevertContext{
			Sender:        common.HexToAddress(inboundSender),
			Asset:         txData.asset,
			Amount:        txData.amount,
			RevertMessage: txData.revertOptions.RevertMessage,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "withdraw pack error")
	}

	tx, _, _, err := signer.Sign(
		data,
		signer.er20CustodyAddress,
		zeroValue,
		txData.gas,
		txData.nonce,
	)
	if err != nil {
		return nil, errors.Wrap(err, "sign withdrawAndRevert error")
	}

	return tx, nil
}

func (signer *Signer) signZetaConnectorWithdrawRevert(
	inboundSender string,
	txData *OutboundData,
) (*ethtypes.Transaction, error) {
	data, err := connectorABI.Pack(
		"withdrawAndRevert",
		txData.to,
		txData.amount,
		txData.message,
		revert.RevertContext{
			Sender:        common.HexToAddress(inboundSender),
			Asset:         txData.asset,
			Amount:        txData.amount,
			RevertMessage: txData.revertOptions.RevertMessage,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "withdraw pack error")
	}

	tx, _, _, err := signer.Sign(
		data,
		signer.zetaConnectorAddress,
		zeroValue,
		txData.gas,
		txData.nonce,
	)
	if err != nil {
		return nil, errors.Wrap(err, "sign withdrawAndRevert error")
	}

	return tx, nil
}
