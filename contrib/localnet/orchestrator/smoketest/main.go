package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeta-chain/zetacore/contracts/evm/ZetaConnectorEth"
	"github.com/zeta-chain/zetacore/contracts/evm/ZetaEth"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"google.golang.org/grpc"
	"math/big"
	"time"
)

var (
	DeployerAddress    = ethcommon.HexToAddress("0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC")
	DeployerPrivateKey = "d87baf7bf6dc560a252596678c12e41f7d1682837f05b29d411bc3f78ae2c263"
	TSSAddress         = ethcommon.HexToAddress("0xF421292cb0d3c97b90EEEADfcD660B893592c6A2")
	BLOCK              = 6 * time.Second // should be 2x block time
)

func main() {
	ethclient, err := ethclient.Dial("http://eth:8545")
	if err != nil {
		panic(err)
	}
	bn, err := ethclient.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}
	chainID, err := ethclient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("ChainID: %d, Current block number: %d\n", chainID, bn)
	bal, err := ethclient.BalanceAt(context.TODO(), DeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deployer address: %s, balance: %d Ether\n", DeployerAddress.Hex(), bal.Div(bal, big.NewInt(1e18)))

	// ==================== Deploying contracts ====================
	deployerPrivkey, err := crypto.HexToECDSA(DeployerPrivateKey)
	if err != nil {
		panic(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(deployerPrivkey, chainID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Step 0: Check the nonce of deployer address\n")
	nonce, err := ethclient.PendingNonceAt(context.Background(), DeployerAddress)
	if err != nil {
		panic(err)
	}
	if nonce != 0 {
		panic(fmt.Sprintf("nonce of deployer address should be 0, but got %d", nonce))
	}
	fmt.Printf("Step 1: Deploying ZetaEth contract\n")
	zetaEthAddr, tx, ZetaEth, err := ZetaEth.DeployZetaEth(auth, ethclient, big.NewInt(21_000_000_000))
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZetaEth contract address: %s, tx hash: %s\n", zetaEthAddr.Hex(), tx.Hash().Hex())
	time.Sleep(BLOCK)
	receipt, err := ethclient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZetaEth contract receipt: contract address %s, status %d\n", receipt.ContractAddress, receipt.Status)
	bal2, err := ZetaEth.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deployer address: %s, balance: %d ZetaEth\n", DeployerAddress.Hex(), bal2.Div(bal2, big.NewInt(1e18)))
	connectorEthAddr, tx, ConnectorEth, err := ZetaConnectorEth.DeployZetaConnectorEth(auth, ethclient, zetaEthAddr,
		TSSAddress, DeployerAddress, DeployerAddress)
	if err != nil {
		panic(err)
	}
	time.Sleep(BLOCK)
	receipt, err = ethclient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZetaConnectorEth contract address: %s, tx hash: %s\n", connectorEthAddr.Hex(), tx.Hash().Hex())
	fmt.Printf("ZetaConnectorEth contract receipt: contract address %s, status %d\n", receipt.ContractAddress, receipt.Status)
	_ = ConnectorEth

	// ==================== Interacting with contracts ====================
	time.Sleep(10 * time.Second)
	fmt.Printf("Step 2: Interacting with ZetaEth contract\n")
	fmt.Printf("Approving ConnectorEth to spend deployer's ZetaEth\n")
	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(10)) // 10 Zeta
	tx, err = ZetaEth.Approve(auth, connectorEthAddr, amount)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Approve tx hash: %s\n", tx.Hash().Hex())
	time.Sleep(BLOCK)
	fmt.Printf("Calling ConnectorEth.Send\n")
	tx, err = ConnectorEth.Send(auth, ZetaConnectorEth.ZetaInterfacesSendInput{
		DestinationChainId:  big.NewInt(1337), // in dev mode, GOERLI has chainid 1337
		DestinationAddress:  DeployerAddress.Bytes(),
		DestinationGasLimit: big.NewInt(250_000),
		Message:             nil,
		ZetaValueAndGas:     amount,
		ZetaParams:          nil,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Send tx hash: %s\n", tx.Hash().Hex())
	time.Sleep(BLOCK)
	receipt, err = ethclient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Send tx receipt: status %d\n", receipt.Status)
	fmt.Printf("  Logs:\n")
	for _, log := range receipt.Logs {
		sentLog, err := ConnectorEth.ParseZetaSent(*log)
		if err == nil {
			fmt.Printf("    Dest Addr: %s\n", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
			fmt.Printf("    Dest Chain: %d\n", sentLog.DestinationChainId)
			fmt.Printf("    Dest Gas: %d\n", sentLog.DestinationGasLimit)
			fmt.Printf("    Zeta Value: %d\n", sentLog.ZetaValueAndGas)
		}
	}

	grpcConn, err := grpc.Dial(
		"zetacore0:9090",
		grpc.WithInsecure(),
	)
	cctxClient := types.NewQueryClient(grpcConn)

	for {
		time.Sleep(5 * time.Second)
		res, err := cctxClient.CctxAll(context.Background(), &types.QueryAllCctxRequest{})
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		if len(res.CrossChainTx) != 1 {
			fmt.Printf("Waiting for CrossChainTx to appear...len cctx %d\n", len(res.CrossChainTx))
			continue
		}
		if res.CrossChainTx[0].CctxStatus.Status != types.CctxStatus_OutboundMined {
			fmt.Printf("Waiting for CrossChainTx to be mined...status %s\n", res.CrossChainTx[0].CctxStatus.Status)
			continue
		}
		fmt.Printf("CrossChainTx found: %v\n", res.CrossChainTx[0])
		break
	}
}
