package keeper_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgServer_DeploySystemContracts(t *testing.T) {
	t.Run("admin can deploy system contracts", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.FungibleKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group2)

		res, err := msgServer.DeploySystemContracts(ctx, types.NewMsgDeploySystemContracts(admin))
		require.NoError(t, err)
		require.NotNil(t, res)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, ethcommon.HexToAddress(res.UniswapV2Factory))
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, ethcommon.HexToAddress(res.Wzeta))
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, ethcommon.HexToAddress(res.UniswapV2Router))
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, ethcommon.HexToAddress(res.ConnectorZEVM))
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, ethcommon.HexToAddress(res.SystemContract))
	})

	t.Run("non-admin cannot deploy system contracts", func(t *testing.T) {
		k, ctx, _, zk := keepertest.FungibleKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		nonadmin := sample.AccAddress()
		setAdminPolicies(ctx, zk, nonadmin, observertypes.Policy_Type_group1)

		_, err := msgServer.DeploySystemContracts(ctx, types.NewMsgDeploySystemContracts(nonadmin))
		require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
	})

	t.Run("should fail if contract deployment fails", func(t *testing.T) {
		k, ctx, _, zk := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group2)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		mockEVMKeeper.On(
			"EstimateGas",
			mock.Anything,
			mock.Anything,
		).Return(nil, errors.New("failed to estimate gas"))

		_, err := msgServer.DeploySystemContracts(ctx, types.NewMsgDeploySystemContracts(admin))
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to deploy")
	})
}
