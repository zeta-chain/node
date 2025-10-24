package legacy

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts-evm/pkg/zetaconnector.eth.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	cctxtypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestMessagePassingRevertFailExternalChains tests message passing with failing revert between external EVM chains
// TODO: Use two external EVM chains for these tests
// https://github.com/zeta-chain/node/issues/2185
func TestMessagePassingRevertFailExternalChains(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse the amount
	amount := utils.ParseBigInt(r, args[0])

	chainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

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
		Message: []byte(
			"revert",
		), // non-empty message will cause revert, because the dest address is not a contract
		ZetaValueAndGas: amount,
		ZetaParams:      nil,
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
	// expect revert tx to fail
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	receipt, err = r.EVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(cctx.GetCurrentOutboundParam().Hash))
	require.NoError(r, err)

	// expect revert tx to fail as well
	require.Equal(r, ethtypes.ReceiptStatusFailed, receipt.Status)
	utils.RequireCCTXStatus(r, cctx, cctxtypes.CctxStatus_Aborted)
}
