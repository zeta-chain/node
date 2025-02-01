package chains

import (
	"fmt"

	"github.com/zeta-chain/node/cmd/zetatool/cctx"
	"github.com/zeta-chain/node/cmd/zetatool/context"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func CheckInBoundTx(ctx *context.Context) (cctx.CCTXDetails, error) {
	var (
		inboundHash    = ctx.GetInboundHash()
		cctxDetails    = cctx.NewCCTXDetails()
		zetacoreClient = ctx.GetZetaCoreClient()
		goCtx          = ctx.GetContext()
	)

	inboundHashToCCTX, err := zetacoreClient.Crosschain.InboundHashToCctx(
		goCtx, &crosschaintypes.QueryGetInboundHashToCctxRequest{
			InboundHash: inboundHash,
		})
	if err != nil {
		return cctxDetails, fmt.Errorf("inbound chain is zetachain , cctx should be available in the same block: %w", err)
	}
	if len(inboundHashToCCTX.InboundHashToCctx.CctxIndex) == 0 {
		return cctxDetails, fmt.Errorf("inbound hash does not have any cctx linked %s", inboundHash)
	}

	if len(inboundHashToCCTX.InboundHashToCctx.CctxIndex) > 1 {
		return cctxDetails, fmt.Errorf("inbound hash more than one cctx %s", inboundHash)
	}

	cctxDetails.CCCTXIdentifier = inboundHashToCCTX.InboundHashToCctx.CctxIndex[0]
	cctxDetails.Status = cctx.PendingOutbound
	return cctxDetails, nil
}
