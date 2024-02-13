package bitcoin

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/common"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
)

const (
	satoshiPerBitcoin       = 1e8
	bytesPerKB              = 1000
	bytesEmptyTx            = 10  // an empty tx is about 10 bytes
	bytesPerInput           = 41  // each input is about 41 bytes
	bytesPerOutput          = 31  // each output is about 31 bytes
	bytes1stWitness         = 110 // the 1st witness incurs about 110 bytes and it may vary
	bytesPerWitness         = 108 // each additional witness incurs about 108 bytes and it may vary
	defaultDepositorFeeRate = 20  // 20 sat/byte is the default depositor fee rate
)

var (
	BtcOutTxBytesDepositor  uint64
	BtcOutTxBytesWithdrawer uint64
	DefaultDepositorFee     float64
)

func init() {
	BtcOutTxBytesDepositor = SegWitTxSizeDepositor()   // 68vB, the outtx size incurred by the depositor
	BtcOutTxBytesWithdrawer = SegWitTxSizeWithdrawer() // 171vB, the outtx size incurred by the withdrawer

	// default depositor fee calculation is based on a fixed fee rate of 20 sat/byte just for simplicity.
	// In reality, the fee rate on UTXO deposit is different from the fee rate when the UTXO is spent.
	DefaultDepositorFee = DepositorFee(defaultDepositorFeeRate) // 0.00001360 (20 * 68vB / 100000000)
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

// CalcBlockAvgFeeRate calculates the average gas rate (in sat/vByte) for a given block
func CalcBlockAvgFeeRate(blockVb *btcjson.GetBlockVerboseTxResult, netParams *chaincfg.Params) (int64, error) {
	// sanity check
	if len(blockVb.Tx) == 0 {
		return 0, errors.New("block has no transactions")
	}
	if len(blockVb.Tx) == 1 {
		return 0, nil // only coinbase tx, it happens
	}
	txCoinbase := &blockVb.Tx[0]
	if blockVb.Weight < blockchain.WitnessScaleFactor {
		return 0, fmt.Errorf("block weight %d too small", blockVb.Weight)
	}
	if blockVb.Weight < txCoinbase.Weight {
		return 0, fmt.Errorf("block weight %d less than coinbase tx weight %d", blockVb.Weight, txCoinbase.Weight)
	}
	if blockVb.Height <= 0 || blockVb.Height > math.MaxInt32 {
		return 0, fmt.Errorf("invalid block height %d", blockVb.Height)
	}

	// make sure the first tx is coinbase tx
	txBytes, err := hex.DecodeString(txCoinbase.Hex)
	if err != nil {
		return 0, fmt.Errorf("failed to decode coinbase tx %s", txCoinbase.Txid)
	}
	tx, err := btcutil.NewTxFromBytes(txBytes)
	if err != nil {
		return 0, fmt.Errorf("failed to parse coinbase tx %s", txCoinbase.Txid)
	}
	if !blockchain.IsCoinBaseTx(tx.MsgTx()) {
		return 0, fmt.Errorf("first tx %s is not coinbase tx", txCoinbase.Txid)
	}

	// calculate fees earned by the miner
	btcEarned := int64(0)
	for _, out := range tx.MsgTx().TxOut {
		if out.Value > 0 {
			btcEarned += out.Value
		}
	}
	// #nosec G701 checked above
	subsidy := blockchain.CalcBlockSubsidy(int32(blockVb.Height), netParams)
	if btcEarned < subsidy {
		return 0, fmt.Errorf("miner earned %d, less than subsidy %d", btcEarned, subsidy)
	}
	txsFees := btcEarned - subsidy

	// sum up weight of all txs (<= 4 MWU)
	txsWeight := int32(0)
	for i, tx := range blockVb.Tx {
		// coinbase doesn't pay fees, so we exclude it
		if i > 0 && tx.Weight > 0 {
			txsWeight += tx.Weight
		}
	}

	// calculate average fee rate.
	vBytes := txsWeight / blockchain.WitnessScaleFactor
	return txsFees / int64(vBytes), nil
}

// CalcDepositorFee calculates the depositor fee for a given block
func CalcDepositorFee(blockVb *btcjson.GetBlockVerboseTxResult, chainID int64, netParams *chaincfg.Params, logger zerolog.Logger) float64 {
	// use dynamic fee or default
	dynamicFee := true

	// use default fee for regnet
	if common.IsBitcoinRegnet(chainID) {
		dynamicFee = false
	}
	// mainnet dynamic fee takes effect only after a planned upgrade height
	if common.IsBitcoinMainnet(chainID) && blockVb.Height < DynamicDepositorFeeHeight {
		dynamicFee = false
	}
	if !dynamicFee {
		return DefaultDepositorFee
	}

	// calculate deposit fee rate
	feeRate, err := CalcBlockAvgFeeRate(blockVb, netParams)
	if err != nil {
		feeRate = defaultDepositorFeeRate // use default fee rate if calculation fails, should not happen
		logger.Error().Err(err).Msgf("cannot calculate fee rate for block %d", blockVb.Height)
	}
	feeRate = feeRate * common.DefaultGasPriceMultiplier
	return DepositorFee(feeRate)
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
