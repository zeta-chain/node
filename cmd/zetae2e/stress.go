package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"sort"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	"github.com/zeta-chain/zetacore/app"
	zetae2econfig "github.com/zeta-chain/zetacore/cmd/zetae2e/config"
	"github.com/zeta-chain/zetacore/cmd/zetae2e/local"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
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
		Short: "Run Stress Test",
		Run:   StressTest,
	}

	StressCmd.Flags().StringVar(&stressTestArgs.deployerAddress, "addr", "0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC", "--addr <eth address>")
	StressCmd.Flags().StringVar(&stressTestArgs.deployerPrivateKey, "privKey", "d87baf7bf6dc560a252596678c12e41f7d1682837f05b29d411bc3f78ae2c263", "--privKey <eth private key>")
	StressCmd.Flags().StringVar(&stressTestArgs.network, "network", "LOCAL", "--network TESTNET")
	StressCmd.Flags().Int64Var(&stressTestArgs.txnInterval, "tx-interval", 500, "--tx-interval [TIME_INTERVAL_MILLISECONDS]")
	StressCmd.Flags().BoolVar(&stressTestArgs.contractsDeployed, "contracts-deployed", false, "--contracts-deployed=false")
	StressCmd.Flags().StringVar(&stressTestArgs.config, local.FlagConfigFile, "", "config file to use for the E2E test")
	StressCmd.Flags().Bool(flagVerbose, false, "set to true to enable verbose logging")

	local.DeployerAddress = ethcommon.HexToAddress(stressTestArgs.deployerAddress)

	return StressCmd
}

func StressTest(cmd *cobra.Command, _ []string) {
	testStartTime := time.Now()
	defer func() {
		fmt.Println("E2E test took", time.Since(testStartTime))
	}()
	go func() {
		time.Sleep(StressTestTimeout)
		fmt.Println("E2E test timed out after", StressTestTimeout)
		os.Exit(1)
	}()

	// set account prefix to zeta
	cosmosConf := sdk.GetConfig()
	cosmosConf.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	cosmosConf.Seal()

	// initialize E2E tests config
	conf, err := local.GetConfig(cmd)
	if err != nil {
		panic(err)
	}

	// Initialize clients ----------------------------------------------------------------
	evmClient, err := ethclient.Dial(conf.RPCs.EVM)
	if err != nil {
		panic(err)
	}

	bal, err := evmClient.BalanceAt(context.TODO(), local.DeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deployer address: %s, balance: %d Wei\n", local.DeployerAddress.Hex(), bal)

	grpcConn, err := grpc.Dial(conf.RPCs.ZetaCoreGRPC, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	cctxClient := crosschaintypes.NewQueryClient(grpcConn)
	// -----------------------------------------------------------------------------------

	// Wait for Genesis and keygen to be completed if network is local. ~ height 30
	if stressTestArgs.network == "LOCAL" {
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
	}

	// initialize context
	ctx, cancel := context.WithCancel(context.Background())

	verbose, err := cmd.Flags().GetBool(flagVerbose)
	if err != nil {
		panic(err)
	}
	logger := runner.NewLogger(verbose, color.FgWhite, "setup")

	// initialize E2E test runner
	e2eTest, err := zetae2econfig.RunnerFromConfig(
		ctx,
		"deployer",
		cancel,
		conf,
		local.DeployerAddress,
		local.DeployerPrivateKey,
		utils.FungibleAdminName,
		local.FungibleAdminMnemonic,
		logger,
	)
	if err != nil {
		panic(err)
	}

	// setup TSS addresses
	if err := e2eTest.SetTSSAddresses(); err != nil {
		panic(err)
	}

	e2eTest.SetupEVM(stressTestArgs.contractsDeployed, true)

	// If stress test is running on local docker environment
	if stressTestArgs.network == "LOCAL" {
		// deploy and set zevm contract
		e2eTest.SetZEVMContracts()

		// deposit on ZetaChain
		e2eTest.DepositEther(false)
		e2eTest.DepositZeta()
	} else if stressTestArgs.network == "TESTNET" {
		ethZRC20Addr, err := e2eTest.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(5))
		if err != nil {
			panic(err)
		}
		e2eTest.ETHZRC20Addr = ethZRC20Addr
		e2eTest.ETHZRC20, err = zrc20.NewZRC20(e2eTest.ETHZRC20Addr, e2eTest.ZEVMClient)
		if err != nil {
			panic(err)
		}
	} else {
		err := errors.New("invalid network argument: " + stressTestArgs.network)
		panic(err)
	}

	// Check zrc20 balance of Deployer address
	ethZRC20Balance, err := e2eTest.ETHZRC20.BalanceOf(nil, local.DeployerAddress)
	if err != nil {
		panic(err)
	}
	fmt.Printf("eth zrc20 balance: %s Wei \n", ethZRC20Balance.String())

	//Pre-approve ETH withdraw on ZEVM
	fmt.Printf("approving ETH ZRC20...\n")
	ethZRC20 := e2eTest.ETHZRC20
	tx, err := ethZRC20.Approve(e2eTest.ZEVMAuth, e2eTest.ETHZRC20Addr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(ctx, e2eTest.ZEVMClient, tx, logger, e2eTest.ReceiptTimeout)
	fmt.Printf("eth zrc20 approve receipt: status %d\n", receipt.Status)

	// Get current nonce on zevm for DeployerAddress - Need to keep track of nonce at client level
	blockNum, err := e2eTest.ZEVMClient.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}

	// #nosec G701 e2eTest - always in range
	nonce, err := e2eTest.ZEVMClient.NonceAt(context.Background(), local.DeployerAddress, big.NewInt(int64(blockNum)))
	if err != nil {
		panic(err)
	}

	// #nosec G701 e2e - always in range
	zevmNonce = big.NewInt(int64(nonce))

	// -------------- TEST BEGINS ------------------

	fmt.Println("**** STRESS TEST BEGINS ****")
	fmt.Println("	1. Periodically Withdraw ETH from ZEVM to EVM")
	fmt.Println("	2. Display Network metrics to monitor performance [Num Pending outbound tx], [Num Trackers]")

	e2eTest.WG.Add(2)
	go WithdrawCCtx(e2eTest)       // Withdraw from ZEVM to EVM
	go EchoNetworkMetrics(e2eTest) // Display Network metrics periodically to monitor performance

	e2eTest.WG.Wait()
}

// WithdrawCCtx withdraw ETHZRC20 from ZEVM to EVM
func WithdrawCCtx(runner *runner.E2ERunner) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(stressTestArgs.txnInterval))
	for {
		select {
		case <-ticker.C:
			WithdrawETHZRC20(runner)
		}
	}
}

func EchoNetworkMetrics(runner *runner.E2ERunner) {
	ticker := time.NewTicker(time.Second * StatInterval)
	var queue = make([]uint64, 0)
	var numTicks = 0
	var totalMinedTxns = uint64(0)
	var previousMinedTxns = uint64(0)
	chainID, err := getChainID(runner.EVMClient)

	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-ticker.C:
			numTicks++
			// Get all pending outbound transactions
			cctxResp, err := runner.CctxClient.ListPendingCctx(context.Background(), &crosschaintypes.QueryListPendingCctxRequest{
				ChainId: chainID.Int64(),
			})
			if err != nil {
				continue
			}
			sends := cctxResp.CrossChainTx
			sort.Slice(sends, func(i, j int) bool {
				return sends[i].GetCurrentOutboundParam().TssNonce < sends[j].GetCurrentOutboundParam().TssNonce
			})
			if len(sends) > 0 {
				fmt.Printf("pending nonces %d to %d\n", sends[0].GetCurrentOutboundParam().TssNonce, sends[len(sends)-1].GetCurrentOutboundParam().TssNonce)
			} else {
				continue
			}
			//
			// Get all trackers
			trackerResp, err := runner.CctxClient.OutboundTrackerAll(context.Background(), &crosschaintypes.QueryAllOutboundTrackerRequest{})
			if err != nil {
				continue
			}

			currentMinedTxns := sends[0].GetCurrentOutboundParam().TssNonce
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
			numTrackers := len(trackerResp.OutboundTracker)

			fmt.Println("Network Stat => Num of Pending cctx: ", numPending, "Num active trackers: ", numTrackers, "Tx Rate: ", rate, " tx/min")
		}
	}
}

func WithdrawETHZRC20(runner *runner.E2ERunner) {
	defer func() {
		zevmNonce.Add(zevmNonce, big.NewInt(1))
	}()

	ethZRC20 := runner.ETHZRC20

	runner.ZEVMAuth.Nonce = zevmNonce
	_, err := ethZRC20.Withdraw(runner.ZEVMAuth, local.DeployerAddress.Bytes(), big.NewInt(100))
	if err != nil {
		panic(err)
	}
}

// Get ETH based chain ID
func getChainID(client *ethclient.Client) (*big.Int, error) {
	return client.ChainID(context.Background())
}
