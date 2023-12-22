package smoketests

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	zrc20 "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	testcontract "github.com/zeta-chain/zetacore/testutil/contracts"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestERC20Withdraw(sm *runner.SmokeTestRunner) {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()

	bal, err := sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	fmt.Printf("balance of deployer on USDT ZRC20: %d\n", bal)
	supply, err := sm.USDTZRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("supply of USDT ZRC20: %d\n", supply)

	gasZRC20, gasFee, err := sm.USDTZRC20.WithdrawGasFee(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("gasZRC20: %s, gasFee: %d\n", gasZRC20.Hex(), gasFee)

	ethZRC20, err := zrc20.NewZRC20(gasZRC20, sm.ZevmClient)
	if err != nil {
		panic(err)
	}
	bal, err = ethZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	fmt.Printf("balance of deployer on ETH ZRC20: %d\n", bal)
	if bal.Int64() <= 0 {
		panic("not enough ETH ZRC20 balance!")
	}

	utils.LoudPrintf("Withdraw USDT ZRC20\n")
	WithdrawERC20(sm, ethZRC20)

	utils.LoudPrintf("Multiple withdraws USDT ZRC20\n")
	MultipleWithdraws(sm, ethZRC20)
}

func WithdrawERC20(sm *runner.SmokeTestRunner, ethZRC20 *zrc20.ZRC20) {
	// approve
	tx, err := ethZRC20.Approve(sm.ZevmAuth, sm.USDTZRC20Addr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	fmt.Printf("eth zrc20 approve receipt: status %d\n", receipt.Status)

	// withdraw
	tx, err = sm.USDTZRC20.Withdraw(sm.ZevmAuth, sm.DeployerAddress.Bytes(), big.NewInt(100))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	fmt.Printf("Receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)
	for _, log := range receipt.Logs {
		event, err := sm.USDTZRC20.ParseWithdrawal(*log)
		if err != nil {
			continue
		}
		fmt.Printf("  logs: from %s, to %x, value %d, gasfee %d\n", event.From.Hex(), event.To, event.Value, event.Gasfee)
	}

	// verify the withdraw value
	cctx := utils.WaitCctxMinedByInTxHash(receipt.TxHash.Hex(), sm.CctxClient)
	verifyTransferAmountFromCCTX(sm, cctx, 100)
}

func MultipleWithdraws(sm *runner.SmokeTestRunner, ethZRC20 *zrc20.ZRC20) {
	// deploy withdrawer
	withdrawerAddr, _, withdrawer, err := testcontract.DeployWithdrawer(sm.ZevmAuth, sm.ZevmClient)
	if err != nil {
		panic(err)
	}

	// approve
	tx, err := sm.USDTZRC20.Approve(sm.ZevmAuth, withdrawerAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	fmt.Printf("USDT ZRC20 approve receipt: status %d\n", receipt.Status)

	// approve gas token
	tx, err = ethZRC20.Approve(sm.ZevmAuth, withdrawerAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	if receipt.Status == 0 {
		panic("approve gas token failed")
	}
	fmt.Printf("eth zrc20 approve receipt: status %d\n", receipt.Status)

	// check the balance
	bal, err := sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	fmt.Printf("balance of deployer on USDT ZRC20: %d\n", bal)

	if bal.Int64() < 1000 {
		panic("not enough USDT ZRC20 balance!")
	}

	// withdraw
	tx, err = withdrawer.RunWithdraws(
		sm.ZevmAuth,
		sm.DeployerAddress.Bytes(),
		sm.USDTZRC20Addr,
		big.NewInt(100),
		big.NewInt(10),
	)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	if receipt.Status == 0 {
		panic("withdraw failed")
	}
	fmt.Printf("Withdraws receipt: status %d\n", receipt.Status)

	cctxs := utils.WaitCctxsMinedByInTxHash(tx.Hash().Hex(), sm.CctxClient, 10)
	if len(cctxs) != 10 {
		panic(fmt.Sprintf("cctxs length is not correct: %d", len(cctxs)))
	}

	// verify the withdraw value
	for _, cctx := range cctxs {
		verifyTransferAmountFromCCTX(sm, cctx, 100)
	}
}

// verifyTransferAmountFromCCTX verifies the transfer amount from the CCTX on Goerli
func verifyTransferAmountFromCCTX(sm *runner.SmokeTestRunner, cctx *crosschaintypes.CrossChainTx, amount int64) {
	fmt.Printf("outTx hash %s\n", cctx.GetCurrentOutTxParam().OutboundTxHash)

	receipt, err := sm.GoerliClient.TransactionReceipt(
		context.Background(),
		ethcommon.HexToHash(cctx.GetCurrentOutTxParam().OutboundTxHash),
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)
	for _, log := range receipt.Logs {
		event, err := sm.USDTERC20.ParseTransfer(*log)
		if err != nil {
			continue
		}
		fmt.Printf("  logs: from %s, to %s, value %d\n", event.From.Hex(), event.To.Hex(), event.Value)
		if event.Value.Int64() != amount {
			panic("value is not correct")
		}
	}
}
