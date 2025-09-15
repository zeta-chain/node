package e2etests

import (
	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/testutil/sample"
)

// TestSuiDepositRestrictedAddress tests a deposit to a restricted address that won't be observed by the observers
func TestSuiDepositRestrictedAddress(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)
	amount := utils.ParseBigInt(r, args[0])

	// ARRANGE
	// Given restricted receiver
	receiver := ethcommon.HexToAddress(sample.RestrictedEVMAddressTest)

	// balance before
	oldBalance, err := r.SUIZRC20.BalanceOf(&bind.CallOpts{}, receiver)
	require.NoError(r, err)

	// ACT
	// perform the deposit
	resp := r.SuiDepositSUI(r.SuiGateway.PackageID(), receiver, math.NewUintFromBigInt(amount))
	r.Logger.Info("Sui restricted deposit tx: %s", resp.Digest)

	// wait enough time
	r.WaitForBlocks(5)

	// no cctx should be created
	utils.EnsureNoCctxMinedByInboundHash(r.Ctx, resp.Digest, r.CctxClient)

	// receiver balance should not change
	newBalance, err := r.SUIZRC20.BalanceOf(&bind.CallOpts{}, receiver)
	require.NoError(r, err)
	require.Equal(r, oldBalance.Uint64(), newBalance.Uint64())
}
