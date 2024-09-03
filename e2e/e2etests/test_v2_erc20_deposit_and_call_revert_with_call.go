package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

const payloadMessageDepositOnRevertERC20 = "this is a test ERC20 deposit and call on revert"

func TestV2ERC20DepositAndCallRevertWithCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	require.True(r, ok, "Invalid amount specified for TestV2ERC20DepositAndCallRevertWithCall")

	r.ApproveERC20OnEVM(r.GatewayEVMAddr)

	r.AssertTestDAppEVMCalled(false, payloadMessageDepositOnRevertERC20, amount)

	// perform the deposit
	tx := r.V2ERC20DepositAndCall(r.TestDAppV2ZEVMAddr, amount, []byte("revert"), gatewayevm.RevertOptions{
		RevertAddress:    r.TestDAppV2EVMAddr,
		CallOnRevert:     true,
		RevertMessage:    []byte(payloadMessageDepositOnRevertERC20),
		OnRevertGasLimit: big.NewInt(200000),
	})

	// wait for the cctx to be reverted
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	require.Equal(r, crosschaintypes.CctxStatus_Reverted, cctx.CctxStatus.Status)

	// check the payload was received on the contract
	r.AssertTestDAppEVMCalled(true, payloadMessageDepositOnRevertERC20, big.NewInt(0))
}
