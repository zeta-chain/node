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
)

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
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, receipt.TxHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		panic(fmt.Errorf("cctx status is not OutboundMined"))
	}

	// get bitcoin tx according to the outTxHash in cctx
	outTxHash := cctx.GetCurrentOutTxParam().OutboundTxHash
	hash, err := chainhash.NewHashFromStr(outTxHash)
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
