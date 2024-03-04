package stub

import (
	"math/big"

	"cosmossdk.io/math"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/go-tss/blame"
	"github.com/zeta-chain/zetacore/common"
	cctxtypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

var _ interfaces.ZetaCoreBridger = &ZetaCoreBridge{}

type ZetaCoreBridge struct {
	zetaChain common.Chain
}

func (z ZetaCoreBridge) PostVoteInbound(_, _ uint64, _ *cctxtypes.MsgVoteOnObservedInboundTx) (string, string, error) {
	return "", "", nil
}

func (z ZetaCoreBridge) PostVoteOutbound(_ string, _ string, _ uint64, _ uint64, _ *big.Int, _ uint64, _ *big.Int, _ common.ReceiveStatus, _ common.Chain, _ uint64, _ common.CoinType) (string, string, error) {
	return "", "", nil
}

func (z ZetaCoreBridge) PostGasPrice(_ common.Chain, _ uint64, _ string, _ uint64) (string, error) {
	return "", nil
}

func (z ZetaCoreBridge) PostAddBlockHeader(_ int64, _ []byte, _ int64, _ common.HeaderData) (string, error) {
	return "", nil
}

func (z ZetaCoreBridge) GetBlockHeaderStateByChain(_ int64) (observerTypes.QueryGetBlockHeaderStateResponse, error) {
	return observerTypes.QueryGetBlockHeaderStateResponse{}, nil
}

func (z ZetaCoreBridge) PostBlameData(_ *blame.Blame, _ int64, _ string) (string, error) {
	return "", nil
}

func (z ZetaCoreBridge) AddTxHashToOutTxTracker(_ int64, _ uint64, _ string, _ *common.Proof, _ string, _ int64) (string, error) {
	return "", nil
}

func (z ZetaCoreBridge) GetKeys() *keys.Keys {
	return &keys.Keys{}
}

func (z ZetaCoreBridge) GetBlockHeight() (int64, error) {
	return 0, nil
}

func (z ZetaCoreBridge) GetZetaBlockHeight() (int64, error) {
	return 0, nil
}

func (z ZetaCoreBridge) GetLastBlockHeightByChain(_ common.Chain) (*cctxtypes.LastBlockHeight, error) {
	return &cctxtypes.LastBlockHeight{}, nil
}

func (z ZetaCoreBridge) ListPendingCctx(_ int64) ([]*cctxtypes.CrossChainTx, uint64, error) {
	return []*cctxtypes.CrossChainTx{}, 0, nil
}

func (z ZetaCoreBridge) GetPendingNoncesByChain(_ int64) (observerTypes.PendingNonces, error) {
	return observerTypes.PendingNonces{}, nil
}

func (z ZetaCoreBridge) GetCctxByNonce(_ int64, _ uint64) (*cctxtypes.CrossChainTx, error) {
	return &cctxtypes.CrossChainTx{}, nil
}

func (z ZetaCoreBridge) GetOutTxTracker(_ common.Chain, _ uint64) (*cctxtypes.OutTxTracker, error) {
	return &cctxtypes.OutTxTracker{}, nil
}

func (z ZetaCoreBridge) GetAllOutTxTrackerByChain(_ int64, _ interfaces.Order) ([]cctxtypes.OutTxTracker, error) {
	return []cctxtypes.OutTxTracker{}, nil
}

func (z ZetaCoreBridge) GetCrosschainFlags() (observerTypes.CrosschainFlags, error) {
	return observerTypes.CrosschainFlags{}, nil
}

func (z ZetaCoreBridge) GetObserverList() ([]string, error) {
	return []string{}, nil
}

func (z ZetaCoreBridge) GetKeyGen() (*observerTypes.Keygen, error) {
	return &observerTypes.Keygen{}, nil
}

func (z ZetaCoreBridge) GetBtcTssAddress(_ int64) (string, error) {
	return testutils.TSSAddressBTCMainnet, nil
}

func (z ZetaCoreBridge) GetInboundTrackersForChain(_ int64) ([]cctxtypes.InTxTracker, error) {
	return []cctxtypes.InTxTracker{}, nil
}

func (z ZetaCoreBridge) GetLogger() *zerolog.Logger {
	return nil
}

func (z ZetaCoreBridge) ZetaChain() common.Chain {
	return z.zetaChain
}

func (z ZetaCoreBridge) Pause() {
}

func (z ZetaCoreBridge) Unpause() {
}

func (z ZetaCoreBridge) GetZetaHotKeyBalance() (math.Int, error) {
	return math.NewInt(0), nil
}

func NewZetaCoreBridge() *ZetaCoreBridge {
	zetaChain, err := common.ZetaChainFromChainID("zetachain_7000-1")
	if err != nil {
		panic(err)
	}
	return &ZetaCoreBridge{zetaChain: zetaChain}
}
