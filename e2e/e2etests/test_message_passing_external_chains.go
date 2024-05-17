package e2etests

import (
	"fmt"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	cctxtypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// TestMessagePassingExternalChains tests message passing between external EVM chains
// TODO: Use two external EVM chains for these tests
// https://github.com/zeta-chain/node/issues/2185
func TestMessagePassingExternalChains(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestMessagePassing requires exactly one argument for the amount.")
	}

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestMessagePassing.")
	}

	chainID, err := r.EVMClient.ChainID(r.Ctx)
	if err != nil {
		panic(err)
	}

	r.Logger.Info("Approving ConnectorEth to spend deployer's ZetaEth")
	auth := r.EVMAuth
	tx, err := r.ZetaEth.Approve(auth, r.ConnectorEthAddr, amount)
	if err != nil {
		panic(err)
	}

	r.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("tx failed")
	}
	r.Logger.Info("Approve tx receipt: %d", receipt.Status)
	r.Logger.Info("Calling ConnectorEth.Send")
	tx, err = r.ConnectorEth.Send(auth, zetaconnectoreth.ZetaInterfacesSendInput{
		DestinationChainId:  chainID,
		DestinationAddress:  r.DeployerAddress.Bytes(),
		DestinationGasLimit: big.NewInt(400_000),
		Message:             nil,
		ZetaValueAndGas:     amount,
		ZetaParams:          nil,
	})
	if err != nil {
		panic(err)
	}

	r.Logger.Info("ConnectorEth.Send tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("tx failed")
	}
	r.Logger.Info("ConnectorEth.Send tx receipt: status %d", receipt.Status)
	r.Logger.Info("  Logs:")
	for _, log := range receipt.Logs {
		sentLog, err := r.ConnectorEth.ParseZetaSent(*log)
		if err == nil {
			r.Logger.Info("    Dest Addr: %s", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
			r.Logger.Info("    Dest Chain: %d", sentLog.DestinationChainId)
			r.Logger.Info("    Dest Gas: %d", sentLog.DestinationGasLimit)
			r.Logger.Info("    Zeta Value: %d", sentLog.ZetaValueAndGas)
		}
	}

	r.Logger.Info("Waiting for ConnectorEth.Send CCTX to be mined...")
	r.Logger.Info("  INTX hash: %s", receipt.TxHash.String())
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx.CctxStatus.Status != cctxtypes.CctxStatus_OutboundMined {
		panic(fmt.Sprintf(
			"expected cctx status to be %s; got %s, message %s",
			cctxtypes.CctxStatus_OutboundMined,
			cctx.CctxStatus.Status.String(),
			cctx.CctxStatus.StatusMessage,
		))
	}
	receipt, err = r.EVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(cctx.GetCurrentOutTxParam().OutboundTxHash))
	if err != nil {
		panic(err)
	}
	if receipt.Status != 1 {
		panic("tx failed")
	}
	for _, log := range receipt.Logs {
		event, err := r.ConnectorEth.ParseZetaReceived(*log)
		if err == nil {
			r.Logger.Info("Received ZetaSent event:")
			r.Logger.Info("  Dest Addr: %s", event.DestinationAddress)
			r.Logger.Info("  Zeta Value: %d", event.ZetaValue)
			r.Logger.Info("  src chainid: %d", event.SourceChainId)
			if event.ZetaValue.Cmp(cctx.GetCurrentOutTxParam().Amount.BigInt()) != 0 {
				panic("Zeta value mismatch")
			}
		}
	}
}
