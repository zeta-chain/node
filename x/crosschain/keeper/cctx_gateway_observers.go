package keeper

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/chains"
	cctxerror "github.com/zeta-chain/node/pkg/errors"
	"github.com/zeta-chain/node/x/crosschain/types"
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
) (newCCTXStatus types.CctxStatus, err error) {
	tmpCtx, commit := ctx.CacheContext()
	outboundReceiverChainID := config.CCTX.GetCurrentOutboundParam().ReceiverChainId
	// TODO (https://github.com/zeta-chain/node/issues/1010): workaround for this bug
	noEthereumTxEvent := chains.IsZetaChain(
		config.CCTX.InboundParams.SenderChainId,
		c.crosschainKeeper.GetAuthorityKeeper().GetAdditionalChainList(ctx),
	)

	err = func() error {
		// If ShouldPayGas flag is set during ValidateInbound PayGasAndUpdateCctx should be called
		// which will set GasPrice and Amount. Otherwise, use median gas price and InboundParams amount.
		if config.ShouldPayGas {
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
		} else {
			gasPrice, priorityFee, found := c.crosschainKeeper.GetMedianGasValues(ctx, config.CCTX.GetCurrentOutboundParam().ReceiverChainId)
			if !found {
				return fmt.Errorf("gasprice not found for %d", config.CCTX.GetCurrentOutboundParam().ReceiverChainId)
			}
			gasLimit := sdkmath.NewUint(config.CCTX.GetCurrentOutboundParam().CallOptions.GasLimit)
			config.CCTX.GetCurrentOutboundParam().GasPrice = gasPrice.String()
			config.CCTX.GetCurrentOutboundParam().GasPriorityFee = priorityFee.String()
			config.CCTX.GetCurrentOutboundParam().Amount = config.CCTX.InboundParams.Amount

			// Calculate the gas fee paid by the user and add this to the CCTX
			// This value does not include the protocol fee paid by the user,
			// as that should not be included in any calculations for stability pool or refund
			gasFeePaidByUser := gasLimit.Mul(gasPrice)
			config.CCTX.GetCurrentOutboundParam().UserGasFeePaid = gasFeePaidByUser
		}
		return c.crosschainKeeper.SetObserverOutboundInfo(tmpCtx, outboundReceiverChainID, config.CCTX)
	}()
	if err != nil {
		// do not commit anything here as the CCTX should be aborted, and use ctx for processing the abort
		c.crosschainKeeper.ProcessAbort(ctx, config.CCTX, types.StatusMessages{
			StatusMessage:        "outbound failed unable to process",
			ErrorMessageOutbound: cctxerror.NewCCTXErrorJSONMessage("Unable to create outbound", err),
		})
		return types.CctxStatus_Aborted, err
	}
	commit()
	return types.CctxStatus_PendingOutbound, nil
}
