package observer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/testutils"
)

func TestObserveMempool(t *testing.T) {
	var (
		sampleTx1 = btcutil.NewTx(wire.NewMsgTx(wire.TxVersion))
		sampleTx2 = btcutil.NewTx(wire.NewMsgTx(2))
	)

	txFinder := func(tx *btcutil.Tx, nonce uint64, errMsg string) pendingTxFinder {
		return func(ctx context.Context) (*btcutil.Tx, uint64, error) {
			var err error
			if errMsg != "" {
				err = errors.New(errMsg)
			}

			return tx, nonce, err
		}
	}

	mockStuckRPC := func(stuck bool, stuckFor time.Duration, errMsg string) func(ts *testSuite) {
		var err error
		if errMsg != "" {
			err = errors.New(errMsg)
		}

		return func(ts *testSuite) {
			ts.client.
				On("IsTxStuckInMempool", mock.Anything, mock.Anything, mock.Anything).
				Maybe().
				Return(stuck, stuckFor, err)
		}
	}

	for _, tt := range []struct {
		name        string
		txFinder    pendingTxFinder
		txChecker   func(ts *testSuite)
		oldStuckTx  *LastStuckOutbound
		expectedTx  *LastStuckOutbound
		errContains string
	}{
		{
			name:       "should set last stuck tx successfully",
			txFinder:   txFinder(sampleTx1, 1, ""),
			txChecker:  mockStuckRPC(true, 30*time.Minute, ""),
			oldStuckTx: nil,
			expectedTx: newLastStuckOutbound(1, sampleTx1, 30*time.Minute),
		},
		{
			name:       "should update last stuck tx successfully",
			txFinder:   txFinder(sampleTx2, 2, ""),
			txChecker:  mockStuckRPC(true, 40*time.Minute, ""),
			oldStuckTx: newLastStuckOutbound(1, sampleTx1, 30*time.Minute),
			expectedTx: newLastStuckOutbound(2, sampleTx2, 40*time.Minute),
		},
		{
			name:       "should clear last stuck tx successfully",
			txFinder:   txFinder(sampleTx1, 1, ""),
			txChecker:  mockStuckRPC(false, 1*time.Minute, ""),
			oldStuckTx: newLastStuckOutbound(1, sampleTx1, 30*time.Minute),
			expectedTx: nil,
		},
		{
			name:       "do nothing if unable to find last pending tx",
			txFinder:   txFinder(nil, 0, "txFinder failed"),
			expectedTx: nil,
		},
		{
			name:        "should return error if txChecker failed",
			txFinder:    txFinder(sampleTx1, 1, ""),
			txChecker:   mockStuckRPC(false, 0, "txChecker failed"),
			errContains: "cannot determine",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			ctx := context.Background()
			ts := newTestSuite(t, chains.BitcoinMainnet)

			if tt.txFinder != nil {
				ctx = withPendingTxFinder(ctx, tt.txFinder)
			}

			if tt.txChecker != nil {
				tt.txChecker(ts)
			}

			if tt.oldStuckTx != nil {
				ts.setLastStuckOutbound(tt.oldStuckTx)
			}

			// ACT
			err := ts.ObserveMempool(ctx)

			// ASSERT
			if tt.errContains != "" {
				require.ErrorContains(t, err, tt.errContains)
				return
			}

			require.NoError(t, err)

			// mimic access from Signer's side.
			stuckTx, _ := ts.LastStuckOutbound()
			require.Equal(t, tt.expectedTx, stuckTx)
		})
	}
}

func Test_GetLastPendingOutbound(t *testing.T) {
	sampleTx := btcutil.NewTx(wire.NewMsgTx(wire.TxVersion))
	tssAddr := testutils.TSSAddressBTCMainnet

	// makeUTXOsWithNonceMark is a helper function to create UTXOs with nonce mark
	makeUTXOsWithNonceMark := func(tssAddr string, vout uint32, nonce uint64) []btcjson.ListUnspentResult {
		return []btcjson.ListUnspentResult{
			{
				TxID:    sampleTx.MsgTx().TxID(),
				Vout:    vout,
				Address: tssAddr,
				Amount:  float64(chains.NonceMarkAmount(nonce)) / btcutil.SatoshiPerBitcoin,
			},
		}
	}

	tests := []struct {
		name          string
		txNonce       uint64
		pendingNonce  uint64
		pendingNonces *observertypes.PendingNonces
		utxos         []btcjson.ListUnspentResult
		tx            *btcutil.Tx
		saveTx        bool
		includeTx     bool
		failMempool   bool
		failGetRawTx  bool
		expectedTx    *btcutil.Tx
		expectedNonce uint64
		errMsg        string
	}{
		{
			name:         "should return last broadcasted (and stuck) outbound",
			txNonce:      9,
			pendingNonce: 9, // tx is stuck, pending nonce hasn't increased to 10 yet
			pendingNonces: &observertypes.PendingNonces{
				NonceLow:  9,
				NonceHigh: 10,
			},
			utxos:         makeUTXOsWithNonceMark(tssAddr, 0, 9), // nonce mark utxo is produced even if tx is stuck
			tx:            sampleTx,
			saveTx:        true,
			includeTx:     false, // tx is stuck, cannot be included
			expectedTx:    sampleTx,
			expectedNonce: 9,
		},
		{
			name:         "should return last included (no stuck) outbound",
			txNonce:      9,
			pendingNonce: 10, // pending nonce already increased to 10 when tx is included (not stuck)
			pendingNonces: &observertypes.PendingNonces{
				NonceLow:  9,
				NonceHigh: 10,
			},
			utxos:         makeUTXOsWithNonceMark(tssAddr, 0, 9),
			tx:            sampleTx,
			saveTx:        false, // both true and false are fine
			includeTx:     true,
			expectedTx:    sampleTx,
			expectedNonce: 9,
		},
		{
			name:          "return error if pending nonce is zero",
			pendingNonce:  0,
			expectedTx:    nil,
			expectedNonce: 0,
			errMsg:        "pending nonce is zero",
		},
		{
			name:          "return error if GetPendingNoncesByChain failed",
			pendingNonce:  10,
			saveTx:        true,
			includeTx:     false,
			expectedTx:    nil,
			expectedNonce: 0,
			errMsg:        "GetPendingNoncesByChain failed",
		},
		{
			name:         "return error if no last tx found",
			txNonce:      9,
			pendingNonce: 9, // this zetaclient completely lost track of the tx, pending nonce hasn't increased to 10 yet
			pendingNonces: &observertypes.PendingNonces{
				NonceLow:  9,
				NonceHigh: 10,
			},
			saveTx:        false, // lost track of the tx hash in db
			includeTx:     false, // lost track of the tx hash in outbound tracker
			expectedTx:    nil,
			expectedNonce: 0,
			errMsg:        "last tx not found",
		},
		{
			name:         "return error if GetMempoolEntry failed",
			txNonce:      9,
			pendingNonce: 10,
			pendingNonces: &observertypes.PendingNonces{
				NonceLow:  9,
				NonceHigh: 10,
			},
			tx:            sampleTx,
			saveTx:        true,
			includeTx:     false,
			failMempool:   true,
			expectedTx:    nil,
			expectedNonce: 0,
			errMsg:        "not in mempool",
		},
		{
			name:         "return error if FetchUTXOs failed",
			txNonce:      9,
			pendingNonce: 10,
			pendingNonces: &observertypes.PendingNonces{
				NonceLow:  9,
				NonceHigh: 10,
			},
			tx:            sampleTx,
			saveTx:        true,
			includeTx:     false,
			expectedTx:    nil,
			expectedNonce: 0,
			errMsg:        "unable to fetch UTXOs: failed",
		},
		{
			name:         "return error if unable to find nonce-mark UTXO",
			txNonce:      9,
			pendingNonce: 10,
			pendingNonces: &observertypes.PendingNonces{
				NonceLow:  9,
				NonceHigh: 10,
			},
			utxos:         makeUTXOsWithNonceMark(tssAddr, 1, 9), // wrong output index, should be 0
			tx:            sampleTx,
			saveTx:        true,
			includeTx:     false,
			expectedTx:    nil,
			expectedNonce: 0,
			errMsg:        "findNonceMarkUTXO failed",
		},
		{
			name:         "return error if GetRawTransactionByStr failed",
			txNonce:      9,
			pendingNonce: 10,
			pendingNonces: &observertypes.PendingNonces{
				NonceLow:  9,
				NonceHigh: 10,
			},
			utxos:         makeUTXOsWithNonceMark(tssAddr, 0, 9),
			tx:            sampleTx,
			saveTx:        false,
			includeTx:     true,
			failGetRawTx:  true,
			expectedTx:    nil,
			expectedNonce: 0,
			errMsg:        "GetRawTransactionByStr failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create observer
			ob := newTestSuite(t, chains.BitcoinMainnet)

			// mock up pending nonce
			if tt.pendingNonce > 0 {
				mockAndRefreshPendingNonce(t, ob, "abc", tt.pendingNonce)
			}

			// mock zetacore client response
			if tt.pendingNonces != nil {
				ob.zetacore.On("GetPendingNoncesByChain", mock.Anything, mock.Anything).
					Maybe().
					Return(*tt.pendingNonces, nil)
			} else {
				empty := observertypes.PendingNonces{}
				ob.zetacore.On("GetPendingNoncesByChain", mock.Anything, mock.Anything).Maybe().Return(empty, errors.New("failed"))
			}

			if tt.tx != nil {
				// save tx to simulate broadcasted tx
				if tt.saveTx {
					ob.SaveBroadcastedTx(tt.tx.MsgTx().TxID(), tt.txNonce)
				}

				// include tx to simulate included tx
				if tt.includeTx {
					ob.SetIncludedTx(tt.txNonce, &btcjson.GetTransactionResult{TxID: tt.tx.MsgTx().TxID()})
				}
			}

			// mock btc client response
			if tt.utxos != nil {
				ob.client.On("ListUnspentMinMaxAddresses", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Maybe().
					Return(tt.utxos, nil)
			} else {
				ob.client.On("ListUnspentMinMaxAddresses", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return(nil, errors.New("failed"))
			}
			if !tt.failMempool {
				ob.client.On("GetMempoolEntry", mock.Anything, mock.Anything).Maybe().Return(nil, nil)
			} else {
				ob.client.On("GetMempoolEntry", mock.Anything, mock.Anything).Maybe().Return(nil, errors.New("failed"))
			}
			if tt.tx != nil && !tt.failGetRawTx {
				ob.client.On("GetRawTransactionByStr", mock.Anything, mock.Anything).Maybe().Return(tt.tx, nil)
			} else {
				ob.client.On("GetRawTransactionByStr", mock.Anything, mock.Anything).Maybe().Return(nil, errors.New("failed"))
			}

			ctx := context.Background()
			lastTx, lastNonce, err := ob.getLastPendingOutbound(ctx)

			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
				require.Nil(t, lastTx)
				require.Zero(t, lastNonce)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedTx, lastTx)
			require.Equal(t, tt.expectedNonce, lastNonce)
		})
	}
}

// mockAndRefreshPendingNonce is a helper function to set pending nonce as 'nonce' by mocking
// all the mocked RPCs are one-time use and will be called only once
func mockAndRefreshPendingNonce(t *testing.T, s *testSuite, txid string, nonce uint64) {
	ctx := context.Background()
	cctx := sample.CrossChainTx(t, "0x123")
	cctx.GetCurrentOutboundParam().Hash = txid
	// #nosec G115 always positive in test
	nonceInt := int64(nonce)
	pendingNonces := &observertypes.PendingNonces{NonceLow: nonceInt, NonceHigh: nonceInt}

	s.zetacore.On("GetPendingNoncesByChain", mock.Anything, mock.Anything).Return(*pendingNonces, nil).Once()
	s.client.On("ListUnspentMinMaxAddresses", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil).
		Once()

	// fetch utxos and refresh pending nonce
	err := s.FetchUTXOs(ctx)
	require.NoError(t, err)
}

func withPendingTxFinder(ctx context.Context, fn pendingTxFinder) context.Context {
	return context.WithValue(ctx, pendingTxFinderKey{}, fn)
}
