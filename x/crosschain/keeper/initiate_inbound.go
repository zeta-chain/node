package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// InitiateOutbound processes the inbound CCTX.
// It does a conditional dispatch to observer or zevm CCTX gateway based on the receiver chain
// which handle the state changes and error handling.
func (k Keeper) InitiateOutbound(ctx sdk.Context, cctx *types.CrossChainTx) error {
	chainInfo := chains.GetChainFromChainID(cctx.GetCurrentOutboundParam().ReceiverChainId)

	return k.cctxGateways[chainInfo.CctxGateway].InitiateOutbound(ctx, cctx)
}
