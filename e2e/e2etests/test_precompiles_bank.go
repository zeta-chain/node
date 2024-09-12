package e2etests

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/precompiles/bank"
)

func TestPrecompilesBank(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	// Increase the gasLimit. It's required because of the gas consumed by precompiled functions.
	previousGasLimit := r.ZEVMAuth.GasLimit
	r.ZEVMAuth.GasLimit = 10_000_000
	defer func() {
		r.ZEVMAuth.GasLimit = previousGasLimit
	}()

	spender, bankAddr := r.EVMAddress(), bank.ContractAddress

	// Deposit and approve 50 WZETA for the test.
	approveAmount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(50))
	r.DepositAndApproveWZeta(approveAmount)
	fmt.Printf("DEBUG: approveAmount %s\n", approveAmount.String())

	initialBalance, err := r.WZeta.BalanceOf(&bind.CallOpts{Context: r.Ctx}, spender)
	fmt.Printf("DEBUG: initialBalance %s\n", initialBalance.String())
	require.NoError(r, err, "Error approving allowance for bank contract")
	require.EqualValues(r, approveAmount, initialBalance, "spender balance should be 50")

	// Create a bank contract caller.
	bankContract, err := bank.NewIBank(bank.ContractAddress, r.ZEVMClient)
	require.NoError(r, err, "Failed to create bank contract caller")

	// Get the balance of the spender in coins "zevm/WZetaAddr". This calls bank.balanceOf().
	// BalanceOf will convert the ZRC20 address to a Cosmos denom formatted as "zevm/0x12345".
	retBalanceOf, err := bankContract.BalanceOf(&bind.CallOpts{Context: r.Ctx}, r.WZetaAddr, spender)
	fmt.Printf("DEBUG: initial bank.balanceOf() %s\n", retBalanceOf.String())
	require.NoError(r, err, "Error calling bank.balanceOf()")
	require.EqualValues(r, uint64(0), retBalanceOf.Uint64(), "Initial cosmos coins balance has to be 0")

	// Allow the bank contract to spend 25 WZeta tokens.
	tx, err := r.WZeta.Approve(r.ZEVMAuth, bankAddr, big.NewInt(25))
	require.NoError(r, err, "Error approving allowance for bank contract")
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	require.EqualValues(r, uint64(1), receipt.Status, "approve allowance tx failed")

	// Check the allowance of the bank in WZeta tokens. Should be 25.
	balance, err := r.WZeta.Allowance(&bind.CallOpts{Context: r.Ctx}, spender, bankAddr)
	fmt.Printf("DEBUG: bank allowance %s\n", balance.String())
	require.NoError(r, err, "Error retrieving bank allowance")
	require.EqualValues(r, uint64(25), balance.Uint64(), "Error allowance for bank contract")

	// Call deposit with 25 coins.
	tx, err = bankContract.Deposit(r.ZEVMAuth, r.WZetaAddr, big.NewInt(25))
	fmt.Printf("DEBUG: bank.deposit() tx hash %s\n", tx.Hash().String())
	require.NoError(r, err, "Error calling bank.deposit()")

	r.Logger.Info("Waiting for 5 blocks")
	r.WaitForBlocks(5)
	fmt.Printf("DEBUG: bank.deposit() tx %+v\n", tx)

	// Check the balance of the spender in coins "zevm/WZetaAddr".
	retBalanceOf, err = bankContract.BalanceOf(&bind.CallOpts{Context: r.Ctx}, r.WZetaAddr, spender)
	fmt.Printf("DEBUG: final bank.balanceOf() tx %s\n", retBalanceOf.String())
	require.NoError(r, err, "Error calling bank.balanceOf()")
	require.EqualValues(r, uint64(25), retBalanceOf.Uint64(), "balanceOf result has to be 25")

	// Check the balance of the spender in r.WZeta, should be 100 less.
	finalBalance, err := r.WZeta.BalanceOf(&bind.CallOpts{Context: r.Ctx}, spender)
	fmt.Printf("DEBUG: final WZeta balance %s\n", finalBalance.String())
	require.NoError(r, err, "Error retrieving final owner balance")
	require.EqualValues(
		r,
		initialBalance.Uint64()-25, // expected
		finalBalance.Uint64(),      // actual
		"Final balance should be initial - 25",
	)
}
