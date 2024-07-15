package e2etests

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestERC20DepositAndCallRefund(r *runner.E2ERunner, _ []string) {
	// Get the initial balance of the deployer
	initialBal, err := r.ERC20ZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)

	r.Logger.Info("Sending a deposit that should revert without a liquidity pool makes the cctx aborted")

	amount := big.NewInt(1e4)

	// send the deposit
	inboundHash, err := sendInvalidERC20Deposit(r, amount)
	require.NoError(r, err)

	// There is no liquidity pool, therefore the cctx should abort
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, inboundHash, r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, types.CctxStatus_Aborted)
	require.False(r, cctx.CctxStatus.IsAbortRefunded, "expected cctx status to be not refunded")

	r.Logger.CCTX(*cctx, "deposit")
	r.Logger.Info("Refunding the cctx via admin")

	msg := types.NewMsgRefundAbortedCCTX(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		cctx.Index,
		r.EVMAddress().String(),
	)

	_, err = r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, msg)
	require.NoError(r, err)

	// Check that the erc20 in the aborted cctx was refunded on ZetaChain
	newBalance, err := r.ERC20ZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)

	expectedBalance := initialBal.Add(initialBal, amount)
	require.Equal(
		r,
		0,
		newBalance.Cmp(expectedBalance),
		"expected balance to be %s after refund; got %s",
		expectedBalance.String(),
		newBalance.String(),
	)
	r.Logger.Info("CCTX has been aborted on ZetaChain")

	// test refund when there is a liquidity pool
	r.Logger.Info("Sending a deposit that should revert with a liquidity pool")

	r.Logger.Info("Creating the liquidity pool USTD/ZETA")
	err = createZetaERC20LiquidityPool(r)
	require.NoError(r, err)

	r.Logger.Info("Liquidity pool created")

	erc20Balance, err := r.ERC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)

	// send the deposit
	amount = big.NewInt(1e7)
	inboundHash, err = sendInvalidERC20Deposit(r, amount)
	require.NoError(r, err)

	erc20BalanceAfterSend := big.NewInt(0).Sub(erc20Balance, amount)

	// there is a liquidity pool, therefore the cctx should revert
	// the revert tx creation will fail because the sender, used as the recipient, is not defined in the cctx
	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, inboundHash, r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, types.CctxStatus_Reverted)

	// get revert tx
	revertTxHash := cctx.GetCurrentOutboundParam().Hash
	receipt, err := r.EVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(revertTxHash))
	require.NoError(r, err)
	utils.RequireTxSuccessful(r, receipt)

	// check that the erc20 in the reverted cctx was refunded on EVM
	erc20BalanceAfterRefund, err := r.ERC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)

	// the new balance must be higher than the previous one because of the revert refund
	require.Equal(
		r,
		-1,
		erc20BalanceAfterSend.Cmp(erc20BalanceAfterRefund),
		"expected balance to be higher after refund than after send %s < %s",
		erc20BalanceAfterSend.String(),
		erc20BalanceAfterRefund.String(),
	)

	// it must also be lower than the previous balance + the amount because of the gas fee for the revert tx
	require.Equal(
		r,
		-1,
		erc20BalanceAfterRefund.Cmp(erc20Balance),
		"expected balance to be lower after refund than before send %s < %s",
		erc20BalanceAfterRefund.String(),
		erc20Balance.String(),
	)

	r.Logger.Info("ERC20 CCTX successfully reverted")
	r.Logger.Info("\tbalance before refund: %s", erc20Balance.String())
	r.Logger.Info("\tamount: %s", amount.String())
	r.Logger.Info("\tbalance after refund: %s", erc20BalanceAfterRefund.String())
}

func createZetaERC20LiquidityPool(r *runner.E2ERunner) error {
	amount := big.NewInt(1e10)
	txHash := r.DepositERC20WithAmountAndMessage(r.EVMAddress(), amount, []byte{})
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
		r.EVMAddress(),
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
		r.EVMAddress().Bytes(),
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
