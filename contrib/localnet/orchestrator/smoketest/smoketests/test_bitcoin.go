package smoketests

import (
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
)

func TestBitcoinWithdraw(sm *runner.SmokeTestRunner) {
	// withdraw 0.1 BTC from ZRC20 to BTC address
	// first, approve the ZRC20 contract to spend 1 BTC from the deployer address
	WithdrawBitcoin(sm)
}

func WithdrawBitcoin(sm *runner.SmokeTestRunner) {
	amount := big.NewInt(0.1 * btcutil.SatoshiPerBitcoin)

	// approve the ZRC20 contract to spend 1 BTC from the deployer address
	tx, err := sm.BTCZRC20.Approve(sm.ZevmAuth, sm.BTCZRC20Addr, big.NewInt(amount.Int64()*2)) // approve more to cover withdraw fee
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.ZevmClient, tx, sm.Logger)
	if receipt.Status != 1 {
		panic(fmt.Errorf("approve receipt status is not 1"))
	}

	// mine blocks
	stop := sm.MineBlocks()

	// withdraw 0.1 BTC from ZRC20 to BTC address
	tx, err = sm.BTCZRC20.Withdraw(sm.ZevmAuth, []byte(sm.BTCDeployerAddress.EncodeAddress()), amount)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.Ctx, sm.ZevmClient, tx, sm.Logger)
	if receipt.Status != 1 {
		panic(fmt.Errorf("withdraw receipt status is not 1"))
	}

	// mine 10 blocks to confirm the withdraw tx
	_, err = sm.BtcRPCClient.GenerateToAddress(10, sm.BTCDeployerAddress, nil)
	if err != nil {
		panic(err)
	}

	cctx := utils.WaitCctxMinedByInTxHash(sm.Ctx, receipt.TxHash.Hex(), sm.CctxClient, sm.Logger)
	outTxHash := cctx.GetCurrentOutTxParam().OutboundTxHash
	hash, err := chainhash.NewHashFromStr(outTxHash)
	if err != nil {
		panic(err)
	}

	rawTx, err := sm.BtcRPCClient.GetRawTransactionVerbose(hash)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("raw tx:")
	sm.Logger.Info("  TxIn: %d", len(rawTx.Vin))
	for idx, txIn := range rawTx.Vin {
		sm.Logger.Info("  TxIn %d:", idx)
		sm.Logger.Info("    TxID:Vout:  %s:%d", txIn.Txid, txIn.Vout)
		sm.Logger.Info("    ScriptSig: %s", txIn.ScriptSig.Hex)
	}
	sm.Logger.Info("  TxOut: %d", len(rawTx.Vout))
	for _, txOut := range rawTx.Vout {
		sm.Logger.Info("  TxOut %d:", txOut.N)
		sm.Logger.Info("    Value: %.8f", txOut.Value)
		sm.Logger.Info("    ScriptPubKey: %s", txOut.ScriptPubKey.Hex)
	}

	// stop mining
	stop <- struct{}{}
}

// WithdrawBitcoinMultipleTimes ...
// TODO: define smoke test
// https://github.com/zeta-chain/node-private/issues/79
//func WithdrawBitcoinMultipleTimes(sm *runner.SmokeTestRunner, repeat int64) {
//	totalAmount := big.NewInt(int64(0.1 * 1e8))
//
//	// #nosec G701 smoketest - always in range
//	amount := big.NewInt(int64(0.1 * 1e8 / float64(repeat)))
//
//	// check if the deposit is successful
//	BTCZRC20Addr, err := sm.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(common.BtcRegtestChain().ChainId))
//	if err != nil {
//		panic(err)
//	}
//	sm.Logger.Info("BTCZRC20Addr: %s", BTCZRC20Addr.Hex())
//	BTCZRC20, err := zrc20.NewZRC20(BTCZRC20Addr, sm.ZevmClient)
//	if err != nil {
//		panic(err)
//	}
//	balance, err := BTCZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
//	if err != nil {
//		panic(err)
//	}
//	if balance.Cmp(totalAmount) < 0 {
//		panic(fmt.Errorf("not enough balance in ZRC20 contract"))
//	}
//	// approve the ZRC20 contract to spend 1 BTC from the deployer address
//	{
//		// approve more to cover withdraw fee
//		tx, err := BTCZRC20.Approve(sm.ZevmAuth, BTCZRC20Addr, totalAmount.Mul(totalAmount, big.NewInt(100)))
//		if err != nil {
//			panic(err)
//		}
//		receipt := utils.MustWaitForTxReceipt(sm.ZevmClient, tx, sm.Logger)
//		sm.Logger.Info("approve receipt: status %d", receipt.Status)
//		if receipt.Status != 1 {
//			panic(fmt.Errorf("approve receipt status is not 1"))
//		}
//	}
//	go func() {
//		for {
//			time.Sleep(3 * time.Second)
//			_, err = sm.BtcRPCClient.GenerateToAddress(1, sm.BTCDeployerAddress, nil)
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
//		sm.Logger.Info("withdraw gas fee: %d", gasFee)
//		tx, err := BTCZRC20.Withdraw(sm.ZevmAuth, []byte(sm.BTCDeployerAddress.EncodeAddress()), amount)
//		if err != nil {
//			panic(err)
//		}
//		receipt := utils.MustWaitForTxReceipt(sm.ZevmClient, tx, sm.Logger)
//		sm.Logger.Info("withdraw receipt: status %d", receipt.Status)
//		if receipt.Status != 1 {
//			panic(fmt.Errorf("withdraw receipt status is not 1"))
//		}
//		_, err = sm.BtcRPCClient.GenerateToAddress(10, sm.BTCDeployerAddress, nil)
//		if err != nil {
//			panic(err)
//		}
//		cctx := utils.WaitCctxMinedByInTxHash(receipt.TxHash.Hex(), sm.CctxClient, sm.Logger)
//		outTxHash := cctx.GetCurrentOutTxParam().OutboundTxHash
//		hash, err := chainhash.NewHashFromStr(outTxHash)
//		if err != nil {
//			panic(err)
//		}
//
//		rawTx, err := sm.BtcRPCClient.GetRawTransactionVerbose(hash)
//		if err != nil {
//			panic(err)
//		}
//		sm.Logger.Info("raw tx:")
//		sm.Logger.Info("  TxIn: %d", len(rawTx.Vin))
//		for idx, txIn := range rawTx.Vin {
//			sm.Logger.Info("  TxIn %d:", idx)
//			sm.Logger.Info("    TxID:Vout:  %s:%d", txIn.Txid, txIn.Vout)
//			sm.Logger.Info("    ScriptSig: %s", txIn.ScriptSig.Hex)
//		}
//		sm.Logger.Info("  TxOut: %d", len(rawTx.Vout))
//		for _, txOut := range rawTx.Vout {
//			sm.Logger.Info("  TxOut %d:", txOut.N)
//			sm.Logger.Info("    Value: %.8f", txOut.Value)
//			sm.Logger.Info("    ScriptPubKey: %s", txOut.ScriptPubKey.Hex)
//		}
//	}
//}

// DepositBTCRefund ...
// TODO: define smoke test
// https://github.com/zeta-chain/node-private/issues/79
//func DepositBTCRefund(sm *runner.SmokeTestRunner) {
//	sm.Logger.InfoLoud("Deposit BTC with invalid memo; should be refunded")
//	btc := sm.BtcRPCClient
//	utxos, err := sm.BtcRPCClient.ListUnspent()
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
//	sm.Logger.Info("ListUnspent:")
//	sm.Logger.Info("  spendableAmount: %f", spendableAmount)
//	sm.Logger.Info("  spendableUTXOs: %d", spendableUTXOs)
//	sm.Logger.Info("Now sending two txs to TSS address...")
//	_, err = sm.SendToTSSFromDeployerToDeposit(sm.BTCTSSAddress, 1.1, utxos[:2], btc, sm.BTCDeployerAddress)
//	if err != nil {
//		panic(err)
//	}
//	_, err = sm.SendToTSSFromDeployerToDeposit(sm.BTCTSSAddress, 0.05, utxos[2:4], btc, sm.BTCDeployerAddress)
//	if err != nil {
//		panic(err)
//	}
//
//	sm.Logger.Info("testing if the deposit into BTC ZRC20 is successful...")
//
//	// check if the deposit is successful
//	initialBalance, err := sm.BTCZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
//	if err != nil {
//		panic(err)
//	}
//	for {
//		time.Sleep(3 * time.Second)
//		balance, err := sm.BTCZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
//		if err != nil {
//			panic(err)
//		}
//		diff := big.NewInt(0)
//		diff.Sub(balance, initialBalance)
//		if diff.Cmp(big.NewInt(1.15*btcutil.SatoshiPerBitcoin)) != 0 {
//			sm.Logger.Info("waiting for BTC balance to show up in ZRC contract... current bal %d", balance)
//		} else {
//			sm.Logger.Info("BTC balance is in ZRC contract! Success")
//			break
//		}
//	}
//}
