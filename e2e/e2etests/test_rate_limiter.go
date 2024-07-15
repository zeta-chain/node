package e2etests

import (
	"context"
	"fmt"
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// WithdrawType is the type of withdraw to perform in the test
type withdrawType string

const (
	withdrawTypeZETA  withdrawType = "ZETA"
	withdrawTypeETH   withdrawType = "ETH"
	withdrawTypeERC20 withdrawType = "ERC20"

	rateLimiterWithdrawNumber = 5
)

func TestRateLimiter(r *runner.E2ERunner, _ []string) {
	r.Logger.Info("TestRateLimiter")

	// rateLimiterFlags are the rate limiter flags for the test
	rateLimiterFlags := crosschaintypes.RateLimiterFlags{
		Enabled: true,
		Rate:    sdk.NewUint(1e17).MulUint64(5), // 0.5 ZETA this value is used so rate is reached
		Window:  10,
		Conversions: []crosschaintypes.Conversion{
			{
				Zrc20: r.ETHZRC20Addr.Hex(),
				Rate:  sdk.NewDec(2), // 1 ETH = 2 ZETA
			},
			{
				Zrc20: r.ERC20ZRC20Addr.Hex(),
				Rate:  sdk.NewDec(1).QuoInt64(2), // 2 USDC = 1 ZETA
			},
		},
	}

	// these are the amounts for the withdraws for the different types
	// currently these are arbitrary values that can be fine-tuned for manual testing of rate limiter
	// TODO: define more rigorous assertions with proper values
	// https://github.com/zeta-chain/node/issues/2090
	zetaAmount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(3))
	ethAmount := big.NewInt(1e18)
	erc20Amount := big.NewInt(1e6)

	// approve tokens for the tests
	require.NoError(r, approveTokens(r))

	// add liquidity in the pool to prevent high slippage in WZETA/gas pair
	require.NoError(r, addZetaGasLiquidity(r))

	// Set the rate limiter to 0.5ZETA per 10 blocks
	// These rate limiter flags will only allow to process 1 withdraw per 10 blocks
	r.Logger.Info("setting up rate limiter flags")
	require.NoError(r, setupRateLimiterFlags(r, rateLimiterFlags))

	// Test with rate limiter
	// TODO: define proper assertion to check the rate limiter is working
	// https://github.com/zeta-chain/node/issues/2090
	r.Logger.Print("rate limiter enabled")
	require.NoError(r, createAndWaitWithdraws(r, withdrawTypeZETA, zetaAmount))
	require.NoError(r, createAndWaitWithdraws(r, withdrawTypeETH, ethAmount))
	require.NoError(r, createAndWaitWithdraws(r, withdrawTypeERC20, erc20Amount))

	// Disable rate limiter
	r.Logger.Info("disabling rate limiter")
	require.NoError(r, setupRateLimiterFlags(r, crosschaintypes.RateLimiterFlags{Enabled: false}))

	// Test without rate limiter again and try again ZETA withdraws
	r.Logger.Print("rate limiter disabled")
	require.NoError(r, createAndWaitWithdraws(r, withdrawTypeZETA, zetaAmount))
}

// createAndWaitWithdraws performs RateLimiterWithdrawNumber withdraws
func createAndWaitWithdraws(r *runner.E2ERunner, withdrawType withdrawType, withdrawAmount *big.Int) error {
	startTime := time.Now()

	r.Logger.Print("starting %d %s withdraws", rateLimiterWithdrawNumber, withdrawType)

	// Perform RateLimiterWithdrawNumber withdraws to log time for completion
	txs := make([]*ethtypes.Transaction, rateLimiterWithdrawNumber)
	for i := 0; i < rateLimiterWithdrawNumber; i++ {
		// create a new withdraw depending on the type
		switch withdrawType {
		case withdrawTypeZETA:
			txs[i] = r.WithdrawZeta(withdrawAmount, true)
		case withdrawTypeETH:
			txs[i] = r.WithdrawEther(withdrawAmount)
		case withdrawTypeERC20:
			txs[i] = r.WithdrawERC20(withdrawAmount)
		default:
			return fmt.Errorf("invalid withdraw type: %s", withdrawType)
		}
	}

	// start a error group to wait for all the withdraws to be mined
	g, ctx := errgroup.WithContext(r.Ctx)
	for i, tx := range txs {
		// capture the loop variables
		tx, i := tx, i

		// start a goroutine to wait for the withdraw to be mined
		g.Go(func() error {
			return waitForWithdrawMined(ctx, r, tx, i, startTime)
		})
	}

	// wait for all the withdraws to be mined
	if err := g.Wait(); err != nil {
		return err
	}

	duration := time.Now().Sub(startTime).Seconds()
	block, err := r.ZEVMClient.BlockNumber(r.Ctx)
	if err != nil {
		return fmt.Errorf("error getting block number: %w", err)
	}
	r.Logger.Print("all withdraws completed in %vs at block %d", duration, block)

	return nil
}

// waitForWithdrawMined waits for a withdraw to be mined
// we first wait to get the receipt
// NOTE: this could be a more general function but we define it here for this test because we emit in the function logs specific to this test
func waitForWithdrawMined(
	ctx context.Context,
	r *runner.E2ERunner,
	tx *ethtypes.Transaction,
	index int,
	startTime time.Time,
) error {
	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		return fmt.Errorf(
			"expected cctx status to be %s; got %s, message %s",
			crosschaintypes.CctxStatus_OutboundMined,
			cctx.CctxStatus.Status.String(),
			cctx.CctxStatus.StatusMessage,
		)
	}

	// record the time for completion
	duration := time.Now().Sub(startTime).Seconds()
	block, err := r.ZEVMClient.BlockNumber(ctx)
	if err != nil {
		return err
	}
	r.Logger.Print("cctx %d mined in %vs at block %d", index, duration, block)

	return nil
}

// setupRateLimiterFlags sets up the rate limiter flags with flags defined in the test
func setupRateLimiterFlags(r *runner.E2ERunner, flags crosschaintypes.RateLimiterFlags) error {
	adminAddr, err := r.ZetaTxServer.GetAccountAddressFromName(utils.OperationalPolicyName)
	if err != nil {
		return err
	}
	_, err = r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, crosschaintypes.NewMsgUpdateRateLimiterFlags(
		adminAddr,
		flags,
	))
	if err != nil {
		return err
	}

	return nil
}

// addZetaGasLiquidity adds liquidity to the ZETA/gas pool
func addZetaGasLiquidity(r *runner.E2ERunner) error {
	// use 10 ZETA and 10 ETH for the liquidity
	// this will be sufficient for the tests
	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(10))
	approveAmount := big.NewInt(0).Mul(amount, big.NewInt(10))

	// approve uniswap router to spend gas
	txETHZRC20Approve, err := r.ETHZRC20.Approve(r.ZEVMAuth, r.UniswapV2RouterAddr, approveAmount)
	if err != nil {
		return fmt.Errorf("error approving ZETA: %w", err)
	}

	// wait for the tx to be mined
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, txETHZRC20Approve, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		return fmt.Errorf("approve failed")
	}

	// approve uniswap router to spend ZETA
	txZETAApprove, err := r.WZeta.Approve(r.ZEVMAuth, r.UniswapV2RouterAddr, approveAmount)
	if err != nil {
		return fmt.Errorf("error approving ZETA: %w", err)
	}

	// wait for the tx to be mined
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, txZETAApprove, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		return fmt.Errorf("approve failed")
	}

	// add liquidity in the pool to prevent high slippage in WZETA/gas pair
	r.ZEVMAuth.Value = amount
	txAddLiquidity, err := r.UniswapV2Router.AddLiquidityETH(
		r.ZEVMAuth,
		r.ETHZRC20Addr,
		amount,
		big.NewInt(1e18),
		big.NewInt(1e18),
		r.EVMAddress(),
		big.NewInt(time.Now().Add(10*time.Minute).Unix()),
	)
	if err != nil {
		return fmt.Errorf("error adding liquidity: %w", err)
	}
	r.ZEVMAuth.Value = big.NewInt(0)

	// wait for the tx to be mined
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, txAddLiquidity, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		return fmt.Errorf("add liquidity failed")
	}

	return nil
}

// approveTokens approves the tokens for the tests
func approveTokens(r *runner.E2ERunner) error {
	// deposit and approve 50 WZETA for the tests
	approveAmount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(50))
	r.DepositAndApproveWZeta(approveAmount)

	// approve ETH for withdraws
	tx, err := r.ETHZRC20.Approve(r.ZEVMAuth, r.ETHZRC20Addr, approveAmount)
	if err != nil {
		return fmt.Errorf("error approving ETH: %w", err)
	}
	r.Logger.EVMTransaction(*tx, "approve")

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		return fmt.Errorf("eth approve failed")
	}
	r.Logger.EVMReceipt(*receipt, "approve")

	// approve ETH for ERC20 withdraws (this is for the gas fees)
	tx, err = r.ETHZRC20.Approve(r.ZEVMAuth, r.ERC20ZRC20Addr, approveAmount)
	if err != nil {
		return fmt.Errorf("error approving ERC20: %w", err)
	}

	r.Logger.EVMTransaction(*tx, "approve")

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		return fmt.Errorf("erc 20 approve failed")
	}
	r.Logger.EVMReceipt(*receipt, "approve")

	return nil
}
