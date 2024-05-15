package e2etests

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	cctxtypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// TestMessagePassingRevertFailExternalChains tests message passing with failing revert between external EVM chains
// TODO: Use two external EVM chains for these tests
// https://github.com/zeta-chain/node/issues/2185
func TestMessagePassingRevertFailExternalChains(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestMessagePassingRevertFail requires exactly one argument for the amount.")
	}

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestMessagePassingRevertFail.")
	}

	chainID, err := r.EVMClient.ChainID(r.Ctx)
	if err != nil {
		panic(err)
	}

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
		Message:             []byte("revert"), // non-empty message will cause revert, because the dest address is not a contract
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

	// expect revert tx to fail
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	receipt, err = r.EVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(cctx.GetCurrentOutTxParam().OutboundTxHash))
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
