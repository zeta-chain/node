package smoketests

import (
	"fmt"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	zrc20 "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
)

func TestBitcoinWithdraw(sm *runner.SmokeTestRunner) {
	startTime := time.Now()
	defer func() {
		fmt.Printf("Bitcoin withdraw took %s\n", time.Since(startTime))
	}()
	utils.LoudPrintf("Testing Bitcoin ZRC20 Withdraw...\n")
	// withdraw 0.1 BTC from ZRC20 to BTC address
	// first, approve the ZRC20 contract to spend 1 BTC from the deployer address
	WithdrawBitcoin(sm)
}

func WithdrawBitcoin(sm *runner.SmokeTestRunner) {
	amount := big.NewInt(0.1 * btcutil.SatoshiPerBitcoin)

	// check if the deposit is successful
	BTCZRC20Addr, err := sm.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(common.BtcRegtestChain().ChainId))
	if err != nil {
		panic(err)
	}
	fmt.Printf("BTCZRC20Addr: %s\n", BTCZRC20Addr.Hex())
	BTCZRC20, err := zrc20.NewZRC20(BTCZRC20Addr, sm.ZevmClient)
	if err != nil {
		panic(err)
	}
	balance, err := BTCZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	if balance.Cmp(amount) < 0 {
		panic(fmt.Errorf("not enough balance in ZRC20 contract"))
	}
	// approve the ZRC20 contract to spend 1 BTC from the deployer address
	{
		tx, err := BTCZRC20.Approve(sm.ZevmAuth, BTCZRC20Addr, big.NewInt(amount.Int64()*2)) // approve more to cover withdraw fee
		if err != nil {
			panic(err)
		}
		receipt := utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
		fmt.Printf("approve receipt: status %d\n", receipt.Status)
		if receipt.Status != 1 {
			panic(fmt.Errorf("approve receipt status is not 1"))
		}
	}
	go func() {
		for {
			time.Sleep(5 * time.Second)
			_, err = sm.BtcRPCClient.GenerateToAddress(1, sm.BTCDeployerAddress, nil)
			if err != nil {
				panic(err)
			}
		}
	}()
	// withdraw 0.1 BTC from ZRC20 to BTC address
	{
		_, gasFee, err := BTCZRC20.WithdrawGasFee(&bind.CallOpts{})
		if err != nil {
			panic(err)
		}
		fmt.Printf("withdraw gas fee: %d\n", gasFee)
		tx, err := BTCZRC20.Withdraw(sm.ZevmAuth, []byte(sm.BTCDeployerAddress.EncodeAddress()), amount)
		if err != nil {
			panic(err)
		}
		receipt := utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
		fmt.Printf("withdraw receipt: status %d\n", receipt.Status)
		if receipt.Status != 1 {
			panic(fmt.Errorf("withdraw receipt status is not 1"))
		}
		_, err = sm.BtcRPCClient.GenerateToAddress(10, sm.BTCDeployerAddress, nil)
		if err != nil {
			panic(err)
		}
		cctx := utils.WaitCctxMinedByInTxHash(receipt.TxHash.Hex(), sm.CctxClient)
		outTxHash := cctx.GetCurrentOutTxParam().OutboundTxHash
		hash, err := chainhash.NewHashFromStr(outTxHash)
		if err != nil {
			panic(err)
		}

		rawTx, err := sm.BtcRPCClient.GetRawTransactionVerbose(hash)
		if err != nil {
			panic(err)
		}
		fmt.Printf("raw tx:\n")
		fmt.Printf("  TxIn: %d\n", len(rawTx.Vin))
		for idx, txIn := range rawTx.Vin {
			fmt.Printf("  TxIn %d:\n", idx)
			fmt.Printf("    TxID:Vout:  %s:%d\n", txIn.Txid, txIn.Vout)
			fmt.Printf("    ScriptSig: %s\n", txIn.ScriptSig.Hex)
		}
		fmt.Printf("  TxOut: %d\n", len(rawTx.Vout))
		for _, txOut := range rawTx.Vout {
			fmt.Printf("  TxOut %d:\n", txOut.N)
			fmt.Printf("    Value: %.8f\n", txOut.Value)
			fmt.Printf("    ScriptPubKey: %s\n", txOut.ScriptPubKey.Hex)
		}
	}
}

// WithdrawBitcoinMultipleTimes ...
// TODO: define smoke test
// https://github.com/zeta-chain/node-private/issues/79
func WithdrawBitcoinMultipleTimes(sm *runner.SmokeTestRunner, repeat int64) {
	totalAmount := big.NewInt(int64(0.1 * 1e8))

	// #nosec G701 smoketest - always in range
	amount := big.NewInt(int64(0.1 * 1e8 / float64(repeat)))

	// check if the deposit is successful
	BTCZRC20Addr, err := sm.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(common.BtcRegtestChain().ChainId))
	if err != nil {
		panic(err)
	}
	fmt.Printf("BTCZRC20Addr: %s\n", BTCZRC20Addr.Hex())
	BTCZRC20, err := zrc20.NewZRC20(BTCZRC20Addr, sm.ZevmClient)
	if err != nil {
		panic(err)
	}
	balance, err := BTCZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	if balance.Cmp(totalAmount) < 0 {
		panic(fmt.Errorf("not enough balance in ZRC20 contract"))
	}
	// approve the ZRC20 contract to spend 1 BTC from the deployer address
	{
		// approve more to cover withdraw fee
		tx, err := BTCZRC20.Approve(sm.ZevmAuth, BTCZRC20Addr, totalAmount.Mul(totalAmount, big.NewInt(100)))
		if err != nil {
			panic(err)
		}
		receipt := utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
		fmt.Printf("approve receipt: status %d\n", receipt.Status)
		if receipt.Status != 1 {
			panic(fmt.Errorf("approve receipt status is not 1"))
		}
	}
	go func() {
		for {
			time.Sleep(3 * time.Second)
			_, err = sm.BtcRPCClient.GenerateToAddress(1, sm.BTCDeployerAddress, nil)
			if err != nil {
				panic(err)
			}
		}
	}()
	// withdraw 0.1 BTC from ZRC20 to BTC address
	for i := int64(0); i < repeat; i++ {
		_, gasFee, err := BTCZRC20.WithdrawGasFee(&bind.CallOpts{})
		if err != nil {
			panic(err)
		}
		fmt.Printf("withdraw gas fee: %d\n", gasFee)
		tx, err := BTCZRC20.Withdraw(sm.ZevmAuth, []byte(sm.BTCDeployerAddress.EncodeAddress()), amount)
		if err != nil {
			panic(err)
		}
		receipt := utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
		fmt.Printf("withdraw receipt: status %d\n", receipt.Status)
		if receipt.Status != 1 {
			panic(fmt.Errorf("withdraw receipt status is not 1"))
		}
		_, err = sm.BtcRPCClient.GenerateToAddress(10, sm.BTCDeployerAddress, nil)
		if err != nil {
			panic(err)
		}
		cctx := utils.WaitCctxMinedByInTxHash(receipt.TxHash.Hex(), sm.CctxClient)
		outTxHash := cctx.GetCurrentOutTxParam().OutboundTxHash
		hash, err := chainhash.NewHashFromStr(outTxHash)
		if err != nil {
			panic(err)
		}

		rawTx, err := sm.BtcRPCClient.GetRawTransactionVerbose(hash)
		if err != nil {
			panic(err)
		}
		fmt.Printf("raw tx:\n")
		fmt.Printf("  TxIn: %d\n", len(rawTx.Vin))
		for idx, txIn := range rawTx.Vin {
			fmt.Printf("  TxIn %d:\n", idx)
			fmt.Printf("    TxID:Vout:  %s:%d\n", txIn.Txid, txIn.Vout)
			fmt.Printf("    ScriptSig: %s\n", txIn.ScriptSig.Hex)
		}
		fmt.Printf("  TxOut: %d\n", len(rawTx.Vout))
		for _, txOut := range rawTx.Vout {
			fmt.Printf("  TxOut %d:\n", txOut.N)
			fmt.Printf("    Value: %.8f\n", txOut.Value)
			fmt.Printf("    ScriptPubKey: %s\n", txOut.ScriptPubKey.Hex)
		}
	}
}

// DepositBTCRefund ...
// TODO: define smoke test
// https://github.com/zeta-chain/node-private/issues/79
func DepositBTCRefund(sm *runner.SmokeTestRunner) {
	utils.LoudPrintf("Deposit BTC with invalid memo; should be refunded\n")
	btc := sm.BtcRPCClient
	utxos, err := sm.BtcRPCClient.ListUnspent()
	if err != nil {
		panic(err)
	}
	spendableAmount := 0.0
	spendableUTXOs := 0
	for _, utxo := range utxos {
		if utxo.Spendable {
			spendableAmount += utxo.Amount
			spendableUTXOs++
		}
	}
	fmt.Printf("ListUnspent:\n")
	fmt.Printf("  spendableAmount: %f\n", spendableAmount)
	fmt.Printf("  spendableUTXOs: %d\n", spendableUTXOs)
	fmt.Printf("Now sending two txs to TSS address...\n")
	_, err = sm.SendToTSSFromDeployerToDeposit(sm.BTCTSSAddress, 1.1, utxos[:2], btc, sm.BTCDeployerAddress)
	if err != nil {
		panic(err)
	}
	_, err = sm.SendToTSSFromDeployerToDeposit(sm.BTCTSSAddress, 0.05, utxos[2:4], btc, sm.BTCDeployerAddress)
	if err != nil {
		panic(err)
	}

	fmt.Printf("testing if the deposit into BTC ZRC20 is successful...\n")

	// check if the deposit is successful
	BTCZRC20Addr, err := sm.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(common.BtcRegtestChain().ChainId))
	if err != nil {
		panic(err)
	}
	sm.BTCZRC20Addr = BTCZRC20Addr
	fmt.Printf("BTCZRC20Addr: %s\n", BTCZRC20Addr.Hex())
	BTCZRC20, err := zrc20.NewZRC20(BTCZRC20Addr, sm.ZevmClient)
	if err != nil {
		panic(err)
	}
	sm.BTCZRC20 = BTCZRC20
	initialBalance, err := BTCZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	for {
		time.Sleep(5 * time.Second)
		balance, err := BTCZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
		if err != nil {
			panic(err)
		}
		diff := big.NewInt(0)
		diff.Sub(balance, initialBalance)
		if diff.Cmp(big.NewInt(1.15*btcutil.SatoshiPerBitcoin)) != 0 {
			fmt.Printf("waiting for BTC balance to show up in ZRC contract... current bal %d\n", balance)
		} else {
			fmt.Printf("BTC balance is in ZRC contract! Success\n")
			break
		}
	}
}
