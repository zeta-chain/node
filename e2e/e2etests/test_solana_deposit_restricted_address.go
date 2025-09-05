package e2etests

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

func TestSolanaDepositRestricted(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// parse restricted address
	receiverRestricted := ethcommon.HexToAddress(args[0])

	// parse deposit amount (in lamports)
	depositAmount := utils.ParseBigInt(r, args[1])

	// execute the deposit transaction
	sig := r.SOLDepositAndCall(nil, receiverRestricted, depositAmount, nil, nil)

	// wait for 5 zeta blocks
	r.WaitForBlocks(5)

	// no cctx should be created
	utils.EnsureNoCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient)
}
