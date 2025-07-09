package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

// TestMigrateConnectorFunds tests the migration of funds from the old ZetaConnectorEth (V1) to the new ZetaConnectorNative (V2)
func TestMigrateConnectorFunds(r *runner.E2ERunner, _ []string) {
	r.Logger.Print("Migrating connector funds from V1 to V2")
	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
		utils.RequireTxSuccessful(r, receipt, failMessage)
	}

	// Pause the connectors before migration
	pauseConnectors(r, ensureTxReceipt)
	updateTSSAddress(r, r.EVMAddress(), ensureTxReceipt)

	// Transfer all funds from the old connector to the new connector
	balanceTransferred := transferAllFunds(r, ensureTxReceipt)
	verifyMigrationSuccess(r, balanceTransferred)

	// Unpause the connectors after migration and reset the TSS address
	unpauseConnectors(r, ensureTxReceipt)
	updateTSSAddress(r, r.TSSAddress, ensureTxReceipt)
}

func updateTSSAddress(
	r *runner.E2ERunner,
	tssAddress common.Address,
	ensureTxReceipt func(*ethtypes.Transaction, string),
) {
	updateTssTx, err := r.ConnectorEth.UpdateTssAddress(r.EVMAuth, tssAddress)
	require.NoError(r, err)
	ensureTxReceipt(updateTssTx, "ZetaConnectorEth TSS address update failed")
}

func transferAllFunds(r *runner.E2ERunner, ensureTxReceipt func(*ethtypes.Transaction, string)) *big.Int {
	chainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)
	totalAmount, err := r.ConnectorEth.GetLockedAmount(&bind.CallOpts{})
	require.NoError(r, err)

	// Transfer all funds
	// message should be empty so that there is no call triggered on the new connector
	transferTx, err := r.ConnectorEth.OnReceive(
		r.EVMAuth,
		[]byte{},              // empty zetaTxSenderAddress
		chainID,               // sourceChainId
		r.ConnectorNativeAddr, // destinationAddress
		totalAmount,           // zetaValue
		[]byte{},              // empty message
		[32]byte{},            // empty internalSendHash
	)
	require.NoError(r, err)
	ensureTxReceipt(transferTx, "Fund transfer failed")
	return totalAmount
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

// verifyMigrationSuccess verifies that the migration was completed successfully
func verifyMigrationSuccess(r *runner.E2ERunner, expectedBalance *big.Int) {
	newConnectorBalance, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, r.ConnectorNativeAddr)
	require.NoError(r, err, "BalanceOf failed for new connector")

	require.Zero(r,
		expectedBalance.Cmp(newConnectorBalance),
		"Migration failed: expected %s, got %s",
		expectedBalance.String(), newConnectorBalance.String())

	r.Logger.Print("✅ Migration verification successful: %s ZETA tokens migrated", newConnectorBalance.String())
}
