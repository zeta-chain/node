package observer_test

import (
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	clientcommon "github.com/zeta-chain/node/zetaclient/common"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

func TestParseScriptFromWitness(t *testing.T) {
	t.Run("decode script ok", func(t *testing.T) {
		witness := [3]string{
			"3a4b32aef0e6ecc62d185594baf4df186c6d48ec15e72515bf81c1bcc1f04c758f4d54486bc2e7c280e649761d9084dbd2e7cdfb20708a7f8d0f82e5277bba2b",
			"20888269c4f0b7f6fe95d0cba364e2b1b879d9b00735d19cfab4b8d87096ce2b3cac00634d0802000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004c50000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000068",
			"c0888269c4f0b7f6fe95d0cba364e2b1b879d9b00735d19cfab4b8d87096ce2b3c",
		}
		expected := "20888269c4f0b7f6fe95d0cba364e2b1b879d9b00735d19cfab4b8d87096ce2b3cac00634d0802000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004c50000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000068"

		script := observer.ParseScriptFromWitness(witness[:], log.Logger)
		require.NotNil(t, script)
		require.Equal(t, hex.EncodeToString(script), expected)
	})

	t.Run("no witness", func(t *testing.T) {
		witness := [0]string{}
		script := observer.ParseScriptFromWitness(witness[:], log.Logger)
		require.Nil(t, script)
	})

	t.Run("ignore key spending path", func(t *testing.T) {
		witness := [1]string{
			"134896c42cd95680b048845847c8054756861ffab7d4abab72f6508d67d1ec0c590287ec2161dd7884983286e1cd56ce65c08a24ee0476ede92678a93b1b180c",
		}
		script := observer.ParseScriptFromWitness(witness[:], log.Logger)
		require.Nil(t, script)
	})
}

func TestGetBtcEventFromInscription(t *testing.T) {
	// load archived inbound P2WPKH raw result
	// https://mempool.space/tx/847139aa65aa4a5ee896375951cbf7417cfc8a4d6f277ec11f40cd87319f04aa
	txHash := "847139aa65aa4a5ee896375951cbf7417cfc8a4d6f277ec11f40cd87319f04aa"
	chain := chains.BitcoinMainnet

	tssAddress := testutils.TSSAddressBTCMainnet
	blockNumber := uint64(835640)
	net := &chaincfg.MainNetParams
	// 2.992e-05, see avgFeeRate https://mempool.space/api/v1/blocks/835640
	depositorFee := bitcoin.DepositorFee(22 * clientcommon.BTCOutboundGasPriceMultiplier)

	t.Run("decode OP_RETURN ok", func(t *testing.T) {
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)

		// https://mempool.space/tx/c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697
		preHash := "c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Vout = 2

		memo, _ := hex.DecodeString(tx.Vout[1].ScriptPubKey.Hex[4:])
		eventExpected := &observer.BTCInboundEvent{
			FromAddress: "bc1q68kxnq52ahz5vd6c8czevsawu0ux9nfrzzrh6e",
			ToAddress:   tssAddress,
			Value:       tx.Vout[0].Value - depositorFee,
			MemoBytes:   memo,
			BlockNumber: blockNumber,
			TxHash:      tx.Txid,
		}

		// load previous raw tx so so mock rpc client can return it
		rpcClient := createRPCClientAndLoadTx(t, chain.ChainId, preHash)

		// get BTC event
		event, err := observer.GetBtcEventWithWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			depositorFee,
		)
		require.NoError(t, err)
		require.Equal(t, eventExpected, event)
	})

	t.Run("decode inscription ok", func(t *testing.T) {
		txHash2 := "37777defed8717c581b4c0509329550e344bdc14ac38f71fc050096887e535c8"
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash2, false)

		preHash := "c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697"
		rpcClient := createRPCClientAndLoadTx(t, chain.ChainId, preHash)

		eventExpected := &observer.BTCInboundEvent{
			FromAddress: "bc1qm24wp577nk8aacckv8np465z3dvmu7ry45el6y",
			ToAddress:   tssAddress,
			Value:       tx.Vout[0].Value - depositorFee,
			MemoBytes:   make([]byte, 600),
			BlockNumber: blockNumber,
			TxHash:      tx.Txid,
		}

		// get BTC event
		event, err := observer.GetBtcEventWithWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			depositorFee,
		)
		require.NoError(t, err)
		require.Equal(t, event, eventExpected)
	})

	t.Run("decode inscription ok - mainnet", func(t *testing.T) {
		// The input data is from the below mainnet, but output is modified for test case
		txHash2 := "7a57f987a3cb605896a5909d9ef2bf7afbf0c78f21e4118b85d00d9e4cce0c2c"
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash2, false)

		preHash := "c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Sequence = 2
		rpcClient := createRPCClientAndLoadTx(t, chain.ChainId, preHash)

		memo, _ := hex.DecodeString(
			"72f080c854647755d0d9e6f6821f6931f855b9acffd53d87433395672756d58822fd143360762109ab898626556b1c3b8d3096d2361f1297df4a41c1b429471a9aa2fc9be5f27c13b3863d6ac269e4b587d8389f8fd9649859935b0d48dea88cdb40f20c",
		)
		eventExpected := &observer.BTCInboundEvent{
			FromAddress: "bc1qm24wp577nk8aacckv8np465z3dvmu7ry45el6y",
			ToAddress:   tssAddress,
			Value:       tx.Vout[0].Value - depositorFee,
			MemoBytes:   memo,
			BlockNumber: blockNumber,
			TxHash:      tx.Txid,
		}

		// get BTC event
		event, err := observer.GetBtcEventWithWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			depositorFee,
		)
		require.NoError(t, err)
		require.Equal(t, event, eventExpected)
	})

	t.Run("should skip tx if receiver address is not TSS address", func(t *testing.T) {
		// load tx and modify receiver address to any non-tss address: bc1qw8wrek2m7nlqldll66ajnwr9mh64syvkt67zlu
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)
		tx.Vout[0].ScriptPubKey.Hex = "001471dc3cd95bf4fe0fb7ffd6bb29b865ddf5581196"

		// get BTC event
		rpcClient := mocks.NewMockBTCRPCClient()
		event, err := observer.GetBtcEventWithWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			depositorFee,
		)
		require.NoError(t, err)
		require.Nil(t, event)
	})

	t.Run("should skip tx if amount is less than depositor fee", func(t *testing.T) {
		// load tx and modify amount to less than depositor fee
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)
		tx.Vout[0].Value = depositorFee - 1.0/1e8 // 1 satoshi less than depositor fee

		// get BTC event
		rpcClient := mocks.NewMockBTCRPCClient()
		event, err := observer.GetBtcEventWithWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			depositorFee,
		)
		require.NoError(t, err)
		require.Nil(t, event)
	})

	t.Run("should return error if RPC client fails to get raw tx", func(t *testing.T) {
		// load tx and leave rpc client without preloaded tx
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)
		rpcClient := mocks.NewMockBTCRPCClient()

		// get BTC event
		event, err := observer.GetBtcEventWithWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			depositorFee,
		)
		require.Error(t, err)
		require.Nil(t, event)
	})

	t.Run("should return error if RPC client fails to get raw tx", func(t *testing.T) {
		// load tx and leave rpc client without preloaded tx
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)
		rpcClient := mocks.NewMockBTCRPCClient()

		// get BTC event
		event, err := observer.GetBtcEventWithWitness(
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			log.Logger,
			net,
			depositorFee,
		)
		require.Error(t, err)
		require.Nil(t, event)
	})
}
