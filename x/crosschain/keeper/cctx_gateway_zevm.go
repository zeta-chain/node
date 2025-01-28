package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	cctxerror "github.com/zeta-chain/node/pkg/errors"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// CCTXGatewayZEVM is implementation of CCTXGateway interface for ZEVM
type CCTXGatewayZEVM struct {
	crosschainKeeper Keeper
}

// NewCCTXGatewayZEVM returns new instance of CCTXGatewayZEVM
func NewCCTXGatewayZEVM(crosschainKeeper Keeper) CCTXGatewayZEVM {
	return CCTXGatewayZEVM{
		crosschainKeeper: crosschainKeeper,
	}
}

// InitiateOutbound handles evm deposit and immediately validates pending outbound
func (c CCTXGatewayZEVM) InitiateOutbound(
	ctx sdk.Context,
	config InitiateOutboundConfig,
) (newCCTXStatus types.CctxStatus, err error) {
	// abort if CCTX inbound observation status indicates failure
	// this is specifically for Bitcoin inbound error 'INSUFFICIENT_DEPOSITOR_FEE'
	if config.CCTX.InboundParams.Status == types.InboundStatus_INSUFFICIENT_DEPOSITOR_FEE {

		config.CCTX.SetAbort(types.StatusMessages{StatusMessage: fmt.Sprintf("observation failed, reason %s", types.InboundStatus_INSUFFICIENT_DEPOSITOR_FEE.String())})
		return types.CctxStatus_Aborted, nil
	}

	// process the deposit
	tmpCtx, commit := ctx.CacheContext()
	isContractReverted, err := c.crosschainKeeper.HandleEVMDeposit(tmpCtx, config.CCTX)

	if err != nil && !isContractReverted {
		// exceptional case; internal error; should abort CCTX
		config.CCTX.SetAbort(types.StatusMessages{
			StatusMessage:        "outbound failed but the universal contract did not revert",
			ErrorMessageOutbound: cctxerror.NewCCTXErrorJSONMessage("failed to deposit tokens in ZEVM", err),
		})
		return types.CctxStatus_Aborted, err
	}

	newCCTXStatus = c.crosschainKeeper.ValidateOutboundZEVM(ctx, config.CCTX, err, isContractReverted)
	if newCCTXStatus == types.CctxStatus_OutboundMined || newCCTXStatus == types.CctxStatus_PendingRevert {
		commit()
	}

	return newCCTXStatus, nil
}
