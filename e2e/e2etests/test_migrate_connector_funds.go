package e2etests

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/txserver"
	"github.com/zeta-chain/node/e2e/utils"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// TestMigrateConnectorFunds tests the migration of funds from the old ZetaConnectorEth (V1) to the new ZetaConnectorNative (V2)
func TestMigrateConnectorFunds(r *runner.E2ERunner, _ []string) {
	r.Logger.Print("Migrating connector funds from V1 to V2")

	// Define the common transaction receipt handler
	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
		utils.RequireTxSuccessful(r, receipt, failMessage)
	}

	if err := addMigrationAuthorization(r); err != nil {
		require.NoError(r, err)
	}

	pauseConnectors(r, ensureTxReceipt)
	defer unpauseConnectors(r, ensureTxReceipt)

	chainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)
	balance, err := r.ConnectorEth.GetLockedAmount(&bind.CallOpts{})
	require.NoError(r, err, "ZetaConnectorEth GetLockedAmount failed")

	if balance.Int64() == 0 {
		r.Logger.Print("No funds to migrate from old connector")
		return
	}

	// Perform the migration
	migrationIndex := performFundsMigration(r, chainID.Int64(), balance)

	// Verify migration success
	verifyMigrationSuccess(r, migrationIndex, balance)

	// Update chain parameters to use new connector
	// This would disable the old connector and thus stop the V1 flow from working.
	// V1 : Call the connector directly
	// V2 : Call the gateway
	// updateChainParams(r, chainID.Int64())
}

// addMigrationAuthorization adds the necessary authorization for migration
func addMigrationAuthorization(r *runner.E2ERunner) error {
	msgAddAuthorization := authoritytypes.NewMsgAddAuthorization(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		"/zetachain.zetacore.crosschain.MsgMigrateConnectorFunds",
		authoritytypes.PolicyType_groupAdmin,
	)
	_, err := r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, msgAddAuthorization)
	return err
}

// pauseConnectors pauses both V1 and V2 connectors
func pauseConnectors(r *runner.E2ERunner, ensureTxReceipt func(*ethtypes.Transaction, string)) {
	pauseV2Tx, err := r.ConnectorNative.Pause(r.EVMAuth)
	require.NoError(r, err)
	ensureTxReceipt(pauseV2Tx, "ConnectorNative(V2) pause failed")

	pauseV1Tx, err := r.ConnectorEth.Pause(r.EVMAuth)
	require.NoError(r, err)
	ensureTxReceipt(pauseV1Tx, "ZetaConnectorEth(V1) pause failed")
}

// unpauseConnectors unpauses both V1 and V2 connectors
func unpauseConnectors(r *runner.E2ERunner, ensureTxReceipt func(*ethtypes.Transaction, string)) {
	unpauseV2Tx, err := r.ConnectorNative.Unpause(r.EVMAuth)
	require.NoError(r, err)
	ensureTxReceipt(unpauseV2Tx, "ConnectorNative(V2) unpause failed")

	unpauseV1Tx, err := r.ConnectorEth.Unpause(r.EVMAuth)
	require.NoError(r, err)
	ensureTxReceipt(unpauseV1Tx, "ZetaConnectorEth(V1) unpause failed")
}

// performFundsMigration executes the actual funds migration and returns the CCTX index
func performFundsMigration(r *runner.E2ERunner, chainID int64, balance *big.Int) string {
	msgMigrateConnectorFunds := crosschaintypes.NewMsgMigrateConnectorFunds(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		chainID,
		r.ConnectorNativeAddr.Hex(),
		sdkmath.NewUintFromBigInt(balance),
	)

	res, err := r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, msgMigrateConnectorFunds)
	require.NoError(r, err)

	// Extract migration event
	event, ok := txserver.EventOfType[*crosschaintypes.EventConnectorFundsMigration](res.Events)
	require.True(r, ok, "no EventConnectorFundsMigration in %s", res.TxHash)

	return event.CctxIndex
}

// verifyMigrationSuccess verifies that the migration was completed successfully
func verifyMigrationSuccess(r *runner.E2ERunner, cctxIndex string, expectedBalance *big.Int) {
	// Get CCTX details
	cctxRes, err := r.CctxClient.Cctx(r.Ctx, &crosschaintypes.QueryGetCctxRequest{Index: cctxIndex})
	require.NoError(r, err)

	cctx := cctxRes.CrossChainTx
	r.Logger.CCTX(*cctx, "migration")

	// Wait for the CCTX to be mined
	r.WaitForMinedCCTXFromIndex(cctxIndex)

	// Check if the new connector has the funds
	newConnectorBalance, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, r.ConnectorNativeAddr)
	require.NoError(r, err, "BalanceOf failed for new connector")

	// Verify that the migration was successful
	require.Equal(r, expectedBalance, newConnectorBalance,
		"Migration failed: old connector balance (%s) != new connector balance (%s)",
		expectedBalance.String(), newConnectorBalance.String())

	r.Logger.Print("✅ Migration verification successful: %s ZETA tokens migrated", newConnectorBalance.String())
}

// updateChainParams updates the chain parameters to use the new connector address
func updateChainParams(r *runner.E2ERunner, chainID int64) {
	params, err := r.ObserverClient.GetChainParamsForChain(r.Ctx, &observertypes.QueryGetChainParamsForChainRequest{
		ChainId: chainID,
	})
	require.NoError(r, err)

	newChainParams := params.GetChainParams()
	newChainParams.ConnectorContractAddress = r.ConnectorNativeAddr.Hex()

	msgUpdateChainParams := observertypes.NewMsgUpdateChainParams(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		newChainParams)

	_, err = r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, msgUpdateChainParams)
	require.NoError(r, err)
}
