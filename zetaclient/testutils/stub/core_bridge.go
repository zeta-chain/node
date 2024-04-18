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
	cctxtypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

const ErrMsgPaused = "zeta core bridge is paused"

var _ interfaces.ZetaCoreBridger = &MockZetaCoreBridge{}

type MockZetaCoreBridge struct {
	paused    bool
	zetaChain chains.Chain
}

func NewMockZetaCoreBridge() *MockZetaCoreBridge {
	zetaChain, err := chains.ZetaChainFromChainID("zetachain_7000-1")
	if err != nil {
		panic(err)
	}
	return &MockZetaCoreBridge{
		paused:    false,
		zetaChain: zetaChain,
	}
}

func (z *MockZetaCoreBridge) PostVoteInbound(_, _ uint64, _ *cctxtypes.MsgVoteOnObservedInboundTx) (string, string, error) {
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

func (z *MockZetaCoreBridge) AddTxHashToOutTxTracker(_ int64, _ uint64, _ string, _ *proofs.Proof, _ string, _ int64) (string, error) {
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

func (z *MockZetaCoreBridge) GetLastBlockHeightByChain(_ chains.Chain) (*cctxtypes.LastBlockHeight, error) {
	if z.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return &cctxtypes.LastBlockHeight{}, nil
}

func (z *MockZetaCoreBridge) ListPendingCctx(_ int64) ([]*cctxtypes.CrossChainTx, uint64, error) {
	if z.paused {
		return nil, 0, errors.New(ErrMsgPaused)
	}
	return []*cctxtypes.CrossChainTx{}, 0, nil
}

func (z *MockZetaCoreBridge) ListPendingCctxWithRatelimit() ([]*cctxtypes.CrossChainTx, uint64, bool, error) {
	if z.paused {
		return nil, 0, false, errors.New(ErrMsgPaused)
	}
	return []*cctxtypes.CrossChainTx{}, 0, false, nil
}

func (z *MockZetaCoreBridge) GetPendingNoncesByChain(_ int64) (observerTypes.PendingNonces, error) {
	if z.paused {
		return observerTypes.PendingNonces{}, errors.New(ErrMsgPaused)
	}
	return observerTypes.PendingNonces{}, nil
}

func (z *MockZetaCoreBridge) GetCctxByNonce(_ int64, _ uint64) (*cctxtypes.CrossChainTx, error) {
	if z.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return &cctxtypes.CrossChainTx{}, nil
}

func (z *MockZetaCoreBridge) GetOutTxTracker(_ chains.Chain, _ uint64) (*cctxtypes.OutTxTracker, error) {
	if z.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return &cctxtypes.OutTxTracker{}, nil
}

func (z *MockZetaCoreBridge) GetAllOutTxTrackerByChain(_ int64, _ interfaces.Order) ([]cctxtypes.OutTxTracker, error) {
	if z.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return []cctxtypes.OutTxTracker{}, nil
}

func (z *MockZetaCoreBridge) GetCrosschainFlags() (observerTypes.CrosschainFlags, error) {
	if z.paused {
		return observerTypes.CrosschainFlags{}, errors.New(ErrMsgPaused)
	}
	return observerTypes.CrosschainFlags{}, nil
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

func (z *MockZetaCoreBridge) GetInboundTrackersForChain(_ int64) ([]cctxtypes.InTxTracker, error) {
	if z.paused {
		return nil, errors.New(ErrMsgPaused)
	}
	return []cctxtypes.InTxTracker{}, nil
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
