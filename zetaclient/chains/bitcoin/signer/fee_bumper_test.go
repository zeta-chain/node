package signer_test

import (
	"testing"

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
				10,     // average fee rate 10 sat/vbyte
				"",     // no error
			),
			expected: &signer.CPFPFeeBumper{
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
	msgTx := testutils.LoadBTCMsgTx(t, TestDataDir, chain.ChainId, txid).Copy()

	// cleanMsgTx is a helper function to clean witness data
	cleanMsgTx := func(tx *wire.MsgTx) *wire.MsgTx {
		newTx := tx.Copy()
		for idx := range newTx.TxIn {
			newTx.TxIn[idx].Witness = wire.TxWitness{}
		}
		return newTx
	}

	tests := []struct {
		name           string
		feeBumper      *signer.CPFPFeeBumper
		errMsg         string
		additionalFees int64
		expectedTx     *wire.MsgTx
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
			additionalFees: 5790,
			expectedTx: func() *wire.MsgTx {
				// deduct additional fees
				newTx := cleanMsgTx(msgTx)
				newTx.TxOut[2].Value -= 5790
				return newTx
			}(),
		},
		{
			name: "should cover min relay fees",
			feeBumper: &signer.CPFPFeeBumper{
				Tx:          btcutil.NewTx(msgTx),
				MinRelayFee: 0.00002, // min relay fee will be 579vB * 2 = 1158 sats
				CCTXRate:    6,
				LiveRate:    8,
				TotalFees:   2895,
				TotalVSize:  579,
				AvgFeeRate:  5,
			},
			additionalFees: 1158,
			expectedTx: func() *wire.MsgTx {
				// deduct additional fees
				newTx := cleanMsgTx(msgTx)
				newTx.TxOut[2].Value -= 1158
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
			additionalFees: 5790 + constant.BTCWithdrawalDustAmount - 1, // 6789
			expectedTx: func() *wire.MsgTx {
				// give up all reserved bump fees
				newTx := cleanMsgTx(msgTx)
				newTx.TxOut = newTx.TxOut[:2]
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
				CCTXRate:   56, // 56 < 47 * 120%
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newTx, additionalFees, err := tt.feeBumper.BumpTxFee()
			if tt.errMsg != "" {
				require.Nil(t, newTx)
				require.Zero(t, additionalFees)
				require.ErrorContains(t, err, tt.errMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTx, newTx)
				require.Equal(t, tt.additionalFees, additionalFees)
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
				10,     // average fee rate 10 sat/vbyte
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

	return func(interfaces.BTCRPCClient, string) (int64, float64, int64, int64, error) {
		return totalTxs, totalFees, totalVSize, avgFeeRate, err
	}
}
