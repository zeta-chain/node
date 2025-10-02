package signer

import (
	"context"
	"encoding/hex"
	"strconv"
	"testing"

	"cosmossdk.io/math"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/testutil/sample"
	cc "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/ton/encoder"
	"github.com/zeta-chain/node/zetaclient/chains/ton/rpc"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/mode"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

func TestSigner(t *testing.T) {
	// ARRANGE
	ts := newTestSuite(t)

	// Given TON signer
	signer := New(ts.baseSigner, ts.rpc, ts.gw)

	// Given a sample TON receiver
	receiver := ton.MustParseAccountID("0QAyaVdkvWSuax8luWhDXY_0X9Am1ASWlJz4OI7M-jqcM5wK")

	const (
		zetaHeight = 123
		outboundID = "abc123"
		nonce      = 2
	)

	amount := tonCoins(t, "5")

	// Given CCTX#1
	cctx1 := sample.CrossChainTx(t, "123")
	cctx1.InboundParams.CoinType = coin.CoinType_Gas
	cctx1.OutboundParams = []*cc.OutboundParams{{
		Receiver:        receiver.ToRaw(),
		ReceiverChainId: ts.chain.ChainId,
		CoinType:        coin.CoinType_Gas,
		Amount:          amount,
		TssNonce:        nonce,
	}}

	// Given expected withdrawal
	withdrawal := toncontracts.Withdrawal{
		Recipient: receiver,
		Amount:    amount,
		Seqno:     nonce,
	}

	ts.Sign(&withdrawal)

	// Given CCTX#2 (invalid receiver -1)
	cctx2 := sample.CrossChainTx(t, "456")
	cctx2.InboundParams.CoinType = coin.CoinType_Gas
	cctx2.OutboundParams = []*cc.OutboundParams{{
		Receiver:        "-1:3eb9fbd865df68188ea84d6615086c26a7b5912a60bc55fded2cdb029b67ff0e",
		ReceiverChainId: ts.chain.ChainId,
		CoinType:        coin.CoinType_Gas,
		Amount:          amount,
		TssNonce:        nonce + 1,
	}}

	// Given expected increase_seqno
	increaseSeqno := toncontracts.IncreaseSeqno{
		ReasonCode: uint32(InvalidWorkchain),
		Seqno:      nonce + 1,
	}

	ts.Sign(&increaseSeqno)

	// Given expected rpc calls
	lt, hash := uint64(400), decodeHash(t, "df8a01053f50a74503dffe6802f357bf0e665bd1f3d082faccfebdea93cddfeb")
	ts.OnGetAccountState(ts.gw.AccountID(), rpc.Account{LastTxLT: lt, LastTxHash: hash})

	ts.OnSendMessage(0, nil)

	var (
		withdrawalTx    = sample.TONWithdrawal(t, ts.gw.AccountID(), withdrawal)
		increaseSeqnoTx = sample.TONIncreaseSeqno(t, ts.gw.AccountID(), increaseSeqno)
	)

	// returns both txs
	ts.OnGetTransactionsSince(
		ts.gw.AccountID(),
		lt,
		ton.Bits256(hash),
		[]ton.Transaction{withdrawalTx, increaseSeqnoTx},
		nil,
	)

	// ACT
	signer.TryProcessOutbound(ts.ctx, cctx1, ts.zetacore, zetaHeight)
	signer.TryProcessOutbound(ts.ctx, cctx2, ts.zetacore, zetaHeight)

	// ASSERT
	// Make sure signer send the tx the chain AND published the outbound tracker
	require.Len(t, ts.trackerBag, 2)

	tracker1 := ts.trackerBag[0]

	require.Equal(t, uint64(nonce), tracker1.nonce)
	require.Equal(t, encoder.EncodeTx(withdrawalTx), tracker1.hash)

	tracker2 := ts.trackerBag[1]
	require.Equal(t, uint64(nonce+1), tracker2.nonce)
	require.Equal(t, encoder.EncodeTx(increaseSeqnoTx), tracker2.hash)
}

type testSuite struct {
	ctx context.Context
	t   *testing.T

	chain       chains.Chain
	chainParams *observertypes.ChainParams

	rpc *mocks.TONRPC

	zetacore *mocks.ZetacoreClient
	tss      *mocks.TSS

	gw         *toncontracts.Gateway
	baseSigner *base.Signer

	trackerBag []testTracker
}

type testTracker struct {
	nonce uint64
	hash  string
}

func newTestSuite(t *testing.T) *testSuite {
	var (
		ctx = context.Background()

		chain       = chains.TONTestnet
		chainParams = sample.ChainParams(chain.ChainId)

		rpc = mocks.NewTONRPC(t)

		tss      = mocks.NewTSS(t)
		zetacore = mocks.NewZetacoreClient(t).WithKeys(&keys.Keys{})

		testLogger = zerolog.New(zerolog.NewTestWriter(t))
		logger     = base.Logger{Std: testLogger, Compliance: testLogger}

		gwAccountID = ton.MustParseAccountID(testutils.GatewayAddresses[chain.ChainId])
	)

	ts := &testSuite{
		ctx: ctx,
		t:   t,

		chain:       chain,
		chainParams: chainParams,

		rpc: rpc,

		zetacore: zetacore,
		tss:      tss,

		gw:         toncontracts.NewGateway(gwAccountID),
		baseSigner: base.NewSigner(chain, tss, logger, mode.StandardMode),
	}

	// Setup mocks
	ts.zetacore.On("Chain").Return(chain).Maybe()

	setupTrackersBag(ts)

	return ts
}

func (ts *testSuite) OnGetAccountState(acc ton.AccountID, state rpc.Account) *mock.Call {
	return ts.rpc.On("GetAccountState", mock.Anything, acc).Return(state, nil)
}

func (ts *testSuite) OnSendMessage(id uint32, err error) *mock.Call {
	return ts.rpc.On("SendMessage", mock.Anything, mock.Anything).Return(id, err)
}

func (ts *testSuite) OnGetTransactionsSince(
	acc ton.AccountID,
	lt uint64,
	hash ton.Bits256,
	txs []ton.Transaction,
	err error,
) *mock.Call {
	return ts.rpc.
		On("GetTransactionsSince", mock.Anything, acc, lt, hash).
		Return(txs, err)
}

func (ts *testSuite) Sign(msg toncontracts.ExternalMsg) {
	hash, err := msg.Hash()
	require.NoError(ts.t, err)

	sig, err := ts.tss.Sign(ts.ctx, hash[:], 0, 0, 0)
	require.NoError(ts.t, err)

	msg.SetSignature(sig)
}

// parses string to TON
func tonCoins(t *testing.T, raw string) math.Uint {
	t.Helper()

	const oneTON = 1_000_000_000

	f, err := strconv.ParseFloat(raw, 64)
	require.NoError(t, err)

	f *= oneTON

	return math.NewUint(uint64(f))
}

func decodeHash(t *testing.T, raw string) tlb.Bits256 {
	t.Helper()

	h, err := hex.DecodeString(raw)
	require.NoError(t, err)

	var hash tlb.Bits256

	copy(hash[:], h)

	return hash
}

func setupTrackersBag(ts *testSuite) {
	catcher := func(args mock.Arguments) {
		require.Equal(ts.t, ts.chain.ChainId, args.Get(1).(int64))
		nonce := args.Get(2).(uint64)
		txHash := args.Get(3).(string)

		ts.t.Logf("Adding outbound tracker: nonce=%d, hash=%s", nonce, txHash)

		ts.trackerBag = append(ts.trackerBag, testTracker{nonce, txHash})
	}

	ts.zetacore.On(
		"PostOutboundTracker",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Maybe().Run(catcher).Return("", nil)
}
