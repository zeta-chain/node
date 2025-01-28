package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/x/emissions/types"
)

func (k Keeper) GetReservesFactor(ctx sdk.Context) sdkmath.LegacyDec {
	reserveAmount := k.GetBankKeeper().GetBalance(ctx, types.EmissionsModuleAddress, config.BaseDenom)
	return sdkmath.LegacyNewDecFromInt(reserveAmount.Amount)
}
