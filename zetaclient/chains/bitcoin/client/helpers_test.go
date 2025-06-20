package client_test

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/testutils/testrpc"
)

func Test_GetTransactionInputSpender(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name               string
		txid               string
		vout               uint32
		spenderPkScriptHex string
		want               string
		rpcErr             string
		errMsg             string
	}{
		{
			name:               "should return the correct input spender address",
			txid:               "45823fe9c92f5bcfd0ab09c2021088a7170f24ed90bcb43debd18efc0b599c1f",
			vout:               0,
			spenderPkScriptHex: "00148831681655d9d379309fe3c03fcee2d076bd0b90",
			want:               "bc1q3qcks9j4m8fhjvylu0qrlnhz6pmt6zuslwxudf",
		},
		{
			name:   "should return an error if unable to get raw transaction",
			txid:   "45823fe9c92f5bcfd0ab09c2021088a7170f24ed90bcb43debd18efc0b599c1f",
			vout:   0,
			want:   "",
			errMsg: "unable to get raw transaction",
		},
		{
			name:               "should return an error if vout index is out of range",
			txid:               "45823fe9c92f5bcfd0ab09c2021088a7170f24ed90bcb43debd18efc0b599c1f",
			vout:               1,
			spenderPkScriptHex: "00148831681655d9d379309fe3c03fcee2d076bd0b90",
			want:               "",
			errMsg:             "out of range",
		},
	}

	// ACT
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			// create mock BTC server
			btcServer, btcConfig := testrpc.NewBtcServer(t)

			// create test suite with BTC server
			ts := newTestSuite(t, btcConfig)

			// mock raw transaction if provided
			if tt.spenderPkScriptHex != "" {
				preTxid := sample.BtcHash().String() // any txid
				msgTx := createMsgTx(t, preTxid, 0, tt.spenderPkScriptHex)
				btcServer.OnSetRawTransaction(t, *msgTx, tt.txid)
			}

			// ACT
			senderAddr, err := ts.GetTransactionInputSpender(ctx, tt.txid, tt.vout)

			// ASSERT
			if tt.errMsg == "" {
				require.NoError(t, err)
				require.Equal(t, tt.want, senderAddr)
			} else {
				require.Empty(t, senderAddr)
				require.Contains(t, err.Error(), tt.errMsg)
			}
		})
	}
}

func Test_GetTransactionInitiator(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name                    string
		txid                    string
		preTxid                 string
		preVout                 uint32
		preTxSpenderPkScriptHex string
		want                    string
		errMsg                  string
	}{
		{
			name:                    "should return the correct initiator address",
			txid:                    "4e53d9bc84a9ede9f5665ae5bcbace3929d13445b4a8319caf80e482126b4d74",
			preTxid:                 "45823fe9c92f5bcfd0ab09c2021088a7170f24ed90bcb43debd18efc0b599c1f",
			preVout:                 0,
			preTxSpenderPkScriptHex: "00148831681655d9d379309fe3c03fcee2d076bd0b90",
			want:                    "bc1q3qcks9j4m8fhjvylu0qrlnhz6pmt6zuslwxudf",
		},
		{
			name:   "should return an error if unable to get raw transaction",
			txid:   "",
			want:   "",
			errMsg: "unable to get raw transaction",
		},
		{
			name:    "should return an error if unable to get first input spender",
			txid:    "4e53d9bc84a9ede9f5665ae5bcbace3929d13445b4a8319caf80e482126b4d74",
			preTxid: "45823fe9c92f5bcfd0ab09c2021088a7170f24ed90bcb43debd18efc0b599c1f",
			preVout: 0,
			want:    "",
			errMsg:  "unable to get transaction input spender",
		},
	}

	// ACT
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			// create mock BTC server
			btcServer, btcConfig := testrpc.NewBtcServer(t)

			// create test suite with BTC server
			ts := newTestSuite(t, btcConfig)

			// mock this tx
			if tt.txid != "" {
				msgTx := createMsgTx(t, tt.preTxid, tt.preVout, "0123") // pkScript is irrelevant
				btcServer.OnSetRawTransaction(t, *msgTx, tt.txid)
			}

			// mock the previous tx
			if tt.preTxSpenderPkScriptHex != "" {
				prevTx := createMsgTx(t, "4567", 0, tt.preTxSpenderPkScriptHex) // only output script matters
				btcServer.OnSetRawTransaction(t, *prevTx, tt.preTxid)
			}

			// ACT
			initiator, err := ts.GetTransactionInitiator(ctx, tt.txid)

			// ASSERT
			if tt.errMsg == "" {
				require.NoError(t, err)
				require.Equal(t, tt.want, initiator)
			} else {
				require.Empty(t, initiator)
				require.Contains(t, err.Error(), tt.errMsg)
			}
		})
	}
}

// createMsgTx creates a new MsgTx with a single input and output.
func createMsgTx(t *testing.T, preTxid string, preVout uint32, outPkScriptHex string) *wire.MsgTx {
	preHash, err := chainhash.NewHashFromStr(preTxid)
	require.NoError(t, err)

	pkScriptBytes, err := hex.DecodeString(outPkScriptHex)
	require.NoError(t, err)

	msgTx := wire.NewMsgTx(wire.TxVersion)
	msgTx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(preHash, preVout), nil, nil))
	msgTx.AddTxOut(wire.NewTxOut(1000000, pkScriptBytes))

	return msgTx
}
