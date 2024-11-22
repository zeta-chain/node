package observer_test

import (
	"bytes"
	"encoding/hex"
	"math"
	"path"
	"strings"
	"testing"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/testutil"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	clientcommon "github.com/zeta-chain/node/zetaclient/common"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
	"github.com/zeta-chain/node/zetaclient/testutils/testrpc"
)

// mockDepositFeeCalculator returns a mock depositor fee calculator that returns the given fee and error.
func mockDepositFeeCalculator(fee float64, err error) bitcoin.DepositorFeeCalculator {
	return func(interfaces.BTCRPCClient, *btcjson.TxRawResult, *chaincfg.Params) (float64, error) {
		return fee, err
	}
}

func TestAvgFeeRateBlock828440(t *testing.T) {
	// load archived block 828440
	var blockVb btcjson.GetBlockVerboseTxResult
	testutils.LoadObjectFromJSONFile(
		t,
		&blockVb,
		path.Join(TestDataDir, testutils.TestDataPathBTC, "block_trimmed_8332_828440.json"),
	)

	// https://mempool.space/block/000000000000000000025ca01d2c1094b8fd3bacc5468cc3193ced6a14618c27
	var blockMb testutils.MempoolBlock
	testutils.LoadObjectFromJSONFile(
		t,
		&blockMb,
		path.Join(TestDataDir, testutils.TestDataPathBTC, "block_mempool.space_8332_828440.json"),
	)

	gasRate, err := bitcoin.CalcBlockAvgFeeRate(&blockVb, &chaincfg.MainNetParams)
	require.NoError(t, err)
	require.Equal(t, int64(blockMb.Extras.AvgFeeRate), gasRate)
}

func TestAvgFeeRateBlock828440Errors(t *testing.T) {
	// load archived block 828440
	var blockVb btcjson.GetBlockVerboseTxResult
	testutils.LoadObjectFromJSONFile(
		t,
		&blockVb,
		path.Join(TestDataDir, testutils.TestDataPathBTC, "block_trimmed_8332_828440.json"),
	)

	t.Run("block has no transactions", func(t *testing.T) {
		emptyVb := btcjson.GetBlockVerboseTxResult{Tx: []btcjson.TxRawResult{}}
		_, err := bitcoin.CalcBlockAvgFeeRate(&emptyVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "block has no transactions")
	})
	t.Run("it's okay if block has only coinbase tx", func(t *testing.T) {
		coinbaseVb := btcjson.GetBlockVerboseTxResult{Tx: []btcjson.TxRawResult{
			blockVb.Tx[0],
		}}
		_, err := bitcoin.CalcBlockAvgFeeRate(&coinbaseVb, &chaincfg.MainNetParams)
		require.NoError(t, err)
	})
	t.Run("tiny block weight should fail", func(t *testing.T) {
		invalidVb := blockVb
		invalidVb.Weight = 3
		_, err := bitcoin.CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "block weight 3 too small")
	})
	t.Run("block weight should not be less than coinbase tx weight", func(t *testing.T) {
		invalidVb := blockVb
		invalidVb.Weight = blockVb.Tx[0].Weight - 1
		_, err := bitcoin.CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "less than coinbase tx weight")
	})
	t.Run("invalid block height should fail", func(t *testing.T) {
		invalidVb := blockVb
		invalidVb.Height = 0
		_, err := bitcoin.CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid block height")

		invalidVb.Height = math.MaxInt32 + 1
		_, err = bitcoin.CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid block height")
	})
	t.Run("failed to decode coinbase tx", func(t *testing.T) {
		invalidVb := blockVb
		invalidVb.Tx = []btcjson.TxRawResult{blockVb.Tx[0], blockVb.Tx[1]}
		invalidVb.Tx[0].Hex = "invalid hex"
		_, err := bitcoin.CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to decode coinbase tx")
	})
	t.Run("1st tx is not coinbase", func(t *testing.T) {
		invalidVb := blockVb
		invalidVb.Tx = []btcjson.TxRawResult{blockVb.Tx[1], blockVb.Tx[0]}
		_, err := bitcoin.CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
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
		_, err = bitcoin.CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "less than subsidy")
	})
}

func Test_GetInboundVoteFromBtcEvent(t *testing.T) {
	// can use any bitcoin chain for testing
	chain := chains.BitcoinMainnet
	params := mocks.MockChainParams(chain.ChainId, 10)

	// create test observer
	ob := MockBTCObserver(t, chain, params, nil)
	zetacoreClient := mocks.NewZetacoreClient(t).WithKeys(&keys.Keys{}).WithZetaChain()
	ob.WithZetacoreClient(zetacoreClient)

	// test cases
	tests := []struct {
		name    string
		event   *observer.BTCInboundEvent
		nilVote bool
	}{
		{
			name: "should return vote for standard memo",
			event: &observer.BTCInboundEvent{
				FromAddress: sample.BtcAddressP2WPKH(t, &chaincfg.MainNetParams),
				// a deposit and call
				MemoBytes: testutil.HexToBytes(
					t,
					"5a0110032d07a9cbd57dcca3e2cf966c88bc874445b6e3b60d68656c6c6f207361746f736869",
				),
			},
		},
		{
			name: "should return vote for legacy memo",
			event: &observer.BTCInboundEvent{
				// raw address + payload
				MemoBytes: testutil.HexToBytes(t, "2d07a9cbd57dcca3e2cf966c88bc874445b6e3b668656c6c6f207361746f736869"),
			},
		},
		{
			name: "should return nil if unable to decode memo",
			event: &observer.BTCInboundEvent{
				// standard memo that carries payload only, receiver address is empty
				MemoBytes: testutil.HexToBytes(t, "5a0110020d68656c6c6f207361746f736869"),
			},
			nilVote: true,
		},
		{
			name: "should return nil on donation message",
			event: &observer.BTCInboundEvent{
				MemoBytes: []byte(constant.DonationMessage),
			},
			nilVote: true,
		},
		{
			name: "should return nil on invalid deposit value",
			event: &observer.BTCInboundEvent{
				Value:     -1, // invalid value
				MemoBytes: testutil.HexToBytes(t, "2d07a9cbd57dcca3e2cf966c88bc874445b6e3b668656c6c6f207361746f736869"),
			},
			nilVote: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := ob.GetInboundVoteFromBtcEvent(tt.event)
			if tt.nilVote {
				require.Nil(t, msg)
			} else {
				require.NotNil(t, msg)
			}
		})
	}
}

func TestGetSenderAddressByVin(t *testing.T) {
	// https://mempool.space/tx/3618e869f9e87863c0f1cc46dbbaa8b767b4a5d6d60b143c2c50af52b257e867
	txHash := "3618e869f9e87863c0f1cc46dbbaa8b767b4a5d6d60b143c2c50af52b257e867"
	chain := chains.BitcoinMainnet
	net := &chaincfg.MainNetParams

	t.Run("should get sender address from tx", func(t *testing.T) {
		// vin from the archived P2WPKH tx
		// https://mempool.space/tx/c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697
		txHash := "c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697"
		rpcClient := testrpc.CreateBTCRPCAndLoadTx(t, TestDataDir, chain.ChainId, txHash)

		// get sender address
		txVin := btcjson.Vin{Txid: txHash, Vout: 2}
		sender, err := observer.GetSenderAddressByVin(rpcClient, txVin, net)
		require.NoError(t, err)
		require.Equal(t, "bc1q68kxnq52ahz5vd6c8czevsawu0ux9nfrzzrh6e", sender)
	})

	t.Run("should return error on invalid txHash", func(t *testing.T) {
		rpcClient := mocks.NewBTCRPCClient(t)
		// use invalid tx hash
		txVin := btcjson.Vin{Txid: "invalid tx hash", Vout: 2}
		sender, err := observer.GetSenderAddressByVin(rpcClient, txVin, net)
		require.Error(t, err)
		require.Empty(t, sender)
	})

	t.Run("should return error when RPC client fails to get raw tx", func(t *testing.T) {
		// create mock rpc client that returns rpc error
		rpcClient := mocks.NewBTCRPCClient(t)
		rpcClient.On("GetRawTransaction", mock.Anything).Return(nil, errors.New("rpc error"))

		// get sender address
		txVin := btcjson.Vin{Txid: txHash, Vout: 2}
		sender, err := observer.GetSenderAddressByVin(rpcClient, txVin, net)
		require.ErrorContains(t, err, "error getting raw transaction")
		require.Empty(t, sender)
	})

	t.Run("should return error on invalid output index", func(t *testing.T) {
		// create mock rpc client with preloaded tx
		rpcClient := testrpc.CreateBTCRPCAndLoadTx(t, TestDataDir, chain.ChainId, txHash)
		// invalid output index
		txVin := btcjson.Vin{Txid: txHash, Vout: 3}
		sender, err := observer.GetSenderAddressByVin(rpcClient, txVin, net)
		require.ErrorContains(t, err, "out of range")
		require.Empty(t, sender)
	})
}

func TestGetBtcEventWithoutWitness(t *testing.T) {
	// load archived inbound P2WPKH raw result
	// https://mempool.space/tx/847139aa65aa4a5ee896375951cbf7417cfc8a4d6f277ec11f40cd87319f04aa
	txHash := "847139aa65aa4a5ee896375951cbf7417cfc8a4d6f277ec11f40cd87319f04aa"
	chain := chains.BitcoinMainnet

	// GetBtcEventWithoutWitness arguments
	tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)
	tssAddress := testutils.TSSAddressBTCMainnet
	blockNumber := uint64(835640)
	net := &chaincfg.MainNetParams

	// fee rate of above tx is 28 sat/vB
	depositorFee := bitcoin.DepositorFee(28 * clientcommon.BTCOutboundGasPriceMultiplier)
	feeCalculator := mockDepositFeeCalculator(depositorFee, nil)

	// expected result
	memo, err := hex.DecodeString(tx.Vout[1].ScriptPubKey.Hex[4:])
	require.NoError(t, err)
	eventExpected := &observer.BTCInboundEvent{
		FromAddress:  "bc1q68kxnq52ahz5vd6c8czevsawu0ux9nfrzzrh6e",
		ToAddress:    tssAddress,
		Value:        tx.Vout[0].Value - depositorFee, // 6192 sataoshis
		DepositorFee: depositorFee,
		MemoBytes:    memo,
		BlockNumber:  blockNumber,
		TxHash:       tx.Txid,
	}

	t.Run("should get BTC inbound event from P2WPKH sender", func(t *testing.T) {
		// https://mempool.space/tx/c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697
		preHash := "c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Vout = 2
		eventExpected.FromAddress = "bc1q68kxnq52ahz5vd6c8czevsawu0ux9nfrzzrh6e"
		// load previous raw tx so so mock rpc client can return it
		rpcClient := testrpc.CreateBTCRPCAndLoadTx(t, TestDataDir, chain.ChainId, preHash)

		// get BTC event
		event, err := observer.GetBtcEventWithoutWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			feeCalculator,
		)
		require.NoError(t, err)
		require.Equal(t, eventExpected, event)
	})

	t.Run("should get BTC inbound event from P2TR sender", func(t *testing.T) {
		// replace vin with a P2TR vin, so the sender address will change
		// https://mempool.space/tx/3618e869f9e87863c0f1cc46dbbaa8b767b4a5d6d60b143c2c50af52b257e867
		preHash := "3618e869f9e87863c0f1cc46dbbaa8b767b4a5d6d60b143c2c50af52b257e867"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Vout = 2
		eventExpected.FromAddress = "bc1px3peqcd60hk7wqyqk36697u9hzugq0pd5lzvney93yzzrqy4fkpq6cj7m3"
		// load previous raw tx so so mock rpc client can return it
		rpcClient := testrpc.CreateBTCRPCAndLoadTx(t, TestDataDir, chain.ChainId, preHash)

		// get BTC event
		event, err := observer.GetBtcEventWithoutWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			feeCalculator,
		)
		require.NoError(t, err)
		require.Equal(t, eventExpected, event)
	})

	t.Run("should get BTC inbound event from P2WSH sender", func(t *testing.T) {
		// replace vin with a P2WSH vin, so the sender address will change
		// https://mempool.space/tx/d13de30b0cc53b5c4702b184ae0a0b0f318feaea283185c1cddb8b341c27c016
		preHash := "d13de30b0cc53b5c4702b184ae0a0b0f318feaea283185c1cddb8b341c27c016"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Vout = 0
		eventExpected.FromAddress = "bc1q79kmcyc706d6nh7tpzhnn8lzp76rp0tepph3hqwrhacqfcy4lwxqft0ppq"
		// load previous raw tx so so mock rpc client can return it
		rpcClient := testrpc.CreateBTCRPCAndLoadTx(t, TestDataDir, chain.ChainId, preHash)

		// get BTC event
		event, err := observer.GetBtcEventWithoutWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			feeCalculator,
		)
		require.NoError(t, err)
		require.Equal(t, eventExpected, event)
	})

	t.Run("should get BTC inbound event from P2SH sender", func(t *testing.T) {
		// replace vin with a P2SH vin, so the sender address will change
		// https://mempool.space/tx/211568441340fd5e10b1a8dcb211a18b9e853dbdf265ebb1c728f9b52813455a
		preHash := "211568441340fd5e10b1a8dcb211a18b9e853dbdf265ebb1c728f9b52813455a"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Vout = 0
		eventExpected.FromAddress = "3MqRRSP76qxdVD9K4cfFnVtSLVwaaAjm3t"
		// load previous raw tx so so mock rpc client can return it
		rpcClient := testrpc.CreateBTCRPCAndLoadTx(t, TestDataDir, chain.ChainId, preHash)

		// get BTC event
		event, err := observer.GetBtcEventWithoutWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			feeCalculator,
		)
		require.NoError(t, err)
		require.Equal(t, eventExpected, event)
	})

	t.Run("should get BTC inbound event from P2PKH sender", func(t *testing.T) {
		// replace vin with a P2PKH vin, so the sender address will change
		// https://mempool.space/tx/781fc8d41b476dbceca283ebff9573fda52c8fdbba5e78152aeb4432286836a7
		preHash := "781fc8d41b476dbceca283ebff9573fda52c8fdbba5e78152aeb4432286836a7"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Vout = 1
		eventExpected.FromAddress = "1ESQp1WQi7fzSpzCNs2oBTqaUBmNjLQLoV"
		// load previous raw tx so so mock rpc client can return it
		rpcClient := testrpc.CreateBTCRPCAndLoadTx(t, TestDataDir, chain.ChainId, preHash)

		// get BTC event
		event, err := observer.GetBtcEventWithoutWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			feeCalculator,
		)
		require.NoError(t, err)
		require.Equal(t, eventExpected, event)
	})

	t.Run("should skip tx if len(tx.Vout) < 2", func(t *testing.T) {
		// load tx and modify the tx to have only 1 vout
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)
		tx.Vout = tx.Vout[:1]

		// get BTC event
		rpcClient := mocks.NewBTCRPCClient(t)
		event, err := observer.GetBtcEventWithoutWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			feeCalculator,
		)
		require.NoError(t, err)
		require.Nil(t, event)
	})

	t.Run("should skip tx if Vout[0] is not a P2WPKH output", func(t *testing.T) {
		// load tx
		rpcClient := mocks.NewBTCRPCClient(t)
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)

		// modify the tx to have Vout[0] a P2SH output
		tx.Vout[0].ScriptPubKey.Hex = strings.Replace(tx.Vout[0].ScriptPubKey.Hex, "0014", "a914", 1)
		event, err := observer.GetBtcEventWithoutWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			feeCalculator,
		)
		require.NoError(t, err)
		require.Nil(t, event)

		// append 1 byte to script to make it longer than 22 bytes
		tx.Vout[0].ScriptPubKey.Hex = tx.Vout[0].ScriptPubKey.Hex + "00"
		event, err = observer.GetBtcEventWithoutWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			feeCalculator,
		)
		require.NoError(t, err)
		require.Nil(t, event)
	})

	t.Run("should skip tx if receiver address is not TSS address", func(t *testing.T) {
		// load tx and modify receiver address to any non-tss address: bc1qw8wrek2m7nlqldll66ajnwr9mh64syvkt67zlu
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)
		tx.Vout[0].ScriptPubKey.Hex = "001471dc3cd95bf4fe0fb7ffd6bb29b865ddf5581196"

		// get BTC event
		rpcClient := mocks.NewBTCRPCClient(t)
		event, err := observer.GetBtcEventWithoutWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			feeCalculator,
		)
		require.NoError(t, err)
		require.Nil(t, event)
	})

	t.Run("should return error if RPC failed to calculate depositor fee", func(t *testing.T) {
		// load tx
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)

		// get BTC event
		rpcClient := mocks.NewBTCRPCClient(t)
		event, err := observer.GetBtcEventWithoutWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			mockDepositFeeCalculator(0.0, errors.New("rpc error")),
		)
		require.ErrorContains(t, err, "rpc error")
		require.Nil(t, event)
	})

	t.Run("should skip tx if amount is less than depositor fee", func(t *testing.T) {
		// load tx and modify amount to less than depositor fee
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)
		tx.Vout[0].Value = depositorFee - 1.0/1e8 // 1 satoshi less than depositor fee

		// get BTC event
		rpcClient := mocks.NewBTCRPCClient(t)
		event, err := observer.GetBtcEventWithoutWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			feeCalculator,
		)
		require.NoError(t, err)
		require.Nil(t, event)
	})

	t.Run("should skip tx if 2nd vout is not OP_RETURN", func(t *testing.T) {
		// load tx and modify memo OP_RETURN to OP_1
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)
		tx.Vout[1].ScriptPubKey.Hex = strings.Replace(tx.Vout[1].ScriptPubKey.Hex, "6a", "51", 1)

		// get BTC event
		rpcClient := mocks.NewBTCRPCClient(t)
		event, err := observer.GetBtcEventWithoutWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			feeCalculator,
		)
		require.NoError(t, err)
		require.Nil(t, event)
	})

	t.Run("should skip tx if memo decoding fails", func(t *testing.T) {
		// load tx and modify memo length to be 1 byte less than actual
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)
		tx.Vout[1].ScriptPubKey.Hex = strings.Replace(tx.Vout[1].ScriptPubKey.Hex, "6a14", "6a13", 1)

		// get BTC event
		rpcClient := mocks.NewBTCRPCClient(t)
		event, err := observer.GetBtcEventWithoutWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			feeCalculator,
		)
		require.NoError(t, err)
		require.Nil(t, event)
	})

	t.Run("should skip tx if sender address is empty", func(t *testing.T) {
		// https://mempool.space/tx/c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697
		preVout := uint32(2)
		preHash := "c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Vout = preVout

		// create mock rpc client
		rpcClient := mocks.NewBTCRPCClient(t)

		// load archived MsgTx and modify previous input script to invalid
		msgTx := testutils.LoadBTCMsgTx(t, TestDataDir, chain.ChainId, preHash)
		msgTx.TxOut[preVout].PkScript = []byte{0x00, 0x01}

		// mock rpc response to return invalid tx msg
		rpcClient.On("GetRawTransaction", mock.Anything).Return(btcutil.NewTx(msgTx), nil)

		// get BTC event
		event, err := observer.GetBtcEventWithoutWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			feeCalculator,
		)
		require.NoError(t, err)
		require.Nil(t, event)
	})
}

func TestGetBtcEventErrors(t *testing.T) {
	// load archived inbound P2WPKH raw result
	// https://mempool.space/tx/847139aa65aa4a5ee896375951cbf7417cfc8a4d6f277ec11f40cd87319f04aa
	txHash := "847139aa65aa4a5ee896375951cbf7417cfc8a4d6f277ec11f40cd87319f04aa"
	chain := chains.BitcoinMainnet
	net := &chaincfg.MainNetParams
	tssAddress := testutils.TSSAddressBTCMainnet
	blockNumber := uint64(835640)

	// fee rate of above tx is 28 sat/vB
	depositorFee := bitcoin.DepositorFee(28 * clientcommon.BTCOutboundGasPriceMultiplier)
	feeCalculator := mockDepositFeeCalculator(depositorFee, nil)

	t.Run("should return error on invalid Vout[0] script", func(t *testing.T) {
		// load tx and modify Vout[0] script to invalid script
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)
		tx.Vout[0].ScriptPubKey.Hex = "0014invalid000000000000000000000000000000000"

		// get BTC event
		rpcClient := mocks.NewBTCRPCClient(t)
		event, err := observer.GetBtcEventWithoutWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			feeCalculator,
		)
		require.Error(t, err)
		require.Nil(t, event)
	})

	t.Run("should return error if len(tx.Vin) < 1", func(t *testing.T) {
		// load tx and remove vin
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)
		tx.Vin = nil

		// get BTC event
		rpcClient := mocks.NewBTCRPCClient(t)
		event, err := observer.GetBtcEventWithoutWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			feeCalculator,
		)
		require.ErrorContains(t, err, "no input found")
		require.Nil(t, event)
	})

	t.Run("should return error if RPC client fails to get raw tx", func(t *testing.T) {
		// load tx and leave rpc client without preloaded tx
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)

		// create mock rpc client that returns rpc error
		rpcClient := mocks.NewBTCRPCClient(t)
		rpcClient.On("GetRawTransaction", mock.Anything).Return(nil, errors.New("rpc error"))

		// get BTC event
		event, err := observer.GetBtcEventWithoutWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			feeCalculator,
		)
		require.ErrorContains(t, err, "error getting sender address")
		require.Nil(t, event)
	})
}

func TestGetBtcEvent(t *testing.T) {
	t.Run("should not decode inbound event with witness with mainnet chain", func(t *testing.T) {
		// load archived inbound P2WPKH raw result
		// https://mempool.space/tx/847139aa65aa4a5ee896375951cbf7417cfc8a4d6f277ec11f40cd87319f04aa
		chain := chains.BitcoinMainnet
		tssAddress := testutils.TSSAddressBTCMainnet
		blockNumber := uint64(835640)
		net := &chaincfg.MainNetParams
		// 2.992e-05, see avgFeeRate https://mempool.space/api/v1/blocks/835640
		depositorFee := bitcoin.DepositorFee(22 * clientcommon.BTCOutboundGasPriceMultiplier)
		feeCalculator := mockDepositFeeCalculator(depositorFee, nil)

		txHash2 := "37777defed8717c581b4c0509329550e344bdc14ac38f71fc050096887e535c8"
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash2, false)
		rpcClient := mocks.NewBTCRPCClient(t)
		// get BTC event
		event, err := observer.GetBtcEvent(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			feeCalculator,
		)
		require.NoError(t, err)
		require.Equal(t, (*observer.BTCInboundEvent)(nil), event)
	})

	t.Run("should support legacy BTC inbound event parsing for mainnet", func(t *testing.T) {
		// load archived inbound P2WPKH raw result
		// https://mempool.space/tx/847139aa65aa4a5ee896375951cbf7417cfc8a4d6f277ec11f40cd87319f04aa
		txHash := "847139aa65aa4a5ee896375951cbf7417cfc8a4d6f277ec11f40cd87319f04aa"
		chain := chains.BitcoinMainnet

		// GetBtcEventWithoutWitness arguments
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)
		tssAddress := testutils.TSSAddressBTCMainnet
		blockNumber := uint64(835640)
		net := &chaincfg.MainNetParams

		// fee rate of above tx is 28 sat/vB
		depositorFee := bitcoin.DepositorFee(28 * clientcommon.BTCOutboundGasPriceMultiplier)
		feeCalculator := mockDepositFeeCalculator(depositorFee, nil)

		// expected result
		memo, err := hex.DecodeString(tx.Vout[1].ScriptPubKey.Hex[4:])
		require.NoError(t, err)
		eventExpected := &observer.BTCInboundEvent{
			FromAddress:  "bc1q68kxnq52ahz5vd6c8czevsawu0ux9nfrzzrh6e",
			ToAddress:    tssAddress,
			Value:        tx.Vout[0].Value - depositorFee, // 6192 sataoshis
			DepositorFee: depositorFee,
			MemoBytes:    memo,
			BlockNumber:  blockNumber,
			TxHash:       tx.Txid,
		}

		// https://mempool.space/tx/c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697
		preHash := "c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Vout = 2
		eventExpected.FromAddress = "bc1q68kxnq52ahz5vd6c8czevsawu0ux9nfrzzrh6e"
		// load previous raw tx so so mock rpc client can return it
		rpcClient := testrpc.CreateBTCRPCAndLoadTx(t, TestDataDir, chain.ChainId, preHash)

		// get BTC event
		event, err := observer.GetBtcEvent(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			feeCalculator,
		)
		require.NoError(t, err)
		require.Equal(t, eventExpected, event)
	})
}
