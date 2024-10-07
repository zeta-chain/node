package observer

import (
	"context"
	"strconv"
	"testing"

	"cosmossdk.io/math"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/zeta-chain/node/pkg/chains"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/testutil/sample"
	cctxtypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/ton/liteapi"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

type testSuite struct {
	ctx context.Context
	t   *testing.T

	chain       chains.Chain
	chainParams *observertypes.ChainParams

	liteClient *mocks.LiteClient

	zetacore *mocks.ZetacoreClient
	tss      *mocks.TSS
	database *db.DB

	baseObserver *base.Observer

	votesBag []*cctxtypes.MsgVoteInbound
}

func newTestSuite(t *testing.T) *testSuite {
	var (
		ctx = context.Background()

		chain       = chains.TONTestnet
		chainParams = sample.ChainParams(chain.ChainId)

		liteClient = mocks.NewLiteClient(t)

		tss      = mocks.NewTSSAthens3()
		zetacore = mocks.NewZetacoreClient(t).WithKeys(&keys.Keys{})

		testLogger = zerolog.New(zerolog.NewTestWriter(t))
		logger     = base.Logger{Std: testLogger, Compliance: testLogger}
	)

	database, err := db.NewFromSqliteInMemory(true)
	require.NoError(t, err)

	baseObserver, err := base.NewObserver(
		chain,
		*chainParams,
		zetacore,
		tss,
		1,
		1,
		60,
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

		zetacore: zetacore,
		tss:      tss,
		database: database,

		baseObserver: baseObserver,
	}

	// Setup mocks
	ts.zetacore.On("Chain").Return(chain).Maybe()

	setupVotesBag(ts)

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

func (ts *testSuite) MockGetBlockHeader(id ton.BlockIDExt) *mock.Call {
	// let's pretend that block's masterchain ref has the same seqno
	blockInfo := tlb.BlockInfo{
		BlockInfoPart: tlb.BlockInfoPart{MinRefMcSeqno: id.Seqno},
	}

	return ts.liteClient.
		On("GetBlockHeader", mock.Anything, id, uint32(0)).
		Return(blockInfo, nil)
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
		cctx, ok := vote.(*cctxtypes.MsgVoteInbound)
		require.True(ts.t, ok, "unexpected cctx type")

		ts.votesBag = append(ts.votesBag, cctx)
	}
	ts.zetacore.
		On("PostVoteInbound", ts.ctx, mock.Anything, mock.Anything, mock.Anything).
		Maybe().
		Run(catcher).
		Return("", "", nil) // zeta hash, ballot index, error
}
