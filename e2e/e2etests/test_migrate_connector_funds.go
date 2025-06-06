package e2etests

import (
	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/txserver"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestMigrateConnectorFunds tests the migration of funds from the old ZetaConnectorEth (V1) to the new ZetaConnectorNative (V2)
func TestMigrateConnectorFunds(r *runner.E2ERunner, _ []string) {
	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
		utils.RequireTxSuccessful(r, receipt, failMessage)
	}

	pauseV2Tx, err := r.ConnectorNative.Pause(r.EVMAuth)
	require.NoError(r, err)
	ensureTxReceipt(pauseV2Tx, "ConnectorNative(V2) pause failed")

	pauseV1Tx, err := r.ConnectorEth.Pause(r.EVMAuth)
	require.NoError(r, err)
	ensureTxReceipt(pauseV1Tx, "ZetaConnectorEth(V1) pause failed")

	defer func() {
		unpauseV2Tx, err := r.ConnectorNative.Unpause(r.EVMAuth)
		require.NoError(r, err)
		ensureTxReceipt(unpauseV2Tx, "ConnectorNative(V2) pause failed")

		unpauseV1Tx, err := r.ConnectorEth.Unpause(r.EVMAuth)
		require.NoError(r, err)
		ensureTxReceipt(unpauseV1Tx, "ZetaConnectorEth(V1) pause failed")
	}()

	chainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	balance, err := r.ConnectorEth.GetLockedAmount(&bind.CallOpts{})
	require.NoError(r, err, "ZetaConnectorEth GetLockedAmount failed")

	if balance.Int64() == 0 {
		r.Logger.Print("No funds to migrate from old connector")
		return
	}
	// Migrate funds using Admin CMD
	msgMigrateConnectorFunds := crosschaintypes.NewMsgMigrateConnectorFunds(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		chainID.Int64(),
		r.ConnectorNativeAddr.Hex(),
		sdkmath.NewUintFromBigInt(balance),
	)
	res, err := r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, msgMigrateConnectorFunds)
	require.NoError(r, err)

	// Verify that the funds have been migrated successfully
	event, ok := txserver.EventOfType[*crosschaintypes.EventConnectorFundsMigration](res.Events)
	require.True(r, ok, "no EventERC20CustodyFundsMigration in %s", res.TxHash)

	cctxRes, err := r.CctxClient.Cctx(r.Ctx, &crosschaintypes.QueryGetCctxRequest{Index: event.CctxIndex})
	require.NoError(r, err)

	cctx := cctxRes.CrossChainTx
	r.Logger.CCTX(*cctx, "migration")

	// wait for the cctx to be mined
	r.WaitForMinedCCTXFromIndex(event.CctxIndex)

	// check if the new connector has the funds
	newConnectorBalance, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, r.ConnectorNativeAddr)
	require.NoError(r, err, "BalanceOf failed for new connector")

	// Verify that the migration was successful
	require.Equal(r, balance, newConnectorBalance,
		"Migration failed: old connector balance (%s) != new connector balance (%s)",
		balance.String(), newConnectorBalance.String())

	r.Logger.Print("✅ Migration verification successful: %s ZETA tokens migrated", newConnectorBalance.String())
}
