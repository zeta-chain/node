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

func (k Keeper) SetNonceToCCTX(ctx sdk.Context, cctx types.CrossChainTx, tssPubkey string) {
	k.setNonceToCCTX(ctx, cctx, tssPubkey)
}

func (k Keeper) GetNextCctxCounter(ctx sdk.Context) uint64 {
	return k.getNextCctxCounter(ctx)
}
