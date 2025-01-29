package signer_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/testutils"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/signer"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

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
			inAmounts, err := signer.AddTxInputs(tx, tt.utxos)

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
	)
	signer := signer.New(
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
			inAmounts, err := signer.AddTxInputs(tx, utxos)
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
