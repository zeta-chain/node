package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/zetacore/contracts/evm/erc20"
	contracts "github.com/zeta-chain/zetacore/contracts/zevm"
	"math/big"
	"time"
)

func (sm *SmokeTest) TestERC20Withdraw() {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	LoudPrintf("Withdraw USDT ZRC20\n")
	zevmClient := sm.zevmClient
	goerliClient := sm.goerliClient
	cctxClient := sm.cctxClient

	usdtZRC20, err := contracts.NewZRC20(ethcommon.HexToAddress(USDTZRC20Addr), zevmClient)
	if err != nil {
		panic(err)
	}
	bal, err := usdtZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	if err != nil {
		panic(err)
	}
	fmt.Printf("balance of deployer on USDT ZRC20: %d\n", bal)
	supply, err := usdtZRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("supply of USDT ZRC20: %d\n", supply)
	if bal.Int64() != 1e6 {
		panic("balance is not correct")
	}

	gasZRC20, gasFee, err := usdtZRC20.WithdrawGasFee(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("gasZRC20: %s, gasFee: %d\n", gasZRC20.Hex(), gasFee)

	ethZRC20, err := contracts.NewZRC20(gasZRC20, zevmClient)
	if err != nil {
		panic(err)
	}
	bal, err = ethZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	if err != nil {
		panic(err)
	}
	fmt.Printf("balance of deployer on ETH ZRC20: %d\n", bal)
	if bal.Int64() <= 0 {
		panic("not enough ETH ZRC20 balance!")
	}

	chainID, err := zevmClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	deployerPrivkey, err := crypto.HexToECDSA(DeployerPrivateKey)
	if err != nil {
		panic(err)
	}
	zauth, err := bind.NewKeyedTransactorWithChainID(deployerPrivkey, chainID)
	if err != nil {
		panic("zauth error")
	}
	// Approve
	tx, err := ethZRC20.Approve(zauth, ethcommon.HexToAddress(USDTZRC20Addr), big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	time.Sleep(BLOCK)
	receipt, err := zevmClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Printf("eth zrc20 approve receipt: status %d\n", receipt.Status)
	// Withdraw
	tx, err = usdtZRC20.Withdraw(zauth, DeployerAddress.Bytes(), big.NewInt(100))
	if err != nil {
		panic(err)
	}
	time.Sleep(BLOCK)
	receipt, err = zevmClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)
	for _, log := range receipt.Logs {
		event, err := usdtZRC20.ParseWithdrawal(*log)
		if err != nil {
			continue
		}
		fmt.Printf("  logs: from %s, to %x, value %d, gasfee %d\n", event.From.Hex(), event.To, event.Value, event.Gasfee)
	}

	sm.wg.Add(1)
	go func() {
		defer sm.wg.Done()
		cctx := WaitCctxMinedByInTxHash(receipt.TxHash.Hex(), cctxClient)
		fmt.Printf("outTx hash %s\n", cctx.OutboundTxParams.OutboundTxHash)

		USDTERC20, err := erc20.NewUSDT(ethcommon.HexToAddress(USDTERC20Addr), goerliClient)
		if err != nil {
			panic(err)
		}
		bal, err = USDTERC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
		if err != nil {
			panic(err)
		}
		fmt.Printf("USDT ERC20 bal: %d\n", bal)
	}()
	sm.wg.Wait()
}
