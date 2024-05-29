package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/authority/keeper"
	"github.com/zeta-chain/zetacore/x/authority/types"
)

func TestMsgServer_UpdateAuthorizations(t *testing.T) {
	t.Run("can't update authorizations with invalid signer", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)

		msg := types.MsgUpdateAuthorizations{
			Signer:            sample.AccAddress(),
			AuthorizationList: sample.AuthorizationList("sample"),
		}

		_, err := msgServer.UpdateAuthorizations(sdk.WrapSDKContext(ctx), &msg)
		require.ErrorIs(t, err, govtypes.ErrInvalidSigner)
	})

	//t.Run("can update authorizations", func(t *testing.T) {
	//	k, ctx := keepertest.AuthorityKeeper(t)
	//	msgServer := keeper.NewMsgServerImpl(*k)
	//	require.NoError(t, k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList()))
	//
	//	authorizationList := sample.AuthorizationList("sample")
	//	msg := types.MsgUpdateAuthorizations{
	//		Signer:            keepertest.AuthorityGovAddress.String(),
	//		AuthorizationList: authorizationList,
	//	}
	//
	//	_, err := msgServer.UpdateAuthorizations(sdk.WrapSDKContext(ctx), &msg)
	//	require.NoError(t, err)
	//
	//	// Check authorization list is set
	//	got, found := k.GetAuthorizationList(ctx)
	//	require.True(t, found)
	//	require.Equal(t, append(types.DefaultAuthorizationsList().Authorizations, authorizationList.Authorizations...), got.Authorizations)
	//})

	t.Run("can add new authorizations when authorizations are not set", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)

		authorizationList := sample.AuthorizationList("sample")
		msg := types.MsgUpdateAuthorizations{
			Signer:            keepertest.AuthorityGovAddress.String(),
			AuthorizationList: authorizationList,
		}

		_, err := msgServer.UpdateAuthorizations(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)

		// Check authorization list is set
		got, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		require.Equal(t, authorizationList, got)
	})
}
