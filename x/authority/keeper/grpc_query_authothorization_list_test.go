package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/x/authority/keeper"
	"github.com/zeta-chain/node/x/authority/types"
)

func TestKeeper_Authorization(t *testing.T) {
	k, ctx := keepertest.AuthorityKeeper(t)
	authorizationList := types.AuthorizationList{Authorizations: []types.Authorization{
		{
			MsgUrl:           "ABC",
			AuthorizedPolicy: types.PolicyType_groupOperational,
		},
		{
			MsgUrl:           "DEF",
			AuthorizedPolicy: types.PolicyType_groupAdmin,
		},
	}}

	tt := []struct {
		name                 string
		setAuthorizationList bool
		req                  *types.QueryAuthorizationRequest
		expectedResponse     *types.QueryAuthorizationResponse
		expecterErrorString  string
	}{
		{
			name:                 "successfully get authorization",
			setAuthorizationList: true,
			req: &types.QueryAuthorizationRequest{
				MsgUrl: "ABC",
			},
			expectedResponse: &types.QueryAuthorizationResponse{
				Authorization: types.Authorization{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
			},
			expecterErrorString: "",
		},
		{
			name:                 "invalid request",
			setAuthorizationList: true,
			req:                  nil,
			expectedResponse:     nil,
			expecterErrorString:  "invalid request",
		},
		{
			name:                 "invalid msg url",
			setAuthorizationList: true,
			req: &types.QueryAuthorizationRequest{
				MsgUrl: "",
			},
			expectedResponse:    nil,
			expecterErrorString: "message URL cannot be empty",
		},
		{
			name:                 "authorization not found",
			setAuthorizationList: true,
			req: &types.QueryAuthorizationRequest{
				MsgUrl: "GHI",
			},
			expectedResponse:    nil,
			expecterErrorString: "authorization not found",
		},
		{
			name:                 "authorization list not found",
			setAuthorizationList: false,
			req: &types.QueryAuthorizationRequest{
				MsgUrl: "ABC",
			},
			expectedResponse:    nil,
			expecterErrorString: "authorization list not found",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			removeAuthorizationList(ctx, *k)
			if tc.setAuthorizationList {
				k.SetAuthorizationList(ctx, authorizationList)
			}
			response, err := k.Authorization(ctx, tc.req)
			require.Equal(t, tc.expectedResponse, response)
			if tc.expecterErrorString != "" {
				require.ErrorContains(t, err, tc.expecterErrorString)
			}
		})
	}
}

func TestKeeper_AuthorizationList(t *testing.T) {
	k, ctx := keepertest.AuthorityKeeper(t)
	authorizationList := types.AuthorizationList{Authorizations: []types.Authorization{
		{
			MsgUrl:           "ABC",
			AuthorizedPolicy: types.PolicyType_groupOperational,
		},
		{
			MsgUrl:           "DEF",
			AuthorizedPolicy: types.PolicyType_groupAdmin,
		},
	}}
	tt := []struct {
		name                 string
		setAuthorizationList bool
		req                  *types.QueryAuthorizationListRequest
		expectedResponse     *types.QueryAuthorizationListResponse
		expecterErrorString  string
	}{
		{
			name:                 "successfully get authorization list",
			setAuthorizationList: true,
			req:                  &types.QueryAuthorizationListRequest{},
			expectedResponse: &types.QueryAuthorizationListResponse{
				AuthorizationList: authorizationList,
			},
			expecterErrorString: "",
		},
		{
			name:                 "invalid request",
			setAuthorizationList: true,
			req:                  nil,
			expectedResponse:     nil,
			expecterErrorString:  "invalid request",
		},
		{
			name:                 "authorization list not found",
			setAuthorizationList: false,
			req:                  &types.QueryAuthorizationListRequest{},
			expectedResponse:     nil,
			expecterErrorString:  "authorization list not found",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			removeAuthorizationList(ctx, *k)
			if tc.setAuthorizationList {
				k.SetAuthorizationList(ctx, authorizationList)
			}
			response, err := k.AuthorizationList(ctx, tc.req)
			require.Equal(t, tc.expectedResponse, response)
			if tc.expecterErrorString != "" {
				require.ErrorContains(t, err, tc.expecterErrorString)
			}
		})

	}
}

func removeAuthorizationList(ctx sdk.Context, k keeper.Keeper) {
	store := prefix.NewStore(ctx.KVStore(k.GetStoreKey()), types.KeyPrefix(types.AuthorizationListKey))
	store.Delete([]byte{0})
}
