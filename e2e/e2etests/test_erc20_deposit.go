package e2etests

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

func TestERC20Deposit(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestERC20Deposit requires exactly one argument for the amount.")
	}

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestERC20Deposit.")
	}

	hash := r.DepositERC20WithAmountAndMessage(r.DeployerAddress, amount, []byte{})

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, hash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
}

func TestERC20DepositRestricted(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestERC20DepositRestricted requires exactly one argument for the amount.")
	}

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestERC20DepositRestricted.")
	}

	// deposit ERC20 to restricted address
	r.DepositERC20WithAmountAndMessage(ethcommon.HexToAddress(testutils.RestrictedEVMAddressTest), amount, []byte{})
}
