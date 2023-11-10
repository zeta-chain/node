package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/vault"
	"github.com/zeta-chain/zetacore/testutil/sample"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

func (sm *SmokeTest) TestPauseZRC20() {
	LoudPrintf("Test ZRC20 pause\n")
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()

	// Setup vault used to test zrc20 interactions
	fmt.Println("Deploying vault")
	vaultAddr, _, vaultContract, err := vault.DeployVault(sm.zevmAuth, sm.zevmClient)
	if err != nil {
		panic(err)
	}
	// Approving vault to spend ZRC20
	tx, err := sm.ETHZRC20.Approve(sm.zevmAuth, vaultAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt := MustWaitForTxReceipt(sm.zevmClient, tx)
	if receipt.Status == 0 {
		panic("Vault approval should succeed")
	}
	tx, err = sm.BTCZRC20.Approve(sm.zevmAuth, vaultAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	if receipt.Status == 0 {
		panic("Vault approval should succeed")
	}

	// Pause ETH ZRC20
	fmt.Println("Pausing ETH")
	msg := fungibletypes.NewMsgUpdateZRC20PausedStatus(
		FungibleAdminAddress,
		[]string{sm.ETHZRC20Addr.Hex()},
		fungibletypes.UpdatePausedStatusAction_PAUSE,
	)
	res, err := sm.zetaTxServer.BroadcastTx(FungibleAdminName, msg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("pause zrc20 tx hash: %s\n", res.TxHash)

	// Fetch and check pause status
	fcRes, err := sm.fungibleClient.ForeignCoins(context.Background(), &fungibletypes.QueryGetForeignCoinsRequest{
		Index: sm.ETHZRC20Addr.Hex(),
	})
	if err != nil {
		panic(err)
	}
	if !fcRes.GetForeignCoins().Paused {
		panic("ETH should be paused")
	} else {
		fmt.Printf("ETH is paused\n")
	}

	// Try operations with ETH ZRC20
	fmt.Println("Can no longer do operations on ETH ZRC20")
	tx, err = sm.ETHZRC20.Transfer(sm.zevmAuth, sample.EthAddress(), big.NewInt(1e5))
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	if receipt.Status == 1 {
		panic("transfer should fail")
	}
	tx, err = sm.ETHZRC20.Burn(sm.zevmAuth, big.NewInt(1e5))
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	if receipt.Status == 1 {
		panic("burn should fail")
	}

	// Operation on a contract that interact with ETH ZRC20 should fail
	fmt.Println("Vault contract can no longer interact with ETH ZRC20")
	tx, err = vaultContract.Deposit(sm.zevmAuth, sm.ETHZRC20Addr, big.NewInt(1e5))
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	if receipt.Status == 1 {
		panic("deposit should fail")
	}
	fmt.Println("Operations all failed")

	// Check we can still interact with BTC ZRC20
	fmt.Println("Check other ZRC20 can still be operated")
	tx, err = sm.BTCZRC20.Transfer(sm.zevmAuth, sample.EthAddress(), big.NewInt(1e3))
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	if receipt.Status == 0 {
		panic("BTC transfer should succeed")
	}
	tx, err = vaultContract.Deposit(sm.zevmAuth, sm.BTCZRC20Addr, big.NewInt(1e3))
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	if receipt.Status == 0 {
		panic("BTC vault deposit should succeed")
	}

	// Unpause ETH ZRC20
	fmt.Println("Unpausing ETH")
	msg = fungibletypes.NewMsgUpdateZRC20PausedStatus(
		FungibleAdminAddress,
		[]string{sm.ETHZRC20Addr.Hex()},
		fungibletypes.UpdatePausedStatusAction_UNPAUSE,
	)
	res, err = sm.zetaTxServer.BroadcastTx(FungibleAdminName, msg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("unpause zrc20 tx hash: %s\n", res.TxHash)

	// Fetch and check pause status
	fcRes, err = sm.fungibleClient.ForeignCoins(context.Background(), &fungibletypes.QueryGetForeignCoinsRequest{
		Index: sm.ETHZRC20Addr.Hex(),
	})
	if err != nil {
		panic(err)
	}
	if fcRes.GetForeignCoins().Paused {
		panic("ETH should be unpaused")
	} else {
		fmt.Printf("ETH is unpaused\n")
	}

	// Try operations with ETH ZRC20
	fmt.Println("Can do operations on ETH ZRC20 again")
	tx, err = sm.ETHZRC20.Transfer(sm.zevmAuth, sample.EthAddress(), big.NewInt(1e5))
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	if receipt.Status == 0 {
		panic("transfer should succeed")
	}
	tx, err = sm.ETHZRC20.Burn(sm.zevmAuth, big.NewInt(1e5))
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	if receipt.Status == 0 {
		panic("burn should succeed")
	}

	// Can deposit tokens into the vault again
	tx, err = vaultContract.Deposit(sm.zevmAuth, sm.ETHZRC20Addr, big.NewInt(1e5))
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	if receipt.Status == 0 {
		panic("deposit should succeed")
	}

	fmt.Println("Operations all succeeded")
}
