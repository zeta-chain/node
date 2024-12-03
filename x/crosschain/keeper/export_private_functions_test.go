package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// These functions are exported for testing purposes

func (k Keeper) UpdateZetaAccounting(ctx sdk.Context, cctx types.CrossChainTx) {
	k.updateZetaAccounting(ctx, cctx)
}

func (k Keeper) UpdateInboundHashToCCTX(ctx sdk.Context, cctx types.CrossChainTx) {
	k.updateInboundHashToCCTX(ctx, cctx)
}

func (k Keeper) SetNonceToCCTXMapping(ctx sdk.Context, cctx types.CrossChainTx, tssPubkey string) {
	k.setNonceToCCTXMapping(ctx, cctx, tssPubkey)
}
