package runner

func (r *E2ERunner) RunSetup(testLegacy bool) {
	r.UpgradeGatewayEVM()
	r.UpgradeGatewayZEVM()
	r.DeployCoreRegistry()
	r.ActivateChainsOnRegistry()
	r.DeployZetaConnectorNative()
	r.UpdateProtocolContractsInChainParams(testLegacy)
	r.DeployTestDAppV2ZEVM()
	r.DeployTestDAppV2EVM()
}
