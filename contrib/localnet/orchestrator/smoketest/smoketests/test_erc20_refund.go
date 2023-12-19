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
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	utils.LoudPrintf("Deposit a non-gas ZRC20 into ZEVM and call a contract that reverts; should refund on ZetaChain if no liquidity pool, should refund on origin if liquidity pool\n")

	// Get the initial balance of the deployer
	initialBal, err := sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}

	fmt.Println("Sending a deposit that should revert without a liquidity pool makes the cctx aborted")

	amount := big.NewInt(1e4)

	// send the deposit
	inTxHash, err := sendInvalidUSDTDeposit(sm, amount)
	if err != nil {
		panic(err)
	}

	// There is no liquidity pool, therefore the cctx should abort
	cctx := utils.WaitCctxMinedByInTxHash(inTxHash, sm.CctxClient)
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
	fmt.Println("CCTX has been aborted and the erc20 has been refunded on ZetaChain")

	amount = big.NewInt(1e7)
	goerliBalance, err := sm.USDTERC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}

	fmt.Println("Sending a deposit that should revert with a liquidity pool")

	fmt.Println("Creating the liquidity pool USTD/ZETA")
	err = createZetaERC20LiquidityPool(sm)
	if err != nil {
		panic(err)
	}
	fmt.Println("Liquidity pool created")

	// send the deposit
	inTxHash, err = sendInvalidUSDTDeposit(sm, amount)
	if err != nil {
		panic(err)
	}

	// there is a liquidity pool, therefore the cctx should revert
	cctx = utils.WaitCctxMinedByInTxHash(inTxHash, sm.CctxClient)

	// the revert tx creation will fail because the sender, used as the recipient, is not defined in the cctx
	if cctx.CctxStatus.Status != types.CctxStatus_Reverted {
		panic(fmt.Sprintf("expected cctx status to be PendingRevert; got %s", cctx.CctxStatus.Status))
	}

	// get revert tx
	revertTxHash := cctx.GetCurrentOutTxParam().OutboundTxHash
	_, _, err = sm.GoerliClient.TransactionByHash(context.Background(), ethcommon.HexToHash(revertTxHash))
	if err != nil {
		panic(err)
	}
	receipt, err := sm.GoerliClient.TransactionReceipt(context.Background(), ethcommon.HexToHash(revertTxHash))
	if err != nil {
		panic(err)
	}
	if receipt.Status == 0 {
		panic("expected the revert tx receipt to have status 1; got 0")
	}

	// check that the erc20 in the reverted cctx was refunded on Goerli
	newGoerliBalance, err := sm.USDTERC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	// the new balance must be higher than the previous one because of the revert refund
	if goerliBalance.Cmp(newGoerliBalance) != -1 {
		panic(fmt.Sprintf("expected balance to be higher than %s after refund; got %s", goerliBalance.String(), newGoerliBalance.String()))
	}
	// it must also be lower than the previous balance + the amount because of the gas fee for the revert tx
	balancePlusAmount := goerliBalance.Add(goerliBalance, amount)
	if newGoerliBalance.Cmp(balancePlusAmount) != -1 {
		panic(fmt.Sprintf("expected balance to be lower than %s after refund; got %s", balancePlusAmount.String(), newGoerliBalance.String()))
	}

	fmt.Println("ERC20 CCTX successfully reverted")
	fmt.Println("\tbalance before refund: ", goerliBalance.String())
	fmt.Println("\tamount: ", amount.String())
	fmt.Println("\tbalance after refund: ", newGoerliBalance.String())
}

func createZetaERC20LiquidityPool(sm *runner.SmokeTestRunner) error {
	amount := big.NewInt(1e10)
	txHash := sm.DepositERC20(amount, []byte{})
	utils.WaitCctxMinedByInTxHash(txHash.Hex(), sm.CctxClient)

	tx, err := sm.USDTZRC20.Approve(sm.ZevmAuth, sm.UniswapV2RouterAddr, big.NewInt(1e10))
	if err != nil {
		return err
	}
	receipt := utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
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
	receipt = utils.MustWaitForTxReceipt(sm.ZevmClient, tx)
	if receipt.Status == 0 {
		return fmt.Errorf("add liquidity failed")
	}

	return nil
}

func sendInvalidUSDTDeposit(sm *runner.SmokeTestRunner, amount *big.Int) (string, error) {
	// send the tx
	USDT := sm.USDTERC20
	tx, err := USDT.Mint(sm.GoerliAuth, amount)
	if err != nil {
		return "", err
	}
	receipt := utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
	fmt.Printf("Mint receipt tx hash: %s\n", tx.Hash().Hex())

	tx, err = USDT.Approve(sm.GoerliAuth, sm.ERC20CustodyAddr, amount)
	if err != nil {
		return "", err
	}
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
	fmt.Printf("USDT Approve receipt tx hash: %s\n", tx.Hash().Hex())

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

	fmt.Printf("GOERLI tx sent: %s; to %s, nonce %d\n", tx.Hash().String(), tx.To().Hex(), tx.Nonce())
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
	if receipt.Status == 0 {
		return "", errors.New("expected the tx receipt to have status 1; got 0")
	}
	fmt.Printf("GOERLI tx receipt: %d\n", receipt.Status)
	fmt.Printf("  tx hash: %s\n", receipt.TxHash.String())
	fmt.Printf("  to: %s\n", tx.To().String())
	fmt.Printf("  value: %d\n", tx.Value())
	fmt.Printf("  block num: %d\n", receipt.BlockNumber)

	return tx.Hash().Hex(), nil
}
