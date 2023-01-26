package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	eth "github.com/ethereum/go-ethereum/common"
	ethcommon "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/pkg/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

func (k Keeper) DepositCoinZeta(ctx sdk.Context, to eth.Address, amount *big.Int) error {
	zetaToAddress := sdk.AccAddress(to.Bytes())
	return k.MintZetaToEVMAccount(ctx, zetaToAddress, amount)
}

func (k Keeper) DepositCoin(ctx sdk.Context, to eth.Address, amount *big.Int, senderChain string, message string, contract eth.Address, data []byte, coinType common.CoinType, asset string) (*evmtypes.MsgEthereumTxResponse, bool, error) {
	var tx *evmtypes.MsgEthereumTxResponse
	withdrawMessage := false
	var Zrc20Contract ethcommon.Address
	var coin fungibletypes.ForeignCoins
	if coinType == common.CoinType_Gas {
		var found bool
		coin, found = k.GetGasCoinForForeignCoin(ctx, senderChain)
		if !found {
			return tx, false, types.ErrGasCoinNotFound
		}
	} else {
		foreignCoinList := k.GetAllForeignCoinsForChain(ctx, senderChain)
		for _, foreignCoin := range foreignCoinList {
			if foreignCoin.Erc20ContractAddress == asset && foreignCoin.ForeignChain == senderChain {
				coin = foreignCoin
				break
			}
		}
	}
	Zrc20Contract = eth.HexToAddress(coin.Zrc20ContractAddress)
	if len(message) == 0 { // no message; transfer
		var txNoWithdraw *evmtypes.MsgEthereumTxResponse
		txNoWithdraw, err := k.DepositZRC20(ctx, Zrc20Contract, to, amount)
		if err != nil {
			return tx, false, errors.Wrap(types.ErrUnableToDepositZRC20, err.Error())
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
			return tx, false, errors.Wrap(types.ErrUnableToDepositZRC20, err.Error())
		}
		withdrawMessage = true
		tx = txWithWithdraw
	}
	return tx, withdrawMessage, nil
}
