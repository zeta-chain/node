package e2etests

import (
	"fmt"
	"github.com/zeta-chain/zetacore/e2e/contracts/testdapp"
	"github.com/zeta-chain/zetacore/e2e/runner"
	utils2 "github.com/zeta-chain/zetacore/e2e/utils"
	"math/big"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"
	cctxtypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMessagePassing(sm *runner.E2ERunner) {
	chainID, err := sm.GoerliClient.ChainID(sm.Ctx)
	if err != nil {
		panic(err)
	}

	sm.Logger.Info("Approving ConnectorEth to spend deployer's ZetaEth")
	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(10)) // 10 Zeta
	auth := sm.GoerliAuth
	tx, err := sm.ZetaEth.Approve(auth, sm.ConnectorEthAddr, amount)
	if err != nil {
		panic(err)
	}

	sm.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())
	receipt := utils2.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, tx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("tx failed")
	}
	sm.Logger.Info("Approve tx receipt: %d", receipt.Status)
	sm.Logger.Info("Calling ConnectorEth.Send")
	tx, err = sm.ConnectorEth.Send(auth, zetaconnectoreth.ZetaInterfacesSendInput{
		DestinationChainId:  chainID,
		DestinationAddress:  sm.DeployerAddress.Bytes(),
		DestinationGasLimit: big.NewInt(400_000),
		Message:             nil,
		ZetaValueAndGas:     amount,
		ZetaParams:          nil,
	})
	if err != nil {
		panic(err)
	}

	sm.Logger.Info("ConnectorEth.Send tx hash: %s", tx.Hash().Hex())
	receipt = utils2.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, tx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("tx failed")
	}
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
	cctx := utils2.WaitCctxMinedByInTxHash(sm.Ctx, receipt.TxHash.String(), sm.CctxClient, sm.Logger, sm.CctxTimeout)
	if cctx.CctxStatus.Status != cctxtypes.CctxStatus_OutboundMined {
		panic(fmt.Sprintf(
			"expected cctx status to be %s; got %s, message %s",
			cctxtypes.CctxStatus_OutboundMined,
			cctx.CctxStatus.Status.String(),
			cctx.CctxStatus.StatusMessage,
		))
	}
	receipt, err = sm.GoerliClient.TransactionReceipt(sm.Ctx, ethcommon.HexToHash(cctx.GetCurrentOutTxParam().OutboundTxHash))
	if err != nil {
		panic(err)
	}
	if receipt.Status != 1 {
		panic("tx failed")
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

func TestMessagePassingRevertFail(sm *runner.E2ERunner) {
	chainID, err := sm.GoerliClient.ChainID(sm.Ctx)
	if err != nil {
		panic(err)
	}

	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(10)) // 10 Zeta
	auth := sm.GoerliAuth
	tx, err := sm.ZetaEth.Approve(auth, sm.ConnectorEthAddr, amount)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())
	receipt := utils2.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, tx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("tx failed")
	}
	sm.Logger.Info("Approve tx receipt: %d", receipt.Status)
	sm.Logger.Info("Calling ConnectorEth.Send")
	tx, err = sm.ConnectorEth.Send(auth, zetaconnectoreth.ZetaInterfacesSendInput{
		DestinationChainId:  chainID,
		DestinationAddress:  sm.DeployerAddress.Bytes(),
		DestinationGasLimit: big.NewInt(400_000),
		Message:             []byte("revert"), // non-empty message will cause revert, because the dest address is not a contract
		ZetaValueAndGas:     amount,
		ZetaParams:          nil,
	})
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("ConnectorEth.Send tx hash: %s", tx.Hash().Hex())
	receipt = utils2.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, tx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("tx failed")
	}
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
	cctx := utils2.WaitCctxMinedByInTxHash(sm.Ctx, receipt.TxHash.String(), sm.CctxClient, sm.Logger, sm.CctxTimeout)
	receipt, err = sm.GoerliClient.TransactionReceipt(sm.Ctx, ethcommon.HexToHash(cctx.GetCurrentOutTxParam().OutboundTxHash))
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

func TestMessagePassingRevertSuccess(sm *runner.E2ERunner) {
	chainID, err := sm.GoerliClient.ChainID(sm.Ctx)
	if err != nil {
		panic(err)
	}

	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(10)) // 10 Zeta
	auth := sm.GoerliAuth

	tx, err := sm.ZetaEth.Approve(auth, sm.TestDAppAddr, amount)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())

	receipt := utils2.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, tx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("tx failed")
	}
	sm.Logger.Info("Approve tx receipt: %d", receipt.Status)

	sm.Logger.Info("Calling TestDApp.SendHello on contract address %s", sm.TestDAppAddr.Hex())
	testDApp, err := testdapp.NewTestDApp(sm.TestDAppAddr, sm.GoerliClient)
	if err != nil {
		panic(err)
	}

	res2, err := sm.BankClient.SupplyOf(sm.Ctx, &banktypes.QuerySupplyOfRequest{
		Denom: "azeta",
	})
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("$$$ Before: SUPPLY OF AZETA: %d", res2.Amount.Amount)

	tx, err = testDApp.SendHelloWorld(auth, sm.TestDAppAddr, chainID, amount, true)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("TestDApp.SendHello tx hash: %s", tx.Hash().Hex())
	receipt = utils2.MustWaitForTxReceipt(sm.Ctx, sm.GoerliClient, tx, sm.Logger, sm.ReceiptTimeout)
	sm.Logger.Info("TestDApp.SendHello tx receipt: status %d", receipt.Status)

	cctx := utils2.WaitCctxMinedByInTxHash(sm.Ctx, receipt.TxHash.String(), sm.CctxClient, sm.Logger, sm.CctxTimeout)
	if cctx.CctxStatus.Status != cctxtypes.CctxStatus_Reverted {
		panic("expected cctx to be reverted")
	}
	outTxHash := cctx.GetCurrentOutTxParam().OutboundTxHash
	receipt, err = sm.GoerliClient.TransactionReceipt(sm.Ctx, ethcommon.HexToHash(outTxHash))
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
	res3, err := sm.BankClient.SupplyOf(sm.Ctx, &banktypes.QuerySupplyOfRequest{
		Denom: "azeta",
	})
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("$$$ After: SUPPLY OF AZETA: %d", res3.Amount.Amount.BigInt())
	sm.Logger.Info("$$$ Diff: SUPPLY OF AZETA: %d", res3.Amount.Amount.Sub(res2.Amount.Amount).BigInt())
}
