package keeper

import (
	"github.com/pkg/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	eth "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

func (k Keeper) DepositCoinZeta(ctx sdk.Context, to eth.Address, amount *big.Int) error {
	zetaToAddress := sdk.AccAddress(to.Bytes())
	return k.MintZetaToEVMAccount(ctx, zetaToAddress, amount)
}

func (k Keeper) ZRC20DepositAndCallContract(ctx sdk.Context, to eth.Address, amount *big.Int, senderChain *common.Chain, message string, contract eth.Address, data []byte, coinType common.CoinType, asset string) (*evmtypes.MsgEthereumTxResponse, bool, error) {
	var evmTxResponse *evmtypes.MsgEthereumTxResponse
	withdrawMessage := false
	var Zrc20Contract eth.Address
	var coin fungibletypes.ForeignCoins
	if coinType == common.CoinType_Gas {
		var found bool
		coin, found = k.GetGasCoinForForeignCoin(ctx, senderChain.ChainId)
		if !found {
			return nil, false, types.ErrGasCoinNotFound
		}
	} else {
		foreignCoinList := k.GetAllForeignCoinsForChain(ctx, senderChain.ChainId)
		found := false
		for _, foreignCoin := range foreignCoinList {
			if foreignCoin.Asset == asset && foreignCoin.ForeignChainId == senderChain.ChainId {
				coin = foreignCoin
				found = true
				break
			}
		}
		if !found {
			return nil, false, types.ErrForeignCoinNotFound
		}
	}
	Zrc20Contract = eth.HexToAddress(coin.Zrc20ContractAddress)
	if len(message) == 0 { // no message; transfer
		var txNoWithdraw *evmtypes.MsgEthereumTxResponse
		txNoWithdraw, err := k.DepositZRC20(ctx, Zrc20Contract, to, amount)
		if err != nil {
			return txNoWithdraw, false, errors.Wrap(types.ErrUnableToDepositZRC20, err.Error())
		}
		evmTxResponse = txNoWithdraw
	} else { // non-empty message = [contractaddress, calldata]
		var txWithWithdraw *evmtypes.MsgEthereumTxResponse
		var err error
		if len(data) == 0 {
			txWithWithdraw, err = k.DepositZRC20(ctx, Zrc20Contract, contract, amount)
		} else {
			txWithWithdraw, err = k.DepositZRC20AndCallContract(ctx, Zrc20Contract, contract, amount, data)
		}
		if err != nil {
			return txWithWithdraw, false, errors.Wrap(types.ErrUnableToDepositZRC20, err.Error())
		}
		withdrawMessage = true
		evmTxResponse = txWithWithdraw
	}
	return evmTxResponse, withdrawMessage, nil

}
