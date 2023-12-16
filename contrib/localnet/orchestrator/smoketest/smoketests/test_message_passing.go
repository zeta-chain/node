package smoketests

import (
	"context"
	"fmt"
	"math/big"
	"time"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/testdapp"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	cctxtypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMessagePassing(sm *runner.SmokeTestRunner) {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	// ==================== Interacting with contracts ====================
	time.Sleep(10 * time.Second)
	utils.LoudPrintf("Goerli->Goerli Message Passing (Sending ZETA only)\n")
	fmt.Printf("Approving ConnectorEth to spend deployer's ZetaEth\n")
	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(10)) // 10 Zeta
	auth := sm.GoerliAuth
	tx, err := sm.ZetaEth.Approve(auth, sm.ConnectorEthAddr, amount)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Approve tx hash: %s\n", tx.Hash().Hex())
	receipt := utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
	fmt.Printf("Approve tx receipt: %d\n", receipt.Status)
	fmt.Printf("Calling ConnectorEth.Send\n")
	tx, err = sm.ConnectorEth.Send(auth, zetaconnectoreth.ZetaInterfacesSendInput{
		DestinationChainId:  big.NewInt(1337), // in dev mode, GOERLI has chainid 1337
		DestinationAddress:  sm.DeployerAddress.Bytes(),
		DestinationGasLimit: big.NewInt(250_000),
		Message:             nil,
		ZetaValueAndGas:     amount,
		ZetaParams:          nil,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("ConnectorEth.Send tx hash: %s\n", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
	fmt.Printf("ConnectorEth.Send tx receipt: status %d\n", receipt.Status)
	fmt.Printf("  Logs:\n")
	for _, log := range receipt.Logs {
		sentLog, err := sm.ConnectorEth.ParseZetaSent(*log)
		if err == nil {
			fmt.Printf("    Dest Addr: %s\n", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
			fmt.Printf("    Dest Chain: %d\n", sentLog.DestinationChainId)
			fmt.Printf("    Dest Gas: %d\n", sentLog.DestinationGasLimit)
			fmt.Printf("    Zeta Value: %d\n", sentLog.ZetaValueAndGas)
		}
	}
	sm.WG.Add(1)
	go func() {
		defer sm.WG.Done()
		fmt.Printf("Waiting for ConnectorEth.Send CCTX to be mined...\n")
		fmt.Printf("  INTX hash: %s\n", receipt.TxHash.String())
		cctx := utils.WaitCctxMinedByInTxHash(receipt.TxHash.String(), sm.CctxClient)
		receipt, err := sm.GoerliClient.TransactionReceipt(context.Background(), ethcommon.HexToHash(cctx.GetCurrentOutTxParam().OutboundTxHash))
		if err != nil {
			panic(err)
		}
		for _, log := range receipt.Logs {
			event, err := sm.ConnectorEth.ParseZetaReceived(*log)
			if err == nil {
				fmt.Printf("Received ZetaSent event:\n")
				fmt.Printf("  Dest Addr: %s\n", event.DestinationAddress)
				fmt.Printf("  Zeta Value: %d\n", event.ZetaValue)
				fmt.Printf("  src chainid: %d\n", event.SourceChainId)
				if event.ZetaValue.Cmp(cctx.GetCurrentOutTxParam().Amount.BigInt()) != 0 {
					panic("Zeta value mismatch")
				}
			}
		}
	}()
	sm.WG.Wait()
}

func TestMessagePassingRevertFail(sm *runner.SmokeTestRunner) {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	// ==================== Interacting with contracts ====================
	utils.LoudPrintf("Goerli->Goerli Message Passing (revert fail)\n")
	fmt.Printf("Approving ConnectorEth to spend deployer's ZetaEth\n")
	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(10)) // 10 Zeta
	auth := sm.GoerliAuth
	tx, err := sm.ZetaEth.Approve(auth, sm.ConnectorEthAddr, amount)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Approve tx hash: %s\n", tx.Hash().Hex())
	receipt := utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
	fmt.Printf("Approve tx receipt: %d\n", receipt.Status)
	fmt.Printf("Calling ConnectorEth.Send\n")
	tx, err = sm.ConnectorEth.Send(auth, zetaconnectoreth.ZetaInterfacesSendInput{
		DestinationChainId:  big.NewInt(1337), // in dev mode, GOERLI has chainid 1337
		DestinationAddress:  sm.DeployerAddress.Bytes(),
		DestinationGasLimit: big.NewInt(250_000),
		Message:             []byte("revert"), // non-empty message will cause revert, because the dest address is not a contract
		ZetaValueAndGas:     amount,
		ZetaParams:          nil,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("ConnectorEth.Send tx hash: %s\n", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
	fmt.Printf("ConnectorEth.Send tx receipt: status %d\n", receipt.Status)
	fmt.Printf("  Logs:\n")
	for _, log := range receipt.Logs {
		sentLog, err := sm.ConnectorEth.ParseZetaSent(*log)
		if err == nil {
			fmt.Printf("    Dest Addr: %s\n", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
			fmt.Printf("    Dest Chain: %d\n", sentLog.DestinationChainId)
			fmt.Printf("    Dest Gas: %d\n", sentLog.DestinationGasLimit)
			fmt.Printf("    Zeta Value: %d\n", sentLog.ZetaValueAndGas)
		}
	}
	sm.WG.Add(1)
	go func() {
		defer sm.WG.Done()
		fmt.Printf("Waiting for ConnectorEth.Send CCTX to be mined...\n")
		cctx := utils.WaitCctxMinedByInTxHash(receipt.TxHash.String(), sm.CctxClient)
		receipt, err := sm.GoerliClient.TransactionReceipt(context.Background(), ethcommon.HexToHash(cctx.GetCurrentOutTxParam().OutboundTxHash))
		if err != nil {
			panic(err)
		}
		// expect revert tx to fail as well
		if receipt.Status != 0 {
			panic("expected revert tx to fail")
		}
		if cctx.CctxStatus.Status != cctxtypes.CctxStatus_Aborted {
			panic("expected cctx to be aborted")
		}
	}()
	sm.WG.Wait()
}

func TestMessagePassingRevertSuccess(sm *runner.SmokeTestRunner) {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	// ==================== Interacting with contracts ====================
	utils.LoudPrintf("Goerli->Goerli Message Passing (revert success)\n")
	fmt.Printf("Approving TestDApp to spend deployer's ZetaEth\n")
	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(10)) // 10 Zeta
	auth := sm.GoerliAuth
	tx, err := sm.ZetaEth.Approve(auth, sm.TestDAppAddr, amount)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Approve tx hash: %s\n", tx.Hash().Hex())
	receipt := utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
	fmt.Printf("Approve tx receipt: %d\n", receipt.Status)
	fmt.Printf("Calling TestDApp.SendHello on contract address %s\n", sm.TestDAppAddr.Hex())
	testDApp, err := testdapp.NewTestDApp(sm.TestDAppAddr, sm.GoerliClient)
	if err != nil {
		panic(err)
	}

	res2, err := sm.BankClient.SupplyOf(context.Background(), &banktypes.QuerySupplyOfRequest{
		Denom: "azeta",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("$$$ Before: SUPPLY OF AZETA: %d\n", res2.Amount.Amount)

	tx, err = testDApp.SendHelloWorld(auth, sm.TestDAppAddr, big.NewInt(1337), amount, true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("TestDApp.SendHello tx hash: %s\n", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, tx)
	fmt.Printf("TestDApp.SendHello tx receipt: status %d\n", receipt.Status)

	cctx := utils.WaitCctxMinedByInTxHash(receipt.TxHash.String(), sm.CctxClient)
	if cctx.CctxStatus.Status != cctxtypes.CctxStatus_Reverted {
		panic("expected cctx to be reverted")
	}
	outTxHash := cctx.GetCurrentOutTxParam().OutboundTxHash
	receipt, err = sm.GoerliClient.TransactionReceipt(context.Background(), ethcommon.HexToHash(outTxHash))
	if err != nil {
		panic(err)
	}
	for _, log := range receipt.Logs {
		event, err := sm.ConnectorEth.ParseZetaReverted(*log)
		if err == nil {
			fmt.Printf("ZetaReverted event: \n")
			fmt.Printf("  Dest Addr: %s\n", ethcommon.BytesToAddress(event.DestinationAddress).Hex())
			fmt.Printf("  Dest Chain: %d\n", event.DestinationChainId)
			fmt.Printf("  RemainingZetaValue: %d\n", event.RemainingZetaValue)
			fmt.Printf("  Message: %x\n", event.Message)
		}
	}
	res3, err := sm.BankClient.SupplyOf(context.Background(), &banktypes.QuerySupplyOfRequest{
		Denom: "azeta",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("$$$ After: SUPPLY OF AZETA: %d\n", res3.Amount.Amount.BigInt())
	fmt.Printf("$$$ Diff: SUPPLY OF AZETA: %d\n", res3.Amount.Amount.Sub(res2.Amount.Amount).BigInt())
}
