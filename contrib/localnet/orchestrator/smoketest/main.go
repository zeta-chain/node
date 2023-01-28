package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/zetacore/contracts/evm/erc20custody"
	"github.com/zeta-chain/zetacore/contracts/evm/zetaconnectoreth"
	"github.com/zeta-chain/zetacore/contracts/evm/zetaeth"
	contracts "github.com/zeta-chain/zetacore/contracts/zevm"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/erc20"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	"math/big"
	"os"
	"sync"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"google.golang.org/grpc"
)

var (
	DeployerAddress    = ethcommon.HexToAddress("0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC")
	DeployerPrivateKey = "d87baf7bf6dc560a252596678c12e41f7d1682837f05b29d411bc3f78ae2c263"
	TSSAddress         = ethcommon.HexToAddress("0xF421292cb0d3c97b90EEEADfcD660B893592c6A2")
	BLOCK              = 5 * time.Second // should be 2x block time
	BigZero            = big.NewInt(0)
	SmokeTestTimeout   = 10 * time.Minute // smoke test fails if timeout is reached
	USDTZRC20Addr      = "0x7c8dDa80bbBE1254a7aACf3219EBe1481c6E01d7"
	USDTERC20Addr      = "0xff3135df4F2775f4091b81f4c7B6359CfA07862a"
	ERC20CustodyAddr   = "0xD28D6A0b8189305551a0A8bd247a6ECa9CE781Ca"
	HexToAddress       = ethcommon.HexToAddress
)

type SmokeTest struct {
	zevmClient       *ethclient.Client
	goerliClient     *ethclient.Client
	cctxClient       types.QueryClient
	fungibleClient   fungibletypes.QueryClient
	wg               sync.WaitGroup
	ZetaEth          *zetaeth.ZetaEth
	ZetaEthAddr      ethcommon.Address
	ConnectorEth     *zetaconnectoreth.ZetaConnectorEth
	ConnectorEthAddr ethcommon.Address
	goerliAuth       *bind.TransactOpts
	zevmAuth         *bind.TransactOpts

	ERC20CustodyAddr ethcommon.Address
	ERC20Custody     *erc20custody.ERC20Custody
	USDTERC20Addr    ethcommon.Address
	USDTERC20        *erc20.USDT
	USDTZRC20Addr    ethcommon.Address
	USDTZRC20        *contracts.ZRC20
	ETHZRC20Addr     ethcommon.Address
	ETHZRC20         *contracts.ZRC20
}

func NewSmokeTest(goerliClient *ethclient.Client, zevmClient *ethclient.Client,
	cctxClient types.QueryClient, fungibleClient fungibletypes.QueryClient,
	goerliAuth *bind.TransactOpts, zevmAuth *bind.TransactOpts) *SmokeTest {
	return &SmokeTest{
		zevmClient:     zevmClient,
		goerliClient:   goerliClient,
		cctxClient:     cctxClient,
		fungibleClient: fungibleClient,
		wg:             sync.WaitGroup{},
		goerliAuth:     goerliAuth,
		zevmAuth:       zevmAuth,
	}
}

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

	bal, err := goerliClient.BalanceAt(context.TODO(), DeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deployer address: %s, balance: %d Ether\n", DeployerAddress.Hex(), bal.Div(bal, big.NewInt(1e18)))

	chainid, err := goerliClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	deployerPrivkey, err := crypto.HexToECDSA(DeployerPrivateKey)
	if err != nil {
		panic(err)
	}
	goerliAuth, err := bind.NewKeyedTransactorWithChainID(deployerPrivkey, chainid)
	if err != nil {
		panic(err)
	}
	// get the clients for tests
	var zevmClient *ethclient.Client
	for {
		time.Sleep(5 * time.Second)
		fmt.Printf("dialing zevm client: http://zetacore0:8545\n")
		zevmClient, err = ethclient.Dial("http://zetacore0:8545")
		if err != nil {
			continue
		}
		break
	}
	chainid, err = zevmClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	zevmAuth, err := bind.NewKeyedTransactorWithChainID(deployerPrivkey, chainid)
	if err != nil {
		panic(err)
	}
	grpcConn, err := grpc.Dial("zetacore0:9090", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	cctxClient := types.NewQueryClient(grpcConn)
	fungibleClient := fungibletypes.NewQueryClient(grpcConn)

	smokeTest := NewSmokeTest(goerliClient, zevmClient, cctxClient, fungibleClient, goerliAuth, zevmAuth)
	// The following deployment must happen here and in this order, please do not change
	// ==================== Deploying contracts ====================
	startTime := time.Now()
	smokeTest.TestSetupZetaTokenAndConnectorContracts()
	smokeTest.TestDepositEtherIntoZRC20()
	smokeTest.TestSendZetaIn()
	fmt.Printf("## Essential tests takes %s\n", time.Since(startTime))
	// The following tests are optional tests; comment out the ones you don't want to run
	// temporarily to reduce dev/test cycle turnaround time
	smokeTest.TestERC20Deposit()
	smokeTest.TestERC20Withdraw()
	smokeTest.TestSendZetaOut()
	smokeTest.TestMessagePassing()

	// add your dev test here
	smokeTest.TestMyTest()

	smokeTest.wg.Wait()
}
