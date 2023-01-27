package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeta-chain/zetacore/contracts/zevm"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	"math/big"
	"sync"
	"time"
)

func TestSendZetaOut(zevmClient *ethclient.Client, goerliClient *ethclient.Client, cctxClient types.QueryClient, fungibleClient fungibletypes.QueryClient) {
	LoudPrintf("Step 4: Sending ZETA from ZEVM to Ethereum\n")
	ConnectorZEVMAddr := ethcommon.HexToAddress("0x239e96c8f17C85c30100AC26F635Ea15f23E9c67")
	ConnectorZEVM, err := zevm.NewZetaConnectorZEVM(ConnectorZEVMAddr, zevmClient)
	if err != nil {
		panic(err)
	}
	//SystemContractAddr := ethcommon.HexToAddress("0x91d18e54DAf4F677cB28167158d6dd21F6aB3921")
	wzetaAddr := ethcommon.HexToAddress("0x5F0b1a82749cb4E2278EC87F8BF6B618dC71a8bf")
	wzeta, err := zevm.NewWZETA(wzetaAddr, zevmClient)
	if err != nil {
		panic(err)
	}
	zchainid, err := zevmClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("zevm chainid: %d\n", zchainid)
	deployerPrivkey, err := crypto.HexToECDSA(DeployerPrivateKey)
	if err != nil {
		panic(err)
	}
	zauth, err := bind.NewKeyedTransactorWithChainID(deployerPrivkey, zchainid)
	if err != nil {
		panic(err)
	}
	zauth.Value = big.NewInt(1e18)
	tx, err := wzeta.Deposit(zauth)
	zauth.Value = BigZero
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deposit tx hash: %s\n", tx.Hash().Hex())
	time.Sleep(12 * time.Second)
	receipt, err := zevmClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deposit tx receipt: status %d\n", receipt.Status)
	tx, err = wzeta.Approve(zauth, ConnectorZEVMAddr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	fmt.Printf("wzeta.approve tx hash: %s\n", tx.Hash().Hex())
	time.Sleep(12 * time.Second)
	receipt, err = zevmClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Printf("approve tx receipt: status %d\n", receipt.Status)
	tx, err = ConnectorZEVM.Send(zauth, zevm.ZetaInterfacesSendInput{
		DestinationChainId:  big.NewInt(1337),
		DestinationAddress:  DeployerAddress.Bytes(),
		DestinationGasLimit: big.NewInt(250_000),
		Message:             nil,
		ZetaValueAndGas:     big.NewInt(1e17),
		ZetaParams:          nil,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("send tx hash: %s\n", tx.Hash().Hex())
	time.Sleep(12 * time.Second)
	receipt, err = zevmClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Printf("send tx receipt: status %d\n", receipt.Status)
	fmt.Printf("  Logs:\n")
	for _, log := range receipt.Logs {
		sentLog, err := ConnectorZEVM.ParseZetaSent(*log)
		if err == nil {
			fmt.Printf("    Dest Addr: %s\n", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
			fmt.Printf("    Dest Chain: %d\n", sentLog.DestinationChainId)
			fmt.Printf("    Dest Gas: %d\n", sentLog.DestinationGasLimit)
			fmt.Printf("    Zeta Value: %d\n", sentLog.ZetaValueAndGas)
		}
	}
	fmt.Printf("waiting for cctx status to change to final...\n")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		var index string
		for {
			time.Sleep(5 * time.Second)
			res, err := cctxClient.InTxHashToCctx(context.Background(), &types.QueryGetInTxHashToCctxRequest{
				InTxHash: tx.Hash().Hex(),
			})
			if err != nil {
				fmt.Printf("No CCTX found for inTxHash %s\n", tx.Hash().Hex())
				continue
			}
			index = res.InTxHashToCctx.CctxIndex
			fmt.Printf("Found CCTX for inTxHash %s: %s\n", tx.Hash().Hex(), index)
			break
		}
		for {
			time.Sleep(5 * time.Second)
			res, err := cctxClient.Cctx(context.Background(), &types.QueryGetCctxRequest{
				Index: index,
			})
			if err != nil {
				fmt.Printf("No CCTX found for index %s\n", index)
				continue
			}
			if res.CrossChainTx.CctxStatus.Status != types.CctxStatus_OutboundMined {
				fmt.Printf("Found CCTX for index %s: status %s\n", index, res.CrossChainTx.CctxStatus.Status)
				continue
			}
			if res.CrossChainTx.CctxStatus.Status == types.CctxStatus_OutboundMined {
				fmt.Printf("Found CCTX for index %s: status %s; success\n", index, res.CrossChainTx.CctxStatus.Status)
				break
			}
		}
	}()
	wg.Wait()
}
