package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
)

func TestSolanaWithdraw(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// print balance of from address
	solZRC20 := r.SOLZRC20
	balance, err := solZRC20.BalanceOf(&bind.CallOpts{}, r.ZEVMAuth.From)
	require.NoError(r, err)
	r.Logger.Info("from address %s balance of SOL before: %d", r.ZEVMAuth.From, balance)

	// parse withdraw amount (in lamports), approve amount is 1 SOL
	approvedAmount := new(big.Int).SetUint64(solana.LAMPORTS_PER_SOL)
	withdrawAmount, ok := new(big.Int).SetString(args[0], 10)
	require.True(r, ok, "Invalid withdrawal amount specified for TestSolanaWithdraw.")
	require.Equal(
		r,
		-1,
		withdrawAmount.Cmp(approvedAmount),
		"Withdrawal amount must be less than the approved amount (1e9).",
	)

	// load deployer private key
	privkey := solana.MustPrivateKeyFromBase58(r.Account.SolanaPrivateKey.String())

	// withdraw
	r.WithdrawSOLZRC20(privkey.PublicKey(), withdrawAmount)

	// print balance of from address after withdraw
	balance, err = solZRC20.BalanceOf(&bind.CallOpts{}, r.ZEVMAuth.From)
	require.NoError(r, err)
	r.Logger.Info("from address %s balance of SOL after: %d", r.ZEVMAuth.From, balance)
}
