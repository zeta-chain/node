package bitcoin

import (
	"bytes"
	"encoding/hex"
	"math"
	"math/big"
	"path"
	"strings"
	"sync"
	"testing"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/testutil/sample"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	appcontext "github.com/zeta-chain/zetacore/zetaclient/app_context"
	clientcommon "github.com/zeta-chain/zetacore/zetaclient/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	corecontext "github.com/zeta-chain/zetacore/zetaclient/core_context"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/stub"
)

func MockBTCClientMainnet() *BTCChainClient {
	cfg := config.NewConfig()
	coreContext := corecontext.NewZetaCoreContext(cfg)

	return &BTCChainClient{
		chain:       chains.BtcMainnetChain,
		zetaClient:  stub.NewMockZetaCoreBridge(),
		Tss:         stub.NewTSSMainnet(),
		coreContext: coreContext,
	}
}

// createRPCClientAndLoadTx is a helper function to load raw tx and feed it to mock rpc client
func createRPCClientAndLoadTx(t *testing.T, chainId int64, txHash string) *stub.MockBTCRPCClient {
	// file name for the archived MsgTx
	nameMsgTx := path.Join("../", testutils.TestDataPathBTC, testutils.FileNameBTCMsgTx(chainId, txHash))

	// load archived MsgTx
	var msgTx wire.MsgTx
	testutils.LoadObjectFromJSONFile(t, &msgTx, nameMsgTx)
	tx := btcutil.NewTx(&msgTx)

	// feed tx to mock rpc client
	rpcClient := stub.NewMockBTCRPCClient()
	rpcClient.WithRawTransaction(tx)
	return rpcClient
}

func TestNewBitcoinClient(t *testing.T) {
	t.Run("should return error because zetacore doesn't update core context", func(t *testing.T) {
		cfg := config.NewConfig()
		coreContext := corecontext.NewZetaCoreContext(cfg)
		appContext := appcontext.NewAppContext(coreContext, cfg)
		chain := chains.BtcMainnetChain
		bridge := stub.NewMockZetaCoreBridge()
		tss := stub.NewMockTSS(chains.BtcTestNetChain, sample.EthAddress().String(), "")
		loggers := clientcommon.ClientLogger{}
		btcCfg := cfg.BitcoinConfig
		ts := metrics.NewTelemetryServer()

		client, err := NewBitcoinClient(appContext, chain, bridge, tss, tempSQLiteDbPath, loggers, btcCfg, ts)
		require.ErrorContains(t, err, "btc chains params not initialized")
		require.Nil(t, client)
	})
}

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
	testutils.LoadObjectFromJSONFile(t, &blockVb, path.Join("../", testutils.TestDataPathBTC, "block_trimmed_8332_828440.json"))

	// https://mempool.space/block/000000000000000000025ca01d2c1094b8fd3bacc5468cc3193ced6a14618c27
	var blockMb testutils.MempoolBlock
	testutils.LoadObjectFromJSONFile(t, &blockMb, path.Join("../", testutils.TestDataPathBTC, "block_mempool.space_8332_828440.json"))

	gasRate, err := CalcBlockAvgFeeRate(&blockVb, &chaincfg.MainNetParams)
	require.NoError(t, err)
	require.Equal(t, int64(blockMb.Extras.AvgFeeRate), gasRate)
}

func TestAvgFeeRateBlock828440Errors(t *testing.T) {
	// load archived block 828440
	var blockVb btcjson.GetBlockVerboseTxResult
	testutils.LoadObjectFromJSONFile(t, &blockVb, path.Join("../", testutils.TestDataPathBTC, "block_trimmed_8332_828440.json"))

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
	testutils.LoadObjectFromJSONFile(t, &blockVb, path.Join("../", testutils.TestDataPathBTC, "block_trimmed_8332_828440.json"))
	avgGasRate := float64(32.0)
	// #nosec G701 test - always in range
	gasRate := int64(avgGasRate * clientcommon.BTCOutboundGasPriceMultiplier)
	dynamicFee828440 := DepositorFee(gasRate)

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

func TestCheckTSSVout(t *testing.T) {
	// the archived outtx raw result file and cctx file
	// https://blockstream.info/tx/030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0
	chain := chains.BtcMainnetChain
	chainID := chain.ChainId
	nonce := uint64(148)

	// create mainnet mock client
	btcClient := MockBTCClientMainnet()

	t.Run("valid TSS vout should pass", func(t *testing.T) {
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()
		err := btcClient.checkTSSVout(params, rawResult.Vout)
		require.NoError(t, err)
	})
	t.Run("should fail if vout length < 2 or > 3", func(t *testing.T) {
		_, cctx := testutils.LoadBTCTxRawResultNCctx(t, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		err := btcClient.checkTSSVout(params, []btcjson.Vout{{}})
		require.ErrorContains(t, err, "invalid number of vouts")

		err = btcClient.checkTSSVout(params, []btcjson.Vout{{}, {}, {}, {}})
		require.ErrorContains(t, err, "invalid number of vouts")
	})
	t.Run("should fail on invalid TSS vout", func(t *testing.T) {
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		// invalid TSS vout
		rawResult.Vout[0].ScriptPubKey.Hex = "invalid script"
		err := btcClient.checkTSSVout(params, rawResult.Vout)
		require.Error(t, err)
	})
	t.Run("should fail if vout 0 is not to the TSS address", func(t *testing.T) {
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		// not TSS address, bc1qh297vdt8xq6df5xae9z8gzd4jsu9a392mp0dus
		rawResult.Vout[0].ScriptPubKey.Hex = "0014ba8be635673034d4d0ddc9447409b594385ec4aa"
		err := btcClient.checkTSSVout(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match TSS address")
	})
	t.Run("should fail if vout 0 not match nonce mark", func(t *testing.T) {
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		// not match nonce mark
		rawResult.Vout[0].Value = 0.00000147
		err := btcClient.checkTSSVout(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match nonce-mark amount")
	})
	t.Run("should fail if vout 1 is not to the receiver address", func(t *testing.T) {
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		// not receiver address, bc1qh297vdt8xq6df5xae9z8gzd4jsu9a392mp0dus
		rawResult.Vout[1].ScriptPubKey.Hex = "0014ba8be635673034d4d0ddc9447409b594385ec4aa"
		err := btcClient.checkTSSVout(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match params receiver")
	})
	t.Run("should fail if vout 1 not match payment amount", func(t *testing.T) {
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		// not match payment amount
		rawResult.Vout[1].Value = 0.00011000
		err := btcClient.checkTSSVout(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match params amount")
	})
	t.Run("should fail if vout 2 is not to the TSS address", func(t *testing.T) {
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		// not TSS address, bc1qh297vdt8xq6df5xae9z8gzd4jsu9a392mp0dus
		rawResult.Vout[2].ScriptPubKey.Hex = "0014ba8be635673034d4d0ddc9447409b594385ec4aa"
		err := btcClient.checkTSSVout(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match TSS address")
	})
}

func TestCheckTSSVoutCancelled(t *testing.T) {
	// the archived outtx raw result file and cctx file
	// https://blockstream.info/tx/030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0
	chain := chains.BtcMainnetChain
	chainID := chain.ChainId
	nonce := uint64(148)

	// create mainnet mock client
	btcClient := MockBTCClientMainnet()

	t.Run("valid TSS vout should pass", func(t *testing.T) {
		// remove change vout to simulate cancelled tx
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, chainID, nonce)
		rawResult.Vout[1] = rawResult.Vout[2]
		rawResult.Vout = rawResult.Vout[:2]
		params := cctx.GetCurrentOutboundParam()

		err := btcClient.checkTSSVoutCancelled(params, rawResult.Vout)
		require.NoError(t, err)
	})
	t.Run("should fail if vout length < 1 or > 2", func(t *testing.T) {
		_, cctx := testutils.LoadBTCTxRawResultNCctx(t, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		err := btcClient.checkTSSVoutCancelled(params, []btcjson.Vout{})
		require.ErrorContains(t, err, "invalid number of vouts")

		err = btcClient.checkTSSVoutCancelled(params, []btcjson.Vout{{}, {}, {}})
		require.ErrorContains(t, err, "invalid number of vouts")
	})
	t.Run("should fail if vout 0 is not to the TSS address", func(t *testing.T) {
		// remove change vout to simulate cancelled tx
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, chainID, nonce)
		rawResult.Vout[1] = rawResult.Vout[2]
		rawResult.Vout = rawResult.Vout[:2]
		params := cctx.GetCurrentOutboundParam()

		// not TSS address, bc1qh297vdt8xq6df5xae9z8gzd4jsu9a392mp0dus
		rawResult.Vout[0].ScriptPubKey.Hex = "0014ba8be635673034d4d0ddc9447409b594385ec4aa"
		err := btcClient.checkTSSVoutCancelled(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match TSS address")
	})
	t.Run("should fail if vout 0 not match nonce mark", func(t *testing.T) {
		// remove change vout to simulate cancelled tx
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, chainID, nonce)
		rawResult.Vout[1] = rawResult.Vout[2]
		rawResult.Vout = rawResult.Vout[:2]
		params := cctx.GetCurrentOutboundParam()

		// not match nonce mark
		rawResult.Vout[0].Value = 0.00000147
		err := btcClient.checkTSSVoutCancelled(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match nonce-mark amount")
	})
	t.Run("should fail if vout 1 is not to the TSS address", func(t *testing.T) {
		// remove change vout to simulate cancelled tx
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, chainID, nonce)
		rawResult.Vout[1] = rawResult.Vout[2]
		rawResult.Vout[1].N = 1 // swap vout index
		rawResult.Vout = rawResult.Vout[:2]
		params := cctx.GetCurrentOutboundParam()

		// not TSS address, bc1qh297vdt8xq6df5xae9z8gzd4jsu9a392mp0dus
		rawResult.Vout[1].ScriptPubKey.Hex = "0014ba8be635673034d4d0ddc9447409b594385ec4aa"
		err := btcClient.checkTSSVoutCancelled(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match TSS address")
	})
}

func TestGetSenderAddressByVin(t *testing.T) {
	chain := chains.BtcMainnetChain
	net := &chaincfg.MainNetParams

	t.Run("should get sender address from P2TR tx", func(t *testing.T) {
		// vin from the archived P2TR tx
		// https://mempool.space/tx/3618e869f9e87863c0f1cc46dbbaa8b767b4a5d6d60b143c2c50af52b257e867
		txHash := "3618e869f9e87863c0f1cc46dbbaa8b767b4a5d6d60b143c2c50af52b257e867"
		rpcClient := createRPCClientAndLoadTx(t, chain.ChainId, txHash)

		// get sender address
		txVin := btcjson.Vin{Txid: txHash, Vout: 2}
		sender, err := GetSenderAddressByVin(rpcClient, txVin, net)
		require.NoError(t, err)
		require.Equal(t, "bc1px3peqcd60hk7wqyqk36697u9hzugq0pd5lzvney93yzzrqy4fkpq6cj7m3", sender)
	})
	t.Run("should get sender address from P2WSH tx", func(t *testing.T) {
		// vin from the archived P2WSH tx
		// https://mempool.space/tx/d13de30b0cc53b5c4702b184ae0a0b0f318feaea283185c1cddb8b341c27c016
		txHash := "d13de30b0cc53b5c4702b184ae0a0b0f318feaea283185c1cddb8b341c27c016"
		rpcClient := createRPCClientAndLoadTx(t, chain.ChainId, txHash)

		// get sender address
		txVin := btcjson.Vin{Txid: txHash, Vout: 0}
		sender, err := GetSenderAddressByVin(rpcClient, txVin, net)
		require.NoError(t, err)
		require.Equal(t, "bc1q79kmcyc706d6nh7tpzhnn8lzp76rp0tepph3hqwrhacqfcy4lwxqft0ppq", sender)
	})
	t.Run("should get sender address from P2WPKH tx", func(t *testing.T) {
		// vin from the archived P2WPKH tx
		// https://mempool.space/tx/c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697
		txHash := "c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697"
		rpcClient := createRPCClientAndLoadTx(t, chain.ChainId, txHash)

		// get sender address
		txVin := btcjson.Vin{Txid: txHash, Vout: 2}
		sender, err := GetSenderAddressByVin(rpcClient, txVin, net)
		require.NoError(t, err)
		require.Equal(t, "bc1q68kxnq52ahz5vd6c8czevsawu0ux9nfrzzrh6e", sender)
	})
	t.Run("should get sender address from P2SH tx", func(t *testing.T) {
		// vin from the archived P2SH tx
		// https://mempool.space/tx/211568441340fd5e10b1a8dcb211a18b9e853dbdf265ebb1c728f9b52813455a
		txHash := "211568441340fd5e10b1a8dcb211a18b9e853dbdf265ebb1c728f9b52813455a"
		rpcClient := createRPCClientAndLoadTx(t, chain.ChainId, txHash)

		// get sender address
		txVin := btcjson.Vin{Txid: txHash, Vout: 0}
		sender, err := GetSenderAddressByVin(rpcClient, txVin, net)
		require.NoError(t, err)
		require.Equal(t, "3MqRRSP76qxdVD9K4cfFnVtSLVwaaAjm3t", sender)
	})
	t.Run("should get sender address from P2PKH tx", func(t *testing.T) {
		// vin from the archived P2PKH tx
		// https://mempool.space/tx/781fc8d41b476dbceca283ebff9573fda52c8fdbba5e78152aeb4432286836a7
		txHash := "781fc8d41b476dbceca283ebff9573fda52c8fdbba5e78152aeb4432286836a7"
		rpcClient := createRPCClientAndLoadTx(t, chain.ChainId, txHash)

		// get sender address
		txVin := btcjson.Vin{Txid: txHash, Vout: 1}
		sender, err := GetSenderAddressByVin(rpcClient, txVin, net)
		require.NoError(t, err)
		require.Equal(t, "1ESQp1WQi7fzSpzCNs2oBTqaUBmNjLQLoV", sender)
	})
	t.Run("should get empty sender address on unknown script", func(t *testing.T) {
		// vin from the archived P2PKH tx
		// https://mempool.space/tx/781fc8d41b476dbceca283ebff9573fda52c8fdbba5e78152aeb4432286836a7
		txHash := "781fc8d41b476dbceca283ebff9573fda52c8fdbba5e78152aeb4432286836a7"
		nameMsgTx := path.Join("../", testutils.TestDataPathBTC, testutils.FileNameBTCMsgTx(chain.ChainId, txHash))
		var msgTx wire.MsgTx
		testutils.LoadObjectFromJSONFile(t, &msgTx, nameMsgTx)

		// modify script to unknown script
		msgTx.TxOut[1].PkScript = []byte{0x00, 0x01, 0x02, 0x03} // can be any invalid script bytes
		tx := btcutil.NewTx(&msgTx)

		// feed tx to mock rpc client
		rpcClient := stub.NewMockBTCRPCClient()
		rpcClient.WithRawTransaction(tx)

		// get sender address
		txVin := btcjson.Vin{Txid: txHash, Vout: 1}
		sender, err := GetSenderAddressByVin(rpcClient, txVin, net)
		require.NoError(t, err)
		require.Empty(t, sender)
	})
}

func TestGetSenderAddressByVinErrors(t *testing.T) {
	// https://mempool.space/tx/3618e869f9e87863c0f1cc46dbbaa8b767b4a5d6d60b143c2c50af52b257e867
	txHash := "3618e869f9e87863c0f1cc46dbbaa8b767b4a5d6d60b143c2c50af52b257e867"
	chain := chains.BtcMainnetChain
	net := &chaincfg.MainNetParams

	t.Run("should get sender address from P2TR tx", func(t *testing.T) {
		rpcClient := stub.NewMockBTCRPCClient()
		// use invalid tx hash
		txVin := btcjson.Vin{Txid: "invalid tx hash", Vout: 2}
		sender, err := GetSenderAddressByVin(rpcClient, txVin, net)
		require.Error(t, err)
		require.Empty(t, sender)
	})
	t.Run("should return error when RPC client fails to get raw tx", func(t *testing.T) {
		// create mock rpc client without preloaded tx
		rpcClient := stub.NewMockBTCRPCClient()
		txVin := btcjson.Vin{Txid: txHash, Vout: 2}
		sender, err := GetSenderAddressByVin(rpcClient, txVin, net)
		require.ErrorContains(t, err, "error getting raw transaction")
		require.Empty(t, sender)
	})
	t.Run("should return error on invalid output index", func(t *testing.T) {
		// create mock rpc client with preloaded tx
		rpcClient := createRPCClientAndLoadTx(t, chain.ChainId, txHash)
		// invalid output index
		txVin := btcjson.Vin{Txid: txHash, Vout: 3}
		sender, err := GetSenderAddressByVin(rpcClient, txVin, net)
		require.ErrorContains(t, err, "out of range")
		require.Empty(t, sender)
	})
}

func TestGetBtcEvent(t *testing.T) {
	// load archived intx P2WPKH raw result
	// https://mempool.space/tx/847139aa65aa4a5ee896375951cbf7417cfc8a4d6f277ec11f40cd87319f04aa
	txHash := "847139aa65aa4a5ee896375951cbf7417cfc8a4d6f277ec11f40cd87319f04aa"
	chain := chains.BtcMainnetChain

	// GetBtcEvent arguments
	tx := testutils.LoadBTCInboundRawResult(t, chain.ChainId, txHash, false)
	tssAddress := testutils.TSSAddressBTCMainnet
	blockNumber := uint64(835640)
	net := &chaincfg.MainNetParams
	// 2.992e-05, see avgFeeRate https://mempool.space/api/v1/blocks/835640
	depositorFee := DepositorFee(22 * clientcommon.BTCOutboundGasPriceMultiplier)

	// expected result
	memo, err := hex.DecodeString(tx.Vout[1].ScriptPubKey.Hex[4:])
	require.NoError(t, err)
	eventExpected := &BTCInboundEvent{
		FromAddress: "bc1q68kxnq52ahz5vd6c8czevsawu0ux9nfrzzrh6e",
		ToAddress:   tssAddress,
		Value:       tx.Vout[0].Value - depositorFee, // 7008 sataoshis
		MemoBytes:   memo,
		BlockNumber: blockNumber,
		TxHash:      tx.Txid,
	}

	t.Run("should get BTC intx event from P2WPKH sender", func(t *testing.T) {
		// https://mempool.space/tx/c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697
		preHash := "c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Vout = 2
		eventExpected.FromAddress = "bc1q68kxnq52ahz5vd6c8czevsawu0ux9nfrzzrh6e"
		// load previous raw tx so so mock rpc client can return it
		rpcClient := createRPCClientAndLoadTx(t, chain.ChainId, preHash)

		// get BTC event
		event, err := GetBtcEvent(rpcClient, *tx, tssAddress, blockNumber, log.Logger, net, depositorFee)
		require.NoError(t, err)
		require.Equal(t, eventExpected, event)
	})
	t.Run("should get BTC intx event from P2TR sender", func(t *testing.T) {
		// replace vin with a P2TR vin, so the sender address will change
		// https://mempool.space/tx/3618e869f9e87863c0f1cc46dbbaa8b767b4a5d6d60b143c2c50af52b257e867
		preHash := "3618e869f9e87863c0f1cc46dbbaa8b767b4a5d6d60b143c2c50af52b257e867"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Vout = 2
		eventExpected.FromAddress = "bc1px3peqcd60hk7wqyqk36697u9hzugq0pd5lzvney93yzzrqy4fkpq6cj7m3"
		// load previous raw tx so so mock rpc client can return it
		rpcClient := createRPCClientAndLoadTx(t, chain.ChainId, preHash)

		// get BTC event
		event, err := GetBtcEvent(rpcClient, *tx, tssAddress, blockNumber, log.Logger, net, depositorFee)
		require.NoError(t, err)
		require.Equal(t, eventExpected, event)
	})
	t.Run("should get BTC intx event from P2WSH sender", func(t *testing.T) {
		// replace vin with a P2WSH vin, so the sender address will change
		// https://mempool.space/tx/d13de30b0cc53b5c4702b184ae0a0b0f318feaea283185c1cddb8b341c27c016
		preHash := "d13de30b0cc53b5c4702b184ae0a0b0f318feaea283185c1cddb8b341c27c016"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Vout = 0
		eventExpected.FromAddress = "bc1q79kmcyc706d6nh7tpzhnn8lzp76rp0tepph3hqwrhacqfcy4lwxqft0ppq"
		// load previous raw tx so so mock rpc client can return it
		rpcClient := createRPCClientAndLoadTx(t, chain.ChainId, preHash)

		// get BTC event
		event, err := GetBtcEvent(rpcClient, *tx, tssAddress, blockNumber, log.Logger, net, depositorFee)
		require.NoError(t, err)
		require.Equal(t, eventExpected, event)
	})
	t.Run("should get BTC intx event from P2SH sender", func(t *testing.T) {
		// replace vin with a P2SH vin, so the sender address will change
		// https://mempool.space/tx/211568441340fd5e10b1a8dcb211a18b9e853dbdf265ebb1c728f9b52813455a
		preHash := "211568441340fd5e10b1a8dcb211a18b9e853dbdf265ebb1c728f9b52813455a"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Vout = 0
		eventExpected.FromAddress = "3MqRRSP76qxdVD9K4cfFnVtSLVwaaAjm3t"
		// load previous raw tx so so mock rpc client can return it
		rpcClient := createRPCClientAndLoadTx(t, chain.ChainId, preHash)

		// get BTC event
		event, err := GetBtcEvent(rpcClient, *tx, tssAddress, blockNumber, log.Logger, net, depositorFee)
		require.NoError(t, err)
		require.Equal(t, eventExpected, event)
	})
	t.Run("should get BTC intx event from P2PKH sender", func(t *testing.T) {
		// replace vin with a P2PKH vin, so the sender address will change
		// https://mempool.space/tx/781fc8d41b476dbceca283ebff9573fda52c8fdbba5e78152aeb4432286836a7
		preHash := "781fc8d41b476dbceca283ebff9573fda52c8fdbba5e78152aeb4432286836a7"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Vout = 1
		eventExpected.FromAddress = "1ESQp1WQi7fzSpzCNs2oBTqaUBmNjLQLoV"
		// load previous raw tx so so mock rpc client can return it
		rpcClient := createRPCClientAndLoadTx(t, chain.ChainId, preHash)

		// get BTC event
		event, err := GetBtcEvent(rpcClient, *tx, tssAddress, blockNumber, log.Logger, net, depositorFee)
		require.NoError(t, err)
		require.Equal(t, eventExpected, event)
	})
	t.Run("should skip tx if len(tx.Vout) < 2", func(t *testing.T) {
		// load tx and modify the tx to have only 1 vout
		tx := testutils.LoadBTCInboundRawResult(t, chain.ChainId, txHash, false)
		tx.Vout = tx.Vout[:1]

		// get BTC event
		rpcClient := stub.NewMockBTCRPCClient()
		event, err := GetBtcEvent(rpcClient, *tx, tssAddress, blockNumber, log.Logger, net, depositorFee)
		require.NoError(t, err)
		require.Nil(t, event)
	})
	t.Run("should skip tx if Vout[0] is not a P2WPKH output", func(t *testing.T) {
		// load tx
		rpcClient := stub.NewMockBTCRPCClient()
		tx := testutils.LoadBTCInboundRawResult(t, chain.ChainId, txHash, false)

		// modify the tx to have Vout[0] a P2SH output
		tx.Vout[0].ScriptPubKey.Hex = strings.Replace(tx.Vout[0].ScriptPubKey.Hex, "0014", "a914", 1)
		event, err := GetBtcEvent(rpcClient, *tx, tssAddress, blockNumber, log.Logger, net, depositorFee)
		require.NoError(t, err)
		require.Nil(t, event)

		// append 1 byte to script to make it longer than 22 bytes
		tx.Vout[0].ScriptPubKey.Hex = tx.Vout[0].ScriptPubKey.Hex + "00"
		event, err = GetBtcEvent(rpcClient, *tx, tssAddress, blockNumber, log.Logger, net, depositorFee)
		require.NoError(t, err)
		require.Nil(t, event)
	})
	t.Run("should skip tx if receiver address is not TSS address", func(t *testing.T) {
		// load tx and modify receiver address to any non-tss address: bc1qw8wrek2m7nlqldll66ajnwr9mh64syvkt67zlu
		tx := testutils.LoadBTCInboundRawResult(t, chain.ChainId, txHash, false)
		tx.Vout[0].ScriptPubKey.Hex = "001471dc3cd95bf4fe0fb7ffd6bb29b865ddf5581196"

		// get BTC event
		rpcClient := stub.NewMockBTCRPCClient()
		event, err := GetBtcEvent(rpcClient, *tx, tssAddress, blockNumber, log.Logger, net, depositorFee)
		require.NoError(t, err)
		require.Nil(t, event)
	})
	t.Run("should skip tx if amount is less than depositor fee", func(t *testing.T) {
		// load tx and modify amount to less than depositor fee
		tx := testutils.LoadBTCInboundRawResult(t, chain.ChainId, txHash, false)
		tx.Vout[0].Value = depositorFee - 1.0/1e8 // 1 satoshi less than depositor fee

		// get BTC event
		rpcClient := stub.NewMockBTCRPCClient()
		event, err := GetBtcEvent(rpcClient, *tx, tssAddress, blockNumber, log.Logger, net, depositorFee)
		require.NoError(t, err)
		require.Nil(t, event)
	})
	t.Run("should skip tx if 2nd vout is not OP_RETURN", func(t *testing.T) {
		// load tx and modify memo OP_RETURN to OP_1
		tx := testutils.LoadBTCInboundRawResult(t, chain.ChainId, txHash, false)
		tx.Vout[1].ScriptPubKey.Hex = strings.Replace(tx.Vout[1].ScriptPubKey.Hex, "6a", "51", 1)

		// get BTC event
		rpcClient := stub.NewMockBTCRPCClient()
		event, err := GetBtcEvent(rpcClient, *tx, tssAddress, blockNumber, log.Logger, net, depositorFee)
		require.NoError(t, err)
		require.Nil(t, event)
	})
	t.Run("should skip tx if memo decoding fails", func(t *testing.T) {
		// load tx and modify memo length to be 1 byte less than actual
		tx := testutils.LoadBTCInboundRawResult(t, chain.ChainId, txHash, false)
		tx.Vout[1].ScriptPubKey.Hex = strings.Replace(tx.Vout[1].ScriptPubKey.Hex, "6a14", "6a13", 1)

		// get BTC event
		rpcClient := stub.NewMockBTCRPCClient()
		event, err := GetBtcEvent(rpcClient, *tx, tssAddress, blockNumber, log.Logger, net, depositorFee)
		require.NoError(t, err)
		require.Nil(t, event)
	})
}

func TestGetBtcEventErrors(t *testing.T) {
	// load archived intx P2WPKH raw result
	// https://mempool.space/tx/847139aa65aa4a5ee896375951cbf7417cfc8a4d6f277ec11f40cd87319f04aa
	txHash := "847139aa65aa4a5ee896375951cbf7417cfc8a4d6f277ec11f40cd87319f04aa"
	chain := chains.BtcMainnetChain
	net := &chaincfg.MainNetParams
	tssAddress := testutils.TSSAddressBTCMainnet
	blockNumber := uint64(835640)
	depositorFee := DepositorFee(22 * clientcommon.BTCOutboundGasPriceMultiplier)

	t.Run("should return error on invalid Vout[0] script", func(t *testing.T) {
		// load tx and modify Vout[0] script to invalid script
		tx := testutils.LoadBTCInboundRawResult(t, chain.ChainId, txHash, false)
		tx.Vout[0].ScriptPubKey.Hex = "0014invalid000000000000000000000000000000000"

		// get BTC event
		rpcClient := stub.NewMockBTCRPCClient()
		event, err := GetBtcEvent(rpcClient, *tx, tssAddress, blockNumber, log.Logger, net, depositorFee)
		require.Error(t, err)
		require.Nil(t, event)
	})
	t.Run("should return error if len(tx.Vin) < 1", func(t *testing.T) {
		// load tx and remove vin
		tx := testutils.LoadBTCInboundRawResult(t, chain.ChainId, txHash, false)
		tx.Vin = nil

		// get BTC event
		rpcClient := stub.NewMockBTCRPCClient()
		event, err := GetBtcEvent(rpcClient, *tx, tssAddress, blockNumber, log.Logger, net, depositorFee)
		require.Error(t, err)
		require.Nil(t, event)
	})
	t.Run("should return error if RPC client fails to get raw tx", func(t *testing.T) {
		// load tx and leave rpc client without preloaded tx
		tx := testutils.LoadBTCInboundRawResult(t, chain.ChainId, txHash, false)
		rpcClient := stub.NewMockBTCRPCClient()

		// get BTC event
		event, err := GetBtcEvent(rpcClient, *tx, tssAddress, blockNumber, log.Logger, net, depositorFee)
		require.Error(t, err)
		require.Nil(t, event)
	})
}
