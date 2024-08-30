package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMsgServer_UpdateTssAddress(t *testing.T) {
	t.Run("should fail if not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})
		admin := sample.AccAddress()

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msgServer := keeper.NewMsgServerImpl(*k)

		msg := crosschaintypes.MsgUpdateTssAddress{
			Creator:   admin,
			TssPubkey: "",
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.UpdateTssAddress(ctx, &msg)
		require.Error(t, err)
	})

	t.Run("should fail if tss not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		msgServer := keeper.NewMsgServerImpl(*k)

		msg := crosschaintypes.MsgUpdateTssAddress{
			Creator:   admin,
			TssPubkey: "",
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.UpdateTssAddress(ctx, &msg)
		require.Error(t, err)
	})

	t.Run("successfully update tss address", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		tssOld := sample.Tss()
		tssNew := sample.Tss()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		msgServer := keeper.NewMsgServerImpl(*k)

		k.GetObserverKeeper().SetTSSHistory(ctx, tssOld)
		k.GetObserverKeeper().SetTSSHistory(ctx, tssNew)
		k.GetObserverKeeper().SetTSS(ctx, tssOld)
		for _, chain := range k.GetChainsSupportingTSSMigration(ctx) {
			index := chain.Name + "_migration_tx_index"
			k.GetObserverKeeper().SetFundMigrator(ctx, types.TssFundMigratorInfo{
				ChainId:            chain.ChainId,
				MigrationCctxIndex: sample.GetCctxIndexFromString(index),
			})
			cctx := sample.CrossChainTx(t, index)
			cctx.CctxStatus.Status = crosschaintypes.CctxStatus_OutboundMined
			k.SetCrossChainTx(ctx, *cctx)
		}
		require.Equal(
			t,
			len(k.GetObserverKeeper().GetAllTssFundMigrators(ctx)),
			len(k.GetChainsSupportingTSSMigration(ctx)),
		)

		msg := crosschaintypes.MsgUpdateTssAddress{
			Creator:   admin,
			TssPubkey: tssNew.TssPubkey,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.UpdateTssAddress(ctx, &msg)
		require.NoError(t, err)
		tss, found := k.GetObserverKeeper().GetTSS(ctx)
		require.True(t, found)
		require.Equal(t, tssNew, tss)
		migrators := k.GetObserverKeeper().GetAllTssFundMigrators(ctx)
		require.Equal(t, 0, len(migrators))
	})

	t.Run("new tss has not been added to tss history", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		tssOld := sample.Tss()
		tssNew := sample.Tss()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msgServer := keeper.NewMsgServerImpl(*k)

		k.GetObserverKeeper().SetTSSHistory(ctx, tssOld)
		k.GetObserverKeeper().SetTSS(ctx, tssOld)
		for _, chain := range k.GetChainsSupportingTSSMigration(ctx) {
			index := chain.Name + "_migration_tx_index"
			k.GetObserverKeeper().SetFundMigrator(ctx, types.TssFundMigratorInfo{
				ChainId:            chain.ChainId,
				MigrationCctxIndex: sample.GetCctxIndexFromString(index),
			})
			cctx := sample.CrossChainTx(t, index)
			cctx.CctxStatus.Status = crosschaintypes.CctxStatus_OutboundMined
			k.SetCrossChainTx(ctx, *cctx)
		}
		require.Equal(
			t,
			len(k.GetObserverKeeper().GetAllTssFundMigrators(ctx)),
			len(k.GetChainsSupportingTSSMigration(ctx)),
		)

		msg := crosschaintypes.MsgUpdateTssAddress{
			Creator:   admin,
			TssPubkey: tssNew.TssPubkey,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.UpdateTssAddress(ctx, &msg)
		require.ErrorContains(t, err, "tss pubkey has not been generated")
		require.ErrorIs(t, err, crosschaintypes.ErrUnableToUpdateTss)
		tss, found := k.GetObserverKeeper().GetTSS(ctx)
		require.True(t, found)
		require.Equal(t, tssOld, tss)
		require.Equal(
			t,
			len(k.GetObserverKeeper().GetAllTssFundMigrators(ctx)),
			len(k.GetChainsSupportingTSSMigration(ctx)),
		)
	})

	t.Run("old tss pubkey provided", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		tssOld := sample.Tss()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msgServer := keeper.NewMsgServerImpl(*k)

		k.GetObserverKeeper().SetTSSHistory(ctx, tssOld)
		k.GetObserverKeeper().SetTSS(ctx, tssOld)
		for _, chain := range k.GetChainsSupportingTSSMigration(ctx) {
			index := chain.Name + "_migration_tx_index"
			k.GetObserverKeeper().SetFundMigrator(ctx, types.TssFundMigratorInfo{
				ChainId:            chain.ChainId,
				MigrationCctxIndex: sample.GetCctxIndexFromString(index),
			})
			cctx := sample.CrossChainTx(t, index)
			cctx.CctxStatus.Status = crosschaintypes.CctxStatus_OutboundMined
			k.SetCrossChainTx(ctx, *cctx)
		}
		require.Equal(
			t,
			len(k.GetObserverKeeper().GetAllTssFundMigrators(ctx)),
			len(k.GetChainsSupportingTSSMigration(ctx)),
		)

		msg := crosschaintypes.MsgUpdateTssAddress{
			Creator:   admin,
			TssPubkey: tssOld.TssPubkey,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.UpdateTssAddress(ctx, &msg)
		require.ErrorContains(t, err, "no new tss address has been generated")
		require.ErrorIs(t, err, crosschaintypes.ErrUnableToUpdateTss)
		tss, found := k.GetObserverKeeper().GetTSS(ctx)
		require.True(t, found)
		require.Equal(t, tssOld, tss)
		require.Equal(
			t,
			len(k.GetObserverKeeper().GetAllTssFundMigrators(ctx)),
			len(k.GetChainsSupportingTSSMigration(ctx)),
		)
	})

	t.Run("unable to update tss when not enough migrators are present", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		tssOld := sample.Tss()
		tssNew := sample.Tss()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		msgServer := keeper.NewMsgServerImpl(*k)

		k.GetObserverKeeper().SetTSSHistory(ctx, tssOld)
		k.GetObserverKeeper().SetTSSHistory(ctx, tssNew)
		k.GetObserverKeeper().SetTSS(ctx, tssOld)
		setSupportedChain(ctx, zk, getValidEthChainIDWithIndex(t, 0), getValidEthChainIDWithIndex(t, 1))

		// set a single migrator while there are 2 supported chains
		chain := k.GetChainsSupportingTSSMigration(ctx)[0]
		index := chain.Name + "_migration_tx_index"
		k.GetObserverKeeper().SetFundMigrator(ctx, types.TssFundMigratorInfo{
			ChainId:            chain.ChainId,
			MigrationCctxIndex: sample.GetCctxIndexFromString(index),
		})
		cctx := sample.CrossChainTx(t, index)
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_OutboundMined
		k.SetCrossChainTx(ctx, *cctx)
		require.Equal(t, len(k.GetObserverKeeper().GetAllTssFundMigrators(ctx)), 1)

		msg := crosschaintypes.MsgUpdateTssAddress{
			Creator:   admin,
			TssPubkey: tssNew.TssPubkey,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.UpdateTssAddress(ctx, &msg)
		require.ErrorContains(
			t,
			err,
			"cannot update tss address incorrect number of migrations have been created and completed: unable to update TSS address",
		)
		require.ErrorIs(t, err, crosschaintypes.ErrUnableToUpdateTss)
		tss, found := k.GetObserverKeeper().GetTSS(ctx)
		require.True(t, found)
		require.Equal(t, tssOld, tss)
		migrators := k.GetObserverKeeper().GetAllTssFundMigrators(ctx)
		require.Equal(t, 1, len(migrators))
	})

	t.Run("unable to update tss when pending cctx is present", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		tssOld := sample.Tss()
		tssNew := sample.Tss()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msgServer := keeper.NewMsgServerImpl(*k)

		k.GetObserverKeeper().SetTSSHistory(ctx, tssOld)
		k.GetObserverKeeper().SetTSSHistory(ctx, tssNew)
		k.GetObserverKeeper().SetTSS(ctx, tssOld)
		setSupportedChain(ctx, zk, getValidEthChainIDWithIndex(t, 0), getValidEthChainIDWithIndex(t, 1))

		for _, chain := range k.GetChainsSupportingTSSMigration(ctx) {
			index := chain.Name + "_migration_tx_index"
			k.GetObserverKeeper().SetFundMigrator(ctx, types.TssFundMigratorInfo{
				ChainId:            chain.ChainId,
				MigrationCctxIndex: sample.GetCctxIndexFromString(index),
			})
			cctx := sample.CrossChainTx(t, index)
			cctx.CctxStatus.Status = crosschaintypes.CctxStatus_PendingOutbound
			k.SetCrossChainTx(ctx, *cctx)
		}
		require.Equal(
			t,
			len(k.GetObserverKeeper().GetAllTssFundMigrators(ctx)),
			len(k.GetObserverKeeper().GetSupportedChains(ctx)),
		)

		msg := crosschaintypes.MsgUpdateTssAddress{
			Creator:   admin,
			TssPubkey: tssNew.TssPubkey,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.UpdateTssAddress(ctx, &msg)
		require.ErrorContains(t, err, "cannot update tss address while there are pending migrations")
		require.ErrorIs(t, err, crosschaintypes.ErrUnableToUpdateTss)
		tss, found := k.GetObserverKeeper().GetTSS(ctx)
		require.True(t, found)
		require.Equal(t, tssOld, tss)
		migrators := k.GetObserverKeeper().GetAllTssFundMigrators(ctx)
		require.Equal(t, len(k.GetObserverKeeper().GetSupportedChains(ctx)), len(migrators))
	})

	t.Run("unable to update tss cctx is not present", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		tssOld := sample.Tss()
		tssNew := sample.Tss()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		msgServer := keeper.NewMsgServerImpl(*k)

		k.GetObserverKeeper().SetTSSHistory(ctx, tssOld)
		k.GetObserverKeeper().SetTSSHistory(ctx, tssNew)
		k.GetObserverKeeper().SetTSS(ctx, tssOld)
		setSupportedChain(ctx, zk, getValidEthChainIDWithIndex(t, 0), getValidEthChainIDWithIndex(t, 1))

		for _, chain := range k.GetChainsSupportingTSSMigration(ctx) {
			index := chain.Name + "_migration_tx_index"
			k.GetObserverKeeper().SetFundMigrator(ctx, types.TssFundMigratorInfo{
				ChainId:            chain.ChainId,
				MigrationCctxIndex: sample.GetCctxIndexFromString(index),
			})
		}
		require.Equal(
			t,
			len(k.GetObserverKeeper().GetAllTssFundMigrators(ctx)),
			len(k.GetObserverKeeper().GetSupportedChains(ctx)),
		)

		msg := crosschaintypes.MsgUpdateTssAddress{
			Creator:   admin,
			TssPubkey: tssNew.TssPubkey,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.UpdateTssAddress(ctx, &msg)
		require.ErrorContains(t, err, "migration cross chain tx not found")
		require.ErrorIs(t, err, crosschaintypes.ErrUnableToUpdateTss)
		tss, found := k.GetObserverKeeper().GetTSS(ctx)
		require.True(t, found)
		require.Equal(t, tssOld, tss)
		migrators := k.GetObserverKeeper().GetAllTssFundMigrators(ctx)
		require.Equal(t, len(k.GetObserverKeeper().GetSupportedChains(ctx)), len(migrators))
	})
}

func TestKeeper_GetChainsSupportingTSSMigration(t *testing.T) {
	t.Run("should return only ethereum and bitcoin chains", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{})
		chainList := chains.ExternalChainList([]chains.Chain{})
		var chainParamsList types.ChainParamsList
		for _, chain := range chainList {
			chainParamsList.ChainParams = append(
				chainParamsList.ChainParams,
				sample.ChainParamsSupported(chain.ChainId),
			)
		}
		zk.ObserverKeeper.SetChainParamsList(ctx, chainParamsList)

		chainsSupportingMigration := k.GetChainsSupportingTSSMigration(ctx)
		for _, chain := range chainsSupportingMigration {
			require.NotEqual(t, chain.Consensus, chains.Consensus_solana_consensus)
			require.NotEqual(t, chain.Consensus, chains.Consensus_op_stack)
			require.NotEqual(t, chain.Consensus, chains.Consensus_tendermint)
			require.Equal(t, chain.IsExternal, true)
		}
	})
}
