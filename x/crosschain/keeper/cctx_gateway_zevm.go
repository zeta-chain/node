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

/*
InitiateOutbound handles evm deposit and call ValidateOutbound.
TODO (https://github.com/zeta-chain/node/issues/2278): move remaining of this comment to ValidateOutbound once it's added.

  - If the deposit is successful, the CCTX status is changed to OutboundMined.

  - If the deposit returns an internal error i.e if HandleEVMDeposit() returns an error, but isContractReverted is false, the CCTX status is changed to Aborted.

  - If the deposit is reverted, the function tries to create a revert cctx with status PendingRevert.

  - If the creation of revert tx also fails it changes the status to Aborted.

Note : Aborted CCTXs are not refunded in this function. The refund is done using a separate refunding mechanism.
We do not return an error from this function , as all changes need to be persisted to the state.
Instead we use a temporary context to make changes and then commit the context on for the happy path ,i.e cctx is set to OutboundMined.
New CCTX status after preprocessing is returned.
*/
func (c CCTXGatewayZEVM) InitiateOutbound(ctx sdk.Context, cctx *types.CrossChainTx) (newCCTXStatus types.CctxStatus) {
	isContractReverted, err := c.crosschainKeeper.HandleEVMDeposit(ctx, cctx)

	return c.crosschainKeeper.ValidateOutboundZEVM(ctx, cctx, err, isContractReverted)
}
