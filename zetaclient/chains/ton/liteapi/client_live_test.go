package liteapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/config"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/zeta-chain/node/zetaclient/common"
)

func TestClient(t *testing.T) {
	if !common.LiveTestEnabled() {
		t.Skip("Live tests are disabled")
	}

	var (
		ctx    = context.Background()
		client = New(mustCreateClient(t))
	)

	t.Run("GetFirstTransaction", func(t *testing.T) {
		t.Run("Account doesn't exist", func(t *testing.T) {
			// ARRANGE
			accountID, err := ton.ParseAccountID("0:55798cb7b87168251a7c39f6806b8c202f6caa0f617a76f4070b3fdacfd056a2")
			require.NoError(t, err)

			// ACT
			tx, scrolled, err := client.GetFirstTransaction(ctx, accountID)

			// ASSERT
			require.ErrorContains(t, err, "account is not active")
			require.Zero(t, scrolled)
			require.Nil(t, tx)
		})

		t.Run("All good", func(t *testing.T) {
			// ARRANGE
			// Given sample account id (a dev wallet)
			// https://tonviewer.com/UQCVlMcZ7EyV9maDsvscoLCd5KQfb7CHukyNJluWpMzlD0vr?section=transactions
			accountID, err := ton.ParseAccountID("UQCVlMcZ7EyV9maDsvscoLCd5KQfb7CHukyNJluWpMzlD0vr")
			require.NoError(t, err)

			// Given expected hash for the first tx
			const expect = "b73df4853ca02a040df46f56635d6b8f49b554d5f556881ab389111bbfce4498"

			// as of 2024-09-18
			const expectedTransactions = 23

			start := time.Now()

			// ACT
			tx, scrolled, err := client.GetFirstTransaction(ctx, accountID)

			finish := time.Since(start)

			// ASSERT
			require.NoError(t, err)

			assert.GreaterOrEqual(t, scrolled, expectedTransactions)
			assert.Equal(t, expect, tx.Hash().Hex())

			t.Logf("Time taken %s; transactions scanned: %d", finish.String(), scrolled)
		})
	})

	t.Run("GetTransactionsUntil", func(t *testing.T) {
		// ARRANGE
		// Given sample account id (dev wallet)
		// https://tonviewer.com/UQCVlMcZ7EyV9maDsvscoLCd5KQfb7CHukyNJluWpMzlD0vr?section=transactions
		accountID, err := ton.ParseAccountID("UQCVlMcZ7EyV9maDsvscoLCd5KQfb7CHukyNJluWpMzlD0vr")
		require.NoError(t, err)

		const getUntilLT = uint64(48645164000001)
		const getUntilHash = `2e107215e634bbc3492bdf4b1466d59432623295072f59ab526d15737caa9531`

		// as of 2024-09-20
		const expectedTX = 3

		var hash ton.Bits256
		require.NoError(t, hash.FromHex(getUntilHash))

		start := time.Now()

		// ACT
		// https://tonviewer.com/UQCVlMcZ7EyV9maDsvscoLCd5KQfb7CHukyNJluWpMzlD0vr?section=transactions
		txs, err := client.GetTransactionsSince(ctx, accountID, getUntilLT, hash)

		finish := time.Since(start)

		// ASSERT
		require.NoError(t, err)

		t.Logf("Time taken %s; transactions fetched: %d", finish.String(), len(txs))
		for _, tx := range txs {
			printTx(t, tx)
		}

		mustContainTX(t, txs, "a6672a0e80193c1f705ef1cf45a5883441b8252523b1d08f7656c80e400c74a8")
		assert.GreaterOrEqual(t, len(txs), expectedTX)
	})

	t.Run("GetBlockHeader", func(t *testing.T) {
		// ARRANGE
		// Given sample account id (dev wallet)
		// https://tonscan.org/address/UQCVlMcZ7EyV9maDsvscoLCd5KQfb7CHukyNJluWpMzlD0vr
		accountID, err := ton.ParseAccountID("UQCVlMcZ7EyV9maDsvscoLCd5KQfb7CHukyNJluWpMzlD0vr")
		require.NoError(t, err)

		const getUntilLT = uint64(48645164000001)
		const getUntilHash = `2e107215e634bbc3492bdf4b1466d59432623295072f59ab526d15737caa9531`

		var hash ton.Bits256
		require.NoError(t, hash.FromHex(getUntilHash))

		txs, err := client.GetTransactions(ctx, 1, accountID, getUntilLT, hash)
		require.NoError(t, err)
		require.Len(t, txs, 1)

		// Given a block
		blockID := txs[0].BlockID

		// ACT
		header, err := client.GetBlockHeader(ctx, blockID, 0)

		// ASSERT
		require.NoError(t, err)
		require.NotZero(t, header.MinRefMcSeqno)
		require.Equal(t, header.MinRefMcSeqno, header.MasterRef.Master.SeqNo)
	})
}

func mustCreateClient(t *testing.T) *liteapi.Client {
	client, err := liteapi.NewClient(
		liteapi.WithConfigurationFile(mustFetchConfig(t)),
		liteapi.WithDetectArchiveNodes(),
	)

	require.NoError(t, err)

	return client
}

func mustFetchConfig(t *testing.T) config.GlobalConfigurationFile {
	// archival light client for mainnet
	const url = "https://api.tontech.io/ton/archive-mainnet.autoconf.json"

	res, err := http.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	defer res.Body.Close()

	conf, err := config.ParseConfig(res.Body)
	require.NoError(t, err)

	return *conf
}

func mustContainTX(t *testing.T, txs []ton.Transaction, hash string) {
	var h ton.Bits256
	require.NoError(t, h.FromHex(hash))

	for _, tx := range txs {
		if tx.Hash() == tlb.Bits256(h) {
			return
		}
	}

	t.Fatalf("transaction %q not found", hash)
}

func printTx(t *testing.T, tx ton.Transaction) {
	b, err := json.MarshalIndent(simplifyTx(tx), "", "  ")
	require.NoError(t, err)

	t.Logf("TX %s", string(b))
}

func simplifyTx(tx ton.Transaction) map[string]any {
	return map[string]any{
		"block":            fmt.Sprintf("shard: %d, seqno: %d", tx.BlockID.Shard, tx.BlockID.Seqno),
		"hash":             tx.Hash().Hex(),
		"logicalTime":      tx.Lt,
		"unixTime":         time.Unix(int64(tx.Transaction.Now), 0).UTC().String(),
		"outMessagesCount": tx.OutMsgCnt,
		// "inMessageInfo":    tx.Msgs.InMsg.Value.Value.Info.IntMsgInfo,
		// "outMessages":      tx.Msgs.OutMsgs,
	}
}
