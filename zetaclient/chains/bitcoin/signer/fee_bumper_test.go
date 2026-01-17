package signer

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
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

func Test_NewCPFPFeeBumper(t *testing.T) {
	tests := []struct {
		name        string
		chain       chains.Chain
		client      *mocks.BitcoinClient
		tx          *btcutil.Tx
		cctxRate    uint64
		liveRate    uint64
		minRelayFee float64
		txsAndFees  *client.MempoolTxsAndFees
		errMsg      string
		expected    *CPFPFeeBumper
	}{
		{
			chain:       chains.BitcoinMainnet,
			name:        "should create new CPFPFeeBumper successfully",
			client:      mocks.NewBitcoinClient(t),
			tx:          btcutil.NewTx(wire.NewMsgTx(wire.TxVersion)),
			cctxRate:    10,
			liveRate:    12,
			minRelayFee: 0.00001,
			txsAndFees: &client.MempoolTxsAndFees{
				TotalTxs:   2,     // 2 stuck TSS txs
				TotalFees:  10000, // total fees 0.0001 BTC
				TotalVSize: 1000,  // total vsize 1000
				AvgFeeRate: 10,    // average fee rate 10 sat/vB
			},
			expected: &CPFPFeeBumper{
				ctx:         context.Background(),
				chain:       chains.BitcoinMainnet,
				tx:          btcutil.NewTx(wire.NewMsgTx(wire.TxVersion)),
				minRelayFee: 0.00001,
				cctxRate:    10,
				liveRate:    12,
				txsAndFees: client.MempoolTxsAndFees{
					TotalTxs:   2,
					TotalFees:  10000,
					TotalVSize: 1000,
					AvgFeeRate: 10,
				},
				logger: log.Logger,
			},
		},
		{
			chain:    chains.BitcoinMainnet,
			name:     "should fail if unable to estimate smart fee",
			client:   mocks.NewBitcoinClient(t),
			tx:       btcutil.NewTx(wire.NewMsgTx(wire.TxVersion)),
			liveRate: 0,
			errMsg:   "GetEstimatedFeeRate failed",
		},
		{
			chain:      chains.BitcoinMainnet,
			name:       "should fail if unable to fetch mempool txs info",
			client:     mocks.NewBitcoinClient(t),
			tx:         btcutil.NewTx(wire.NewMsgTx(wire.TxVersion)),
			liveRate:   12,
			txsAndFees: nil,
			errMsg:     "unable to fetch mempool txs info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			// mock RPC fee rate
			if tt.liveRate > 0 {
				tt.client.On("GetEstimatedFeeRate", mock.Anything, mock.Anything).Return(tt.liveRate, nil)
			} else {
				tt.client.On("GetEstimatedFeeRate", mock.Anything, mock.Anything).Return(uint64(0), errors.New("rpc error"))
			}

			// mock mempool txs information
			if tt.txsAndFees != nil {
				tt.client.On("GetMempoolTxsAndFees", mock.Anything, mock.Anything).Maybe().Return(*tt.txsAndFees, nil)
			} else {
				tt.client.On("GetMempoolTxsAndFees", mock.Anything, mock.Anything).Maybe().Return(client.MempoolTxsAndFees{}, errors.New("rpc error"))
			}

			// ACT
			bumper, err := NewCPFPFeeBumper(
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
				bumper.bitcoinClient = nil // ignore the RPC
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
		name      string
		feeBumper *CPFPFeeBumper
		expected  BumpResult
		errMsg    string
	}{
		{
			name: "should bump tx fee successfully",
			feeBumper: &CPFPFeeBumper{
				tx:          btcutil.NewTx(msgTx),
				minRelayFee: 0.00001,
				cctxRate:    55,
				liveRate:    67,
				txsAndFees: client.MempoolTxsAndFees{
					TotalFees:  27213,
					TotalVSize: 579,
					AvgFeeRate: 47,
				},
				logger: log.Logger,
			},
			expected: BumpResult{
				NewTx: func() *wire.MsgTx {
					// deduct additional fees
					newTx := CopyMsgTxNoWitness(msgTx)
					newTx.TxOut[2].Value -= 4632
					return newTx
				}(),
				AdditionalFees: 4632, // 579*55 - 27213
				NewFeeRate:     55,
			},
		},
		{
			name: "should give up all reserved bump fees",
			feeBumper: &CPFPFeeBumper{
				tx: func() *btcutil.Tx {
					// modify reserved bump fees to barely cover bump fees
					newTx := msgTx.Copy()
					newTx.TxOut[2].Value = 57*579 - 27213 + constant.BTCWithdrawalDustAmount - 1 // 6789
					return btcutil.NewTx(newTx)
				}(),
				minRelayFee: 0.00001,
				cctxRate:    57,
				liveRate:    67,
				txsAndFees: client.MempoolTxsAndFees{
					TotalFees:  27213,
					TotalVSize: 579,
					AvgFeeRate: 47,
				},
				logger: log.Logger,
			},
			expected: BumpResult{
				NewTx: func() *wire.MsgTx {
					// give up all reserved bump fees
					newTx := CopyMsgTxNoWitness(msgTx)
					newTx.TxOut = newTx.TxOut[:2]
					return newTx
				}(),
				AdditionalFees: 6789, // same as the reserved value in 2nd output
				NewFeeRate:     59,   // (27213 + 6789) / 579 ≈ 59
			},
		},
		{
			name: "should set new gas rate to 'gasRateCap'",
			feeBumper: &CPFPFeeBumper{
				tx:          btcutil.NewTx(msgTx),
				minRelayFee: 0.00001,
				cctxRate:    101, // > 100
				liveRate:    120,
				txsAndFees: client.MempoolTxsAndFees{
					TotalFees:  27213,
					TotalVSize: 579,
					AvgFeeRate: 47,
				},
				logger: log.Logger,
			},
			expected: BumpResult{
				NewTx: func() *wire.MsgTx {
					// deduct additional fees
					newTx := CopyMsgTxNoWitness(msgTx)
					newTx.TxOut[2].Value -= 30687
					return newTx
				}(),
				AdditionalFees: 30687, // (100-47)*579
				NewFeeRate:     100,
			},
		},
		{
			name: "should fail if original tx has no reserved bump fees",
			feeBumper: &CPFPFeeBumper{
				tx: func() *btcutil.Tx {
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
			feeBumper: &CPFPFeeBumper{
				tx:          btcutil.NewTx(msgTx),
				minRelayFee: 0.00002, // min relay fee will be 579vB * 2 = 1158 sats
				cctxRate:    6,
				liveRate:    7,
				txsAndFees: client.MempoolTxsAndFees{
					TotalFees:  2895,
					TotalVSize: 579,
					AvgFeeRate: 5,
				},
				logger: log.Logger,
			},
			errMsg: "lower than min relay fees",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.feeBumper.BumpTxFee()
			if tt.errMsg != "" {
				require.Nil(t, result.NewTx)
				require.ErrorContains(t, err, tt.errMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
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
		copyTx := CopyMsgTxNoWitness(msgTx)

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
			CopyMsgTxNoWitness(nil)
		}, "should panic on nil input")
	})
}
