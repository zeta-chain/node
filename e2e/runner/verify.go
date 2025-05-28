package runner

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"

	solanacontracts "github.com/zeta-chain/node/pkg/contracts/solana"
)

// EVMVerifyOutboundTransferAmount verifies the transfer amount on EVM chain for given outbound hash
func (r *E2ERunner) EVMVerifyOutboundTransferAmount(outboundHash string, amount int64) {
	receipt, err := r.EVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(outboundHash))
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

// SolanaVerifyWithdrawalAmount verifies the withdrawn amount on Solana for given outbound hash
func (r *E2ERunner) SolanaVerifyWithdrawalAmount(outboundHash string, amount uint64) {
	// convert txHash to signature
	sig, err := solana.SignatureFromBase58(outboundHash)
	require.NoError(r, err)

	// query transaction by signature
	txResult, err := r.SolanaClient.GetTransaction(r.Ctx, sig, &rpc.GetTransactionOpts{
		Commitment: rpc.CommitmentConfirmed,
	})
	require.NoError(r, err)

	// unmarshal transaction
	tx, err := txResult.Transaction.GetTransaction()
	require.NoError(r, err)

	// 1st instruction is the withdraw
	instruction := tx.Message.Instructions[0]
	instWithdraw, err := solanacontracts.ParseInstructionWithdraw(instruction)
	require.NoError(r, err)

	// verify the amount
	require.Equal(r, amount, instWithdraw.TokenAmount(), "withdraw amount is not correct")
}
