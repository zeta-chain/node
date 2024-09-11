package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/precompiles/bank"
)

func TestPrecompilesBank(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	owner, spender := r.ZEVMAuth.From, bank.ContractAddress

	bankContract, err := bank.NewIBank(bank.ContractAddress, r.ZEVMClient)
	require.NoError(r, err, "Failed to create bank contract caller")

	// Get the initial balance of the owner in ERC20ZRC20 tokens.
	ownerERC20InitialBalance, err := r.ERC20ZRC20.BalanceOf(&bind.CallOpts{}, owner)
	require.NoError(r, err, "Error retrieving initial owner balance")

	// Get the balance of the owner in coins "zevm/0x12345".
	// BalanceOf will convert the ZRC20 address to a Cosmos denom formatted as "zevm/0x12345".
	retBalanceOf, err := bankContract.BalanceOf(nil, r.ERC20ZRC20Addr, owner)
	require.NoError(r, err, "Error calling bank.balanceOf")
	require.EqualValues(r, uint64(0), retBalanceOf.Uint64(), "balanceOf result has to be 0")

	// Allow the bank contract to spend 100 ERC20ZRC20 tokens.
	tx, err := r.ERC20ZRC20.Approve(r.ZEVMAuth, spender, big.NewInt(100))
	require.NoError(r, err, "Error approving allowance for bank contract")

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	require.EqualValues(r, 1, receipt.Status, "Error in the approve allowance transaction")

	// Check the allowance of the bank in ERC20ZRC20 tokens. Should be 100.
	balance, err := r.ERC20ZRC20.Allowance(&bind.CallOpts{}, owner, spender)
	require.NoError(r, err, "Error retrieving bank allowance")
	require.EqualValues(r, uint64(100), balance.Uint64(), "Error allowance for bank contract")

	// Call deposit with 100 coins.
	_, err = bankContract.Deposit(r.ZEVMAuth, r.ERC20ZRC20Addr, big.NewInt(100))
	require.NoError(r, err, "Error calling bank.deposit")

	// Check the balance of the owner in coins "zevm/0x12345".
	retBalanceOf, err = bankContract.BalanceOf(nil, r.ERC20ZRC20Addr, owner)
	require.NoError(r, err, "Error calling balanceOf")
	require.EqualValues(r, uint64(100), retBalanceOf.Uint64(), "balanceOf result has to be 100")

	// Check the balance of the owner in r.ERC20ZRC20Addr, should be 100 less.
	ownerERC20FinalBalance, err := r.ERC20ZRC20.BalanceOf(&bind.CallOpts{}, owner)
	require.NoError(r, err, "Error retrieving final owner balance")
	require.EqualValues(
		r,
		ownerERC20InitialBalance.Uint64()-100, // expected
		ownerERC20FinalBalance.Uint64(),       // actual
		"Final balance should be initial - 100",
	)
}
