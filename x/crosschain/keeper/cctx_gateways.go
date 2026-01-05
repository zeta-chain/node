package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// CCTXGateway is interface implemented by every gateway. It is one of interfaces used for communication
// between CCTX gateways and crosschain module, and it is called by crosschain module.

// Note on gas payements:
// - CCTXGateway_observers : This is the gateway used for all connected chains.The outbound processing needs gas fee payment for execution on the destination chain.
// - CCTXGateway_zevm : This is the gateway used only for ZEVM.The outbound processing does not need gas fee payment for execution on zevm.
type CCTXGateway interface {
	// InitiateOutbound initiates a new outbound, this tells the CCTXGateway to carry out the action to execute the outbound.
	// It is the only entry point to initiate an outbound and it returns new CCTX status after it is completed.
	InitiateOutbound(ctx sdk.Context, config InitiateOutboundConfig) (newCCTXStatus types.CctxStatus, err error)
}

var cctxGateways map[chains.CCTXGateway]CCTXGateway

// ResolveCCTXGateway respolves cctx gateway implementation based on provided cctx gateway
func ResolveCCTXGateway(c chains.CCTXGateway, keeper Keeper) (CCTXGateway, bool) {
	cctxGateways = map[chains.CCTXGateway]CCTXGateway{
		chains.CCTXGateway_observers: NewCCTXGatewayObservers(keeper),
		chains.CCTXGateway_zevm:      NewCCTXGatewayZEVM(keeper),
	}

	cctxGateway, ok := cctxGateways[c]
	return cctxGateway, ok
}
