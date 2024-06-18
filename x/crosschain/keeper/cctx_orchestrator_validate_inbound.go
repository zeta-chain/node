package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// ValidateInbound is the only entry-point to create new CCTX (eg. when observers voting is done or new inbound event is detected).
// It creates new CCTX object and calls InitiateOutbound method.
func (k Keeper) ValidateInbound(
	ctx sdk.Context,
	msg *types.MsgVoteInbound,
	shouldPayGas bool,
) (*types.CrossChainTx, error) {
	tss, tssFound := k.zetaObserverKeeper.GetTSS(ctx)
	if !tssFound {
		return nil, types.ErrCannotFindTSSKeys
	}

	// Do not process if inbound is disabled
	if !k.zetaObserverKeeper.IsInboundEnabled(ctx) {
		return nil, observertypes.ErrInboundDisabled
	}

	// create a new CCTX from the inbound message. The status of the new CCTX is set to PendingInbound.
	cctx, err := types.NewCCTX(ctx, *msg, tss.TssPubkey)
	if err != nil {
		return nil, err
	}

	// Initiate outbound, the process function manages the state commit and cctx status change.
	// If the process fails, the changes to the evm state are rolled back.
	_, err = k.InitiateOutbound(ctx, InitiateOutboundConfig{
		CCTX:         &cctx,
		ShouldPayGas: shouldPayGas,
	})
	if err != nil {
		return nil, err
	}

	inCctxIndex, ok := ctx.Value("inCctxIndex").(string)
	if ok {
		cctx.InboundParams.ObservedHash = inCctxIndex
	}
	k.SetCctxAndNonceToCctxAndInboundHashToCctx(ctx, cctx)

	return &cctx, nil
}
