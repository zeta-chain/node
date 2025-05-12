package observer

import (
	"context"
	"strconv"
	"testing"

	"cosmossdk.io/math"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
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
	"github.com/zeta-chain/node/zetaclient/chains/ton/liteapi"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
	"github.com/zeta-chain/node/zetaclient/testutils/testlog"
)

type testSuite struct {
	ctx context.Context
	t   *testing.T

	chain       chains.Chain
	chainParams *observertypes.ChainParams

	gateway    *toncontracts.Gateway
	liteClient *mocks.LiteClient

	zetacore *mocks.ZetacoreClient
	tss      *mocks.TSS
	database *db.DB
	logger   *testlog.Log

	baseObserver *base.Observer

	votesBag   []*cc.MsgVoteInbound
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

		gateway = toncontracts.NewGateway(ton.MustParseAccountID(
			"0:997d889c815aeac21c47f86ae0e38383efc3c3463067582f6263ad48c5a1485b",
		))

		liteClient = mocks.NewLiteClient(t)

		tss      = mocks.NewTSS(t)
		zetacore = mocks.NewZetacoreClient(t).WithKeys(&keys.Keys{
			OperatorAddress: sample.Bech32AccAddress(),
		})

		testLogger = testlog.New(t)
		logger     = base.Logger{Std: testLogger.Logger, Compliance: testLogger.Logger}
	)

	database, err := db.NewFromSqliteInMemory(true)
	require.NoError(t, err)

	baseObserver, err := base.NewObserver(
		chain,
		*chainParams,
		zetacore,
		tss,
		1,
		nil,
		database,
		logger,
	)

	require.NoError(t, err)

	ts := &testSuite{
		ctx: ctx,
		t:   t,

		chain:       chain,
		chainParams: chainParams,

		liteClient: liteClient,
		gateway:    gateway,

		zetacore: zetacore,
		tss:      tss,
		database: database,
		logger:   testLogger,

		baseObserver: baseObserver,
	}

	// Setup mocks
	ts.zetacore.On("Chain").Return(chain).Maybe()

	setupVotesBag(ts)
	setupTrackersBag(ts)

	return ts
}

func (ts *testSuite) SetupLastScannedTX(gw ton.AccountID) ton.Transaction {
	lastScannedTX := sample.TONDonation(ts.t, gw, toncontracts.Donation{
		Sender: sample.GenerateTONAccountID(),
		Amount: tonCoins(ts.t, "1"),
	})

	txHash := liteapi.TransactionHashToString(lastScannedTX.Lt, ton.Bits256(lastScannedTX.Hash()))

	ts.baseObserver.WithLastTxScanned(txHash)
	require.NoError(ts.t, ts.baseObserver.WriteLastTxScannedToDB(txHash))

	return lastScannedTX
}

func (ts *testSuite) OnGetFirstTransaction(acc ton.AccountID, tx *ton.Transaction, scanned int, err error) *mock.Call {
	return ts.liteClient.
		On("GetFirstTransaction", ts.ctx, acc).
		Return(tx, scanned, err)
}

func (ts *testSuite) MockGetTransaction(acc ton.AccountID, tx ton.Transaction) *mock.Call {
	return ts.liteClient.
		On("GetTransaction", mock.Anything, acc, tx.Lt, ton.Bits256(tx.Hash())).
		Return(tx, nil)
}

func (ts *testSuite) MockCCTXByNonce(cctx *cc.CrossChainTx) *mock.Call {
	nonce := cctx.GetCurrentOutboundParam().TssNonce

	return ts.zetacore.On("GetCctxByNonce", ts.ctx, ts.chain.ChainId, nonce).Return(cctx, nil)
}

func (ts *testSuite) OnGetTransactionsSince(
	acc ton.AccountID,
	lt uint64,
	hash ton.Bits256,
	txs []ton.Transaction,
	err error,
) *mock.Call {
	return ts.liteClient.
		On("GetTransactionsSince", mock.Anything, acc, lt, hash).
		Return(txs, err)
}

func (ts *testSuite) OnGetAllOutboundTrackerByChain(trackers []cc.OutboundTracker) *mock.Call {
	return ts.zetacore.
		On("GetAllOutboundTrackerByChain", mock.Anything, ts.chain.ChainId, mock.Anything).
		Return(trackers, nil)
}

func (ts *testSuite) MockGetBlockHeader(id ton.BlockIDExt) *mock.Call {
	// let's pretend that block's masterchain ref has the same seqno
	blockInfo := tlb.BlockInfo{
		BlockInfoPart: tlb.BlockInfoPart{MinRefMcSeqno: id.Seqno},
	}

	return ts.liteClient.
		On("GetBlockHeader", mock.Anything, id, uint32(0)).
		Return(blockInfo, nil)
}

func (ts *testSuite) MockGetCctxByHash() *mock.Call {
	return ts.zetacore.
		On("GetCctxByHash", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))
}

func (ts *testSuite) OnGetInboundTrackersForChain(trackers []cc.InboundTracker) *mock.Call {
	return ts.zetacore.
		On("GetInboundTrackersForChain", mock.Anything, ts.chain.ChainId).
		Return(trackers, nil)
}

func (ts *testSuite) TxToInboundTracker(tx ton.Transaction) cc.InboundTracker {
	return cc.InboundTracker{
		ChainId:  ts.chain.ChainId,
		TxHash:   liteapi.TransactionToHashString(tx),
		CoinType: coin.CoinType_Gas,
	}
}

type signable interface {
	Hash() ([32]byte, error)
	SetSignature([65]byte)
	Signer() (eth.Address, error)
}

func (ts *testSuite) sign(msg signable) {
	hash, err := msg.Hash()
	require.NoError(ts.t, err)

	sig, err := ts.tss.Sign(ts.ctx, hash[:], 0, 0, 0)
	require.NoError(ts.t, err)

	msg.SetSignature(sig)

	// double check
	evmSigner, err := msg.Signer()
	require.NoError(ts.t, err)
	require.Equal(ts.t, ts.tss.PubKey().AddressEVM().String(), evmSigner.String())
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

func setupVotesBag(ts *testSuite) {
	catcher := func(args mock.Arguments) {
		vote := args.Get(3)
		cctx, ok := vote.(*cc.MsgVoteInbound)
		require.True(ts.t, ok, "unexpected cctx type")

		ts.votesBag = append(ts.votesBag, cctx)
	}
	ts.zetacore.
		On("PostVoteInbound", ts.ctx, mock.Anything, mock.Anything, mock.Anything).
		Maybe().
		Run(catcher).
		Return("", "", nil) // zeta hash, ballot index, error
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
