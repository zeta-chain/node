package signer

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/erc20custody.sol"

	"github.com/zeta-chain/node/pkg/constant"
)

// SignAdminTx signs a admin cmd transaction based on the given command
func (signer *Signer) SignAdminTx(
	ctx context.Context,
	txData *OutboundData,
	cmd string,
	params string,
) (*ethtypes.Transaction, error) {
	switch cmd {
	case constant.CmdWhitelistERC20:
		return signer.signWhitelistERC20Cmd(ctx, txData, params)
	case constant.CmdMigrateERC20CustodyFunds:
		return signer.signMigrateERC20CustodyFundsCmd(ctx, txData, params)
	case constant.CmdUpdateERC20CustodyPauseStatus:
		return signer.signUpdateERC20CustodyPauseStatusCmd(ctx, txData, params)
	case constant.CmdMigrateTssFunds:
		return signer.signMigrateTssFundsCmd(ctx, txData)
	}
	return nil, fmt.Errorf("SignAdminTx: unknown command %s", cmd)
}

// signWhitelistERC20Cmd signs a whitelist command for ERC20 token
func (signer *Signer) signWhitelistERC20Cmd(
	ctx context.Context,
	txData *OutboundData,
	params string,
) (*ethtypes.Transaction, error) {
	erc20 := ethcommon.HexToAddress(params)
	if erc20 == (ethcommon.Address{}) {
		return nil, fmt.Errorf("SignAdminTx: invalid erc20 address %s", params)
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
		txData.zetaHeight,
		txData.cctxHeight,
	)
	if err != nil {
		return nil, errors.Wrap(err, "sign whitelist error")
	}
	return tx, nil
}

// signMigrateERC20CustodyFundsCmd signs a migrate ERC20 custody funds command
func (signer *Signer) signMigrateERC20CustodyFundsCmd(
	ctx context.Context,
	txData *OutboundData,
	params string,
) (*ethtypes.Transaction, error) {
	paramsArray := strings.Split(params, ",")
	if len(paramsArray) != 3 {
		return nil, fmt.Errorf("signMigrateERC20CustodyFundsCmd: invalid params %s", params)
	}
	newCustody := ethcommon.HexToAddress(paramsArray[0])
	erc20 := ethcommon.HexToAddress(paramsArray[1])
	amount, ok := new(big.Int).SetString(paramsArray[2], 10)
	if !ok {
		return nil, fmt.Errorf("signMigrateERC20CustodyFundsCmd: invalid amount %s", paramsArray[2])
	}

	custodyAbi, err := erc20custody.ERC20CustodyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	data, err := custodyAbi.Pack("withdraw", newCustody, erc20, amount)
	if err != nil {
		return nil, errors.Wrap(err, "withdraw pack error")
	}

	tx, _, _, err := signer.Sign(
		ctx,
		data,
		txData.to,
		zeroValue,
		txData.gas,
		txData.nonce,
		txData.zetaHeight,
		txData.cctxHeight,
	)
	if err != nil {
		return nil, errors.Wrap(err, "signMigrateERC20CustodyFundsCmd error")
	}
	return tx, nil
}

// signUpdateERC20CustodyPauseStatusCmd signs a update ERC20 custody pause status command
func (signer *Signer) signUpdateERC20CustodyPauseStatusCmd(
	ctx context.Context,
	txData *OutboundData,
	params string,
) (*ethtypes.Transaction, error) {
	custodyAbi, err := erc20custody.ERC20CustodyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	// select the action
	// NOTE: we could directly do Pack(params)
	// having this logic here is more explicit and restrict the possible values
	var functionName string
	switch params {
	case constant.OptionPause:
		functionName = "pause"
	case constant.OptionUnpause:
		functionName = "unpause"
	default:
		return nil, fmt.Errorf("signUpdateERC20CustodyPauseStatusCmd: invalid params %s", params)
	}

	data, err := custodyAbi.Pack(functionName)
	if err != nil {
		return nil, errors.Wrap(err, "pause/unpause pack error")
	}

	tx, _, _, err := signer.Sign(
		ctx,
		data,
		txData.to,
		zeroValue,
		txData.gas,
		txData.nonce,
		txData.zetaHeight,
		txData.cctxHeight,
	)
	if err != nil {
		return nil, errors.Wrap(err, "signUpdateERC20CustodyPauseStatusCmd error")
	}
	return tx, nil
}

// signMigrateTssFundsCmd signs a migrate TSS funds command
func (signer *Signer) signMigrateTssFundsCmd(ctx context.Context, txData *OutboundData) (*ethtypes.Transaction, error) {
	tx, _, _, err := signer.Sign(
		ctx,
		nil,
		txData.to,
		txData.amount,
		txData.gas,
		txData.nonce,
		txData.zetaHeight,
		txData.cctxHeight,
	)
	if err != nil {
		return nil, errors.Wrap(err, "signMigrateTssFundsCmd error")
	}
	return tx, nil
}
