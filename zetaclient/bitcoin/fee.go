package bitcoin

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/bitcoin"
	clientcommon "github.com/zeta-chain/zetacore/zetaclient/common"

	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
)

const (
	bytesPerKB              = 1000
	bytesEmptyTx            = 10  // an empty tx is 10 bytes
	bytesPerInput           = 41  // each input is 41 bytes
	bytesPerOutputP2TR      = 43  // each P2TR output is 43 bytes
	bytesPerOutputP2WSH     = 43  // each P2WSH output is 43 bytes
	bytesPerOutputP2WPKH    = 31  // each P2WPKH output is 31 bytes
	bytesPerOutputP2SH      = 32  // each P2SH output is 32 bytes
	bytesPerOutputP2PKH     = 34  // each P2PKH output is 34 bytes
	bytesPerOutputAvg       = 37  // average size of all above types of outputs (36.6 bytes)
	bytes1stWitness         = 110 // the 1st witness incurs about 110 bytes and it may vary
	bytesPerWitness         = 108 // each additional witness incurs about 108 bytes and it may vary
	defaultDepositorFeeRate = 20  // 20 sat/byte is the default depositor fee rate

	outTxBytesMin = uint64(239)  // 239vB == EstimateSegWitTxSize(2, 2, toP2WPKH)
	outTxBytesMax = uint64(1543) // 1543v == EstimateSegWitTxSize(21, 2, toP2TR)
	outTxBytesAvg = uint64(245)  // 245vB is a suggested gas limit for zeta core
)

var (
	BtcOutTxBytesDepositor  uint64
	BtcOutTxBytesWithdrawer uint64
	DefaultDepositorFee     float64
)

func init() {
	BtcOutTxBytesDepositor = OuttxSizeDepositor()   // 68vB, the outtx size incurred by the depositor
	BtcOutTxBytesWithdrawer = OuttxSizeWithdrawer() // 177vB, the outtx size incurred by the withdrawer

	// default depositor fee calculation is based on a fixed fee rate of 20 sat/byte just for simplicity.
	// In reality, the fee rate on UTXO deposit is different from the fee rate when the UTXO is spent.
	DefaultDepositorFee = DepositorFee(defaultDepositorFeeRate) // 0.00001360 (20 * 68vB / 100000000)
}

// FeeRateToSatPerByte converts a fee rate in BTC/KB to sat/byte.
func FeeRateToSatPerByte(rate float64) *big.Int {
	// #nosec G701 always in range
	satPerKB := new(big.Int).SetInt64(int64(rate * btcutil.SatoshiPerBitcoin))
	return new(big.Int).Div(satPerKB, big.NewInt(bytesPerKB))
}

// WiredTxSize calculates the wired tx size in bytes
func WiredTxSize(numInputs uint64, numOutputs uint64) uint64 {
	// Version 4 bytes + LockTime 4 bytes + Serialized varint size for the
	// number of transaction inputs and outputs.
	// #nosec G701 always positive
	return uint64(8 + wire.VarIntSerializeSize(numInputs) + wire.VarIntSerializeSize(numOutputs))
}

// EstimateOuttxSize estimates the size of a outtx in vBytes
func EstimateOuttxSize(numInputs uint64, payees []btcutil.Address) uint64 {
	if numInputs == 0 {
		return 0
	}
	// #nosec G701 always positive
	numOutputs := 2 + uint64(len(payees))
	bytesWiredTx := WiredTxSize(numInputs, numOutputs)
	bytesInput := numInputs * bytesPerInput
	bytesOutput := uint64(2) * bytesPerOutputP2WPKH // new nonce mark, change

	// calculate the size of the outputs to payees
	bytesToPayees := uint64(0)
	for _, to := range payees {
		bytesToPayees += GetOutputSizeByAddress(to)
	}
	// calculate the size of the witness
	bytesWitness := bytes1stWitness + (numInputs-1)*bytesPerWitness
	// https://github.com/bitcoin/bips/blob/master/bip-0141.mediawiki#transaction-size-calculations
	// Calculation for signed SegWit tx: blockchain.GetTransactionWeight(tx) / 4
	return bytesWiredTx + bytesInput + bytesOutput + bytesToPayees + bytesWitness/blockchain.WitnessScaleFactor
}

// GetOutputSizeByAddress returns the size of a tx output in bytes by the given address
func GetOutputSizeByAddress(to btcutil.Address) uint64 {
	switch addr := to.(type) {
	case *bitcoin.AddressTaproot:
		if addr == nil {
			return 0
		}
		return bytesPerOutputP2TR
	case *btcutil.AddressWitnessScriptHash:
		if addr == nil {
			return 0
		}
		return bytesPerOutputP2WSH
	case *btcutil.AddressWitnessPubKeyHash:
		if addr == nil {
			return 0
		}
		return bytesPerOutputP2WPKH
	case *btcutil.AddressScriptHash:
		if addr == nil {
			return 0
		}
		return bytesPerOutputP2SH
	case *btcutil.AddressPubKeyHash:
		if addr == nil {
			return 0
		}
		return bytesPerOutputP2PKH
	default:
		return bytesPerOutputP2WPKH
	}
}

// OuttxSizeDepositor returns outtx size (68vB) incurred by the depositor
func OuttxSizeDepositor() uint64 {
	return bytesPerInput + bytesPerWitness/blockchain.WitnessScaleFactor
}

// OuttxSizeWithdrawer returns outtx size (177vB) incurred by the withdrawer (1 input, 3 outputs)
func OuttxSizeWithdrawer() uint64 {
	bytesWiredTx := WiredTxSize(1, 3)
	bytesInput := uint64(1) * bytesPerInput         // nonce mark
	bytesOutput := uint64(2) * bytesPerOutputP2WPKH // 2 P2WPKH outputs: new nonce mark, change
	bytesOutput += bytesPerOutputAvg                // 1 output to withdrawer's address
	return bytesWiredTx + bytesInput + bytesOutput + bytes1stWitness/blockchain.WitnessScaleFactor
}

// DepositorFee calculates the depositor fee in BTC for a given sat/byte fee rate
// Note: the depositor fee is charged in order to cover the cost of spending the deposited UTXO in the future
func DepositorFee(satPerByte int64) float64 {
	return float64(satPerByte) * float64(BtcOutTxBytesDepositor) / btcutil.SatoshiPerBitcoin
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
	// #nosec G701 always in range
	feeRate = int64(float64(feeRate) * clientcommon.BTCOuttxGasPriceMultiplier)
	return DepositorFee(feeRate)
}
