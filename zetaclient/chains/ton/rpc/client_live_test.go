package rpc

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/zeta-chain/node/zetaclient/chains/ton/config"
	"github.com/zeta-chain/node/zetaclient/common"
	// "github.com/zeta-chain/node/zetaclient/common"
)

func TestLiveClient(t *testing.T) {
	// todo
	// if !common.LiveTestEnabled() {
	// t.Skip("live test is disabled")
	// }

	endpoint := os.Getenv(common.EnvTONRPC)
	if endpoint == "" {
		endpoint = "https://testnet.toncenter.com/api/v2/"
	}

	ctx := context.Background()

	client := New(endpoint)

	t.Run("HealthCheck", func(t *testing.T) {
		// Involves getMasterchainInfo and getBlockHeader
		blockTime, err := client.HealthCheck(ctx)

		require.NoError(t, err)
		require.Less(t, time.Since(blockTime), 10*time.Second)

		t.Logf("blockTime: %s, since: %s", blockTime, time.Since(blockTime))
	})

	t.Run("GetAccountState", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			accountID := ton.MustParseAccountID("EQB6TUFJZyaq2yJ89NMTyVkS8f5sx0LBjr3jBv9ZiB2IFjrk")

			acc, err := client.GetAccountState(ctx, accountID)

			require.NoError(t, err)
			require.Equal(t, accountID, acc.ID)
			require.Equal(t, AccountStateActive, acc.State)
			require.NotZero(t, acc.Balance)
			require.NotEmpty(t, acc.LastTxHash)
			require.NotZero(t, acc.LastTxLT)

			t.Logf("account: %+v", acc)
		})

		t.Run("NotExists", func(t *testing.T) {
			accountID := ton.MustParseAccountID("0:7a4d41496726aadb227cf4d313c95912f1fe6cc742c18ebde306ff59881d8000")

			acc, err := client.GetAccountState(ctx, accountID)

			require.NoError(t, err)
			require.Equal(t, accountID, acc.ID)
			require.Equal(t, AccountStateNotExists, acc.State)
			require.Zero(t, acc.Balance)
			require.Empty(t, acc.LastTxHash)
			require.Empty(t, acc.LastTxLT)

			t.Logf("account: %+v", acc)
		})
	})

	t.Run("GetConfigParam", func(t *testing.T) {
		// Get gas config
		cell, err := client.GetConfigParam(ctx, 21)

		require.NoError(t, err)
		require.NotNil(t, cell)

		// Parse it
		var cfg tlb.ConfigParam21
		require.NoError(t, tlb.Unmarshal(cell, &cfg))

		gasPrice, err := config.ParseGasPrice(cfg.GasLimitsPrices)
		require.NoError(t, err)

		t.Logf("gasPrice: %d", gasPrice)
	})

	t.Run("GetTransactions", func(t *testing.T) {
		accountID := ton.MustParseAccountID("EQB6TUFJZyaq2yJ89NMTyVkS8f5sx0LBjr3jBv9ZiB2IFjrk")

		txs, err := client.GetTransactions(ctx, 10, accountID, 0, ton.Bits256{})
		require.NoError(t, err)
		require.NotEmpty(t, txs)

		for _, tx := range txs {
			printTx(t, tx)
		}
	})
}

func printTx(t *testing.T, tx ton.Transaction) {
	b, err := json.MarshalIndent(simplifyTx(tx), "", "  ")
	require.NoError(t, err)

	t.Logf("TX %s", string(b))
}

func simplifyTx(tx ton.Transaction) map[string]any {
	return map[string]any{
		"hash":             tx.Hash().Hex(),
		"lt":               tx.Lt,
		"time":             time.Unix(int64(tx.Transaction.Now), 0).UTC().String(),
		"outMessagesCount": tx.OutMsgCnt,
		"gasUsed":          tx.TotalFees.Grams,
	}
}
