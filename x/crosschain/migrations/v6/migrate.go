package v6

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// crosschainKeeper is an interface to prevent cyclic dependency
type crosschainKeeper interface {
	GetStoreKey() storetypes.StoreKey
	GetCodec() codec.Codec
	GetAllCrossChainTx(ctx sdk.Context) []types.CrossChainTx

	SetCrossChainTx(ctx sdk.Context, cctx types.CrossChainTx)
	AddFinalizedInbound(ctx sdk.Context, inboundTxHash string, senderChainID int64, height uint64)

	SetZetaAccounting(ctx sdk.Context, accounting types.ZetaAccounting)
}

// MigrateStore migrates the x/crosschain module state from the consensus version 4 to 5
// It resets the aborted zeta amount to use the inbound tx amount instead in situations where the outbound cctx is never created.
func MigrateStore(ctx sdk.Context, crosschainKeeper crosschainKeeper, observerKeeper types.ObserverKeeper) error {

	return nil
}
