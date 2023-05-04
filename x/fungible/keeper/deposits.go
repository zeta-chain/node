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

func (k Keeper) DepositCoin(ctx sdk.Context, to eth.Address, amount *big.Int, senderChain *common.Chain, message string, contract eth.Address, data []byte, coinType common.CoinType, asset string) (*evmtypes.MsgEthereumTxResponse, error) {
	var tx *evmtypes.MsgEthereumTxResponse
	var Zrc20Contract eth.Address
	var coin fungibletypes.ForeignCoins
	if coinType == common.CoinType_Gas {
		var found bool
		coin, found = k.GetGasCoinForForeignCoin(ctx, senderChain.ChainId)
		if !found {
			return tx, types.ErrGasCoinNotFound
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
			return nil, types.ErrForeignCoinNotFound
		}
	}
	Zrc20Contract = eth.HexToAddress(coin.Zrc20ContractAddress)
	var err error
	if len(message) == 0 { // no message; transfer
		tx, err = k.DepositZRC20(ctx, Zrc20Contract, to, amount)
		if err != nil {
			return tx, errors.Wrap(types.ErrUnableToDepositZRC20, err.Error())
		}
	} else { // non-empty message = [contractaddress, calldata]
		if len(data) == 0 {
			tx, err = k.DepositZRC20(ctx, Zrc20Contract, contract, amount)
		} else {
			tx, err = k.DepositZRC20AndCallContract(ctx, Zrc20Contract, contract, amount, data)
		}
		if err != nil {
			return tx, errors.Wrap(types.ErrUnableToDepositZRC20, err.Error())
		}
	}
	return tx, nil

}
