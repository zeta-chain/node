package zetacore

import (
	"context"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/crypto"
	"github.com/zeta-chain/node/testutil/sample"

	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

func Test_GetForeignCoinsFromAsset(t *testing.T) {
	erc20Asset := sample.EthAddress().Hex()

	tests := []struct {
		name    string
		chainID int64
		asset   string
		errMsg  string
	}{
		{
			name:    "get ERC20 foreign coins from asset",
			chainID: 1,
			asset:   erc20Asset,
		},
		{
			name:    "get Gas foreign coins from zero-address asset",
			chainID: 1,
			asset:   constant.EVMZeroAddress,
		},
		{
			name:    "get Gas foreign coins from empty asset",
			chainID: 1,
			asset:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			// construct foreign coin
			assetAddress := ethcommon.HexToAddress(tt.asset)
			asset := assetAddress.Hex()
			if crypto.IsEmptyAddress(assetAddress) {
				asset = ""
			}
			fCoins := sample.ForeignCoins(t, "0x123")
			fCoins.Asset = asset
			fCoins.ForeignChainId = tt.chainID

			// mock RPC server
			method := "/zetachain.zetacore.fungible.Query/ForeignCoinsFromAsset"
			mockRequest := fungibletypes.QueryGetForeignCoinsFromAssetRequest{
				ChainId: fCoins.ForeignChainId,
				Asset:   fCoins.Asset,
			}
			mockResponse := &fungibletypes.QueryGetForeignCoinsFromAssetResponse{
				ForeignCoins: fCoins,
			}
			setupMockServer(t, fungibletypes.RegisterQueryServer, method, mockRequest, mockResponse)
			client := setupZetacoreClient(
				t,
				withDefaultObserverKeys(),
				withCometBFT(mocks.NewSDKClientWithErr(t, nil, 0)),
			)

			// ACT
			ctx := context.Background()
			resp, err := client.GetForeignCoinsFromAsset(ctx, tt.chainID, tt.asset)

			// ASSERT
			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
				require.Equal(t, fungibletypes.ForeignCoins{}, resp)
				return
			}
			require.NoError(t, err)
			require.Equal(t, fCoins, resp)
		})
	}
}

func Test_GetForeignCoinsFromAsset1(t *testing.T) {
	// Given input and output
	fCoinERC20 := sample.ForeignCoins(t, "0x123")
	input := fungibletypes.QueryGetForeignCoinsFromAssetRequest{
		ChainId: fCoinERC20.ForeignChainId,
		Asset:   fCoinERC20.Asset,
	}
	expectedOutput := &fungibletypes.QueryGetForeignCoinsFromAssetResponse{
		ForeignCoins: fCoinERC20,
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
	resp, err := client.GetForeignCoinsFromAsset(ctx, fCoinERC20.ForeignChainId, fCoinERC20.Asset)

	// ASSERT
	require.NoError(t, err)
	require.Equal(t, expectedOutput.ForeignCoins, resp)
}
