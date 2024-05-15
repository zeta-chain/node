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

func TestMessagePassingEVMtoZEVMRevertFail(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestMessagePassingEVMtoZEVMRevertFail requires exactly one argument for the amount.")
	}

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestMessagePassingEVMtoZEVMRevertFail.")
	}

	// Deploying a test contract not containing a logic for reverting the cctx
	testDappNoRevertEVMAddr, tx, testDappNoRevertEVM, err := testdappnorevert.DeployTestDAppNoRevert(
		r.EVMAuth,
		r.EVMClient,
		r.ConnectorEthAddr,
		r.ZetaEthAddr,
	)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("TestDAppNoRevertEVM deployed at: %s", testDappNoRevertEVMAddr.Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "deploy TestDAppNoRevert")
	if receipt.Status == 0 {
		panic("deploy TestDAppNoRevert failed")
	}

	// Set destination details
	zEVMChainID, err := r.ZEVMClient.ChainID(r.Ctx)
	if err != nil {
		panic(err)
	}

	destinationAddress := r.ZevmTestDAppAddr

	// Contract call originates from EVM chain
	tx, err = r.ZetaEth.Approve(r.EVMAuth, testDappNoRevertEVMAddr, amount)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("tx failed")
	}
	r.Logger.Info("Approve tx receipt: %d", receipt.Status)

	// Get ZETA balance before test
	previousBalanceZEVM, err := r.WZeta.BalanceOf(&bind.CallOpts{}, r.ZevmTestDAppAddr)
	if err != nil {
		panic(err)
	}
	previousBalanceEVM, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, testDappNoRevertEVMAddr)
	if err != nil {
		panic(err)
	}

	// Send message with doRevert
	tx, err = testDappNoRevertEVM.SendHelloWorld(r.EVMAuth, destinationAddress, zEVMChainID, amount, true)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("TestDAppNoRevert.SendHello tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)

	// New inbound message picked up by zeta-clients and voted on by observers to initiate a contract call on zEVM which would revert the transaction
	// A revert transaction is created and gets fialized on the original sender chain.
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx.CctxStatus.Status != cctxtypes.CctxStatus_Aborted {
		panic("expected cctx to be reverted")
	}

	// Check ZETA balance on ZEVM TestDApp and check new balance is previous balance
	newBalanceZEVM, err := r.WZeta.BalanceOf(&bind.CallOpts{}, r.ZevmTestDAppAddr)
	if err != nil {
		panic(err)
	}
	if newBalanceZEVM.Cmp(previousBalanceZEVM) != 0 {
		panic(fmt.Sprintf(
			"expected new balance to be %s, got %s",
			previousBalanceZEVM.String(),
			newBalanceZEVM.String()),
		)
	}

	// Check ZETA balance on EVM TestDApp and check new balance is previous balance
	newBalanceEVM, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, testDappNoRevertEVMAddr)
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
}
