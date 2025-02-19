package zetacore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/testutil/sample"

	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

func Test_GetForeignCoinsFromAsset(t *testing.T) {
	// Given input and output
	fCoin := sample.ForeignCoins(t, "0x123")
	input := fungibletypes.QueryGetForeignCoinsFromAssetRequest{
		ChainId: fCoin.ForeignChainId,
		Asset:   fCoin.Asset,
	}
	expectedOutput := &fungibletypes.QueryGetForeignCoinsFromAssetResponse{
		ForeignCoins: fCoin,
	}

	// ARRANGE
	ctx := context.Background()
	method := "/zetachain.zetacore.fungible.Query/ForeignCoinsFromAsset"
	setupMockServer(t, fungibletypes.RegisterQueryServer, method, input, expectedOutput)
	client := setupZetacoreClient(
		t,
		withDefaultObserverKeys(),
		withCometBFT(mocks.NewSDKClientWithErr(t, nil, 0)),
	)

	// ACT
	resp, err := client.GetForeignCoinsFromAsset(ctx, fCoin.ForeignChainId, fCoin.Asset)

	// ASSERT
	require.NoError(t, err)
	require.Equal(t, expectedOutput.ForeignCoins, resp)
}
