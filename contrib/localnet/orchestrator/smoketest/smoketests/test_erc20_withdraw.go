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
		sm.Logger.Info("test finishes in %s", time.Since(startTime))
	}()

	bal, err := sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("balance of deployer on USDT ZRC20: %d", bal)
	supply, err := sm.USDTZRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("supply of USDT ZRC20: %d", supply)

	gasZRC20, gasFee, err := sm.USDTZRC20.WithdrawGasFee(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("gasZRC20: %s, gasFee: %d", gasZRC20.Hex(), gasFee)

	ethZRC20, err := zrc20.NewZRC20(gasZRC20, sm.ZevmClient)
	if err != nil {
		panic(err)
	}
	bal, err = ethZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("balance of deployer on ETH ZRC20: %d", bal)
	if bal.Int64() <= 0 {
		panic("not enough ETH ZRC20 balance!")
	}

	sm.Logger.InfoLoud("Withdraw USDT ZRC20")
	WithdrawERC20(sm, ethZRC20)

	sm.Logger.InfoLoud("Multiple withdraws USDT ZRC20")
	MultipleWithdraws(sm, ethZRC20)
}

func WithdrawERC20(sm *runner.SmokeTestRunner, ethZRC20 *zrc20.ZRC20) {
	// approve
	tx, err := ethZRC20.Approve(sm.ZevmAuth, sm.USDTZRC20Addr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(sm.ZevmClient, tx, sm.Logger)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	sm.Logger.Info("eth zrc20 approve receipt: status %d", receipt.Status)

	// withdraw
	tx, err = sm.USDTZRC20.Withdraw(sm.ZevmAuth, sm.DeployerAddress.Bytes(), big.NewInt(100))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx, sm.Logger)
	sm.Logger.Info("Receipt txhash %s status %d", receipt.TxHash, receipt.Status)
	for _, log := range receipt.Logs {
		event, err := sm.USDTZRC20.ParseWithdrawal(*log)
		if err != nil {
			continue
		}
		sm.Logger.Info("  logs: from %s, to %x, value %d, gasfee %d", event.From.Hex(), event.To, event.Value, event.Gasfee)
	}

	// verify the withdraw value
	cctx := utils.WaitCctxMinedByInTxHash(receipt.TxHash.Hex(), sm.CctxClient, sm.Logger)
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
	receipt := utils.MustWaitForTxReceipt(sm.ZevmClient, tx, sm.Logger)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	sm.Logger.Info("USDT ZRC20 approve receipt: status %d", receipt.Status)

	// approve gas token
	tx, err = ethZRC20.Approve(sm.ZevmAuth, withdrawerAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx, sm.Logger)
	if receipt.Status == 0 {
		panic("approve gas token failed")
	}
	sm.Logger.Info("eth zrc20 approve receipt: status %d", receipt.Status)

	// check the balance
	bal, err := sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("balance of deployer on USDT ZRC20: %d", bal)

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
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx, sm.Logger)
	if receipt.Status == 0 {
		panic("withdraw failed")
	}
	sm.Logger.Info("Withdraws receipt: status %d", receipt.Status)

	cctxs := utils.WaitCctxsMinedByInTxHash(tx.Hash().Hex(), sm.CctxClient, 10, sm.Logger)
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
	sm.Logger.Info("outTx hash %s", cctx.GetCurrentOutTxParam().OutboundTxHash)

	receipt, err := sm.GoerliClient.TransactionReceipt(
		context.Background(),
		ethcommon.HexToHash(cctx.GetCurrentOutTxParam().OutboundTxHash),
	)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("Receipt txhash %s status %d", receipt.TxHash, receipt.Status)
	for _, log := range receipt.Logs {
		event, err := sm.USDTERC20.ParseTransfer(*log)
		if err != nil {
			continue
		}
		sm.Logger.Info("  logs: from %s, to %s, value %d", event.From.Hex(), event.To.Hex(), event.Value)
		if event.Value.Int64() != amount {
			panic("value is not correct")
		}
	}
}
