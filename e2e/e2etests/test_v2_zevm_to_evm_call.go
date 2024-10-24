package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

const payloadMessageEVMAuthenticatedCall = "this is a test EVM authenticated call payload"

func TestV2ZEVMToEVMCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	r.AssertTestDAppEVMCalled(false, payloadMessageEVMAuthenticatedCall, big.NewInt(0))

	// necessary approval for fee payment
	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	// perform the authenticated call
	tx := r.V2ZEVMToEMVCall(
		r.TestDAppV2EVMAddr,
		[]byte(payloadMessageEVMAuthenticatedCall),
		gatewayzevm.RevertOptions{
			OnRevertGasLimit: big.NewInt(0),
		},
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "call")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	// check the payload was received on the contract
	r.AssertTestDAppEVMCalled(true, payloadMessageEVMAuthenticatedCall, big.NewInt(0))

	// check expected sender was used
	senderForMsg, err := r.TestDAppV2EVM.SenderWithMessage(&bind.CallOpts{}, []byte(payloadMessageEVMAuthenticatedCall))
	require.NoError(r, err)
	require.Equal(r, r.ZEVMAuth.From, senderForMsg)
}
