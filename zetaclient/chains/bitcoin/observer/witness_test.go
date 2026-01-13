package observer

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
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

		script := ParseScriptFromWitness(witness[:], log.Logger)
		require.NotNil(t, script)
		require.Equal(t, hex.EncodeToString(script), expected)
	})

	t.Run("no witness", func(t *testing.T) {
		witness := [0]string{}
		script := ParseScriptFromWitness(witness[:], log.Logger)
		require.Nil(t, script)
	})

	t.Run("ignore key spending path", func(t *testing.T) {
		witness := [1]string{
			"134896c42cd95680b048845847c8054756861ffab7d4abab72f6508d67d1ec0c590287ec2161dd7884983286e1cd56ce65c08a24ee0476ede92678a93b1b180c",
		}
		script := ParseScriptFromWitness(witness[:], log.Logger)
		require.Nil(t, script)
	})
}

func TestGetBtcEventWithWitness(t *testing.T) {
	ctx := context.Background()

	// load archived inbound P2WPKH raw result
	// https://mempool.space/tx/847139aa65aa4a5ee896375951cbf7417cfc8a4d6f277ec11f40cd87319f04aa
	txHash := "847139aa65aa4a5ee896375951cbf7417cfc8a4d6f277ec11f40cd87319f04aa"
	chain := chains.BitcoinMainnet

	tssAddress := testutils.TSSAddressBTCMainnet
	blockNumber := uint64(835640)
	net := &chaincfg.MainNetParams

	// fee rate of above tx is 28 sat/vB, apply gas rate multiplier to get depositor fee
	gasRateMultiplier := observertypes.DefaultBTCOutboundGasPriceMultiplier.MustFloat64()
	// #nosec G115 always in range
	depositorFee := common.DepositorFee(int64(28 * gasRateMultiplier))
	feeCalculator := mockDepositFeeCalculator(depositorFee, nil)

	t.Run("decode OP_RETURN ok", func(t *testing.T) {
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)

		// mock up the input
		// https://mempool.space/tx/c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697
		preHash := "c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Vout = 2

		sender := "bc1q68kxnq52ahz5vd6c8czevsawu0ux9nfrzzrh6e"
		memo, _ := hex.DecodeString(tx.Vout[1].ScriptPubKey.Hex[4:])
		eventExpected := &BTCInboundEvent{
			FromAddress:  sender,
			ToAddress:    tssAddress,
			Value:        tx.Vout[0].Value - depositorFee,
			DepositorFee: depositorFee,
			MemoBytes:    memo,
			BlockNumber:  blockNumber,
			TxHash:       tx.Txid,
		}

		// mock up rpc client to return sender address
		rpcClient := mocks.NewBitcoinClient(t)
		rpcClient.On("GetTransactionInputSpender", mock.Anything, preHash, mock.Anything, mock.Anything).
			Return(sender, nil)

		// get BTC event
		event, err := GetBtcEventWithWitness(
			ctx,
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			gasRateMultiplier,
			log.Logger,
			net,
			feeCalculator,
		)
		require.NoError(t, err)
		require.Equal(t, eventExpected, event)
	})

	t.Run("it's ok if no memo provided", func(t *testing.T) {
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)

		// mock up the input
		// https://mempool.space/tx/c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697
		preHash := "c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Vout = 2

		// mock up the output
		// remove OP_RETURN output to simulate no memo provided
		tx.Vout[1] = tx.Vout[2]
		tx.Vout = tx.Vout[:2]

		sender := "bc1q68kxnq52ahz5vd6c8czevsawu0ux9nfrzzrh6e"
		eventExpected := &BTCInboundEvent{
			FromAddress:  "bc1q68kxnq52ahz5vd6c8czevsawu0ux9nfrzzrh6e",
			ToAddress:    tssAddress,
			Value:        tx.Vout[0].Value - depositorFee,
			DepositorFee: depositorFee,
			MemoBytes:    []byte("no memo found"),
			BlockNumber:  blockNumber,
			TxHash:       tx.Txid,
		}

		// mock up rpc client to return sender address
		rpcClient := mocks.NewBitcoinClient(t)
		rpcClient.On("GetTransactionInputSpender", mock.Anything, preHash, mock.Anything, mock.Anything).
			Return(sender, nil)

		// get BTC event
		event, err := GetBtcEventWithWitness(
			ctx,
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			gasRateMultiplier,
			log.Logger,
			net,
			feeCalculator,
		)
		require.NoError(t, err)
		require.Equal(t, eventExpected, event)
	})

	t.Run("should return failed status if amount is less than depositor fee", func(t *testing.T) {
		// load tx and modify amount to less than depositor fee
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)
		tx.Vout[0].Value = depositorFee - 1.0/1e8 // 1 satoshi less than depositor fee

		// mock up the input
		// https://mempool.space/tx/c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697
		preHash := "c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Vout = 2

		sender := "bc1q68kxnq52ahz5vd6c8czevsawu0ux9nfrzzrh6e"
		memo, _ := hex.DecodeString(tx.Vout[1].ScriptPubKey.Hex[4:])
		depositedAmount := tx.Vout[0].Value
		eventExpected := &BTCInboundEvent{
			FromAddress:  "bc1q68kxnq52ahz5vd6c8czevsawu0ux9nfrzzrh6e",
			ToAddress:    tssAddress,
			Value:        0.0,
			DepositorFee: depositorFee,
			MemoBytes:    memo,
			BlockNumber:  blockNumber,
			TxHash:       tx.Txid,
			Status:       types.InboundStatus_INSUFFICIENT_DEPOSITOR_FEE,
			ErrorMessage: fmt.Sprintf(
				"deposited amount %v is less than depositor fee %v",
				depositedAmount,
				depositorFee,
			),
		}

		// mock up rpc client to return sender address
		rpcClient := mocks.NewBitcoinClient(t)
		rpcClient.On("GetTransactionInputSpender", mock.Anything, preHash, mock.Anything, mock.Anything).
			Return(sender, nil)

		// get BTC event
		event, err := GetBtcEventWithWitness(
			ctx,
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			gasRateMultiplier,
			log.Logger,
			net,
			feeCalculator,
		)
		require.NoError(t, err)
		require.Equal(t, eventExpected, event)
	})

	t.Run("decode inscription ok", func(t *testing.T) {
		txHash2 := "37777defed8717c581b4c0509329550e344bdc14ac38f71fc050096887e535c8"
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash2, false)

		// mock up the input
		// https://mempool.space/tx/c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697
		preHash := "c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Vout = 2

		sender := "bc1q68kxnq52ahz5vd6c8czevsawu0ux9nfrzzrh6e"
		eventExpected := &BTCInboundEvent{
			FromAddress:  sender,
			ToAddress:    tssAddress,
			Value:        tx.Vout[0].Value - depositorFee,
			DepositorFee: depositorFee,
			MemoBytes:    make([]byte, 600),
			BlockNumber:  blockNumber,
			TxHash:       tx.Txid,
		}

		// mock up rpc client to return sender address
		rpcClient := mocks.NewBitcoinClient(t)
		anyCommitAddress := "bc1qf0stfhtxnchsvgpnhhmmxlx9722sfatq02uhn4"
		rpcClient.On("GetTransactionInputSpender", mock.Anything, preHash, mock.Anything, mock.Anything).
			Return(anyCommitAddress, nil)
		rpcClient.On("GetTransactionInitiator", mock.Anything, preHash, mock.Anything, mock.Anything).
			Return(sender, nil)

		// get BTC event
		event, err := GetBtcEventWithWitness(
			ctx,
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			gasRateMultiplier,
			log.Logger,
			net,
			feeCalculator,
		)
		require.NoError(t, err)
		require.Equal(t, eventExpected, event)
	})

	t.Run("decode inscription ok - mainnet", func(t *testing.T) {
		// The input data is from the below mainnet, but output is modified for test case
		txHash2 := "7a57f987a3cb605896a5909d9ef2bf7afbf0c78f21e4118b85d00d9e4cce0c2c"
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash2, false)

		// mock up the input
		// https://mempool.space/tx/c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697
		preHash := "c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Vout = 2

		sender := "bc1q68kxnq52ahz5vd6c8czevsawu0ux9nfrzzrh6e"
		memo, _ := hex.DecodeString(
			"72f080c854647755d0d9e6f6821f6931f855b9acffd53d87433395672756d58822fd143360762109ab898626556b1c3b8d3096d2361f1297df4a41c1b429471a9aa2fc9be5f27c13b3863d6ac269e4b587d8389f8fd9649859935b0d48dea88cdb40f20c",
		)
		eventExpected := &BTCInboundEvent{
			FromAddress:  sender,
			ToAddress:    tssAddress,
			Value:        tx.Vout[0].Value - depositorFee,
			DepositorFee: depositorFee,
			MemoBytes:    memo,
			BlockNumber:  blockNumber,
			TxHash:       tx.Txid,
		}

		// mock up rpc client to return sender address
		rpcClient := mocks.NewBitcoinClient(t)
		anyCommitAddress := "bc1qf0stfhtxnchsvgpnhhmmxlx9722sfatq02uhn4"
		rpcClient.On("GetTransactionInputSpender", mock.Anything, preHash, mock.Anything, mock.Anything).
			Return(anyCommitAddress, nil)
		rpcClient.On("GetTransactionInitiator", mock.Anything, preHash, mock.Anything, mock.Anything).
			Return(sender, nil)

		// get BTC event
		event, err := GetBtcEventWithWitness(
			ctx,
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			gasRateMultiplier,
			log.Logger,
			net,
			feeCalculator,
		)
		require.NoError(t, err)
		require.Equal(t, event, eventExpected)
	})

	t.Run("should skip tx if Vout[0] is not a valid P2WPKH output", func(t *testing.T) {
		// load tx
		rpcClient := mocks.NewBitcoinClient(t)
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)

		// modify the tx to have Vout[0] a P2SH output
		tx.Vout[0].ScriptPubKey.Hex = strings.Replace(tx.Vout[0].ScriptPubKey.Hex, "0014", "a914", 1)
		event, err := GetBtcEventWithWitness(
			ctx,
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			gasRateMultiplier,
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
		rpcClient := mocks.NewBitcoinClient(t)
		event, err := GetBtcEventWithWitness(
			ctx,
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			gasRateMultiplier,
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

		// mock up the input
		// https://mempool.space/tx/c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697
		preHash := "c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Vout = 2

		// mock up rpc client to return sender address
		sender := "bc1q68kxnq52ahz5vd6c8czevsawu0ux9nfrzzrh6e"
		rpcClient := mocks.NewBitcoinClient(t)
		rpcClient.On("GetTransactionInputSpender", mock.Anything, preHash, mock.Anything, mock.Anything).
			Return(sender, nil)

		// get BTC event
		event, err := GetBtcEventWithWitness(
			ctx,
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			gasRateMultiplier,
			log.Logger,
			net,
			mockDepositFeeCalculator(0.0, errors.New("rpc error")),
		)
		require.ErrorContains(t, err, "rpc error")
		require.Nil(t, event)
	})

	t.Run("should return error if unable to get sender address", func(t *testing.T) {
		// load tx
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)

		// mock up rpc client to return rpc error
		rpcClient := mocks.NewBitcoinClient(t)
		rpcClient.On("GetTransactionInputSpender", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return("", errors.New("rpc error"))

		// get BTC event
		event, err := GetBtcEventWithWitness(
			ctx,
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			gasRateMultiplier,
			log.Logger,
			net,
			feeCalculator,
		)
		require.ErrorContains(t, err, "rpc error")
		require.Nil(t, event)
	})

	t.Run("should skip tx if sender address is empty", func(t *testing.T) {
		// load tx
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)

		// https://mempool.space/tx/c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697
		preVout := uint32(2)
		preHash := "c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Vout = preVout

		// create mock rpc client to return empty sender address
		rpcClient := mocks.NewBitcoinClient(t)
		rpcClient.On("GetTransactionInputSpender", mock.Anything, preHash, mock.Anything, mock.Anything).Return("", nil)

		// get BTC event
		event, err := GetBtcEventWithWitness(
			ctx,
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			gasRateMultiplier,
			log.Logger,
			net,
			feeCalculator,
		)
		require.NoError(t, err)
		require.Nil(t, event)
	})

	t.Run("should skip tx if sender address is TSS address", func(t *testing.T) {
		// load tx
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)

		// mock up the input
		// https://mempool.space/tx/c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697
		preHash := "c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Vout = 2

		// mock up rpc client to return TSS pkScript
		rpcClient := mocks.NewBitcoinClient(t)
		rpcClient.On("GetTransactionInputSpender", mock.Anything, preHash, mock.Anything, mock.Anything).
			Return(tssAddress, nil)

		// get BTC event
		event, err := GetBtcEventWithWitness(
			ctx,
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			gasRateMultiplier,
			log.Logger,
			net,
			feeCalculator,
		)
		require.NoError(t, err)
		require.Nil(t, event)
	})

	t.Run("should return failed status if amount is less than depositor fee", func(t *testing.T) {
		// load tx and modify amount to less than depositor fee
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)
		tx.Vout[0].Value = depositorFee - 1.0/1e8 // 1 satoshi less than depositor fee

		// mock up the input
		// https://mempool.space/tx/c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697
		preHash := "c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Vout = 2

		sender := "bc1q68kxnq52ahz5vd6c8czevsawu0ux9nfrzzrh6e"
		memo, _ := hex.DecodeString(tx.Vout[1].ScriptPubKey.Hex[4:])
		depositedAmount := tx.Vout[0].Value
		eventExpected := &BTCInboundEvent{
			FromAddress:  sender,
			ToAddress:    tssAddress,
			Value:        0.0,
			DepositorFee: depositorFee,
			MemoBytes:    memo,
			BlockNumber:  blockNumber,
			TxHash:       tx.Txid,
			Status:       types.InboundStatus_INSUFFICIENT_DEPOSITOR_FEE,
			ErrorMessage: fmt.Sprintf(
				"deposited amount %v is less than depositor fee %v",
				depositedAmount,
				depositorFee,
			),
		}

		// mock up rpc client to return sender address
		rpcClient := mocks.NewBitcoinClient(t)
		rpcClient.On("GetTransactionInputSpender", mock.Anything, preHash, mock.Anything, mock.Anything).
			Return(sender, nil)

		// get BTC event
		event, err := GetBtcEventWithWitness(
			ctx,
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			gasRateMultiplier,
			log.Logger,
			net,
			feeCalculator,
		)
		require.NoError(t, err)
		require.Equal(t, eventExpected, event)
	})

	t.Run("should return error if unable to get inscription initiator", func(t *testing.T) {
		txHash := "37777defed8717c581b4c0509329550e344bdc14ac38f71fc050096887e535c8"
		tx := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)

		// mock up the input
		// https://mempool.space/tx/c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697
		preHash := "c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697"
		tx.Vin[0].Txid = preHash
		tx.Vin[0].Vout = 2

		// mock up rpc client to return error when getting inscription initiator
		rpcClient := mocks.NewBitcoinClient(t)
		anyCommitAddress := "bc1qf0stfhtxnchsvgpnhhmmxlx9722sfatq02uhn4"
		rpcClient.On("GetTransactionInputSpender", mock.Anything, preHash, mock.Anything, mock.Anything).
			Return(anyCommitAddress, nil)
		rpcClient.On("GetTransactionInitiator", mock.Anything, preHash, mock.Anything, mock.Anything).
			Return("", errors.New("rpc error"))

		// get BTC event
		event, err := GetBtcEventWithWitness(
			ctx,
			rpcClient,
			*tx,
			tssAddress,
			blockNumber,
			gasRateMultiplier,
			log.Logger,
			net,
			feeCalculator,
		)
		require.ErrorContains(t, err, "rpc error")
		require.Nil(t, event)
	})
}
