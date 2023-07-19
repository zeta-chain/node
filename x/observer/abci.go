package observer

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	lastBlockObserverCount, found := k.GetLastBlockObserverCount(ctx)
	if !found {
		ctx.Logger().Error("LastBlockObserverCount not found at height", ctx.BlockHeight())
	}

	allObservers := k.GetAllObserverMappers(ctx)
	totalObserverCount := 0
	for _, observer := range allObservers {
		totalObserverCount += len(observer.ObserverList)
	}
	if len(allObservers) != int(lastBlockObserverCount.Count) {
		ctx.Logger().Error("LastBlockObserverCount does not match the number of observers found at current height", ctx.BlockHeight())
		k.SetPermissionFlags(ctx, types.PermissionFlags{IsInboundEnabled: false})
	}

}
