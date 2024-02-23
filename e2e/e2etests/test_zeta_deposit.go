package e2etests

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

func TestZetaDeposit(r *runner.E2ERunner) {
	// Deposit 1 Zeta
	hash := r.DepositZetaWithAmount(r.DeployerAddress, big.NewInt(1e18))

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, hash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")

	// Deposit 1 Zeta to restricted address
	r.DepositZetaWithAmount(ethcommon.HexToAddress(testutils.RestrictedEVMAddressTest), big.NewInt(1e18))
}
