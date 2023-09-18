//go:build PRIVNET
// +build PRIVNET

package main

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func (sm *SmokeTest) TestERC20DepositAndCallRefund() {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	LoudPrintf("Deposit a non-gas ZRC20 into ZEVM and call a contract that reverts; should refund on ZetaChain if no liquidity pool, should refund on origin if liquidity pool\n")

	// TODO: uncomment once erc20 deposit contains the sender address
	//https://github.com/zeta-chain/node/issues/1140
	//// Get the initial balance of the deployer
	//initialBal, err := sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	//if err != nil {
	//	panic(err)
	//}

	fmt.Println("Sending a deposit that should revert without a liquidity pool makes the cctx aborted")

	amount := big.NewInt(1e7)

	// send the deposit
	inTxHash, err := sm.sendInvalidUSDTDeposit(amount)
	if err != nil {
		panic(err)
	}

	// There is no liquidity pool, therefore the cctx should abort
	cctx := WaitCctxMinedByInTxHash(inTxHash, sm.cctxClient)
	if cctx.CctxStatus.Status != types.CctxStatus_Aborted {
		panic(fmt.Sprintf("expected cctx status to be Aborted; got %s", cctx.CctxStatus.Status))
	}

	// TODO: uncomment once erc20 deposit contains the sender address
	//https://github.com/zeta-chain/node/issues/1140
	//// Check that the erc20 in the aborted cctx was refunded on ZetaChain
	//newBalance, err := sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	//if err != nil {
	//	panic(err)
	//}
	//expectedBalance := initialBal.Add(initialBal, amount)
	//if newBalance.Cmp(expectedBalance) != 0 {
	//	panic(fmt.Sprintf("expected balance to be %s after refund; got %s", expectedBalance.String(), newBalance.String()))
	//}

	fmt.Println("Sending a deposit that should revert with a liquidity pool")

	fmt.Println("Creating the liquidity pool USTD/ZETA")
	err = sm.createZetaERC20LiquidityPool()
	if err != nil {
		panic(err)
	}
	fmt.Println("Liquidity pool created")

	// send the deposit
	inTxHash, err = sm.sendInvalidUSDTDeposit(amount)
	if err != nil {
		panic(err)
	}

	// There is a liquidity pool, therefore the cctx should revert
	cctx = WaitCctxMinedByInTxHash(inTxHash, sm.cctxClient)

	// The revert tx creation will fail because the sender, used as the recipient, is not defined in the cctx
	// TODO: uncomment once erc20 deposit contains the sender address
	//https://github.com/zeta-chain/node/issues/1140
	//if cctx.CctxStatus.Status != types.CctxStatus_Reverted {
	//	panic(fmt.Sprintf("expected cctx status to be PendingRevert; got %s", cctx.CctxStatus.Status))
	//}
	//
	//// get revert tx
	//revertTxHash := cctx.GetCurrentOutTxParam().OutboundTxHash
	//tx, _, err := sm.goerliClient.TransactionByHash(context.Background(), ethcommon.HexToHash(revertTxHash))
	//if err != nil {
	//	panic(err)
	//}
	//receipt, err := sm.goerliClient.TransactionReceipt(context.Background(), ethcommon.HexToHash(revertTxHash))
	//if err != nil {
	//	panic(err)
	//}
	//if receipt.Status == 0 {
	//	panic("expected the revert tx receipt to have status 1; got 0")
	//}
	//
	//if *tx.To() != DeployerAddress {
	//	panic(fmt.Sprintf("expected tx to %s; got %s", DeployerAddress.Hex(), tx.To().Hex()))
	//}
}

func (sm *SmokeTest) createZetaERC20LiquidityPool() error {
	amount := big.NewInt(1e10)
	txHash := sm.DepositERC20(amount, []byte{})
	WaitCctxMinedByInTxHash(txHash.Hex(), sm.cctxClient)

	tx, err := sm.USDTZRC20.Approve(sm.zevmAuth, sm.UniswapV2RouterAddr, big.NewInt(1e10))
	if err != nil {
		return err
	}
	receipt := MustWaitForTxReceipt(sm.zevmClient, tx)
	if receipt.Status == 0 {
		return errors.New("approve failed")
	}

	previousValue := sm.zevmAuth.Value
	sm.zevmAuth.Value = big.NewInt(1e10)
	tx, err = sm.UniswapV2Router.AddLiquidityETH(
		sm.zevmAuth,
		sm.USDTZRC20Addr,
		amount,
		BigZero,
		BigZero,
		DeployerAddress,
		big.NewInt(time.Now().Add(10*time.Minute).Unix()),
	)
	sm.zevmAuth.Value = previousValue
	if err != nil {
		return err
	}
	receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	if receipt.Status == 0 {
		return errors.New("add liquidity failed")
	}

	return nil
}

func (sm *SmokeTest) sendInvalidUSDTDeposit(amount *big.Int) (string, error) {
	// send the tx
	USDT := sm.USDTERC20
	tx, err := USDT.Mint(sm.goerliAuth, amount)
	if err != nil {
		return "", err
	}
	receipt := MustWaitForTxReceipt(sm.goerliClient, tx)
	fmt.Printf("Mint receipt tx hash: %s\n", tx.Hash().Hex())

	tx, err = USDT.Approve(sm.goerliAuth, sm.ERC20CustodyAddr, amount)
	if err != nil {
		return "", err
	}
	receipt = MustWaitForTxReceipt(sm.goerliClient, tx)
	fmt.Printf("USDT Approve receipt tx hash: %s\n", tx.Hash().Hex())

	tx, err = sm.ERC20Custody.Deposit(
		sm.goerliAuth,
		DeployerAddress.Bytes(),
		sm.USDTERC20Addr,
		amount,
		[]byte("this is an invalid msg that will cause the contract to revert"),
	)
	if err != nil {
		return "", err
	}

	fmt.Printf("GOERLI tx sent: %s; to %s, nonce %d\n", tx.Hash().String(), tx.To().Hex(), tx.Nonce())
	receipt = MustWaitForTxReceipt(sm.goerliClient, tx)
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
