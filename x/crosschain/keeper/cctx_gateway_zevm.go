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
func (c CCTXGatewayZEVM) InitiateOutbound(ctx sdk.Context, cctx *types.CrossChainTx) (newCCTXStatus types.CctxStatus) {
	tmpCtx, _ := ctx.CacheContext()
	isContractReverted, err := c.crosschainKeeper.HandleEVMDeposit(tmpCtx, cctx)

	if err != nil && !isContractReverted {
		// exceptional case; internal error; should abort CCTX
		cctx.SetAbort(err.Error())
		return types.CctxStatus_Aborted
	}

	cctx.SetPendingOutbound("")
	return c.crosschainKeeper.ValidateOutboundZEVM(ctx, cctx, err, isContractReverted)
}
