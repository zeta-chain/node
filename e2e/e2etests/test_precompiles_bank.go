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

	previousGasLimit := r.ZEVMAuth.GasLimit
	r.ZEVMAuth.GasLimit = 10_000_000
	defer func() {
		r.ZEVMAuth.GasLimit = previousGasLimit
	}()

	// Set owner and spender for legibility.
	owner, spender := r.EVMAddress(), bank.ContractAddress
	fmt.Println("owner ", owner.String())
	fmt.Println("spender ", spender.String())
	fmt.Println("ERC20ZRC20 ", r.ERC20ZRC20Addr.String())

	// Fund owner with 200 token.
	tx, err := r.ERC20ZRC20.Transfer(r.ZEVMAuth, owner, big.NewInt(200))
	require.NoError(r, err, "Error funding owner with ERC20ZRC20")

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	fmt.Printf("funding owner tx receipt: %+v\n", receipt)
	utils.RequireTxSuccessful(r, receipt, "funding owner tx")

	// Create a bank contract caller.
	bankContract, err := bank.NewIBank(bank.ContractAddress, r.ZEVMClient)
	require.NoError(r, err, "Failed to create bank contract caller")

	// Get the initial balance of the owner in ERC20ZRC20 tokens. Should be 200.
	ownerERC20InitialBalance, err := r.ERC20ZRC20.BalanceOf(&bind.CallOpts{Context: r.Ctx}, owner)
	require.NoError(r, err, "Error retrieving initial owner balance")
	require.EqualValues(r, uint64(0), ownerERC20InitialBalance.Uint64(), "Initial ERC20ZRC20 has to be 200")
	fmt.Println("owner balance ERC20: ", ownerERC20InitialBalance)

	// Get the balance of the owner in coins "zevm/0x12345". This calls bank.balanceOf().
	// BalanceOf will convert the ZRC20 address to a Cosmos denom formatted as "zevm/0x12345".
	retBalanceOf, err := bankContract.BalanceOf(&bind.CallOpts{Context: r.Ctx}, r.ERC20ZRC20Addr, owner)
	require.NoError(r, err, "Error calling bank.balanceOf()")
	require.EqualValues(r, uint64(0), retBalanceOf.Uint64(), "Initial cosmos coins balance has to be 0")
	fmt.Println("owner balance zevm/coin: ", retBalanceOf)

	// Allow the bank contract to spend 100 ERC20ZRC20 tokens.
	tx, err = r.ERC20ZRC20.Approve(r.ZEVMAuth, spender, big.NewInt(100))
	require.NoError(r, err, "Error approving allowance for bank contract")
	fmt.Printf("approve allowance tx: %+v\n", tx)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "approve allowance tx")
	fmt.Printf("approve allowance tx receipt: %+v\n", receipt)

	// Check the allowance of the bank in ERC20ZRC20 tokens. Should be 100.
	balance, err := r.ERC20ZRC20.Allowance(&bind.CallOpts{Context: r.Ctx}, owner, spender)
	require.NoError(r, err, "Error retrieving bank allowance")
	require.EqualValues(r, uint64(100), balance.Uint64(), "Error allowance for bank contract")
	fmt.Printf("bank allowance: %v\n", balance)

	// Call deposit with 100 coins.
	tx, err = bankContract.Deposit(r.ZEVMAuth, r.ERC20ZRC20Addr, big.NewInt(100))
	require.NoError(r, err, "Error calling bank.deposit")
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	fmt.Printf("Deposit tx: %+v\n", tx)

	// Check the balance of the owner in coins "zevm/0x12345".
	retBalanceOf, err = bankContract.BalanceOf(nil, r.ERC20ZRC20Addr, owner)
	require.NoError(r, err, "Error calling balanceOf")
	require.EqualValues(r, uint64(100), retBalanceOf.Uint64(), "balanceOf result has to be 100")
	fmt.Printf("owner balance zevm/coin (should increase): %+v\n", retBalanceOf)

	// Check the balance of the owner in r.ERC20ZRC20Addr, should be 100 less.
	ownerERC20FinalBalance, err := r.ERC20ZRC20.BalanceOf(&bind.CallOpts{Context: r.Ctx}, owner)
	require.NoError(r, err, "Error retrieving final owner balance")
	fmt.Printf("owner final ERC20 balance (should decrease): %+v\n", retBalanceOf)
	require.EqualValues(
		r,
		ownerERC20InitialBalance.Uint64()-100, // expected
		ownerERC20FinalBalance.Uint64(),       // actual
		"Final balance should be initial - 100",
	)
}
