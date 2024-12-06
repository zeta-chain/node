package runner

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"

	solanacontracts "github.com/zeta-chain/node/pkg/contracts/solana"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// VerifyTransferAmountFromCCTX verifies the transfer amount from the CCTX on EVM
func (r *E2ERunner) VerifyTransferAmountFromCCTX(cctx *crosschaintypes.CrossChainTx, amount int64) {
	r.Logger.Info("outTx hash %s", cctx.GetCurrentOutboundParam().Hash)

	receipt, err := r.EVMClient.TransactionReceipt(
		r.Ctx,
		ethcommon.HexToHash(cctx.GetCurrentOutboundParam().Hash),
	)
	require.NoError(r, err)

	r.Logger.Info("Receipt txhash %s status %d", receipt.TxHash, receipt.Status)

	for _, log := range receipt.Logs {
		event, err := r.ERC20.ParseTransfer(*log)
		if err != nil {
			continue
		}
		r.Logger.Info("  logs: from %s, to %s, value %d", event.From.Hex(), event.To.Hex(), event.Value)
		require.Equal(r, amount, event.Value.Int64(), "value is not correct")
	}
}

// VerifySolanaWithdrawalAmountFromCCTX verifies the withdrawn amount on Solana for given CCTX
func (r *E2ERunner) VerifySolanaWithdrawalAmountFromCCTX(cctx *crosschaintypes.CrossChainTx, amount uint64) {
	txHash := cctx.GetCurrentOutboundParam().Hash
	r.Logger.Info("outbound hash %s", txHash)

	// convert txHash to signature
	sig, err := solana.SignatureFromBase58(txHash)
	require.NoError(r, err)

	// query transaction by signature
	txResult, err := r.SolanaClient.GetTransaction(r.Ctx, sig, &rpc.GetTransactionOpts{})
	require.NoError(r, err)

	// unmarshal transaction
	tx, err := txResult.Transaction.GetTransaction()
	require.NoError(r, err)

	// 1st instruction is the withdraw
	instruction := tx.Message.Instructions[0]
	instWithdrae, err := solanacontracts.ParseInstructionWithdraw(instruction)
	require.NoError(r, err)

	// verify the amount
	require.Equal(r, amount, instWithdrae.TokenAmount(), "withdraw amount is not correct")
}
