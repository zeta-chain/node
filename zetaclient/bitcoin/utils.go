package bitcoin

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/blockchain"

	"github.com/btcsuite/btcd/txscript"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/btcsuite/btcd/wire"
)

const (
	satoshiPerBitcoin = 1e8
	bytesPerKB        = 1000
	bytesPerInput     = 41  // each input is about 41 bytes
	bytesPerOutput    = 31  // each output is about 31 bytes
	bytes1stWitness   = 110 // the 1st witness incurs about 110 bytes and it may vary
	bytesPerWitness   = 108 // each additional witness incurs about 108 bytes and it may vary
)

var (
	BtcOutTxBytesDepositor  uint64
	BtcOutTxBytesWithdrawer uint64
	BtcDepositorFeeMin      float64
)

func init() {
	BtcOutTxBytesDepositor = SegWitTxSizeDepositor()   // 68vB, the outtx size incurred by the depositor
	BtcOutTxBytesWithdrawer = SegWitTxSizeWithdrawer() // 171vB, the outtx size incurred by the withdrawer

	// depositor fee calculation is based on a fixed fee rate of 20 sat/byte just for simplicity.
	// In reality, the fee rate on UTXO deposit is different from the fee rate when the UTXO is spent.
	BtcDepositorFeeMin = DepositorFee(20) // 0.00001360 (20 * 68vB / 100000000), the minimum deposit fee in BTC for 20 sat/byte
}

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

// FeeRateToSatPerByte converts a fee rate in BTC/KB to sat/byte.
func FeeRateToSatPerByte(rate float64) *big.Int {
	// #nosec G701 always in range
	satPerKB := new(big.Int).SetInt64(int64(rate * satoshiPerBitcoin))
	return new(big.Int).Div(satPerKB, big.NewInt(bytesPerKB))
}

// WiredTxSize calculates the wired tx size in bytes
func WiredTxSize(numInputs uint64, numOutputs uint64) uint64 {
	// Version 4 bytes + LockTime 4 bytes + Serialized varint size for the
	// number of transaction inputs and outputs.
	// #nosec G701 always positive
	return uint64(8 + wire.VarIntSerializeSize(numInputs) + wire.VarIntSerializeSize(numOutputs))
}

// EstimateSegWitTxSize estimates SegWit tx size
func EstimateSegWitTxSize(numInputs uint64, numOutputs uint64) uint64 {
	if numInputs == 0 {
		return 0
	}
	bytesWiredTx := WiredTxSize(numInputs, numOutputs)
	bytesInput := numInputs * bytesPerInput
	bytesOutput := numOutputs * bytesPerOutput
	bytesWitness := bytes1stWitness + (numInputs-1)*bytesPerWitness
	// https://github.com/bitcoin/bips/blob/master/bip-0141.mediawiki#transaction-size-calculations
	// Calculation for signed SegWit tx: blockchain.GetTransactionWeight(tx) / 4
	return bytesWiredTx + bytesInput + bytesOutput + bytesWitness/blockchain.WitnessScaleFactor
}

// SegWitTxSizeDepositor returns SegWit tx size (68vB) incurred by the depositor
func SegWitTxSizeDepositor() uint64 {
	return bytesPerInput + bytesPerWitness/blockchain.WitnessScaleFactor
}

// SegWitTxSizeWithdrawer returns SegWit tx size (171vB) incurred by the withdrawer (1 input, 3 outputs)
func SegWitTxSizeWithdrawer() uint64 {
	bytesWiredTx := WiredTxSize(1, 3)
	bytesInput := uint64(1) * bytesPerInput   // nonce mark
	bytesOutput := uint64(3) * bytesPerOutput // 3 outputs: new nonce mark, payment, change
	return bytesWiredTx + bytesInput + bytesOutput + bytes1stWitness/blockchain.WitnessScaleFactor
}

// DepositorFee calculates the depositor fee in BTC for a given sat/byte fee rate
// Note: the depositor fee is charged in order to cover the cost of spending the deposited UTXO in the future
func DepositorFee(satPerByte int64) float64 {
	return float64(satPerByte) * float64(BtcOutTxBytesDepositor) / satoshiPerBitcoin
}

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

func PayToWitnessPubKeyHashScript(pubKeyHash []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(pubKeyHash).Script()
}

type DynamicTicker struct {
	name     string
	interval uint64
	impl     *time.Ticker
}

func NewDynamicTicker(name string, interval uint64) (*DynamicTicker, error) {
	if interval <= 0 {
		return nil, fmt.Errorf("non-positive ticker interval %d for %s", interval, name)
	}

	return &DynamicTicker{
		name:     name,
		interval: interval,
		impl:     time.NewTicker(time.Duration(interval) * time.Second),
	}, nil
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
