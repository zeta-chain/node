package signer

import (
	"context"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/zetacore/pkg/constant"
	"math/big"
	"strings"
)

// SignCommandTx signs a transaction based on the given command includes:
//
//	cmd_whitelist_erc20
//	cmd_migrate_er20_custody_funds
//	cmd_migrate_tss_funds
func (signer *Signer) SignCommandTx(
	ctx context.Context,
	txData *OutboundData,
	cmd string,
	params string,
) (*ethtypes.Transaction, error) {
	switch cmd {
	case constant.CmdWhitelistERC20:
		return signer.SignWhitelistERC20Cmd(ctx, txData, params)
	case constant.CmdMigrateERC20CustodyFunds:
		return signer.SignMigrateERC20CustodyFundsCmd(ctx, txData, params)
	case constant.CmdMigrateTssFunds:
		return signer.SignMigrateTssFundsCmd(ctx, txData)
	}
	return nil, fmt.Errorf("SignCommandTx: unknown command %s", cmd)
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
		return nil, fmt.Errorf("whitelist pack error: %w", err)
	}

	tx, _, _, err := signer.Sign(
		ctx,
		data,
		txData.to,
		zeroValue,
		txData.gasLimit,
		txData.gasPrice,
		txData.outboundParams.TssNonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("sign whitelist error: %w", err)
	}
	return tx, nil
}

// SignMigrateERC20CustodyFundsCmd signs a migrate ERC20 custody funds command
func (signer *Signer) SignMigrateERC20CustodyFundsCmd(ctx context.Context, txData *OutboundData, params string) (*ethtypes.Transaction, error) {
	paramsArray := strings.Split(params, ":")
	if len(paramsArray) != 3 {
		return nil, fmt.Errorf("SignMigrateERC20CustodyFundsCmd: invalid params %s", params)
	}
	newCustody := ethcommon.HexToAddress(paramsArray[0])
	erc20 := ethcommon.HexToAddress(paramsArray[1])
	amount, ok := new(big.Int).SetString(paramsArray[2], 10)
	if !ok {
		return nil, fmt.Errorf("SignMigrateERC20CustodyFundsCmd: invalid amount %s", paramsArray[2])
	}

	custodyAbi, err := erc20custody.ERC20CustodyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	data, err := custodyAbi.Pack("withdraw", newCustody, erc20, amount)

	tx, _, _, err := signer.Sign(
		ctx,
		data,
		txData.to,
		zeroValue,
		txData.gasLimit,
		txData.gasPrice,
		txData.outboundParams.TssNonce,
		txData.height,
	)
	if err != nil {
		return nil, fmt.Errorf("SignMigrateERC20CustodyFundsCmd error: %w", err)
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
