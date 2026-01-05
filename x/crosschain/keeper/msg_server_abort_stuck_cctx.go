package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/crosschain/types"
)

const (
	// AbortMessage is the message to abort a stuck CCTX
	AbortMessage = "cctx aborted through MsgAbortStuckCCTX"
)

// AbortStuckCCTX aborts a stuck CCTX
// Authorized: admin policy group 2
func (k msgServer) AbortStuckCCTX(
	goCtx context.Context,
	msg *types.MsgAbortStuckCCTX,
) (*types.MsgAbortStuckCCTXResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if authorized
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, errors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// check if the cctx exists
	cctx, found := k.GetCrossChainTx(ctx, msg.CctxIndex)
	if !found {
		return nil, types.ErrCannotFindCctx
	}

	// check if the cctx is pending
	if !cctx.CctxStatus.Status.IsPending() {
		return nil, types.ErrStatusNotPending
	}

	// update the status
	k.ProcessAbort(ctx, &cctx, types.StatusMessages{
		StatusMessage: AbortMessage,
	})

	// Set all outbound params to Executed to ensure nonces are removed from pending nonces
	for _, outboundParams := range cctx.OutboundParams {
		outboundParams.TxFinalizationStatus = types.TxFinalizationStatus_Executed
	}

	// Save out outbound,
	// We do not need to provide the tss-pubkey as NonceToCctx is not updated / New outbound is not added
	k.SaveOutbound(ctx, &cctx, "")

	return &types.MsgAbortStuckCCTXResponse{}, nil
}
