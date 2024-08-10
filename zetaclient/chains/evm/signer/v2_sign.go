package signer

import (
	"context"
	"fmt"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	erc20custodyv2 "github.com/zeta-chain/protocol-contracts/v2/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayevm.sol"
)

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
