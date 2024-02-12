package zetaclient

import (
	"math/big"
	"path"
	"sync"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

func TestConfirmationThreshold(t *testing.T) {
	client := &BitcoinChainClient{Mu: &sync.Mutex{}}
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
	// get archived block 828440
	var blockVb btcjson.GetBlockVerboseTxResult
	err := testutils.LoadObjectFromJSONFile(&blockVb, path.Join(testutils.TestDataPath, "bitcoin_block_828440.json"))
	require.NoError(t, err)

	// https://mempool.space/block/000000000000000000025ca01d2c1094b8fd3bacc5468cc3193ced6a14618c27
	var blockMb testutils.MempoolBlock
	err = testutils.LoadObjectFromJSONFile(&blockMb, path.Join(testutils.TestDataPath, "mempool.space_block_828440.json"))
	require.NoError(t, err)

	gasRate, err := CalcBlockAvgFeeRate(&blockVb, &chaincfg.MainNetParams)
	require.NoError(t, err)
	require.Equal(t, int64(blockMb.Extras.AvgFeeRate), gasRate)
}

func TestCalcDepositorFee828440(t *testing.T) {
	// get archived block 828440
	var blockVb btcjson.GetBlockVerboseTxResult
	err := testutils.LoadObjectFromJSONFile(&blockVb, path.Join(testutils.TestDataPath, "bitcoin_block_828440.json"))
	require.NoError(t, err)

	// should return default fee if it's a regtest block
	fee := CalcDepositorFee(&blockVb, 18444, &chaincfg.RegressionNetParams, log.Logger)
	require.Equal(t, DefaultDepositorFee, fee)

	// should return dynamic fee if it's a testnet block
	fee = CalcDepositorFee(&blockVb, 18332, &chaincfg.TestNet3Params, log.Logger)
	require.NotEqual(t, DefaultDepositorFee, fee)

	// mainnet should return default fee before upgrade height
	blockVb.Height = DynamicDepositorFeeHeight - 1
	fee = CalcDepositorFee(&blockVb, 8332, &chaincfg.MainNetParams, log.Logger)
	require.Equal(t, DefaultDepositorFee, fee)

	// mainnet should return dynamic fee after upgrade height
	blockVb.Height = DynamicDepositorFeeHeight
	fee = CalcDepositorFee(&blockVb, 8332, &chaincfg.MainNetParams, log.Logger)
	require.NotEqual(t, DefaultDepositorFee, fee)
}
