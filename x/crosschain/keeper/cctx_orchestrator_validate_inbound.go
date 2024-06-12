package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func (k Keeper) ValidateInboundObservers(
	ctx sdk.Context,
	msg *types.MsgVoteInbound,
	payGas bool,
) (*types.CrossChainTx, error) {
	tss, tssFound := k.zetaObserverKeeper.GetTSS(ctx)
	if !tssFound {
		return nil, types.ErrCannotFindTSSKeys
	}
	// create a new CCTX from the inbound message.The status of the new CCTX is set to PendingInbound.
	cctx, err := types.NewCCTX(ctx, *msg, tss.TssPubkey)
	if err != nil {
		return nil, err
	}
	// Initiate outbound, the process function manages the state commit and cctx status change.
	// If the process fails, the changes to the evm state are rolled back.
	_, err = k.InitiateOutbound(ctx, InitiateOutboundConfig{
		CCTX:   &cctx,
		PayGas: payGas,
	})
	if err != nil {
		return nil, err
	}

	return &cctx, nil
}
