package bitcoin

import (
	"bytes"
	"encoding/hex"
	"math"
	"math/big"
	"path"
	"sync"
	"testing"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

func TestConfirmationThreshold(t *testing.T) {
	client := &BTCChainClient{Mu: &sync.Mutex{}}
	t.Run("should return confirmations in chain param", func(t *testing.T) {
		client.SetChainParams(observertypes.ChainParams{ConfirmationCount: 3})
		require.Equal(t, int64(3), client.ConfirmationsThreshold(big.NewInt(1000)))
	})

	t.Run("should return big value confirmations", func(t *testing.T) {
		client.SetChainParams(observertypes.ChainParams{ConfirmationCount: 3})
		require.Equal(t, int64(bigValueConfirmationCount), client.ConfirmationsThreshold(big.NewInt(bigValueSats)))
	})

	t.Run("big value confirmations is the upper cap", func(t *testing.T) {
		client.SetChainParams(observertypes.ChainParams{ConfirmationCount: bigValueConfirmationCount + 1})
		require.Equal(t, int64(bigValueConfirmationCount), client.ConfirmationsThreshold(big.NewInt(1000)))
	})
}

func TestAvgFeeRateBlock828440(t *testing.T) {
	// load archived block 828440
	var blockVb btcjson.GetBlockVerboseTxResult
	err := testutils.LoadObjectFromJSONFile(&blockVb, path.Join(testutils.TestDataPath, "bitcoin_block_trimmed_828440.json"))
	require.NoError(t, err)

	// https://mempool.space/block/000000000000000000025ca01d2c1094b8fd3bacc5468cc3193ced6a14618c27
	var blockMb testutils.MempoolBlock
	err = testutils.LoadObjectFromJSONFile(&blockMb, path.Join(testutils.TestDataPath, "mempool.space_block_828440.json"))
	require.NoError(t, err)

	gasRate, err := CalcBlockAvgFeeRate(&blockVb, &chaincfg.MainNetParams)
	require.NoError(t, err)
	require.Equal(t, int64(blockMb.Extras.AvgFeeRate), gasRate)
}

func TestAvgFeeRateBlock828440Errors(t *testing.T) {
	// load archived block 828440
	var blockVb btcjson.GetBlockVerboseTxResult
	err := testutils.LoadObjectFromJSONFile(&blockVb, path.Join(testutils.TestDataPath, "bitcoin_block_trimmed_828440.json"))
	require.NoError(t, err)

	t.Run("block has no transactions", func(t *testing.T) {
		emptyVb := btcjson.GetBlockVerboseTxResult{Tx: []btcjson.TxRawResult{}}
		_, err := CalcBlockAvgFeeRate(&emptyVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "block has no transactions")
	})
	t.Run("it's okay if block has only coinbase tx", func(t *testing.T) {
		coinbaseVb := btcjson.GetBlockVerboseTxResult{Tx: []btcjson.TxRawResult{
			blockVb.Tx[0],
		}}
		_, err := CalcBlockAvgFeeRate(&coinbaseVb, &chaincfg.MainNetParams)
		require.NoError(t, err)
	})
	t.Run("tiny block weight should fail", func(t *testing.T) {
		invalidVb := blockVb
		invalidVb.Weight = 3
		_, err := CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "block weight 3 too small")
	})
	t.Run("block weight should not be less than coinbase tx weight", func(t *testing.T) {
		invalidVb := blockVb
		invalidVb.Weight = blockVb.Tx[0].Weight - 1
		_, err := CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "less than coinbase tx weight")
	})
	t.Run("invalid block height should fail", func(t *testing.T) {
		invalidVb := blockVb
		invalidVb.Height = 0
		_, err := CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid block height")

		invalidVb.Height = math.MaxInt32 + 1
		_, err = CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid block height")
	})
	t.Run("failed to decode coinbase tx", func(t *testing.T) {
		invalidVb := blockVb
		invalidVb.Tx = []btcjson.TxRawResult{blockVb.Tx[0], blockVb.Tx[1]}
		invalidVb.Tx[0].Hex = "invalid hex"
		_, err := CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to decode coinbase tx")
	})
	t.Run("1st tx is not coinbase", func(t *testing.T) {
		invalidVb := blockVb
		invalidVb.Tx = []btcjson.TxRawResult{blockVb.Tx[1], blockVb.Tx[0]}
		_, err := CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "not coinbase tx")
	})
	t.Run("miner earned less than subsidy", func(t *testing.T) {
		invalidVb := blockVb
		coinbaseTxBytes, err := hex.DecodeString(blockVb.Tx[0].Hex)
		require.NoError(t, err)
		coinbaseTx, err := btcutil.NewTxFromBytes(coinbaseTxBytes)
		require.NoError(t, err)
		msgTx := coinbaseTx.MsgTx()

		// reduce subsidy by 1 satoshi
		for i := range msgTx.TxOut {
			if i == 0 {
				msgTx.TxOut[i].Value = blockchain.CalcBlockSubsidy(int32(blockVb.Height), &chaincfg.MainNetParams) - 1
			} else {
				msgTx.TxOut[i].Value = 0
			}
		}
		// calculate fee rate on modified coinbase tx
		var buf bytes.Buffer
		err = msgTx.Serialize(&buf)
		require.NoError(t, err)
		invalidVb.Tx[0].Hex = hex.EncodeToString(buf.Bytes())
		_, err = CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "less than subsidy")
	})
}

func TestCalcDepositorFee828440(t *testing.T) {
	// load archived block 828440
	var blockVb btcjson.GetBlockVerboseTxResult
	err := testutils.LoadObjectFromJSONFile(&blockVb, path.Join(testutils.TestDataPath, "bitcoin_block_trimmed_828440.json"))
	require.NoError(t, err)
	dynamicFee828440 := DepositorFee(32 * common.DefaultGasPriceMultiplier)

	// should return default fee if it's a regtest block
	fee := CalcDepositorFee(&blockVb, 18444, &chaincfg.RegressionNetParams, log.Logger)
	require.Equal(t, DefaultDepositorFee, fee)

	// should return dynamic fee if it's a testnet block
	fee = CalcDepositorFee(&blockVb, 18332, &chaincfg.TestNet3Params, log.Logger)
	require.NotEqual(t, DefaultDepositorFee, fee)
	require.Equal(t, dynamicFee828440, fee)

	// mainnet should return default fee before upgrade height
	blockVb.Height = DynamicDepositorFeeHeight - 1
	fee = CalcDepositorFee(&blockVb, 8332, &chaincfg.MainNetParams, log.Logger)
	require.Equal(t, DefaultDepositorFee, fee)

	// mainnet should return dynamic fee after upgrade height
	blockVb.Height = DynamicDepositorFeeHeight
	fee = CalcDepositorFee(&blockVb, 8332, &chaincfg.MainNetParams, log.Logger)
	require.NotEqual(t, DefaultDepositorFee, fee)
	require.Equal(t, dynamicFee828440, fee)
}
