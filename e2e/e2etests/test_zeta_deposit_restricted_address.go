package e2etests

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

func TestZetaDepositRestricted(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	require.True(r, ok, "Invalid amount specified for TestZetaDepositRestricted.")

	// Deposit amount to restricted address
	r.DepositZetaWithAmount(ethcommon.HexToAddress(testutils.RestrictedEVMAddressTest), amount)
}
