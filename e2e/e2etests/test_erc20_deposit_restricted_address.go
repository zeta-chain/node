package e2etests

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

func TestERC20DepositRestricted(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	require.True(r, ok)

	// deposit ERC20 to restricted address
	r.DepositERC20WithAmountAndMessage(ethcommon.HexToAddress(testutils.RestrictedEVMAddressTest), amount, []byte{})
}
