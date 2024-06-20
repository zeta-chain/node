package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// CCTXGatewayObservers is implementation of CCTXGateway interface for observers
type CCTXGatewayObservers struct {
	crosschainKeeper Keeper
}

// NewCCTXGatewayObservers returns new instance of CCTXGatewayObservers
func NewCCTXGatewayObservers(crosschainKeeper Keeper) CCTXGatewayObservers {
	return CCTXGatewayObservers{
		crosschainKeeper: crosschainKeeper,
	}
}

/*
InitiateOutbound updates the store so observers can use the PendingCCTX query:

  - If preprocessing of outbound is successful, the CCTX status is changed to PendingOutbound.

  - if preprocessing of outbound, such as paying the gas fee for the destination fails, the state is reverted to aborted

    We do not return an error from this function, as all changes need to be persisted to the state.

    Instead, we use a temporary context to make changes and then commit the context on for the happy path, i.e cctx is set to PendingOutbound.
    New CCTX status after preprocessing is returned.
*/
func (c CCTXGatewayObservers) InitiateOutbound(
	ctx sdk.Context,
	cctx *types.CrossChainTx,
) (newCCTXStatus types.CctxStatus) {
	tmpCtx, commit := ctx.CacheContext()
	outboundReceiverChainID := cctx.GetCurrentOutboundParam().ReceiverChainId
	err := func() error {
		err := c.crosschainKeeper.PayGasAndUpdateCctx(
			tmpCtx,
			outboundReceiverChainID,
			cctx,
			cctx.InboundParams.Amount,
			false,
		)
		if err != nil {
			return err
		}
		return c.crosschainKeeper.UpdateNonce(tmpCtx, outboundReceiverChainID, cctx)
	}()
	if err != nil {
		// do not commit anything here as the CCTX should be aborted
		cctx.SetAbort(err.Error())
		return types.CctxStatus_Aborted
	}
	commit()
	return types.CctxStatus_PendingOutbound
}
