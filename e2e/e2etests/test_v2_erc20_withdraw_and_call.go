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

func TestV2ERC20WithdrawAndCall(r *runner.E2ERunner, _ []string) {
	previousGasLimit := r.ZEVMAuth.GasLimit
	r.ZEVMAuth.GasLimit = 10000000
	defer func() {
		r.ZEVMAuth.GasLimit = previousGasLimit
	}()

	// called with 0 amount since onCall implementation is for TestDappV2 is simple and generic without decoding the payload and amount handling for erc20
	// and purpose of test is to verify that onCall is called with correct sender and payload
	amount := big.NewInt(0)

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
		[]byte(payloadMessageWithdrawAuthenticatedCallERC20),
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
