package stub

import (
	"errors"
	"math/big"

	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"

	"cosmossdk.io/math"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/go-tss/blame"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	"github.com/zeta-chain/zetacore/testutil/sample"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

const ErrMsgPaused = "zeta core bridge is paused"
const ErrMsgRPCFailed = "rpc failed"

var _ interfaces.ZetaCoreBridger = &MockZetaCoreBridge{}

type MockZetaCoreBridge struct {
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

func NewMockZetaCoreBridge() *MockZetaCoreBridge {
	zetaChain, err := chains.ZetaChainFromChainID("zetachain_7000-1")
	if err != nil {
		panic(err)
	}
	return &MockZetaCoreBridge{
		paused:       false,
		zetaChain:    zetaChain,
		pendingCctxs: map[int64][]*crosschaintypes.CrossChainTx{},
	}
}

func (z *MockZetaCoreBridge) PostVoteInbound(_, _ uint64, _ *crosschaintypes.MsgVoteInbound) (string, string, error) {
	if z.paused {
		return "", "", errors.New(ErrMsgPaused)
	}
	return "", "", nil
}

func (z *MockZetaCoreBridge) PostVoteOutbound(_ string, _ string, _ uint64, _ uint64, _ *big.Int, _ uint64, _ *big.Int, _ chains.ReceiveStatus, _ chains.Chain, _ uint64, _ coin.CoinType) (string, string, error) {
	if z.paused {
		return "", "", errors.New(ErrMsgPaused)
	}
	return sample.Hash().Hex(), "", nil
}

func (z *MockZetaCoreBridge) PostGasPrice(_ chains.Chain, _ uint64, _ string, _ uint64) (string, error) {
	if z.paused {
		return "", errors.New(ErrMsgPaused)
	}
	return "", nil
}

func (z *MockZetaCoreBridge) PostVoteBlockHeader(_ int64, _ []byte, _ int64, _ proofs.HeaderData) (string, error) {
	if z.paused {
		return "", errors.New(ErrMsgPaused)
	}
	return "", nil
}

func (z *MockZetaCoreBridge) GetBlockHeaderChainState(_ int64) (lightclienttypes.QueryGetChainStateResponse, error) {
	if z.paused {
		return lightclienttypes.QueryGetChainStateResponse{}, errors.New(ErrMsgPaused)
	}
	return lightclienttypes.QueryGetChainStateResponse{}, nil
}

func (z *MockZetaCoreBridge) PostBlameData(_ *blame.Blame, _ int64, _ string) (string, error) {
	if z.paused {
		return "", errors.New(ErrMsgPaused)
	}
	return "", nil
}

func (z *MockZetaCoreBridge) AddTxHashToOutboundTracker(_ int64, _ uint64, _ string, _ *proofs.Proof, _ string, _ int64) (string, error) {
	if z.paused {
		return "", errors.New(ErrMsgPaused)
	}
	return "", nil
}

func (z *MockZetaCoreBridge) GetKeys() *keys.Keys {
	return &keys.Keys{}
}

func (z *MockZetaCoreBridge) GetBlockHeight() (int64, error) {
	if z.paused {
		return 0, errors.New(ErrMsgPaused)
	}
	return 0, nil
}

func (z *MockZetaCoreBridge) GetZetaBlockHeight() (int64, error) {
	if z.paused {
		return 0, errors.New(ErrMsgPaused)
	}
	return 0, nil
}

func (z *MockZetaCoreBridge) GetLastBlockHeightByChain(_ chains.Chain) (*crosschaintypes.LastBlockHeight, error) {
	if z.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return &crosschaintypes.LastBlockHeight{}, nil
}

func (z *MockZetaCoreBridge) GetRateLimiterInput(_ int64) (crosschaintypes.QueryRateLimiterInputResponse, error) {
	if z.paused {
		return crosschaintypes.QueryRateLimiterInputResponse{}, errors.New(ErrMsgPaused)
	}
	if z.input == nil {
		return crosschaintypes.QueryRateLimiterInputResponse{}, errors.New(ErrMsgRPCFailed)
	}
	return *z.input, nil
}

func (z *MockZetaCoreBridge) ListPendingCctx(chainID int64) ([]*crosschaintypes.CrossChainTx, uint64, error) {
	if z.paused {
		return nil, 0, errors.New(ErrMsgPaused)
	}
	return z.pendingCctxs[chainID], 0, nil
}

func (z *MockZetaCoreBridge) ListPendingCctxWithinRatelimit() ([]*crosschaintypes.CrossChainTx, uint64, int64, string, bool, error) {
	if z.paused {
		return nil, 0, 0, "", false, errors.New(ErrMsgPaused)
	}
	return []*crosschaintypes.CrossChainTx{}, 0, 0, "", false, nil
}

func (z *MockZetaCoreBridge) GetPendingNoncesByChain(_ int64) (observerTypes.PendingNonces, error) {
	if z.paused {
		return observerTypes.PendingNonces{}, errors.New(ErrMsgPaused)
	}
	return observerTypes.PendingNonces{}, nil
}

func (z *MockZetaCoreBridge) GetCctxByNonce(_ int64, _ uint64) (*crosschaintypes.CrossChainTx, error) {
	if z.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return &crosschaintypes.CrossChainTx{}, nil
}

func (z *MockZetaCoreBridge) GetOutboundTracker(_ chains.Chain, _ uint64) (*crosschaintypes.OutboundTracker, error) {
	if z.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return &crosschaintypes.OutboundTracker{}, nil
}

func (z *MockZetaCoreBridge) GetAllOutboundTrackerByChainbound(_ int64, _ interfaces.Order) ([]crosschaintypes.OutboundTracker, error) {
	if z.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return []crosschaintypes.OutboundTracker{}, nil
}

func (z *MockZetaCoreBridge) GetCrosschainFlags() (observerTypes.CrosschainFlags, error) {
	if z.paused {
		return observerTypes.CrosschainFlags{}, errors.New(ErrMsgPaused)
	}
	return observerTypes.CrosschainFlags{}, nil
}

func (z *MockZetaCoreBridge) GetRateLimiterFlags() (crosschaintypes.RateLimiterFlags, error) {
	if z.paused {
		return crosschaintypes.RateLimiterFlags{}, errors.New(ErrMsgPaused)
	}
	if z.rateLimiterFlags == nil {
		return crosschaintypes.RateLimiterFlags{}, errors.New(ErrMsgRPCFailed)
	}
	return *z.rateLimiterFlags, nil
}

func (z *MockZetaCoreBridge) GetObserverList() ([]string, error) {
	if z.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return []string{}, nil
}

func (z *MockZetaCoreBridge) GetKeyGen() (*observerTypes.Keygen, error) {
	if z.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return &observerTypes.Keygen{}, nil
}

func (z *MockZetaCoreBridge) GetBtcTssAddress(_ int64) (string, error) {
	if z.paused {
		return "", errors.New(ErrMsgPaused)
	}
	return testutils.TSSAddressBTCMainnet, nil
}

func (z *MockZetaCoreBridge) GetInboundTrackersForChain(_ int64) ([]crosschaintypes.InboundTracker, error) {
	if z.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return []crosschaintypes.InboundTracker{}, nil
}

func (z *MockZetaCoreBridge) GetLogger() *zerolog.Logger {
	return nil
}

func (z *MockZetaCoreBridge) ZetaChain() chains.Chain {
	return z.zetaChain
}

func (z *MockZetaCoreBridge) Pause() {
	z.paused = true
}

func (z *MockZetaCoreBridge) Unpause() {
	z.paused = false
}

func (z *MockZetaCoreBridge) GetZetaHotKeyBalance() (math.Int, error) {
	if z.paused {
		return math.NewInt(0), errors.New(ErrMsgPaused)
	}
	return math.NewInt(0), nil
}

// ----------------------------------------------------------------------------
// Feed data to the mock zeta bridge for testing
// ----------------------------------------------------------------------------

func (z *MockZetaCoreBridge) WithPendingCctx(chainID int64, cctxs []*crosschaintypes.CrossChainTx) *MockZetaCoreBridge {
	z.pendingCctxs[chainID] = cctxs
	return z
}

func (z *MockZetaCoreBridge) WithRateLimiterFlags(flags *crosschaintypes.RateLimiterFlags) *MockZetaCoreBridge {
	z.rateLimiterFlags = flags
	return z
}

func (z *MockZetaCoreBridge) WithRateLimiterInput(input *crosschaintypes.QueryRateLimiterInputResponse) *MockZetaCoreBridge {
	z.input = input
	return z
}
