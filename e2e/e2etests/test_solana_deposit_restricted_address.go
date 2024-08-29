package e2etests

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
)

func TestSolanaDepositRestricted(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// parse restricted address
	receiverRestricted := ethcommon.HexToAddress(args[0])

	// parse deposit amount (in lamports)
	depositAmount := parseBigInt(r, args[1])

	// execute the deposit transaction
	r.SOLDepositAndCall(nil, receiverRestricted, depositAmount, nil)
}
