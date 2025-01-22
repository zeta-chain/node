package e2etests

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/x/crosschain/types"
	zetabitcoin "github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
)

func TestBitcoinDepositAndAbortWithLowDepositFee(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	// ARRANGE
	// Given small amount
	depositAmount := zetabitcoin.DefaultDepositorFee - float64(1)/btcutil.SatoshiPerBitcoin

	// ACT
	txHash := r.DepositBTCWithAmount(depositAmount, nil, false)

	// ASSERT
	// cctx status should be aborted
	cctx := utils.WaitCctxAbortedByInboundHash(r.Ctx, r, txHash.String(), r.CctxClient)
	r.Logger.CCTX(cctx, "deposit aborted")

	// check cctx details
	require.Equal(r, cctx.InboundParams.Amount.Uint64(), uint64(0))
	require.Equal(r, cctx.GetCurrentOutboundParam().Amount.Uint64(), uint64(0))

	// check cctx error message
	require.Equal(r, cctx.CctxStatus.ErrorMessage, types.InboundStatus_insufficient_depositor_fee.String())
}
