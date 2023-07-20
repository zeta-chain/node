package observer

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"math"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	lastBlockObserverCount, found := k.GetLastBlockObserverCount(ctx)
	if !found {
		ctx.Logger().Error("LastBlockObserverCount not found at height", ctx.BlockHeight())
	}

	allObservers := k.GetAllObserverMappers(ctx)
	totalObserverCountCurrentBlock := 0
	for _, observer := range allObservers {
		totalObserverCountCurrentBlock += len(observer.ObserverList)
	}
	if totalObserverCountCurrentBlock < 0 {
		ctx.Logger().Error("TotalObserverCount is negative at height", ctx.BlockHeight())
		return
	}
	if totalObserverCountCurrentBlock != int(lastBlockObserverCount.Count) {
		ctx.Logger().Error("LastBlockObserverCount does not match the number of observers found at current height", ctx.BlockHeight())
		k.SetPermissionFlags(ctx, types.PermissionFlags{IsInboundEnabled: false})
		k.SetKeygen(ctx, types.Keygen{BlockNumber: math.MaxInt64})
	}

	k.SetLastBlockObserverCount(ctx, &types.LastBlockObserverCount{Count: uint64(totalObserverCountCurrentBlock)})

}
