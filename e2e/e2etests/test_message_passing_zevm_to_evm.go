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

func TestMessagePassingZEVMtoEVM(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestMessagePassingZEVMtoEVM requires exactly one argument for the amount.")
	}

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestMessagePassingZEVMtoEVM.")
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

	// Get previous balances
	previousBalanceEVM, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, r.EvmTestDAppAddr)
	if err != nil {
		panic(err)
	}
	previousBalanceZEVM, err := r.WZeta.BalanceOf(&bind.CallOpts{}, r.ZEVMAuth.From)
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
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx.CctxStatus.Status != cctxtypes.CctxStatus_OutboundMined {
		panic("expected cctx to be outbound_mined")
	}

	// On finalization the Tss calls the onReceive function which in turn calls the onZetaMessage function on the destination contract.
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
		panic(fmt.Sprintf("expected HelloWorld event, logs: %+v", receipt.Logs))
	}

	// Check ZETA balance on EVM TestDApp and check new balance between previous balance and previous balance + amount
	// Contract receive less than the amount because of the gas fee to pay
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

	// Check ZETA balance on ZEVM TestDApp and check new balance is previous balance - amount
	newBalanceZEVM, err := r.WZeta.BalanceOf(&bind.CallOpts{}, r.ZEVMAuth.From)
	if err != nil {
		panic(err)
	}
	if newBalanceZEVM.Cmp(big.NewInt(0).Sub(previousBalanceZEVM, amount)) != 0 {
		panic(fmt.Sprintf(
			"expected new balance to be %s, got %s",
			big.NewInt(0).Sub(previousBalanceZEVM, amount).String(),
			newBalanceZEVM.String()),
		)
	}
}
