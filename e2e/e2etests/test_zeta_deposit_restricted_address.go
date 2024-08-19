package e2etests

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

func TestZetaDepositRestricted(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse the deposit amount
	amount := parseBigInt(r, args[0])

	// Deposit amount to restricted address
	r.DepositZetaWithAmount(ethcommon.HexToAddress(testutils.RestrictedEVMAddressTest), amount)
}
