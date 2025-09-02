package v5

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/fungible/types"
)

type fungibleKeeper interface {
	GetSystemContract(ctx sdk.Context) (types.SystemContract, bool)
	SetSystemContract(ctx sdk.Context, sc types.SystemContract)
}

// MigrateStore migrates the store from consensus version 4 to 5
// It sets the default value for GatewayGasLimit to the previously hardcoded value DefaultGatewayGasLimit
func MigrateStore(ctx sdk.Context, fungibleKeeper fungibleKeeper) error {
	system := types.SystemContract{}
	systemContract, found := fungibleKeeper.GetSystemContract(ctx)
	if found {
		system = systemContract
	}
	system.GatewayGasLimit = sdkmath.NewIntFromBigInt(types.GatewayGasLimit)
	fungibleKeeper.SetSystemContract(ctx, system)
	return nil
}
