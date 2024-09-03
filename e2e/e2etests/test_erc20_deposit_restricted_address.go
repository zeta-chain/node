package e2etests

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/zetaclient/testutils"
)

func TestERC20DepositRestricted(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse the deposit amount
	amount := parseBigInt(r, args[0])

	// deposit ERC20 to restricted address
	r.DepositERC20WithAmountAndMessage(ethcommon.HexToAddress(testutils.RestrictedEVMAddressTest), amount, []byte{})
}
