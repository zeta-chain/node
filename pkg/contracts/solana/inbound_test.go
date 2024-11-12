package solana

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gagliardetto/solana-go/rpc"
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
	txHash := "MS3MPLN7hkbyCZFwKqXcg8fmEvQMD74fN6Ps2LSWXJoRxPW5ehaxBorK9q1JFVbqnAvu9jXm6ertj7kT7HpYw1j"
	chain := chains.SolanaDevnet

	txResult := LoadSolanaInboundTxResult(t, txHash)
	tx, err := txResult.Transaction.GetTransaction()
	require.NoError(t, err)

	require.NoError(t, err)

	// create observer
	chainParams := sample.ChainParams(chain.ChainId)
	chainParams.GatewayAddress = testutils.GatewayAddresses[chain.ChainId]
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
		deposit, err := ParseInboundAsDeposit(tx, 0, txResult.Slot)
		require.NoError(t, err)

		// check result
		require.EqualValues(t, expectedDeposit, deposit)
	})
}

func Test_ParseInboundAsDepositSPL(t *testing.T) {
	txHash := "aY8yLDze6nHSRi7L5REozKAZY1aAyPJ6TfibiqQL5JGwgSBkYux5z5JfXs5ed8LZqpXUy4VijoU3x15mBd66ZGE"
	chain := chains.SolanaDevnet

	txResult := LoadSolanaInboundTxResult(t, txHash)
	tx, err := txResult.Transaction.GetTransaction()
	require.NoError(t, err)

	// create observer
	chainParams := sample.ChainParams(chain.ChainId)
	chainParams.GatewayAddress = testutils.GatewayAddresses[chain.ChainId]

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
		deposit, err := ParseInboundAsDepositSPL(tx, 0, txResult.Slot)
		require.NoError(t, err)

		// check result
		require.EqualValues(t, expectedDeposit, deposit)
	})
}
