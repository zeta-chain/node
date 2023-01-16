package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	eth "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/pkg/errors"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"math/big"
)

func (k Keeper) DepositCoinZeta(ctx sdk.Context, to eth.Address, amount *big.Int) error {
	zetaToAddress := sdk.AccAddress(to.Bytes())
	return k.MintZetaToEVMAccount(ctx, zetaToAddress, amount)
}

func (k Keeper) DepositCoinGas(ctx sdk.Context, to eth.Address, amount *big.Int, senderChain string, message string, contract eth.Address, data []byte) (*evmtypes.MsgEthereumTxResponse, error) {
	var tx *evmtypes.MsgEthereumTxResponse
	gasCoin, found := k.GetGasCoinForForeignCoin(ctx, senderChain)
	if !found {
		return tx, types.ErrGasCoinNotFound
	}
	Zrc20Contract := eth.HexToAddress(gasCoin.Zrc20ContractAddress)

	if len(message) == 0 { // no message; transfer
		var txNoWithdraw *evmtypes.MsgEthereumTxResponse
		txNoWithdraw, err := k.DepositZRC20(ctx, Zrc20Contract, to, amount)
		if err != nil {
			return tx, errors.Wrap(types.ErrUnableToDepositZRC20, err.Error())
		}
		tx = txNoWithdraw
	} else { // non-empty message = [contractaddress, calldata]
		var txWithWithdraw *evmtypes.MsgEthereumTxResponse
		var err error
		if len(data) == 0 {
			txWithWithdraw, err = k.DepositZRC20(ctx, Zrc20Contract, contract, amount)
		} else {
			txWithWithdraw, err = k.DepositZRC20AndCallContract(ctx, Zrc20Contract, contract, amount, data)
		}
		if err != nil {
			return tx, errors.Wrap(types.ErrUnableToDepositZRC20, err.Error())
		}
		tx = txWithWithdraw
	}
	return tx, nil
}
