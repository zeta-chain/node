package signer

import (
	"context"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
	"github.com/zeta-chain/node/zetaclient/testutils"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
)

func Test_SignRBFTx(t *testing.T) {
	// https://mempool.space/tx/030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0
	chain := chains.BitcoinMainnet
	txid := "030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0"
	msgTx := testutils.LoadBTCMsgTx(t, TestDataDir, chain.ChainId, txid)

	// inputs of the transaction
	type prevTx struct {
		hash   *chainhash.Hash
		vout   uint32
		amount int64
	}
	preTxs := []prevTx{
		{
			hash: hashFromTXID(
				t,
				"efca302a18bd8cebb3b8afef13e98ecaac47157755a62ab241ef3848140cfe92",
			), vout: 0, amount: 2147,
		},
		{
			hash: hashFromTXID(
				t,
				"efca302a18bd8cebb3b8afef13e98ecaac47157755a62ab241ef3848140cfe92",
			), vout: 2, amount: 28240703,
		},
		{
			hash: hashFromTXID(
				t,
				"3dc005eb0c1d393e717070ea84aa13e334a458a4fb7c7f9f98dbf8b231b5ceef",
			), vout: 0, amount: 10000,
		},
		{
			hash: hashFromTXID(
				t,
				"74c3aca825f3b21b82ee344d939c40d4c1e836a89c18abbd521bfa69f5f6e5d7",
			), vout: 0, amount: 10000,
		},
		{
			hash: hashFromTXID(
				t,
				"87264cef0e581f4aab3c99c53221bec3219686b48088d651a8cf8a98e4c2c5bf",
			), vout: 0, amount: 10000,
		},
		{
			hash: hashFromTXID(
				t,
				"5af24933973df03d96624ae1341d79a860e8dbc2ffc841420aa6710f3abc0074",
			), vout: 0, amount: 1200000,
		},
		{
			hash: hashFromTXID(
				t,
				"b85755938ac026b2d13e5fbacf015288f453712b4eb4a02d7e4c98ee76ada530",
			), vout: 0, amount: 9610000,
		},
	}

	// test cases
	tests := []struct {
		name       string
		chain      chains.Chain
		lastTx     *btcutil.Tx
		preTxs     []prevTx
		txData     OutboundData
		liveRate   uint64
		txsAndFees *client.MempoolTxsAndFees
		errMsg     string
		expectedTx *wire.MsgTx
	}{
		{
			name:     "should sign RBF tx successfully",
			chain:    chains.BitcoinMainnet,
			lastTx:   btcutil.NewTx(msgTx.Copy()),
			preTxs:   preTxs,
			txData:   mkTxData(t, 0.00001, "57"), // 57 sat/vB as cctx rate
			liveRate: 59,                         // 59 sat/vB
			txsAndFees: &client.MempoolTxsAndFees{
				TotalTxs:   1,     // 1 stuck tx
				TotalFees:  27213, // fees: 0.00027213 BTC
				TotalVSize: 579,   // size: 579 vByte
				AvgFeeRate: 47,    // rate: 47 sat/vB
			},
			expectedTx: func() *wire.MsgTx {
				// deduct additional fees
				newTx := CopyMsgTxNoWitness(msgTx)
				newTx.TxOut[2].Value -= 5790
				return newTx
			}(),
		},
		{
			name:   "should return error if fee rate is not bumped by zetacore yet",
			chain:  chains.BitcoinMainnet,
			lastTx: btcutil.NewTx(msgTx.Copy()),
			txData: mkTxData(t, 0.00001, ""), // empty gas priority fee, not bumped yet
			errMsg: "fee rate is not bumped by zetacore yet",
		},
		{
			name:       "should return error if unable to create fee bumper",
			chain:      chains.BitcoinMainnet,
			lastTx:     btcutil.NewTx(msgTx.Copy()),
			txData:     mkTxData(t, 0.00001, "57"),
			txsAndFees: nil, // no mempool txs info provided
			errMsg:     "NewCPFPFeeBumper failed",
		},
		{
			name:  "should return error if unable to bump tx fee",
			chain: chains.BitcoinMainnet,
			lastTx: func() *btcutil.Tx {
				txCopy := msgTx.Copy()
				txCopy.TxOut = txCopy.TxOut[:2] // remove reserved bump fees to cause error
				return btcutil.NewTx(txCopy)
			}(),
			txData:   mkTxData(t, 0.00001, "57"), // 57 sat/vB as cctx rate
			liveRate: 99,                         // 99 sat/vB is much higher than ccxt rate
			txsAndFees: &client.MempoolTxsAndFees{
				TotalTxs:   1,     // 1 stuck tx
				TotalFees:  27213, // fees: 0.00027213 BTC
				TotalVSize: 579,   // size: 579 vByte
				AvgFeeRate: 47,    // rate: 47 sat/vB
			},
			errMsg: "BumpTxFee failed",
		},
		{
			name:     "should return error if unable to get previous tx",
			chain:    chains.BitcoinMainnet,
			lastTx:   btcutil.NewTx(msgTx.Copy()),
			txData:   mkTxData(t, 0.00001, "57"), // 57 sat/vB as cctx rate
			preTxs:   nil,                        // no previous info provided
			liveRate: 59,                         // 59 sat/vB
			txsAndFees: &client.MempoolTxsAndFees{
				TotalTxs:   1,     // 1 stuck tx
				TotalFees:  27213, // fees: 0.00027213 BTC
				TotalVSize: 579,   // size: 579 vByte
				AvgFeeRate: 47,    // rate: 47 sat/vB
			},
			errMsg: "unable to get previous tx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			// setup signer
			s := newTestSuite(t, tt.chain)

			// mock RPC live fee rate
			if tt.liveRate > 0 {
				s.client.On("GetEstimatedFeeRate", mock.Anything, mock.Anything, mock.Anything).Return(tt.liveRate, nil)
			} else {
				s.client.On("GetEstimatedFeeRate", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(uint64(0), errors.New("rpc error"))
			}

			// mock mempool txs information
			if tt.txsAndFees != nil {
				s.client.On("GetMempoolTxsAndFees", mock.Anything, mock.Anything).Return(*tt.txsAndFees, nil)
			} else {
				s.client.On("GetMempoolTxsAndFees", mock.Anything, mock.Anything).Maybe().Return(client.MempoolTxsAndFees{}, nil)
			}

			// mock RPC transactions
			if tt.preTxs != nil {
				// mock first two inputs they belong to same tx
				mockMsg := wire.NewMsgTx(wire.TxVersion)
				mockMsg.TxOut = make([]*wire.TxOut, 3)
				for _, preTx := range tt.preTxs[:2] {
					mockMsg.TxOut[preTx.vout] = wire.NewTxOut(preTx.amount, []byte{})
				}
				s.client.On("GetRawTransaction", mock.Anything, tt.preTxs[0].hash).
					Maybe().
					Return(btcutil.NewTx(mockMsg), nil)

				// mock other inputs
				for _, preTx := range tt.preTxs[2:] {
					mockMsg := wire.NewMsgTx(wire.TxVersion)
					mockMsg.TxOut = make([]*wire.TxOut, 3)
					mockMsg.TxOut[preTx.vout] = wire.NewTxOut(preTx.amount, []byte{})

					s.client.On("GetRawTransaction", mock.Anything, preTx.hash).
						Maybe().
						Return(btcutil.NewTx(mockMsg), nil)
				}
			} else {
				s.client.On("GetRawTransaction", mock.Anything, mock.Anything).Maybe().Return(nil, errors.New("rpc error"))
			}

			// ACT
			// sign tx
			ctx := context.Background()
			newTx, err := s.SignRBFTx(ctx, &tt.txData, tt.lastTx)
			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
				return
			}

			// ASSERT
			require.NoError(t, err)

			// check tx signature
			for i := range newTx.TxIn {
				require.Len(t, newTx.TxIn[i].Witness, 2)
			}
		})
	}
}

// mkTxData creates a new outbound data for testing
func mkTxData(t *testing.T, minRelayFee float64, latestFeeRate string) OutboundData {
	net := &chaincfg.MainNetParams
	cctx := sample.CrossChainTx(t, "0x123")
	cctx.InboundParams.CoinType = coin.CoinType_Gas
	cctx.GetCurrentOutboundParam().GasPrice = "1"
	cctx.GetCurrentOutboundParam().GasPriorityFee = latestFeeRate
	cctx.GetCurrentOutboundParam().Receiver = sample.BTCAddressP2WPKH(t, sample.Rand(), net).String()
	cctx.GetCurrentOutboundParam().ReceiverChainId = chains.BitcoinMainnet.ChainId
	cctx.GetCurrentOutboundParam().Amount = sdkmath.NewUint(1e7) // 0.1 BTC

	txData, err := NewOutboundData(cctx, 1, minRelayFee, false, zerolog.Nop())
	require.NoError(t, err)
	return *txData
}
