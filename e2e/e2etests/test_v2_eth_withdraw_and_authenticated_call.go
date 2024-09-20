package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	testgatewayzevmcaller "github.com/zeta-chain/node/pkg/contracts/testgatewayzevmcaller"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

const payloadMessageAuthenticatedWithdrawETH = "this is a test ETH withdraw and authenticated call payload"
const payloadMessageAuthenticatedWithdrawETHThroughContract = "this is a test ETH withdraw and authenticated call payload through contract"

func TestV2ETHWithdrawAndAuthenticatedCall(r *runner.E2ERunner, args []string) {
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

	// set expected sender
	tx, err := r.TestDAppV2EVM.SetExpectedOnCallSender(r.EVMAuth, r.ZEVMAuth.From)
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)

	// perform the withdraw
	tx = r.V2ETHWithdrawAndAuthenticatedCall(
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
	senderForMsg, err := r.TestDAppV2EVM.SenderWithMessage(&bind.CallOpts{}, []byte(payloadMessageAuthenticatedWithdrawETH))
	require.NoError(r, err)
	require.Equal(r, r.ZEVMAuth.From, senderForMsg)

	// deploy caller contract and send it gas zrc20 to pay gas fee
	gatewayCallerAddr, tx, gatewayCaller, err := testgatewayzevmcaller.DeployTestGatewayZEVMCaller(r.ZEVMAuth, r.ZEVMClient, r.GatewayZEVMAddr, r.WZetaAddr)
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	tx, err = r.ETHZRC20.Transfer(r.ZEVMAuth, gatewayCallerAddr, big.NewInt(100000000000000000))
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// set expected sender
	tx, err = r.TestDAppV2EVM.SetExpectedOnCallSender(r.EVMAuth, gatewayCallerAddr)
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)

	// perform the authenticated call
	tx = r.V2ETHWithdrawAndAuthenticatedCallThroughContract(gatewayCaller, r.TestDAppV2EVMAddr,
		amount,
		[]byte(payloadMessageAuthenticatedWithdrawETHThroughContract),
		testgatewayzevmcaller.RevertOptions{OnRevertGasLimit: big.NewInt(0)})

	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	r.AssertTestDAppEVMCalled(true, payloadMessageAuthenticatedWithdrawETHThroughContract, amount)

	// check expected sender was used
	senderForMsg, err = r.TestDAppV2EVM.SenderWithMessage(&bind.CallOpts{}, []byte(payloadMessageAuthenticatedWithdrawETHThroughContract))
	require.NoError(r, err)
	require.Equal(r, gatewayCallerAddr, senderForMsg)

	// set expected sender to wrong one
	tx, err = r.TestDAppV2EVM.SetExpectedOnCallSender(r.EVMAuth, r.ZEVMAuth.From)
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)

	// repeat authenticated call through contract, should revert because of wrong sender
	tx = r.V2ETHWithdrawAndAuthenticatedCallThroughContract(gatewayCaller, r.TestDAppV2EVMAddr,
		amount,
		[]byte(payloadMessageAuthenticatedWithdrawETHThroughContract),
		testgatewayzevmcaller.RevertOptions{OnRevertGasLimit: big.NewInt(0)})

	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	require.Equal(r, crosschaintypes.CctxStatus_Reverted, cctx.CctxStatus.Status)
}
