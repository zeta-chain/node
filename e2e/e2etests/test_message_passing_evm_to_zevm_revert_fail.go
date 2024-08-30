package e2etests

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/contracts/testdappnorevert"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	cctxtypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestMessagePassingEVMtoZEVMRevertFail(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse the amount
	amount := parseBigInt(r, args[0])

	// Deploying a test contract not containing a logic for reverting the cctx
	testDappNoRevertEVMAddr, tx, testDappNoRevertEVM, err := testdappnorevert.DeployTestDAppNoRevert(
		r.EVMAuth,
		r.EVMClient,
		r.ConnectorEthAddr,
		r.ZetaEthAddr,
	)
	require.NoError(r, err)

	r.Logger.Info("TestDAppNoRevertEVM deployed at: %s", testDappNoRevertEVMAddr.Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	r.Logger.EVMReceipt(*receipt, "deploy TestDAppNoRevert")

	// Set destination details
	zEVMChainID, err := r.ZEVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	destinationAddress := r.ZevmTestDAppAddr

	// Contract call originates from EVM chain
	tx, err = r.ZetaEth.Approve(r.EVMAuth, testDappNoRevertEVMAddr, amount)
	require.NoError(r, err)

	r.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	r.Logger.Info("Approve tx receipt: %d", receipt.Status)

	// Get ZETA balance before test
	previousBalanceZEVM, err := r.WZeta.BalanceOf(&bind.CallOpts{}, r.ZevmTestDAppAddr)
	require.NoError(r, err)

	previousBalanceEVM, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, testDappNoRevertEVMAddr)
	require.NoError(r, err)

	// Send message with doRevert
	tx, err = testDappNoRevertEVM.SendHelloWorld(r.EVMAuth, destinationAddress, zEVMChainID, amount, true)
	require.NoError(r, err)

	r.Logger.Info("TestDAppNoRevert.SendHello tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)

	// New inbound message picked up by zeta-clients and voted on by observers to initiate a contract call on zEVM which would revert the transaction
	// A revert transaction is created and gets finalized on the original sender chain.
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, cctxtypes.CctxStatus_Aborted)

	// Check ZETA balance on ZEVM TestDApp and check new balance is previous balance
	newBalanceZEVM, err := r.WZeta.BalanceOf(&bind.CallOpts{}, r.ZevmTestDAppAddr)
	require.NoError(r, err)
	require.Equal(
		r,
		0,
		newBalanceZEVM.Cmp(previousBalanceZEVM),
		"expected new balance to be %s, got %s",
		previousBalanceZEVM.String(),
		newBalanceZEVM.String(),
	)

	// Check ZETA balance on EVM TestDApp and check new balance is previous balance
	newBalanceEVM, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, testDappNoRevertEVMAddr)
	require.NoError(r, err)
	require.Equal(r, 0, newBalanceEVM.Cmp(previousBalanceEVM))
}
