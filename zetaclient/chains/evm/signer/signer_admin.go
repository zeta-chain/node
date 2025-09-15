package signer

import (
	"context"
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/zeta-chain/protocol-contracts/pkg/erc20custody.sol"

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
	case constant.CmdMigrateTSSFunds:
		return signer.signMigrateTSSFundsCmd(ctx, txData)
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
		txData.height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "sign whitelist error")
	}
	return tx, nil
}

// signMigrateTSSFundsCmd signs a migrate TSS funds command
func (signer *Signer) signMigrateTSSFundsCmd(ctx context.Context, txData *OutboundData) (*ethtypes.Transaction, error) {
	tx, _, _, err := signer.Sign(
		ctx,
		nil,
		txData.to,
		txData.amount,
		txData.gas,
		txData.nonce,
		txData.height,
	)
	if err != nil {
		return nil, errors.Wrap(err, "signMigrateTSSFundsCmd error")
	}
	return tx, nil
}
