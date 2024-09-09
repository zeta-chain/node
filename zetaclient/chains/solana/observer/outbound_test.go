package observer_test

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/chains/solana/observer"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

const (
	// gatewayAddressDevnet is the gateway address on devnet for testing
	GatewayAddressTest = "2kJndCL9NBR36ySiQ4bmArs4YgWQu67LmCDfLzk5Gb7s"

	// withdrawTxTest is an archived withdraw tx result on devnet for testing
	// https://explorer.solana.com/tx/5iBYjBYCphzjHKfmPwddMWpV2RNssmzk9Z8NNmV9Rei71pZKBTEVdkmUeyXfn7eWbV8932uSsPfBxgA7UgERNTvq?cluster=devnet
	withdrawTxTest = "5iBYjBYCphzjHKfmPwddMWpV2RNssmzk9Z8NNmV9Rei71pZKBTEVdkmUeyXfn7eWbV8932uSsPfBxgA7UgERNTvq"

	// withdrawFailedTxTest is an archived failed withdraw tx result on devnet for testing
	// https://explorer.solana.com/tx/5nFUQgNSdqTd4aPS4a1xNcbehj19hDzuQLfBqFRj8g7BJdESVY6hFuTFPWFuV6aWAfzEMfVfCdNu9DfzVp5FsHg5?cluster=devnet
	withdrawFailedTxTest = "5nFUQgNSdqTd4aPS4a1xNcbehj19hDzuQLfBqFRj8g7BJdESVY6hFuTFPWFuV6aWAfzEMfVfCdNu9DfzVp5FsHg5"

	// tssAddressTest is the TSS address for testing
	tssAddressTest = "0x05C7dBdd1954D59c9afaB848dA7d8DD3F35e69Cd"
)

// createTestObserver creates a test observer for testing
func createTestObserver(
	t *testing.T,
	chain chains.Chain,
	solClient interfaces.SolanaRPCClient,
	tss interfaces.TSSSigner,
) *observer.Observer {
	database, err := db.NewFromSqliteInMemory(true)
	require.NoError(t, err)

	// create observer
	chainParams := sample.ChainParams(chain.ChainId)
	chainParams.GatewayAddress = GatewayAddressTest
	ob, err := observer.NewObserver(chain, solClient, *chainParams, nil, tss, 60, database, base.DefaultLogger(), nil)
	require.NoError(t, err)

	return ob
}

func Test_CheckFinalizedTx(t *testing.T) {
	// the test chain and transaction hash
	chain := chains.SolanaDevnet
	txHash := withdrawTxTest
	txHashFailed := withdrawFailedTxTest
	txSig := solana.MustSignatureFromBase58(txHash)
	coinType := coin.CoinType_Gas
	nonce := uint64(0)

	// load archived outbound tx result
	txResult := testutils.LoadSolanaOutboundTxResult(t, TestDataDir, chain.ChainId, txHash)

	// mock GetTransaction result
	solClient := mocks.NewSolanaRPCClient(t)
	solClient.On("GetTransaction", mock.Anything, txSig, mock.Anything).Return(txResult, nil)

	// mock TSS
	tss := mocks.NewMockTSS(chain, tssAddressTest, "")

	// create observer with and TSS
	ob := createTestObserver(t, chain, solClient, tss)
	ctx := context.Background()

	t.Run("should successfully check finalized tx", func(t *testing.T) {
		tx, finalized := ob.CheckFinalizedTx(ctx, txHash, nonce, coinType)
		require.True(t, finalized)
		require.NotNil(t, tx)
	})

	t.Run("should return error on invalid tx hash", func(t *testing.T) {
		tx, finalized := ob.CheckFinalizedTx(ctx, "invalid_hash_1234", nonce, coinType)
		require.False(t, finalized)
		require.Nil(t, tx)
	})

	t.Run("should return error on GetTransaction error", func(t *testing.T) {
		// mock GetTransaction error
		client := mocks.NewSolanaRPCClient(t)
		client.On("GetTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))

		// create observer
		ob := createTestObserver(t, chain, client, tss)

		tx, finalized := ob.CheckFinalizedTx(ctx, txHash, nonce, coinType)
		require.False(t, finalized)
		require.Nil(t, tx)
	})

	t.Run("should return error on if transaction is failed", func(t *testing.T) {
		// load archived outbound tx result which is failed due to nonce mismatch
		failedResult := testutils.LoadSolanaOutboundTxResult(t, TestDataDir, chain.ChainId, txHashFailed)

		// mock GetTransaction result with failed status
		client := mocks.NewSolanaRPCClient(t)
		client.On("GetTransaction", mock.Anything, txSig, mock.Anything).Return(failedResult, nil)

		// create observer
		ob := createTestObserver(t, chain, client, tss)

		tx, finalized := ob.CheckFinalizedTx(ctx, txHash, nonce, coinType)
		require.False(t, finalized)
		require.Nil(t, tx)
	})

	t.Run("should return error on ParseGatewayInstruction error", func(t *testing.T) {
		// use CoinType_Zeta to cause ParseGatewayInstruction error
		tx, finalized := ob.CheckFinalizedTx(ctx, txHash, nonce, coin.CoinType_Zeta)
		require.False(t, finalized)
		require.Nil(t, tx)
	})

	t.Run("should return error on ECDSA signer mismatch", func(t *testing.T) {
		// create observer with other TSS address
		tssOther := mocks.NewMockTSS(chain, sample.EthAddress().String(), "")
		ob := createTestObserver(t, chain, solClient, tssOther)

		tx, finalized := ob.CheckFinalizedTx(ctx, txHash, nonce, coinType)
		require.False(t, finalized)
		require.Nil(t, tx)
	})

	t.Run("should return error on nonce mismatch", func(t *testing.T) {
		// use different nonce
		tx, finalized := ob.CheckFinalizedTx(ctx, txHash, nonce+1, coinType)
		require.False(t, finalized)
		require.Nil(t, tx)
	})
}

func Test_ParseGatewayInstruction(t *testing.T) {
	// the test chain and transaction hash
	chain := chains.SolanaDevnet
	txHash := withdrawTxTest
	txAmount := uint64(890880)

	// gateway address
	gatewayID, err := solana.PublicKeyFromBase58(GatewayAddressTest)
	require.NoError(t, err)

	t.Run("should parse gateway instruction", func(t *testing.T) {
		// load archived outbound tx result
		txResult := testutils.LoadSolanaOutboundTxResult(t, TestDataDir, chain.ChainId, txHash)

		// parse gateway instruction
		inst, err := observer.ParseGatewayInstruction(txResult, gatewayID, coin.CoinType_Gas)
		require.NoError(t, err)

		// check sender, nonce and amount
		sender, err := inst.Signer()
		require.NoError(t, err)
		require.Equal(t, tssAddressTest, sender.String())
		require.EqualValues(t, inst.GatewayNonce(), 0)
		require.EqualValues(t, inst.TokenAmount(), txAmount)
	})

	t.Run("should return error on invalid number of instructions", func(t *testing.T) {
		// load and unmarshal archived transaction
		txResult := testutils.LoadSolanaOutboundTxResult(t, TestDataDir, chain.ChainId, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)

		// remove all instructions
		tx.Message.Instructions = nil

		inst, err := observer.ParseGatewayInstruction(txResult, gatewayID, coin.CoinType_Gas)
		require.ErrorContains(t, err, "want 1 instruction, got 0")
		require.Nil(t, inst)
	})

	t.Run("should return error on invalid program id index", func(t *testing.T) {
		// load and unmarshal archived transaction
		txResult := testutils.LoadSolanaOutboundTxResult(t, TestDataDir, chain.ChainId, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)

		// set invalid program id index (out of range)
		tx.Message.Instructions[0].ProgramIDIndex = 4

		inst, err := observer.ParseGatewayInstruction(txResult, gatewayID, coin.CoinType_Gas)
		require.ErrorContains(t, err, "error getting program ID")
		require.Nil(t, inst)
	})

	t.Run("should return error when invoked program is not gateway", func(t *testing.T) {
		// load and unmarshal archived transaction
		txResult := testutils.LoadSolanaOutboundTxResult(t, TestDataDir, chain.ChainId, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)

		// set invalid program id index (pda account index)
		tx.Message.Instructions[0].ProgramIDIndex = 1

		inst, err := observer.ParseGatewayInstruction(txResult, gatewayID, coin.CoinType_Gas)
		require.ErrorContains(t, err, "not matching gatewayID")
		require.Nil(t, inst)
	})

	t.Run("should return error when instruction parsing fails", func(t *testing.T) {
		// load and unmarshal archived transaction
		txResult := testutils.LoadSolanaOutboundTxResult(t, TestDataDir, chain.ChainId, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)

		// set invalid instruction data to cause parsing error
		tx.Message.Instructions[0].Data = []byte("invalid instruction data")

		inst, err := observer.ParseGatewayInstruction(txResult, gatewayID, coin.CoinType_Gas)
		require.Error(t, err)
		require.Nil(t, inst)
	})

	t.Run("should return error on unsupported coin type", func(t *testing.T) {
		// load and unmarshal archived transaction
		txResult := testutils.LoadSolanaOutboundTxResult(t, TestDataDir, chain.ChainId, txHash)

		inst, err := observer.ParseGatewayInstruction(txResult, gatewayID, coin.CoinType_ERC20)
		require.ErrorContains(t, err, "unsupported outbound coin type")
		require.Nil(t, inst)
	})
}

func Test_ParseInstructionWithdraw(t *testing.T) {
	// the test chain and transaction hash
	chain := chains.SolanaDevnet
	txHash := withdrawTxTest
	txAmount := uint64(890880)

	t.Run("should parse instruction withdraw", func(t *testing.T) {
		// load and unmarshal archived transaction
		txResult := testutils.LoadSolanaOutboundTxResult(t, TestDataDir, chain.ChainId, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)

		instruction := tx.Message.Instructions[0]
		inst, err := contracts.ParseInstructionWithdraw(instruction)
		require.NoError(t, err)

		// check sender, nonce and amount
		sender, err := inst.Signer()
		require.NoError(t, err)
		require.Equal(t, tssAddressTest, sender.String())
		require.EqualValues(t, inst.GatewayNonce(), 0)
		require.EqualValues(t, inst.TokenAmount(), txAmount)
	})

	t.Run("should return error on invalid instruction data", func(t *testing.T) {
		// load and unmarshal archived transaction
		txResult := testutils.LoadSolanaOutboundTxResult(t, TestDataDir, chain.ChainId, txHash)
		txFake, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)

		// set invalid instruction data
		instruction := txFake.Message.Instructions[0]
		instruction.Data = []byte("invalid instruction data")

		inst, err := contracts.ParseInstructionWithdraw(instruction)
		require.ErrorContains(t, err, "error deserializing instruction")
		require.Nil(t, inst)
	})

	t.Run("should return error on discriminator mismatch", func(t *testing.T) {
		// load and unmarshal archived transaction
		txResult := testutils.LoadSolanaOutboundTxResult(t, TestDataDir, chain.ChainId, txHash)
		txFake, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)

		// overwrite discriminator (first 8 bytes)
		instruction := txFake.Message.Instructions[0]
		fakeDiscriminator := "b712469c946da12100980d0000000000"
		fakeDiscriminatorBytes, err := hex.DecodeString(fakeDiscriminator)
		require.NoError(t, err)
		copy(instruction.Data, fakeDiscriminatorBytes)

		inst, err := contracts.ParseInstructionWithdraw(instruction)
		require.ErrorContains(t, err, "not a withdraw instruction")
		require.Nil(t, inst)
	})
}
