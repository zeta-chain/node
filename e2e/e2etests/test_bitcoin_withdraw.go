package e2etests

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/pkg/chains"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

func TestBitcoinWithdrawSegWit(r *runner.E2ERunner, args []string) {
	// check length of arguments
	if len(args) != 2 {
		panic("TestBitcoinWithdrawSegWit requires two arguments: [receiver, amount]")
	}
	r.SetBtcAddress(r.Name, false)

	// parse arguments
	defaultReceiver := r.BTCDeployerAddress.EncodeAddress()
	receiver, amount := parseBitcoinWithdrawArgs(args, defaultReceiver)
	_, ok := receiver.(*btcutil.AddressWitnessPubKeyHash)
	if !ok {
		panic("Invalid receiver address specified for TestBitcoinWithdrawSegWit.")
	}

	withdrawBTCZRC20(r, receiver, amount)
}

func TestBitcoinWithdrawTaproot(r *runner.E2ERunner, args []string) {
	// check length of arguments
	if len(args) != 2 {
		panic("TestBitcoinWithdrawTaproot requires two arguments: [receiver, amount]")
	}
	r.SetBtcAddress(r.Name, false)

	// parse arguments and withdraw BTC
	defaultReceiver := "bcrt1pqqqsyqcyq5rqwzqfpg9scrgwpugpzysnzs23v9ccrydpk8qarc0sj9hjuh"
	receiver, amount := parseBitcoinWithdrawArgs(args, defaultReceiver)
	_, ok := receiver.(*chains.AddressTaproot)
	if !ok {
		panic("Invalid receiver address specified for TestBitcoinWithdrawTaproot.")
	}

	withdrawBTCZRC20(r, receiver, amount)
}

func TestBitcoinWithdrawLegacy(r *runner.E2ERunner, args []string) {
	// check length of arguments
	if len(args) != 2 {
		panic("TestBitcoinWithdrawLegacy requires two arguments: [receiver, amount]")
	}
	r.SetBtcAddress(r.Name, false)

	// parse arguments and withdraw BTC
	defaultReceiver := "mxpYha3UJKUgSwsAz2qYRqaDSwAkKZ3YEY"
	receiver, amount := parseBitcoinWithdrawArgs(args, defaultReceiver)
	_, ok := receiver.(*btcutil.AddressPubKeyHash)
	if !ok {
		panic("Invalid receiver address specified for TestBitcoinWithdrawLegacy.")
	}

	withdrawBTCZRC20(r, receiver, amount)
}

func TestBitcoinWithdrawP2WSH(r *runner.E2ERunner, args []string) {
	// check length of arguments
	if len(args) != 2 {
		panic("TestBitcoinWithdrawP2WSH requires two arguments: [receiver, amount]")
	}
	r.SetBtcAddress(r.Name, false)

	// parse arguments and withdraw BTC
	defaultReceiver := "bcrt1qm9mzhyky4w853ft2ms6dtqdyyu3z2tmrq8jg8xglhyuv0dsxzmgs2f0sqy"
	receiver, amount := parseBitcoinWithdrawArgs(args, defaultReceiver)
	_, ok := receiver.(*btcutil.AddressWitnessScriptHash)
	if !ok {
		panic("Invalid receiver address specified for TestBitcoinWithdrawP2WSH.")
	}

	withdrawBTCZRC20(r, receiver, amount)
}

func TestBitcoinWithdrawP2SH(r *runner.E2ERunner, args []string) {
	// check length of arguments
	if len(args) != 2 {
		panic("TestBitcoinWithdrawP2SH requires two arguments: [receiver, amount]")
	}
	r.SetBtcAddress(r.Name, false)

	// parse arguments and withdraw BTC
	defaultReceiver := "2N6AoUj3KPS7wNGZXuCckh8YEWcSYNsGbqd"
	receiver, amount := parseBitcoinWithdrawArgs(args, defaultReceiver)
	_, ok := receiver.(*btcutil.AddressScriptHash)
	if !ok {
		panic("Invalid receiver address specified for TestBitcoinWithdrawP2SH.")
	}

	withdrawBTCZRC20(r, receiver, amount)
}

func TestBitcoinWithdrawRestricted(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestBitcoinWithdrawRestricted requires exactly one argument for the amount.")
	}

	withdrawalAmount, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		panic("Invalid withdrawal amount specified for TestBitcoinWithdrawRestricted.")
	}

	withdrawalAmountSat, err := btcutil.NewAmount(withdrawalAmount)
	if err != nil {
		panic(err)
	}
	amount := big.NewInt(int64(withdrawalAmountSat))

	r.SetBtcAddress(r.Name, false)

	withdrawBitcoinRestricted(r, amount)
}

func parseBitcoinWithdrawArgs(args []string, defaultReceiver string) (btcutil.Address, *big.Int) {
	// parse receiver address
	var err error
	var receiver btcutil.Address
	if args[0] == "" {
		// use the default receiver
		receiver, err = chains.DecodeBtcAddress(defaultReceiver, chains.BtcRegtestChain.ChainId)
		if err != nil {
			panic("Invalid default receiver address specified for TestBitcoinWithdraw.")
		}
	} else {
		receiver, err = chains.DecodeBtcAddress(args[0], chains.BtcRegtestChain.ChainId)
		if err != nil {
			panic("Invalid receiver address specified for TestBitcoinWithdraw.")
		}
	}
	// parse the withdrawal amount
	withdrawalAmount, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		panic("Invalid withdrawal amount specified for TestBitcoinWithdraw.")
	}
	withdrawalAmountSat, err := btcutil.NewAmount(withdrawalAmount)
	if err != nil {
		panic(err)
	}
	amount := big.NewInt(int64(withdrawalAmountSat))

	return receiver, amount
}

func withdrawBTCZRC20(r *runner.E2ERunner, to btcutil.Address, amount *big.Int) *btcjson.TxRawResult {
	tx, err := r.BTCZRC20.Approve(r.ZEVMAuth, r.BTCZRC20Addr, big.NewInt(amount.Int64()*2)) // approve more to cover withdraw fee
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		panic(fmt.Errorf("approve receipt status is not 1"))
	}

	// mine blocks
	stop := r.MineBlocks()

	// withdraw 'amount' of BTC from ZRC20 to BTC address
	tx, err = r.BTCZRC20.Withdraw(r.ZEVMAuth, []byte(to.EncodeAddress()), amount)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		panic(fmt.Errorf("withdraw receipt status is not 1"))
	}

	// mine 10 blocks to confirm the withdraw tx
	_, err = r.BtcRPCClient.GenerateToAddress(10, to, nil)
	if err != nil {
		panic(err)
	}

	// get cctx and check status
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, receipt.TxHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		panic(fmt.Errorf("cctx status is not OutboundMined"))
	}

	// get bitcoin tx according to the outboundHash in cctx
	outboundHash := cctx.GetCurrentOutboundParam().Hash
	hash, err := chainhash.NewHashFromStr(outboundHash)
	if err != nil {
		panic(err)
	}

	rawTx, err := r.BtcRPCClient.GetRawTransactionVerbose(hash)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("raw tx:")
	r.Logger.Info("  TxIn: %d", len(rawTx.Vin))
	for idx, txIn := range rawTx.Vin {
		r.Logger.Info("  TxIn %d:", idx)
		r.Logger.Info("    TxID:Vout:  %s:%d", txIn.Txid, txIn.Vout)
		r.Logger.Info("    ScriptSig: %s", txIn.ScriptSig.Hex)
	}
	r.Logger.Info("  TxOut: %d", len(rawTx.Vout))
	for _, txOut := range rawTx.Vout {
		r.Logger.Info("  TxOut %d:", txOut.N)
		r.Logger.Info("    Value: %.8f", txOut.Value)
		r.Logger.Info("    ScriptPubKey: %s", txOut.ScriptPubKey.Hex)
	}

	// stop mining
	stop <- struct{}{}

	return rawTx
}

func withdrawBitcoinRestricted(r *runner.E2ERunner, amount *big.Int) {
	// use restricted BTC P2WPKH address
	addressRestricted, err := chains.DecodeBtcAddress(testutils.RestrictedBtcAddressTest, chains.BtcRegtestChain.ChainId)
	if err != nil {
		panic(err)
	}

	// the cctx should be cancelled
	rawTx := withdrawBTCZRC20(r, addressRestricted, amount)
	if len(rawTx.Vout) != 2 {
		panic(fmt.Errorf("BTC cancelled outtx rawTx.Vout should have 2 outputs"))
	}
}

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
//		cctx := config.WaitCctxMinedByInboundHash(receipt.TxHash.Hex(), r.CctxClient, r.Logger)
//		outboundHash := cctx.GetCurrentOutboundParam().Hash
//		hash, err := chainhash.NewHashFromStr(outboundHash)
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
