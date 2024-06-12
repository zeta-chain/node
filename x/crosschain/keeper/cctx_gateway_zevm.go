package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
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
) (newCCTXStatus types.CctxStatus) {
	tmpCtx, commit := ctx.CacheContext()
	isContractReverted, err := c.crosschainKeeper.HandleEVMDeposit(tmpCtx, config.CCTX)

	if err != nil && !isContractReverted {
		// exceptional case; internal error; should abort CCTX
		config.CCTX.SetAbort(err.Error())
		return types.CctxStatus_Aborted
	}

	config.CCTX.SetPendingOutbound("")
	newCCTXStatus = c.crosschainKeeper.ValidateOutboundZEVM(ctx, config.CCTX, err, isContractReverted)
	if newCCTXStatus == types.CctxStatus_OutboundMined {
		commit()
	}

	return newCCTXStatus
}
