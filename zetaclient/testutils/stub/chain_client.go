package stub

import (
	"github.com/rs/zerolog"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
)

// ----------------------------------------------------------------------------
// EVMClient
// ----------------------------------------------------------------------------
var _ interfaces.ChainClient = (*EVMClient)(nil)

// EVMClient is a mock of evm chain client for testing
type EVMClient struct {
	ChainParams observertypes.ChainParams
}

func NewEVMClient(chainParams *observertypes.ChainParams) *EVMClient {
	return &EVMClient{
		ChainParams: *chainParams,
	}
}

func (s *EVMClient) Start() {
}

func (s *EVMClient) Stop() {
}

func (s *EVMClient) IsOutboundProcessed(_ *crosschaintypes.CrossChainTx, _ zerolog.Logger) (bool, bool, error) {
	return false, false, nil
}

func (s *EVMClient) SetChainParams(chainParams observertypes.ChainParams) {
	s.ChainParams = chainParams
}

func (s *EVMClient) GetChainParams() observertypes.ChainParams {
	return s.ChainParams
}

func (s *EVMClient) GetTxID(_ uint64) string {
	return ""
}

func (s *EVMClient) WatchInboundTracker() {
}

// ----------------------------------------------------------------------------
// BTCClient
// ----------------------------------------------------------------------------
var _ interfaces.ChainClient = (*BTCClient)(nil)

// BTCClient is a mock of btc chain client for testing
type BTCClient struct {
	ChainParams observertypes.ChainParams
}

func NewBTCClient(chainParams *observertypes.ChainParams) *BTCClient {
	return &BTCClient{
		ChainParams: *chainParams,
	}
}

func (s *BTCClient) Start() {
}

func (s *BTCClient) Stop() {
}

func (s *BTCClient) IsOutboundProcessed(_ *crosschaintypes.CrossChainTx, _ zerolog.Logger) (bool, bool, error) {
	return false, false, nil
}

func (s *BTCClient) SetChainParams(chainParams observertypes.ChainParams) {
	s.ChainParams = chainParams
}

func (s *BTCClient) GetChainParams() observertypes.ChainParams {
	return s.ChainParams
}

func (s *BTCClient) GetTxID(_ uint64) string {
	return ""
}

func (s *BTCClient) WatchInboundTracker() {
}
