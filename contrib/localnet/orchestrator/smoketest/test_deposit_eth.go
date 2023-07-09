//go:build PRIVNET
// +build PRIVNET

package main

import (
	"context"
	"fmt"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	zrc20 "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	"github.com/zeta-chain/zetacore/common"
)

// this tests sending ZETA out of ZetaChain to Ethereum
func (sm *SmokeTest) TestDepositEtherIntoZRC20() {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	goerliClient := sm.goerliClient
	LoudPrintf("Deposit Ether into ZEVM\n")
	bn, err := goerliClient.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("GOERLI block number: %d\n", bn)
	bal, err := goerliClient.BalanceAt(context.Background(), DeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("GOERLI deployer balance: %s\n", bal.String())
	nonce, err := goerliClient.PendingNonceAt(context.Background(), DeployerAddress)
	if err != nil {
		panic(err)
	}
	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	gasLimit := uint64(21000)                // in units
	gasPrice, err := goerliClient.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}
	tx := ethtypes.NewTransaction(nonce, TSSAddress, value, gasLimit, gasPrice, nil)
	chainID, err := goerliClient.NetworkID(context.Background())
	if err != nil {
		panic(err)
	}
	deployerPrivkey, err := crypto.HexToECDSA(DeployerPrivateKey)
	if err != nil {
		panic(err)
	}

	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), deployerPrivkey)
	if err != nil {
		panic(err)
	}
	err = goerliClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("GOERLI tx sent: %s; to %s, nonce %d\n", signedTx.Hash().String(), signedTx.To().Hex(), signedTx.Nonce())
	receipt := MustWaitForTxReceipt(sm.goerliClient, signedTx)
	fmt.Printf("GOERLI tx receipt: %d\n", receipt.Status)
	fmt.Printf("  tx hash: %s\n", receipt.TxHash.String())
	fmt.Printf("  to: %s\n", signedTx.To().String())
	fmt.Printf("  value: %d\n", signedTx.Value())
	fmt.Printf("  block num: %d\n", receipt.BlockNumber)

	c := make(chan any)
	sm.wg.Add(1)
	go func() {
		defer sm.wg.Done()
		WaitCctxMinedByInTxHash(signedTx.Hash().Hex(), sm.cctxClient)
		c <- 0
	}()
	sm.wg.Add(1)
	go func() {
		defer sm.wg.Done()
		<-c

		systemContract := sm.SystemContract
		if err != nil {
			panic(err)
		}
		ethZRC20Addr, err := systemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(common.GoerliChain().ChainId))
		if err != nil {
			panic(err)
		}
		sm.ETHZRC20Addr = ethZRC20Addr
		fmt.Printf("eth zrc20 address: %s\n", ethZRC20Addr.String())
		ethZRC20, err := zrc20.NewZRC20(ethZRC20Addr, sm.zevmClient)
		if err != nil {
			panic(err)
		}
		sm.ETHZRC20 = ethZRC20
		ethZRC20Balance, err := ethZRC20.BalanceOf(nil, DeployerAddress)
		if err != nil {
			panic(err)
		}
		fmt.Printf("eth zrc20 balance: %s\n", ethZRC20Balance.String())
		if ethZRC20Balance.Cmp(value) != 0 {
			fmt.Printf("eth zrc20 bal wanted %d, got %d\n", value, ethZRC20Balance)
			panic("bal mismatch")
		}
	}()
	sm.wg.Wait()
}

func (sm *SmokeTest) TestDepositAndCallRefund() {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	LoudPrintf("Deposit ZRC20 into ZEVM and call a contract that reverts; should refund\n")

	goerliClient := sm.goerliClient
	bn, err := goerliClient.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("GOERLI block number: %d\n", bn)
	bal, err := goerliClient.BalanceAt(context.Background(), DeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("GOERLI deployer balance: %s\n", bal.String())
	nonce, err := goerliClient.PendingNonceAt(context.Background(), DeployerAddress)
	if err != nil {
		panic(err)
	}
	value := big.NewInt(100000000000000000) // in wei (1 eth)
	gasLimit := uint64(23000)               // in units
	gasPrice, err := goerliClient.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}
	data := append(sm.BTCZRC20Addr.Bytes(), []byte("hello sailors")...) // this data
	tx := ethtypes.NewTransaction(nonce, TSSAddress, value, gasLimit, gasPrice, data)
	chainID, err := goerliClient.NetworkID(context.Background())
	if err != nil {
		panic(err)
	}
	deployerPrivkey, err := crypto.HexToECDSA(DeployerPrivateKey)
	if err != nil {
		panic(err)
	}

	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), deployerPrivkey)
	if err != nil {
		panic(err)
	}
	err = goerliClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("GOERLI tx sent: %s; to %s, nonce %d\n", signedTx.Hash().String(), signedTx.To().Hex(), signedTx.Nonce())
	receipt := MustWaitForTxReceipt(sm.goerliClient, signedTx)
	fmt.Printf("GOERLI tx receipt: %d\n", receipt.Status)
	fmt.Printf("  tx hash: %s\n", receipt.TxHash.String())
	fmt.Printf("  to: %s\n", signedTx.To().String())
	fmt.Printf("  value: %d\n", signedTx.Value())
	fmt.Printf("  block num: %d\n", receipt.BlockNumber)

	c := make(chan any)
	sm.wg.Add(1)
	go func() {
		defer sm.wg.Done()
		cctx := WaitCctxMinedByInTxHash(signedTx.Hash().Hex(), sm.cctxClient)
		fmt.Printf("cctx status message: %s", cctx.CctxStatus.StatusMessage)
		if cctx.CctxStatus.Status != types.CctxStatus_Reverted {
			panic(fmt.Sprintf("expected cctx status PendingRevert; got %s", cctx.CctxStatus.Status))
		}
		c <- 0
	}()
	sm.wg.Add(1)
	go func() {
		defer sm.wg.Done()
		<-c

		systemContract := sm.SystemContract
		if err != nil {
			panic(err)
		}
		ethZRC20Addr, err := systemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(common.GoerliChain().ChainId))
		if err != nil {
			panic(err)
		}
		sm.ETHZRC20Addr = ethZRC20Addr
		fmt.Printf("eth zrc20 address: %s\n", ethZRC20Addr.String())
		ethZRC20, err := zrc20.NewZRC20(ethZRC20Addr, sm.zevmClient)
		if err != nil {
			panic(err)
		}
		sm.ETHZRC20 = ethZRC20
		ethZRC20Balance, err := ethZRC20.BalanceOf(nil, DeployerAddress)
		if err != nil {
			panic(err)
		}
		fmt.Printf("eth zrc20 balance: %s\n", ethZRC20Balance.String())
		if ethZRC20Balance.Cmp(value) != 0 {
			fmt.Printf("eth zrc20 bal wanted %d, got %d\n", value, ethZRC20Balance)
			panic("bal mismatch")
		}
	}()
	sm.wg.Wait()
}
