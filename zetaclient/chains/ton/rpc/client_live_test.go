package rpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	// "github.com/zeta-chain/node/zetaclient/common"
)

func TestLiveClient(t *testing.T) {
	// todo
	// if !common.LiveTestEnabled() {
	// t.Skip("live test is disabled")
	// }

	const endpoint = "https://testnet.toncenter.com/api/v2/"

	ctx := context.Background()

	client := New(endpoint)

	t.Run("GetMasterchainInfo", func(t *testing.T) {
		info, err := client.GetMasterchainInfo(ctx)

		require.NoError(t, err)

		require.Greater(t, info.Seqno, uint64(31622452))
		require.NotEmpty(t, info.RootHash)
		require.NotEmpty(t, info.FileHash)
	})

}
