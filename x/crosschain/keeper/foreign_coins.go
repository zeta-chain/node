package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	fungibleModuleTypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

func (k Keeper) GetAllForeignCoins(ctx sdk.Context) ([]fungibleModuleTypes.ForeignCoins, error) {
	chains := k.zetaObserverKeeper.GetParams(ctx).GetSupportedChains()
	var fCoins []fungibleModuleTypes.ForeignCoins
	for _, chain := range chains {
		fCoins = append(fCoins, k.fungibleKeeper.GetAllForeignCoinsForChain(ctx, chain.ChainId)...)
	}
	return fCoins, nil
}
