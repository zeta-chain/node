package v8

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

type observerKeeper interface {
	GetParamsIfExists(ctx sdk.Context) (params types.Params)
	SetParams(ctx sdk.Context, params types.Params)
}

// MigrateStore performs in-place store migrations from v6 to v7
func MigrateStore(ctx sdk.Context, observerKeeper observerKeeper) error {
	ctx.Logger().Info("Migrating observer store from v6 to v7")
	params := observerKeeper.GetParamsIfExists(ctx)
	for _, ob := range params.ObserverParams {
		chain := chains.GetChainFromChainID(ob.Chain.ChainId)
		ob.Chain = chain
	}
	observerKeeper.SetParams(ctx, params)
	return nil
}
