package mocks

import (
	"context"

	cc "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
)

type DummyObserver struct{}

func (ob *DummyObserver) VoteOutboundIfConfirmed(_ context.Context, _ *cc.CrossChainTx) (bool, error) {
	return false, nil
}

func (ob *DummyObserver) Start(_ context.Context)                    {}
func (ob *DummyObserver) Stop()                                      {}
func (ob *DummyObserver) SetChainParams(_ observertypes.ChainParams) {}
func (ob *DummyObserver) ChainParams() (_ observertypes.ChainParams) { return }

// ----------------------------------------------------------------------------
// EVMObserver
// ----------------------------------------------------------------------------
var _ interfaces.ChainObserver = (*EVMObserver)(nil)

// EVMObserver is a mock of evm chain observer for testing
type EVMObserver struct {
	DummyObserver
	chainParams observertypes.ChainParams
}

func NewEVMObserver(chainParams *observertypes.ChainParams) *EVMObserver {
	return &EVMObserver{chainParams: *chainParams}
}

func (ob *EVMObserver) SetChainParams(chainParams observertypes.ChainParams) {
	ob.chainParams = chainParams
}

func (ob *EVMObserver) ChainParams() observertypes.ChainParams {
	return ob.chainParams
}

// ----------------------------------------------------------------------------
// BTCObserver
// ----------------------------------------------------------------------------
var _ interfaces.ChainObserver = (*BTCObserver)(nil)

// BTCObserver is a mock of btc chain observer for testing
type BTCObserver struct {
	DummyObserver
	chainParams observertypes.ChainParams
}

func NewBTCObserver(chainParams *observertypes.ChainParams) *BTCObserver {
	return &BTCObserver{chainParams: *chainParams}
}

func (ob *BTCObserver) SetChainParams(chainParams observertypes.ChainParams) {
	ob.chainParams = chainParams
}

func (ob *BTCObserver) ChainParams() observertypes.ChainParams {
	return ob.chainParams
}

// ----------------------------------------------------------------------------
// SolanaObserver
// ----------------------------------------------------------------------------
var _ interfaces.ChainObserver = (*SolanaObserver)(nil)

// SolanaObserver is a mock of solana chain observer for testing
type SolanaObserver struct {
	DummyObserver
	chainParams observertypes.ChainParams
}

func NewSolanaObserver(chainParams *observertypes.ChainParams) *SolanaObserver {
	return &SolanaObserver{chainParams: *chainParams}
}

func (ob *SolanaObserver) SetChainParams(chainParams observertypes.ChainParams) {
	ob.chainParams = chainParams
}

func (ob *SolanaObserver) ChainParams() observertypes.ChainParams {
	return ob.chainParams
}

// ----------------------------------------------------------------------------
// SolanaObserver
// ----------------------------------------------------------------------------
var _ interfaces.ChainObserver = (*TONObserver)(nil)

// TONObserver is a mock of TON chain observer for testing
type TONObserver struct {
	DummyObserver
	chainParams observertypes.ChainParams
}

func NewTONObserver(chainParams *observertypes.ChainParams) *TONObserver {
	return &TONObserver{chainParams: *chainParams}
}

func (ob *TONObserver) SetChainParams(chainParams observertypes.ChainParams) {
	ob.chainParams = chainParams
}

func (ob *TONObserver) ChainParams() observertypes.ChainParams {
	return ob.chainParams
}
