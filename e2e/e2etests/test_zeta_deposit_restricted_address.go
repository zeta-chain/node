package e2etests

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

func TestZetaDepositRestricted(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestZetaDepositRestricted requires exactly one argument for the amount.")
	}

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestZetaDepositRestricted.")
	}

	// Deposit amount to restricted address
	r.DepositZetaWithAmount(ethcommon.HexToAddress(testutils.RestrictedEVMAddressTest), amount)
}
