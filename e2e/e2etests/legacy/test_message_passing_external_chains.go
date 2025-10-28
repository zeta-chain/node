package legacy

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts-evm/pkg/zetaconnector.eth.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	cctxtypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestMessagePassingExternalChains tests message passing between external EVM chains
// TODO: Use two external EVM chains for these tests
// https://github.com/zeta-chain/node/issues/2185
func TestMessagePassingExternalChains(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse the amount
	amount := utils.ParseBigInt(r, args[0])

	chainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	r.Logger.Info("Approving ConnectorEth to spend deployer's ZetaEth")
	auth := r.EVMAuth

	tx, err := r.ZetaEth.Approve(auth, r.ConnectorEthAddr, amount)
	require.NoError(r, err)

	r.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	r.Logger.Info("Approve tx receipt: %d", receipt.Status)
	r.Logger.Info("Calling ConnectorEth.Send")
	tx, err = r.ConnectorEth.Send(auth, zetaconnectoreth.ZetaInterfacesSendInput{
		DestinationChainId:  chainID,
		DestinationAddress:  r.EVMAddress().Bytes(),
		DestinationGasLimit: big.NewInt(400_000),
		Message:             nil,
		ZetaValueAndGas:     amount,
		ZetaParams:          nil,
	})
	require.NoError(r, err)

	r.Logger.Info("ConnectorEth.Send tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

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

	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, cctxtypes.CctxStatus_OutboundMined)

	receipt, err = r.EVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(cctx.GetCurrentOutboundParam().Hash))
	require.NoError(r, err)
	utils.RequireTxSuccessful(r, receipt)

	for _, log := range receipt.Logs {
		event, err := r.ConnectorEth.ParseZetaReceived(*log)
		if err == nil {
			r.Logger.Info("Received ZetaSent event:")
			r.Logger.Info("  Dest Addr: %s", event.DestinationAddress)
			r.Logger.Info("  Zeta Value: %d", event.ZetaValue)
			r.Logger.Info("  src chainid: %d", event.SourceChainId)

			comp := event.ZetaValue.Cmp(cctx.GetCurrentOutboundParam().Amount.BigInt())
			require.Equal(r, 0, comp, "Zeta value mismatch")
		}
	}
}
