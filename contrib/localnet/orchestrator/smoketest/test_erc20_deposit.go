package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	testcontract "github.com/zeta-chain/zetacore/testutil/contracts"
)

func (sm *SmokeTest) TestERC20Deposit() {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	LoudPrintf("Deposit USDT ERC20 into ZEVM\n")

	initialBal, err := sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	if err != nil {
		panic(err)
	}
	txhash := sm.DepositERC20(big.NewInt(1e9), []byte{})
	WaitCctxMinedByInTxHash(txhash.Hex(), sm.cctxClient)

	bal, err := sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	if err != nil {
		panic(err)
	}

	diff := big.NewInt(0)
	diff.Sub(bal, initialBal)

	fmt.Printf("balance of deployer on USDT ZRC20: %d\n", bal)
	supply, err := sm.USDTZRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("supply of USDT ZRC20: %d\n", supply)
	if diff.Int64() != 1e9 {
		panic("balance is not correct")
	}

	LoudPrintf("Same-transaction multiple deposit USDT ERC20 into ZEVM\n")
	initialBal, err = sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	if err != nil {
		panic(err)
	}
	txhash = sm.MultipleDeposits(big.NewInt(1e9), big.NewInt(10))
	cctxs := WaitCctxsMinedByInTxHash(txhash.Hex(), sm.cctxClient)
	if len(cctxs) != 10 {
		panic(fmt.Sprintf("cctxs length is not correct: %d", len(cctxs)))
	}

	// check new balance is increased by 1e9 * 10
	bal, err = sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	if err != nil {
		panic(err)
	}
	diff = big.NewInt(0).Sub(bal, initialBal)
	if diff.Int64() != 1e10 {
		panic(fmt.Sprintf("balance difference is not correct: %d", diff.Int64()))
	}
}

func (sm *SmokeTest) DepositERC20(amount *big.Int, msg []byte) ethcommon.Hash {
	USDT := sm.USDTERC20
	tx, err := USDT.Mint(sm.goerliAuth, big.NewInt(1e10))
	if err != nil {
		panic(err)
	}
	receipt := MustWaitForTxReceipt(sm.goerliClient, tx)
	fmt.Printf("Mint receipt tx hash: %s\n", tx.Hash().Hex())

	tx, err = USDT.Approve(sm.goerliAuth, sm.ERC20CustodyAddr, big.NewInt(1e10))
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.goerliClient, tx)
	fmt.Printf("USDT Approve receipt tx hash: %s\n", tx.Hash().Hex())

	tx, err = sm.ERC20Custody.Deposit(sm.goerliAuth, DeployerAddress.Bytes(), sm.USDTERC20Addr, amount, msg)
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.goerliClient, tx)
	if receipt.Status == 0 {
		panic("deposit failed")
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
	fmt.Printf("gas limit %d\n", sm.zevmAuth.GasLimit)
	return tx.Hash()
}

func (sm *SmokeTest) MultipleDeposits(amount, count *big.Int) ethcommon.Hash {
	// deploy depositor
	depositorAddr, _, depositor, err := testcontract.DeployDepositor(sm.goerliAuth, sm.goerliClient, sm.ERC20CustodyAddr)
	if err != nil {
		panic(err)
	}

	// mint
	tx, err := sm.USDTERC20.Mint(sm.goerliAuth, big.NewInt(0).Mul(amount, count))
	if err != nil {
		panic(err)
	}
	receipt := MustWaitForTxReceipt(sm.goerliClient, tx)
	if receipt.Status == 0 {
		panic("mint failed")
	}
	fmt.Printf("Mint receipt tx hash: %s\n", tx.Hash().Hex())

	// approve
	tx, err = sm.USDTERC20.Approve(sm.goerliAuth, depositorAddr, big.NewInt(1e10))
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.goerliClient, tx)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	fmt.Printf("USDT Approve receipt tx hash: %s\n", tx.Hash().Hex())

	// deposit
	tx, err = depositor.RunDeposits(sm.goerliAuth, DeployerAddress.Bytes(), sm.USDTERC20Addr, amount, []byte{}, count)
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.goerliClient, tx)
	if receipt.Status == 0 {
		panic("deposits failed")
	}
	fmt.Printf("Deposits receipt tx hash: %s\n", tx.Hash().Hex())

	for _, log := range receipt.Logs {
		event, err := sm.ERC20Custody.ParseDeposited(*log)
		if err != nil {
			continue
		}
		fmt.Printf("Multiple deposit event: \n")
		fmt.Printf("  Recipient address: %x, \n", event.Recipient)
		fmt.Printf("  ERC20 address: %s, \n", event.Asset.Hex())
		fmt.Printf("  Amount: %d, \n", event.Amount)
		fmt.Printf("  Message: %x, \n", event.Message)
	}
	fmt.Printf("gas limit %d\n", sm.zevmAuth.GasLimit)
	return tx.Hash()
}
