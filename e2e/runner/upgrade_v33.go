package runner

func (r *E2ERunner) RunSetup(testLegacy bool) {
	//ensureReceiptEVM := func(tx *ethtypes.Transaction, failMessage string) {
	//	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	//	msg := "tx %s receipt status is not successful: %s"
	//	require.Equal(
	//		r,
	//		ethtypes.ReceiptStatusSuccessful,
	//		receipt.Status,
	//		msg,
	//		receipt.TxHash.String(),
	//		failMessage,
	//	)
	//}
	//ensureReceiptZEVM := func(tx *ethtypes.Transaction, failMessage string) {
	//	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	//	msg := "tx %s receipt status is not successful: %s"
	//	require.Equal(
	//		r,
	//		ethtypes.ReceiptStatusSuccessful,
	//		receipt.Status,
	//		msg,
	//		receipt.TxHash.String(),
	//		failMessage,
	//	)
	//}
	r.UpgradeGatewayEVM()
	r.UpgradeGatewayZEVM()
	r.DeployZetaConnectorNative()
	r.UpdateProtocolContractsInChainParams(testLegacy)
	r.DeployTestDAppV2ZEVM()
	r.DeployTestDAppV2EVM()
}
