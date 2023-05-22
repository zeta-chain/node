package zetaclient

import (
	"errors"
	"fmt"
	"math"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/txscript"
)

const (
	satoshiPerBitcoin = 1e8
)

func getSatoshis(btc float64) (int64, error) {
	// The amount is only considered invalid if it cannot be represented
	// as an integer type.  This may happen if f is NaN or +-Infinity.
	// BTC max amount is 21 mil and its at least 10^(-8) or one satoshi.
	switch {
	case math.IsNaN(btc):
		fallthrough
	case math.IsInf(btc, 1):
		fallthrough
	case math.IsInf(btc, -1):
		return 0, errors.New("invalid bitcoin amount")
	case btc > 21000000.0:
		return 0, errors.New("exceeded max bitcoin amount")
	case btc < 0.00000001:
		return 0, errors.New("cannot be less than 1 satoshi")
	}
	return round(btc * satoshiPerBitcoin), nil
}

func round(f float64) int64 {
	if f < 0 {
		return int64(f - 0.5)
	}
	return int64(f + 0.5)
}

func payToWitnessPubKeyHashScript(pubKeyHash []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(pubKeyHash).Script()
}

func utxoKey(utxo btcjson.ListUnspentResult) string {
	return fmt.Sprintf("%s_%d", utxo.TxID, utxo.Vout)
}
