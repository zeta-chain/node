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

const payloadMessageAuthenticatedWithdrawETH = "this is a test ETH withdraw and authenticated call payload"

func TestV2ETHWithdrawAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	previousGasLimit := r.ZEVMAuth.GasLimit
	r.ZEVMAuth.GasLimit = 10000000
	defer func() {
		r.ZEVMAuth.GasLimit = previousGasLimit
	}()

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	require.True(r, ok, "Invalid amount specified for TestV2ETHWithdrawAndCall")

	r.AssertTestDAppEVMCalled(false, payloadMessageAuthenticatedWithdrawETH, amount)

	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	// perform the withdraw
	tx := r.V2ETHWithdrawAndCall(
		r.TestDAppV2EVMAddr,
		amount,
		[]byte(payloadMessageAuthenticatedWithdrawETH),
		gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	r.AssertTestDAppEVMCalled(true, payloadMessageAuthenticatedWithdrawETH, amount)

	// check expected sender was used
	senderForMsg, err := r.TestDAppV2EVM.SenderWithMessage(
		&bind.CallOpts{},
		[]byte(payloadMessageAuthenticatedWithdrawETH),
	)
	require.NoError(r, err)
	require.Equal(r, r.ZEVMAuth.From, senderForMsg)
}
