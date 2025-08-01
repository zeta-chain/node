package runner

import (
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/utils"
)

func (r *E2ERunner) RunSetup(testLegacy bool) {
	ensureReceiptEVM := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
		msg := "receipt status is not successful: %s"
		require.Equal(
			r,
			ethtypes.ReceiptStatusSuccessful,
			receipt.Status,
			msg,
			receipt.TxHash.String(),
		)
	}
	ensureReceiptZEVM := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		msg := "receipt status is not successful: %s"
		require.Equal(
			r,
			ethtypes.ReceiptStatusSuccessful,
			receipt.Status,
			msg,
			receipt.TxHash.String(),
		)
	}
	r.UpgradeGatewayEVM()
	r.UpgradeGatewayZEVM()
	r.DeployZetaConnectorNative(ensureReceiptEVM)
	r.UpdateProtocolContractsInChainParams(testLegacy)
	r.SetupZEVMTestDappV2(ensureReceiptZEVM)
	r.DeployTestDAppV2(ensureReceiptEVM)
}
