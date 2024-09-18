package liteapi

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/config"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/ton"
	"github.com/zeta-chain/node/zetaclient/common"
)

func TestClient(t *testing.T) {
	if !common.LiveTestEnabled() {
		t.Skip("Live tests are disabled")
	}

	var (
		ctx    = context.Background()
		client = &Client{Client: mustCreateClient(t)}
	)

	t.Run("GetFirstTransaction", func(t *testing.T) {
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
