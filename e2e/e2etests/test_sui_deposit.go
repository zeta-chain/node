package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

func TestSuiDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])
	_ = amount

	r.SUIDeposit(r.EVMAddress())
}
