package smoketests

import (
	"context"
	"math/big"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/testdapp"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	cctxtypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMessagePassing(sm *runner.SmokeTestRunner) {
	sm.Logger.Info("Approving ConnectorEth to spend deployer's ZetaEth")
	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(10)) // 10 Zeta
	auth := sm.GoerliAuth
	tx, err := sm.ZetaEth.Approve(auth, sm.ConnectorEthAddr, amount)
	if err != nil {
		panic(err)
	}

	sm.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())
	receipt := utils.MustWaitForTxReceipt(sm.GoerliClient, tx, sm.Logger)
	sm.Logger.Info("Approve tx receipt: %d", receipt.Status)
	sm.Logger.Info("Calling ConnectorEth.Send")
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

	sm.Logger.Info("ConnectorEth.Send tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, tx, sm.Logger)
	sm.Logger.Info("ConnectorEth.Send tx receipt: status %d", receipt.Status)
	sm.Logger.Info("  Logs:")
	for _, log := range receipt.Logs {
		sentLog, err := sm.ConnectorEth.ParseZetaSent(*log)
		if err == nil {
			sm.Logger.Info("    Dest Addr: %s", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
			sm.Logger.Info("    Dest Chain: %d", sentLog.DestinationChainId)
			sm.Logger.Info("    Dest Gas: %d", sentLog.DestinationGasLimit)
			sm.Logger.Info("    Zeta Value: %d", sentLog.ZetaValueAndGas)
		}
	}

	sm.Logger.Info("Waiting for ConnectorEth.Send CCTX to be mined...")
	sm.Logger.Info("  INTX hash: %s", receipt.TxHash.String())
	cctx := utils.WaitCctxMinedByInTxHash(receipt.TxHash.String(), sm.CctxClient, sm.Logger)
	receipt, err = sm.GoerliClient.TransactionReceipt(context.Background(), ethcommon.HexToHash(cctx.GetCurrentOutTxParam().OutboundTxHash))
	if err != nil {
		panic(err)
	}
	for _, log := range receipt.Logs {
		event, err := sm.ConnectorEth.ParseZetaReceived(*log)
		if err == nil {
			sm.Logger.Info("Received ZetaSent event:")
			sm.Logger.Info("  Dest Addr: %s", event.DestinationAddress)
			sm.Logger.Info("  Zeta Value: %d", event.ZetaValue)
			sm.Logger.Info("  src chainid: %d", event.SourceChainId)
			if event.ZetaValue.Cmp(cctx.GetCurrentOutTxParam().Amount.BigInt()) != 0 {
				panic("Zeta value mismatch")
			}
		}
	}
}

func TestMessagePassingRevertFail(sm *runner.SmokeTestRunner) {
	sm.Logger.Info("Approving ConnectorEth to spend deployer's ZetaEth")

	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(10)) // 10 Zeta
	auth := sm.GoerliAuth
	tx, err := sm.ZetaEth.Approve(auth, sm.ConnectorEthAddr, amount)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())
	receipt := utils.MustWaitForTxReceipt(sm.GoerliClient, tx, sm.Logger)
	sm.Logger.Info("Approve tx receipt: %d", receipt.Status)
	sm.Logger.Info("Calling ConnectorEth.Send")
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
	sm.Logger.Info("ConnectorEth.Send tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, tx, sm.Logger)
	sm.Logger.Info("ConnectorEth.Send tx receipt: status %d", receipt.Status)
	sm.Logger.Info("  Logs:")
	for _, log := range receipt.Logs {
		sentLog, err := sm.ConnectorEth.ParseZetaSent(*log)
		if err == nil {
			sm.Logger.Info("    Dest Addr: %s", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
			sm.Logger.Info("    Dest Chain: %d", sentLog.DestinationChainId)
			sm.Logger.Info("    Dest Gas: %d", sentLog.DestinationGasLimit)
			sm.Logger.Info("    Zeta Value: %d", sentLog.ZetaValueAndGas)
		}
	}

	// expect revert tx to fail
	cctx := utils.WaitCctxMinedByInTxHash(receipt.TxHash.String(), sm.CctxClient, sm.Logger)
	receipt, err = sm.GoerliClient.TransactionReceipt(context.Background(), ethcommon.HexToHash(cctx.GetCurrentOutTxParam().OutboundTxHash))
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
}

func TestMessagePassingRevertSuccess(sm *runner.SmokeTestRunner) {
	sm.Logger.Info("Approving TestDApp to spend deployer's ZetaEth")

	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(10)) // 10 Zeta
	auth := sm.GoerliAuth

	tx, err := sm.ZetaEth.Approve(auth, sm.TestDAppAddr, amount)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(sm.GoerliClient, tx, sm.Logger)
	sm.Logger.Info("Approve tx receipt: %d", receipt.Status)

	sm.Logger.Info("Calling TestDApp.SendHello on contract address %s", sm.TestDAppAddr.Hex())
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
	sm.Logger.Info("$$$ Before: SUPPLY OF AZETA: %d", res2.Amount.Amount)

	tx, err = testDApp.SendHelloWorld(auth, sm.TestDAppAddr, big.NewInt(1337), amount, true)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("TestDApp.SendHello tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(sm.GoerliClient, tx, sm.Logger)
	sm.Logger.Info("TestDApp.SendHello tx receipt: status %d", receipt.Status)

	cctx := utils.WaitCctxMinedByInTxHash(receipt.TxHash.String(), sm.CctxClient, sm.Logger)
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
			sm.Logger.Info("ZetaReverted event: ")
			sm.Logger.Info("  Dest Addr: %s", ethcommon.BytesToAddress(event.DestinationAddress).Hex())
			sm.Logger.Info("  Dest Chain: %d", event.DestinationChainId)
			sm.Logger.Info("  RemainingZetaValue: %d", event.RemainingZetaValue)
			sm.Logger.Info("  Message: %x", event.Message)
		}
	}
	res3, err := sm.BankClient.SupplyOf(context.Background(), &banktypes.QuerySupplyOfRequest{
		Denom: "azeta",
	})
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("$$$ After: SUPPLY OF AZETA: %d", res3.Amount.Amount.BigInt())
	sm.Logger.Info("$$$ Diff: SUPPLY OF AZETA: %d", res3.Amount.Amount.Sub(res2.Amount.Amount).BigInt())
}
