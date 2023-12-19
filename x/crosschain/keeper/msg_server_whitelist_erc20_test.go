package keeper_test

import (
	"fmt"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	crosschainkeeper "github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestKeeper_WhitelistERC20(t *testing.T) {
	t.Run("can deploy and whitelist an erc20", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		msgServer := crosschainkeeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := getValidEthChainID(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)

		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foobar", "FOOBAR")
		k.GetObserverKeeper().SetTssAndUpdateNonce(ctx, sample.Tss())
		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     chainID,
			MedianIndex: 0,
			Prices:      []uint64{1},
		})

		erc20Address := sample.EthAddress().Hex()
		res, err := msgServer.WhitelistERC20(ctx, &types.MsgWhitelistERC20{
			Creator:      admin,
			Erc20Address: erc20Address,
			ChainId:      chainID,
			Name:         "foo",
			Symbol:       "FOO",
			Decimals:     18,
			GasLimit:     100000,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		zrc20 := res.Zrc20Address
		cctxIndex := res.CctxIndex

		// check zrc20 and cctx created
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, ethcommon.HexToAddress(zrc20))
		fc, found := zk.FungibleKeeper.GetForeignCoins(ctx, zrc20)
		require.True(t, found)
		require.EqualValues(t, "foo", fc.Name)
		require.EqualValues(t, erc20Address, fc.Asset)
		cctx, found := k.GetCrossChainTx(ctx, cctxIndex)
		require.True(t, found)
		require.EqualValues(t, fmt.Sprintf("%s:%s", common.CmdWhitelistERC20, erc20Address), cctx.RelayedMessage)

		// check gas limit is set
		gasLimit, err := zk.FungibleKeeper.QueryGasLimit(ctx, ethcommon.HexToAddress(zrc20))
		require.NoError(t, err)
		require.Equal(t, uint64(100000), gasLimit.Uint64())

		// Ensure that whitelist a new erc20 create a cctx with a different index
		res, err = msgServer.WhitelistERC20(ctx, &types.MsgWhitelistERC20{
			Creator:      admin,
			Erc20Address: sample.EthAddress().Hex(),
			ChainId:      chainID,
			Name:         "bar",
			Symbol:       "BAR",
			Decimals:     18,
			GasLimit:     100000,
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.NotEqual(t, cctxIndex, res.CctxIndex)
	})

	t.Run("should fail if not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		msgServer := crosschainkeeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		_, err := msgServer.WhitelistERC20(ctx, &types.MsgWhitelistERC20{
			Creator:      sample.AccAddress(),
			Erc20Address: sample.EthAddress().Hex(),
			ChainId:      getValidEthChainID(t),
			Name:         "foo",
			Symbol:       "FOO",
			Decimals:     18,
			GasLimit:     100000,
		})
		require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
	})

	t.Run("should fail if invalid erc20 address", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		msgServer := crosschainkeeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)

		_, err := msgServer.WhitelistERC20(ctx, &types.MsgWhitelistERC20{
			Creator:      admin,
			Erc20Address: "invalid",
			ChainId:      getValidEthChainID(t),
			Name:         "foo",
			Symbol:       "FOO",
			Decimals:     18,
			GasLimit:     100000,
		})
		require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
	})

	t.Run("should fail if foreign coin already exists for the asset", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		msgServer := crosschainkeeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)

		chainID := getValidEthChainID(t)
		asset := sample.EthAddress().Hex()
		fc := sample.ForeignCoins(t, sample.EthAddress().Hex())
		fc.Asset = asset
		fc.ForeignChainId = chainID
		zk.FungibleKeeper.SetForeignCoins(ctx, fc)

		_, err := msgServer.WhitelistERC20(ctx, &types.MsgWhitelistERC20{
			Creator:      admin,
			Erc20Address: asset,
			ChainId:      chainID,
			Name:         "foo",
			Symbol:       "FOO",
			Decimals:     18,
			GasLimit:     100000,
		})
		require.ErrorIs(t, err, fungibletypes.ErrForeignCoinAlreadyExist)
	})

	t.Run("should fail if no tss set", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		msgServer := crosschainkeeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := getValidEthChainID(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)

		erc20Address := sample.EthAddress().Hex()
		_, err := msgServer.WhitelistERC20(ctx, &types.MsgWhitelistERC20{
			Creator:      admin,
			Erc20Address: erc20Address,
			ChainId:      chainID,
			Name:         "foo",
			Symbol:       "FOO",
			Decimals:     18,
			GasLimit:     100000,
		})
		require.ErrorIs(t, err, types.ErrCannotFindTSSKeys)
	})

	t.Run("should fail if nox valid chain ID", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		msgServer := crosschainkeeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		k.GetObserverKeeper().SetTssAndUpdateNonce(ctx, sample.Tss())

		erc20Address := sample.EthAddress().Hex()
		_, err := msgServer.WhitelistERC20(ctx, &types.MsgWhitelistERC20{
			Creator:      admin,
			Erc20Address: erc20Address,
			ChainId:      10000,
			Name:         "foo",
			Symbol:       "FOO",
			Decimals:     18,
			GasLimit:     100000,
		})
		require.ErrorIs(t, err, types.ErrInvalidChainID)
	})
}
