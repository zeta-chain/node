package zetaclient

import (
	"errors"
	"math"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/txscript"
	"github.com/rs/zerolog"
)

const (
	satoshiPerBitcoin = 1e8
)

func getSatoshis(btc float64) (int64, error) {
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
	return round(btc * satoshiPerBitcoin), nil
}

func round(f float64) int64 {
	if f < 0 {
		// #nosec G701 always in range
		return int64(f - 0.5)
	}
	// #nosec G701 always in range
	return int64(f + 0.5)
}

// roundGasPriceUp rounds up the gasPrice to the nearest multiple of base
func roundGasPriceUp(gasPrice *big.Int, base int64) *big.Int {
	oneUnit := big.NewInt(base) // e.g. 1 Gwei
	mod := new(big.Int)
	mod.Mod(gasPrice, oneUnit)
	if mod.Cmp(big.NewInt(0)) == 0 { // gasPrice is already a multiple of base
		return gasPrice
	}
	return new(big.Int).Add(gasPrice, new(big.Int).Sub(oneUnit, mod))
}

func payToWitnessPubKeyHashScript(pubKeyHash []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(pubKeyHash).Script()
}

type DynamicTicker struct {
	name     string
	interval uint64
	impl     *time.Ticker
}

func NewDynamicTicker(name string, interval uint64) *DynamicTicker {
	return &DynamicTicker{
		name:     name,
		interval: interval,
		impl:     time.NewTicker(time.Duration(interval) * time.Second),
	}
}

func (t *DynamicTicker) C() <-chan time.Time {
	return t.impl.C
}

func (t *DynamicTicker) UpdateInterval(newInterval uint64, logger zerolog.Logger) {
	if newInterval > 0 && t.interval != newInterval {
		t.impl.Stop()
		oldInterval := t.interval
		t.interval = newInterval
		t.impl = time.NewTicker(time.Duration(t.interval) * time.Second)
		logger.Info().Msgf("%s ticker interval changed from %d to %d", t.name, oldInterval, newInterval)
	}
}

func (t *DynamicTicker) Stop() {
	t.impl.Stop()
}
