package signer_test

import (
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/signer"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

func Test_NewCPFPFeeBumper(t *testing.T) {
	tests := []struct {
		name                string
		chain               chains.Chain
		client              *mocks.BTCRPCClient
		tx                  *btcutil.Tx
		cctxRate            int64
		liveRate            float64
		minRelayFee         float64
		memplTxsInfoFetcher signer.MempoolTxsInfoFetcher
		errMsg              string
		expected            *signer.CPFPFeeBumper
	}{
		{
			chain:       chains.BitcoinMainnet,
			name:        "should create new CPFPFeeBumper successfully",
			client:      mocks.NewBTCRPCClient(t),
			tx:          btcutil.NewTx(wire.NewMsgTx(wire.TxVersion)),
			cctxRate:    10,
			liveRate:    0.00012,
			minRelayFee: 0.00001,
			memplTxsInfoFetcher: makeMempoolTxsInfoFetcher(
				2,      // 2 stuck TSS txs
				0.0001, // total fees 0.0001 BTC
				1000,   // total vsize 1000
				10,     // average fee rate 10 sat/vB
				"",     // no error
			),
			expected: &signer.CPFPFeeBumper{
				Chain:       chains.BitcoinMainnet,
				Tx:          btcutil.NewTx(wire.NewMsgTx(wire.TxVersion)),
				MinRelayFee: 0.00001,
				CCTXRate:    10,
				LiveRate:    12,
				TotalTxs:    2,
				TotalFees:   10000,
				TotalVSize:  1000,
				AvgFeeRate:  10,
			},
		},
		{
			chain:    chains.BitcoinMainnet,
			name:     "should fail when mempool txs info fetcher returns error",
			client:   mocks.NewBTCRPCClient(t),
			tx:       btcutil.NewTx(wire.NewMsgTx(wire.TxVersion)),
			liveRate: 0.00012,
			//memplTxsInfoFetcher: makeMempoolTxsInfoFetcher(0, 0.0, 0, 0, "rpc error"),
			memplTxsInfoFetcher: makeMempoolTxsInfoFetcher(
				2,      // 2 stuck TSS txs
				0.0001, // total fees 0.0001 BTC
				1000,   // total vsize 1000
				10,     // average fee rate 10 sat/vbyte
				"err",  // no error
			),
			errMsg: "unable to fetch mempool txs info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mock RPC fee rate
			tt.client.On("EstimateSmartFee", mock.Anything, mock.Anything).Return(&btcjson.EstimateSmartFeeResult{
				FeeRate: &tt.liveRate,
			}, nil)

			bumper, err := signer.NewCPFPFeeBumper(
				tt.chain,
				tt.client,
				tt.memplTxsInfoFetcher,
				tt.tx,
				tt.cctxRate,
				tt.minRelayFee,
				log.Logger,
			)
			if tt.errMsg != "" {
				require.Nil(t, bumper)
				require.ErrorContains(t, err, tt.errMsg)
			} else {
				bumper.Client = nil // ignore the RPC client
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
				CCTXRate:    57,
				LiveRate:    60,
				TotalFees:   27213,
				TotalVSize:  579,
				AvgFeeRate:  47,
			},
			additionalFees:  5790,
			expectedNewRate: 57,
			expectedNewTx: func() *wire.MsgTx {
				// deduct additional fees
				newTx := signer.CopyMsgTxNoWitness(msgTx)
				newTx.TxOut[2].Value -= 5790
				return newTx
			}(),
		},
		{
			name: "should give up all reserved bump fees",
			feeBumper: &signer.CPFPFeeBumper{
				Tx: func() *btcutil.Tx {
					// modify reserved bump fees to barely cover bump fees
					newTx := msgTx.Copy()
					newTx.TxOut[2].Value = 5790 + constant.BTCWithdrawalDustAmount - 1
					return btcutil.NewTx(newTx)
				}(),
				MinRelayFee: 0.00001,
				CCTXRate:    57,
				LiveRate:    60,
				TotalFees:   27213,
				TotalVSize:  579,
				AvgFeeRate:  47,
			},
			additionalFees:  5790 + constant.BTCWithdrawalDustAmount - 1, // 6789
			expectedNewRate: 59,                                          // (27213 + 6789) / 579 ≈ 59
			expectedNewTx: func() *wire.MsgTx {
				// give up all reserved bump fees
				newTx := signer.CopyMsgTxNoWitness(msgTx)
				newTx.TxOut = newTx.TxOut[:2]
				return newTx
			}(),
		},
		{
			name: "should cap new gas rate to 'gasRateCap'",
			feeBumper: &signer.CPFPFeeBumper{
				Tx:          btcutil.NewTx(msgTx),
				MinRelayFee: 0.00001,
				CCTXRate:    101, // > 100
				LiveRate:    120,
				TotalFees:   27213,
				TotalVSize:  579,
				AvgFeeRate:  47,
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
			name: "should hold on RBF if CCTX rate is lower than minimum bumpeed rate",
			feeBumper: &signer.CPFPFeeBumper{
				Tx:         btcutil.NewTx(msgTx),
				CCTXRate:   55, // 56 < 47 * 120%
				AvgFeeRate: 47,
			},
			errMsg: "lower than the min bumped rate",
		},
		{
			name: "should hold on RBF if live rate is much higher than CCTX rate",
			feeBumper: &signer.CPFPFeeBumper{
				Tx:         btcutil.NewTx(msgTx),
				CCTXRate:   57,
				LiveRate:   70, // 70 > 57 * 120%
				AvgFeeRate: 47,
			},
			errMsg: "much higher than the cctx rate",
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
	liveRate := 0.00012
	mockClient := mocks.NewBTCRPCClient(t)
	mockClient.On("EstimateSmartFee", mock.Anything, mock.Anything).Return(&btcjson.EstimateSmartFeeResult{
		FeeRate: &liveRate,
	}, nil)

	tests := []struct {
		name                string
		client              *mocks.BTCRPCClient
		tx                  *btcutil.Tx
		memplTxsInfoFetcher signer.MempoolTxsInfoFetcher
		expected            *signer.CPFPFeeBumper
		errMsg              string
	}{
		{
			name:   "should fetch fee bump info successfully",
			client: mockClient,
			tx:     btcutil.NewTx(wire.NewMsgTx(wire.TxVersion)),
			memplTxsInfoFetcher: makeMempoolTxsInfoFetcher(
				2,      // 2 stuck TSS txs
				0.0001, // total fees 0.0001 BTC
				1000,   // total vsize 1000
				10,     // average fee rate 10 sat/vB
				"",     // no error
			),
			expected: &signer.CPFPFeeBumper{
				Tx:         btcutil.NewTx(wire.NewMsgTx(wire.TxVersion)),
				LiveRate:   12,
				TotalTxs:   2,
				TotalFees:  10000,
				TotalVSize: 1000,
				AvgFeeRate: 10,
			},
		},
		{
			name: "should fail if unable to estimate smart fee",
			client: func() *mocks.BTCRPCClient {
				client := mocks.NewBTCRPCClient(t)
				client.On("EstimateSmartFee", mock.Anything, mock.Anything).Return(nil, errors.New("rpc error"))
				return client
			}(),
			errMsg: "GetEstimatedFeeRate failed",
		},
		{
			name:                "should fail if unable to fetch mempool txs info",
			client:              mockClient,
			tx:                  btcutil.NewTx(wire.NewMsgTx(wire.TxVersion)),
			memplTxsInfoFetcher: makeMempoolTxsInfoFetcher(0, 0.0, 0, 0, "rpc error"),
			errMsg:              "unable to fetch mempool txs info",
		},
		{
			name:                "should fail on invalid total fees",
			client:              mockClient,
			tx:                  btcutil.NewTx(wire.NewMsgTx(wire.TxVersion)),
			memplTxsInfoFetcher: makeMempoolTxsInfoFetcher(2, 21000000.1, 1000, 10, ""), // fee exceeds max BTC supply
			errMsg:              "cannot convert total fees",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bumper := &signer.CPFPFeeBumper{
				Client: tt.client,
				Tx:     tt.tx,
			}
			err := bumper.FetchFeeBumpInfo(tt.memplTxsInfoFetcher, log.Logger)

			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
			} else {
				bumper.Client = nil // ignore the RPC client
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

// makeMempoolTxsInfoFetcher is a helper function to create a mock MempoolTxsInfoFetcher
func makeMempoolTxsInfoFetcher(
	totalTxs int64,
	totalFees float64,
	totalVSize int64,
	avgFeeRate int64,
	errMsg string,
) signer.MempoolTxsInfoFetcher {
	var err error
	if errMsg != "" {
		err = errors.New(errMsg)
	}

	return func(interfaces.BTCRPCClient, string, time.Duration) (int64, float64, int64, int64, error) {
		return totalTxs, totalFees, totalVSize, avgFeeRate, err
	}
}
