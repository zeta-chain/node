package e2etests

import (
	"fmt"
	"math/big"

	"github.com/zeta-chain/zetacore/e2e/contracts/testdapp"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	cctxtypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMessagePassingZEVM(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestMessagePassing requires exactly one argument for the amount.")
	}

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestMessagePassing.")
	}

	zEVMChainID, err := r.ZEVMClient.ChainID(r.Ctx)
	if err != nil {
		panic(err)
	}
	destinationAddress := r.ZevmTestDAppAddr

	//Use TestDapp to call the Send function on the EVM connector to create a message
	auth := r.EVMAuth

	tx, err := r.ZetaEth.Approve(auth, r.EvmTestDAppAddr, amount)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("tx failed")
	}
	r.Logger.Info("Approve tx receipt: %d", receipt.Status)

	testDAppEVM, err := testdapp.NewTestDApp(r.EvmTestDAppAddr, r.EVMClient)
	if err != nil {
		panic(err)
	}

	tx, err = testDAppEVM.SendHelloWorld(auth, destinationAddress, zEVMChainID, amount, false)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("TestDApp.SendHello tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)

	r.Logger.Print(fmt.Sprintf("Waiting for cctx , intx : %s", receipt.TxHash.String()))

	// New inbound message picked up by zeta-clients and voted on by observers to initiate a contract call on zEVM
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx.CctxStatus.Status != cctxtypes.CctxStatus_OutboundMined {
		panic("expected cctx to be outbound_mined")
	}
	r.Logger.Print("🔄 Cctx mined by outbound chain zevm", cctx.Index)

}
