package runner

import (
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayevm.sol"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/utils"
)

// UpgradeGatewaysAndERC20Custody upgrades gateways and ERC20Custody contracts
// It deploys new contract implementation with the current imported artifacts and upgrades the contract
func (r *E2ERunner) UpgradeGatewaysAndERC20Custody() {
	r.UpgradeGatewayZEVM()
	r.UpgradeGatewayEVM()
	r.UpgradeERC20Custody()
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

// UpgradeERC20Custody upgrades the ERC20Custody contract
func (r *E2ERunner) UpgradeERC20Custody() {
	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
		r.requireTxSuccessful(receipt, failMessage+" tx hash: "+tx.Hash().Hex())
	}

	r.Logger.Info("Upgrading ERC20Custody contract")
	// Deploy the new erc20Custody contract implementation
	newImplementationAddress, txDeploy, _, err := erc20custody.DeployERC20Custody(r.EVMAuth, r.EVMClient)
	require.NoError(r, err)
	ensureTxReceipt(txDeploy, "New ERC20Custody implementation deployment failed")

	// Upgrade
	txUpgrade, err := r.ERC20CustodyV2.UpgradeToAndCall(r.EVMAuth, newImplementationAddress, []byte{})
	require.NoError(r, err)
	ensureTxReceipt(txUpgrade, "ERC20Custody upgrade failed")
}
