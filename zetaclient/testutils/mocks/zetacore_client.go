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
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

const ErrMsgPaused = "zetacore client is paused"
const ErrMsgRPCFailed = "rpc failed"

var _ interfaces.ZetacoreClient = &MockZetaCoreClient{}

type MockZetaCoreClient struct {
	paused    bool
	zetaChain chains.Chain

	// the mock data for testing
	// pending cctxs
	pendingCctxs map[int64][]*crosschaintypes.CrossChainTx

	// rate limiter flags
	rateLimiterFlags *crosschaintypes.RateLimiterFlags

	// rate limiter input
	input *crosschaintypes.QueryRateLimiterInputResponse
}

func NewMockZetaCoreClient() *MockZetaCoreClient {
	zetaChain, err := chains.ZetaChainFromChainID("zetachain_7000-1")
	if err != nil {
		panic(err)
	}
	return &MockZetaCoreClient{
		paused:       false,
		zetaChain:    zetaChain,
		pendingCctxs: map[int64][]*crosschaintypes.CrossChainTx{},
	}
}

func (z *MockZetaCoreClient) PostVoteInbound(_, _ uint64, _ *crosschaintypes.MsgVoteOnObservedInboundTx) (string, string, error) {
	if z.paused {
		return "", "", errors.New(ErrMsgPaused)
	}
	return "", "", nil
}

func (z *MockZetaCoreClient) PostVoteOutbound(_ string, _ string, _ uint64, _ uint64, _ *big.Int, _ uint64, _ *big.Int, _ chains.ReceiveStatus, _ chains.Chain, _ uint64, _ coin.CoinType) (string, string, error) {
	if z.paused {
		return "", "", errors.New(ErrMsgPaused)
	}
	return sample.Hash().Hex(), "", nil
}

func (z *MockZetaCoreClient) PostGasPrice(_ chains.Chain, _ uint64, _ string, _ uint64) (string, error) {
	if z.paused {
		return "", errors.New(ErrMsgPaused)
	}
	return "", nil
}

func (z *MockZetaCoreClient) PostVoteBlockHeader(_ int64, _ []byte, _ int64, _ proofs.HeaderData) (string, error) {
	if z.paused {
		return "", errors.New(ErrMsgPaused)
	}
	return "", nil
}

func (z *MockZetaCoreClient) GetBlockHeaderChainState(_ int64) (lightclienttypes.QueryGetChainStateResponse, error) {
	if z.paused {
		return lightclienttypes.QueryGetChainStateResponse{}, errors.New(ErrMsgPaused)
	}
	return lightclienttypes.QueryGetChainStateResponse{}, nil
}

func (z *MockZetaCoreClient) PostBlameData(_ *blame.Blame, _ int64, _ string) (string, error) {
	if z.paused {
		return "", errors.New(ErrMsgPaused)
	}
	return "", nil
}

func (z *MockZetaCoreClient) AddTxHashToOutTxTracker(_ int64, _ uint64, _ string, _ *proofs.Proof, _ string, _ int64) (string, error) {
	if z.paused {
		return "", errors.New(ErrMsgPaused)
	}
	return "", nil
}

func (z *MockZetaCoreClient) Chain() chains.Chain {
	return z.zetaChain
}

func (z *MockZetaCoreClient) GetLogger() *zerolog.Logger {
	return nil
}

func (z *MockZetaCoreClient) GetKeys() *keys.Keys {
	return &keys.Keys{}
}

func (z *MockZetaCoreClient) GetKeyGen() (*observerTypes.Keygen, error) {
	if z.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return &observerTypes.Keygen{}, nil
}

func (z *MockZetaCoreClient) GetBlockHeight() (int64, error) {
	if z.paused {
		return 0, errors.New(ErrMsgPaused)
	}
	return 0, nil
}

func (z *MockZetaCoreClient) GetLastBlockHeightByChain(_ chains.Chain) (*crosschaintypes.LastBlockHeight, error) {
	if z.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return &crosschaintypes.LastBlockHeight{}, nil
}

func (z *MockZetaCoreClient) GetRateLimiterInput(_ int64) (crosschaintypes.QueryRateLimiterInputResponse, error) {
	if z.paused {
		return crosschaintypes.QueryRateLimiterInputResponse{}, errors.New(ErrMsgPaused)
	}
	if z.input == nil {
		return crosschaintypes.QueryRateLimiterInputResponse{}, errors.New(ErrMsgRPCFailed)
	}
	return *z.input, nil
}

func (z *MockZetaCoreClient) ListPendingCctx(chainID int64) ([]*crosschaintypes.CrossChainTx, uint64, error) {
	if z.paused {
		return nil, 0, errors.New(ErrMsgPaused)
	}
	return z.pendingCctxs[chainID], 0, nil
}

func (z *MockZetaCoreClient) ListPendingCctxWithinRatelimit() ([]*crosschaintypes.CrossChainTx, uint64, int64, string, bool, error) {
	if z.paused {
		return nil, 0, 0, "", false, errors.New(ErrMsgPaused)
	}
	return []*crosschaintypes.CrossChainTx{}, 0, 0, "", false, nil
}

func (z *MockZetaCoreClient) GetPendingNoncesByChain(_ int64) (observerTypes.PendingNonces, error) {
	if z.paused {
		return observerTypes.PendingNonces{}, errors.New(ErrMsgPaused)
	}
	return observerTypes.PendingNonces{}, nil
}

func (z *MockZetaCoreClient) GetCctxByNonce(_ int64, _ uint64) (*crosschaintypes.CrossChainTx, error) {
	if z.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return &crosschaintypes.CrossChainTx{}, nil
}

func (z *MockZetaCoreClient) GetOutTxTracker(_ chains.Chain, _ uint64) (*crosschaintypes.OutTxTracker, error) {
	if z.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return &crosschaintypes.OutTxTracker{}, nil
}

func (z *MockZetaCoreClient) GetAllOutTxTrackerByChain(_ int64, _ interfaces.Order) ([]crosschaintypes.OutTxTracker, error) {
	if z.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return []crosschaintypes.OutTxTracker{}, nil
}

func (z *MockZetaCoreClient) GetCrosschainFlags() (observerTypes.CrosschainFlags, error) {
	if z.paused {
		return observerTypes.CrosschainFlags{}, errors.New(ErrMsgPaused)
	}
	return observerTypes.CrosschainFlags{}, nil
}

func (z *MockZetaCoreClient) GetRateLimiterFlags() (crosschaintypes.RateLimiterFlags, error) {
	if z.paused {
		return crosschaintypes.RateLimiterFlags{}, errors.New(ErrMsgPaused)
	}
	if z.rateLimiterFlags == nil {
		return crosschaintypes.RateLimiterFlags{}, errors.New(ErrMsgRPCFailed)
	}
	return *z.rateLimiterFlags, nil
}

func (z *MockZetaCoreClient) GetObserverList() ([]string, error) {
	if z.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return []string{}, nil
}

func (z *MockZetaCoreClient) GetBtcTssAddress(_ int64) (string, error) {
	if z.paused {
		return "", errors.New(ErrMsgPaused)
	}
	return testutils.TSSAddressBTCMainnet, nil
}

func (z *MockZetaCoreClient) GetInboundTrackersForChain(_ int64) ([]crosschaintypes.InTxTracker, error) {
	if z.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return []crosschaintypes.InTxTracker{}, nil
}

func (z *MockZetaCoreClient) Pause() {
	z.paused = true
}

func (z *MockZetaCoreClient) Unpause() {
	z.paused = false
}

func (z *MockZetaCoreClient) GetZetaHotKeyBalance() (math.Int, error) {
	if z.paused {
		return math.NewInt(0), errors.New(ErrMsgPaused)
	}
	return math.NewInt(0), nil
}

// ----------------------------------------------------------------------------
// Feed data to the mock zetacore client for testing
// ----------------------------------------------------------------------------

func (z *MockZetaCoreClient) WithPendingCctx(chainID int64, cctxs []*crosschaintypes.CrossChainTx) *MockZetaCoreClient {
	z.pendingCctxs[chainID] = cctxs
	return z
}

func (z *MockZetaCoreClient) WithRateLimiterFlags(flags *crosschaintypes.RateLimiterFlags) *MockZetaCoreClient {
	z.rateLimiterFlags = flags
	return z
}

func (z *MockZetaCoreClient) WithRateLimiterInput(input *crosschaintypes.QueryRateLimiterInputResponse) *MockZetaCoreClient {
	z.input = input
	return z
}
