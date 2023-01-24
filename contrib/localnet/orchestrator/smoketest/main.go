package main

import (
	"context"
	"fmt"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeta-chain/zetacore/contracts/evm/zetaconnectoreth"
	"github.com/zeta-chain/zetacore/contracts/evm/zetaeth"
	"github.com/zeta-chain/zetacore/contracts/zevm"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"google.golang.org/grpc"
)

var (
	DeployerAddress    = ethcommon.HexToAddress("0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC")
	DeployerPrivateKey = "d87baf7bf6dc560a252596678c12e41f7d1682837f05b29d411bc3f78ae2c263"
	TSSAddress         = ethcommon.HexToAddress("0xF421292cb0d3c97b90EEEADfcD660B893592c6A2")
	BLOCK              = 6 * time.Second // should be 2x block time
	BigZero            = big.NewInt(0)
	SmokeTestTimeout   = 10 * time.Minute // smoke test fails if timeout is reached
	USDTZRC20Addr      = "0x0cbe0dF132a6c6B4a2974Fa1b7Fb953CF0Cc798a"
	USDTERC20Addr      = "0x92339c9Cf464c96E63A4104f3cb97ca336Ea4cE1"
	ERC20CustodyAddr   = "0x0e141A7e7C0A7E15E7d22713Fc0a6187515Fa9BF"
)

func main() {
	testStartTime := time.Now()
	defer func() {
		fmt.Println("Smoke test took", time.Since(testStartTime))
	}()
	go func() {
		time.Sleep(SmokeTestTimeout)
		fmt.Println("Smoke test timed out after", SmokeTestTimeout)
		os.Exit(1)
	}()
	goerliClient, err := ethclient.Dial("http://eth:8545")
	if err != nil {
		panic(err)
	}
	bn, err := goerliClient.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}
	chainID, err := goerliClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("ChainID: %d, Current block number: %d\n", chainID, bn)
	bal, err := goerliClient.BalanceAt(context.TODO(), DeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deployer address: %s, balance: %d Ether\n", DeployerAddress.Hex(), bal.Div(bal, big.NewInt(1e18)))

	// The following deployment must happen here and in this order, please do not change
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
	nonce, err := goerliClient.PendingNonceAt(context.Background(), DeployerAddress)
	if err != nil {
		panic(err)
	}
	if nonce != 0 {
		panic(fmt.Sprintf("nonce of deployer address should be 0, but got %d", nonce))
	}
	fmt.Printf("Step 1: Deploying ZetaEth contract\n")
	zetaEthAddr, tx, ZetaEth, err := zetaeth.DeployZetaEth(auth, goerliClient, big.NewInt(21_000_000_000))
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZetaEth contract address: %s, tx hash: %s\n", zetaEthAddr.Hex(), tx.Hash().Hex())
	time.Sleep(BLOCK)
	receipt, err := goerliClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZetaEth contract receipt: contract address %s, status %d\n", receipt.ContractAddress, receipt.Status)
	bal2, err := ZetaEth.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deployer address: %s, balance: %d ZetaEth\n", DeployerAddress.Hex(), bal2.Div(bal2, big.NewInt(1e18)))
	connectorEthAddr, tx, ConnectorEth, err := zetaconnectoreth.DeployZetaConnectorEth(auth, goerliClient, zetaEthAddr,
		TSSAddress, DeployerAddress, DeployerAddress)
	if err != nil {
		panic(err)
	}
	time.Sleep(BLOCK)
	receipt, err = goerliClient.TransactionReceipt(context.Background(), tx.Hash())
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
	tx, err = ConnectorEth.Send(auth, zetaconnectoreth.ZetaInterfacesSendInput{
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
	receipt, err = goerliClient.TransactionReceipt(context.Background(), tx.Hash())
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
	if err != nil {
		panic(err)
	}
	cctxClient := types.NewQueryClient(grpcConn)
	fungibleClient := fungibletypes.NewQueryClient(grpcConn)

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
	//wg.Wait() // allow the tests to run in parallel

	// ==================== Sending ZETA to ZetaChain ===================
	amount = big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(100)) // 100 Zeta
	fmt.Printf("Step 3: Sending ZETA to ZetaChain\n")
	tx, err = ZetaEth.Approve(auth, connectorEthAddr, amount)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Approve tx hash: %s\n", tx.Hash().Hex())
	time.Sleep(BLOCK)
	tx, err = ConnectorEth.Send(auth, zetaconnectoreth.ZetaInterfacesSendInput{
		DestinationChainId:  big.NewInt(101), // in dev mode, 101 is the  zEVM ChainID
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
	receipt, err = goerliClient.TransactionReceipt(context.Background(), tx.Hash())
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
			fmt.Printf("    Block Num: %d\n", log.BlockNumber)
		}
	}

	zevmClient, err := ethclient.Dial("http://zetacore0:8545")
	if err != nil {
		panic(err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			time.Sleep(5 * time.Second)
			bn, _ := zevmClient.BlockNumber(context.Background())
			bal, _ := zevmClient.BalanceAt(context.Background(), DeployerAddress, big.NewInt(int64(bn)))
			fmt.Printf("Zeta block %d, Deployer Zeta balance: %d\n", bn, bal)

			if bal.Int64() > 0 {
				fmt.Printf("Positive zeta balance; success!\n")
				break
			}
		}
	}()
	wg.Wait()
	TestERC20Deposit(goerliClient, zevmClient, cctxClient, fungibleClient)
	TestDepositEtherIntoZRC20(goerliClient, zevmClient, cctxClient, fungibleClient)
	TestERC20Withdraw(goerliClient, zevmClient, cctxClient, fungibleClient)
	// ==================== Add your tests here ====================

	// ==================== Sending ZETA to Ethereum ===================
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
	zauth, err := bind.NewKeyedTransactorWithChainID(deployerPrivkey, zchainid)

	zauth.Value = big.NewInt(1e18)
	tx, err = wzeta.Deposit(zauth)
	zauth.Value = BigZero
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deposit tx hash: %s\n", tx.Hash().Hex())
	time.Sleep(12 * time.Second)
	receipt, err = zevmClient.TransactionReceipt(context.Background(), tx.Hash())
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
}

// wait until cctx is mined; returns the cctxIndex
func WaitCctxMinedByInTxHash(inTxHash string, cctxClient types.QueryClient) *types.CrossChainTx {
	var cctxIndex string
	for {
		time.Sleep(5 * time.Second)
		res, err := cctxClient.InTxHashToCctx(context.Background(), &types.QueryGetInTxHashToCctxRequest{InTxHash: inTxHash})
		if err != nil {
			continue
		}
		cctxIndex = res.InTxHashToCctx.CctxIndex
		fmt.Printf("Deposit receipt cctx index: %s\n", cctxIndex)
		break
	}
	for {
		time.Sleep(5 * time.Second)
		res, err := cctxClient.Cctx(context.Background(), &types.QueryGetCctxRequest{Index: cctxIndex})
		if err != nil || res.CrossChainTx.CctxStatus.Status != types.CctxStatus_OutboundMined {
			fmt.Printf("Deposit receipt cctx status: %s\n", res.CrossChainTx.CctxStatus.Status.String())
			continue
		}
		fmt.Printf("Deposit receipt cctx status: %+v; success\n", res.CrossChainTx.CctxStatus.Status.String())
		return res.CrossChainTx
	}

}
