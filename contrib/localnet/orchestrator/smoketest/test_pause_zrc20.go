package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/zeta-chain/zetacore/testutil/sample"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

func (sm *SmokeTest) TestPauseZRC20() {
	LoudPrintf("Test ZRC20 pause\n")
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()

	bal, err := sm.ETHZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	if err != nil {
		panic(err)
	}
	fmt.Printf("balance of deployer on ETH ZRC20: %s\n", bal.String())

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
	tx, err := sm.ETHZRC20.Transfer(sm.zevmAuth, sample.EthAddress(), big.NewInt(1e5))
	if err != nil {
		panic(err)
	}
	receipt := MustWaitForTxReceipt(sm.zevmClient, tx)
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
	fmt.Println("Operations all failed")

	fmt.Println("Check other ZRC20 can still be operated")
	tx, err = sm.BTCZRC20.Transfer(sm.zevmAuth, sample.EthAddress(), big.NewInt(1e3))
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	if receipt.Status == 0 {
		panic("BTC transfer should succeed")
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
	fmt.Println("Operations all succeeded")
}
