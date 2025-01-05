package signer

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/zetaclient/testutils"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

func TestAddWithdrawTxOutputs(t *testing.T) {
	// Create test signer and receiver address
	signer := NewSigner(
		chains.BitcoinMainnet,
		mocks.NewBTCRPCClient(t),
		mocks.NewTSS(t).FakePubKey(testutils.TSSPubKeyMainnet),
		base.DefaultLogger(),
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
		name      string
		tx        *wire.MsgTx
		to        btcutil.Address
		total     float64
		amount    float64
		nonceMark int64
		fees      int64
		cancelTx  bool
		fail      bool
		message   string
		txout     []*wire.TxOut
	}{
		{
			name:      "should add outputs successfully",
			tx:        wire.NewMsgTx(wire.TxVersion),
			to:        to,
			total:     1.00012000,
			amount:    0.2,
			nonceMark: 10000,
			fees:      2000,
			fail:      false,
			txout: []*wire.TxOut{
				{Value: 10000, PkScript: tssScript},
				{Value: 20000000, PkScript: toScript},
				{Value: 80000000, PkScript: tssScript},
			},
		},
		{
			name:      "should add outputs without change successfully",
			tx:        wire.NewMsgTx(wire.TxVersion),
			to:        to,
			total:     0.20012000,
			amount:    0.2,
			nonceMark: 10000,
			fees:      2000,
			fail:      false,
			txout: []*wire.TxOut{
				{Value: 10000, PkScript: tssScript},
				{Value: 20000000, PkScript: toScript},
			},
		},
		{
			name:      "should cancel tx successfully",
			tx:        wire.NewMsgTx(wire.TxVersion),
			to:        to,
			total:     1.00012000,
			amount:    0.2,
			nonceMark: 10000,
			fees:      2000,
			cancelTx:  true,
			fail:      false,
			txout: []*wire.TxOut{
				{Value: 10000, PkScript: tssScript},
				{Value: 100000000, PkScript: tssScript},
			},
		},
		{
			name:   "should fail on invalid amount",
			tx:     wire.NewMsgTx(wire.TxVersion),
			to:     to,
			total:  1.00012000,
			amount: -0.5,
			fail:   true,
		},
		{
			name:   "should fail when total < amount",
			tx:     wire.NewMsgTx(wire.TxVersion),
			to:     to,
			total:  0.00012000,
			amount: 0.2,
			fail:   true,
		},
		{
			name:      "should fail when total < fees + amount + nonce",
			tx:        wire.NewMsgTx(wire.TxVersion),
			to:        to,
			total:     0.20011000,
			amount:    0.2,
			nonceMark: 10000,
			fees:      2000,
			fail:      true,
			message:   "remainder value is negative",
		},
		{
			name:      "should not produce duplicate nonce mark",
			tx:        wire.NewMsgTx(wire.TxVersion),
			to:        to,
			total:     0.20022000, //  0.2 + fee + nonceMark * 2
			amount:    0.2,
			nonceMark: 10000,
			fees:      2000,
			fail:      false,
			txout: []*wire.TxOut{
				{Value: 10000, PkScript: tssScript},
				{Value: 20000000, PkScript: toScript},
				{Value: 9999, PkScript: tssScript}, // nonceMark - 1
			},
		},
		{
			name:      "should not produce dust change to TSS self",
			tx:        wire.NewMsgTx(wire.TxVersion),
			to:        to,
			total:     0.20012999, // 0.2 + fee + nonceMark
			amount:    0.2,
			nonceMark: 10000,
			fees:      2000,
			fail:      false,
			txout: []*wire.TxOut{ // 3rd output 999 is dust and removed
				{Value: 10000, PkScript: tssScript},
				{Value: 20000000, PkScript: toScript},
			},
		},
		{
			name:      "should fail on invalid to address",
			tx:        wire.NewMsgTx(wire.TxVersion),
			to:        nil,
			total:     1.00012000,
			amount:    0.2,
			nonceMark: 10000,
			fees:      2000,
			fail:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := signer.AddWithdrawTxOutputs(tt.tx, tt.to, tt.total, tt.amount, tt.nonceMark, tt.fees, tt.cancelTx)
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
