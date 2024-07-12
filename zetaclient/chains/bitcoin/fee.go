package bitcoin

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/pkg/chains"
	clientcommon "github.com/zeta-chain/zetacore/zetaclient/common"
)

const (
	bytesPerKB              = 1000
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

	OutboundBytesMin = uint64(239)  // 239vB == EstimateSegWitTxSize(2, 2, toP2WPKH)
	OutboundBytesMax = uint64(1543) // 1543v == EstimateSegWitTxSize(21, 2, toP2TR)
	OutboundBytesAvg = uint64(245)  // 245vB is a suggested gas limit for zetacore

	DynamicDepositorFeeHeight = 834500 // DynamicDepositorFeeHeight contains the starting height (Bitcoin mainnet) from which dynamic depositor fee will take effect
)

var (
	// BtcOutboundBytesDepositor is the outbound size incurred by the depositor: 68vB
	BtcOutboundBytesDepositor = OutboundSizeDepositor()

	// BtcOutboundBytesWithdrawer is the outbound size incurred by the withdrawer: 177vB
	BtcOutboundBytesWithdrawer = OutboundSizeWithdrawer()

	// DefaultDepositorFee is the default depositor fee is 0.00001360 BTC (20 * 68vB / 100000000)
	// default depositor fee calculation is based on a fixed fee rate of 20 sat/byte just for simplicity.
	DefaultDepositorFee = DepositorFee(defaultDepositorFeeRate)
)

// FeeRateToSatPerByte converts a fee rate in BTC/KB to sat/byte.
func FeeRateToSatPerByte(rate float64) *big.Int {
	// #nosec G115 always in range
	satPerKB := new(big.Int).SetInt64(int64(rate * btcutil.SatoshiPerBitcoin))
	return new(big.Int).Div(satPerKB, big.NewInt(bytesPerKB))
}

// WiredTxSize calculates the wired tx size in bytes
func WiredTxSize(numInputs uint64, numOutputs uint64) uint64 {
	// Version 4 bytes + LockTime 4 bytes + Serialized varint size for the
	// number of transaction inputs and outputs.
	// #nosec G115 always positive
	return uint64(8 + wire.VarIntSerializeSize(numInputs) + wire.VarIntSerializeSize(numOutputs))
}

// EstimateOutboundSize estimates the size of an outbound in vBytes
func EstimateOutboundSize(numInputs uint64, payees []btcutil.Address) (uint64, error) {
	if numInputs == 0 {
		return 0, nil
	}
	// #nosec G115 always positive
	numOutputs := 2 + uint64(len(payees))
	bytesWiredTx := WiredTxSize(numInputs, numOutputs)
	bytesInput := numInputs * bytesPerInput
	bytesOutput := uint64(2) * bytesPerOutputP2WPKH // new nonce mark, change

	// calculate the size of the outputs to payees
	bytesToPayees := uint64(0)
	for _, to := range payees {
		sizeOutput, err := GetOutputSizeByAddress(to)
		if err != nil {
			return 0, err
		}
		bytesToPayees += sizeOutput
	}

	// calculate the size of the witness
	bytesWitness := bytes1stWitness + (numInputs-1)*bytesPerWitness

	// https://github.com/bitcoin/bips/blob/master/bip-0141.mediawiki#transaction-size-calculations
	// Calculation for signed SegWit tx: blockchain.GetTransactionWeight(tx) / 4
	return bytesWiredTx + bytesInput + bytesOutput + bytesToPayees + bytesWitness/blockchain.WitnessScaleFactor, nil
}

// GetOutputSizeByAddress returns the size of a tx output in bytes by the given address
func GetOutputSizeByAddress(to btcutil.Address) (uint64, error) {
	switch addr := to.(type) {
	case *chains.AddressTaproot:
		if addr == nil {
			return 0, nil
		}
		return bytesPerOutputP2TR, nil
	case *btcutil.AddressWitnessScriptHash:
		if addr == nil {
			return 0, nil
		}
		return bytesPerOutputP2WSH, nil
	case *btcutil.AddressWitnessPubKeyHash:
		if addr == nil {
			return 0, nil
		}
		return bytesPerOutputP2WPKH, nil
	case *btcutil.AddressScriptHash:
		if addr == nil {
			return 0, nil
		}
		return bytesPerOutputP2SH, nil
	case *btcutil.AddressPubKeyHash:
		if addr == nil {
			return 0, nil
		}
		return bytesPerOutputP2PKH, nil
	default:
		return 0, fmt.Errorf("cannot get output size for address type %T", to)
	}
}

// OutboundSizeDepositor returns outbound size (68vB) incurred by the depositor
func OutboundSizeDepositor() uint64 {
	return bytesPerInput + bytesPerWitness/blockchain.WitnessScaleFactor
}

// OutboundSizeWithdrawer returns outbound size (177vB) incurred by the withdrawer (1 input, 3 outputs)
func OutboundSizeWithdrawer() uint64 {
	bytesWiredTx := WiredTxSize(1, 3)
	bytesInput := uint64(1) * bytesPerInput         // nonce mark
	bytesOutput := uint64(2) * bytesPerOutputP2WPKH // 2 P2WPKH outputs: new nonce mark, change
	bytesOutput += bytesPerOutputAvg                // 1 output to withdrawer's address

	return bytesWiredTx + bytesInput + bytesOutput + bytes1stWitness/blockchain.WitnessScaleFactor
}

// DepositorFee calculates the depositor fee in BTC for a given sat/byte fee rate
// Note: the depositor fee is charged in order to cover the cost of spending the deposited UTXO in the future
func DepositorFee(satPerByte int64) float64 {
	return float64(satPerByte) * float64(BtcOutboundBytesDepositor) / btcutil.SatoshiPerBitcoin
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
	// #nosec G115 checked above
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
func CalcDepositorFee(
	blockVb *btcjson.GetBlockVerboseTxResult,
	chainID int64,
	netParams *chaincfg.Params,
	logger zerolog.Logger,
) float64 {
	// use default fee for regnet
	if chains.IsBitcoinRegnet(chainID) {
		return DefaultDepositorFee
	}
	// mainnet dynamic fee takes effect only after a planned upgrade height
	if chains.IsBitcoinMainnet(chainID) && blockVb.Height < DynamicDepositorFeeHeight {
		return DefaultDepositorFee
	}

	// calculate deposit fee rate
	feeRate, err := CalcBlockAvgFeeRate(blockVb, netParams)
	if err != nil {
		feeRate = defaultDepositorFeeRate // use default fee rate if calculation fails, should not happen
		logger.Error().Err(err).Msgf("cannot calculate fee rate for block %d", blockVb.Height)
	}

	// #nosec G115 always in range
	feeRate = int64(float64(feeRate) * clientcommon.BTCOutboundGasPriceMultiplier)

	return DepositorFee(feeRate)
}
