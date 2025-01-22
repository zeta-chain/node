package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

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
	// abort if CCTX already contains an initial error message from inbound vote msg
	if strings.Contains(config.CCTX.CctxStatus.ErrorMessage, types.InboundStatus_insufficient_depositor_fee.String()) {
		config.CCTX.SetAbort("observation failed", "")
		return types.CctxStatus_Aborted, nil
	}

	// process the deposit
	tmpCtx, commit := ctx.CacheContext()
	isContractReverted, err := c.crosschainKeeper.HandleEVMDeposit(tmpCtx, config.CCTX)

	if err != nil && !isContractReverted {
		// exceptional case; internal error; should abort CCTX
		config.CCTX.SetAbort(
			"error during deposit that is not smart contract revert",
			err.Error())
		return types.CctxStatus_Aborted, err
	}

	newCCTXStatus = c.crosschainKeeper.ValidateOutboundZEVM(ctx, config.CCTX, err, isContractReverted)
	if newCCTXStatus == types.CctxStatus_OutboundMined || newCCTXStatus == types.CctxStatus_PendingRevert {
		commit()
	}

	return newCCTXStatus, nil
}
