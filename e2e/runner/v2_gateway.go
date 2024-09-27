package runner

import (
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayevm.sol"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/utils"
)

// UpgradeGateways upgrades the GatewayEVM and GatewayZEVM contracts
// It deploy new gateway contract implementation with the current imported artifacts and upgrade the gateway contract
func (r *E2ERunner) UpgradeGateways() {
	r.UpgradeGatewayZEVM()
	r.UpgradeGatewayEVM()
}

// UpgradeGatewayZEVM upgrades the GatewayZEVM contract
func (r *E2ERunner) UpgradeGatewayZEVM() {
	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		r.requireTxSuccessful(receipt, failMessage+" tx hash: "+tx.Hash().Hex())
	}

	r.Logger.Info("Upgrading Gateway ZEVM contract")
	// Deploy the new gateway contract implementation
	newImplementationAddress, txDeploy, _, err := gatewayzevm.DeployGatewayZEVM(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	ensureTxReceipt(txDeploy, "New GatewayZEVM implementation deployment failed")

	// Upgrade
	txUpgrade, err := r.GatewayZEVM.UpgradeToAndCall(r.ZEVMAuth, newImplementationAddress, []byte{})
	require.NoError(r, err)
	ensureTxReceipt(txUpgrade, "GatewayZEVM upgrade failed")
}

// UpgradeGatewayEVM upgrades the GatewayEVM contract
func (r *E2ERunner) UpgradeGatewayEVM() {
	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
		r.requireTxSuccessful(receipt, failMessage+" tx hash: "+tx.Hash().Hex())
	}

	r.Logger.Info("Upgrading Gateway EVM contract")
	// Deploy the new gateway contract implementation
	newImplementationAddress, txDeploy, _, err := gatewayevm.DeployGatewayEVM(r.EVMAuth, r.EVMClient)
	require.NoError(r, err)
	ensureTxReceipt(txDeploy, "New GatewayEVM implementation deployment failed")

	// Upgrade
	txUpgrade, err := r.GatewayEVM.UpgradeToAndCall(r.EVMAuth, newImplementationAddress, []byte{})
	require.NoError(r, err)
	ensureTxReceipt(txUpgrade, "GatewayEVM upgrade failed")
}
