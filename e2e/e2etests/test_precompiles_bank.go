package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/precompiles/bank"
)

func TestPrecompilesBank(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	bankContract, err := bank.NewIBank(bank.ContractAddress, r.ZEVMClient)
	require.NoError(r, err, "Failed to create bank contract caller")

	// Get the balance of the user_precompile in coins "zevm/0x12345"
	// BalanceOf will convert the ZRC20 address to a Cosmos denom formatted as "zevm/0x12345".
	retBalanceOf, err := bankContract.BalanceOf(nil, r.ERC20ZRC20Addr, r.ZEVMAuth.From)
	require.NoError(r, err, "Error calling balanceOf")
	require.EqualValues(r, uint64(0), retBalanceOf.Uint64(), "balanceOf result has to be 0")

	// Allow the bank contract to spend 100 coins
	tx, err := r.ERC20ZRC20.Approve(r.ZEVMAuth, bank.ContractAddress, big.NewInt(100))
	require.NoError(r, err, "Error approving bank contract")

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	require.EqualValues(r, 1, receipt.Status, "Error approving allowance for bank contract")

	// Call deposit with 100 coins
	// _, err = bankContract.Deposit(r.ZEVMAuth, r.ERC20ZRC20Addr, big.NewInt(0))
	// require.NoError(r, err, "Error calling deposit")

	// Check the balance of the user_precompile in coins "zevm/0x12345"

	// Check the balance of the user_precompile in r.ERC20ZRC20Addr

}
