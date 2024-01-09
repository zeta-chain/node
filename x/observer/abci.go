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

	allObservers, found := k.GetObserverSet(ctx)
	if !found {
		ctx.Logger().Error("ObserverSet not found at height", ctx.BlockHeight())
		return
	}
	totalObserverCountCurrentBlock := allObservers.LenUint()
	if totalObserverCountCurrentBlock < 0 {
		ctx.Logger().Error("TotalObserverCount is negative at height", ctx.BlockHeight())
		return
	}
	// #nosec G701 always in range
	if totalObserverCountCurrentBlock == lastBlockObserverCount.Count {
		return
	}
	ctx.Logger().Error("LastBlockObserverCount does not match the number of observers found at current height", ctx.BlockHeight())
	for _, observer := range allObservers.ObserverList {
		ctx.Logger().Error("Observer :  ", observer)
	}
	// #nosec G701 always in range

	k.DisableInboundOnly(ctx)
	k.SetKeygen(ctx, types.Keygen{BlockNumber: math.MaxInt64})
	// #nosec G701 always positive
	k.SetLastObserverCount(ctx, &types.LastObserverCount{Count: totalObserverCountCurrentBlock, LastChangeHeight: ctx.BlockHeight()})
}
