package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	contracts "github.com/zeta-chain/zetacore/contracts/zevm"
	"math/big"
	"time"
)

func (sm *SmokeTest) TestERC20Deposit() {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	LoudPrintf("Deposit USDT ERC20 into ZEVM\n")
	USDT := sm.USDTERC20
	tx, err := USDT.Mint(sm.goerliAuth, big.NewInt(1e10))
	if err != nil {
		panic(err)
	}
	time.Sleep(BLOCK)
	receipt, err := sm.goerliClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Mint receipt tx hash: %s\n", tx.Hash().Hex())

	tx, err = USDT.Approve(sm.goerliAuth, sm.ERC20CustodyAddr, big.NewInt(1e10))
	if err != nil {
		panic(err)
	}
	time.Sleep(BLOCK)
	receipt, err = sm.goerliClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Printf("USDT Approve receipt tx hash: %s\n", tx.Hash().Hex())

	tx, err = sm.ERC20Custody.Deposit(sm.goerliAuth, DeployerAddress.Bytes(), sm.USDTERC20Addr, big.NewInt(1e6), nil)
	if err != nil {
		panic(err)
	}
	time.Sleep(BLOCK)
	receipt, err = sm.goerliClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deposit receipt tx hash: %s, status %d\n", receipt.TxHash.Hex(), receipt.Status)
	for _, log := range receipt.Logs {
		event, err := sm.ERC20Custody.ParseDeposited(*log)
		if err != nil {
			continue
		}
		fmt.Printf("Deposited event: \n")
		fmt.Printf("  Recipient address: %x, \n", event.Recipient)
		fmt.Printf("  ERC20 address: %s, \n", event.Asset.Hex())
		fmt.Printf("  Amount: %d, \n", event.Amount)
		fmt.Printf("  Message: %x, \n", event.Message)
	}
	WaitCctxMinedByInTxHash(tx.Hash().Hex(), sm.cctxClient)

	usdtZRC20, err := contracts.NewZRC20(ethcommon.HexToAddress(USDTZRC20Addr), sm.zevmClient)
	if err != nil {
		panic(err)
	}
	bal, err := usdtZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	if err != nil {
		panic(err)
	}
	fmt.Printf("balance of deployer on USDT ZRC20: %d\n", bal)
	supply, err := usdtZRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("supply of USDT ZRC20: %d\n", supply)
	if bal.Int64() != 1e6 {
		panic("balance is not correct")
	}
}
