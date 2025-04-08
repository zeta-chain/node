package e2etests

import (
	"math/big"
	"strings"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestSolanaDepositAndCallRevertWithDust tests Solana deposit and call that reverts with a dust amount in the revert outbound.
func TestSolanaDepositAndCallRevertWithDust(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	// deposit the rent exempt amount which will result in a dust amount (after fee deduction) in the revert outbound
	depositAmount := big.NewInt(constant.SolanaWalletRentExempt)

	// ACT
	// execute the deposit and call transaction
	nonExistReceiver := sample.EthAddress()
	data := []byte("dust lamports should abort cctx")
	sig := r.SOLDepositAndCall(nil, nonExistReceiver, depositAmount, data, nil)

	// ASSERT
	// Now we want to make sure cctx is aborted.
	cctx := utils.WaitCctxAbortedByInboundHash(r.Ctx, r, sig.String(), r.CctxClient)
	require.True(r, cctx.GetCurrentOutboundParam().Amount.Uint64() < constant.SolanaWalletRentExempt)
	require.True(
		r,
		strings.Contains(cctx.CctxStatus.ErrorMessageRevert, crosschaintypes.ErrInvalidWithdrawalAmount.Error()),
	)
}
