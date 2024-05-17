package e2etests

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

// TestEtherWithdrawRestricted tests the withdrawal to a restricted receiver address
func TestEtherWithdrawRestricted(r *runner.E2ERunner, args []string) {
	approvedAmount := big.NewInt(1e18)
	if len(args) != 1 {
		panic("TestEtherWithdrawRestricted requires exactly one argument for the withdrawal amount.")
	}

	withdrawalAmount, ok := new(big.Int).SetString(args[0], 10)
	if !ok {
		panic("Invalid withdrawal amount specified for TestEtherWithdrawRestricted.")
	}

	if withdrawalAmount.Cmp(approvedAmount) >= 0 {
		panic("Withdrawal amount must be less than the approved amount (1e18).")
	}

	// approve
	tx, err := r.ETHZRC20.Approve(r.ZEVMAuth, r.ETHZRC20Addr, approvedAmount)
	if err != nil {
		panic(err)
	}
	r.Logger.EVMTransaction(*tx, "approve")

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	r.Logger.EVMReceipt(*receipt, "approve")

	// withdraw
	restrictedAddress := ethcommon.HexToAddress(testutils.RestrictedEVMAddressTest)
	tx, err = r.ETHZRC20.Withdraw(r.ZEVMAuth, restrictedAddress.Bytes(), withdrawalAmount)
	if err != nil {
		panic(err)
	}
	r.Logger.EVMTransaction(*tx, "withdraw to restricted address")

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("withdraw failed")
	}
	r.Logger.EVMReceipt(*receipt, "withdraw")
	r.Logger.ZRC20Withdrawal(r.ETHZRC20, *receipt, "withdraw")

	// verify the withdraw value
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, receipt.TxHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		panic("cctx status is not outbound mined")
	}

	// the cctx should be cancelled with zero value
	verifyTransferAmountFromCCTX(r, cctx, 0)
}
