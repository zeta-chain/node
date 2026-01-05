package runner

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/contracts/testdapp"
	"github.com/zeta-chain/node/e2e/contracts/zevmswap"
	e2eutils "github.com/zeta-chain/node/e2e/utils"
)

// SetupLegacyZEVMContracts sets up the legacy contracts on ZEVM
// In particular it deploys test contracts used with the protocol contracts v1
func (r *E2ERunner) SetupLegacyZEVMContracts() {
	r.Logger.Print("⚙️ Setting up ZEVM legacy protocol contracts")
	// deploy TestDApp contract on zEVM
	appAddr, txApp, _, err := testdapp.DeployTestDApp(
		r.ZEVMAuth,
		r.ZEVMClient,
		r.ConnectorZEVMAddr,
		r.WZetaAddr,
	)
	require.NoError(r, err)

	r.ZevmTestDAppAddr = appAddr
	r.Logger.Info("TestDApp Zevm contract address: %s, tx hash: %s", appAddr.Hex(), txApp.Hash().Hex())

	// deploy ZEVMSwapApp
	zevmSwapAppAddr, txZEVMSwapApp, zevmSwapApp, err := zevmswap.DeployZEVMSwapApp(
		r.ZEVMAuth,
		r.ZEVMClient,
		r.UniswapV2RouterAddr,
		r.SystemContractAddr,
	)
	require.NoError(r, err)

	receipt := e2eutils.MustWaitForTxReceipt(
		r.Ctx,
		r.ZEVMClient,
		txZEVMSwapApp,
		r.Logger,
		r.ReceiptTimeout,
	)
	r.requireTxSuccessful(receipt, "ZEVMSwapApp deployment failed")

	r.ZEVMSwapAppAddr = zevmSwapAppAddr
	r.ZEVMSwapApp = zevmSwapApp
}
