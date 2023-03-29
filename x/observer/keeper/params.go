package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

func (k Keeper) UpdateClientParamsForChain(ctx sdk.Context, chain common.Chain, newClientParams *types.ClientParams) error {
	params := k.GetParams(ctx)
	for _, p := range params.ObserverParams {
		if p.Chain.IsEqual(chain) {
			p.ClientParams = newClientParams
		}
	}
	k.SetParams(ctx, params)
	return nil
}
