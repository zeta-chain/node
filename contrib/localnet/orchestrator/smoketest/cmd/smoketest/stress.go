package main

import (
	"context"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/app"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/txserver"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	"math/big"
	"os"
	"sort"
	"time"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc"
)

const (
	StatInterval      = 5
	StressTestTimeout = 100 * time.Minute
)

var (
	zevmNonce = big.NewInt(1)
)

type stressArguments struct {
	ethURL             string
	grpcURL            string
	zevmURL            string
	deployerAddress    string
	deployerPrivateKey string
	network            string
	txnInterval        int64
	contractsDeployed  bool
	config             string
}

var stressTestArgs = stressArguments{}

func NewStressTestCmd() *cobra.Command {
	var StressCmd = &cobra.Command{
		Use:   "stress",
		Short: "Run Local Stress Test",
		Run:   StressTest,
	}

	StressCmd.Flags().StringVar(&stressTestArgs.ethURL, "ethURL", "http://eth:8545", "--ethURL http://eth:8545")
	StressCmd.Flags().StringVar(&stressTestArgs.grpcURL, "grpcURL", "zetacore0:9090", "--grpcURL zetacore0:9090")
	StressCmd.Flags().StringVar(&stressTestArgs.deployerAddress, "addr", "0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC", "--addr <eth address>")
	StressCmd.Flags().StringVar(&stressTestArgs.deployerPrivateKey, "privKey", "d87baf7bf6dc560a252596678c12e41f7d1682837f05b29d411bc3f78ae2c263", "--privKey <eth private key>")
	StressCmd.Flags().StringVar(&stressTestArgs.network, "network", "", "--network TESTNET")
	StressCmd.Flags().Int64Var(&stressTestArgs.txnInterval, "tx-interval", 500, "--tx-interval [TIME_INTERVAL_MILLISECONDS]")
	StressCmd.Flags().BoolVar(&stressTestArgs.contractsDeployed, "contracts-deployed", false, "--contracts-deployed=false")
	StressCmd.Flags().StringVar(&stressTestArgs.config, flagConfigFile, "", "config file to use for the smoketest")

	DeployerAddress = ethcommon.HexToAddress(stressTestArgs.deployerAddress)

	return StressCmd
}

func StressTest(cmd *cobra.Command, _ []string) {
	testStartTime := time.Now()
	defer func() {
		fmt.Println("Smoke test took", time.Since(testStartTime))
	}()
	go func() {
		time.Sleep(StressTestTimeout)
		fmt.Println("Smoke test timed out after", StressTestTimeout)
		os.Exit(1)
	}()

	// set account prefix to zeta
	cosmosConf := sdk.GetConfig()
	cosmosConf.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	cosmosConf.Seal()

	// initialize smoke tests config
	conf, err := getConfig(cmd)
	if err != nil {
		panic(err)
	}

	// Initialize clients ----------------------------------------------------------------
	goerliClient, err := ethclient.Dial(stressTestArgs.ethURL)
	if err != nil {
		panic(err)
	}

	bal, err := goerliClient.BalanceAt(context.TODO(), DeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deployer address: %s, balance: %d Wei\n", DeployerAddress.Hex(), bal)

	chainid, err := goerliClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	deployerPrivkey, err := crypto.HexToECDSA(stressTestArgs.deployerPrivateKey)
	if err != nil {
		panic(err)
	}
	goerliAuth, err := bind.NewKeyedTransactorWithChainID(deployerPrivkey, chainid)
	if err != nil {
		panic(err)
	}

	grpcConn, err := grpc.Dial(stressTestArgs.grpcURL, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	cctxClient := crosschaintypes.NewQueryClient(grpcConn)
	fungibleClient := fungibletypes.NewQueryClient(grpcConn)
	bankClient := banktypes.NewQueryClient(grpcConn)
	authClient := authtypes.NewQueryClient(grpcConn)
	observerClient := observertypes.NewQueryClient(grpcConn)
	// -----------------------------------------------------------------------------------

	// Wait for Genesis and keygen to be completed. ~ height 30
	time.Sleep(20 * time.Second)
	for {
		time.Sleep(5 * time.Second)
		response, err := cctxClient.LastZetaHeight(context.Background(), &crosschaintypes.QueryLastZetaHeightRequest{})
		if err != nil {
			fmt.Printf("cctxClient.LastZetaHeight error: %s", err)
			continue
		}
		if response.Height >= 30 {
			break
		}
		fmt.Printf("Last ZetaHeight: %d\n", response.Height)
	}

	// initialize client to send messages to ZetaChain
	zetaTxServer, err := txserver.NewZetaTxServer(
		conf.RPCs.ZetaCoreRPC,
		[]string{utils.FungibleAdminName},
		[]string{FungibleAdminMnemonic},
		conf.ZetaChainID,
	)
	if err != nil {
		panic(err)
	}

	// get the clients for tests
	var zevmClient *ethclient.Client
	for {
		time.Sleep(5 * time.Second)
		fmt.Printf("dialing zevm client: %s\n", conf.RPCs.Zevm)
		zevmClient, err = ethclient.Dial(conf.RPCs.Zevm)
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

	// initialize smoke test runner
	smokeTest := runner.NewSmokeTestRunner(
		DeployerAddress,
		DeployerPrivateKey,
		FungibleAdminMnemonic,
		goerliClient,
		zevmClient,
		cctxClient,
		zetaTxServer,
		fungibleClient,
		authClient,
		bankClient,
		observerClient,
		goerliAuth,
		zevmAuth,
		nil,
	)

	// setup TSS addresses
	smokeTest.SetTSSAddresses()

	smokeTest.SetupEVM(stressTestArgs.contractsDeployed)

	fmt.Printf("Stress test arguments: %+v", stressTestArgs)

	// If stress test is running on local docker environment
	if stressTestArgs.network == "" {
		// deploy and set zevm contract
		smokeTest.SetZEVMContracts()

		// deposit on ZetaChain
		smokeTest.DepositEtherIntoZRC20()
		smokeTest.SendZetaIn()
	} else if stressTestArgs.network == "TESTNET" {
		ethZRC20Addr, err := smokeTest.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(5))
		if err != nil {
			panic(err)
		}
		smokeTest.ETHZRC20Addr = ethZRC20Addr
		smokeTest.ETHZRC20, err = zrc20.NewZRC20(smokeTest.ETHZRC20Addr, smokeTest.ZevmClient)
		if err != nil {
			panic(err)
		}
	} else {
		err := errors.New("invalid network argument: " + stressTestArgs.network)
		panic(err)
	}

	// Check zrc20 balance of Deployer address
	ethZRC20Balance, err := smokeTest.ETHZRC20.BalanceOf(nil, DeployerAddress)
	if err != nil {
		panic(err)
	}
	fmt.Printf("eth zrc20 balance: %s Wei \n", ethZRC20Balance.String())

	//Pre-approve ETH withdraw on ZEVM
	fmt.Printf("approving ETH ZRC20...\n")
	ethZRC20 := smokeTest.ETHZRC20
	tx, err := ethZRC20.Approve(smokeTest.ZevmAuth, smokeTest.ETHZRC20Addr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(smokeTest.ZevmClient, tx)
	fmt.Printf("eth zrc20 approve receipt: status %d\n", receipt.Status)

	// Get current nonce on zevm for DeployerAddress - Need to keep track of nonce at client level
	blockNum, err := smokeTest.ZevmClient.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}

	// #nosec G701 smoketest - always in range
	nonce, err := smokeTest.ZevmClient.NonceAt(context.Background(), DeployerAddress, big.NewInt(int64(blockNum)))
	if err != nil {
		panic(err)
	}

	// #nosec G701 smoketest - always in range
	zevmNonce = big.NewInt(int64(nonce))

	// -------------- TEST BEGINS ------------------

	fmt.Println("**** STRESS TEST BEGINS ****")
	fmt.Println("	1. Periodically Withdraw ETH from ZEVM to EVM - goerli")
	fmt.Println("	2. Display Network metrics to monitor performance [Num Pending outbound tx], [Num Trackers]")

	smokeTest.WG.Add(2)
	go WithdrawCCtx(smokeTest)       // Withdraw USDT from ZEVM to EVM - goerli
	go EchoNetworkMetrics(smokeTest) // Display Network metrics periodically to monitor performance

	smokeTest.WG.Wait()
}

// WithdrawCCtx withdraw USDT from ZEVM to EVM
func WithdrawCCtx(sm *runner.SmokeTestRunner) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(stressTestArgs.txnInterval))
	for {
		select {
		case <-ticker.C:
			WithdrawETHZRC20(sm)
		}
	}
}

func EchoNetworkMetrics(sm *runner.SmokeTestRunner) {
	ticker := time.NewTicker(time.Second * StatInterval)
	var queue = make([]uint64, 0)
	var numTicks = 0
	var totalMinedTxns = uint64(0)
	var previousMinedTxns = uint64(0)

	for {
		select {
		case <-ticker.C:
			numTicks++
			// Get all pending outbound transactions
			cctxResp, err := sm.CctxClient.CctxListPending(context.Background(), &crosschaintypes.QueryListCctxPendingRequest{
				ChainId: getChainID(),
			})
			if err != nil {
				continue
			}
			sends := cctxResp.CrossChainTx
			sort.Slice(sends, func(i, j int) bool {
				return sends[i].GetCurrentOutTxParam().OutboundTxTssNonce < sends[j].GetCurrentOutTxParam().OutboundTxTssNonce
			})
			if len(sends) > 0 {
				fmt.Printf("pending nonces %d to %d\n", sends[0].GetCurrentOutTxParam().OutboundTxTssNonce, sends[len(sends)-1].GetCurrentOutTxParam().OutboundTxTssNonce)
			} else {
				continue
			}
			//
			// Get all trackers
			trackerResp, err := sm.CctxClient.OutTxTrackerAll(context.Background(), &crosschaintypes.QueryAllOutTxTrackerRequest{})
			if err != nil {
				continue
			}

			currentMinedTxns := sends[0].GetCurrentOutTxParam().OutboundTxTssNonce
			newMinedTxCnt := currentMinedTxns - previousMinedTxns
			previousMinedTxns = currentMinedTxns

			// Add new mined txn count to queue and remove the oldest entry
			queue = append(queue, newMinedTxCnt)
			if numTicks > 60/StatInterval {
				totalMinedTxns -= queue[0]
				queue = queue[1:]
				numTicks = 60/StatInterval + 1 //prevent overflow
			}

			//Calculate rate -> tx/min
			totalMinedTxns += queue[len(queue)-1]
			rate := totalMinedTxns

			numPending := len(cctxResp.CrossChainTx)
			numTrackers := len(trackerResp.OutTxTracker)

			fmt.Println("Network Stat => Num of Pending cctx: ", numPending, "Num active trackers: ", numTrackers, "Tx Rate: ", rate, " tx/min")
		}
	}
}

func WithdrawETHZRC20(sm *runner.SmokeTestRunner) {
	defer func() {
		zevmNonce.Add(zevmNonce, big.NewInt(1))
	}()

	ethZRC20 := sm.ETHZRC20

	sm.ZevmAuth.Nonce = zevmNonce
	_, err := ethZRC20.Withdraw(sm.ZevmAuth, DeployerAddress.Bytes(), big.NewInt(100))
	if err != nil {
		panic(err)
	}
}

// Get ETH based chain ID - Build flags are conflicting with current lib, need to do this manually
func getChainID() int64 {
	switch stressTestArgs.network {
	case "PRIVNET":
		return 1337
	case "TESTNET":
		return 5
	case "MAINNET":
		return 1
	default:
		return 1337
	}
}
