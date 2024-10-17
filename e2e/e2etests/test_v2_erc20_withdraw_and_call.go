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

const payloadMessageWithdrawAuthenticatedCallERC20 = "this is a test ERC20 withdraw and authenticated call payload"

func TestV2ERC20WithdrawAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	require.True(r, ok, "Invalid amount specified for TestV2ERC20WithdrawAndCall")

	r.AssertTestDAppEVMCalled(false, payloadMessageWithdrawAuthenticatedCallERC20, amount)

	r.ApproveERC20ZRC20(r.GatewayZEVMAddr)
	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	// set expected sender
	tx, err := r.TestDAppV2EVM.SetExpectedOnCallSender(r.EVMAuth, r.ZEVMAuth.From)
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)

	// perform the withdraw
	tx = r.V2ERC20WithdrawAndCall(
		r.TestDAppV2EVMAddr,
		amount,
		r.EncodeERC20Call(r.ERC20Addr, amount, payloadMessageWithdrawAuthenticatedCallERC20),
		gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	r.AssertTestDAppEVMCalled(true, payloadMessageWithdrawAuthenticatedCallERC20, amount)

	// check expected sender was used
	senderForMsg, err := r.TestDAppV2EVM.SenderWithMessage(
		&bind.CallOpts{},
		[]byte(payloadMessageAuthenticatedWithdrawETH),
	)
	require.NoError(r, err)
	require.Equal(r, r.ZEVMAuth.From, senderForMsg)
}
