package ballot

import (
	"fmt"

	"github.com/zeta-chain/node/cmd/zetatool/cctx"
	"github.com/zeta-chain/node/cmd/zetatool/chains"

	"github.com/zeta-chain/node/cmd/zetatool/context"
)

func GetBallotIdentifier(ctx *context.Context) (cctx.CCTXDetails, error) {
	var (
		ballotIdentifierMessage = cctx.CCTXDetails{}
		inboundChain            = ctx.GetInboundChain()
		err                     error
	)

	switch {
	case inboundChain.IsZetaChain():
		{
			ballotIdentifierMessage, err = chains.CheckInBoundTx(ctx)
			if err != nil {
				return ballotIdentifierMessage, fmt.Errorf(
					"failed to get inbound ballot for zeta chain %d, %w",
					inboundChain.ChainId,
					err,
				)
			}
		}

	case inboundChain.IsEVMChain():
		{
			ballotIdentifierMessage, err = chains.EvmInboundBallotIdentifier(ctx)
			if err != nil {
				return ballotIdentifierMessage, fmt.Errorf(
					"failed to get inbound ballot for evm chain %d, %w",
					inboundChain.ChainId,
					err,
				)
			}
		}
	case inboundChain.IsBitcoinChain():
		{
			ballotIdentifierMessage, err = btcInboundBallotIdentifier(ctx)
			if err != nil {
				return ballotIdentifierMessage, fmt.Errorf(
					"failed to get inbound ballot for bitcoin chain %d, %w",
					inboundChain.ChainId,
					err,
				)
			}
		}
	case inboundChain.IsSolanaChain():
		{
			ballotIdentifierMessage, err = solanaInboundBallotIdentifier(ctx)
			if err != nil {
				return ballotIdentifierMessage, fmt.Errorf(
					"failed to get inbound ballot for solana chain %d, %w",
					inboundChain.ChainId,
					err,
				)
			}
		}
	default:
		ballotIdentifierMessage.Message = "Chain not supported"
	}

	return ballotIdentifierMessage, nil
}
