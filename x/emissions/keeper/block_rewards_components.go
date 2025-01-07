package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkmath "cosmossdk.io/math"
	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/x/emissions/types"
)

func (k Keeper) GetReservesFactor(ctx sdk.Context) sdkmath.LegacyDec {
	reserveAmount := k.GetBankKeeper().GetBalance(ctx, types.EmissionsModuleAddress, config.BaseDenom)
	return sdkmath.LegacyNewDecFromInt(reserveAmount.Amount)
}
