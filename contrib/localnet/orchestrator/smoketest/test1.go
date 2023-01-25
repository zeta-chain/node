package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	contracts "github.com/zeta-chain/zetacore/contracts/zevm"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	"math/big"
	"time"
)

// this tests sending ZETA out of ZetaChain to Ethereum
func test1(zevmClient *ethclient.Client, goerliClient *ethclient.Client, cctxClient types.QueryClient, fungibleClient fungibletypes.QueryClient) {

	fmt.Printf("======= test1 =================")
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

	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), deployerPrivkey)
	err = goerliClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("GOERLI tx sent: %s\n", signedTx.Hash().String())
	time.Sleep(BLOCK)
	receipt, err := goerliClient.TransactionReceipt(context.Background(), signedTx.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Printf("GOERLI tx receipt: %d\n", receipt.Status)
	time.Sleep(BLOCK)
	cctxIndex := ""
	for {
		time.Sleep(5 * time.Second)
		res, err := cctxClient.InTxHashToCctx(context.Background(), &types.QueryGetInTxHashToCctxRequest{InTxHash: signedTx.Hash().String()})
		if err != nil {
			fmt.Printf("waiting for cctx from intxhash %s\n", tx.Hash().String())
			continue
		}
		fmt.Printf("cctx found: %s\n", res.InTxHashToCctx.CctxIndex)
		cctxIndex = res.InTxHashToCctx.CctxIndex
		break
	}
	for {
		time.Sleep(5 * time.Second)
		res, err := cctxClient.Cctx(context.Background(), &types.QueryGetCctxRequest{Index: cctxIndex})
		if err != nil {
			fmt.Printf("waiting for cctx %s: status %s\n", cctxIndex, res.CrossChainTx.CctxStatus.Status)
			continue
		}
		fmt.Printf("cctx found: %s\n", res.CrossChainTx.CctxStatus)
		if res.CrossChainTx.CctxStatus.Status == types.CctxStatus_OutboundMined {
			fmt.Printf("cctx %s is mined\n", cctxIndex)
			break
		}
	}
	systemContractAddr, err := fungibleClient.SystemContract(context.Background(), &fungibletypes.QueryGetSystemContractRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("system contract address: %s\n", systemContractAddr.SystemContract.SystemContract)
	addr := ethcommon.HexToAddress(systemContractAddr.SystemContract.SystemContract)
	systemContract, err := contracts.NewSystemContract(addr, zevmClient)
	if err != nil {
		panic(err)
	}
	ethZRC20Addr, err := systemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(5))
	if err != nil {
		panic(err)
	}
	fmt.Printf("eth zrc20 address: %s\n", ethZRC20Addr.String())
	ethZRC20, err := contracts.NewZRC20(ethZRC20Addr, zevmClient)
	if err != nil {
		panic(err)
	}
	ethZRC20Balance, err := ethZRC20.BalanceOf(nil, DeployerAddress)
	if err != nil {
		panic(err)
	}
	fmt.Printf("eth zrc20 balance: %s\n", ethZRC20Balance.String())
}
