package v4

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

type ObserverKeeper interface {
	GetCrosschainFlags(ctx sdk.Context) (val types.CrosschainFlags, found bool)
	SetCrosschainFlags(ctx sdk.Context, crosschainFlags types.CrosschainFlags)
}

func MigrateStore(ctx sdk.Context, k ObserverKeeper) error {
	k.SetCrosschainFlags(ctx, *types.DefaultCrosschainFlags())
	return nil
}
