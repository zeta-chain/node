package runner

import (
	"fmt"
	"os"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayevm.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"
	"golang.org/x/mod/semver"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/utils"
)

// UpgradeGatewayOptions is the options for the gateway upgrade tests
type UpgradeGatewayOptions struct {
	TestSolana bool
	TestSui    bool
}

// UpgradeGatewaysAndERC20Custody upgrades gateways and ERC20Custody contracts
// It deploys new contract implementation with the current imported artifacts and upgrades the contract
func (r *E2ERunner) UpgradeGatewaysAndERC20Custody() {
	r.UpgradeGatewayZEVM()
	r.UpgradeGatewayEVM()
	r.UpgradeERC20Custody()

	// Ensure gateway use 4M for gas limit
	err := r.ZetaTxServer.UpdateGatewayGasLimit(GatewayGasLimit)
	require.NoError(r, err)
}

// RunGatewayUpgradeTestsExternalChains runs the gateway upgrade tests for external chains
func (r *E2ERunner) RunGatewayUpgradeTestsExternalChains(conf config.Config, opts UpgradeGatewayOptions) {
	// Skip upgrades if this is the second run of the upgrade tests

	if opts.TestSolana {
		r.SolanaVerifyGatewayContractsUpgrade(conf.AdditionalAccounts.UserSolana.SolanaPrivateKey.String())
	}

	if opts.TestSui {
		r.SuiVerifyGatewayPackageUpgrade()
	}
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
	txUpgrade, err := r.ERC20Custody.UpgradeToAndCall(r.EVMAuth, newImplementationAddress, []byte{})
	require.NoError(r, err)
	ensureTxReceipt(txUpgrade, "ERC20Custody upgrade failed")
}

func (r *E2ERunner) AssertAfterUpgrade(assertVersion string, assertFunc func()) {
	version := r.GetZetacoredVersion()
	versionMajorIsZero := semver.Major(version) == "v0"
	oldVersion := fmt.Sprintf("v%s", os.Getenv("OLD_VERSION"))

	// run these assertions only on the second run of the upgrade
	if !r.IsRunningUpgrade() || !versionMajorIsZero || checkVersion(assertVersion, oldVersion) {
		return
	}
	r.Logger.Print("🏃 Running assertions after upgrade for version: %s", assertVersion)
	assertFunc()
}

// AddPreUpgradeHandler adds a handler to run any logic before an upgrade
func (r *E2ERunner) AddPreUpgradeHandler(upgradeFrom string, preHandler func()) {
	currentVersion := r.GetZetacoredVersion()
	// run these assertions only on the first run of the upgrade
	if !r.IsRunningUpgrade() || checkVersion(upgradeFrom, currentVersion) {
		return
	}
	preHandler()
}

// AddPostUpgradeHandler adds a handler to run any logic after and upgrade to enable tests to be executed
// Note This is handler is not related to the cosmos-sdk upgrade handler in any way
func (r *E2ERunner) AddPostUpgradeHandler(upgradeFrom string, postHandler func()) {
	//version := r.GetZetacoredVersion()
	//versionMajorIsZero := semver.Major(version) == "v0"
	//oldVersion := fmt.Sprintf("v%s", os.Getenv("OLD_VERSION"))
	//
	//// Run the handler only if this is the second run of the upgrade tests
	//if !r.IsRunningUpgrade() || !r.IsRunningTssMigration() || !versionMajorIsZero ||
	//	checkVersion(upgradeFrom, oldVersion) {
	//	return
	//}

	postHandler()
}

func checkVersion(upgradeFromm, oldVersion string) bool {
	//return semver.Major(upgradeFromm) != semver.Major(oldVersion)
	return false
}
