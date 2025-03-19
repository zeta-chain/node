package signer_test

import (
	"context"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/signer"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

// mempoolTxsInfo is a helper struct to contain mempool txs information
type mempoolTxsInfo struct {
	totalTxs   int64
	totalFees  float64
	totalVSize int64
	avgFeeRate int64
}

func newMempoolTxsInfo(totalTxs int64, totalFees float64, totalVSize int64, avgFeeRate int64) *mempoolTxsInfo {
	return &mempoolTxsInfo{
		totalTxs:   totalTxs,
		totalFees:  totalFees,
		totalVSize: totalVSize,
		avgFeeRate: avgFeeRate,
	}
}

func Test_NewCPFPFeeBumper(t *testing.T) {
	tests := []struct {
		name         string
		chain        chains.Chain
		client       *mocks.BitcoinClient
		tx           *btcutil.Tx
		cctxRate     int64
		liveRate     int64
		minRelayFee  float64
		memplTxsInfo *mempoolTxsInfo
		errMsg       string
		expected     *signer.CPFPFeeBumper
	}{
		{
			chain:       chains.BitcoinMainnet,
			name:        "should create new CPFPFeeBumper successfully",
			client:      mocks.NewBitcoinClient(t),
			tx:          btcutil.NewTx(wire.NewMsgTx(wire.TxVersion)),
			cctxRate:    10,
			liveRate:    12,
			minRelayFee: 0.00001,
			memplTxsInfo: newMempoolTxsInfo(
				2,      // 2 stuck TSS txs
				0.0001, // total fees 0.0001 BTC
				1000,   // total vsize 1000
				10,     // average fee rate 10 sat/vB
			),
			expected: &signer.CPFPFeeBumper{
				Ctx:         context.Background(),
				Chain:       chains.BitcoinMainnet,
				Tx:          btcutil.NewTx(wire.NewMsgTx(wire.TxVersion)),
				MinRelayFee: 0.00001,
				CCTXRate:    10,
				LiveRate:    12,
				TotalTxs:    2,
				TotalFees:   10000,
				TotalVSize:  1000,
				AvgFeeRate:  10,
				Logger:      log.Logger,
			},
		},
		{
			chain:        chains.BitcoinMainnet,
			name:         "should fail when mempool txs info fetcher returns error",
			client:       mocks.NewBitcoinClient(t),
			tx:           btcutil.NewTx(wire.NewMsgTx(wire.TxVersion)),
			liveRate:     12,
			memplTxsInfo: nil,
			errMsg:       "unable to fetch mempool txs info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			// mock RPC fee rate
			tt.client.On("GetEstimatedFeeRate", mock.Anything, mock.Anything, mock.Anything).Return(tt.liveRate, nil)

			// mock mempool txs information
			if tt.memplTxsInfo != nil {
				tt.client.On("GetTotalMempoolParentsSizeNFees", mock.Anything, mock.Anything, mock.Anything).
					Return(tt.memplTxsInfo.totalTxs, tt.memplTxsInfo.totalFees, tt.memplTxsInfo.totalVSize, tt.memplTxsInfo.avgFeeRate, nil)
			} else {
				v := int64(0)
				tt.client.On("GetTotalMempoolParentsSizeNFees", mock.Anything, mock.Anything, mock.Anything).Return(v, 0.0, v, v, errors.New("rpc error"))
			}

			// ACT
			bumper, err := signer.NewCPFPFeeBumper(
				context.Background(),
				tt.client,
				tt.chain,
				tt.tx,
				tt.cctxRate,
				tt.minRelayFee,
				log.Logger,
			)

			// ASSERT
			if tt.errMsg != "" {
				require.Nil(t, bumper)
				require.ErrorContains(t, err, tt.errMsg)
			} else {
				bumper.RPC = nil // ignore the RPC
				require.NoError(t, err)
				require.Equal(t, tt.expected, bumper)
			}
		})
	}
}

func Test_BumpTxFee(t *testing.T) {
	// https://mempool.space/tx/030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0
	chain := chains.BitcoinMainnet
	txid := "030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0"
	msgTx := testutils.LoadBTCMsgTx(t, TestDataDir, chain.ChainId, txid)

	tests := []struct {
		name            string
		feeBumper       *signer.CPFPFeeBumper
		additionalFees  int64
		expectedNewRate int64
		expectedNewTx   *wire.MsgTx
		errMsg          string
	}{
		{
			name: "should bump tx fee successfully",
			feeBumper: &signer.CPFPFeeBumper{
				Tx:          btcutil.NewTx(msgTx),
				MinRelayFee: 0.00001,
				CCTXRate:    55,
				LiveRate:    67,
				TotalFees:   27213,
				TotalVSize:  579,
				AvgFeeRate:  47,
				Logger:      log.Logger,
			},
			additionalFees:  4632, // 579*55 - 27213
			expectedNewRate: 55,
			expectedNewTx: func() *wire.MsgTx {
				// deduct additional fees
				newTx := signer.CopyMsgTxNoWitness(msgTx)
				newTx.TxOut[2].Value -= 4632
				return newTx
			}(),
		},
		{
			name: "should give up all reserved bump fees",
			feeBumper: &signer.CPFPFeeBumper{
				Tx: func() *btcutil.Tx {
					// modify reserved bump fees to barely cover bump fees
					newTx := msgTx.Copy()
					newTx.TxOut[2].Value = 57*579 - 27213 + constant.BTCWithdrawalDustAmount - 1 // 6789
					return btcutil.NewTx(newTx)
				}(),
				MinRelayFee: 0.00001,
				CCTXRate:    57,
				LiveRate:    67,
				TotalFees:   27213,
				TotalVSize:  579,
				AvgFeeRate:  47,
				Logger:      log.Logger,
			},
			additionalFees:  6789, // same as the reserved value in 2nd output
			expectedNewRate: 59,   // (27213 + 6789) / 579 ≈ 59
			expectedNewTx: func() *wire.MsgTx {
				// give up all reserved bump fees
				newTx := signer.CopyMsgTxNoWitness(msgTx)
				newTx.TxOut = newTx.TxOut[:2]
				return newTx
			}(),
		},
		{
			name: "should set new gas rate to 'gasRateCap'",
			feeBumper: &signer.CPFPFeeBumper{
				Tx:          btcutil.NewTx(msgTx),
				MinRelayFee: 0.00001,
				CCTXRate:    101, // > 100
				LiveRate:    120,
				TotalFees:   27213,
				TotalVSize:  579,
				AvgFeeRate:  47,
				Logger:      log.Logger,
			},
			additionalFees:  30687, // (100-47)*579
			expectedNewRate: 100,
			expectedNewTx: func() *wire.MsgTx {
				// deduct additional fees
				newTx := signer.CopyMsgTxNoWitness(msgTx)
				newTx.TxOut[2].Value -= 30687
				return newTx
			}(),
		},
		{
			name: "should fail if original tx has no reserved bump fees",
			feeBumper: &signer.CPFPFeeBumper{
				Tx: func() *btcutil.Tx {
					// remove the change output
					newTx := msgTx.Copy()
					newTx.TxOut = newTx.TxOut[:2]
					return btcutil.NewTx(newTx)
				}(),
			},
			errMsg: "original tx has no reserved bump fees",
		},
		{
			name: "should hold on RBF if additional fees is lower than min relay fees",
			feeBumper: &signer.CPFPFeeBumper{
				Tx:          btcutil.NewTx(msgTx),
				MinRelayFee: 0.00002, // min relay fee will be 579vB * 2 = 1158 sats
				CCTXRate:    6,
				LiveRate:    7,
				TotalFees:   2895,
				TotalVSize:  579,
				AvgFeeRate:  5,
				Logger:      log.Logger,
			},
			errMsg: "lower than min relay fees",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newTx, additionalFees, newRate, err := tt.feeBumper.BumpTxFee()
			if tt.errMsg != "" {
				require.Nil(t, newTx)
				require.Zero(t, additionalFees)
				require.ErrorContains(t, err, tt.errMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedNewTx, newTx)
				require.Equal(t, tt.additionalFees, additionalFees)
				require.Equal(t, tt.expectedNewRate, newRate)
			}
		})
	}
}

func Test_FetchFeeBumpInfo(t *testing.T) {
	liveRate := int64(12)

	tests := []struct {
		name         string
		tx           *btcutil.Tx
		liveRate     int64
		memplTxsInfo *mempoolTxsInfo
		expected     *signer.CPFPFeeBumper
		errMsg       string
	}{
		{
			name:     "should fetch fee bump info successfully",
			tx:       btcutil.NewTx(wire.NewMsgTx(wire.TxVersion)),
			liveRate: 12,
			memplTxsInfo: newMempoolTxsInfo(
				2,      // 2 stuck TSS txs
				0.0001, // total fees 0.0001 BTC
				1000,   // total vsize 1000
				10,     // average fee rate 10 sat/vB
			),
			expected: &signer.CPFPFeeBumper{
				Tx:         btcutil.NewTx(wire.NewMsgTx(wire.TxVersion)),
				LiveRate:   12,
				TotalTxs:   2,
				TotalFees:  10000,
				TotalVSize: 1000,
				AvgFeeRate: 10,
				Logger:     log.Logger,
			},
		},
		{
			name:     "should fail if unable to estimate smart fee",
			liveRate: 0,
			errMsg:   "GetEstimatedFeeRate failed",
		},
		{
			name:         "should fail if unable to fetch mempool txs info",
			liveRate:     12,
			tx:           btcutil.NewTx(wire.NewMsgTx(wire.TxVersion)),
			memplTxsInfo: nil,
			errMsg:       "unable to fetch mempool txs info",
		},
		{
			name:         "should fail on invalid total fees",
			liveRate:     12,
			tx:           btcutil.NewTx(wire.NewMsgTx(wire.TxVersion)),
			memplTxsInfo: newMempoolTxsInfo(2, 21000000.1, 1000, 10), // fee exceeds max BTC supply
			errMsg:       "cannot convert total fees",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			// mock RPC fee rate
			client := mocks.NewBitcoinClient(t)
			if tt.liveRate > 0 {
				client.On("GetEstimatedFeeRate", mock.Anything, mock.Anything, mock.Anything).Return(liveRate, nil)
			} else {
				client.On("GetEstimatedFeeRate", mock.Anything, mock.Anything, mock.Anything).Return(int64(0), errors.New("rpc error"))
			}

			// mock mempool txs information
			if tt.memplTxsInfo != nil {
				client.On("GetTotalMempoolParentsSizeNFees", mock.Anything, mock.Anything, mock.Anything).
					Return(tt.memplTxsInfo.totalTxs, tt.memplTxsInfo.totalFees, tt.memplTxsInfo.totalVSize, tt.memplTxsInfo.avgFeeRate, nil)
			} else {
				v := int64(0)
				client.On("GetTotalMempoolParentsSizeNFees", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(v, 0.0, v, v, errors.New("rpc error"))
			}

			// ACT
			bumper := &signer.CPFPFeeBumper{
				RPC:    client,
				Tx:     tt.tx,
				Logger: log.Logger,
			}
			err := bumper.FetchFeeBumpInfo()

			// ASSERT
			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
			} else {
				bumper.RPC = nil // ignore the RPC client
				require.NoError(t, err)
				require.Equal(t, tt.expected, bumper)
			}
		})
	}
}

func Test_CopyMsgTxNoWitness(t *testing.T) {
	t.Run("should copy tx msg without witness", func(t *testing.T) {
		chain := chains.BitcoinMainnet
		txid := "030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0"
		msgTx := testutils.LoadBTCMsgTx(t, TestDataDir, chain.ChainId, txid)

		// make a non-witness copy
		copyTx := signer.CopyMsgTxNoWitness(msgTx)

		// make another copy and clear witness data manually
		newTx := msgTx.Copy()
		for idx := range newTx.TxIn {
			newTx.TxIn[idx].Witness = wire.TxWitness{}
		}

		// check
		require.Equal(t, newTx, copyTx)
	})

	t.Run("should handle nil input", func(t *testing.T) {
		require.Panics(t, func() {
			signer.CopyMsgTxNoWitness(nil)
		}, "should panic on nil input")
	})
}
