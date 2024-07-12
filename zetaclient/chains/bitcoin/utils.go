package bitcoin

import (
	"encoding/json"
	"math"

	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

// TODO(revamp): Remove utils.go and move the functions to the appropriate files

// PrettyPrintStruct returns a pretty-printed string representation of a struct
func PrettyPrintStruct(val interface{}) (string, error) {
	prettyStruct, err := json.MarshalIndent(
		val,
		"",
		" ",
	)
	if err != nil {
		return "", err
	}
	return string(prettyStruct), nil
}

// GetSatoshis converts a bitcoin amount to satoshis
func GetSatoshis(btc float64) (int64, error) {
	// The amount is only considered invalid if it cannot be represented
	// as an integer type.  This may happen if f is NaN or +-Infinity.
	// BTC max amount is 21 mil and its at least 0 (Note: bitcoin allows creating 0-value outputs)
	switch {
	case math.IsNaN(btc):
		fallthrough
	case math.IsInf(btc, 1):
		fallthrough
	case math.IsInf(btc, -1):
		return 0, errors.New("invalid bitcoin amount")
	case btc > 21000000.0:
		return 0, errors.New("exceeded max bitcoin amount")
	case btc < 0.0:
		return 0, errors.New("cannot be less than zero")
	}
	return round(btc * btcutil.SatoshiPerBitcoin), nil
}

// round rounds a float64 to the nearest integer
func round(f float64) int64 {
	if f < 0 {
		// #nosec G115 always in range
		return int64(f - 0.5)
	}
	// #nosec G115 always in range
	return int64(f + 0.5)
}
