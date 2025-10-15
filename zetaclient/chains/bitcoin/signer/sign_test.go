package signer

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/mode"
	"github.com/zeta-chain/node/zetaclient/testutils"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

func Test_SignWithdrawTx(t *testing.T) {
	net := &chaincfg.MainNetParams

	// make sample cctx
	mkCCTX := func(t *testing.T) *crosschaintypes.CrossChainTx {
		cctx := sample.CrossChainTx(t, "0x123")
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		cctx.GetCurrentOutboundParam().GasPrice = "10"
		cctx.GetCurrentOutboundParam().Receiver = sample.BTCAddressP2WPKH(t, sample.Rand(), net).String()
		cctx.GetCurrentOutboundParam().ReceiverChainId = chains.BitcoinMainnet.ChainId
		cctx.GetCurrentOutboundParam().Amount = sdkmath.NewUint(1e7) // 0.1 BTC
		cctx.GetCurrentOutboundParam().CallOptions = &crosschaintypes.CallOptions{GasLimit: 254}
		cctx.GetCurrentOutboundParam().TssNonce = 0
		return cctx
	}

	// helper function to create tx data
	mkTxData := func(height uint64, minRelayFee float64) OutboundData {
		cctx := mkCCTX(t)
		txData, err := NewOutboundData(cctx, height, minRelayFee, false, zerolog.Nop())
		require.NoError(t, err)
		return *txData
	}

	tests := []struct {
		name           string
		chain          chains.Chain
		txData         OutboundData
		failFetchUTXOs bool
		failSignTx     bool
		fail           bool
	}{
		{
			name:   "should sign withdraw tx successfully",
			chain:  chains.BitcoinMainnet,
			txData: mkTxData(101, 0.00001),
		},
		{
			name:           "should fail if no UTXOs fetched due to RPC error",
			chain:          chains.BitcoinMainnet,
			txData:         mkTxData(101, 0.00001),
			failFetchUTXOs: true,
			fail:           true,
		},
		{
			name:       "should fail if TSS keysign fails",
			chain:      chains.BitcoinMainnet,
			txData:     mkTxData(101, 0.00001),
			failSignTx: true,
			fail:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			// setup signer
			s := newTestSuite(t, tt.chain)
			btcAddress, err := s.TSS().PubKey().AddressBTC(tt.chain.ChainId)
			require.NoError(t, err)
			tssAddress := btcAddress.EncodeAddress()

			// mock up pending nonces
			pendingNonces := observertypes.PendingNonces{}
			s.zetacoreClient.On("GetPendingNoncesByChain", mock.Anything, mock.Anything).
				Maybe().
				Return(pendingNonces, nil)

			// mock up utxos
			utxos := []btcjson.ListUnspentResult{}
			utxos = append(utxos, btcjson.ListUnspentResult{Address: tssAddress, Amount: 1.0, Confirmations: 1})
			if !tt.failFetchUTXOs {
				s.client.On("ListUnspentMinMaxAddresses", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(utxos, nil)
			} else {
				s.client.On("ListUnspentMinMaxAddresses", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("rpc error"))
			}

			// mock up TSS SignBatch error
			if tt.failSignTx {
				s.tss.Pause()
			}

			// ACT
			// sign withdraw tx
			ctx := context.Background()
			tx, err := s.SignWithdrawTx(ctx, &tt.txData, s.observer)

			// ASSERT
			if tt.fail {
				require.Error(t, err)
				require.Nil(t, tx)
				return
			}
			require.NoError(t, err)

			// check tx signature
			for i := range tx.TxIn {
				require.Len(t, tx.TxIn[i].Witness, 2)
			}
		})
	}
}

func Test_AddTxInputs(t *testing.T) {
	r := sample.Rand()
	net := &chaincfg.MainNetParams

	tests := []struct {
		name            string
		utxos           []btcjson.ListUnspentResult
		expectedAmounts []int64
		fail            bool
	}{
		{
			name: "should add tx inputs successfully",
			utxos: []btcjson.ListUnspentResult{
				{
					TxID:    sample.BtcHash().String(),
					Vout:    0,
					Address: sample.BTCAddressP2WPKH(t, r, net).String(),
					Amount:  0.1,
				},
				{
					TxID:    sample.BtcHash().String(),
					Vout:    1,
					Address: sample.BTCAddressP2WPKH(t, r, net).String(),
					Amount:  0.2,
				},
			},
			expectedAmounts: []int64{10000000, 20000000},
		},
		{
			name: "should fail on invalid txid",
			utxos: []btcjson.ListUnspentResult{
				{
					TxID: "invalid txid",
				},
			},
			fail: true,
		},
		{
			name: "should fail on invalid amount",
			utxos: []btcjson.ListUnspentResult{
				{
					TxID:   sample.BtcHash().String(),
					Amount: -0.1,
				},
			},
			fail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create tx msg and add inputs
			tx := wire.NewMsgTx(wire.TxVersion)
			inAmounts, err := AddTxInputs(tx, tt.utxos)

			// assert
			if tt.fail {
				require.Error(t, err)
				require.Nil(t, inAmounts)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAmounts, inAmounts)
			}
		})
	}
}

func Test_AddWithdrawTxOutputs(t *testing.T) {
	// Create test signer and receiver address
	baseSigner := base.NewSigner(
		chains.BitcoinMainnet,
		mocks.NewTSS(t).FakePubKey(testutils.TSSPubKeyMainnet),
		base.DefaultLogger(),
		mode.StandardMode,
	)
	signer := New(
		baseSigner,
		mocks.NewBitcoinClient(t),
	)

	// tss address and script
	tssAddr, err := signer.TSS().PubKey().AddressBTC(chains.BitcoinMainnet.ChainId)
	require.NoError(t, err)
	tssScript, err := txscript.PayToAddrScript(tssAddr)
	require.NoError(t, err)
	fmt.Printf("tss address: %s", tssAddr.EncodeAddress())

	// receiver addresses
	receiver := "bc1qaxf82vyzy8y80v000e7t64gpten7gawewzu42y"
	to, err := chains.DecodeBtcAddress(receiver, chains.BitcoinMainnet.ChainId)
	require.NoError(t, err)
	toScript, err := txscript.PayToAddrScript(to)
	require.NoError(t, err)

	// test cases
	tests := []struct {
		name       string
		tx         *wire.MsgTx
		to         btcutil.Address
		total      float64
		amountSats int64
		nonceMark  int64
		fees       int64
		cancelTx   bool
		fail       bool
		message    string
		txout      []*wire.TxOut
	}{
		{
			name:       "should add outputs successfully",
			tx:         wire.NewMsgTx(wire.TxVersion),
			to:         to,
			total:      1.00012000,
			amountSats: 20000000,
			nonceMark:  10000,
			fees:       2000,
			fail:       false,
			txout: []*wire.TxOut{
				{Value: 10000, PkScript: tssScript},
				{Value: 20000000, PkScript: toScript},
				{Value: 80000000, PkScript: tssScript},
			},
		},
		{
			name:       "should add outputs without change successfully",
			tx:         wire.NewMsgTx(wire.TxVersion),
			to:         to,
			total:      0.20012000,
			amountSats: 20000000,
			nonceMark:  10000,
			fees:       2000,
			fail:       false,
			txout: []*wire.TxOut{
				{Value: 10000, PkScript: tssScript},
				{Value: 20000000, PkScript: toScript},
			},
		},
		{
			name:       "should cancel tx successfully",
			tx:         wire.NewMsgTx(wire.TxVersion),
			to:         to,
			total:      1.00012000,
			amountSats: 20000000,
			nonceMark:  10000,
			fees:       2000,
			cancelTx:   true,
			fail:       false,
			txout: []*wire.TxOut{
				{Value: 10000, PkScript: tssScript},
				{Value: 100000000, PkScript: tssScript},
			},
		},
		{
			name:       "should fail when total < amount",
			tx:         wire.NewMsgTx(wire.TxVersion),
			to:         to,
			total:      0.00012000,
			amountSats: 20000000,
			fail:       true,
		},
		{
			name:       "should fail when total < fees + amount + nonce",
			tx:         wire.NewMsgTx(wire.TxVersion),
			to:         to,
			total:      0.20011000,
			amountSats: 20000000,
			nonceMark:  10000,
			fees:       2000,
			fail:       true,
			message:    "remainder value is negative",
		},
		{
			name:       "should not produce duplicate nonce mark",
			tx:         wire.NewMsgTx(wire.TxVersion),
			to:         to,
			total:      0.20022000, //  0.2 + fee + nonceMark * 2
			amountSats: 20000000,
			nonceMark:  10000,
			fees:       2000,
			fail:       false,
			txout: []*wire.TxOut{
				{Value: 10000, PkScript: tssScript},
				{Value: 20000000, PkScript: toScript},
				{Value: 9999, PkScript: tssScript}, // nonceMark - 1
			},
		},
		{
			name:       "should not produce dust change to TSS self",
			tx:         wire.NewMsgTx(wire.TxVersion),
			to:         to,
			total:      0.20012999, // 0.2 + fee + nonceMark
			amountSats: 20000000,
			nonceMark:  10000,
			fees:       2000,
			fail:       false,
			txout: []*wire.TxOut{ // 3rd output 999 is dust and removed
				{Value: 10000, PkScript: tssScript},
				{Value: 20000000, PkScript: toScript},
			},
		},
		{
			name:       "should fail on invalid to address",
			tx:         wire.NewMsgTx(wire.TxVersion),
			to:         nil,
			total:      1.00012000,
			amountSats: 20000000,
			nonceMark:  10000,
			fees:       2000,
			fail:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := signer.AddWithdrawTxOutputs(
				tt.tx,
				tt.to,
				tt.total,
				tt.amountSats,
				tt.nonceMark,
				tt.fees,
				tt.cancelTx,
			)
			if tt.fail {
				require.Error(t, err)
				if tt.message != "" {
					require.ErrorContains(t, err, tt.message)
				}
				return
			} else {
				require.NoError(t, err)
				require.True(t, reflect.DeepEqual(tt.txout, tt.tx.TxOut))
			}
		})
	}
}

func Test_SignTx(t *testing.T) {
	tests := []struct {
		name    string
		chain   chains.Chain
		net     *chaincfg.Params
		inputs  []float64
		outputs []int64
		height  uint64
		nonce   uint64
	}{
		{
			name:  "should sign tx successfully",
			chain: chains.BitcoinMainnet,
			net:   &chaincfg.MainNetParams,
			inputs: []float64{
				0.0001,
				0.0002,
			},
			outputs: []int64{
				5000,
				20000,
			},
			nonce: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup signer
			s := newTestSuite(t, tt.chain)
			address, err := s.TSS().PubKey().AddressBTC(tt.chain.ChainId)
			require.NoError(t, err)

			// create tx msg
			tx := wire.NewMsgTx(wire.TxVersion)

			// add inputs
			utxos := []btcjson.ListUnspentResult{}
			for i, amount := range tt.inputs {
				utxos = append(utxos, btcjson.ListUnspentResult{
					TxID:    sample.BtcHash().String(),
					Vout:    uint32(i),
					Address: address.EncodeAddress(),
					Amount:  amount,
				})
			}
			inAmounts, err := AddTxInputs(tx, utxos)
			require.NoError(t, err)
			require.Len(t, inAmounts, len(tt.inputs))

			// add outputs
			r := sample.Rand()
			for _, amount := range tt.outputs {
				pkScript := sample.BTCAddressP2WPKHScript(t, r, tt.net)
				tx.AddTxOut(wire.NewTxOut(amount, pkScript))
			}

			// sign tx
			ctx := context.Background()
			err = s.SignTx(ctx, tx, inAmounts, tt.height, tt.nonce)
			require.NoError(t, err)

			// check tx signature
			for i := range tx.TxIn {
				require.Len(t, tx.TxIn[i].Witness, 2)
			}
		})
	}
}
