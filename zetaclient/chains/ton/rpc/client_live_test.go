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
)

func TestLiveClient(t *testing.T) {
	if !common.LiveTestEnabled() {
		t.Skip("live test is disabled")
	}

	endpoint := os.Getenv(common.EnvTONRPC)
	if endpoint == "" {
		endpoint = "https://testnet.toncenter.com/api/v2/"
	}

	gatewayTestnet := ton.MustParseAccountID("EQB6TUFJZyaq2yJ89NMTyVkS8f5sx0LBjr3jBv9ZiB2IFjrk")

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
			accountID := gatewayTestnet

			acc, err := client.GetAccountState(ctx, accountID)

			require.NoError(t, err)
			require.Equal(t, accountID, acc.ID)
			require.Equal(t, tlb.AccountActive, acc.Status)
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
			require.Equal(t, tlb.AccountNone, acc.Status)
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
		accountID := gatewayTestnet

		txs, err := client.GetTransactions(ctx, 10, accountID, 0, ton.Bits256{})
		require.NoError(t, err)
		require.NotEmpty(t, txs)

		for _, tx := range txs {
			printTx(t, tx)
		}
	})

	t.Run("GetTransactionsSince", func(t *testing.T) {
		// ARRANGE
		// Given testnet gateway
		accountID := gatewayTestnet

		// Given its last 3 txs
		txs, err := client.GetTransactions(ctx, 3, accountID, 0, ton.Bits256{})
		require.NoError(t, err)
		require.Len(t, txs, 3)

		for i := 0; i < 2; i++ {
			// ensure that GetTransactions orders TXs by DESC
			require.Greater(t, txs[i].Lt, txs[i+1].Lt)
		}

		t.Logf("GetTransactions")
		for _, tx := range txs {
			printTx(t, tx)
		}

		// ACT
		// Get all txs since last-3
		txs2, err := client.GetTransactionsSince(ctx, accountID, txs[2].Lt, ton.Bits256(txs[2].Hash()))

		// ASSERT
		require.NoError(t, err)
		require.Len(t, txs2, 2)

		t.Logf("GetTransactionsSince")
		for _, tx := range txs2 {
			printTx(t, tx)
		}

		// now the pagination should be ASC
		require.Equal(t, txs[0].Lt, txs2[1].Lt)
		require.Equal(t, txs[1].Lt, txs2[0].Lt)

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
