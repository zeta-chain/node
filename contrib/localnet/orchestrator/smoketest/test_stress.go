//go:build PRIVNET
// +build PRIVNET

package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	types2 "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	"google.golang.org/grpc"
	"math/big"
	"os"
	"sort"
	"time"
)

const (
	StatInterval      = 5
	WithdrawInterval  = 1000
	StressTestTimeout = 1000 * time.Minute
)

var (
	zevmNonce = big.NewInt(1)
)

var StressCmd = &cobra.Command{
	Use:   "stress",
	Short: "Run Local Stress Test",
	Run:   StressTest,
}

type stressArguments struct {
	ethURL             string
	grpcURL            string
	zevmURL            string
	deployerAddress    string
	deployerPrivateKey string
	local              bool
}

var stressTestArgs = stressArguments{}

func init() {
	RootCmd.AddCommand(StressCmd)
	StressCmd.Flags().StringVar(&stressTestArgs.ethURL, "ethURL", "http://eth:8545", "--ethURL http://eth:8545")
	StressCmd.Flags().StringVar(&stressTestArgs.grpcURL, "grpcURL", "zetacore0:9090", "--grpcURL zetacore0:9090")
	StressCmd.Flags().StringVar(&stressTestArgs.zevmURL, "zevmURL", "http://zetacore0:8545", "--zevmURL http://zetacore0:8545")
	StressCmd.Flags().StringVar(&stressTestArgs.deployerAddress, "addr", "0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC", "--addr <eth address>")
	StressCmd.Flags().StringVar(&stressTestArgs.deployerPrivateKey, "privKey", "d87baf7bf6dc560a252596678c12e41f7d1682837f05b29d411bc3f78ae2c263", "--privKey <eth private key>")
	StressCmd.Flags().BoolVar(&stressTestArgs.local, "local", true, "--local")

	DeployerAddress = ethcommon.HexToAddress(stressTestArgs.deployerAddress)
}

func StressTest(_ *cobra.Command, _ []string) {
	testStartTime := time.Now()
	defer func() {
		fmt.Println("Smoke test took", time.Since(testStartTime))
	}()
	go func() {
		time.Sleep(StressTestTimeout)
		fmt.Println("Smoke test timed out after", StressTestTimeout)
		os.Exit(1)
	}()

	goerliClient, err := ethclient.Dial(stressTestArgs.ethURL)
	if err != nil {
		panic(err)
	}

	bal, err := goerliClient.BalanceAt(context.TODO(), DeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deployer address: %s, balance: %d Ether\n", DeployerAddress.Hex(), bal.Div(bal, big.NewInt(1e18)))

	chainid, err := goerliClient.ChainID(context.Background())
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

	cctxClient := types.NewQueryClient(grpcConn)
	fungibleClient := fungibletypes.NewQueryClient(grpcConn)

	// Wait for Genesis and keygen to be completed. ~ height 30
	time.Sleep(20 * time.Second)
	for {
		time.Sleep(5 * time.Second)
		response, err := cctxClient.LastZetaHeight(context.Background(), &types.QueryLastZetaHeightRequest{})
		if err != nil {
			fmt.Printf("cctxClient.LastZetaHeight error: %s", err)
			continue
		}
		if response.Height >= 30 {
			break
		}
		fmt.Printf("Last ZetaHeight: %d\n", response.Height)
	}

	// get the clients for tests
	var zevmClient *ethclient.Client
	for {
		time.Sleep(5 * time.Second)
		fmt.Printf("dialing zevm client: %s\n", stressTestArgs.zevmURL)
		zevmClient, err = ethclient.Dial(stressTestArgs.zevmURL)
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

	smokeTest := NewSmokeTest(goerliClient, zevmClient, cctxClient, fungibleClient, goerliAuth, zevmAuth, nil)

	// If stress test is running on local docker environment
	if stressTestArgs.local {
		smokeTest.TestSetupZetaTokenAndConnectorAndZEVMContracts()
		smokeTest.TestDepositEtherIntoZRC20()
		smokeTest.TestSendZetaIn()
	}

	//Pre-approve USDT withdraw on ZEVM
	fmt.Printf("approving ETH ZRC20...\n")
	ethZRC20 := smokeTest.ETHZRC20
	tx, err := ethZRC20.Approve(smokeTest.zevmAuth, smokeTest.ETHZRC20Addr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt := MustWaitForTxReceipt(smokeTest.zevmClient, tx)
	fmt.Printf("eth zrc20 approve receipt: status %d\n", receipt.Status)

	// Get current nonce on zevm for DeployerAddress - Need to keep track of nonce at client level
	blockNum, err := smokeTest.zevmClient.BlockNumber(context.Background())
	nonce, err := smokeTest.zevmClient.NonceAt(context.Background(), DeployerAddress, big.NewInt(int64(blockNum)))
	if err != nil {
		panic(err)
	}
	zevmNonce = big.NewInt(int64(nonce))

	// -------------- TEST BEGINS ------------------

	fmt.Println("**** STRESS TEST BEGINS ****")
	fmt.Println("	1. Periodically Withdraw ETH from ZEVM to EVM - goerli")
	fmt.Println("	2. Display Network metrics to monitor performance [Num Pending outbound tx], [Num Trackers]")

	smokeTest.wg.Add(2)
	go smokeTest.WithdrawCCtx()       // Withdraw USDT from ZEVM to EVM - goerli
	go smokeTest.EchoNetworkMetrics() // Display Network metrics periodically to monitor performance

	smokeTest.wg.Wait()
}

// WithdrawCCtx withdraw USDT from ZEVM to EVM
func (sm *SmokeTest) WithdrawCCtx() {
	ticker := time.NewTicker(time.Millisecond * WithdrawInterval)
	for {
		select {
		case <-ticker.C:
			sm.WithdrawETHZRC20()
		}
	}
}

func (sm *SmokeTest) EchoNetworkMetrics() {
	echoStartTime := time.Now()
	ticker := time.NewTicker(time.Second * StatInterval)
	for {
		select {
		case <-ticker.C:
			// Get all pending outbound transactions
			cctxResp, err := sm.cctxClient.CctxAllPending(context.Background(), &types2.QueryAllCctxPendingRequest{
				ChainId: uint64(common.GoerliChain().ChainId),
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
			}
			//
			// Get all trackers
			trackerResp, err := sm.cctxClient.OutTxTrackerAll(context.Background(), &types2.QueryAllOutTxTrackerRequest{})
			if err != nil {
				continue
			}

			elapsedTime := time.Since(echoStartTime).Minutes()
			rate := float64(sends[0].GetCurrentOutTxParam().OutboundTxTssNonce) / elapsedTime

			numPending := len(cctxResp.CrossChainTx)
			numTrackers := len(trackerResp.OutTxTracker)

			fmt.Println("Network Stat => Num of Pending cctx: ", numPending, "Num active trackers: ", numTrackers, "Tx Rate: ", rate, " tx/min")
		}
	}
}

func (sm *SmokeTest) WithdrawETHZRC20() {
	defer func() {
		zevmNonce.Add(zevmNonce, big.NewInt(1))
	}()

	ethZRC20 := sm.ETHZRC20

	sm.zevmAuth.Nonce = zevmNonce
	_, err := ethZRC20.Withdraw(sm.zevmAuth, DeployerAddress.Bytes(), big.NewInt(100))
	if err != nil {
		panic(err)
	}
}
