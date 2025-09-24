package runner

import (
	"math/big"

	"github.com/stretchr/testify/require"
)

const V36Version = "v36.0.0"

func (r *E2ERunner) RunSetup(testLegacy bool) {
	r.UpgradeGatewayEVM()
	updateAdditionalFeeTx, err := r.GatewayEVM.UpdateAdditionalActionFee(r.EVMAuth, big.NewInt(2e5))
	require.NoError(r, err)
	r.ensureTxReceiptEVM(updateAdditionalFeeTx, "Updating additional fee failed")
	r.UpgradeGatewayZEVM()
	//r.UpdateProtocolContractsInChainParams(testLegacy)
	r.DeployCoreRegistry()
	r.ActivateChainsOnRegistry()
	r.DeployTestDAppV2ZEVM()
	r.DeployTestDAppV2EVM()
}
