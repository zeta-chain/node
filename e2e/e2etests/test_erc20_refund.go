package e2etests

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestERC20DepositAndCallRefund(r *runner.E2ERunner) {
	// Get the initial balance of the deployer
	initialBal, err := r.USDTZRC20.BalanceOf(&bind.CallOpts{}, r.DeployerAddress)
	if err != nil {
		panic(err)
	}

	r.Logger.Info("Sending a deposit that should revert without a liquidity pool makes the cctx aborted")

	amount := big.NewInt(1e4)

	// send the deposit
	inTxHash, err := sendInvalidUSDTDeposit(r, amount)
	if err != nil {
		panic(err)
	}

	// There is no liquidity pool, therefore the cctx should abort
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, inTxHash, r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	if cctx.CctxStatus.Status != types.CctxStatus_Aborted {
		panic(fmt.Sprintf("expected cctx status to be Aborted; got %s", cctx.CctxStatus.Status))
	}

	// Check that the erc20 in the aborted cctx was refunded on ZetaChain
	newBalance, err := r.USDTZRC20.BalanceOf(&bind.CallOpts{}, r.DeployerAddress)
	if err != nil {
		panic(err)
	}
	expectedBalance := initialBal.Add(initialBal, amount)
	if newBalance.Cmp(expectedBalance) != 0 {
		panic(fmt.Sprintf("expected balance to be %s after refund; got %s", expectedBalance.String(), newBalance.String()))
	}
	r.Logger.Info("CCTX has been aborted and the erc20 has been refunded on ZetaChain")

	// test refund when there is a liquidity pool
	r.Logger.Info("Sending a deposit that should revert with a liquidity pool")

	r.Logger.Info("Creating the liquidity pool USTD/ZETA")
	err = createZetaERC20LiquidityPool(r)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("Liquidity pool created")

	goerliBalance, err := r.USDTERC20.BalanceOf(&bind.CallOpts{}, r.DeployerAddress)
	if err != nil {
		panic(err)
	}

	// send the deposit
	amount = big.NewInt(1e7)
	inTxHash, err = sendInvalidUSDTDeposit(r, amount)
	if err != nil {
		panic(err)
	}
	goerliBalanceAfterSend := big.NewInt(0).Sub(goerliBalance, amount)

	// there is a liquidity pool, therefore the cctx should revert
	cctx = utils.WaitCctxMinedByInTxHash(r.Ctx, inTxHash, r.CctxClient, r.Logger, r.CctxTimeout)

	// the revert tx creation will fail because the sender, used as the recipient, is not defined in the cctx
	if cctx.CctxStatus.Status != types.CctxStatus_Reverted {
		panic(fmt.Sprintf(
			"expected cctx status to be PendingRevert; got %s, aborted message: %s",
			cctx.CctxStatus.Status,
			cctx.CctxStatus.StatusMessage,
		))
	}

	// get revert tx
	revertTxHash := cctx.GetCurrentOutTxParam().OutboundTxHash
	receipt, err := r.GoerliClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(revertTxHash))
	if err != nil {
		panic(err)
	}
	if receipt.Status == 0 {
		panic("expected the revert tx receipt to have status 1; got 0")
	}

	// check that the erc20 in the reverted cctx was refunded on Goerli
	goerliBalanceAfterRefund, err := r.USDTERC20.BalanceOf(&bind.CallOpts{}, r.DeployerAddress)
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

	r.Logger.Info("ERC20 CCTX successfully reverted")
	r.Logger.Info("\tbalance before refund: %s", goerliBalance.String())
	r.Logger.Info("\tamount: %s", amount.String())
	r.Logger.Info("\tbalance after refund: %s", goerliBalanceAfterRefund.String())
}

func createZetaERC20LiquidityPool(r *runner.E2ERunner) error {
	amount := big.NewInt(1e10)
	txHash := r.DepositERC20WithAmountAndMessage(amount, []byte{})
	utils.WaitCctxMinedByInTxHash(r.Ctx, txHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)

	tx, err := r.USDTZRC20.Approve(r.ZevmAuth, r.UniswapV2RouterAddr, big.NewInt(1e10))
	if err != nil {
		return err
	}
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		return errors.New("approve failed")
	}

	previousValue := r.ZevmAuth.Value
	r.ZevmAuth.Value = big.NewInt(1e10)
	tx, err = r.UniswapV2Router.AddLiquidityETH(
		r.ZevmAuth,
		r.USDTZRC20Addr,
		amount,
		big.NewInt(0),
		big.NewInt(0),
		r.DeployerAddress,
		big.NewInt(time.Now().Add(10*time.Minute).Unix()),
	)
	r.ZevmAuth.Value = previousValue
	if err != nil {
		return err
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		return fmt.Errorf("add liquidity failed")
	}

	return nil
}

func sendInvalidUSDTDeposit(r *runner.E2ERunner, amount *big.Int) (string, error) {
	USDT := r.USDTERC20
	tx, err := USDT.Approve(r.GoerliAuth, r.ERC20CustodyAddr, amount)
	if err != nil {
		return "", err
	}
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.GoerliClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.Info("USDT Approve receipt tx hash: %s", tx.Hash().Hex())

	tx, err = r.ERC20Custody.Deposit(
		r.GoerliAuth,
		r.DeployerAddress.Bytes(),
		r.USDTERC20Addr,
		amount,
		[]byte("this is an invalid msg that will cause the contract to revert"),
	)
	if err != nil {
		return "", err
	}

	r.Logger.Info("GOERLI tx sent: %s; to %s, nonce %d", tx.Hash().String(), tx.To().Hex(), tx.Nonce())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.GoerliClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		return "", errors.New("expected the tx receipt to have status 1; got 0")
	}
	r.Logger.Info("GOERLI tx receipt: %d", receipt.Status)
	r.Logger.Info("  tx hash: %s", receipt.TxHash.String())
	r.Logger.Info("  to: %s", tx.To().String())
	r.Logger.Info("  value: %d", tx.Value())
	r.Logger.Info("  block num: %d", receipt.BlockNumber)

	return tx.Hash().Hex(), nil
}
