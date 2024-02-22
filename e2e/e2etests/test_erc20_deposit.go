package e2etests

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

func TestERC20Deposit(r *runner.E2ERunner) {
	hash := r.DepositERC20WithAmountAndMessage(r.DeployerAddress, big.NewInt(100000), []byte{})

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, hash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")

	// deposit ERC20 to banned address
	r.DepositERC20WithAmountAndMessage(ethcommon.HexToAddress(testutils.BannedEVMAddressTest), big.NewInt(100000), []byte{})
}
