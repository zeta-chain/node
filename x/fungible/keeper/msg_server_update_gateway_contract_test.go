package keeper_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/zrc20.sol"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/fungible/keeper"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestKeeper_UpdateGatewayContract(t *testing.T) {
	t.Run(
		"can update the gateway contract address stored in the module and update address in ZRC20s",
		func(t *testing.T) {
			// ARRANGE
			k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
				UseAuthorityMock: true,
			})

			msgServer := keeper.NewMsgServerImpl(*k)
			k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
			admin := sample.AccAddress()

			authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
			authorityMock.On("GetAdditionalChainList", ctx).Return([]chains.Chain{})

			// setup gas coins for two chains
			defaultChains := chains.DefaultChainsList()
			require.True(t, len(defaultChains) > 1)
			require.NotNil(t, defaultChains[0])
			require.NotNil(t, defaultChains[1])
			chainID1 := defaultChains[0].ChainId
			chainID2 := defaultChains[1].ChainId
			_, _, _, connectorAddr, systemContractAddr := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
			gas1 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID1, "foo", "foo")
			gas2 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID2, "bar", "bar")
			queryZRC20Gateway := func(contract common.Address) string {
				abi, err := zrc20.ZRC20MetaData.GetAbi()
				require.NoError(t, err)
				res, err := k.CallEVM(
					ctx,
					*abi,
					types.ModuleAddressEVM,
					contract,
					keeper.BigIntZero,
					nil,
					false,
					false,
					"gatewayAddress",
				)
				require.NoError(t, err)
				unpacked, err := abi.Unpack("gatewayAddress", res.Ret)
				require.NoError(t, err)
				address, ok := unpacked[0].(common.Address)
				require.True(t, ok)
				return address.Hex()
			}

			// new gateway address
			newGatewayAddr := sample.EthAddress()
			require.NotEqual(t, newGatewayAddr.Hex(), queryZRC20Gateway(gas1))
			require.NotEqual(t, newGatewayAddr.Hex(), queryZRC20Gateway(gas2))

			msg := types.NewMsgUpdateGatewayContract(admin, newGatewayAddr.Hex())
			keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)

			// ACT
			_, err := msgServer.UpdateGatewayContract(ctx, msg)

			// ASSERT
			require.NoError(t, err)
			sc, found := k.GetSystemContract(ctx)
			require.True(t, found)

			// gateway is updated
			require.EqualValues(t, newGatewayAddr.Hex(), sc.Gateway)

			// system contract and connector remain the same
			require.EqualValues(t, systemContractAddr.Hex(), sc.SystemContract)
			require.EqualValues(t, connectorAddr.Hex(), sc.ConnectorZevm)

			// gateway address in ZRC20s is updated
			require.EqualValues(t, newGatewayAddr.Hex(), queryZRC20Gateway(gas1))
			require.EqualValues(t, newGatewayAddr.Hex(), queryZRC20Gateway(gas2))
		},
	)

	t.Run(
		"can update and overwrite the gateway contract if system contract state variable not found",
		func(t *testing.T) {
			// ARRANGE
			k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
				UseAuthorityMock: true,
			})

			msgServer := keeper.NewMsgServerImpl(*k)
			k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
			admin := sample.AccAddress()

			authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

			newGatewayAddr := sample.EthAddress()

			_, found := k.GetSystemContract(ctx)
			require.False(t, found)

			msg := types.NewMsgUpdateGatewayContract(admin, newGatewayAddr.Hex())
			keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)

			// ACT
			_, err := msgServer.UpdateGatewayContract(ctx, msg)

			// ASSERT
			require.NoError(t, err)
			sc, found := k.GetSystemContract(ctx)
			require.True(t, found)

			// gateway is updated
			require.EqualValues(t, newGatewayAddr.Hex(), sc.Gateway)

			// system contract and connector are not updated
			require.EqualValues(t, "", sc.SystemContract)
			require.EqualValues(t, "", sc.ConnectorZevm)
		},
	)

	t.Run("should prevent update the gateway contract if not admin", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		msg := types.NewMsgUpdateGatewayContract(admin, sample.EthAddress().Hex())
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, authoritytypes.ErrUnauthorized)

		// ACT
		_, err := msgServer.UpdateGatewayContract(ctx, msg)

		// ASSERT
		require.Error(t, err)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("should prevent update the gateway contract if invalid gateway address", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		msg := types.NewMsgUpdateGatewayContract(admin, "invalid")
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)

		// ACT
		_, err := msgServer.UpdateGatewayContract(ctx, msg)

		// ASSERT
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid gateway contract address")
	})
}
