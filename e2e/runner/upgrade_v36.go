package runner

import (
	"math/big"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/utils"
)

const V36Version = "v36.0.0"

func (r *E2ERunner) RunSetup(testLegacy bool) {
	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
		r.requireTxSuccessful(receipt, failMessage)
	}

	r.UpgradeGatewayEVM()
	updateAdditionalFeeTx, err := r.GatewayEVM.UpdateAdditionalActionFee(r.EVMAuth, big.NewInt(2e5))
	require.NoError(r, err)
	ensureTxReceipt(updateAdditionalFeeTx, "Updating additional fee failed")
	r.UpgradeGatewayZEVM()
	r.UpdateProtocolContractsInChainParams(testLegacy)
	r.DeployTestDAppV2ZEVM()
	r.DeployTestDAppV2EVM()
}
