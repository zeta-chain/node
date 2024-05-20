package e2etests

// DepositBTCRefund ...
// TODO: define e2e test
// https://github.com/zeta-chain/node-private/issues/79
//func DepositBTCRefund(r *runner.E2ERunner) {
//	r.Logger.InfoLoud("Deposit BTC with invalid memo; should be refunded")
//	btc := r.BtcRPCClient
//	utxos, err := r.BtcRPCClient.ListUnspent()
//	if err != nil {
//		panic(err)
//	}
//	spendableAmount := 0.0
//	spendableUTXOs := 0
//	for _, utxo := range utxos {
//		if utxo.Spendable {
//			spendableAmount += utxo.Amount
//			spendableUTXOs++
//		}
//	}
//	r.Logger.Info("ListUnspent:")
//	r.Logger.Info("  spendableAmount: %f", spendableAmount)
//	r.Logger.Info("  spendableUTXOs: %d", spendableUTXOs)
//	r.Logger.Info("Now sending two txs to TSS address...")
//	_, err = r.SendToTSSFromDeployerToDeposit(r.BTCTSSAddress, 1.1, utxos[:2], btc, r.BTCDeployerAddress)
//	if err != nil {
//		panic(err)
//	}
//	_, err = r.SendToTSSFromDeployerToDeposit(r.BTCTSSAddress, 0.05, utxos[2:4], btc, r.BTCDeployerAddress)
//	if err != nil {
//		panic(err)
//	}
//
//	r.Logger.Info("testing if the deposit into BTC ZRC20 is successful...")
//
//	// check if the deposit is successful
//	initialBalance, err := r.BTCZRC20.BalanceOf(&bind.CallOpts{}, r.DeployerAddress)
//	if err != nil {
//		panic(err)
//	}
//	for {
//		time.Sleep(3 * time.Second)
//		balance, err := r.BTCZRC20.BalanceOf(&bind.CallOpts{}, r.DeployerAddress)
//		if err != nil {
//			panic(err)
//		}
//		diff := big.NewInt(0)
//		diff.Sub(balance, initialBalance)
//		if diff.Cmp(big.NewInt(1.15*btcutil.SatoshiPerBitcoin)) != 0 {
//			r.Logger.Info("waiting for BTC balance to show up in ZRC contract... current bal %d", balance)
//		} else {
//			r.Logger.Info("BTC balance is in ZRC contract! Success")
//			break
//		}
//	}
//}
