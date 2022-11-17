package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	fmt.Println(params.BallotThresholds)
	k.paramstore.GetParamSet(ctx, &params)
	return
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	fmt.Println("Setting Params", params.String())
	k.paramstore.SetParamSet(ctx, &params)
}
