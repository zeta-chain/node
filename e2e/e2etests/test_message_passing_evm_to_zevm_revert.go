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

func TestMessagePassingEVMtoZEVMRevert(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestMessagePassingEVMtoZEVMRevert requires exactly one argument for the amount.")
	}

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestMessagePassingEVMtoZEVMRevert.")
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

	// Get ZETA balance before test
	previousBalanceZEVM, err := r.WZeta.BalanceOf(&bind.CallOpts{}, r.ZevmTestDAppAddr)
	if err != nil {
		panic(err)
	}
	previousBalanceEVM, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, r.EvmTestDAppAddr)
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
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx.CctxStatus.Status != cctxtypes.CctxStatus_Reverted {
		panic("expected cctx to be reverted")
	}

	// On finalization the Tss address calls the onRevert function which in turn calls the onZetaRevert function on the sender contract
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
		panic(fmt.Sprintf("expected Reverted HelloWorld event, logs: %+v", receipt.Logs))
	}

	// Check ZETA balance on ZEVM TestDApp and check new balance is previous balance
	newBalanceZEVM, err := r.WZeta.BalanceOf(&bind.CallOpts{}, r.ZevmTestDAppAddr)
	if err != nil {
		panic(err)
	}
	if newBalanceZEVM.Cmp(previousBalanceZEVM) != 0 {
		panic(fmt.Sprintf("expected new balance to be %s, got %s", previousBalanceZEVM.String(), newBalanceZEVM.String()))
	}

	// Check ZETA balance on EVM TestDApp and check new balance is between previous balance and previous balance + amount
	// New balance is increased because ZETA are sent from the sender but sent back to the contract
	// New balance is less than previous balance + amount because of the gas fee to pay
	newBalanceEVM, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, r.EvmTestDAppAddr)
	if err != nil {
		panic(err)
	}
	previousBalanceAndAmountEVM := big.NewInt(0).Add(previousBalanceEVM, amount)

	// check higher than previous balance and lower than previous balance + amount
	if newBalanceEVM.Cmp(previousBalanceEVM) <= 0 || newBalanceEVM.Cmp(previousBalanceAndAmountEVM) > 0 {
		panic(fmt.Sprintf(
			"expected new balance to be between %s and %s, got %s",
			previousBalanceEVM.String(),
			previousBalanceAndAmountEVM.String(),
			newBalanceEVM.String()),
		)
	}
}
