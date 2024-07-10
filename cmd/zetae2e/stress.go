package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"sort"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	"google.golang.org/grpc"

	"github.com/zeta-chain/zetacore/app"
	zetae2econfig "github.com/zeta-chain/zetacore/cmd/zetae2e/config"
	"github.com/zeta-chain/zetacore/cmd/zetae2e/local"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/testutil"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

const (
	StatInterval      = 5
	StressTestTimeout = 100 * time.Minute
)

var (
	zevmNonce = big.NewInt(1)
)

type stressArguments struct {
	network           string
	txnInterval       int64
	contractsDeployed bool
	config            string
}

var stressTestArgs = stressArguments{}

var noError = testutil.NoError

func NewStressTestCmd() *cobra.Command {
	var StressCmd = &cobra.Command{
		Use:   "stress",
		Short: "Run Stress Test",
		Run:   StressTest,
	}

	StressCmd.Flags().StringVar(&stressTestArgs.network, "network", "LOCAL", "--network TESTNET")
	StressCmd.Flags().
		Int64Var(&stressTestArgs.txnInterval, "tx-interval", 500, "--tx-interval [TIME_INTERVAL_MILLISECONDS]")
	StressCmd.Flags().
		BoolVar(&stressTestArgs.contractsDeployed, "contracts-deployed", false, "--contracts-deployed=false")
	StressCmd.Flags().StringVar(&stressTestArgs.config, local.FlagConfigFile, "", "config file to use for the E2E test")
	StressCmd.Flags().Bool(flagVerbose, false, "set to true to enable verbose logging")

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
	conf := must(local.GetConfig(cmd))

	deployerAccount := conf.DefaultAccount

	// Initialize clients ----------------------------------------------------------------
	evmClient := must(ethclient.Dial(conf.RPCs.EVM))
	bal := must(evmClient.BalanceAt(context.TODO(), deployerAccount.EVMAddress(), nil))

	fmt.Printf("Deployer address: %s, balance: %d Wei\n", deployerAccount.EVMAddress().Hex(), bal)

	grpcConn := must(grpc.Dial(conf.RPCs.ZetaCoreGRPC, grpc.WithInsecure()))

	cctxClient := crosschaintypes.NewQueryClient(grpcConn)
	// -----------------------------------------------------------------------------------

	// Wait for Genesis and keygen to be completed if network is local. ~ height 30
	if stressTestArgs.network == "LOCAL" {
		time.Sleep(20 * time.Second)
		for {
			time.Sleep(5 * time.Second)
			response, err := cctxClient.LastZetaHeight(
				context.Background(),
				&crosschaintypes.QueryLastZetaHeightRequest{},
			)
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

	verbose := must(cmd.Flags().GetBool(flagVerbose))
	logger := runner.NewLogger(verbose, color.FgWhite, "setup")

	// initialize E2E test runner
	e2eTest := must(zetae2econfig.RunnerFromConfig(
		ctx,
		"deployer",
		cancel,
		conf,
		conf.DefaultAccount,
		logger,
	))

	// setup TSS addresses
	noError(e2eTest.SetTSSAddresses())
	e2eTest.SetupEVM(stressTestArgs.contractsDeployed, true)

	// If stress test is running on local docker environment
	switch stressTestArgs.network {
	case "LOCAL":
		// deploy and set zevm contract
		e2eTest.SetZEVMContracts()

		// deposit on ZetaChain
		e2eTest.DepositEther(false)
		e2eTest.DepositZeta()
	case "TESTNET":
		ethZRC20Addr := must(e2eTest.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(5)))
		e2eTest.ETHZRC20Addr = ethZRC20Addr

		e2eTest.ETHZRC20 = must(zrc20.NewZRC20(e2eTest.ETHZRC20Addr, e2eTest.ZEVMClient))
	default:
		noError(errors.New("invalid network argument: " + stressTestArgs.network))
	}

	// Check zrc20 balance of Deployer address
	ethZRC20Balance := must(e2eTest.ETHZRC20.BalanceOf(nil, deployerAccount.EVMAddress()))
	fmt.Printf("eth zrc20 balance: %s Wei \n", ethZRC20Balance.String())

	//Pre-approve ETH withdraw on ZEVM
	fmt.Println("approving ETH ZRC20...")
	ethZRC20 := e2eTest.ETHZRC20
	tx := must(ethZRC20.Approve(e2eTest.ZEVMAuth, e2eTest.ETHZRC20Addr, big.NewInt(1e18)))

	receipt := utils.MustWaitForTxReceipt(e2eTest.Ctx, e2eTest.ZEVMClient, tx, logger, e2eTest.ReceiptTimeout)
	fmt.Printf("eth zrc20 approve receipt: status %d\n", receipt.Status)

	// Get current nonce on zevm for DeployerAddress - Need to keep track of nonce at client level
	blockNum := must(e2eTest.ZEVMClient.BlockNumber(ctx))

	// #nosec G115 e2eTest - always in range
	nonce := must(e2eTest.ZEVMClient.NonceAt(ctx, deployerAccount.EVMAddress(), big.NewInt(int64(blockNum))))

	// #nosec G115 e2e - always in range
	zevmNonce = big.NewInt(int64(nonce))

	// -------------- TEST BEGINS ------------------

	fmt.Println("**** STRESS TEST BEGINS ****")
	fmt.Println("	1. Periodically Withdraw ETH from ZEVM to EVM")
	fmt.Println("	2. Display Network metrics to monitor performance [Num Pending outbound tx], [Num Trackers]")

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()

		// Withdraw from ZEVM to EVM
		WithdrawCCtx(e2eTest)
	}()

	go func() {
		defer wg.Done()

		// Display Network metrics periodically to monitor performance
		EchoNetworkMetrics(e2eTest)
	}()

	wg.Wait()
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

func EchoNetworkMetrics(r *runner.E2ERunner) {
	var (
		ticker            = time.NewTicker(time.Second * StatInterval)
		queue             = make([]uint64, 0)
		numTicks          int
		totalMinedTxns    uint64
		previousMinedTxns uint64
		chainID           = must(getChainID(r.EVMClient))
	)

	for {
		select {
		case <-ticker.C:
			numTicks++
			// Get all pending outbound transactions
			cctxResp, err := r.CctxClient.ListPendingCctx(
				context.Background(),
				&crosschaintypes.QueryListPendingCctxRequest{
					ChainId: chainID.Int64(),
				},
			)
			if err != nil {
				continue
			}
			sends := cctxResp.CrossChainTx
			sort.Slice(sends, func(i, j int) bool {
				return sends[i].GetCurrentOutboundParam().TssNonce < sends[j].GetCurrentOutboundParam().TssNonce
			})
			if len(sends) > 0 {
				fmt.Printf(
					"pending nonces %d to %d\n",
					sends[0].GetCurrentOutboundParam().TssNonce,
					sends[len(sends)-1].GetCurrentOutboundParam().TssNonce,
				)
			} else {
				continue
			}
			//
			// Get all trackers
			trackerResp, err := r.CctxClient.OutboundTrackerAll(
				context.Background(),
				&crosschaintypes.QueryAllOutboundTrackerRequest{},
			)
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

			fmt.Println(
				"Network Stat => Num of Pending cctx: ",
				numPending,
				"Num active trackers: ",
				numTrackers,
				"Tx Rate: ",
				rate,
				" tx/min",
			)
		}
	}
}

func WithdrawETHZRC20(r *runner.E2ERunner) {
	defer func() {
		zevmNonce.Add(zevmNonce, big.NewInt(1))
	}()

	ethZRC20 := r.ETHZRC20
	r.ZEVMAuth.Nonce = zevmNonce

	must(ethZRC20.Withdraw(r.ZEVMAuth, r.EVMAddress().Bytes(), big.NewInt(100)))
}

// Get ETH based chain ID
func getChainID(client *ethclient.Client) (*big.Int, error) {
	return client.ChainID(context.Background())
}

func must[T any](v T, err error) T {
	return testutil.Must(v, err)
}
