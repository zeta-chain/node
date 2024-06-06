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

// InitiateOutbound immediately calls ValidateOutbound because cctx can be validated without votes
func (c CCTXGatewayZEVM) InitiateOutbound(ctx sdk.Context, cctx *types.CrossChainTx) (newCCTXStatus types.CctxStatus) {
	return c.crosschainKeeper.ValidateOutboundZEVM(ctx, cctx)
}
