package solana

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
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
	txHash := "8UeJoxY6RbMg6bffsUtZ9f79vSnd4HCRdk5EQgNbAEDYQWXNraiKDtGDZBLp91oyF5eQyWdv6pEwW1vcitiB4px"
	chain := chains.SolanaDevnet

	txResult := LoadSolanaInboundTxResult(t, txHash)
	tx, err := txResult.Transaction.GetTransaction()
	require.NoError(t, err)

	// create observer
	chainParams := sample.ChainParams(chain.ChainId)
	chainParams.GatewayAddress = testutils.OldSolanaGatewayAddressDevnet
	require.NoError(t, err)

	// expected result
	// solana e2e deployer account
	sender := "37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ"
	// solana e2e user evm account
	require.NoError(t, err)
	expectedDeposit := &Deposit{
		Sender:           sender,
		Receiver:         "0x103FD9224F00ce3013e95629e52DFc31D805D68d",
		Amount:           12000000,
		Memo:             []byte{},
		Slot:             txResult.Slot,
		Asset:            "",
		IsCrossChainCall: false,
	}

	t.Run("should parse inbound event deposit SOL", func(t *testing.T) {
		// ACT
		deposit, err := ParseInboundAsDeposit(tx, 0, txResult.Slot)
		require.NoError(t, err)

		// ASSERT
		require.EqualValues(t, expectedDeposit, deposit)
	})

	t.Run("should skip parsing if wrong discriminator", func(t *testing.T) {
		// ARRANGE
		txResult := LoadSolanaInboundTxResult(t, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)

		instruction := tx.Message.Instructions[0]

		// try deserializing instruction as a 'deposit'
		var inst DepositInstructionParams
		err = borsh.Deserialize(&inst, instruction.Data)
		require.NoError(t, err)

		// serialize it back with wrong discriminator
		data, err := borsh.Serialize(DepositInstructionParams{
			Amount:        inst.Amount,
			Discriminator: DiscriminatorDepositSPL,
			Receiver:      inst.Receiver,
		})
		require.NoError(t, err)

		tx.Message.Instructions[0].Data = data

		// ACT
		deposit, err := ParseInboundAsDeposit(tx, 0, txResult.Slot)

		// ASSERT
		require.NoError(t, err)
		require.Nil(t, deposit)
	})

	t.Run("should fail if wrong accounts count", func(t *testing.T) {
		// ARRANGE
		txResult := LoadSolanaInboundTxResult(t, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)

		// append one more account to instruction
		tx.Message.AccountKeys = append(tx.Message.AccountKeys, solana.MustPublicKeyFromBase58(sample.SolanaAddress(t)))
		tx.Message.Instructions[0].Accounts = tx.Message.Instructions[0].Accounts[:len(tx.Message.Instructions[0].Accounts)-1]

		// ACT
		deposit, err := ParseInboundAsDeposit(tx, 0, txResult.Slot)

		// ASSERT
		require.Error(t, err)
		require.Nil(t, deposit)
	})

	t.Run("should fail if first account is not signer", func(t *testing.T) {
		// ARRANGE
		txResult := LoadSolanaInboundTxResult(t, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)

		// switch account places
		tx.Message.Instructions[0].Accounts[0] = 1
		tx.Message.Instructions[0].Accounts[1] = 0

		// ACT
		deposit, err := ParseInboundAsDeposit(tx, 0, txResult.Slot)

		// ASSERT
		require.Error(t, err)
		require.Nil(t, deposit)
	})
}

func Test_ParseInboundAsDepositAndCall(t *testing.T) {
	// ARRANGE
	txHash := "5b7ShhHf8dvUjUBHgMvgH8FFqpfAd7vAGygZLaeNPhugXtY5fatPSACVkn13o7sw6Awob8EJnrwAuiKYqvi7ZkHa"
	chain := chains.SolanaDevnet

	txResult := LoadSolanaInboundTxResult(t, txHash)
	tx, err := txResult.Transaction.GetTransaction()
	require.NoError(t, err)

	// create observer
	chainParams := sample.ChainParams(chain.ChainId)
	chainParams.GatewayAddress = testutils.OldSolanaGatewayAddressDevnet
	require.NoError(t, err)

	// expected result
	// solana e2e deployer account
	sender := "37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ"
	expectedMsg := []byte("hello lamports")
	expectedDeposit := &Deposit{
		Sender:           sender,
		Receiver:         "0x75A06a8C258739dADfe2352D57973deF9ee7A2ba",
		Amount:           1200000,
		Memo:             expectedMsg,
		Slot:             txResult.Slot,
		Asset:            "",
		IsCrossChainCall: true,
	}

	t.Run("should parse inbound event deposit SOL and call", func(t *testing.T) {
		// ACT
		deposit, err := ParseInboundAsDeposit(tx, 0, txResult.Slot)
		require.NoError(t, err)

		// ASSERT
		require.EqualValues(t, expectedDeposit, deposit)
	})

	t.Run("should skip parsing if wrong discriminator", func(t *testing.T) {
		// ARRANGE
		txResult := LoadSolanaInboundTxResult(t, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)

		instruction := tx.Message.Instructions[0]

		// try deserializing instruction as a 'deposit'
		var inst DepositInstructionParams
		err = borsh.Deserialize(&inst, instruction.Data)
		require.NoError(t, err)

		// serialize it back with wrong discriminator
		data, err := borsh.Serialize(DepositInstructionParams{
			Amount:        inst.Amount,
			Discriminator: DiscriminatorDepositSPL,
			Receiver:      inst.Receiver,
		})
		require.NoError(t, err)

		tx.Message.Instructions[0].Data = data

		// ACT
		deposit, err := ParseInboundAsDeposit(tx, 0, txResult.Slot)

		// ASSERT
		require.NoError(t, err)
		require.Nil(t, deposit)
	})

	t.Run("should fail if wrong accounts count", func(t *testing.T) {
		// ARRANGE
		txResult := LoadSolanaInboundTxResult(t, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)

		// append one more account to instruction
		tx.Message.AccountKeys = append(tx.Message.AccountKeys, solana.MustPublicKeyFromBase58(sample.SolanaAddress(t)))
		tx.Message.Instructions[0].Accounts = tx.Message.Instructions[0].Accounts[:len(tx.Message.Instructions[0].Accounts)-1]

		// ACT
		deposit, err := ParseInboundAsDeposit(tx, 0, txResult.Slot)

		// ASSERT
		require.Error(t, err)
		require.Nil(t, deposit)
	})

	t.Run("should fail if first account is not signer", func(t *testing.T) {
		// ARRANGE
		txResult := LoadSolanaInboundTxResult(t, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)

		// switch account places
		tx.Message.Instructions[0].Accounts[0] = 1
		tx.Message.Instructions[0].Accounts[1] = 0

		// ACT
		deposit, err := ParseInboundAsDeposit(tx, 0, txResult.Slot)

		// ASSERT
		require.Error(t, err)
		require.Nil(t, deposit)
	})
}

func Test_ParseInboundAsDepositSPL(t *testing.T) {
	// ARRANGE
	txHash := "5bXSQaq6BY1WhhF3Qm4pLHXxuyM9Mz1MrdMeoCFbimxw4uv11raQgAj4HGULPEQExPKB231rMhm6666dQMwf9fNN"
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
	require.NoError(t, err)
	expectedDeposit := &Deposit{
		Sender:           sender,
		Receiver:         "0x103FD9224F00ce3013e95629e52DFc31D805D68d",
		Amount:           12000000,
		Memo:             []byte{},
		Slot:             txResult.Slot,
		Asset:            "BTmtL9Dh2DcwhPntEbjo3rSWpmz1EhXsmohSC7CGSEWw", // SPL address
		IsCrossChainCall: false,
	}

	t.Run("should parse inbound event deposit SPL", func(t *testing.T) {
		// ACT
		deposit, err := ParseInboundAsDepositSPL(tx, 0, txResult.Slot)
		require.NoError(t, err)

		// ASSERT
		require.EqualValues(t, expectedDeposit, deposit)
	})

	t.Run("should skip parsing if wrong discriminator", func(t *testing.T) {
		// ARRANGE
		txResult := LoadSolanaInboundTxResult(t, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)

		instruction := tx.Message.Instructions[0]

		// try deserializing instruction as a 'deposit_spl'
		var inst DepositSPLInstructionParams
		err = borsh.Deserialize(&inst, instruction.Data)
		require.NoError(t, err)

		// serialize it back with wrong discriminator
		data, err := borsh.Serialize(DepositInstructionParams{
			Amount:        inst.Amount,
			Discriminator: DiscriminatorDeposit,
			Receiver:      inst.Receiver,
		})
		require.NoError(t, err)

		tx.Message.Instructions[0].Data = data

		// ACT
		deposit, err := ParseInboundAsDepositSPL(tx, 0, txResult.Slot)

		// ASSERT
		require.NoError(t, err)
		require.Nil(t, deposit)
	})

	t.Run("should fail if wrong accounts count", func(t *testing.T) {
		// ARRANGE
		txResult := LoadSolanaInboundTxResult(t, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)

		// append one more account to instruction
		tx.Message.AccountKeys = append(tx.Message.AccountKeys, solana.MustPublicKeyFromBase58(sample.SolanaAddress(t)))
		tx.Message.Instructions[0].Accounts = tx.Message.Instructions[0].Accounts[:len(tx.Message.Instructions[0].Accounts)-1]

		// ACT
		deposit, err := ParseInboundAsDepositSPL(tx, 0, txResult.Slot)

		// ASSERT
		require.Error(t, err)
		require.Nil(t, deposit)
	})

	t.Run("should fail if first account is not signer", func(t *testing.T) {
		// ARRANGE
		txResult := LoadSolanaInboundTxResult(t, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)

		// switch account places
		tx.Message.Instructions[0].Accounts[0] = 1
		tx.Message.Instructions[0].Accounts[1] = 0

		// ACT
		deposit, err := ParseInboundAsDepositSPL(tx, 0, txResult.Slot)

		// ASSERT
		require.Error(t, err)
		require.Nil(t, deposit)
	})
}

func Test_ParseInboundAsDepositAndCallSPL(t *testing.T) {
	// ARRANGE
	txHash := "22s5ERRRZmZXAuDMRdwUU33VnWZ7m8NHUZM6hyLH52JQPz5R7mXEkFcvHx88ujq3xDnt3z7sZdZ21JK2FC7vPw1o"
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
	expectedDeposit := &Deposit{
		Sender:           sender,
		Receiver:         "0xd5Fef042019aFAEe2783092d0502bEc0141f67D1",
		Amount:           12000000,
		Memo:             expectedMsg,
		Slot:             txResult.Slot,
		Asset:            "7d4ehzE4WNgithQZMyQFDhmHyN6rQNTEC7re1bsRN7TX", // SPL address,
		IsCrossChainCall: true,
	}

	t.Run("should parse inbound event deposit SPL", func(t *testing.T) {
		// ACT
		deposit, err := ParseInboundAsDepositSPL(tx, 0, txResult.Slot)
		require.NoError(t, err)

		// ASSERT
		require.EqualValues(t, expectedDeposit, deposit)
	})

	t.Run("should skip parsing if wrong discriminator", func(t *testing.T) {
		// ARRANGE
		txResult := LoadSolanaInboundTxResult(t, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)

		instruction := tx.Message.Instructions[0]

		// try deserializing instruction as a 'deposit_spl'
		var inst DepositSPLInstructionParams
		err = borsh.Deserialize(&inst, instruction.Data)
		require.NoError(t, err)

		// serialize it back with wrong discriminator
		data, err := borsh.Serialize(DepositInstructionParams{
			Amount:        inst.Amount,
			Discriminator: DiscriminatorDeposit,
			Receiver:      inst.Receiver,
		})
		require.NoError(t, err)

		tx.Message.Instructions[0].Data = data

		// ACT
		deposit, err := ParseInboundAsDepositSPL(tx, 0, txResult.Slot)

		// ASSERT
		require.NoError(t, err)
		require.Nil(t, deposit)
	})

	t.Run("should fail if wrong accounts count", func(t *testing.T) {
		// ARRANGE
		txResult := LoadSolanaInboundTxResult(t, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)

		// append one more account to instruction
		tx.Message.AccountKeys = append(tx.Message.AccountKeys, solana.MustPublicKeyFromBase58(sample.SolanaAddress(t)))
		tx.Message.Instructions[0].Accounts = tx.Message.Instructions[0].Accounts[:len(tx.Message.Instructions[0].Accounts)-1]

		// ACT
		deposit, err := ParseInboundAsDepositSPL(tx, 0, txResult.Slot)

		// ASSERT
		require.Error(t, err)
		require.Nil(t, deposit)
	})

	t.Run("should fail if first account is not signer", func(t *testing.T) {
		// ARRANGE
		txResult := LoadSolanaInboundTxResult(t, txHash)
		tx, err := txResult.Transaction.GetTransaction()
		require.NoError(t, err)

		// switch account places
		tx.Message.Instructions[0].Accounts[0] = 1
		tx.Message.Instructions[0].Accounts[1] = 0

		// ACT
		deposit, err := ParseInboundAsDepositSPL(tx, 0, txResult.Slot)

		// ASSERT
		require.Error(t, err)
		require.Nil(t, deposit)
	})
}
