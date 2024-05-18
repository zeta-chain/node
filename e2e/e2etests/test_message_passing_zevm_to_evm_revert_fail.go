package e2etests

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/zeta-chain/zetacore/e2e/contracts/testdappnorevert"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	cctxtypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMessagePassingZEVMtoEVMRevertFail(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestMessagePassingZEVMtoEVMRevertFail requires exactly one argument for the amount.")
	}

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestMessagePassingZEVMtoEVMRevertFail.")
	}

	// Deploying a test contract not containing a logic for reverting the cctx
	testDappNoRevertAddr, tx, testDappNoRevert, err := testdappnorevert.DeployTestDAppNoRevert(
		r.ZEVMAuth,
		r.ZEVMClient,
		r.ConnectorZEVMAddr,
		r.WZetaAddr,
	)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("TestDAppNoRevert deployed at: %s", testDappNoRevertAddr.Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "deploy TestDAppNoRevert")
	if receipt.Status == 0 {
		panic("deploy TestDAppNoRevert failed")
	}

	// Set destination details
	EVMChainID, err := r.EVMClient.ChainID(r.Ctx)
	if err != nil {
		panic(err)
	}
	destinationAddress := r.EvmTestDAppAddr

	// Contract call originates from ZEVM chain
	r.ZEVMAuth.Value = amount
	tx, err = r.WZeta.Deposit(r.ZEVMAuth)
	if err != nil {
		panic(err)
	}

	r.ZEVMAuth.Value = big.NewInt(0)
	r.Logger.Info("wzeta deposit tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "wzeta deposit")
	if receipt.Status == 0 {
		panic("deposit failed")
	}

	tx, err = r.WZeta.Approve(r.ZEVMAuth, testDappNoRevertAddr, amount)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("wzeta approve tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "wzeta approve")
	if receipt.Status == 0 {
		panic(fmt.Sprintf("approve failed, logs: %+v", receipt.Logs))
	}

	// Get previous balances to check funds are not minted anywhere when aborted
	previousBalanceEVM, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, r.EvmTestDAppAddr)
	if err != nil {
		panic(err)
	}
	previousBalanceZEVM, err := r.WZeta.BalanceOf(&bind.CallOpts{}, testDappNoRevertAddr)
	if err != nil {
		panic(err)
	}

	// Send message with doRevert
	tx, err = testDappNoRevert.SendHelloWorld(r.ZEVMAuth, destinationAddress, EVMChainID, amount, true)
	if err != nil {
		panic(err)
	}

	r.Logger.Info("TestDAppNoRevert.SendHello tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic(fmt.Sprintf("send failed, logs: %+v", receipt.Logs))
	}

	// The revert tx will fail, the cctx state should be aborted
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx.CctxStatus.Status != cctxtypes.CctxStatus_Aborted {
		panic("expected cctx to be reverted")
	}

	// Check ZETA balance on EVM TestDApp and check new balance is previous balance
	newBalanceEVM, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, r.EvmTestDAppAddr)
	if err != nil {
		panic(err)
	}
	if newBalanceEVM.Cmp(previousBalanceEVM) != 0 {
		panic(fmt.Sprintf(
			"expected new balance to be %s, got %s",
			previousBalanceEVM.String(),
			newBalanceEVM.String()),
		)
	}

	// Check the funds are not minted to the contract as the cctx has been aborted
	newBalanceZEVM, err := r.WZeta.BalanceOf(&bind.CallOpts{}, testDappNoRevertAddr)
	if err != nil {
		panic(err)
	}
	if newBalanceZEVM.Cmp(previousBalanceZEVM) != 0 {
		panic(fmt.Sprintf("expected new balance to be %s, got %s", previousBalanceZEVM.String(), newBalanceZEVM.String()))
	}
}
