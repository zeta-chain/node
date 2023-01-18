package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	fungibleModuleTypes "github.com/zeta-chain/zetacore/x/fungible/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (k Keeper) GetAllForeignCoins(ctx sdk.Context) ([]fungibleModuleTypes.ForeignCoins, error) {
	chains, found := k.zetaObserverKeeper.GetSupportedChains(ctx)
	if !found {
		return nil, zetaObserverTypes.ErrSupportedChains
	}
	var fCoins []fungibleModuleTypes.ForeignCoins
	for _, chain := range chains.ChainList {
		fCoins = append(fCoins, k.fungibleKeeper.GetAllForeignCoinsForChain(ctx, chain.ChainName.String())...)
	}
	return fCoins, nil
}
