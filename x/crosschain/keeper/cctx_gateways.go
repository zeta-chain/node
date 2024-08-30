package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// CCTXGateway is interface implemented by every gateway. It is one of interfaces used for communication
// between CCTX gateways and crosschain module, and it is called by crosschain module.
type CCTXGateway interface {
	// Initiate a new outbound, this tells the CCTXGateway to carry out the action to execute the outbound.
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
