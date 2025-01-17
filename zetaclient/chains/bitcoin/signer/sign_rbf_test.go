package signer_test

import (
	"context"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/signer"
	"github.com/zeta-chain/node/zetaclient/testutils"

	"github.com/zeta-chain/node/pkg/chains"
)

func Test_SignRBFTx(t *testing.T) {
	// https://mempool.space/tx/030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0
	chain := chains.BitcoinMainnet
	nonce := uint64(148)
	cctx := testutils.LoadCctxByNonce(t, chain.ChainId, nonce)
	txid := "030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0"
	msgTx := testutils.LoadBTCMsgTx(t, TestDataDir, chain.ChainId, txid)

	// inputs
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
		name         string
		chain        chains.Chain
		cctx         *crosschaintypes.CrossChainTx
		lastTx       *btcutil.Tx
		preTxs       []prevTx
		minRelayFee  float64
		cctxRate     string
		liveRate     float64
		memplTxsInfo *mempoolTxsInfo
		errMsg       string
		expectedTx   *wire.MsgTx
	}{
		{
			name:        "should sign RBF tx successfully",
			chain:       chains.BitcoinMainnet,
			cctx:        cctx,
			lastTx:      btcutil.NewTx(msgTx.Copy()),
			preTxs:      preTxs,
			minRelayFee: 0.00001,
			cctxRate:    "57",
			liveRate:    0.00059, // 59 sat/vB
			memplTxsInfo: newMempoolTxsInfo(
				1,          // 1 stuck tx
				0.00027213, // fees: 0.00027213 BTC
				579,        // size: 579 vByte
				47,         // rate: 47 sat/vB
			),
			expectedTx: func() *wire.MsgTx {
				// deduct additional fees
				newTx := signer.CopyMsgTxNoWitness(msgTx)
				newTx.TxOut[2].Value -= 5790
				return newTx
			}(),
		},
		{
			name:        "should return error if latest fee rate is not available",
			chain:       chains.BitcoinMainnet,
			cctx:        cctx,
			lastTx:      btcutil.NewTx(msgTx.Copy()),
			minRelayFee: 0.00001,
			cctxRate:    "",
			errMsg:      "invalid fee rate",
		},
		{
			name:         "should return error if unable to create fee bumper",
			chain:        chains.BitcoinMainnet,
			cctx:         cctx,
			lastTx:       btcutil.NewTx(msgTx.Copy()),
			minRelayFee:  0.00001,
			cctxRate:     "57",
			memplTxsInfo: nil,
			errMsg:       "NewCPFPFeeBumper failed",
		},
		{
			name:        "should return error if live rate is too high",
			chain:       chains.BitcoinMainnet,
			cctx:        cctx,
			lastTx:      btcutil.NewTx(msgTx.Copy()),
			minRelayFee: 0.00001,
			cctxRate:    "57",
			liveRate:    0.00099, // 99 sat/vB is much higher than ccxt rate
			memplTxsInfo: newMempoolTxsInfo(
				1,          // 1 stuck tx
				0.00027213, // fees: 0.00027213 BTC
				579,        // size: 579 vByte
				47,         // rate: 47 sat/vB
			),
			errMsg: "BumpTxFee failed",
		},
		{
			name:        "should return error if live rate is too high",
			chain:       chains.BitcoinMainnet,
			cctx:        cctx,
			lastTx:      btcutil.NewTx(msgTx.Copy()),
			minRelayFee: 0.00001,
			cctxRate:    "57",
			liveRate:    0.00059, // 59 sat/vB
			memplTxsInfo: newMempoolTxsInfo(
				1,          // 1 stuck tx
				0.00027213, // fees: 0.00027213 BTC
				579,        // size: 579 vByte
				47,         // rate: 47 sat/vB
			),
			errMsg: "unable to get previous tx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup signer
			s := newTestSuite(t, tt.chain)

			// mock cctx rate
			tt.cctx.GetCurrentOutboundParam().GasPriorityFee = tt.cctxRate

			// mock RPC live fee rate
			if tt.liveRate > 0 {
				s.client.On("EstimateSmartFee", mock.Anything, mock.Anything).
					Maybe().
					Return(&btcjson.EstimateSmartFeeResult{
						FeeRate: &tt.liveRate,
					}, nil)
			} else {
				s.client.On("EstimateSmartFee", mock.Anything, mock.Anything).Maybe().Return(nil, errors.New("rpc error"))
			}

			// mock mempool txs information
			if tt.memplTxsInfo != nil {
				s.client.On("GetTotalMempoolParentsSizeNFees", mock.Anything, mock.Anything, mock.Anything).
					Return(tt.memplTxsInfo.totalTxs, tt.memplTxsInfo.totalFees, tt.memplTxsInfo.totalVSize, tt.memplTxsInfo.avgFeeRate, nil)
			} else {
				s.client.On("GetTotalMempoolParentsSizeNFees", mock.Anything, mock.Anything, mock.Anything).Return(0, 0.0, 0, 0, "rpc error")
			}

			// mock RPC transactions
			if tt.preTxs != nil {
				// mock first two inputs they belong to same tx
				mockMsg := wire.NewMsgTx(wire.TxVersion)
				mockMsg.TxOut = make([]*wire.TxOut, 3)
				for _, preTx := range tt.preTxs[:2] {
					mockMsg.TxOut[preTx.vout] = wire.NewTxOut(preTx.amount, []byte{})
				}
				s.client.On("GetRawTransaction", tt.preTxs[0].hash).Maybe().Return(btcutil.NewTx(mockMsg), nil)

				// mock other inputs
				for _, preTx := range tt.preTxs[2:] {
					mockMsg := wire.NewMsgTx(wire.TxVersion)
					mockMsg.TxOut = make([]*wire.TxOut, 3)
					mockMsg.TxOut[preTx.vout] = wire.NewTxOut(preTx.amount, []byte{})

					s.client.On("GetRawTransaction", preTx.hash).Maybe().Return(btcutil.NewTx(mockMsg), nil)
				}
			} else {
				s.client.On("GetRawTransaction", mock.Anything).Maybe().Return(nil, errors.New("rpc error"))
			}

			// sign tx
			ctx := context.Background()
			newTx, err := s.SignRBFTx(ctx, tt.cctx, 1, tt.lastTx, tt.minRelayFee)
			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
				return
			}
			require.NoError(t, err)

			// check tx signature
			for i := range newTx.TxIn {
				require.Len(t, newTx.TxIn[i].Witness, 2)
			}
		})
	}
}

func hashFromTXID(t *testing.T, txid string) *chainhash.Hash {
	h, err := chainhash.NewHashFromStr(txid)
	require.NoError(t, err)
	return h
}
