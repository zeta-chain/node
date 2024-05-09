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

func TestERC20DepositAndCallRefund(r *runner.E2ERunner, _ []string) {
	// Get the initial balance of the deployer
	initialBal, err := r.ERC20ZRC20.BalanceOf(&bind.CallOpts{}, r.DeployerAddress)
	if err != nil {
		panic(err)
	}

	r.Logger.Info("Sending a deposit that should revert without a liquidity pool makes the cctx aborted")

	amount := big.NewInt(1e4)

	// send the deposit
	inboundHash, err := sendInvalidERC20Deposit(r, amount)
	if err != nil {
		panic(err)
	}

	// There is no liquidity pool, therefore the cctx should abort
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, inboundHash, r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	if cctx.CctxStatus.Status != types.CctxStatus_Aborted {
		panic(fmt.Sprintf("expected cctx status to be Aborted; got %s", cctx.CctxStatus.Status))
	}

	if cctx.CctxStatus.IsAbortRefunded != false {
		panic(fmt.Sprintf("expected cctx status to be not refunded; got %t", cctx.CctxStatus.IsAbortRefunded))
	}

	r.Logger.Info("Refunding the cctx via admin")
	msg := types.NewMsgRefundAbortedCCTX(
		r.ZetaTxServer.GetAccountAddress(0),
		cctx.Index,
		r.DeployerAddress.String())
	_, err = r.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msg)
	if err != nil {
		panic(err)
	}

	// Check that the erc20 in the aborted cctx was refunded on ZetaChain
	newBalance, err := r.ERC20ZRC20.BalanceOf(&bind.CallOpts{}, r.DeployerAddress)
	if err != nil {
		panic(err)
	}
	expectedBalance := initialBal.Add(initialBal, amount)
	if newBalance.Cmp(expectedBalance) != 0 {
		panic(fmt.Sprintf("expected balance to be %s after refund; got %s", expectedBalance.String(), newBalance.String()))
	}
	r.Logger.Info("CCTX has been aborted on ZetaChain")

	// test refund when there is a liquidity pool
	r.Logger.Info("Sending a deposit that should revert with a liquidity pool")

	r.Logger.Info("Creating the liquidity pool USTD/ZETA")
	err = createZetaERC20LiquidityPool(r)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("Liquidity pool created")

	erc20Balance, err := r.ERC20.BalanceOf(&bind.CallOpts{}, r.DeployerAddress)
	if err != nil {
		panic(err)
	}

	// send the deposit
	amount = big.NewInt(1e7)
	inboundHash, err = sendInvalidERC20Deposit(r, amount)
	if err != nil {
		panic(err)
	}
	erc20BalanceAfterSend := big.NewInt(0).Sub(erc20Balance, amount)

	// there is a liquidity pool, therefore the cctx should revert
	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, inboundHash, r.CctxClient, r.Logger, r.CctxTimeout)

	// the revert tx creation will fail because the sender, used as the recipient, is not defined in the cctx
	if cctx.CctxStatus.Status != types.CctxStatus_Reverted {
		panic(fmt.Sprintf(
			"expected cctx status to be PendingRevert; got %s, aborted message: %s",
			cctx.CctxStatus.Status,
			cctx.CctxStatus.StatusMessage,
		))
	}

	// get revert tx
	revertTxHash := cctx.GetCurrentOutboundParam().Hash
	receipt, err := r.EVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(revertTxHash))
	if err != nil {
		panic(err)
	}
	if receipt.Status == 0 {
		panic("expected the revert tx receipt to have status 1; got 0")
	}

	// check that the erc20 in the reverted cctx was refunded on EVM
	erc20BalanceAfterRefund, err := r.ERC20.BalanceOf(&bind.CallOpts{}, r.DeployerAddress)
	if err != nil {
		panic(err)
	}
	// the new balance must be higher than the previous one because of the revert refund
	if erc20BalanceAfterSend.Cmp(erc20BalanceAfterRefund) != -1 {
		panic(fmt.Sprintf(
			"expected balance to be higher after refund than after send %s < %s",
			erc20BalanceAfterSend.String(),
			erc20BalanceAfterRefund.String(),
		))
	}
	// it must also be lower than the previous balance + the amount because of the gas fee for the revert tx
	if erc20BalanceAfterRefund.Cmp(erc20Balance) != -1 {
		panic(fmt.Sprintf(
			"expected balance to be lower after refund than before send %s < %s",
			erc20BalanceAfterRefund.String(),
			erc20Balance.String()),
		)
	}

	r.Logger.Info("ERC20 CCTX successfully reverted")
	r.Logger.Info("\tbalance before refund: %s", erc20Balance.String())
	r.Logger.Info("\tamount: %s", amount.String())
	r.Logger.Info("\tbalance after refund: %s", erc20BalanceAfterRefund.String())
}

func createZetaERC20LiquidityPool(r *runner.E2ERunner) error {
	amount := big.NewInt(1e10)
	txHash := r.DepositERC20WithAmountAndMessage(r.DeployerAddress, amount, []byte{})
	utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)

	tx, err := r.ERC20ZRC20.Approve(r.ZEVMAuth, r.UniswapV2RouterAddr, big.NewInt(1e10))
	if err != nil {
		return err
	}
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		return errors.New("approve failed")
	}

	previousValue := r.ZEVMAuth.Value
	r.ZEVMAuth.Value = big.NewInt(1e10)
	tx, err = r.UniswapV2Router.AddLiquidityETH(
		r.ZEVMAuth,
		r.ERC20ZRC20Addr,
		amount,
		big.NewInt(0),
		big.NewInt(0),
		r.DeployerAddress,
		big.NewInt(time.Now().Add(10*time.Minute).Unix()),
	)
	r.ZEVMAuth.Value = previousValue
	if err != nil {
		return err
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		return fmt.Errorf("add liquidity failed")
	}

	return nil
}

func sendInvalidERC20Deposit(r *runner.E2ERunner, amount *big.Int) (string, error) {
	tx, err := r.ERC20.Approve(r.EVMAuth, r.ERC20CustodyAddr, amount)
	if err != nil {
		return "", err
	}
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.Info("ERC20 Approve receipt tx hash: %s", tx.Hash().Hex())

	tx, err = r.ERC20Custody.Deposit(
		r.EVMAuth,
		r.DeployerAddress.Bytes(),
		r.ERC20Addr,
		amount,
		[]byte("this is an invalid msg that will cause the contract to revert"),
	)
	if err != nil {
		return "", err
	}

	r.Logger.Info("EVM tx sent: %s; to %s, nonce %d", tx.Hash().String(), tx.To().Hex(), tx.Nonce())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		return "", errors.New("expected the tx receipt to have status 1; got 0")
	}
	r.Logger.Info("EVM tx receipt: %d", receipt.Status)
	r.Logger.Info("  tx hash: %s", receipt.TxHash.String())
	r.Logger.Info("  to: %s", tx.To().String())
	r.Logger.Info("  value: %d", tx.Value())
	r.Logger.Info("  block num: %d", receipt.BlockNumber)

	return tx.Hash().Hex(), nil
}
