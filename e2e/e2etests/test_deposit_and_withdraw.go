package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestDepositAndWithdraw makes a depositAndCall that trigger a withdrawal to the origin chain
func TestDepositAndWithdraw(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	require.True(r, ok, "Invalid amount specified for TestV2ETHDepositAndCall")

	withdrawMessage, err := r.TestDAppV2ZEVM.WITHDRAW(&bind.CallOpts{})
	require.NoError(r, err)

	// perform the deposit and call to the TestDAppV2ZEVMAddr
	tx := r.V2ETHDepositAndCall(
		r.TestDAppV2ZEVMAddr,
		amount,
		[]byte(withdrawMessage),
		gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)

	// wait for the cctx to be mined
	cctxDeposit := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctxDeposit, "deposit")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctxDeposit.CctxStatus.Status)

	// first cctx should trigger a new cctx for the withdrawal
	cctxWithdraw := utils.WaitCctxMinedByInboundHash(r.Ctx, cctxDeposit.Index, r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctxWithdraw, "withdraw")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctxWithdraw.CctxStatus.Status)
}
