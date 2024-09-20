package e2etests

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	testgatewayzevmcaller "github.com/zeta-chain/node/pkg/contracts/testgatewayzevmcaller"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

const payloadMessageEVMCall = "this is a test EVM call payload"
const payloadMessageEVMAuthenticatedCall = "this is a test EVM authenticated call payload"

func TestV2ZEVMToEVMCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	r.AssertTestDAppEVMCalled(false, payloadMessageEVMCall, big.NewInt(0))

	// Necessary approval for fee payment
	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	// perform the call
	tx := r.V2ZEVMToEMVCall(r.TestDAppV2EVMAddr, r.EncodeSimpleCall(payloadMessageEVMCall), gatewayzevm.RevertOptions{
		OnRevertGasLimit: big.NewInt(0),
	})

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "call")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	// check the payload was received on the contract
	r.AssertTestDAppEVMCalled(true, payloadMessageEVMCall, big.NewInt(0))
}

func TestV2ZEVMToEVMAuthenticatedCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	r.AssertTestDAppEVMCalled(false, payloadMessageEVMAuthenticatedCall, big.NewInt(0))

	// Necessary approval for fee payment
	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	// perform the authenticated call
	tx := r.V2ZEVMToEMVAuthenticatedCall(r.TestDAppV2EVMAddr, []byte(payloadMessageEVMAuthenticatedCall), gatewayzevm.RevertOptions{
		OnRevertGasLimit: big.NewInt(0),
	})

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "call")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	// check the payload was received on the contract
	r.AssertTestDAppEVMCalled(true, payloadMessageEVMAuthenticatedCall, big.NewInt(0))
	s, _ := r.TestDAppV2EVM.Senders(&bind.CallOpts{}, big.NewInt(0))
	fmt.Println("sender", s.Hex())
	gatewayCallerAddr, tx, gatewayCaller, err := testgatewayzevmcaller.DeployTestGatewayZEVMCaller(r.ZEVMAuth, r.ZEVMClient, r.GatewayZEVMAddr)
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	r.ZEVMAuth.GasLimit = 10000000

	tx, err = r.ETHZRC20.Transfer(r.ZEVMAuth, gatewayCallerAddr, big.NewInt(100000000000000000))
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	tx = r.V2ZEVMToEMVAuthenticatedCallThroughContract(gatewayCaller, r.TestDAppV2EVMAddr, []byte(payloadMessageEVMAuthenticatedCall), testgatewayzevmcaller.RevertOptions{
		OnRevertGasLimit: big.NewInt(0),
	})
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "call")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
	s, _ = r.TestDAppV2EVM.Senders(&bind.CallOpts{}, big.NewInt(1))
	fmt.Println("sender", s.Hex())
}
