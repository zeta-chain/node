package v5

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// crosschainKeeper is an interface to prevent cyclic dependency
type crosschainKeeper interface {
	GetStoreKey() storetypes.StoreKey
	GetCodec() codec.Codec
	GetAllCrossChainTx(ctx sdk.Context) []types.CrossChainTx
	AddFinalizedInbound(ctx sdk.Context, inboundTxHash string, senderChainID int64, height uint64)

	SetZetaAccounting(ctx sdk.Context, accounting types.ZetaAccounting)
}

// MigrateStore migrates the x/crosschain module state from the consensus version 4 to 5
// It resets the aborted zeta amount to use the inbound tx amount instead in situations where the outbound cctx is never created.
func MigrateStore(
	ctx sdk.Context,
	crosschainKeeper crosschainKeeper,
) error {

	ccctxList := crosschainKeeper.GetAllCrossChainTx(ctx)
	abortedAmountZeta := sdkmath.ZeroUint()
	for _, cctx := range ccctxList {
		if cctx.CctxStatus.Status == types.CctxStatus_Aborted && cctx.GetCurrentOutTxParam().CoinType == common.CoinType_Zeta {
			abortedValue := GetAbortedAmount(cctx)
			abortedAmountZeta = abortedAmountZeta.Add(abortedValue)
		}
	}
	crosschainKeeper.SetZetaAccounting(ctx, types.ZetaAccounting{AbortedZetaAmount: abortedAmountZeta})

	return nil
}

func GetAbortedAmount(cctx types.CrossChainTx) sdkmath.Uint {
	if cctx.OutboundTxParams != nil && !cctx.GetCurrentOutTxParam().Amount.IsZero() {
		return cctx.GetCurrentOutTxParam().Amount
	}
	if cctx.InboundTxParams != nil {
		return cctx.InboundTxParams.Amount
	}

	return sdkmath.ZeroUint()
}
