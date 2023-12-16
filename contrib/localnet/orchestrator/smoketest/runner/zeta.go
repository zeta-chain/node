package runner

import (
	"context"
	"fmt"
	"time"

	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
)

func (sm *SmokeTestRunner) SendZetaIn() {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	// ==================== Sending ZETA to ZetaChain ===================
	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(100)) // 100 Zeta
	utils.LoudPrintf("Step 3: Sending ZETA to ZetaChain\n")
	tx, err := sm.ZetaEth.Approve(sm.GoerliAuth, sm.ConnectorEthAddr, amount)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Approve tx hash: %s\n", tx.Hash().Hex())
	receipt := utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
	fmt.Printf("Approve tx receipt: status %d\n", receipt.Status)
	tx, err = sm.ConnectorEth.Send(sm.GoerliAuth, zetaconnectoreth.ZetaInterfacesSendInput{
		DestinationChainId:  big.NewInt(101), // in dev mode, 101 is the  zEVM ChainID
		DestinationAddress:  sm.DeployerAddress.Bytes(),
		DestinationGasLimit: big.NewInt(250_000),
		Message:             nil,
		ZetaValueAndGas:     amount,
		ZetaParams:          nil,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Send tx hash: %s\n", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
	fmt.Printf("Send tx receipt: status %d\n", receipt.Status)
	fmt.Printf("  Logs:\n")
	for _, log := range receipt.Logs {
		sentLog, err := sm.ConnectorEth.ParseZetaSent(*log)
		if err == nil {
			fmt.Printf("    Dest Addr: %s\n", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
			fmt.Printf("    Dest Chain: %d\n", sentLog.DestinationChainId)
			fmt.Printf("    Dest Gas: %d\n", sentLog.DestinationGasLimit)
			fmt.Printf("    Zeta Value: %d\n", sentLog.ZetaValueAndGas)
			fmt.Printf("    Block Num: %d\n", log.BlockNumber)
		}
	}

	sm.WG.Add(1)
	go func() {
		bn, err := sm.ZevmClient.BlockNumber(context.Background())
		if err != nil {
			panic(err)
		}
		// #nosec G701 smoketest - always in range
		initialBal, err := sm.ZevmClient.BalanceAt(context.Background(), sm.DeployerAddress, big.NewInt(int64(bn)))
		if err != nil {
			panic(err)
		}
		fmt.Printf("Zeta block %d, Initial Deployer Zeta balance: %d\n", bn, initialBal)

		defer sm.WG.Done()
		for {
			time.Sleep(5 * time.Second)
			bn, err = sm.ZevmClient.BlockNumber(context.Background())
			if err != nil {
				panic(err)
			}
			// #nosec G701 smoketest - always in range
			bal, err := sm.ZevmClient.BalanceAt(context.Background(), sm.DeployerAddress, big.NewInt(int64(bn)))
			if err != nil {
				panic(err)
			}
			fmt.Printf("Zeta block %d, Deployer Zeta balance: %d\n", bn, bal)

			diff := big.NewInt(0)
			diff.Sub(bal, initialBal)

			if diff.Cmp(amount) == 0 {
				fmt.Printf("Expected zeta balance; success!\n")
				break
			}
		}
	}()
	sm.WG.Wait()
}
