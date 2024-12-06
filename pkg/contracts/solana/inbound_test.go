package solana

import (
	"encoding/hex"
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
	txHash := "MS3MPLN7hkbyCZFwKqXcg8fmEvQMD74fN6Ps2LSWXJoRxPW5ehaxBorK9q1JFVbqnAvu9jXm6ertj7kT7HpYw1j"
	chain := chains.SolanaDevnet

	txResult := LoadSolanaInboundTxResult(t, txHash)
	tx, err := txResult.Transaction.GetTransaction()
	require.NoError(t, err)

	// create observer
	chainParams := sample.ChainParams(chain.ChainId)
	chainParams.GatewayAddress = testutils.OldSolanaGatewayAddressDevnet
	require.NoError(t, err)

	// expected result
	sender := "AS48jKNQsDGkEdDvfwu1QpqjtqbCadrAq9nGXjFmdX3Z"
	expectedDeposit := &Deposit{
		Sender: sender,
		Amount: 100000,
		Memo:   []byte("0x7F8ae2ABb69A558CE6bAd546f25F0464D9e09e5B4955a3F38ff86ae92A914445099caa8eA2B9bA32"),
		Slot:   txResult.Slot,
		Asset:  "",
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
		tx.Message.Instructions[0].Accounts = append(tx.Message.Instructions[0].Accounts, 4)

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
	txHash := "aY8yLDze6nHSRi7L5REozKAZY1aAyPJ6TfibiqQL5JGwgSBkYux5z5JfXs5ed8LZqpXUy4VijoU3x15mBd66ZGE"
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
	// solana e2e user evm account
	expectedMemo, err := hex.DecodeString("103fd9224f00ce3013e95629e52dfc31d805d68d")
	require.NoError(t, err)
	expectedDeposit := &Deposit{
		Sender: sender,
		Amount: 500000,
		Memo:   expectedMemo,
		Slot:   txResult.Slot,
		Asset:  "4GddKQ7baJpMyKna7bPPnhh7UQtpzfSGL1FgZ31hj4mp", // SPL address
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
		tx.Message.Instructions[0].Accounts = append(tx.Message.Instructions[0].Accounts, 4)

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
