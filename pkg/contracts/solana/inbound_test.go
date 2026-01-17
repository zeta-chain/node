package solana_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/testutils"
)

func LoadObjectFromJSONFile(t *testing.T, obj interface{}, filename string) {
	file, err := os.Open(filepath.Clean(filename))
	require.NoError(t, err)
	defer file.Close()

	// read the struct from the file
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&obj)
	require.NoError(t, err)
}

func LoadSolanaInboundTxResult(
	t *testing.T,
	txHash string,
) *rpc.GetTransactionResult {
	txResult := &rpc.GetTransactionResult{}
	LoadObjectFromJSONFile(t, txResult, fmt.Sprintf("testdata/%s.json", txHash))
	return txResult
}

func Test_ParseInboundAsDeposit(t *testing.T) {
	// ARRANGE
	txHash := "QSoSLxcJAFAzxWnHVJ4s2d5k2LyjC83YaLwbMUHYcEvVnCfERsowNb6Nj55GiTXNTbNF9fzF5F8JHUEpAGMrV5k"
	instructionIndex := 2 // first 2 are compute budget instructions
	chain := chains.SolanaDevnet

	txResult := LoadSolanaInboundTxResult(t, txHash)
	tx, err := txResult.Transaction.GetTransaction()
	require.NoError(t, err)

	// create observer
	chainParams := sample.ChainParams(chain.ChainId)
	chainParams.GatewayAddress = testutils.OldSolanaGatewayAddressDevnet

	// expected result
	// solana e2e deployer account
	sender := "37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ"
	expectedDeposit := &contracts.Inbound{
		Sender:           sender,
		Receiver:         "0x103FD9224F00ce3013e95629e52DFc31D805D68d",
		Amount:           24000000,
		Memo:             []byte{},
		Slot:             txResult.Slot,
		Asset:            "",
		IsCrossChainCall: false,
		RevertOptions:    nil,
	}

	t.Run("should parse inbound event deposit SOL", func(t *testing.T) {
		// ACT
		deposit, err := contracts.ParseInboundAsDeposit(tx, tx.Message.Instructions[instructionIndex], txResult.Slot)
		require.NoError(t, err)

		// ASSERT
		require.EqualValues(t, expectedDeposit, deposit)
	})

	t.Run("should skip parsing if wrong discriminator", func(t *testing.T) {
		// ARRANGE
		txResult := LoadSolanaInboundTxResult(t, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)
		instruction := tx.Message.Instructions[instructionIndex]

		// try deserializing instruction as a 'deposit'
		var inst contracts.DepositInstructionParams
		err = borsh.Deserialize(&inst, instruction.Data)
		require.NoError(t, err)

		// serialize it back with wrong discriminator
		data, err := borsh.Serialize(contracts.DepositInstructionParams{
			Amount:        inst.Amount,
			Discriminator: contracts.DiscriminatorDepositSPL,
			Receiver:      inst.Receiver,
		})
		require.NoError(t, err)

		instruction.Data = data

		// ACT
		deposit, err := contracts.ParseInboundAsDeposit(tx, instruction, txResult.Slot)

		// ASSERT
		require.NoError(t, err)
		require.Nil(t, deposit)
	})

	t.Run("should fail if wrong accounts count", func(t *testing.T) {
		// ARRANGE
		txResult := LoadSolanaInboundTxResult(t, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)
		instruction := tx.Message.Instructions[instructionIndex]

		// remove account from instruction
		instruction.Accounts = instruction.Accounts[:len(instruction.Accounts)-1]

		// ACT
		deposit, err := contracts.ParseInboundAsDeposit(tx, instruction, txResult.Slot)

		// ASSERT
		require.Error(t, err)
		require.Nil(t, deposit)
	})
}

func Test_ParseInboundAsDepositAndCall(t *testing.T) {
	// ARRANGE
	txHash := "3M6UAi5siEjcM25ZfCdHCxHtpfD4UhFncBLckQjYGwJMgzYXvUGgQGCpT19irMGAp7uwVV1SGqqMj8f4UhbXYHZL"
	instructionIndex := 2 // first 2 are compute budget instructions
	chain := chains.SolanaDevnet

	txResult := LoadSolanaInboundTxResult(t, txHash)
	tx, err := txResult.Transaction.GetTransaction()
	require.NoError(t, err)

	// create observer
	chainParams := sample.ChainParams(chain.ChainId)
	chainParams.GatewayAddress = testutils.OldSolanaGatewayAddressDevnet

	// expected result
	// solana e2e deployer account
	sender := "37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ"
	expectedMsg := []byte("hello lamports")
	expectedDeposit := &contracts.Inbound{
		Sender:           sender,
		Receiver:         "0x9d443603009ab6922763790A07BCB115E5636cb0",
		Amount:           1200000,
		Memo:             expectedMsg,
		Slot:             txResult.Slot,
		Asset:            "",
		IsCrossChainCall: true,
	}

	t.Run("should parse inbound event deposit SOL and call", func(t *testing.T) {
		// ACT
		deposit, err := contracts.ParseInboundAsDeposit(tx, tx.Message.Instructions[instructionIndex], txResult.Slot)
		require.NoError(t, err)

		// ASSERT
		require.EqualValues(t, expectedDeposit, deposit)
	})

	t.Run("should skip parsing if wrong discriminator", func(t *testing.T) {
		// ARRANGE
		txResult := LoadSolanaInboundTxResult(t, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)
		instruction := tx.Message.Instructions[instructionIndex]

		// try deserializing instruction as a 'deposit'
		var inst contracts.DepositAndCallInstructionParams
		err = borsh.Deserialize(&inst, instruction.Data)
		require.NoError(t, err)

		// serialize it back with wrong discriminator
		data, err := borsh.Serialize(contracts.DepositAndCallInstructionParams{
			Amount:        inst.Amount,
			Discriminator: contracts.DiscriminatorDepositSPL,
			Receiver:      inst.Receiver,
			Memo:          inst.Memo,
		})
		require.NoError(t, err)

		instruction.Data = data

		// ACT
		deposit, err := contracts.ParseInboundAsDeposit(tx, instruction, txResult.Slot)

		// ASSERT
		require.NoError(t, err)
		require.Nil(t, deposit)
	})

	t.Run("should fail if wrong accounts count", func(t *testing.T) {
		// ARRANGE
		txResult := LoadSolanaInboundTxResult(t, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)
		instruction := tx.Message.Instructions[instructionIndex]

		// remove account from instruction
		instruction.Accounts = instruction.Accounts[:len(instruction.Accounts)-1]

		// ACT
		deposit, err := contracts.ParseInboundAsDeposit(tx, instruction, txResult.Slot)

		// ASSERT
		require.Error(t, err)
		require.Nil(t, deposit)
	})
}

func Test_ParseInboundAsDepositSPL(t *testing.T) {
	// ARRANGE
	txHash := "2GxBKbLxsLC25n4EojhiHxXM26rS4AjFFR3z1vZmmoFXYYxu2U6HBtmp8tBfmPe2JosKRPHvFaQMUQzuGav2ZQSv"
	instructionIndex := 2 // first 2 are compute budget instructions
	chain := chains.SolanaDevnet

	txResult := LoadSolanaInboundTxResult(t, txHash)
	tx, err := txResult.Transaction.GetTransaction()
	require.NoError(t, err)

	// create observer
	chainParams := sample.ChainParams(chain.ChainId)
	chainParams.GatewayAddress = testutils.OldSolanaGatewayAddressDevnet

	// expected result
	// solana e2e deployer account
	sender := "37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ"
	expectedDeposit := &contracts.Inbound{
		Sender:           sender,
		Receiver:         "0x103FD9224F00ce3013e95629e52DFc31D805D68d",
		Amount:           24000000,
		Memo:             []byte{},
		Slot:             txResult.Slot,
		Asset:            "CRpcWQYbvZMrpgVJCZrsAASDmw5xX553EYomdMcjDhDT", // SPL address
		IsCrossChainCall: false,
	}

	t.Run("should parse inbound event deposit SPL", func(t *testing.T) {
		// ACT
		deposit, err := contracts.ParseInboundAsDepositSPL(tx, tx.Message.Instructions[instructionIndex], txResult.Slot)
		require.NoError(t, err)

		// ASSERT
		require.EqualValues(t, expectedDeposit, deposit)
	})

	t.Run("should skip parsing if wrong discriminator", func(t *testing.T) {
		// ARRANGE
		txResult := LoadSolanaInboundTxResult(t, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)
		instruction := tx.Message.Instructions[instructionIndex]

		// try deserializing instruction as a 'deposit_spl'
		var inst contracts.DepositSPLInstructionParams
		err = borsh.Deserialize(&inst, instruction.Data)
		require.NoError(t, err)

		// serialize it back with wrong discriminator
		data, err := borsh.Serialize(contracts.DepositInstructionParams{
			Amount:        inst.Amount,
			Discriminator: contracts.DiscriminatorDeposit,
			Receiver:      inst.Receiver,
		})
		require.NoError(t, err)

		instruction.Data = data

		// ACT
		deposit, err := contracts.ParseInboundAsDepositSPL(tx, instruction, txResult.Slot)

		// ASSERT
		require.NoError(t, err)
		require.Nil(t, deposit)
	})

	t.Run("should fail if wrong accounts count", func(t *testing.T) {
		// ARRANGE
		txResult := LoadSolanaInboundTxResult(t, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)
		instruction := tx.Message.Instructions[instructionIndex]

		// remove account from instruction
		instruction.Accounts = instruction.Accounts[:len(instruction.Accounts)-1]

		// ACT
		deposit, err := contracts.ParseInboundAsDepositSPL(tx, instruction, txResult.Slot)

		// ASSERT
		require.Error(t, err)
		require.Nil(t, deposit)
	})
}

func Test_ParseInboundAsDepositAndCallSPL(t *testing.T) {
	// ARRANGE
	txHash := "SPMngkuRWvmgbrWrhbYJvN4uvFeQco8RNxdAPniTsRRpJxzDoSRJqzVeBrPqtCBYVnMm1GQd53Vkkc5FHSc8HJ3"
	instructionIndex := 2 // first 2 are compute budget instructions
	chain := chains.SolanaDevnet

	txResult := LoadSolanaInboundTxResult(t, txHash)
	tx, err := txResult.Transaction.GetTransaction()
	require.NoError(t, err)

	// create observer
	chainParams := sample.ChainParams(chain.ChainId)
	chainParams.GatewayAddress = testutils.OldSolanaGatewayAddressDevnet

	// expected result
	// solana e2e deployer account
	sender := "37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ"
	// example contract deployed during e2e test, read from tx result
	expectedMsg := []byte("hello spl tokens")
	expectedDeposit := &contracts.Inbound{
		Sender:           sender,
		Receiver:         "0x636e279CD671e7208CFF22DEe061C0fbCBF69ba3",
		Amount:           12000000,
		Memo:             expectedMsg,
		Slot:             txResult.Slot,
		Asset:            "CRpcWQYbvZMrpgVJCZrsAASDmw5xX553EYomdMcjDhDT", // SPL address,
		IsCrossChainCall: true,
	}

	t.Run("should parse inbound event deposit SPL", func(t *testing.T) {
		// ACT
		deposit, err := contracts.ParseInboundAsDepositSPL(tx, tx.Message.Instructions[instructionIndex], txResult.Slot)
		require.NoError(t, err)

		// ASSERT
		require.EqualValues(t, expectedDeposit, deposit)
	})

	t.Run("should skip parsing if wrong discriminator", func(t *testing.T) {
		// ARRANGE
		txResult := LoadSolanaInboundTxResult(t, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)
		instruction := tx.Message.Instructions[instructionIndex]

		// try deserializing instruction as a 'deposit_spl'
		var inst contracts.DepositSPLAndCallInstructionParams
		err = borsh.Deserialize(&inst, instruction.Data)
		require.NoError(t, err)

		// serialize it back with wrong discriminator
		data, err := borsh.Serialize(contracts.DepositSPLAndCallInstructionParams{
			Amount:        inst.Amount,
			Discriminator: contracts.DiscriminatorDeposit,
			Receiver:      inst.Receiver,
			Memo:          inst.Memo,
		})
		require.NoError(t, err)

		instruction.Data = data

		// ACT
		deposit, err := contracts.ParseInboundAsDepositSPL(tx, instruction, txResult.Slot)

		// ASSERT
		require.NoError(t, err)
		require.Nil(t, deposit)
	})

	t.Run("should fail if wrong accounts count", func(t *testing.T) {
		// ARRANGE
		txResult := LoadSolanaInboundTxResult(t, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)
		instruction := tx.Message.Instructions[instructionIndex]

		// remove account from instruction
		instruction.Accounts = instruction.Accounts[:len(instruction.Accounts)-1]

		// ACT
		deposit, err := contracts.ParseInboundAsDepositSPL(tx, instruction, txResult.Slot)

		// ASSERT
		require.Error(t, err)
		require.Nil(t, deposit)
	})
}

func Test_ParseInboundAsCall(t *testing.T) {
	// ARRANGE
	txHash := "5MiGNYQGLfeNf7SYGbR9Rocs2Bya3pDaL3c78CYK5mBHqfbNdDyv4UwGYZFQHrmuo8GZh8rSWYaapoTiKZjE5jsZ"
	chain := chains.SolanaDevnet
	instructionIndex := 2 // first 2 are compute budget instructions

	txResult := LoadSolanaInboundTxResult(t, txHash)
	tx, err := txResult.Transaction.GetTransaction()
	require.NoError(t, err)

	// create observer
	chainParams := sample.ChainParams(chain.ChainId)
	chainParams.GatewayAddress = testutils.OldSolanaGatewayAddressDevnet

	// expected result
	// solana e2e deployer account
	sender := "37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ"
	expectedMsg := []byte("hello")
	expectedCall := &contracts.Inbound{
		Sender:           sender,
		Receiver:         "0x9e6932bB8e63C21b4c1e2346CEac12cA0e90b0cf",
		Amount:           0,
		Memo:             expectedMsg,
		Slot:             txResult.Slot,
		Asset:            "",
		IsCrossChainCall: true,
	}

	t.Run("should parse inbound event call", func(t *testing.T) {
		// ACT
		call, err := contracts.ParseInboundAsCall(tx, tx.Message.Instructions[instructionIndex], txResult.Slot)
		require.NoError(t, err)

		// ASSERT
		require.EqualValues(t, expectedCall, call)
	})

	t.Run("should skip parsing if wrong discriminator", func(t *testing.T) {
		// ARRANGE
		txResult := LoadSolanaInboundTxResult(t, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)
		instruction := tx.Message.Instructions[instructionIndex]

		// try deserializing instruction as a 'call'
		var inst contracts.CallInstructionParams
		err = borsh.Deserialize(&inst, instruction.Data)
		require.NoError(t, err)

		// serialize it back with wrong discriminator
		data, err := borsh.Serialize(contracts.CallInstructionParams{
			Discriminator: contracts.DiscriminatorDepositSPL,
			Receiver:      inst.Receiver,
		})
		require.NoError(t, err)

		instruction.Data = data

		// ACT
		call, err := contracts.ParseInboundAsCall(tx, instruction, txResult.Slot)

		// ASSERT
		require.NoError(t, err)
		require.Nil(t, call)
	})

	t.Run("should fail if wrong accounts count", func(t *testing.T) {
		// ARRANGE
		txResult := LoadSolanaInboundTxResult(t, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)
		instruction := tx.Message.Instructions[instructionIndex]

		// remove account from instruction
		instruction.Accounts = instruction.Accounts[:len(instruction.Accounts)-1]

		// ACT
		call, err := contracts.ParseInboundAsCall(tx, instruction, txResult.Slot)

		// ASSERT
		require.Error(t, err)
		require.Nil(t, call)
	})
}
