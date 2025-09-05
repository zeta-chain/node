package e2etests

import (
	"strings"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	zetabitcoin "github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
)

// TestBitcoinDepositAndCallRevertWithDust sends a Bitcoin deposit that reverts with a dust amount in the revert outbound.
func TestBitcoinDepositAndCallRevertWithDust(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	// 0.002 BTC is consumed in a deposit and revert, the dust is set to 1000 satoshis in the protocol
	// Therefore the deposit amount should be 0.002 + 0.000001 = 0.00200100 should trigger the condition
	// As only 100 satoshis are left after the deposit
	const (
		// depositAmount is 0.002001 BTC, chosen to result in 100 satoshis after deposit
		// which is below the dust threshold of 1000 satoshis
		depositAmount = 0.00200100
	)
	amount := depositAmount + zetabitcoin.DefaultDepositorFee

	// ACT
	// Send BTC to TSS address with a dummy memo
	// zetacore should revert cctx if call is made on a non-existing address
	nonExistReceiver := sample.EthAddress()
	anyMemo := append(nonExistReceiver.Bytes(), []byte("gibberish-memo")...)

	// One UTXO is enough to cover the deposit amount
	txHash, err := r.SendToTSSWithMemo(amount, anyMemo)
	require.NoError(r, err)
	require.NotEmpty(r, txHash)

	// ASSERT
	// Now we want to make sure the cctx is aborted with expected error message

	// cctx status would be pending revert if using v21 or before
	cctx := utils.WaitCctxAbortedByInboundHash(r.Ctx, r, txHash.String(), r.CctxClient)

	require.True(r, cctx.GetCurrentOutboundParam().Amount.Uint64() < constant.BTCWithdrawalDustAmount)
	require.True(
		r,
		strings.Contains(cctx.CctxStatus.ErrorMessageRevert, crosschaintypes.ErrInvalidWithdrawalAmount.Error()),
	)
}
