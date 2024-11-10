package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
)

// TestSPLWithdrawAndCreateReceiverAta withdraws spl, but letting gateway to create receiver ata using rent payer
// instead of providing receiver that has it already created
func TestSPLWithdrawAndCreateReceiverAta(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	withdrawAmount := parseBigInt(r, args[0])

	// get SPL ZRC20 balance before withdraw
	zrc20BalanceBefore, err := r.SPLZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SPL before withdraw: %d", zrc20BalanceBefore)

	require.Equal(r, 1, zrc20BalanceBefore.Cmp(withdrawAmount), "Insufficient balance for withdrawal")

	// parse withdraw amount (in lamports), approve amount is 1 SOL
	approvedAmount := new(big.Int).SetUint64(solana.LAMPORTS_PER_SOL)
	require.Equal(
		r,
		-1,
		withdrawAmount.Cmp(approvedAmount),
		"Withdrawal amount must be less than the approved amount (1e9)",
	)

	// create new priv key, with empty ata
	receiverPrivKey, err := solana.NewRandomPrivateKey()
	require.NoError(r, err)

	// verify receiver ata account doesn't exist
	receiverAta, _, err := solana.FindAssociatedTokenAddress(receiverPrivKey.PublicKey(), r.SPLAddr)
	require.NoError(r, err)

	receiverAtaAcc, err := r.SolanaClient.GetAccountInfo(r.Ctx, receiverAta)
	require.Error(r, err)
	require.Nil(r, receiverAtaAcc)

	// withdraw
	r.WithdrawSPLZRC20(receiverPrivKey.PublicKey(), withdrawAmount, approvedAmount)

	// get SPL ZRC20 balance after withdraw
	zrc20BalanceAfter, err := r.SPLZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SPL after withdraw: %d", zrc20BalanceAfter)

	// verify receiver ata was created
	receiverAtaAcc, err = r.SolanaClient.GetAccountInfo(r.Ctx, receiverAta)
	require.NoError(r, err)
	require.NotNil(r, receiverAtaAcc)

	// verify balances are updated
	receiverBalanceAfter, err := r.SolanaClient.GetTokenAccountBalance(r.Ctx, receiverAta, rpc.CommitmentConfirmed)
	require.NoError(r, err)
	r.Logger.Info("receiver balance of SPL after withdraw: %s", receiverBalanceAfter.Value.Amount)

	// verify amount is added to receiver ata
	require.Zero(r, withdrawAmount.Cmp(parseBigInt(r, receiverBalanceAfter.Value.Amount)))

	// verify amount is subtracted on zrc20
	require.Zero(r, new(big.Int).Sub(zrc20BalanceBefore, withdrawAmount).Cmp(zrc20BalanceAfter))
}
