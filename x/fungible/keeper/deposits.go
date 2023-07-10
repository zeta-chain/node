package keeper

import (
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

func (k Keeper) ZRC20DepositAndCallContract(ctx sdk.Context, to eth.Address, amount *big.Int, senderChain *common.Chain,
	message string, contract eth.Address, data []byte, coinType common.CoinType, asset string) (*evmtypes.MsgEthereumTxResponse, error) {
	var Zrc20Contract eth.Address
	var coin fungibletypes.ForeignCoins
	if coinType == common.CoinType_Gas {
		var found bool
		coin, found = k.GetGasCoinForForeignCoin(ctx, senderChain.ChainId)
		if !found {
			return nil, types.ErrGasCoinNotFound
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
	if len(message) == 0 { // no message; transfer
		return k.DepositZRC20(ctx, Zrc20Contract, to, amount)
	}
	// non-empty message = [contractaddress, calldata]
	if len(data) == 0 {
		return k.DepositZRC20(ctx, Zrc20Contract, contract, amount)
	}
	return k.DepositZRC20AndCallContract(ctx, Zrc20Contract, contract, amount, data)

}
