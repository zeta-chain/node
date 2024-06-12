package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/pkg/chains"
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
	config InitiateOutboundConfig,
) (newCCTXStatus types.CctxStatus) {
	tmpCtx, commit := ctx.CacheContext()
	outboundReceiverChainID := config.CCTX.GetCurrentOutboundParam().ReceiverChainId
	// TODO: does this condition make sense?
	noEthereumTxEvent := false
	if chains.IsZetaChain(config.CCTX.InboundParams.SenderChainId) {
		noEthereumTxEvent = true
	}

	err := func() error {
		if config.PayGas {
			err := c.crosschainKeeper.PayGasAndUpdateCctx(
				tmpCtx,
				outboundReceiverChainID,
				config.CCTX,
				config.CCTX.InboundParams.Amount,
				noEthereumTxEvent,
			)
			if err != nil {
				return err
			}
		}
		return c.crosschainKeeper.UpdateNonce(tmpCtx, outboundReceiverChainID, config.CCTX)
	}()
	if err != nil {
		// do not commit anything here as the CCTX should be aborted
		config.CCTX.SetAbort(err.Error())
		return types.CctxStatus_Aborted
	}
	commit()
	config.CCTX.SetPendingOutbound("")
	return types.CctxStatus_PendingOutbound
}
