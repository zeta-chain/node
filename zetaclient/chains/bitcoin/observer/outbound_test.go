package observer

import (
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/mode"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

// MockBTCObserverMainnet creates a mock Bitcoin mainnet observer for testing
func MockBTCObserverMainnet(t *testing.T, tssSigner interfaces.TSSSigner) *Observer {
	// setup mock arguments
	chain := chains.BitcoinMainnet
	params := mocks.MockChainParams(chain.ChainId, 10)

	if tssSigner == nil {
		tssSigner = mocks.NewTSS(t).FakePubKey(testutils.TSSPubKeyMainnet)
	}

	// create mock rpc client
	btcClient := mocks.NewBitcoinClient(t)
	btcClient.On("GetBlockCount", mock.Anything).Return(int64(100), nil)

	database, err := db.NewFromSqliteInMemory(true)
	require.NoError(t, err)

	logger := zerolog.New(zerolog.NewTestWriter(t))
	baseLogger := base.Logger{Std: logger, Compliance: logger}

	baseObserver, err := base.NewObserver(chain, params, zrepo.New(nil, chain, mode.StandardMode),
		tssSigner, 100, nil, database, baseLogger)
	require.NoError(t, err)

	// create Bitcoin observer
	ob, err := New(baseObserver, btcClient, chain)
	require.NoError(t, err)

	return ob
}

func TestCheckTSSVout(t *testing.T) {
	// the archived outbound raw result file and cctx file
	// https://blockstream.info/tx/030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0
	chain := chains.BitcoinMainnet
	chainID := chain.ChainId
	nonce := uint64(148)

	// create mainnet mock client
	ob := MockBTCObserverMainnet(t, nil)

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
	ob := MockBTCObserverMainnet(t, nil)

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
