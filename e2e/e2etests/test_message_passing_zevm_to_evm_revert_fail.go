package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/contracts/testdappnorevert"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	cctxtypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestMessagePassingZEVMtoEVMRevertFail(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse the amount
	amount := parseBigInt(r, args[0])

	// Deploying a test contract not containing a logic for reverting the cctx
	testDappNoRevertAddr, tx, testDappNoRevert, err := testdappnorevert.DeployTestDAppNoRevert(
		r.ZEVMAuth,
		r.ZEVMClient,
		r.ConnectorZEVMAddr,
		r.WZetaAddr,
	)
	require.NoError(r, err)

	r.Logger.Info("TestDAppNoRevert deployed at: %s", testDappNoRevertAddr.Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "deploy TestDAppNoRevert")
	utils.RequireTxSuccessful(r, receipt)

	// Set destination details
	EVMChainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	destinationAddress := r.EvmTestDAppAddr

	// Contract call originates from ZEVM chain
	r.ZEVMAuth.Value = amount
	tx, err = r.WZeta.Deposit(r.ZEVMAuth)
	require.NoError(r, err)

	r.ZEVMAuth.Value = big.NewInt(0)
	r.Logger.Info("wzeta deposit tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "wzeta deposit")
	utils.RequireTxSuccessful(r, receipt)

	tx, err = r.WZeta.Approve(r.ZEVMAuth, testDappNoRevertAddr, amount)
	require.NoError(r, err)

	r.Logger.Info("wzeta approve tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "wzeta approve")
	utils.RequireTxSuccessful(r, receipt)

	// Get previous balances to check funds are not minted anywhere when aborted
	previousBalanceZEVM, err := r.WZeta.BalanceOf(&bind.CallOpts{}, testDappNoRevertAddr)
	require.NoError(r, err)

	// Send message with doRevert
	tx, err = testDappNoRevert.SendHelloWorld(r.ZEVMAuth, destinationAddress, EVMChainID, amount, true)
	require.NoError(r, err)

	r.Logger.Info("TestDAppNoRevert.SendHello tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// The revert tx will fail, the cctx state should be aborted
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, cctxtypes.CctxStatus_Aborted)

	// Check the funds are not minted to the contract as the cctx has been aborted
	newBalanceZEVM, err := r.WZeta.BalanceOf(&bind.CallOpts{}, testDappNoRevertAddr)
	require.NoError(r, err)
	require.Equal(r,
		0,
		newBalanceZEVM.Cmp(previousBalanceZEVM),
		"expected new balance to be %s, got %s",
		previousBalanceZEVM.String(),
		newBalanceZEVM.String(),
	)
}
