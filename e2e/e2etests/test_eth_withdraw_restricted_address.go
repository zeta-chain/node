package e2etests

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

// TestEtherWithdrawRestricted tests the withdrawal to a restricted receiver address
func TestEtherWithdrawRestricted(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	approvedAmount := big.NewInt(1e18)

	withdrawalAmount, ok := new(big.Int).SetString(args[0], 10)
	require.True(r, ok)
	require.True(
		r,
		withdrawalAmount.Cmp(approvedAmount) <= 0,
		"Withdrawal amount must be less than the approved amount (1e18)",
	)

	// approve
	tx, err := r.ETHZRC20.Approve(r.ZEVMAuth, r.ETHZRC20Addr, approvedAmount)
	require.NoError(r, err)

	r.Logger.EVMTransaction(*tx, "approve")

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	r.Logger.EVMReceipt(*receipt, "approve")

	// withdraw
	restrictedAddress := ethcommon.HexToAddress(testutils.RestrictedEVMAddressTest)
	tx, err = r.ETHZRC20.Withdraw(r.ZEVMAuth, restrictedAddress.Bytes(), withdrawalAmount)
	require.NoError(r, err)

	r.Logger.EVMTransaction(*tx, "withdraw to restricted address")

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	r.Logger.EVMReceipt(*receipt, "withdraw")
	r.Logger.ZRC20Withdrawal(r.ETHZRC20, *receipt, "withdraw")

	// verify the withdrawal value
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, receipt.TxHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")

	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// the cctx should be cancelled with zero value
	verifyTransferAmountFromCCTX(r, cctx, 0)
}
