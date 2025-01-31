package tracker

import (
	"fmt"

	"github.com/zeta-chain/node/cmd/zetatool/cctx"
	"github.com/zeta-chain/node/cmd/zetatool/context"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func UpdateCCTXStatus(ctx context.Context, cctxDetails *cctx.CCTXDetails) error {
	var (
		zetacoreClient = ctx.GetZetaCoreClient()
		goCtx          = ctx.GetContext()
	)
	switch cctxDetails.Status {
	case cctx.PendingInboundConfirmation:
		return nil
	default:
		cctx, err := zetacoreClient.GetCctxByHash(goCtx, cctxDetails.CCCTXIdentifier)
		if err != nil {
			return fmt.Errorf("failed to get cctx: %w", err)
		}

		cctxDetails.UpdateStatusFromZetacoreCCTX(cctx.CctxStatus.Status)
	}
	return nil
}

func UpdateHashListForPendingCCTX(ctx context.Context, cctxDetails *cctx.CCTXDetails) error {
	var (
		zetacoreClient = ctx.GetZetaCoreClient()
		goCtx          = ctx.GetContext()
	)

	if !cctxDetails.IsPendingConfirmation() {
		return nil
	}

	CCTX, err := zetacoreClient.GetCctxByHash(goCtx, cctxDetails.CCCTXIdentifier)
	if err != nil {
		return fmt.Errorf("failed to get cctx: %w", err)
	}

	outboundParams := CCTX.GetCurrentOutboundParam()

	tracker, err := zetacoreClient.GetOutboundTracker(goCtx, outboundParams.ReceiverChainId, outboundParams.TssNonce)
	if err != nil {
		switch {
		case err.Error() == "rpc error: code = NotFound desc = not found":
			return err
		case CCTX.CctxStatus.Status == crosschaintypes.CctxStatus_PendingOutbound:
			cctxDetails.Status = cctx.PendingOutboundConfirmation
		case CCTX.CctxStatus.Status == crosschaintypes.CctxStatus_PendingRevert:
			cctxDetails.Status = cctx.PendingRevertConfirmation
		}
	}
	hashList := []string{}
	for _, hash := range tracker.HashList {
		hashList = append(hashList, hash.TxHash)
	}

	cctxDetails.OutboundTrackerHashList = hashList
	return nil
}
