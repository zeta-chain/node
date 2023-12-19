package smoketests

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/vault"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestPauseZRC20(sm *runner.SmokeTestRunner) {
	utils.LoudPrintf("Test ZRC20 pause\n")
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()

	// Setup vault used to test zrc20 interactions
	fmt.Println("Deploying vault")
	vaultAddr, _, vaultContract, err := vault.DeployVault(sm.ZevmAuth, sm.ZevmClient)
	if err != nil {
		panic(err)
	}
	// Approving vault to spend ZRC20
	tx, err := sm.ETHZRC20.Approve(sm.ZevmAuth, vaultAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	if receipt.Status == 0 {
		panic("Vault approval should succeed")
	}
	tx, err = sm.BTCZRC20.Approve(sm.ZevmAuth, vaultAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	if receipt.Status == 0 {
		panic("Vault approval should succeed")
	}

	// Pause ETH ZRC20
	fmt.Println("Pausing ETH")
	msg := fungibletypes.NewMsgUpdateZRC20PausedStatus(
		sm.ZetaTxServer.GetAccountAddress(0),
		[]string{sm.ETHZRC20Addr.Hex()},
		fungibletypes.UpdatePausedStatusAction_PAUSE,
	)
	res, err := sm.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("pause zrc20 tx hash: %s\n", res.TxHash)

	// Fetch and check pause status
	fcRes, err := sm.FungibleClient.ForeignCoins(context.Background(), &fungibletypes.QueryGetForeignCoinsRequest{
		Index: sm.ETHZRC20Addr.Hex(),
	})
	if err != nil {
		panic(err)
	}
	if !fcRes.GetForeignCoins().Paused {
		panic("ETH should be paused")
	}
	fmt.Printf("ETH is paused\n")

	// Try operations with ETH ZRC20
	fmt.Println("Can no longer do operations on ETH ZRC20")
	tx, err = sm.ETHZRC20.Transfer(sm.ZevmAuth, sample.EthAddress(), big.NewInt(1e5))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	if receipt.Status == 1 {
		panic("transfer should fail")
	}
	tx, err = sm.ETHZRC20.Burn(sm.ZevmAuth, big.NewInt(1e5))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	if receipt.Status == 1 {
		panic("burn should fail")
	}

	// Operation on a contract that interact with ETH ZRC20 should fail
	fmt.Printf("Vault contract can no longer interact with ETH ZRC20: %s\n", sm.ETHZRC20Addr.Hex())
	tx, err = vaultContract.Deposit(sm.ZevmAuth, sm.ETHZRC20Addr, big.NewInt(1e5))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	if receipt.Status == 1 {
		panic("deposit should fail")
	}
	fmt.Println("Operations all failed")

	// Check we can still interact with BTC ZRC20
	fmt.Println("Check other ZRC20 can still be operated")
	tx, err = sm.BTCZRC20.Transfer(sm.ZevmAuth, sample.EthAddress(), big.NewInt(1e3))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	if receipt.Status == 0 {
		panic("BTC transfer should succeed")
	}
	tx, err = vaultContract.Deposit(sm.ZevmAuth, sm.BTCZRC20Addr, big.NewInt(1e3))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	if receipt.Status == 0 {
		panic("BTC vault deposit should succeed")
	}

	// Check deposit revert when paused
	signedTx, err := sm.SendEther(sm.TSSAddress, big.NewInt(1e17), nil)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, signedTx)
	if receipt.Status == 0 {
		panic("deposit eth tx failed")
	}
	cctx := utils.WaitCctxMinedByInTxHash(signedTx.Hash().Hex(), sm.CctxClient)
	if cctx.CctxStatus.Status != types.CctxStatus_Reverted {
		panic(fmt.Sprintf("expected cctx status to be Reverted; got %s", cctx.CctxStatus.Status))
	}
	fmt.Println("CCTX has been reverted")

	// Unpause ETH ZRC20
	fmt.Println("Unpausing ETH")
	msg = fungibletypes.NewMsgUpdateZRC20PausedStatus(
		sm.ZetaTxServer.GetAccountAddress(0),
		[]string{sm.ETHZRC20Addr.Hex()},
		fungibletypes.UpdatePausedStatusAction_UNPAUSE,
	)
	res, err = sm.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("unpause zrc20 tx hash: %s\n", res.TxHash)

	// Fetch and check pause status
	fcRes, err = sm.FungibleClient.ForeignCoins(context.Background(), &fungibletypes.QueryGetForeignCoinsRequest{
		Index: sm.ETHZRC20Addr.Hex(),
	})
	if err != nil {
		panic(err)
	}
	if fcRes.GetForeignCoins().Paused {
		panic("ETH should be unpaused")
	}
	fmt.Printf("ETH is unpaused\n")

	// Try operations with ETH ZRC20
	fmt.Println("Can do operations on ETH ZRC20 again")
	tx, err = sm.ETHZRC20.Transfer(sm.ZevmAuth, sample.EthAddress(), big.NewInt(1e5))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	if receipt.Status == 0 {
		panic("transfer should succeed")
	}
	tx, err = sm.ETHZRC20.Burn(sm.ZevmAuth, big.NewInt(1e5))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	if receipt.Status == 0 {
		panic("burn should succeed")
	}

	// Can deposit tokens into the vault again
	tx, err = vaultContract.Deposit(sm.ZevmAuth, sm.ETHZRC20Addr, big.NewInt(1e5))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	if receipt.Status == 0 {
		panic("deposit should succeed")
	}

	fmt.Println("Operations all succeeded")
}
