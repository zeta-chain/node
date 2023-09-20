package observer

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	lastBlockObserverCount, found := k.GetLastObserverCount(ctx)
	if !found {
		ctx.Logger().Error("LastBlockObserverCount not found at height", ctx.BlockHeight())
		return
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
	// #nosec G701 always in range
	if totalObserverCountCurrentBlock == int(lastBlockObserverCount.Count) {
		return
	}
	ctx.Logger().Error("LastBlockObserverCount does not match the number of observers found at current height", ctx.BlockHeight())
	for _, observer := range allObservers {
		ctx.Logger().Error("Observes for | ", observer.ObserverChain.ChainName, ":", observer.ObserverList)
	}
	k.DisableInboundOnly(ctx)
	k.SetKeygen(ctx, types.Keygen{BlockNumber: math.MaxInt64})
	k.SetLastObserverCount(ctx, &types.LastObserverCount{Count: uint64(totalObserverCountCurrentBlock), LastChangeHeight: ctx.BlockHeight()})
}
