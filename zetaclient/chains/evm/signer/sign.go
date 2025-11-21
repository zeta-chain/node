package signer

import (
	"context"
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/erc20custody.sol"
	connectorevm "github.com/zeta-chain/protocol-contracts-evm/pkg/zetaconnector.base.sol"
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
		txData.gas,
		txData.nonce,
	)
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
		txData.gas,
		txData.nonce,
	)
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
		signer.TSS().PubKey().AddressEVM(),
		zeroValue, // zero out the amount to cancel the tx
		txData.gas,
		txData.nonce,
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
		txData.gas,
		txData.nonce,
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
		return nil, errors.Wrap(err, "withdraw pack error")
	}

	tx, _, _, err := signer.Sign(
		ctx,
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

// SignWhitelistERC20Cmd signs a whitelist command for ERC20 token
func (signer *Signer) SignWhitelistERC20Cmd(
	ctx context.Context,
	txData *OutboundData,
	params string,
) (*ethtypes.Transaction, error) {
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
		return nil, errors.Wrap(err, "whitelist pack error")
	}
	tx, _, _, err := signer.Sign(
		ctx,
		data,
		txData.to,
		zeroValue,
		txData.gas,
		txData.nonce,
	)
	if err != nil {
		return nil, errors.Wrap(err, "sign whitelist error")
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
		txData.gas,
		txData.nonce,
	)
	if err != nil {
		return nil, errors.Wrap(err, "SignMigrateTssFundsCmd error")
	}
	return tx, nil
}
