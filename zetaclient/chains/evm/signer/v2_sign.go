package signer

import (
	"context"
	"fmt"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	erc20custodyv2 "github.com/zeta-chain/protocol-contracts/v2/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayevm.sol"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/revert.sol"
)

// signGatewayExecute signs a gateway execute
// used for gas withdrawal and call transaction
// function execute
// address destination,
// bytes calldata data
func (signer *Signer) signGatewayExecute(ctx context.Context, txData *OutboundData) (*ethtypes.Transaction, error) {
	gatewayABI, err := gatewayevm.GatewayEVMMetaData.GetAbi()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get GatewayEVMMetaData ABI")
	}

	data, err := gatewayABI.Pack("execute", txData.to, txData.message)
	if err != nil {
		return nil, fmt.Errorf("execute pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(
		ctx,
		data,
		signer.gatewayAddress,
		txData.amount,
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("sign execute error: %w", err)
	}

	return tx, nil
}

// signGatewayExecuteRevert signs a gateway execute revert
// function executeRevert
// address destination,
// bytes calldata data
func (signer *Signer) signGatewayExecuteRevert(
	ctx context.Context,
	txData *OutboundData,
) (*ethtypes.Transaction, error) {
	gatewayABI, err := gatewayevm.GatewayEVMMetaData.GetAbi()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get GatewayEVMMetaData ABI")
	}

	data, err := gatewayABI.Pack(
		"executeRevert",
		txData.to,
		txData.message,
		revert.RevertContext{
			Asset:         txData.asset,
			Amount:        txData.amount.Uint64(),
			RevertMessage: txData.revertOptions.RevertMessage,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("executeRevert pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(
		ctx,
		data,
		signer.gatewayAddress,
		txData.amount,
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("sign executeRevert error: %w", err)
	}

	return tx, nil
}

// signERC20CustodyWithdraw signs a erc20 withdrawal transaction
// function withdrawAndCall
// address to,
// address token,
// uint256 amount,
func (signer *Signer) signERC20CustodyWithdraw(
	ctx context.Context,
	txData *OutboundData,
) (*ethtypes.Transaction, error) {
	erc20CustodyV2ABI, err := erc20custodyv2.ERC20CustodyMetaData.GetAbi()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get ERC20CustodyMetaData ABI")
	}

	data, err := erc20CustodyV2ABI.Pack("withdraw", txData.to, txData.asset, txData.amount)
	if err != nil {
		return nil, fmt.Errorf("withdraw pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(
		ctx,
		data,
		signer.er20CustodyAddress,
		zeroValue,
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("sign withdraw error: %w", err)
	}

	return tx, nil
}

// signERC20CustodyWithdrawAndCall signs a erc20 withdrawal and call transaction
// function withdrawAndCall
// address token,
// address to,
// uint256 amount,
// bytes calldata data
func (signer *Signer) signERC20CustodyWithdrawAndCall(
	ctx context.Context,
	txData *OutboundData,
) (*ethtypes.Transaction, error) {
	erc20CustodyV2ABI, err := erc20custodyv2.ERC20CustodyMetaData.GetAbi()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get ERC20CustodyMetaData ABI")
	}

	data, err := erc20CustodyV2ABI.Pack("withdrawAndCall", txData.to, txData.asset, txData.amount, txData.message)
	if err != nil {
		return nil, fmt.Errorf("withdraw pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(
		ctx,
		data,
		signer.er20CustodyAddress,
		zeroValue,
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("sign withdrawAndCall error: %w", err)
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
	ctx context.Context,
	txData *OutboundData,
) (*ethtypes.Transaction, error) {
	erc20CustodyV2ABI, err := erc20custodyv2.ERC20CustodyMetaData.GetAbi()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get ERC20CustodyMetaData ABI")
	}

	data, err := erc20CustodyV2ABI.Pack(
		"withdrawAndRevert",
		txData.to,
		txData.asset,
		txData.amount,
		txData.message,
		revert.RevertContext{
			Asset:         txData.asset,
			Amount:        txData.amount.Uint64(),
			RevertMessage: txData.revertOptions.RevertMessage,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("withdraw pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(
		ctx,
		data,
		signer.er20CustodyAddress,
		zeroValue,
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("sign withdrawAndRevert error: %w", err)
	}

	return tx, nil
}
