package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// InitiateOutbound processes the inbound CCTX.
// It does a conditional dispatch to correct CCTX gateway based on the receiver chain
// which handle the state changes and error handling.
func (k Keeper) InitiateOutbound(ctx sdk.Context, cctx *types.CrossChainTx) error {
	receiverChainID := cctx.GetCurrentOutboundParam().ReceiverChainId
	chainInfo := chains.GetChainFromChainID(receiverChainID)
	if chainInfo == nil {
		return fmt.Errorf("chain info not found for %d", receiverChainID)
	}

	cctxGateway, ok := k.cctxGateways[chainInfo.CctxGateway]
	if !ok {
		return fmt.Errorf("CCTXGateway not defined for receiver chain %d", receiverChainID)
	}
	return cctxGateway.InitiateOutbound(ctx, cctx)
}
