package e2etests

import (
	"fmt"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/e2e/contracts/testdapp"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	cctxtypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMessagePassingEVMtoZEVM(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestMessagePassing requires exactly one argument for the amount.")
	}
	amount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestMessagePassing.")
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
	// Call the SendHelloWorld function on the EVM dapp Contract which would in turn create a new send, to be picked up by the zeta-clients
	// set Do revert to false which adds a message to signal the ZEVM zetaReceiver to not revert the transaction
	tx, err = testDAppEVM.SendHelloWorld(r.EVMAuth, destinationAddress, zEVMChainID, amount, false)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("TestDApp.SendHello tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)

	// New inbound message picked up by zeta-clients and voted on by observers to initiate a contract call on zEVM
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx.CctxStatus.Status != cctxtypes.CctxStatus_OutboundMined {
		panic("expected cctx to be outbound_mined")
	}
	r.Logger.Info(fmt.Sprintf("🔄 Cctx mined for contract call chain zevm %s", cctx.Index))

	// On finalization the Fungible module calls the onReceive function which in turn calls the onZetaMessage function on the destination contract
	receipt, err = r.ZEVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(cctx.GetCurrentOutboundParam().Hash))
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
}

func TestMessagePassingEVMtoZEVMRevert(r *runner.E2ERunner, args []string) {

	if len(args) != 1 {
		panic("TestMessagePassingRevert requires exactly one argument for the amount.")
	}

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestMessagePassingRevert.")
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
		panic("tx failed")
	}
	r.Logger.Info("Approve tx receipt: %d", receipt.Status)
	testDAppEVM, err := testdapp.NewTestDApp(r.EvmTestDAppAddr, r.EVMClient)
	if err != nil {
		panic(err)
	}
	// Call the SendHelloWorld function on the EVM dapp Contract which would in turn create a new send, to be picked up by the zeta-clients
	// set Do revert to true which adds a message to signal the ZEVM zetaReceiver to revert the transaction
	tx, err = testDAppEVM.SendHelloWorld(r.EVMAuth, destinationAddress, zEVMChainID, amount, true)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("TestDApp.SendHello tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)

	// New inbound message picked up by zeta-clients and voted on by observers to initiate a contract call on zEVM which would revert the transaction
	// A revert transaction is created and gets fialized on the original sender chain.
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx.CctxStatus.Status != cctxtypes.CctxStatus_Reverted {
		panic("expected cctx to be reverted")
	}

	// On finalization the Tss address calls the onRevert function which in turn calls the onZetaRevert function on the sender contract
	receipt, err = r.EVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(cctx.GetCurrentOutboundParam().Hash))
	if err != nil {
		panic(err)
	}
	if receipt.Status != 1 {
		panic("tx failed")
	}
	receivedHelloWorldEvent := false
	for _, log := range receipt.Logs {
		_, err := testDAppEVM.ParseRevertedHelloWorldEvent(*log)
		if err == nil {
			r.Logger.Info("Received RevertHelloWorld event:")
			receivedHelloWorldEvent = true
		}
	}
	if !receivedHelloWorldEvent {
		panic(fmt.Sprintf("expected Reverted HelloWorld event, logs: %+v", receipt.Logs))
	}
}

func TestMessagePassingZEVMtoEVM(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestMessagePassing requires exactly one argument for the amount.")
	}

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestMessagePassing.")
	}
	// Set destination details
	EVMChainID, err := r.EVMClient.ChainID(r.Ctx)
	if err != nil {
		panic(err)
	}
	destinationAddress := r.EvmTestDAppAddr

	// Contract call originates from ZEVM chain
	r.ZEVMAuth.Value = amount
	tx, err := r.WZeta.Deposit(r.ZEVMAuth)
	if err != nil {
		panic(err)
	}
	r.ZEVMAuth.Value = big.NewInt(0)
	r.Logger.Info("wzeta deposit tx hash: %s", tx.Hash().Hex())
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "wzeta deposit")
	if receipt.Status == 0 {
		panic("deposit failed")
	}
	tx, err = r.WZeta.Approve(r.ZEVMAuth, r.ZevmTestDAppAddr, amount)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("wzeta approve tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "wzeta approve")
	if receipt.Status == 0 {
		panic(fmt.Sprintf("approve failed, logs: %+v", receipt.Logs))
	}
	testDAppZEVM, err := testdapp.NewTestDApp(r.ZevmTestDAppAddr, r.ZEVMClient)
	if err != nil {
		panic(err)
	}
	// Call the SendHelloWorld function on the ZEVM dapp Contract which would in turn create a new send, to be picked up by the zetanode evm hooks
	// set Do revert to false which adds a message to signal the EVM zetaReceiver to not revert the transaction
	tx, err = testDAppZEVM.SendHelloWorld(r.ZEVMAuth, destinationAddress, EVMChainID, amount, false)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("TestDApp.SendHello tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic(fmt.Sprintf("send failed, logs: %+v", receipt.Logs))
	}

	// Transaction is picked up by the zetanode evm hooks and a new contract call is initiated on the EVM chain
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx.CctxStatus.Status != cctxtypes.CctxStatus_OutboundMined {
		panic("expected cctx to be outbound_mined")
	}

	// On finalization the Tss calls the onReceive function which in turn calls the onZetaMessage function on the destination contract.
	receipt, err = r.EVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(cctx.GetCurrentOutboundParam().Hash))
	if err != nil {
		panic(err)
	}
	if receipt.Status != 1 {
		panic("tx failed")
	}
	testDAppEVM, err := testdapp.NewTestDApp(r.EvmTestDAppAddr, r.EVMClient)
	if err != nil {
		panic(err)
	}
	receivedHelloWorldEvent := false
	for _, log := range receipt.Logs {
		_, err := testDAppEVM.ParseHelloWorldEvent(*log)
		if err == nil {
			r.Logger.Info("Received HelloWorld event:")
			receivedHelloWorldEvent = true
		}
	}
	if !receivedHelloWorldEvent {
		panic(fmt.Sprintf("expected HelloWorld event, logs: %+v", receipt.Logs))
	}
}

func TestMessagePassingZEVMtoEVMRevert(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestMessagePassing requires exactly one argument for the amount.")
	}

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestMessagePassing.")
	}

	// Set destination details
	EVMChainID, err := r.EVMClient.ChainID(r.Ctx)
	if err != nil {
		panic(err)
	}
	destinationAddress := r.EvmTestDAppAddr

	// Contract call originates from ZEVM chain
	r.ZEVMAuth.Value = amount
	tx, err := r.WZeta.Deposit(r.ZEVMAuth)
	if err != nil {
		panic(err)
	}
	r.ZEVMAuth.Value = big.NewInt(0)
	r.Logger.Info("wzeta deposit tx hash: %s", tx.Hash().Hex())
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "wzeta deposit")
	if receipt.Status == 0 {
		panic("deposit failed")
	}
	tx, err = r.WZeta.Approve(r.ZEVMAuth, r.ZevmTestDAppAddr, amount)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("wzeta approve tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "wzeta approve")
	if receipt.Status == 0 {
		panic(fmt.Sprintf("approve failed, logs: %+v", receipt.Logs))
	}
	testDAppZEVM, err := testdapp.NewTestDApp(r.ZevmTestDAppAddr, r.ZEVMClient)
	if err != nil {
		panic(err)
	}
	// Call the SendHelloWorld function on the ZEVM dapp Contract which would in turn create a new send, to be picked up by the zetanode evm hooks
	// set Do revert to true which adds a message to signal the EVM zetaReceiver to revert the transaction
	tx, err = testDAppZEVM.SendHelloWorld(r.ZEVMAuth, destinationAddress, EVMChainID, amount, true)
	if err != nil {
		panic(err)
	}

	r.Logger.Info("TestDApp.SendHello tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic(fmt.Sprintf("send failed, logs: %+v", receipt.Logs))
	}

	// New inbound message picked up by zetanode evm hooks and processed directly to initiate a contract call on EVM which would revert the transaction
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx.CctxStatus.Status != cctxtypes.CctxStatus_Reverted {
		panic("expected cctx to be reverted")
	}

	// On finalization the Fungible module calls the onRevert function which in turn calls the onZetaRevert function on the sender contract
	receipt, err = r.ZEVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(cctx.GetCurrentOutboundParam().Hash))
	if err != nil {
		panic(err)
	}
	if receipt.Status != 1 {
		panic("tx failed")
	}
	receivedHelloWorldEvent := false
	for _, log := range receipt.Logs {
		_, err := testDAppZEVM.ParseRevertedHelloWorldEvent(*log)
		if err == nil {
			r.Logger.Info("Received HelloWorld event:")
			receivedHelloWorldEvent = true
		}
	}
	if !receivedHelloWorldEvent {
		panic(fmt.Sprintf("expected Reverted HelloWorld event, logs: %+v", receipt.Logs))
	}
}
