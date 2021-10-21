package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
)

func TestSendMsgServerCreate(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	srv := NewMsgServerImpl(*keeper)
	wctx := sdk.WrapSDKContext(ctx)
	creator := "A"
	for i := 0; i < 5; i++ {
		idx := fmt.Sprintf("%d", i)
		expected := &types.MsgCreateSend{Creator: creator, Index: idx}
		_, err := srv.CreateSend(wctx, expected)
		require.NoError(t, err)
		rst, found := keeper.GetSend(ctx, expected.Index)
		require.True(t, found)
		assert.Equal(t, expected.Creator, rst.Creator)
	}
}

func TestSendMsgServerUpdate(t *testing.T) {
	creator := "A"
	index := "any"

	for _, tc := range []struct {
		desc    string
		request *types.MsgUpdateSend
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgUpdateSend{Creator: creator, Index: index},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgUpdateSend{Creator: "B", Index: index},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "KeyNotFound",
			request: &types.MsgUpdateSend{Creator: creator, Index: "missing"},
			err:     sdkerrors.ErrKeyNotFound,
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			keeper, ctx := setupKeeper(t)
			srv := NewMsgServerImpl(*keeper)
			wctx := sdk.WrapSDKContext(ctx)
			expected := &types.MsgCreateSend{Creator: creator, Index: index}
			_, err := srv.CreateSend(wctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateSend(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, found := keeper.GetSend(ctx, expected.Index)
				require.True(t, found)
				assert.Equal(t, expected.Creator, rst.Creator)
			}
		})
	}
}

func TestSendMsgServerDelete(t *testing.T) {
	creator := "A"
	index := "any"

	for _, tc := range []struct {
		desc    string
		request *types.MsgDeleteSend
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgDeleteSend{Creator: creator, Index: index},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgDeleteSend{Creator: "B", Index: index},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "KeyNotFound",
			request: &types.MsgDeleteSend{Creator: creator, Index: "missing"},
			err:     sdkerrors.ErrKeyNotFound,
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			keeper, ctx := setupKeeper(t)
			srv := NewMsgServerImpl(*keeper)
			wctx := sdk.WrapSDKContext(ctx)

			_, err := srv.CreateSend(wctx, &types.MsgCreateSend{Creator: creator, Index: index})
			require.NoError(t, err)
			_, err = srv.DeleteSend(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := keeper.GetSend(ctx, tc.request.Index)
				require.False(t, found)
			}
		})
	}
}
