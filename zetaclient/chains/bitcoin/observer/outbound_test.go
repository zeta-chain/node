package observer

import (
	"context"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

// the relative path to the testdata directory
var TestDataDir = "../../../"

func Test_VoteOutboundIfConfirmed(t *testing.T) {
	// load archived CCTX
	// https://blockstream.info/tx/030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0
	chainID := chains.BitcoinMainnet.ChainId
	nonce := uint64(148)
	confirmParams := observertypes.ConfirmationParams{
		SafeInboundCount:  1,
		SafeOutboundCount: 1,
	}

	tests := []struct {
		name            string
		cctx            *crosschaintypes.CrossChainTx
		txResult        *btcjson.GetTransactionResult
		txHeight        int64
		lastBlock       uint64
		confirmParams   observertypes.ConfirmationParams
		shouldPostVote  bool
		postVoteError   bool
		continueKeysign bool
		errMsg          string
	}{
		{
			name: "should post vote and return false if outbound is SAFE confirmed",
			cctx: testutils.LoadCctxByNonce(t, chainID, nonce),
			txResult: &btcjson.GetTransactionResult{
				Confirmations: 1,
			},
			txHeight:        100,
			lastBlock:       100,
			confirmParams:   confirmParams,
			shouldPostVote:  true,
			continueKeysign: false,
		},
		{
			name: "should post vote and return false if outbound is FAST confirmed",
			cctx: testutils.LoadCctxByNonce(t, chainID, nonce),
			txResult: &btcjson.GetTransactionResult{
				Confirmations: 1,
			},
			txHeight:  100,
			lastBlock: 100,
			confirmParams: observertypes.ConfirmationParams{
				SafeInboundCount:  1,
				SafeOutboundCount: 2,
				FastOutboundCount: 1,
			},
			shouldPostVote:  true,
			continueKeysign: false,
		},
		{
			name:            "should continue keysign if tx is neither broadcasted nor included",
			cctx:            testutils.LoadCctxByNonce(t, chainID, nonce),
			lastBlock:       100,
			confirmParams:   confirmParams,
			continueKeysign: true,
		},
		{
			name: "should return error if unable to get tx height",
			cctx: testutils.LoadCctxByNonce(t, chainID, nonce),
			txResult: &btcjson.GetTransactionResult{
				Confirmations: 1,
			},
			txHeight:        0, // set to 0 to simulate error
			lastBlock:       100,
			confirmParams:   confirmParams,
			continueKeysign: false,
			errMsg:          "error getting block height by hash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			ob := newTestSuite(t, func(cfg *testSuiteConfig) {
				cfg.LastBlock = tt.lastBlock
				cfg.ConfirmationParams = &tt.confirmParams
			})

			// mock tx inclusion
			if tt.txResult != nil {
				ob.SetIncludedTx(tt.cctx.GetCurrentOutboundParam().TssNonce, tt.txResult)

				// mock tx height
				if tt.txHeight > 0 {
					ob.btcClient.On("GetBlockHeightByStr", mock.Anything, mock.Anything).Return(tt.txHeight, nil)
				} else {
					ob.btcClient.On("GetBlockHeightByStr", mock.Anything, mock.Anything).Return(int64(0), errors.New("error"))
				}
			}

			// mock zetacore
			if tt.shouldPostVote {
				if tt.postVoteError {
					ob.zetacore.On("PostVoteOutbound", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
						Return("", "", errors.New("error"))
				} else {
					ob.zetacore.On("PostVoteOutbound", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("sameZetaHash", "sampleBallot", nil)
				}
			}

			// ACT
			ctx := context.Background()
			continueKeysign, err := ob.VoteOutboundIfConfirmed(ctx, tt.cctx)

			// ASSERT
			if tt.errMsg != "" {
				require.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.continueKeysign, continueKeysign)
		})
	}
}

func TestCheckTSSVout(t *testing.T) {
	// the archived outbound raw result file and cctx file
	// https://blockstream.info/tx/030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0
	chain := chains.BitcoinMainnet
	chainID := chain.ChainId
	nonce := uint64(148)

	// create mainnet mock client
	ob := newTestSuite(t)

	t.Run("valid TSS vout should pass", func(t *testing.T) {
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()
		err := ob.checkTSSVout(params, rawResult.Vout)
		require.NoError(t, err)
	})
	t.Run("should fail if vout length < 2 or > 3", func(t *testing.T) {
		_, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		err := ob.checkTSSVout(params, []btcjson.Vout{{}})
		require.ErrorContains(t, err, "invalid number of vouts")

		err = ob.checkTSSVout(params, []btcjson.Vout{{}, {}, {}, {}})
		require.ErrorContains(t, err, "invalid number of vouts")
	})
	t.Run("should fail on invalid TSS vout", func(t *testing.T) {
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		// invalid TSS vout
		rawResult.Vout[0].ScriptPubKey.Hex = "invalid script"
		err := ob.checkTSSVout(params, rawResult.Vout)
		require.Error(t, err)
	})
	t.Run("should fail if vout 0 is not to the TSS address", func(t *testing.T) {
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		// not TSS address, bc1qh297vdt8xq6df5xae9z8gzd4jsu9a392mp0dus
		rawResult.Vout[0].ScriptPubKey.Hex = "0014ba8be635673034d4d0ddc9447409b594385ec4aa"
		err := ob.checkTSSVout(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match TSS address")
	})
	t.Run("should fail if vout 0 not match nonce mark", func(t *testing.T) {
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		// not match nonce mark
		rawResult.Vout[0].Value = 0.00000147
		err := ob.checkTSSVout(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match nonce-mark amount")
	})
	t.Run("should fail if vout 1 is not to the receiver address", func(t *testing.T) {
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		// not receiver address, bc1qh297vdt8xq6df5xae9z8gzd4jsu9a392mp0dus
		rawResult.Vout[1].ScriptPubKey.Hex = "0014ba8be635673034d4d0ddc9447409b594385ec4aa"
		err := ob.checkTSSVout(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match params receiver")
	})
	t.Run("should fail if vout 1 not match payment amount", func(t *testing.T) {
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		// not match payment amount
		rawResult.Vout[1].Value = 0.00011000
		err := ob.checkTSSVout(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match params amount")
	})
	t.Run("should fail if vout 2 is not to the TSS address", func(t *testing.T) {
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		// not TSS address, bc1qh297vdt8xq6df5xae9z8gzd4jsu9a392mp0dus
		rawResult.Vout[2].ScriptPubKey.Hex = "0014ba8be635673034d4d0ddc9447409b594385ec4aa"
		err := ob.checkTSSVout(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match TSS address")
	})
}

func TestCheckTSSVoutCancelled(t *testing.T) {
	// the archived outbound raw result file and cctx file
	// https://blockstream.info/tx/030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0
	chain := chains.BitcoinMainnet
	chainID := chain.ChainId
	nonce := uint64(148)

	// create mainnet mock client
	ob := newTestSuite(t)

	t.Run("valid TSS vout should pass", func(t *testing.T) {
		// remove change vout to simulate cancelled tx
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		rawResult.Vout[1] = rawResult.Vout[2]
		rawResult.Vout = rawResult.Vout[:2]
		params := cctx.GetCurrentOutboundParam()

		err := ob.checkTSSVoutCancelled(params, rawResult.Vout)
		require.NoError(t, err)
	})
	t.Run("should fail if vout length < 1 or > 2", func(t *testing.T) {
		_, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		err := ob.checkTSSVoutCancelled(params, []btcjson.Vout{})
		require.ErrorContains(t, err, "invalid number of vouts")

		err = ob.checkTSSVoutCancelled(params, []btcjson.Vout{{}, {}, {}})
		require.ErrorContains(t, err, "invalid number of vouts")
	})
	t.Run("should fail if vout 0 is not to the TSS address", func(t *testing.T) {
		// remove change vout to simulate cancelled tx
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		rawResult.Vout[1] = rawResult.Vout[2]
		rawResult.Vout = rawResult.Vout[:2]
		params := cctx.GetCurrentOutboundParam()

		// not TSS address, bc1qh297vdt8xq6df5xae9z8gzd4jsu9a392mp0dus
		rawResult.Vout[0].ScriptPubKey.Hex = "0014ba8be635673034d4d0ddc9447409b594385ec4aa"
		err := ob.checkTSSVoutCancelled(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match TSS address")
	})
	t.Run("should fail if vout 0 not match nonce mark", func(t *testing.T) {
		// remove change vout to simulate cancelled tx
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		rawResult.Vout[1] = rawResult.Vout[2]
		rawResult.Vout = rawResult.Vout[:2]
		params := cctx.GetCurrentOutboundParam()

		// not match nonce mark
		rawResult.Vout[0].Value = 0.00000147
		err := ob.checkTSSVoutCancelled(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match nonce-mark amount")
	})
	t.Run("should fail if vout 1 is not to the TSS address", func(t *testing.T) {
		// remove change vout to simulate cancelled tx
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		rawResult.Vout[1] = rawResult.Vout[2]
		rawResult.Vout[1].N = 1 // swap vout index
		rawResult.Vout = rawResult.Vout[:2]
		params := cctx.GetCurrentOutboundParam()

		// not TSS address, bc1qh297vdt8xq6df5xae9z8gzd4jsu9a392mp0dus
		rawResult.Vout[1].ScriptPubKey.Hex = "0014ba8be635673034d4d0ddc9447409b594385ec4aa"
		err := ob.checkTSSVoutCancelled(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match TSS address")
	})
}

func newTestSuite(t *testing.T, opts ...func(*testSuiteConfig)) *testSuite {
	var cfg testSuiteConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	chain := chains.BitcoinMainnet
	if cfg.chain != nil {
		chain = *cfg.chain
	}

	params := mocks.MockChainParams(chain.ChainId, 10)
	if cfg.ConfirmationParams != nil {
		params.ConfirmationParams = cfg.ConfirmationParams
	}

	btcClient := mocks.NewBitcoinClient(t)
	btcClient.On("GetBlockCount", mock.Anything).Return(int64(100), nil)

	zetacore := mocks.NewZetacoreClient(t).WithKeys(&keys.Keys{OperatorAddress: sample.Bech32AccAddress()})

	tss := mocks.NewTSS(t).FakePubKey(testutils.TSSPubKeyMainnet)

	database, err := db.NewFromSqliteInMemory(true)
	require.NoError(t, err)

	log := zerolog.New(zerolog.NewTestWriter(t)).With().Caller().Logger()
	logger := base.Logger{Std: log, Compliance: log}

	baseObserver, err := base.NewObserver(chain, params, zetacore, tss, 100, nil, database, logger)
	require.NoError(t, err)

	ob, err := New(chain, baseObserver, btcClient)
	require.NoError(t, err)

	ob.WithLastBlock(1)
	if cfg.LastBlock > 0 {
		ob.WithLastBlock(cfg.LastBlock)
	}

	return &testSuite{
		Observer:    ob,
		chainParams: &params,
		tss:         tss,
		zetacore:    zetacore,
		btcClient:   btcClient,
	}
}

type testSuite struct {
	*Observer
	chainParams *observertypes.ChainParams
	tss         *mocks.TSS
	zetacore    *mocks.ZetacoreClient
	btcClient   *mocks.BitcoinClient
}

type testSuiteConfig struct {
	chain              *chains.Chain
	LastBlock          uint64
	ConfirmationParams *observertypes.ConfirmationParams
}
