package smoketests

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestERC20DepositAndCallRefund(sm *runner.SmokeTestRunner) {
	// Get the initial balance of the deployer
	initialBal, err := sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}

	sm.Logger.Info("Sending a deposit that should revert without a liquidity pool makes the cctx aborted")

	amount := big.NewInt(1e4)

	// send the deposit
	inTxHash, err := sendInvalidUSDTDeposit(sm, amount)
	if err != nil {
		panic(err)
	}

	// There is no liquidity pool, therefore the cctx should abort
	cctx := utils.WaitCctxMinedByInTxHash(inTxHash, sm.CctxClient, sm.Logger)
	if cctx.CctxStatus.Status != types.CctxStatus_Aborted {
		panic(fmt.Sprintf("expected cctx status to be Aborted; got %s", cctx.CctxStatus.Status))
	}

	// Check that the erc20 in the aborted cctx was refunded on ZetaChain
	newBalance, err := sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	expectedBalance := initialBal.Add(initialBal, amount)
	if newBalance.Cmp(expectedBalance) != 0 {
		panic(fmt.Sprintf("expected balance to be %s after refund; got %s", expectedBalance.String(), newBalance.String()))
	}
	sm.Logger.Info("CCTX has been aborted and the erc20 has been refunded on ZetaChain")

	// test refund when there is a liquidity pool
	sm.Logger.Info("Sending a deposit that should revert with a liquidity pool")

	sm.Logger.Info("Creating the liquidity pool USTD/ZETA")
	err = createZetaERC20LiquidityPool(sm)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("Liquidity pool created")

	goerliBalance, err := sm.USDTERC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}

	// send the deposit
	amount = big.NewInt(1e7)
	inTxHash, err = sendInvalidUSDTDeposit(sm, amount)
	if err != nil {
		panic(err)
	}
	goerliBalanceAfterSend := big.NewInt(0).Sub(goerliBalance, amount)

	// there is a liquidity pool, therefore the cctx should revert
	cctx = utils.WaitCctxMinedByInTxHash(inTxHash, sm.CctxClient, sm.Logger)

	// the revert tx creation will fail because the sender, used as the recipient, is not defined in the cctx
	if cctx.CctxStatus.Status != types.CctxStatus_Reverted {
		panic(fmt.Sprintf("expected cctx status to be PendingRevert; got %s", cctx.CctxStatus.Status))
	}

	// get revert tx
	revertTxHash := cctx.GetCurrentOutTxParam().OutboundTxHash
	receipt, err := sm.GoerliClient.TransactionReceipt(context.Background(), ethcommon.HexToHash(revertTxHash))
	if err != nil {
		panic(err)
	}
	if receipt.Status == 0 {
		panic("expected the revert tx receipt to have status 1; got 0")
	}

	// check that the erc20 in the reverted cctx was refunded on Goerli
	goerliBalanceAfterRefund, err := sm.USDTERC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	// the new balance must be higher than the previous one because of the revert refund
	if goerliBalanceAfterSend.Cmp(goerliBalanceAfterRefund) != -1 {
		panic(fmt.Sprintf(
			"expected balance to be higher after refund than after send %s < %s",
			goerliBalanceAfterSend.String(),
			goerliBalanceAfterRefund.String(),
		))
	}
	// it must also be lower than the previous balance + the amount because of the gas fee for the revert tx
	if goerliBalanceAfterRefund.Cmp(goerliBalance) != -1 {
		panic(fmt.Sprintf(
			"expected balance to be lower after refund than before send %s < %s",
			goerliBalanceAfterRefund.String(),
			goerliBalance.String()),
		)
	}

	sm.Logger.Info("ERC20 CCTX successfully reverted")
	sm.Logger.Info("\tbalance before refund: %s", goerliBalance.String())
	sm.Logger.Info("\tamount: %s", amount.String())
	sm.Logger.Info("\tbalance after refund: %s", goerliBalanceAfterRefund.String())
}

func createZetaERC20LiquidityPool(sm *runner.SmokeTestRunner) error {
	amount := big.NewInt(1e10)
	txHash := sm.DepositERC20WithAmountAndMessage(amount, []byte{})
	utils.WaitCctxMinedByInTxHash(txHash.Hex(), sm.CctxClient, sm.Logger)

	tx, err := sm.USDTZRC20.Approve(sm.ZevmAuth, sm.UniswapV2RouterAddr, big.NewInt(1e10))
	if err != nil {
		return err
	}
	receipt := utils.MustWaitForTxReceipt(sm.ZevmClient, tx, sm.Logger)
	if receipt.Status == 0 {
		return errors.New("approve failed")
	}

	previousValue := sm.ZevmAuth.Value
	sm.ZevmAuth.Value = big.NewInt(1e10)
	tx, err = sm.UniswapV2Router.AddLiquidityETH(
		sm.ZevmAuth,
		sm.USDTZRC20Addr,
		amount,
		big.NewInt(0),
		big.NewInt(0),
		sm.DeployerAddress,
		big.NewInt(time.Now().Add(10*time.Minute).Unix()),
	)
	sm.ZevmAuth.Value = previousValue
	if err != nil {
		return err
	}
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx, sm.Logger)
	if receipt.Status == 0 {
		return fmt.Errorf("add liquidity failed")
	}

	return nil
}

func sendInvalidUSDTDeposit(sm *runner.SmokeTestRunner, amount *big.Int) (string, error) {
	USDT := sm.USDTERC20
	tx, err := USDT.Approve(sm.GoerliAuth, sm.ERC20CustodyAddr, amount)
	if err != nil {
		return "", err
	}
	receipt := utils.MustWaitForTxReceipt(sm.GoerliClient, tx, sm.Logger)
	sm.Logger.Info("USDT Approve receipt tx hash: %s", tx.Hash().Hex())

	tx, err = sm.ERC20Custody.Deposit(
		sm.GoerliAuth,
		sm.DeployerAddress.Bytes(),
		sm.USDTERC20Addr,
		amount,
		[]byte("this is an invalid msg that will cause the contract to revert"),
	)
	if err != nil {
		return "", err
	}

	sm.Logger.Info("GOERLI tx sent: %s; to %s, nonce %d", tx.Hash().String(), tx.To().Hex(), tx.Nonce())
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, tx, sm.Logger)
	if receipt.Status == 0 {
		return "", errors.New("expected the tx receipt to have status 1; got 0")
	}
	sm.Logger.Info("GOERLI tx receipt: %d", receipt.Status)
	sm.Logger.Info("  tx hash: %s", receipt.TxHash.String())
	sm.Logger.Info("  to: %s", tx.To().String())
	sm.Logger.Info("  value: %d", tx.Value())
	sm.Logger.Info("  block num: %d", receipt.BlockNumber)

	return tx.Hash().Hex(), nil
}
