package observer

import (
	"bytes"
	"context"
	"encoding/hex"
	"math"
	"path"
	"testing"

	cosmosmath "cosmossdk.io/math"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/memo"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/testutil"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/testutils"
)

// mockDepositFeeCalculator returns a mock depositor fee calculator that returns the given fee and error.
func mockDepositFeeCalculator(fee float64, err error) common.DepositorFeeCalculator {
	return func(_ context.Context, _ common.RPC, _ *btcjson.TxRawResult, _ *chaincfg.Params) (float64, error) {
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

	gasRate, err := common.CalcBlockAvgFeeRate(&blockVb, &chaincfg.MainNetParams)
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
		_, err := common.CalcBlockAvgFeeRate(&emptyVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "block has no transactions")
	})
	t.Run("it's okay if block has only coinbase tx", func(t *testing.T) {
		coinbaseVb := btcjson.GetBlockVerboseTxResult{Tx: []btcjson.TxRawResult{
			blockVb.Tx[0],
		}}
		_, err := common.CalcBlockAvgFeeRate(&coinbaseVb, &chaincfg.MainNetParams)
		require.NoError(t, err)
	})
	t.Run("tiny block weight should fail", func(t *testing.T) {
		invalidVb := blockVb
		invalidVb.Weight = 3
		_, err := common.CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "block weight 3 too small")
	})
	t.Run("block weight should not be less than coinbase tx weight", func(t *testing.T) {
		invalidVb := blockVb
		invalidVb.Weight = blockVb.Tx[0].Weight - 1
		_, err := common.CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "less than coinbase tx weight")
	})
	t.Run("invalid block height should fail", func(t *testing.T) {
		invalidVb := blockVb
		invalidVb.Height = 0
		_, err := common.CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid block height")

		invalidVb.Height = math.MaxInt32 + 1
		_, err = common.CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid block height")
	})
	t.Run("failed to decode coinbase tx", func(t *testing.T) {
		invalidVb := blockVb
		invalidVb.Tx = []btcjson.TxRawResult{blockVb.Tx[0], blockVb.Tx[1]}
		invalidVb.Tx[0].Hex = "invalid hex"
		_, err := common.CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to decode coinbase tx")
	})
	t.Run("1st tx is not coinbase", func(t *testing.T) {
		invalidVb := blockVb
		invalidVb.Tx = []btcjson.TxRawResult{blockVb.Tx[1], blockVb.Tx[0]}
		_, err := common.CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
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
		_, err = common.CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "less than subsidy")
	})
}

func Test_GetInboundVoteFromBtcEvent(t *testing.T) {
	r := sample.Rand()

	// can use any bitcoin chain for testing
	chain := chains.BitcoinMainnet

	// create test observer
	ob := newTestSuite(t, chain)
	ob.zetacore.WithKeys(&keys.Keys{}).WithZetaChain()

	// test cases
	tests := []struct {
		name              string
		event             *BTCInboundEvent
		observationStatus crosschaintypes.InboundStatus
		errorMessage      string
		nilVote           bool
	}{
		{
			name: "should return vote for standard memo",
			event: &BTCInboundEvent{
				FromAddress: sample.BTCAddressP2WPKH(t, r, &chaincfg.MainNetParams).String(),
				// a deposit and call
				MemoBytes: testutil.HexToBytes(
					t,
					"5a0110032d07a9cbd57dcca3e2cf966c88bc874445b6e3b60d68656c6c6f207361746f736869",
				),
			},
			observationStatus: crosschaintypes.InboundStatus_SUCCESS,
		},
		{
			name: "should return vote for legacy memo",
			event: &BTCInboundEvent{
				// raw address + payload
				MemoBytes: testutil.HexToBytes(t, "2d07a9cbd57dcca3e2cf966c88bc874445b6e3b668656c6c6f207361746f736869"),
			},
			observationStatus: crosschaintypes.InboundStatus_SUCCESS,
		},
		{
			name: "should return vote for invalid memo",
			event: &BTCInboundEvent{
				// standard memo that carries payload only, receiver address flag is NOT set
				MemoBytes: testutil.HexToBytes(t, "5a0110020d68656c6c6f207361746f736869"),
			},
			observationStatus: crosschaintypes.InboundStatus_INVALID_MEMO,
			errorMessage:      "must set receiver address flag",
		},
		{
			name: "should return vote for invalid legacy memo",
			event: &BTCInboundEvent{
				// only 19 bytes
				MemoBytes: sample.EthAddress().Bytes()[:19],
			},
			observationStatus: crosschaintypes.InboundStatus_INVALID_MEMO,
			errorMessage:      "legacy memo length must be at least 20 bytes",
		},
		{
			name: "should return nil on donation message",
			event: &BTCInboundEvent{
				MemoBytes: []byte(constant.DonationMessage),
			},
			nilVote: true,
		},
		{
			name: "should return nil on invalid deposit value",
			event: &BTCInboundEvent{
				Value:     21000001, // invalid value
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
				require.EqualValues(t, tt.observationStatus, msg.Status)
				require.Contains(t, msg.ErrorMessage, tt.errorMessage)
			}
		})
	}
}

func Test_NewInboundVoteFromEvent_LegacyMemo(t *testing.T) {
	// can use any bitcoin chain for testing
	chain := chains.BitcoinMainnet

	// create test observer
	ob := newTestSuite(t, chain)
	ob.zetacore.WithKeys(&keys.Keys{}).WithZetaChain()

	t.Run("should create new inbound vote msg V2", func(t *testing.T) {
		// create test event
		event := createTestBtcEvent(t, &chaincfg.MainNetParams, []byte("dummy memo"), nil)

		// given receiver and amount
		receiver := sample.EthAddress()
		amountSats := cosmosmath.NewUint(1000)
		event.ToAddress = receiver.Hex()
		event.MsgVoteAmount = amountSats

		// mock SAFE confirmed block
		ob.WithLastBlock(event.BlockNumber + ob.ChainParams().InboundConfirmationSafe())

		// expected vote
		expectedVote := crosschaintypes.MsgVoteInbound{
			Sender:             event.FromAddress,
			SenderChainId:      chain.ChainId,
			TxOrigin:           event.FromAddress,
			Receiver:           event.ToAddress,
			ReceiverChain:      ob.ZetacoreClient().Chain().ChainId,
			Amount:             amountSats,
			Message:            hex.EncodeToString(event.MemoBytes),
			InboundHash:        event.TxHash,
			InboundBlockHeight: event.BlockNumber,
			CallOptions: &crosschaintypes.CallOptions{
				GasLimit: 0,
			},
			CoinType:                coin.CoinType_Gas,
			ProtocolContractVersion: crosschaintypes.ProtocolContractVersion_V2,
			RevertOptions:           crosschaintypes.NewEmptyRevertOptions(), // always empty with legacy memo
			IsCrossChainCall:        true,
			Status:                  crosschaintypes.InboundStatus_SUCCESS,
			ConfirmationMode:        crosschaintypes.ConfirmationMode_SAFE,
		}

		// create new inbound vote V2 with legacy memo
		vote := ob.NewInboundVoteFromEvent(&event)
		require.Equal(t, expectedVote, *vote)
	})
}

func Test_NewInboundVoteFromEvent_StdMemo(t *testing.T) {
	// can use any bitcoin chain for testing
	chain := chains.BitcoinMainnet

	// create test observer
	ob := newTestSuite(t, chain)
	ob.zetacore.WithKeys(&keys.Keys{}).WithZetaChain()

	t.Run("should create new inbound vote msg with standard memo", func(t *testing.T) {
		// create revert options
		r := sample.Rand()
		revertOptions := crosschaintypes.NewEmptyRevertOptions()
		revertOptions.RevertAddress = sample.BTCAddressP2WPKH(t, r, &chaincfg.MainNetParams).String()
		revertOptions.AbortAddress = sample.EthAddress().Hex()
		revertOptions.RevertMessage = []byte("some revert message")

		// create test event
		receiver := sample.EthAddress()
		event := createTestBtcEvent(t, &chaincfg.MainNetParams, []byte("dymmy"), &memo.InboundMemo{
			Header: memo.Header{
				OpCode: memo.OpCodeDepositAndCall,
			},
			FieldsV0: memo.FieldsV0{
				Receiver:      receiver,
				Payload:       []byte("some payload"),
				RevertOptions: revertOptions,
			},
		})

		// given receiver and amount
		amountSats := cosmosmath.NewUint(1000)
		event.MsgVoteAmount = amountSats
		event.ToAddress = receiver.Hex()

		// mock SAFE confirmed block
		ob.WithLastBlock(event.BlockNumber + ob.ChainParams().InboundConfirmationSafe())

		// expected vote
		memoBytesExpected := event.MemoStd.Payload
		expectedVote := crosschaintypes.MsgVoteInbound{
			Sender:             event.FromAddress,
			SenderChainId:      chain.ChainId,
			TxOrigin:           event.FromAddress,
			Receiver:           event.MemoStd.Receiver.Hex(),
			ReceiverChain:      ob.ZetacoreClient().Chain().ChainId,
			Amount:             amountSats,
			Message:            hex.EncodeToString(memoBytesExpected),
			InboundHash:        event.TxHash,
			InboundBlockHeight: event.BlockNumber,
			CallOptions: &crosschaintypes.CallOptions{
				GasLimit: 0,
			},
			CoinType:                coin.CoinType_Gas,
			ProtocolContractVersion: crosschaintypes.ProtocolContractVersion_V2,
			RevertOptions: crosschaintypes.RevertOptions{
				RevertAddress: revertOptions.RevertAddress, // should use revert address
				AbortAddress:  revertOptions.AbortAddress,  // should use abort address
				RevertMessage: revertOptions.RevertMessage, // should use revert message
			},
			IsCrossChainCall: true,
			Status:           crosschaintypes.InboundStatus_SUCCESS,
			ConfirmationMode: crosschaintypes.ConfirmationMode_SAFE,
		}

		// create new inbound vote V2 with standard memo
		vote := ob.NewInboundVoteFromEvent(&event)
		require.Equal(t, expectedVote, *vote)
	})
}
