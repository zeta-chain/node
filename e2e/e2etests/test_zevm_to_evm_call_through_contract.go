package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/contracts/gatewayzevmcaller"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestZEVMToEVMCallThroughContract(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	payload := randomPayload(r)

	r.AssertTestDAppEVMCalled(false, payload, big.NewInt(0))

	// deploy caller contract and send it gas zrc20 to pay gas fee
	gatewayCallerAddr, tx, gatewayCaller, err := gatewayzevmcaller.DeployGatewayZEVMCaller(
		r.ZEVMAuth,
		r.ZEVMClient,
		r.GatewayZEVMAddr,
		r.WZetaAddr,
	)
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	tx, err = r.ETHZRC20.Transfer(r.ZEVMAuth, gatewayCallerAddr, big.NewInt(100000000000000000))
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// perform the authenticated call
	tx = r.ZEVMToEMVCallThroughContract(
		gatewayCaller,
		r.TestDAppV2EVMAddr,
		[]byte(payload),
		gatewayzevmcaller.RevertOptions{
			OnRevertGasLimit: big.NewInt(0),
		},
	)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "call")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	r.AssertTestDAppEVMCalled(true, payload, big.NewInt(0))

	// check expected sender was used
	senderForMsg, err := r.TestDAppV2EVM.SenderWithMessage(
		&bind.CallOpts{},
		[]byte(payload),
	)
	require.NoError(r, err)
	require.Equal(r, gatewayCallerAddr, senderForMsg)
}
