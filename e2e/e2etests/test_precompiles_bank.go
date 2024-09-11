package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/precompiles/bank"
)

func TestPrecompilesBank(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	bankContract, err := bank.NewIBank(bank.ContractAddress, r.ZEVMClient)
	require.NoError(r, err, "Failed to create bank contract caller")

	// Get the balance of the user_precompile in coins "zevm/0x12345"
	// BalanceOf will convert the ZRC20 address to a Cosmos denom formatted as "zevm/0x12345".
	res, err := bankContract.BalanceOf(nil, r.ERC20ZRC20Addr, r.ZEVMAuth.From)
	require.NoError(r, err, "Error calling BalanceOf")
	require.Equal(r, big.NewInt(0), res, "BalanceOf result has to be 0")
}
