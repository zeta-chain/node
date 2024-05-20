package e2etests

// WithdrawBitcoinMultipleTimes ...
// TODO: complete and uncomment E2E test
// https://github.com/zeta-chain/node-private/issues/79
//func WithdrawBitcoinMultipleTimes(r *runner.E2ERunner, repeat int64) {
//	totalAmount := big.NewInt(int64(0.1 * 1e8))
//
//	// #nosec G701 test - always in range
//	amount := big.NewInt(int64(0.1 * 1e8 / float64(repeat)))
//
//	// check if the deposit is successful
//	BTCZRC20Addr, err := r.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(common.BtcRegtestChain.ChainId))
//	if err != nil {
//		panic(err)
//	}
//	r.Logger.Info("BTCZRC20Addr: %s", BTCZRC20Addr.Hex())
//	BTCZRC20, err := zrc20.NewZRC20(BTCZRC20Addr, r.ZEVMClient)
//	if err != nil {
//		panic(err)
//	}
//	balance, err := BTCZRC20.BalanceOf(&bind.CallOpts{}, r.DeployerAddress)
//	if err != nil {
//		panic(err)
//	}
//	if balance.Cmp(totalAmount) < 0 {
//		panic(fmt.Errorf("not enough balance in ZRC20 contract"))
//	}
//	// approve the ZRC20 contract to spend 1 BTC from the deployer address
//	{
//		// approve more to cover withdraw fee
//		tx, err := BTCZRC20.Approve(r.ZEVMAuth, BTCZRC20Addr, totalAmount.Mul(totalAmount, big.NewInt(100)))
//		if err != nil {
//			panic(err)
//		}
//		receipt := config.MustWaitForTxReceipt(r.ZEVMClient, tx, r.Logger)
//		r.Logger.Info("approve receipt: status %d", receipt.Status)
//		if receipt.Status != 1 {
//			panic(fmt.Errorf("approve receipt status is not 1"))
//		}
//	}
//	go func() {
//		for {
//			time.Sleep(3 * time.Second)
//			_, err = r.BtcRPCClient.GenerateToAddress(1, r.BTCDeployerAddress, nil)
//			if err != nil {
//				panic(err)
//			}
//		}
//	}()
//	// withdraw 0.1 BTC from ZRC20 to BTC address
//	for i := int64(0); i < repeat; i++ {
//		_, gasFee, err := BTCZRC20.WithdrawGasFee(&bind.CallOpts{})
//		if err != nil {
//			panic(err)
//		}
//		r.Logger.Info("withdraw gas fee: %d", gasFee)
//		tx, err := BTCZRC20.Withdraw(r.ZEVMAuth, []byte(r.BTCDeployerAddress.EncodeAddress()), amount)
//		if err != nil {
//			panic(err)
//		}
//		receipt := config.MustWaitForTxReceipt(r.ZEVMClient, tx, r.Logger)
//		r.Logger.Info("withdraw receipt: status %d", receipt.Status)
//		if receipt.Status != 1 {
//			panic(fmt.Errorf("withdraw receipt status is not 1"))
//		}
//		_, err = r.BtcRPCClient.GenerateToAddress(10, r.BTCDeployerAddress, nil)
//		if err != nil {
//			panic(err)
//		}
//		cctx := config.WaitCctxMinedByInTxHash(receipt.TxHash.Hex(), r.CctxClient, r.Logger)
//		outTxHash := cctx.GetCurrentOutTxParam().OutboundTxHash
//		hash, err := chainhash.NewHashFromStr(outTxHash)
//		if err != nil {
//			panic(err)
//		}
//
//		rawTx, err := r.BtcRPCClient.GetRawTransactionVerbose(hash)
//		if err != nil {
//			panic(err)
//		}
//		r.Logger.Info("raw tx:")
//		r.Logger.Info("  TxIn: %d", len(rawTx.Vin))
//		for idx, txIn := range rawTx.Vin {
//			r.Logger.Info("  TxIn %d:", idx)
//			r.Logger.Info("    TxID:Vout:  %s:%d", txIn.Txid, txIn.Vout)
//			r.Logger.Info("    ScriptSig: %s", txIn.ScriptSig.Hex)
//		}
//		r.Logger.Info("  TxOut: %d", len(rawTx.Vout))
//		for _, txOut := range rawTx.Vout {
//			r.Logger.Info("  TxOut %d:", txOut.N)
//			r.Logger.Info("    Value: %.8f", txOut.Value)
//			r.Logger.Info("    ScriptPubKey: %s", txOut.ScriptPubKey.Hex)
//		}
//	}
//}
