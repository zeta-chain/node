package keeper

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// TODO (https://github.com/zeta-chain/node/issues/2345): this is just a tmp solution because some flows require gas payment and others don't.
// TBD during implementation of issue above if info can be passed to CCTX constructor somehow.
// and not initialize CCTX using MsgVoteInbound and instead use something like (InboundParams, OutboundParams).
// Also check if msg.Digest can be replaced to calculate index
type InitiateOutboundConfig struct {
	CCTX         *types.CrossChainTx
	ShouldPayGas bool
}

// InitiateOutbound initiates the outbound for the CCTX depending on the CCTX gateway.
// It does a conditional dispatch to correct CCTX gateway based on the receiver chain
// which handles the state changes and error handling.
func (k Keeper) InitiateOutbound(ctx sdk.Context, config InitiateOutboundConfig) (types.CctxStatus, error) {
	receiverChainID := config.CCTX.GetCurrentOutboundParam().ReceiverChainId
	chainInfo, found := chains.GetChainFromChainID(receiverChainID, k.GetAuthorityKeeper().GetAdditionalChainList(ctx))
	if !found {
		return config.CCTX.CctxStatus.Status, cosmoserrors.Wrapf(
			types.ErrInitiatitingOutbound,
			"chain info not found for %d", receiverChainID,
		)
	}

	cctxGateway, found := ResolveCCTXGateway(chainInfo.CctxGateway, k)
	if !found {
		return config.CCTX.CctxStatus.Status, cosmoserrors.Wrapf(
			types.ErrInitiatitingOutbound,
			"CCTXGateway not defined for receiver chain %d", receiverChainID,
		)
	}

	config.CCTX.SetPendingOutbound(types.StatusMessages{StatusMessage: "initiating outbound"})
	return cctxGateway.InitiateOutbound(ctx, config)
}
