package e2etests

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/e2e/contracts/testdapp"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	cctxtypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMessagePassingEVMtoZEVM(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestMessagePassingEVMtoZEVM requires exactly one argument for the amount.")
	}

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestMessagePassingEVMtoZEVM.")
	}

	// Set destination details
	zEVMChainID, err := r.ZEVMClient.ChainID(r.Ctx)
	if err != nil {
		panic(err)
	}

	destinationAddress := r.ZevmTestDAppAddr

	// Contract call originates from EVM chain
	tx, err := r.ZetaEth.Approve(r.EVMAuth, r.EvmTestDAppAddr, amount)
	if err != nil {
		panic(err)
	}

	r.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("approve tx failed")
	}

	r.Logger.Info("Approve tx receipt: %d", receipt.Status)
	testDAppEVM, err := testdapp.NewTestDApp(r.EvmTestDAppAddr, r.EVMClient)
	if err != nil {
		panic(err)
	}

	// Get ZETA balance on ZEVM TestDApp
	previousBalanceZEVM, err := r.WZeta.BalanceOf(&bind.CallOpts{}, r.ZevmTestDAppAddr)
	if err != nil {
		panic(err)
	}
	previousBalanceEVM, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, r.EVMAuth.From)
	if err != nil {
		panic(err)
	}

	// Call the SendHelloWorld function on the EVM dapp Contract which would in turn create a new send, to be picked up by the zeta-clients
	// set Do revert to false which adds a message to signal the ZEVM zetaReceiver to not revert the transaction
	tx, err = testDAppEVM.SendHelloWorld(r.EVMAuth, destinationAddress, zEVMChainID, amount, false)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("TestDApp.SendHello tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)

	// New inbound message picked up by zeta-clients and voted on by observers to initiate a contract call on zEVM
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx.CctxStatus.Status != cctxtypes.CctxStatus_OutboundMined {
		panic("expected cctx to be outbound_mined")
	}
	r.Logger.Info(fmt.Sprintf("🔄 Cctx mined for contract call chain zevm %s", cctx.Index))

	// On finalization the Fungible module calls the onReceive function which in turn calls the onZetaMessage function on the destination contract
	receipt, err = r.ZEVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(cctx.GetCurrentOutTxParam().OutboundTxHash))
	if err != nil {
		panic(err)
	}
	if receipt.Status != 1 {
		panic("tx failed")
	}

	testDAppZEVM, err := testdapp.NewTestDApp(r.ZevmTestDAppAddr, r.ZEVMClient)
	if err != nil {
		panic(err)
	}

	// Check event emitted
	receivedHelloWorldEvent := false
	for _, log := range receipt.Logs {
		_, err := testDAppZEVM.ParseHelloWorldEvent(*log)
		if err == nil {
			r.Logger.Info("Received HelloWorld event")
			receivedHelloWorldEvent = true
		}
	}
	if !receivedHelloWorldEvent {
		panic(fmt.Sprintf("expected HelloWorld event, logs: %+v", receipt.Logs))
	}

	// Check ZETA balance on ZEVM TestDApp and check new balance is previous balance + amount
	newBalanceZEVM, err := r.WZeta.BalanceOf(&bind.CallOpts{}, r.ZevmTestDAppAddr)
	if err != nil {
		panic(err)
	}
	if newBalanceZEVM.Cmp(big.NewInt(0).Add(previousBalanceZEVM, amount)) != 0 {
		panic(fmt.Sprintf(
			"expected new balance to be %s, got %s",
			big.NewInt(0).Add(previousBalanceZEVM, amount).String(),
			newBalanceZEVM.String()),
		)
	}

	// Check ZETA balance on EVM TestDApp and check new balance is previous balance - amount
	newBalanceEVM, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, r.EVMAuth.From)
	if err != nil {
		panic(err)
	}
	if newBalanceEVM.Cmp(big.NewInt(0).Sub(previousBalanceEVM, amount)) != 0 {
		panic(fmt.Sprintf(
			"expected new balance to be %s, got %s",
			big.NewInt(0).Sub(previousBalanceEVM, amount).String(),
			newBalanceEVM.String()),
		)
	}
}
