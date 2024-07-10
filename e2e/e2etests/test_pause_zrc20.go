package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/contracts/vault"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestPauseZRC20(r *runner.E2ERunner, _ []string) {
	// Setup vault used to test zrc20 interactions
	r.Logger.Info("Deploying vault")
	vaultAddr, _, vaultContract, err := vault.DeployVault(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)

	// Approving vault to spend ZRC20
	tx, err := r.ETHZRC20.Approve(r.ZEVMAuth, vaultAddr, big.NewInt(1e18))
	require.NoError(r, err)

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	tx, err = r.ERC20ZRC20.Approve(r.ZEVMAuth, vaultAddr, big.NewInt(1e18))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// Pause ETH ZRC20
	r.Logger.Info("Pausing ETH")
	msgPause := fungibletypes.NewMsgPauseZRC20(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.EmergencyPolicyName),
		[]string{r.ETHZRC20Addr.Hex()},
	)
	res, err := r.ZetaTxServer.BroadcastTx(utils.EmergencyPolicyName, msgPause)
	require.NoError(r, err)
	r.Logger.Info("pause zrc20 tx hash: %s", res.TxHash)

	// Fetch and check pause status
	fcRes, err := r.FungibleClient.ForeignCoins(r.Ctx, &fungibletypes.QueryGetForeignCoinsRequest{
		Index: r.ETHZRC20Addr.Hex(),
	})
	require.NoError(r, err)
	require.True(r, fcRes.GetForeignCoins().Paused, "ETH should be paused")

	r.Logger.Info("ETH is paused")

	// Try operations with ETH ZRC20
	r.Logger.Info("Can no longer do operations on ETH ZRC20")

	tx, err = r.ETHZRC20.Transfer(r.ZEVMAuth, sample.EthAddress(), big.NewInt(1e5))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt)

	tx, err = r.ETHZRC20.Burn(r.ZEVMAuth, big.NewInt(1e5))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt)

	// Operation on a contract that interact with ETH ZRC20 should fail
	r.Logger.Info("Vault contract can no longer interact with ETH ZRC20: %s", r.ETHZRC20Addr.Hex())
	tx, err = vaultContract.Deposit(r.ZEVMAuth, r.ETHZRC20Addr, big.NewInt(1e5))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt)

	r.Logger.Info("Operations all failed")

	// Check we can still interact with ERC20 ZRC20
	r.Logger.Info("Check other ZRC20 can still be operated")

	tx, err = r.ERC20ZRC20.Transfer(r.ZEVMAuth, sample.EthAddress(), big.NewInt(1e3))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	tx, err = vaultContract.Deposit(r.ZEVMAuth, r.ERC20ZRC20Addr, big.NewInt(1e3))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// Check deposit revert when paused
	signedTx, err := r.SendEther(r.TSSAddress, big.NewInt(1e17), nil)
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, signedTx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, signedTx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, types.CctxStatus_Reverted)

	r.Logger.Info("CCTX has been reverted")

	// Unpause ETH ZRC20
	r.Logger.Info("Unpausing ETH")
	msgUnpause := fungibletypes.NewMsgUnpauseZRC20(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		[]string{r.ETHZRC20Addr.Hex()},
	)
	res, err = r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, msgUnpause)
	require.NoError(r, err)

	r.Logger.Info("unpause zrc20 tx hash: %s", res.TxHash)

	// Fetch and check pause status
	fcRes, err = r.FungibleClient.ForeignCoins(r.Ctx, &fungibletypes.QueryGetForeignCoinsRequest{
		Index: r.ETHZRC20Addr.Hex(),
	})
	require.NoError(r, err)
	require.False(r, fcRes.GetForeignCoins().Paused, "ETH should be unpaused")

	r.Logger.Info("ETH is unpaused")

	// Try operations with ETH ZRC20
	r.Logger.Info("Can do operations on ETH ZRC20 again")

	tx, err = r.ETHZRC20.Transfer(r.ZEVMAuth, sample.EthAddress(), big.NewInt(1e5))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	tx, err = r.ETHZRC20.Burn(r.ZEVMAuth, big.NewInt(1e5))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// Can deposit tokens into the vault again
	tx, err = vaultContract.Deposit(r.ZEVMAuth, r.ETHZRC20Addr, big.NewInt(1e5))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	r.Logger.Info("Operations all succeeded")
}
