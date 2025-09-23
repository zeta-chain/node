package common

import (
	"context"
	"encoding/hex"
	"fmt"
	"math"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/zetaclient/common"
)

const (
	// constants related to transaction size calculations
	bytesPerInput        = 41          // each input is 41 bytes
	bytesPerOutputP2TR   = 43          // each P2TR output is 43 bytes
	bytesPerOutputP2WSH  = 43          // each P2WSH output is 43 bytes
	bytesPerOutputP2WPKH = 31          // each P2WPKH output is 31 bytes
	bytesPerOutputP2SH   = 32          // each P2SH output is 32 bytes
	bytesPerOutputP2PKH  = 34          // each P2PKH output is 34 bytes
	bytesPerOutputAvg    = 37          // average size of all above types of outputs (36.6 bytes)
	bytes1stWitness      = 110         // the 1st witness incurs about 110 bytes and it may vary
	bytesPerWitness      = 108         // each additional witness incurs about 108 bytes and it may vary
	OutboundBytesMin     = int64(239)  // 239vB == EstimateOutboundSize(2, 2, toP2WPKH)
	OutboundBytesMax     = int64(1543) // 1543v == EstimateOutboundSize(21, 2, toP2TR)

	// bytesPerKB is the number of vB in a KB
	bytesPerKB = 1000

	// defaultDepositorFeeRate is the default fee rate for depositor fee, 20 sat/vB
	defaultDepositorFeeRate = 20

	// defaultTestnetFeeRate is the default fee rate for testnet, 10 sat/vB
	defaultTestnetFeeRate = 10

	// feeRateCountBackBlocks is the default number of blocks to look back for fee rate estimation
	feeRateCountBackBlocks = 2
)

var (
	// BtcOutboundBytesDepositor is the outbound size incurred by the depositor: 68vB
	BtcOutboundBytesDepositor = OutboundSizeDepositor()

	// BtcOutboundBytesWithdrawer is the outbound size incurred by the withdrawer: 177vB
	// This will be the suggested gas limit used for zetacore
	BtcOutboundBytesWithdrawer = OutboundSizeWithdrawer()

	// DefaultDepositorFee is the default depositor fee is 0.00001360 BTC (20 * 68vB / 100000000)
	// default depositor fee calculation is based on a fixed fee rate of 20 sat/byte just for simplicity.
	DefaultDepositorFee = DepositorFee(defaultDepositorFeeRate)
)

type BitcoinClient interface {
	GetBlockCount(ctx context.Context) (int64, error)
	GetBlockHash(ctx context.Context, blockHeight int64) (*chainhash.Hash, error)
	GetBlockHeader(ctx context.Context, hash *chainhash.Hash) (*wire.BlockHeader, error)
	GetBlockVerbose(ctx context.Context, hash *chainhash.Hash) (*btcjson.GetBlockVerboseTxResult, error)
	GetTransactionFeeAndRate(ctx context.Context, tx *btcjson.TxRawResult) (int64, int64, error)
}

// DepositorFeeCalculator is a function type to calculate the Bitcoin depositor fee
type DepositorFeeCalculator func(context.Context, BitcoinClient, *btcjson.TxRawResult, *chaincfg.Params) (float64, error)

// FeeRateToSatPerByte converts a fee rate from BTC/KB to sat/vB.
func FeeRateToSatPerByte(rate float64) (uint64, error) {
	if rate <= 0 {
		return 0, fmt.Errorf("invalid fee rate %f", rate)
	}
	satPerKB := rate * btcutil.SatoshiPerBitcoin

	// #nosec G115 always positive
	return uint64(satPerKB / bytesPerKB), nil
}

// WiredTxSize calculates the wired tx size in bytes
func WiredTxSize(numInputs uint64, numOutputs uint64) int64 {
	// Version 4 bytes + LockTime 4 bytes + Serialized varint size for the
	// number of transaction inputs and outputs.
	// #nosec G115 always positive
	return int64(8 + wire.VarIntSerializeSize(numInputs) + wire.VarIntSerializeSize(numOutputs))
}

// EstimateOutboundSize estimates the size of an outbound in vBytes
func EstimateOutboundSize(numInputs int64, payees []btcutil.Address) (int64, error) {
	if numInputs <= 0 {
		return 0, nil
	}
	// #nosec G115 always positive
	numOutputs := 2 + uint64(len(payees))
	// #nosec G115 checked positive
	bytesWiredTx := WiredTxSize(uint64(numInputs), numOutputs)
	bytesInput := numInputs * bytesPerInput
	bytesOutput := int64(2) * bytesPerOutputP2WPKH // new nonce mark, change

	// calculate the size of the outputs to payees
	bytesToPayees := int64(0)
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
func GetOutputSizeByAddress(to btcutil.Address) (int64, error) {
	switch addr := to.(type) {
	case *btcutil.AddressTaproot:
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
func OutboundSizeDepositor() int64 {
	return bytesPerInput + bytesPerWitness/blockchain.WitnessScaleFactor
}

// OutboundSizeWithdrawer returns outbound size (177vB) incurred by the withdrawer (1 input, 3 outputs)
func OutboundSizeWithdrawer() int64 {
	bytesWiredTx := WiredTxSize(1, 3)
	bytesInput := int64(1) * bytesPerInput         // nonce mark
	bytesOutput := int64(2) * bytesPerOutputP2WPKH // 2 P2WPKH outputs: new nonce mark, change
	bytesOutput += bytesPerOutputAvg               // 1 output to withdrawer's address

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

// CalcDepositorFee calculates the depositor fee for a given tx result
func CalcDepositorFee(
	ctx context.Context,
	bitcoinClient BitcoinClient,
	rawResult *btcjson.TxRawResult,
	netParams *chaincfg.Params,
) (float64, error) {
	// use default fee for regnet
	if netParams.Name == chaincfg.RegressionNetParams.Name {
		return DefaultDepositorFee, nil
	}

	// get fee rate of the transaction
	_, feeRate, err := bitcoinClient.GetTransactionFeeAndRate(ctx, rawResult)
	if err != nil {
		return 0, errors.Wrapf(err, "error getting fee rate for tx %s", rawResult.Txid)
	}

	// apply gas price multiplier
	// #nosec G115 always in range
	feeRate = int64(float64(feeRate) * common.BTCOutboundGasPriceMultiplier)

	return DepositorFee(feeRate), nil
}

// GetRecentFeeRate gets the highest fee rate from recent blocks
// Note: this method should be used for testnet ONLY
func GetRecentFeeRate(ctx context.Context,
	bitcoinClient BitcoinClient,
	netParams *chaincfg.Params,
) (uint64, error) {
	// should avoid using this method for mainnet
	if netParams.Name == chaincfg.MainNetParams.Name {
		return 0, errors.New("GetRecentFeeRate should not be used for mainnet")
	}

	// get the current block number
	blockNumber, err := bitcoinClient.GetBlockCount(ctx)
	if err != nil {
		return 0, err
	}

	// get the highest fee rate among recent 'countBack' blocks to avoid underestimation
	highestRate := int64(0)
	for i := int64(0); i < feeRateCountBackBlocks; i++ {
		// get the block
		hash, err := bitcoinClient.GetBlockHash(ctx, blockNumber-i)
		if err != nil {
			return 0, err
		}
		block, err := bitcoinClient.GetBlockVerbose(ctx, hash)
		if err != nil {
			return 0, err
		}

		// computes the average fee rate of the block and take the higher rate
		avgFeeRate, err := CalcBlockAvgFeeRate(block, netParams)
		if err != nil {
			return 0, err
		}
		if avgFeeRate > highestRate {
			highestRate = avgFeeRate
		}
	}

	// use 10 sat/byte as default estimation if recent fee rate drops to 0
	if highestRate <= 0 {
		highestRate = defaultTestnetFeeRate
	}

	// #nosec G115 checked positive
	return uint64(highestRate), nil
}
