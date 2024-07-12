package mocks

import (
	"errors"
	"math/big"

	"cosmossdk.io/math"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/go-tss/blame"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	"github.com/zeta-chain/zetacore/testutil/sample"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
	chaininterfaces "github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	keyinterfaces "github.com/zeta-chain/zetacore/zetaclient/keys/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

const ErrMsgPaused = "zetacore client is paused"
const ErrMsgRPCFailed = "rpc failed"

var _ chaininterfaces.ZetacoreClient = &MockZetacoreClient{}

type MockZetacoreClient struct {
	paused    bool
	zetaChain chains.Chain

	// the mock observer keys
	keys keyinterfaces.ObserverKeys

	// the mock data for testing
	// pending cctxs
	pendingCctxs map[int64][]*crosschaintypes.CrossChainTx

	// rate limiter flags
	rateLimiterFlags *crosschaintypes.RateLimiterFlags

	// rate limiter input
	input *crosschaintypes.QueryRateLimiterInputResponse
}

func NewMockZetacoreClient() *MockZetacoreClient {
	return &MockZetacoreClient{
		paused:       false,
		zetaChain:    chains.ZetaChainMainnet,
		pendingCctxs: map[int64][]*crosschaintypes.CrossChainTx{},
	}
}

func (m *MockZetacoreClient) PostVoteInbound(_, _ uint64, _ *crosschaintypes.MsgVoteInbound) (string, string, error) {
	if m.paused {
		return "", "", errors.New(ErrMsgPaused)
	}
	return "", "", nil
}

func (m *MockZetacoreClient) PostVoteOutbound(
	_ string,
	_ string,
	_ uint64,
	_ uint64,
	_ *big.Int,
	_ uint64,
	_ *big.Int,
	_ chains.ReceiveStatus,
	_ chains.Chain,
	_ uint64,
	_ coin.CoinType,
) (string, string, error) {
	if m.paused {
		return "", "", errors.New(ErrMsgPaused)
	}
	return sample.Hash().Hex(), "", nil
}

func (m *MockZetacoreClient) PostGasPrice(_ chains.Chain, _, _, _ uint64) (string, error) {
	if m.paused {
		return "", errors.New(ErrMsgPaused)
	}
	return "", nil
}

func (m *MockZetacoreClient) PostVoteBlockHeader(_ int64, _ []byte, _ int64, _ proofs.HeaderData) (string, error) {
	if m.paused {
		return "", errors.New(ErrMsgPaused)
	}
	return "", nil
}

func (m *MockZetacoreClient) GetBlockHeaderChainState(_ int64) (lightclienttypes.QueryGetChainStateResponse, error) {
	if m.paused {
		return lightclienttypes.QueryGetChainStateResponse{}, errors.New(ErrMsgPaused)
	}
	return lightclienttypes.QueryGetChainStateResponse{}, nil
}

func (m *MockZetacoreClient) PostBlameData(_ *blame.Blame, _ int64, _ string) (string, error) {
	if m.paused {
		return "", errors.New(ErrMsgPaused)
	}
	return "", nil
}

func (m *MockZetacoreClient) AddOutboundTracker(
	_ int64,
	_ uint64,
	_ string,
	_ *proofs.Proof,
	_ string,
	_ int64,
) (string, error) {
	if m.paused {
		return "", errors.New(ErrMsgPaused)
	}
	return "", nil
}

func (m *MockZetacoreClient) Chain() chains.Chain {
	return m.zetaChain
}

func (m *MockZetacoreClient) GetLogger() *zerolog.Logger {
	return nil
}

func (m *MockZetacoreClient) GetKeys() keyinterfaces.ObserverKeys {
	return m.keys
}

func (m *MockZetacoreClient) GetKeyGen() (*observerTypes.Keygen, error) {
	if m.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return &observerTypes.Keygen{}, nil
}

func (m *MockZetacoreClient) GetBlockHeight() (int64, error) {
	if m.paused {
		return 0, errors.New(ErrMsgPaused)
	}
	return 0, nil
}

func (m *MockZetacoreClient) GetRateLimiterInput(_ int64) (crosschaintypes.QueryRateLimiterInputResponse, error) {
	if m.paused {
		return crosschaintypes.QueryRateLimiterInputResponse{}, errors.New(ErrMsgPaused)
	}
	if m.input == nil {
		return crosschaintypes.QueryRateLimiterInputResponse{}, errors.New(ErrMsgRPCFailed)
	}
	return *m.input, nil
}

func (m *MockZetacoreClient) ListPendingCctx(chainID int64) ([]*crosschaintypes.CrossChainTx, uint64, error) {
	if m.paused {
		return nil, 0, errors.New(ErrMsgPaused)
	}
	return m.pendingCctxs[chainID], 0, nil
}

func (m *MockZetacoreClient) ListPendingCctxWithinRatelimit() ([]*crosschaintypes.CrossChainTx, uint64, int64, string, bool, error) {
	if m.paused {
		return nil, 0, 0, "", false, errors.New(ErrMsgPaused)
	}
	return []*crosschaintypes.CrossChainTx{}, 0, 0, "", false, nil
}

func (m *MockZetacoreClient) GetPendingNoncesByChain(_ int64) (observerTypes.PendingNonces, error) {
	if m.paused {
		return observerTypes.PendingNonces{}, errors.New(ErrMsgPaused)
	}
	return observerTypes.PendingNonces{}, nil
}

func (m *MockZetacoreClient) GetCctxByNonce(_ int64, _ uint64) (*crosschaintypes.CrossChainTx, error) {
	if m.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return &crosschaintypes.CrossChainTx{}, nil
}

func (m *MockZetacoreClient) GetOutboundTracker(_ chains.Chain, _ uint64) (*crosschaintypes.OutboundTracker, error) {
	if m.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return &crosschaintypes.OutboundTracker{}, nil
}

func (m *MockZetacoreClient) GetAllOutboundTrackerByChain(
	_ int64,
	_ chaininterfaces.Order,
) ([]crosschaintypes.OutboundTracker, error) {
	if m.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return []crosschaintypes.OutboundTracker{}, nil
}

func (m *MockZetacoreClient) GetCrosschainFlags() (observerTypes.CrosschainFlags, error) {
	if m.paused {
		return observerTypes.CrosschainFlags{}, errors.New(ErrMsgPaused)
	}
	return observerTypes.CrosschainFlags{}, nil
}

func (m *MockZetacoreClient) GetRateLimiterFlags() (crosschaintypes.RateLimiterFlags, error) {
	if m.paused {
		return crosschaintypes.RateLimiterFlags{}, errors.New(ErrMsgPaused)
	}
	if m.rateLimiterFlags == nil {
		return crosschaintypes.RateLimiterFlags{}, errors.New(ErrMsgRPCFailed)
	}
	return *m.rateLimiterFlags, nil
}

func (m *MockZetacoreClient) GetObserverList() ([]string, error) {
	if m.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return []string{}, nil
}

func (m *MockZetacoreClient) GetBtcTssAddress(_ int64) (string, error) {
	if m.paused {
		return "", errors.New(ErrMsgPaused)
	}
	return testutils.TSSAddressBTCMainnet, nil
}

func (m *MockZetacoreClient) GetInboundTrackersForChain(_ int64) ([]crosschaintypes.InboundTracker, error) {
	if m.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return []crosschaintypes.InboundTracker{}, nil
}

func (m *MockZetacoreClient) Pause() {
	m.paused = true
}

func (m *MockZetacoreClient) Unpause() {
	m.paused = false
}

func (m *MockZetacoreClient) GetZetaHotKeyBalance() (math.Int, error) {
	if m.paused {
		return math.NewInt(0), errors.New(ErrMsgPaused)
	}
	return math.NewInt(0), nil
}

// ----------------------------------------------------------------------------
// Feed data to the mock zetacore client for testing
// ----------------------------------------------------------------------------

func (m *MockZetacoreClient) WithKeys(keys keyinterfaces.ObserverKeys) *MockZetacoreClient {
	m.keys = keys
	return m
}

func (m *MockZetacoreClient) WithPendingCctx(chainID int64, cctxs []*crosschaintypes.CrossChainTx) *MockZetacoreClient {
	m.pendingCctxs[chainID] = cctxs
	return m
}

func (m *MockZetacoreClient) WithRateLimiterFlags(flags *crosschaintypes.RateLimiterFlags) *MockZetacoreClient {
	m.rateLimiterFlags = flags
	return m
}

func (m *MockZetacoreClient) WithRateLimiterInput(
	input *crosschaintypes.QueryRateLimiterInputResponse,
) *MockZetacoreClient {
	m.input = input
	return m
}
