package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgServer_UpdateTssAddress(t *testing.T) {
	t.Run("successfully update tss address", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		tssOld := sample.Tss()
		tssNew := sample.Tss()
		k.GetObserverKeeper().SetTSSHistory(ctx, tssOld)
		k.GetObserverKeeper().SetTSSHistory(ctx, tssNew)
		k.GetObserverKeeper().SetTSS(ctx, tssOld)
		for _, chain := range k.GetObserverKeeper().GetSupportedChains(ctx) {
			index := chain.ChainName.String() + "_migration_tx_index"
			k.GetObserverKeeper().SetFundMigrator(ctx, types.TssFundMigratorInfo{
				ChainId:            chain.ChainId,
				MigrationCctxIndex: sample.GetCctxIndexFromString(index),
			})
			cctx := sample.CrossChainTx(t, index)
			cctx.CctxStatus.Status = crosschaintypes.CctxStatus_OutboundMined
			k.SetCrossChainTx(ctx, *cctx)
		}
		assert.Equal(t, len(k.GetObserverKeeper().GetAllTssFundMigrators(ctx)), len(k.GetObserverKeeper().GetSupportedChains(ctx)))
		_, err := msgServer.UpdateTssAddress(ctx, &crosschaintypes.MsgUpdateTssAddress{
			Creator:   admin,
			TssPubkey: tssNew.TssPubkey,
		})
		assert.NoError(t, err)
		tss, found := k.GetObserverKeeper().GetTSS(ctx)
		assert.True(t, found)
		assert.Equal(t, tssNew, tss)
		migrators := k.GetObserverKeeper().GetAllTssFundMigrators(ctx)
		assert.Equal(t, 0, len(migrators))
	})

	t.Run("new tss has not been added to tss history", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		tssOld := sample.Tss()
		tssNew := sample.Tss()
		k.GetObserverKeeper().SetTSSHistory(ctx, tssOld)
		k.GetObserverKeeper().SetTSS(ctx, tssOld)
		for _, chain := range k.GetObserverKeeper().GetSupportedChains(ctx) {
			index := chain.ChainName.String() + "_migration_tx_index"
			k.GetObserverKeeper().SetFundMigrator(ctx, types.TssFundMigratorInfo{
				ChainId:            chain.ChainId,
				MigrationCctxIndex: sample.GetCctxIndexFromString(index),
			})
			cctx := sample.CrossChainTx(t, index)
			cctx.CctxStatus.Status = crosschaintypes.CctxStatus_OutboundMined
			k.SetCrossChainTx(ctx, *cctx)
		}
		assert.Equal(t, len(k.GetObserverKeeper().GetAllTssFundMigrators(ctx)), len(k.GetObserverKeeper().GetSupportedChains(ctx)))
		_, err := msgServer.UpdateTssAddress(ctx, &crosschaintypes.MsgUpdateTssAddress{
			Creator:   admin,
			TssPubkey: tssNew.TssPubkey,
		})
		assert.ErrorContains(t, err, "tss pubkey has not been generated")
		assert.ErrorIs(t, err, crosschaintypes.ErrUnableToUpdateTss)
		tss, found := k.GetObserverKeeper().GetTSS(ctx)
		assert.True(t, found)
		assert.Equal(t, tssOld, tss)
		assert.Equal(t, len(k.GetObserverKeeper().GetAllTssFundMigrators(ctx)), len(k.GetObserverKeeper().GetSupportedChains(ctx)))
	})

	t.Run("old tss pubkey provided", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		tssOld := sample.Tss()
		k.GetObserverKeeper().SetTSSHistory(ctx, tssOld)
		k.GetObserverKeeper().SetTSS(ctx, tssOld)
		for _, chain := range k.GetObserverKeeper().GetSupportedChains(ctx) {
			index := chain.ChainName.String() + "_migration_tx_index"
			k.GetObserverKeeper().SetFundMigrator(ctx, types.TssFundMigratorInfo{
				ChainId:            chain.ChainId,
				MigrationCctxIndex: sample.GetCctxIndexFromString(index),
			})
			cctx := sample.CrossChainTx(t, index)
			cctx.CctxStatus.Status = crosschaintypes.CctxStatus_OutboundMined
			k.SetCrossChainTx(ctx, *cctx)
		}
		assert.Equal(t, len(k.GetObserverKeeper().GetAllTssFundMigrators(ctx)), len(k.GetObserverKeeper().GetSupportedChains(ctx)))
		_, err := msgServer.UpdateTssAddress(ctx, &crosschaintypes.MsgUpdateTssAddress{
			Creator:   admin,
			TssPubkey: tssOld.TssPubkey,
		})
		assert.ErrorContains(t, err, "no new tss address has been generated")
		assert.ErrorIs(t, err, crosschaintypes.ErrUnableToUpdateTss)
		tss, found := k.GetObserverKeeper().GetTSS(ctx)
		assert.True(t, found)
		assert.Equal(t, tssOld, tss)
		assert.Equal(t, len(k.GetObserverKeeper().GetAllTssFundMigrators(ctx)), len(k.GetObserverKeeper().GetSupportedChains(ctx)))
	})

	t.Run("unable to update tss when not enough migrators are present", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		tssOld := sample.Tss()
		tssNew := sample.Tss()

		k.GetObserverKeeper().SetTSSHistory(ctx, tssOld)
		k.GetObserverKeeper().SetTSSHistory(ctx, tssNew)
		k.GetObserverKeeper().SetTSS(ctx, tssOld)
		setSupportedChain(ctx, zk, getValidEthChainIDWithIndex(t, 0), getValidEthChainIDWithIndex(t, 1))

		// set a single migrator while there are 2 supported chains
		chain := k.GetObserverKeeper().GetSupportedChains(ctx)[0]
		index := chain.ChainName.String() + "_migration_tx_index"
		k.GetObserverKeeper().SetFundMigrator(ctx, types.TssFundMigratorInfo{
			ChainId:            chain.ChainId,
			MigrationCctxIndex: sample.GetCctxIndexFromString(index),
		})
		cctx := sample.CrossChainTx(t, index)
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_OutboundMined
		k.SetCrossChainTx(ctx, *cctx)

		require.Equal(t, len(k.GetObserverKeeper().GetAllTssFundMigrators(ctx)), 1)
		_, err := msgServer.UpdateTssAddress(ctx, &crosschaintypes.MsgUpdateTssAddress{
			Creator:   admin,
			TssPubkey: tssNew.TssPubkey,
		})
		require.ErrorContains(t, err, "cannot update tss address not enough migrations have been created and completed")
		require.ErrorIs(t, err, crosschaintypes.ErrUnableToUpdateTss)
		tss, found := k.GetObserverKeeper().GetTSS(ctx)
		require.True(t, found)
		require.Equal(t, tssOld, tss)
		migrators := k.GetObserverKeeper().GetAllTssFundMigrators(ctx)
		assert.Equal(t, 1, len(migrators))
	})

	t.Run("unable to update tss when pending cctx is present", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		tssOld := sample.Tss()
		tssNew := sample.Tss()

		k.GetObserverKeeper().SetTSSHistory(ctx, tssOld)
		k.GetObserverKeeper().SetTSSHistory(ctx, tssNew)
		k.GetObserverKeeper().SetTSS(ctx, tssOld)
		setSupportedChain(ctx, zk, getValidEthChainIDWithIndex(t, 0), getValidEthChainIDWithIndex(t, 1))

		for _, chain := range k.GetObserverKeeper().GetSupportedChains(ctx) {
			index := chain.ChainName.String() + "_migration_tx_index"
			k.GetObserverKeeper().SetFundMigrator(ctx, types.TssFundMigratorInfo{
				ChainId:            chain.ChainId,
				MigrationCctxIndex: sample.GetCctxIndexFromString(index),
			})
			cctx := sample.CrossChainTx(t, index)
			cctx.CctxStatus.Status = crosschaintypes.CctxStatus_PendingOutbound
			k.SetCrossChainTx(ctx, *cctx)
		}
		assert.Equal(t, len(k.GetObserverKeeper().GetAllTssFundMigrators(ctx)), len(k.GetObserverKeeper().GetSupportedChains(ctx)))
		_, err := msgServer.UpdateTssAddress(ctx, &crosschaintypes.MsgUpdateTssAddress{
			Creator:   admin,
			TssPubkey: tssNew.TssPubkey,
		})
		assert.ErrorContains(t, err, "cannot update tss address while there are pending migrations")
		assert.ErrorIs(t, err, crosschaintypes.ErrUnableToUpdateTss)
		tss, found := k.GetObserverKeeper().GetTSS(ctx)
		assert.True(t, found)
		assert.Equal(t, tssOld, tss)
		migrators := k.GetObserverKeeper().GetAllTssFundMigrators(ctx)
		assert.Equal(t, len(k.GetObserverKeeper().GetSupportedChains(ctx)), len(migrators))
	})

	t.Run("unable to update tss cctx is not present", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		tssOld := sample.Tss()
		tssNew := sample.Tss()

		k.GetObserverKeeper().SetTSSHistory(ctx, tssOld)
		k.GetObserverKeeper().SetTSSHistory(ctx, tssNew)
		k.GetObserverKeeper().SetTSS(ctx, tssOld)
		setSupportedChain(ctx, zk, getValidEthChainIDWithIndex(t, 0), getValidEthChainIDWithIndex(t, 1))

		for _, chain := range k.GetObserverKeeper().GetSupportedChains(ctx) {
			index := chain.ChainName.String() + "_migration_tx_index"
			k.GetObserverKeeper().SetFundMigrator(ctx, types.TssFundMigratorInfo{
				ChainId:            chain.ChainId,
				MigrationCctxIndex: sample.GetCctxIndexFromString(index),
			})
		}
		assert.Equal(t, len(k.GetObserverKeeper().GetAllTssFundMigrators(ctx)), len(k.GetObserverKeeper().GetSupportedChains(ctx)))
		_, err := msgServer.UpdateTssAddress(ctx, &crosschaintypes.MsgUpdateTssAddress{
			Creator:   admin,
			TssPubkey: tssNew.TssPubkey,
		})
		assert.ErrorContains(t, err, "migration cross chain tx not found")
		assert.ErrorIs(t, err, crosschaintypes.ErrUnableToUpdateTss)
		tss, found := k.GetObserverKeeper().GetTSS(ctx)
		assert.True(t, found)
		assert.Equal(t, tssOld, tss)
		migrators := k.GetObserverKeeper().GetAllTssFundMigrators(ctx)
		assert.Equal(t, len(k.GetObserverKeeper().GetSupportedChains(ctx)), len(migrators))
	})
}
