package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
)

func TestSolanaWithdraw(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// get ERC20 SOL balance before withdraw
	balanceBefore, err := r.SOLZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SOL before withdraw: %d", balanceBefore)

	// parse withdraw amount (in lamports), approve amount is 1 SOL
	approvedAmount := new(big.Int).SetUint64(solana.LAMPORTS_PER_SOL)
	withdrawAmount := parseBigInt(r, args[0])
	require.Equal(
		r,
		-1,
		withdrawAmount.Cmp(approvedAmount),
		"Withdrawal amount must be less than the approved amount (1e9).",
	)

	// load deployer private key
	privkey, err := solana.PrivateKeyFromBase58(r.Account.SolanaPrivateKey.String())
	require.NoError(r, err)

	// withdraw
	r.WithdrawSOLZRC20(privkey.PublicKey(), withdrawAmount, approvedAmount)

	// get ERC20 SOL balance after withdraw
	balanceAfter, err := r.SOLZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SOL after withdraw: %d", balanceAfter)

	// check if the balance is reduced correctly
	amountReduced := new(big.Int).Sub(balanceBefore, balanceAfter)
	require.True(r, amountReduced.Cmp(withdrawAmount) >= 0, "balance is not reduced correctly")
}
