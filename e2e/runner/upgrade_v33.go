package runner

import (
	"math/big"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/utils"
)

func (r *E2ERunner) RunSetup(testLegacy bool) {
	ensureReceiptEVM := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
		msg := "tx %s receipt status is not successful: %s"
		require.Equal(
			r,
			ethtypes.ReceiptStatusSuccessful,
			receipt.Status,
			msg,
			receipt.TxHash.String(),
			failMessage,
		)
	}
	ensureReceiptZEVM := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		msg := "tx %s receipt status is not successful: %s"
		require.Equal(
			r,
			ethtypes.ReceiptStatusSuccessful,
			receipt.Status,
			msg,
			receipt.TxHash.String(),
			failMessage,
		)
	}
	r.UpgradeGatewayEVM()
	updateAdditionalFeeTx, err := r.GatewayEVM.UpdateAdditionalActionFee(r.EVMAuth, big.NewInt(2e5))
	require.NoError(r, err)
	ensureReceiptEVM(updateAdditionalFeeTx, "Update additional fee failed")
	r.UpgradeGatewayZEVM()
	r.DeployZetaConnectorNative(ensureReceiptEVM)
	r.UpdateProtocolContractsInChainParams(testLegacy)
	r.SetupZEVMTestDappV2(ensureReceiptZEVM)
	r.DeployTestDAppV2(ensureReceiptEVM)
}
