package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func RemoveObserverEvent(ctx sdk.Context, mapper types.ObserverMapper, observerAddress, removalReason string) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.ObserverRemoved,
			sdk.NewAttribute(types.ObserverAddress, observerAddress),
			sdk.NewAttribute(types.RemovalReason, removalReason),
			sdk.NewAttribute(types.ObservevationChain, mapper.ObserverChain.String()),
			sdk.NewAttribute(types.ObservervationType, mapper.ObservationType.String()),
			sdk.NewAttribute(types.ObserverList, PrettyPrintList(mapper.ObserverList)),
		),
	)
}

func AddObserverEvent(ctx sdk.Context, observerAddress string, chain types.ObserverChain, observationType types.ObservationType) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.ObserverAdded,
			sdk.NewAttribute(types.ObserverAddress, observerAddress),
			sdk.NewAttribute(types.ObservevationChain, chain.String()),
			sdk.NewAttribute(types.ObservervationType, observationType.String()),
		),
	)
}

func PrettyPrintList[st any](list []st) string {
	output := ""
	for _, s := range list {
		stringSt := fmt.Sprintf("%v", s)
		output = output + stringSt + ","
	}
	if len(output) <= 0 {
		return output
	}
	output = output[:len(output)-1]
	return output
}
