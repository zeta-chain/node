//go:build PRIVNET
// +build PRIVNET

package main

import (
	"context"
	"fmt"
	types2 "github.com/zeta-chain/zetacore/x/crosschain/types"
	"math/big"
	"time"
)

const (
	StatInterval     = 5
	DepositInterval  = 1
	WithdrawInterval = 1
)

var (
	zevmNonce = big.NewInt(1)
)

func (sm *SmokeTest) StressTestCCTXSwap() {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	LoudPrintf("Stress Testing crosschain swap and ZRC20 withdraw...\n")
	// Firstly, deposit 1.15 BTC into Zeta for liquidity
	// sm.DepositBTC()
	// Secondly, deposit 1000.0 USDT into Zeta for liquidity
	//LoudPrintf("Depositing 1000 USDT & 1.15 BTC for liquidity\n")
	//
	//txhash := sm.DepositERC20(big.NewInt(1e9), []byte{})
	//WaitCctxMinedByInTxHash(txhash.Hex(), sm.cctxClient)

	// Create Uni-swap pair for USDT <-> BTC
	//sm.zevmAuth.GasLimit = 10000000
	//tx, err := sm.UniswapV2Factory.CreatePair(sm.zevmAuth, sm.USDTZRC20Addr, sm.BTCZRC20Addr)
	//if err != nil {
	//	panic(err)
	//}
	//receipt := MustWaitForTxReceipt(sm.zevmClient, tx)
	//usdtBtcPair, err := sm.UniswapV2Factory.GetPair(&bind.CallOpts{}, sm.USDTZRC20Addr, sm.BTCZRC20Addr)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("USDT-BTC pair receipt txhash %s status %d pair addr %s\n", receipt.TxHash, receipt.Status, usdtBtcPair.Hex())

	// Pre-approve 1 BTC and 1 USDT
	//tx, err = sm.USDTZRC20.Approve(sm.zevmAuth, sm.UniswapV2RouterAddr, big.NewInt(1e18))
	//if err != nil {
	//	panic(err)
	//}
	//receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	//fmt.Printf("USDT ZRC20 approval receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)
	//
	//tx, err = sm.BTCZRC20.Approve(sm.zevmAuth, sm.UniswapV2RouterAddr, big.NewInt(1e18))
	//if err != nil {
	//	panic(err)
	//}
	//receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	//fmt.Printf("BTC ZRC20 approval receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)

	//Pre-approve USDT withdraw on ZEVM
	zevmClient := sm.zevmClient
	//usdtZRC20, err := zrc20.NewZRC20(ethcommon.HexToAddress(USDTZRC20Addr), zevmClient)
	//if err != nil {
	//	panic(err)
	//}
	//gasZRC20, _, err := usdtZRC20.WithdrawGasFee(&bind.CallOpts{})
	//if err != nil {
	//	panic(err)
	//}
	fmt.Printf("approving ETH ZRC20...\n")
	ethZRC20 := sm.ETHZRC20
	tx, err := ethZRC20.Approve(sm.zevmAuth, sm.ETHZRC20Addr, big.NewInt(1e18))
	receipt := MustWaitForTxReceipt(zevmClient, tx)
	fmt.Printf("eth zrc20 approve receipt: status %d\n", receipt.Status)

	//tx, err = ethZRC20.Approve(sm.zevmAuth, ethcommon.HexToAddress(USDTZRC20Addr), big.NewInt(1e18))
	//if err != nil {
	//	panic(err)
	//}
	//receipt = MustWaitForTxReceipt(zevmClient, tx)
	//fmt.Printf("eth zrc20 approve receipt: status %d\n", receipt.Status)

	// Add 100 USDT liq and 0.001 BTC
	//bal, err := sm.BTCZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("balance of deployer on BTC ZRC20: %d\n", bal)
	//bal, err = sm.USDTZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("balance of deployer on USDT ZRC20: %d\n", bal)
	//tx, err = sm.UniswapV2Router.AddLiquidity(sm.zevmAuth, sm.USDTZRC20Addr, sm.BTCZRC20Addr, big.NewInt(1e8), big.NewInt(1e8), big.NewInt(1e8), big.NewInt(1e5), DeployerAddress, big.NewInt(time.Now().Add(10*time.Minute).Unix()))
	//if err != nil {
	//	fmt.Printf("Error liq %s", err.Error())
	//	panic(err)
	//}
	//receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	//fmt.Printf("Add liquidity receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)

	//fmt.Printf("Funding contracts ZEVMSwapApp with gas ZRC20s; 1e7 ETH, 1e6 BTC\n")
	//// Fund ZEVMSwapApp with gas ZRC20s
	//tx, err = sm.ETHZRC20.Transfer(sm.zevmAuth, sm.ZEVMSwapAppAddr, big.NewInt(1e7))
	//if err != nil {
	//	panic(err)
	//}
	//receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	//fmt.Printf("  USDT ZRC20 transfer receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)
	//bal1, _ := sm.ETHZRC20.BalanceOf(&bind.CallOpts{}, sm.ZEVMSwapAppAddr)
	//fmt.Printf("  ZEVMSwapApp ETHZRC20 balance %d", bal1)
	//tx, err = sm.BTCZRC20.Transfer(sm.zevmAuth, sm.ZEVMSwapAppAddr, big.NewInt(1e6))
	//if err != nil {
	//	panic(err)
	//}
	//receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	//fmt.Printf("  BTC ZRC20 transfer receipt txhash %s status %d\n", receipt.TxHash, receipt.Status)
	//bal2, _ := sm.BTCZRC20.BalanceOf(&bind.CallOpts{}, sm.ZEVMSwapAppAddr)
	//fmt.Printf("  ZEVMSwapApp BTCZRC20 balance %d", bal2)

	// msg would be [ZEVMSwapAppAddr, memobytes]
	// memobytes is dApp specific; see the contracts/ZEVMSwapApp.sol for details
	//msg := []byte{}
	//msg = append(msg, sm.ZEVMSwapAppAddr.Bytes()...)
	//memobytes, err := sm.ZEVMSwapApp.EncodeMemo(&bind.CallOpts{}, sm.BTCZRC20Addr, []byte(BTCDeployerAddress.EncodeAddress()))
	//
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("memobytes(%d) %x\n", len(memobytes), memobytes)
	//msg = append(msg, memobytes...)

	// Get current nonce on zevm for DeployerAddress - Need to keep track of nonce at client level
	blockNum, err := sm.zevmClient.BlockNumber(context.Background())
	nonce, err := sm.zevmClient.NonceAt(context.Background(), DeployerAddress, big.NewInt(int64(blockNum)))
	if err != nil {
		panic(err)
	}
	zevmNonce = big.NewInt(int64(nonce))

	// -------------- TEST BEGINS ------------------

	fmt.Println("**** STRESS TEST BEGINS ****")
	fmt.Println("	1. Periodically deposit USDT with a memo to swap for BTC")
	fmt.Println("	2. Periodically Withdraw USDT from ZEVM to EVM - goerli")
	fmt.Println("	3. Display Network metrics to monitor performance [Num Pending outbound tx], [Num Trackers]")

	sm.wg.Add(2)
	//go sm.SendCCtx(msg)        // Add goroutine to periodically deposit USDT with a memo to swap for BTC
	go sm.WithdrawCCtx()       // Withdraw USDT from ZEVM to EVM - goerli
	go sm.EchoNetworkMetrics() // Display Network metrics periodically to monitor performance

	sm.wg.Wait()
}

// SendCCtx Send USDT deposit for a BTC swap once every block
func (sm *SmokeTest) SendCCtx(msg []byte) {
	// timeout_commit=2s - Wait 2 seconds before sending next deposit
	ticker := time.NewTicker(time.Second * DepositInterval)
	for {
		select {
		case <-ticker.C:
			sm.DepositUSDTERC20(big.NewInt(8e5), msg)
		}
	}
}

// WithdrawCCtx withdraw USDT from ZEVM to EVM
func (sm *SmokeTest) WithdrawCCtx() {
	ticker := time.NewTicker(time.Second * WithdrawInterval)
	for {
		select {
		case <-ticker.C:
			sm.WithdrawETHZRC20()
		}
	}
}

func (sm *SmokeTest) EchoNetworkMetrics() {
	var minedTxs = map[string]bool{}
	var queue = make([]int, 0)
	var numTicks = 0
	var totalMinedTxns = 0
	ticker := time.NewTicker(time.Second * StatInterval)

	for {
		select {
		case <-ticker.C:
			numTicks++
			// Get all pending outbound transactions
			//==> Check for status == CctxStatus_OutboundMined using sm.cctxClient.CctxAll()
			cctxResp, err := sm.cctxClient.CctxAllPending(context.Background(), &types2.QueryAllCctxPendingRequest{})
			if err != nil {
				continue
			}
			// Get all cross chain transactions
			cctxAll, err := sm.cctxClient.CctxAll(context.Background(), &types2.QueryAllCctxRequest{})
			if err != nil {
				continue
			}
			//
			// Get all trackers
			trackerResp, err := sm.cctxClient.OutTxTrackerAll(context.Background(), &types2.QueryAllOutTxTrackerRequest{})
			if err != nil {
				continue
			}

			newMinedTxCnt := 0
			for _, tx := range cctxAll.CrossChainTx {
				if tx.CctxStatus.Status == types2.CctxStatus_OutboundMined {
					if _, found := minedTxs[tx.Index]; !found {
						minedTxs[tx.Index] = true
						newMinedTxCnt++
					}
				}
			}
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
			fmt.Println("new mined txn count: ", newMinedTxCnt)
		}
	}
}

func (sm *SmokeTest) DepositUSDTERC20(amount *big.Int, msg []byte) {
	_, err := sm.ERC20Custody.Deposit(sm.goerliAuth, DeployerAddress.Bytes(), sm.USDTERC20Addr, amount, msg)
	if err != nil {
		panic(err)
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
