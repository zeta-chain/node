package cctx

import (
	"fmt"

	"github.com/zeta-chain/node/cmd/zetatool/context"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func CheckInBoundTx(ctx *context.Context, cctxDetails *CCTXDetails) error {
	var (
		inboundHash    = ctx.GetInboundHash()
		zetacoreClient = ctx.GetZetaCoreClient()
		goCtx          = ctx.GetContext()
	)

	inboundHashToCCTX, err := zetacoreClient.Crosschain.InboundHashToCctx(
		goCtx, &crosschaintypes.QueryGetInboundHashToCctxRequest{
			InboundHash: inboundHash,
		})
	if err != nil {
		return fmt.Errorf("inbound chain is zetachain , cctx should be available in the same block: %w", err)
	}
	if len(inboundHashToCCTX.InboundHashToCctx.CctxIndex) == 0 {
		return fmt.Errorf("inbound hash does not have any cctx linked %s", inboundHash)
	}

	if len(inboundHashToCCTX.InboundHashToCctx.CctxIndex) > 1 {
		return fmt.Errorf("inbound hash more than one cctx %s", inboundHash)
	}

	cctxDetails.CCTXIdentifier = inboundHashToCCTX.InboundHashToCctx.CctxIndex[0]
	cctxDetails.Status = PendingOutbound
	return nil
}
