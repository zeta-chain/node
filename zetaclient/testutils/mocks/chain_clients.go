package mocks

import (
	"context"

	"github.com/rs/zerolog"

	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
)

// ----------------------------------------------------------------------------
// EVMObserver
// ----------------------------------------------------------------------------
var _ interfaces.ChainObserver = (*EVMObserver)(nil)

// EVMObserver is a mock of evm chain observer for testing
type EVMObserver struct {
	ChainParams observertypes.ChainParams
}

func NewEVMObserver(chainParams *observertypes.ChainParams) *EVMObserver {
	return &EVMObserver{
		ChainParams: *chainParams,
	}
}

func (ob *EVMObserver) Start(_ context.Context) {}
func (ob *EVMObserver) Stop()                   {}

func (ob *EVMObserver) IsOutboundProcessed(
	_ context.Context,
	_ *crosschaintypes.CrossChainTx,
	_ zerolog.Logger,
) (bool, bool, error) {
	return false, false, nil
}

func (ob *EVMObserver) SetChainParams(chainParams observertypes.ChainParams) {
	ob.ChainParams = chainParams
}

func (ob *EVMObserver) GetChainParams() observertypes.ChainParams {
	return ob.ChainParams
}

func (ob *EVMObserver) GetTxID(_ uint64) string {
	return ""
}

func (ob *EVMObserver) WatchInboundTracker(_ context.Context) error {
	return nil
}

// ----------------------------------------------------------------------------
// BTCObserver
// ----------------------------------------------------------------------------
var _ interfaces.ChainObserver = (*BTCObserver)(nil)

// BTCObserver is a mock of btc chain observer for testing
type BTCObserver struct {
	ChainParams observertypes.ChainParams
}

func NewBTCObserver(chainParams *observertypes.ChainParams) *BTCObserver {
	return &BTCObserver{
		ChainParams: *chainParams,
	}
}

func (ob *BTCObserver) Start(_ context.Context) {}

func (ob *BTCObserver) Stop() {}

func (ob *BTCObserver) IsOutboundProcessed(
	_ context.Context,
	_ *crosschaintypes.CrossChainTx,
	_ zerolog.Logger,
) (bool, bool, error) {
	return false, false, nil
}

func (ob *BTCObserver) SetChainParams(chainParams observertypes.ChainParams) {
	ob.ChainParams = chainParams
}

func (ob *BTCObserver) GetChainParams() observertypes.ChainParams {
	return ob.ChainParams
}

func (ob *BTCObserver) GetTxID(_ uint64) string {
	return ""
}

func (ob *BTCObserver) WatchInboundTracker(_ context.Context) error { return nil }
