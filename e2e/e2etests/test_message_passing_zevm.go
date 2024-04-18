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

	r.Logger.Print(fmt.Sprintf("🔄 Successful tx intx : %s", receipt.TxHash.String()))
	// New inbound message picked up by zeta-clients and voted on by observers to initiate a contract call on zEVM
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx.CctxStatus.Status != cctxtypes.CctxStatus_OutboundMined {
		panic("expected cctx to be outbound_mined")
	}
	r.Logger.Print(fmt.Sprintf("🔄 Cctx mined for contract call chain zevm %s", cctx.Index))

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
	receivedHelloWorldEvent := false
	for _, log := range receipt.Logs {
		_, err := testDAppZEVM.ParseHelloWorldEvent(*log)
		if err == nil {
			r.Logger.Info("Received HelloWorld event")
			receivedHelloWorldEvent = true
		}
	}
	if !receivedHelloWorldEvent {
		panic("expected HelloWorld event")
	}
}

func TestMessagePassingZEVMRevert(r *runner.E2ERunner, args []string) {

	if len(args) != 1 {
		panic("TestMessagePassingRevert requires exactly one argument for the amount.")
	}

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestMessagePassingRevert.")
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

	tx, err = testDAppEVM.SendHelloWorld(auth, destinationAddress, zEVMChainID, amount, true)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("TestDApp.SendHello tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)

	r.Logger.Print(fmt.Sprintf("🔄 Revert tx intx : %s", receipt.TxHash.String()))

	// New inbound message picked up by zeta-clients and voted on by observers to initiate a contract call on zEVM
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx.CctxStatus.Status != cctxtypes.CctxStatus_Reverted {
		panic("expected cctx to be reverted")
	}
	r.Logger.Print(fmt.Sprintf("🔄 Cctx mined for revert contract call chain zevm %s", cctx.Index))

	receipt, err = r.EVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(cctx.GetCurrentOutTxParam().OutboundTxHash))
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
		panic("expected HelloWorld event")
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

	EVMChainID, err := r.EVMClient.ChainID(r.Ctx)
	if err != nil {
		panic(err)
	}
	destinationAddress := r.EvmTestDAppAddr

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
	tx, err = testDAppZEVM.SendHelloWorld(r.ZEVMAuth, destinationAddress, EVMChainID, amount, false)
	if err != nil {
		panic(err)
	}

	r.Logger.Info("TestDApp.SendHello tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	if receipt.Status == 0 {
		panic(fmt.Sprintf("send failed, logs: %+v", receipt.Logs))
	}

	r.Logger.Print(fmt.Sprintf("🔄 Successful tx intx : %s", receipt.TxHash.String()))
	// New inbound message picked up by zeta-clients and voted on by observers to initiate a contract call on zEVM
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx.CctxStatus.Status != cctxtypes.CctxStatus_OutboundMined {
		panic("expected cctx to be outbound_mined")
	}
	r.Logger.Print(fmt.Sprintf("🔄 Cctx mined for contract call chain zevm %s", cctx.Index))

	receipt, err = r.EVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(cctx.GetCurrentOutTxParam().OutboundTxHash))
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
		panic("expected HelloWorld event")
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

	EVMChainID, err := r.EVMClient.ChainID(r.Ctx)
	if err != nil {
		panic(err)
	}
	destinationAddress := r.EvmTestDAppAddr

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
	tx, err = testDAppZEVM.SendHelloWorld(r.ZEVMAuth, destinationAddress, EVMChainID, amount, true)
	if err != nil {
		panic(err)
	}
	r.Logger.Print("TestDApp ZEVM address: %s", r.ZevmTestDAppAddr.String())

	r.Logger.Info("TestDApp.SendHello tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	if receipt.Status == 0 {
		panic(fmt.Sprintf("send failed, logs: %+v", receipt.Logs))
	}

	r.Logger.Print(fmt.Sprintf("🔄 Successful tx intx : %s", receipt.TxHash.String()))
	// New inbound message picked up by zeta-clients and voted on by observers to initiate a contract call on zEVM
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx.CctxStatus.Status != cctxtypes.CctxStatus_Reverted {
		panic("expected cctx to be reverted")
	}
	r.Logger.Print(fmt.Sprintf("🔄 Cctx mined for contract call chain zevm %s", cctx.Index))

	receipt, err = r.ZEVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(cctx.GetCurrentOutTxParam().OutboundTxHash))
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
		panic("expected HelloWorld event")
	}
}

// bgGCGUux3roBhJr9PgNaC3DOfLBp5ILuZjUZx2z1abQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQ==
// bgGCGUux3roBhJr9PgNaC3DOfLBp5ILuZjUZx2z1abQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQ==
