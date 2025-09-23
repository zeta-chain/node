package runner

const V36Version = "v36.0.0"

func (r *E2ERunner) RunSetup(testLegacy bool) {
	r.UpgradeGatewayEVM()
	r.UpgradeGatewayZEVM()
	//r.UpdateProtocolContractsInChainParams(testLegacy)
	r.DeployCoreRegistry()
	r.ActivateChainsOnRegistry()
	r.DeployTestDAppV2ZEVM()
	r.DeployTestDAppV2EVM()
}
