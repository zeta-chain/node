package e2etests

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	regularcaller "github.com/zeta-chain/zetacore/precompiles/regular/testutil"
)

func TestPrecompilesRegular(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	_, tx, contract, err := regularcaller.DeployRegularCaller(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err, "Failed to deploy RegularCaller contract")
	r.Logger.Info("Contract deployed")

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	require.Equal(r, receipt.Status, uint64(1), "Failed to deploy RegularCaller contract")
	r.Logger.Info("Deploy transaction successful")

	// Call the Regular contract in the static precompile address.
	ok, err := contract.TestBech32ToHexAddr(&bind.CallOpts{})
	require.NoError(r, err, "Failed to create Regular contract caller")
	require.True(r, ok, "Failed to validate Bech32ToHexAddr function")

	ok, err = contract.TestBech32ify(&bind.CallOpts{})
	require.NoError(r, err, "Failed to create Regular contract caller")
	require.True(r, ok, "Failed to validate Bech32ToHexAddr function")
}
