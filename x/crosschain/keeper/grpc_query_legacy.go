package keeper

import (
	"context"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// OutTxTrackerAll
// Deprecated(v17): use OutboundTrackerAll
func (k Keeper) OutTxTrackerAll(c context.Context, req *types.QueryAllOutboundTrackerRequest) (*types.QueryAllOutboundTrackerResponse, error) {
	return k.OutboundTrackerAll(c, req)
}

// OutTxTrackerAllByChain
// Deprecated(v17): use OutboundTrackerAllByChain
func (k Keeper) OutTxTrackerAllByChain(c context.Context, req *types.QueryAllOutboundTrackerByChainRequest) (*types.QueryAllOutboundTrackerByChainResponse, error) {
	return k.OutboundTrackerAllByChain(c, req)
}

// OutTxTracker
// Deprecated(v17): use OutboundTracker
func (k Keeper) OutTxTracker(c context.Context, req *types.QueryGetOutboundTrackerRequest) (*types.QueryGetOutboundTrackerResponse, error) {
	return k.OutboundTracker(c, req)
}

// InTxTrackerAllByChain
// Deprecated(v17): use InboundTrackerAllByChain
func (k Keeper) InTxTrackerAllByChain(c context.Context, req *types.QueryAllInboundTrackerByChainRequest) (*types.QueryAllInboundTrackerByChainResponse, error) {
	return k.InboundTrackerAllByChain(c, req)
}

// InTxTrackerAll
// Deprecated(v17): use InboundTrackerAll
func (k Keeper) InTxTrackerAll(c context.Context, req *types.QueryAllInboundTrackersRequest) (*types.QueryAllInboundTrackersResponse, error) {
	return k.InboundTrackerAll(c, req)
}

// InTxHashToCctxAll
// Deprecated(v17): use InboundHashToCctxAll
func (k Keeper) InTxHashToCctxAll(c context.Context, req *types.QueryAllInboundHashToCctxRequest) (*types.QueryAllInboundHashToCctxResponse, error) {
	return k.InboundHashToCctxAll(c, req)
}

// InTxHashToCctx
// Deprecated(v17): use InboundHashToCctx
func (k Keeper) InTxHashToCctx(c context.Context, req *types.QueryGetInboundHashToCctxRequest) (*types.QueryGetInboundHashToCctxResponse, error) {
	return k.InboundHashToCctx(c, req)
}

// InTxHashToCctxData
// Deprecated(v17): use InboundHashToCctxData
func (k Keeper) InTxHashToCctxData(c context.Context, req *types.QueryInboundHashToCctxDataRequest) (*types.QueryInboundHashToCctxDataResponse, error) {
	return k.InboundHashToCctxData(c, req)
}
