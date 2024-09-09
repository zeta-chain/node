package e2etests

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestEtherWithdrawRestricted tests the withdrawal to a restricted receiver address
func TestEtherWithdrawRestricted(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	withdrawalAmount, ok := new(big.Int).SetString(args[0], 10)
	require.True(r, ok)

	// approve 1 unit of the gas token to cover the gas fee transfer
	tx, err := r.ETHZRC20.Approve(r.ZEVMAuth, r.ETHZRC20Addr, big.NewInt(1e18))
	require.NoError(r, err)

	r.Logger.EVMTransaction(*tx, "approve")

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	r.Logger.EVMReceipt(*receipt, "approve")

	// withdraw
	restrictedAddress := ethcommon.HexToAddress(sample.RestrictedEVMAddressTest)
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
