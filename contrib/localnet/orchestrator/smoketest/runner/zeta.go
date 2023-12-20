package runner

import (
	"context"
	"time"

	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
)

// DepositZeta deposits ZETA on ZetaChain from the ZETA smart contract on EVM
func (sm *SmokeTestRunner) DepositZeta() {
	sm.Logger.Print("⏳ depositing ZETA into ZEVM")
	startTime := time.Now()
	defer func() {
		sm.Logger.Print("✅ ZETA deposited in %s", time.Since(startTime))
	}()

	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(100)) // 100 Zeta
	sm.Logger.InfoLoud("Sending ZETA to ZetaChain\n")
	tx, err := sm.ZetaEth.Approve(sm.GoerliAuth, sm.ConnectorEthAddr, amount)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())
	receipt := utils.MustWaitForTxReceipt(sm.GoerliClient, tx, sm.Logger)
	sm.Logger.Info("Approve tx receipt: status %d", receipt.Status)
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

	sm.Logger.Info("Send tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, tx, sm.Logger)
	sm.Logger.Info("Send tx receipt: status %d", receipt.Status)
	sm.Logger.Info("  Logs:")
	for _, log := range receipt.Logs {
		sentLog, err := sm.ConnectorEth.ParseZetaSent(*log)
		if err == nil {
			sm.Logger.Info("    Dest Addr: %s", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
			sm.Logger.Info("    Dest Chain: %d", sentLog.DestinationChainId)
			sm.Logger.Info("    Dest Gas: %d", sentLog.DestinationGasLimit)
			sm.Logger.Info("    Zeta Value: %d", sentLog.ZetaValueAndGas)
			sm.Logger.Info("    Block Num: %d", log.BlockNumber)
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
		sm.Logger.Info("Zeta block %d, Initial Deployer Zeta balance: %d", bn, initialBal)

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
			sm.Logger.Info("Zeta block %d, Deployer %s Zeta balance: %d", bn, sm.DeployerAddress.Hex(), bal)

			diff := big.NewInt(0)
			diff.Sub(bal, initialBal)

			if diff.Cmp(amount) == 0 {
				sm.Logger.Info("Expected zeta balance; success!")
				break
			}
		}
	}()
	sm.WG.Wait()
}
