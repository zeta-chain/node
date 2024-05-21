package e2etests

import (
	"fmt"
	"math/big"

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
	if err != nil {
		panic(err)
	}
	// Approving vault to spend ZRC20
	tx, err := r.ETHZRC20.Approve(r.ZEVMAuth, vaultAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("Vault approval should succeed")
	}
	tx, err = r.ERC20ZRC20.Approve(r.ZEVMAuth, vaultAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("Vault approval should succeed")
	}

	// Pause ETH ZRC20
	r.Logger.Info("Pausing ETH")
	msg := fungibletypes.NewMsgUpdateZRC20PausedStatus(
		r.ZetaTxServer.GetAccountAddress(0),
		[]string{r.ETHZRC20Addr.Hex()},
		fungibletypes.UpdatePausedStatusAction_PAUSE,
	)
	res, err := r.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msg)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("pause zrc20 tx hash: %s", res.TxHash)

	// Fetch and check pause status
	fcRes, err := r.FungibleClient.ForeignCoins(r.Ctx, &fungibletypes.QueryGetForeignCoinsRequest{
		Index: r.ETHZRC20Addr.Hex(),
	})
	if err != nil {
		panic(err)
	}
	if !fcRes.GetForeignCoins().Paused {
		panic("ETH should be paused")
	}
	r.Logger.Info("ETH is paused")

	// Try operations with ETH ZRC20
	r.Logger.Info("Can no longer do operations on ETH ZRC20")
	tx, err = r.ETHZRC20.Transfer(r.ZEVMAuth, sample.EthAddress(), big.NewInt(1e5))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 1 {
		panic("transfer should fail")
	}
	tx, err = r.ETHZRC20.Burn(r.ZEVMAuth, big.NewInt(1e5))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 1 {
		panic("burn should fail")
	}

	// Operation on a contract that interact with ETH ZRC20 should fail
	r.Logger.Info("Vault contract can no longer interact with ETH ZRC20: %s", r.ETHZRC20Addr.Hex())
	tx, err = vaultContract.Deposit(r.ZEVMAuth, r.ETHZRC20Addr, big.NewInt(1e5))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 1 {
		panic("deposit should fail")
	}
	r.Logger.Info("Operations all failed")

	// Check we can still interact with ERC20 ZRC20
	r.Logger.Info("Check other ZRC20 can still be operated")
	tx, err = r.ERC20ZRC20.Transfer(r.ZEVMAuth, sample.EthAddress(), big.NewInt(1e3))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("ERC20 ZRC20 transfer should succeed")
	}
	tx, err = vaultContract.Deposit(r.ZEVMAuth, r.ERC20ZRC20Addr, big.NewInt(1e3))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("ERC20 ZRC20 vault deposit should succeed")
	}

	// Check deposit revert when paused
	signedTx, err := r.SendEther(r.TSSAddress, big.NewInt(1e17), nil)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, signedTx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("deposit eth tx failed")
	}
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, signedTx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx.CctxStatus.Status != types.CctxStatus_Reverted {
		panic(fmt.Sprintf("expected cctx status to be Reverted; got %s", cctx.CctxStatus.Status))
	}
	r.Logger.Info("CCTX has been reverted")

	// Unpause ETH ZRC20
	r.Logger.Info("Unpausing ETH")
	msg = fungibletypes.NewMsgUpdateZRC20PausedStatus(
		r.ZetaTxServer.GetAccountAddress(0),
		[]string{r.ETHZRC20Addr.Hex()},
		fungibletypes.UpdatePausedStatusAction_UNPAUSE,
	)
	res, err = r.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msg)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("unpause zrc20 tx hash: %s", res.TxHash)

	// Fetch and check pause status
	fcRes, err = r.FungibleClient.ForeignCoins(r.Ctx, &fungibletypes.QueryGetForeignCoinsRequest{
		Index: r.ETHZRC20Addr.Hex(),
	})
	if err != nil {
		panic(err)
	}
	if fcRes.GetForeignCoins().Paused {
		panic("ETH should be unpaused")
	}
	r.Logger.Info("ETH is unpaused")

	// Try operations with ETH ZRC20
	r.Logger.Info("Can do operations on ETH ZRC20 again")
	tx, err = r.ETHZRC20.Transfer(r.ZEVMAuth, sample.EthAddress(), big.NewInt(1e5))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("transfer should succeed")
	}
	tx, err = r.ETHZRC20.Burn(r.ZEVMAuth, big.NewInt(1e5))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("burn should succeed")
	}

	// Can deposit tokens into the vault again
	tx, err = vaultContract.Deposit(r.ZEVMAuth, r.ETHZRC20Addr, big.NewInt(1e5))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("deposit should succeed")
	}

	r.Logger.Info("Operations all succeeded")
}
