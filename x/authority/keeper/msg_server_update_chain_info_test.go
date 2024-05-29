package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/authority/keeper"
	"github.com/zeta-chain/zetacore/x/authority/types"
)

func TestMsgServer_UpdateChainInfo(t *testing.T) {
	t.Run("can't update chain info if not authorized", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)

		_, err := msgServer.UpdateChainInfo(sdk.WrapSDKContext(ctx), &types.MsgUpdateChainInfo{
			Signer: sample.AccAddress(),
		})
		require.ErrorIs(t, err, types.ErrUnauthorized)
	})

	t.Run("can update chain info", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)

		// Set group admin policy
		admin := sample.AccAddress()
		k.SetPolicies(ctx, types.Policies{
			Items: []*types.Policy{
				{
					PolicyType: types.PolicyType_groupAdmin,
					Address:    admin,
				},
			},
		})

		_, err := msgServer.UpdateChainInfo(sdk.WrapSDKContext(ctx), &types.MsgUpdateChainInfo{
			Signer: admin,
		})
		require.NoError(t, err)
	})
}
