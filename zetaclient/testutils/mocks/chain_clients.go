package mocks

import (
	"context"

	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
)

// ----------------------------------------------------------------------------
// EVMObserver
// ----------------------------------------------------------------------------
var _ interfaces.ChainObserver = (*EVMObserver)(nil)

// EVMObserver is a mock of evm chain observer for testing
type EVMObserver struct {
	chainParams observertypes.ChainParams
}

func NewEVMObserver(chainParams *observertypes.ChainParams) *EVMObserver {
	return &EVMObserver{
		chainParams: *chainParams,
	}
}

func (ob *EVMObserver) Start(_ context.Context) {}
func (ob *EVMObserver) Stop()                   {}

func (ob *EVMObserver) VoteOutboundIfConfirmed(
	_ context.Context,
	_ *crosschaintypes.CrossChainTx,
) (bool, error) {
	return false, nil
}

func (ob *EVMObserver) SetChainParams(chainParams observertypes.ChainParams) {
	ob.chainParams = chainParams
}

func (ob *EVMObserver) ChainParams() observertypes.ChainParams {
	return ob.chainParams
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
	chainParams observertypes.ChainParams
}

func NewBTCObserver(chainParams *observertypes.ChainParams) *BTCObserver {
	return &BTCObserver{
		chainParams: *chainParams,
	}
}

func (ob *BTCObserver) Start(_ context.Context) {}

func (ob *BTCObserver) Stop() {}

func (ob *BTCObserver) VoteOutboundIfConfirmed(
	_ context.Context,
	_ *crosschaintypes.CrossChainTx,
) (bool, error) {
	return false, nil
}

func (ob *BTCObserver) SetChainParams(chainParams observertypes.ChainParams) {
	ob.chainParams = chainParams
}

func (ob *BTCObserver) ChainParams() observertypes.ChainParams {
	return ob.chainParams
}

func (ob *BTCObserver) GetTxID(_ uint64) string {
	return ""
}

func (ob *BTCObserver) WatchInboundTracker(_ context.Context) error { return nil }

// ----------------------------------------------------------------------------
// SolanaObserver
// ----------------------------------------------------------------------------
var _ interfaces.ChainObserver = (*SolanaObserver)(nil)

// SolanaObserver is a mock of solana chain observer for testing
type SolanaObserver struct {
	chainParams observertypes.ChainParams
}

func NewSolanaObserver(chainParams *observertypes.ChainParams) *SolanaObserver {
	return &SolanaObserver{
		chainParams: *chainParams,
	}
}

func (ob *SolanaObserver) Start(_ context.Context) {}

func (ob *SolanaObserver) Stop() {}

func (ob *SolanaObserver) VoteOutboundIfConfirmed(
	_ context.Context,
	_ *crosschaintypes.CrossChainTx,
) (bool, error) {
	return false, nil
}

func (ob *SolanaObserver) SetChainParams(chainParams observertypes.ChainParams) {
	ob.chainParams = chainParams
}

func (ob *SolanaObserver) ChainParams() observertypes.ChainParams {
	return ob.chainParams
}

func (ob *SolanaObserver) GetTxID(_ uint64) string {
	return ""
}

func (ob *SolanaObserver) WatchInboundTracker(_ context.Context) error { return nil }
