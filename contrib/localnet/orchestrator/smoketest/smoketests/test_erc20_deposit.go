package smoketests

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	testcontract "github.com/zeta-chain/zetacore/testutil/contracts"
)

func TestERC20Deposit(sm *runner.SmokeTestRunner) {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	utils.LoudPrintf("Deposit USDT ERC20 into ZEVM\n")

	initialBal, err := sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	txhash := sm.DepositERC20(big.NewInt(1e18), []byte{})
	utils.WaitCctxMinedByInTxHash(txhash.Hex(), sm.CctxClient)

	bal, err := sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
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
	if diff.Int64() != 1e18 {
		panic("balance is not correct")
	}

	utils.LoudPrintf("Same-transaction multiple deposit USDT ERC20 into ZEVM\n")
	initialBal, err = sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	txhash = MultipleDeposits(sm, big.NewInt(1e9), big.NewInt(10))
	cctxs := utils.WaitCctxsMinedByInTxHash(txhash.Hex(), sm.CctxClient, 10)
	if len(cctxs) != 10 {
		panic(fmt.Sprintf("cctxs length is not correct: %d", len(cctxs)))
	}

	// check new balance is increased by 1e9 * 10
	bal, err = sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	diff = big.NewInt(0).Sub(bal, initialBal)
	if diff.Int64() != 1e10 {
		panic(fmt.Sprintf("balance difference is not correct: %d", diff.Int64()))
	}
}

func MultipleDeposits(sm *runner.SmokeTestRunner, amount, count *big.Int) ethcommon.Hash {
	// deploy depositor
	depositorAddr, _, depositor, err := testcontract.DeployDepositor(sm.GoerliAuth, sm.GoerliClient, sm.ERC20CustodyAddr)
	if err != nil {
		panic(err)
	}

	fullAmount := big.NewInt(0).Mul(amount, count)

	// mint
	tx, err := sm.USDTERC20.Mint(sm.GoerliAuth, fullAmount)
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
	if receipt.Status == 0 {
		panic("mint failed")
	}
	fmt.Printf("Mint receipt tx hash: %s\n", tx.Hash().Hex())

	// approve
	tx, err = sm.USDTERC20.Approve(sm.GoerliAuth, depositorAddr, fullAmount)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	fmt.Printf("USDT Approve receipt tx hash: %s\n", tx.Hash().Hex())

	// deposit
	tx, err = depositor.RunDeposits(sm.GoerliAuth, sm.DeployerAddress.Bytes(), sm.USDTERC20Addr, amount, []byte{}, count)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
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
		fmt.Printf("  Amount: %d, \n", event.Amount)
	}
	fmt.Printf("gas limit %d\n", sm.ZevmAuth.GasLimit)
	return tx.Hash()
}
