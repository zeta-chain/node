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
	t.Run("should return only EVM and bitcoin chains for mainnet ", func(t *testing.T) {
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
			// Should not include non-EVM, non-bitcoin chains
			require.NotEqual(t, chain.Consensus, chains.Consensus_solana_consensus,
				"chain %s should not have solana consensus", chain.Name)
			require.NotEqual(t, chain.Consensus, chains.Consensus_tendermint,
				"chain %s should not have tendermint consensus", chain.Name)
			require.NotEqual(t, chain.Consensus, chains.Consensus_sui_consensus,
				"chain %s should not have sui consensus", chain.Name)
			require.NotEqual(t, chain.Consensus, chains.Consensus_catchain_consensus,
				"chain %s should not have catchain consensus", chain.Name)
			require.True(t, chain.IsExternal, "chain %s should be external", chain.Name)
			// Should be EVM or bitcoin
			require.True(t,
				chain.Vm == chains.Vm_evm || chain.Consensus == chains.Consensus_bitcoin,
				"chain %s should be EVM or bitcoin", chain.Name)
		}
	})

	t.Run("should return correct mainnet chains requiring migration of TSS funds", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{})

		// Set up all mainnet chains as supported
		mainnetChains := []chains.Chain{
			chains.Ethereum,         // chain_id: 1, Vm_evm
			chains.BscMainnet,       // chain_id: 56, Vm_evm
			chains.BitcoinMainnet,   // chain_id: 8332, Vm_no_vm, bitcoin consensus
			chains.ZetaChainMainnet, // chain_id: 7000, Vm_evm but not external, zevm gateway (excluded)
			chains.Polygon,          // chain_id: 137, Vm_evm
			chains.BaseMainnet,      // chain_id: 8453, Vm_evm
			chains.SolanaMainnet,    // chain_id: 900, Vm_svm No funds migration needed (excluded)
			chains.ArbitrumMainnet,  // chain_id: 42161, Vm_evm
			chains.AvalancheMainnet, // chain_id: 43114, Vm_evm
			chains.SuiMainnet,       // chain_id: 105, Vm_mvm_sui No funds migration needed (excluded)
			chains.TONMainnet,       // chain_id: 2015140, Vm_tvm No funds migration needed (excluded)
		}

		var chainParamsList types.ChainParamsList
		for _, chain := range mainnetChains {
			chainParamsList.ChainParams = append(
				chainParamsList.ChainParams,
				sample.ChainParamsSupported(chain.ChainId),
			)
		}
		zk.ObserverKeeper.SetChainParamsList(ctx, chainParamsList)

		// Get chains supporting TSS migration
		chainsSupportingMigration := k.GetChainsSupportingTSSMigration(ctx)

		expectedChainIDs := map[int64]bool{
			chains.Ethereum.ChainId:         true,
			chains.BscMainnet.ChainId:       true,
			chains.BitcoinMainnet.ChainId:   true,
			chains.Polygon.ChainId:          true,
			chains.BaseMainnet.ChainId:      true,
			chains.ArbitrumMainnet.ChainId:  true,
			chains.AvalancheMainnet.ChainId: true,
		}

		// Verify the count matches expected
		require.Equal(t, len(expectedChainIDs), len(chainsSupportingMigration),
			"expected %d chains, got %d", len(expectedChainIDs), len(chainsSupportingMigration))

		for _, chain := range chainsSupportingMigration {
			require.True(t, expectedChainIDs[chain.ChainId],
				"unexpected chain in result: %s (chain_id: %d, vm: %s)",
				chain.Name, chain.ChainId, chain.Vm)

			require.True(t, chain.IsExternal, "chain %s should be external", chain.Name)
			require.Equal(t, chains.CCTXGateway_observers, chain.CctxGateway,
				"chain %s should have observers gateway", chain.Name)
			require.True(t,
				chain.Vm == chains.Vm_evm || chain.Consensus == chains.Consensus_bitcoin,
				"chain %s should be EVM or bitcoin, got vm=%s consensus=%s", chain.Name, chain.Vm, chain.Consensus)
		}

		// Verify excluded chains are not in the result
		excludedChainIDs := []int64{
			chains.ZetaChainMainnet.ChainId, // not external, zevm gateway
			chains.SolanaMainnet.ChainId,    // Vm_svm
			chains.SuiMainnet.ChainId,       // Vm_mvm_sui
			chains.TONMainnet.ChainId,       // Vm_tvm
		}

		resultChainIDs := make(map[int64]bool)
		for _, chain := range chainsSupportingMigration {
			resultChainIDs[chain.ChainId] = true
		}

		for _, excludedID := range excludedChainIDs {
			require.False(t, resultChainIDs[excludedID],
				"chain with ID %d should be excluded from TSS migration", excludedID)
		}
	})

	// Testnet: ensure we return the exact expected set of chains supporting migration
	t.Run("should return correct testnet chains supporting TSS migration", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{})

		// Set up all relevant testnet chains as supported using the same chain list used in specs
		testnetChains := []chains.Chain{
			chains.BscTestnet,           // chain_id: 97, Vm_evm
			chains.Sepolia,              // chain_id: 11155111, Vm_evm
			chains.BitcoinSignetTestnet, // chain_id: 18333, Vm_no_vm, bitcoin consensus
			chains.BitcoinTestnet4,      // chain_id: 18334, Vm_no_vm, bitcoin consensus
			chains.Amoy,                 // chain_id: 80002, Vm_evm
			chains.BaseSepolia,          // chain_id: 84532, Vm_evm
			chains.ArbitrumSepolia,      // chain_id: 421614, Vm_evm
			chains.AvalancheTestnet,     // chain_id: 43113, Vm_evm
			chains.WorldTestnet,         // chain_id: 4801, Vm_evm
			// Excluded chains also marked supported to assert exclusion
			chains.ZetaChainTestnet, // not external, zevm gateway (excluded)
			chains.SolanaDevnet,     // Vm_svm (excluded)
			chains.SuiTestnet,       // Vm_mvm_sui (excluded)
			chains.TONTestnet,       // Vm_tvm (excluded)
		}

		var chainParamsList types.ChainParamsList
		for _, chain := range testnetChains {
			chainParamsList.ChainParams = append(
				chainParamsList.ChainParams,
				sample.ChainParamsSupported(chain.ChainId),
			)
		}
		zk.ObserverKeeper.SetChainParamsList(ctx, chainParamsList)

		// Get chains supporting TSS migration
		chainsSupportingMigration := k.GetChainsSupportingTSSMigration(ctx)

		// Expected chains are external EVM or Bitcoin testnet chains (observers gateway)
		expectedChainIDs := map[int64]bool{
			chains.BscTestnet.ChainId:           true,
			chains.Sepolia.ChainId:              true,
			chains.BitcoinSignetTestnet.ChainId: true,
			chains.BitcoinTestnet4.ChainId:      true,
			chains.Amoy.ChainId:                 true,
			chains.BaseSepolia.ChainId:          true,
			chains.ArbitrumSepolia.ChainId:      true,
			chains.AvalancheTestnet.ChainId:     true,
			chains.WorldTestnet.ChainId:         true,
		}

		// Verify the count matches expected
		require.Equal(t, len(expectedChainIDs), len(chainsSupportingMigration),
			"expected %d chains, got %d", len(expectedChainIDs), len(chainsSupportingMigration))

		for _, chain := range chainsSupportingMigration {
			require.True(t, expectedChainIDs[chain.ChainId],
				"unexpected chain in result: %s (chain_id: %d, vm: %s)",
				chain.Name, chain.ChainId, chain.Vm)

			require.True(t, chain.IsExternal, "chain %s should be external", chain.Name)
			require.Equal(t, chains.CCTXGateway_observers, chain.CctxGateway,
				"chain %s should have observers gateway", chain.Name)
			require.True(t,
				chain.Vm == chains.Vm_evm || chain.Consensus == chains.Consensus_bitcoin,
				"chain %s should be EVM or bitcoin, got vm=%s consensus=%s", chain.Name, chain.Vm, chain.Consensus)
		}

		// Verify excluded chains are not in the result
		excludedChainIDs := []int64{
			chains.ZetaChainTestnet.ChainId, // not external, zevm gateway
			chains.SolanaDevnet.ChainId,     // Vm_svm
			chains.SuiTestnet.ChainId,       // Vm_mvm_sui
			chains.TONTestnet.ChainId,       // Vm_tvm
		}

		resultChainIDs := make(map[int64]bool)
		for _, chain := range chainsSupportingMigration {
			resultChainIDs[chain.ChainId] = true
		}

		for _, excludedID := range excludedChainIDs {
			require.False(t, resultChainIDs[excludedID],
				"chain with ID %d should be excluded from TSS migration", excludedID)
		}
	})
}
