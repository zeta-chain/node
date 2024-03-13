package stub

import (
	"github.com/rs/zerolog"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
)

var _ interfaces.ChainClient = (*EVMClient)(nil)

// EVMClient is a mock of evm chain client for testing
type EVMClient struct {
}

func NewEVMClient() *EVMClient {
	return &EVMClient{}
}

func (s *EVMClient) Start() {
}

func (s *EVMClient) Stop() {
}

func (s *EVMClient) IsSendOutTxProcessed(_ *crosschaintypes.CrossChainTx, _ zerolog.Logger) (bool, bool, error) {
	return false, false, nil
}

func (s *EVMClient) SetChainParams(observertypes.ChainParams) {
}

func (s *EVMClient) GetChainParams() observertypes.ChainParams {
	return observertypes.ChainParams{}
}

func (s *EVMClient) GetTxID(_ uint64) string {
	return ""
}

func (s *EVMClient) ExternalChainWatcherForNewInboundTrackerSuggestions() {
}
