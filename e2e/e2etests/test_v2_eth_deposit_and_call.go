package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

const payloadMessageETH = "this is a test ETH deposit and call payload"

func TestV2ETHDepositAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	require.True(r, ok, "Invalid amount specified for TestV2ETHDepositAndCall")

	// perform the deposit and call to the TestDAppV2ZEVMAddr
	tx := r.V2ETHDepositAndCall(r.TestDAppV2ZEVMAddr, amount, []byte(payloadMessageETH))

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit_and_call")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	// check the payload was received on the contract
	message, err := r.TestDAppV2ZEVM.LastMessage(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, payloadMessageETH, message)

	// check the amount was received on the contract
	amountReceived, err := r.TestDAppV2ZEVM.LastAmount(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, amount.String(), amountReceived.String())
}
