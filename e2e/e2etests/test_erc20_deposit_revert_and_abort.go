package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/e2e/contracts/testabort"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestERC20DepositRevertAndAbort(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	r.ApproveERC20OnEVM(r.GatewayEVMAddr)

	// deploy testabort contract
	testAbortAddr, txDeploy, testAbort, err := testabort.DeployTestAbort(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	r.WaitForTxReceiptOnZEVM(txDeploy)

	// perform the deposit
	tx := r.ERC20DepositAndCall(
		sample.EthAddress(), // non-existing address
		big.NewInt(
			1,
		), // a very small amount is passed so the cctx will be aborted as the fee for reverts cannot be paid
		[]byte("revert"),
		gatewayevm.RevertOptions{
			RevertAddress:    r.TestDAppV2EVMAddr,
			CallOnRevert:     true,
			RevertMessage:    []byte("revert"),
			OnRevertGasLimit: big.NewInt(200000),
			AbortAddress:     testAbortAddr,
		},
	)

	// wait for the cctx to be reverted
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit_and_call")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Aborted)

	// check onAbort was called
	aborted, err := testAbort.IsAborted(&bind.CallOpts{})
	require.NoError(r, err)
	require.True(r, aborted)

	// check abort context was passed
	abortContext, err := testAbort.GetAbortedWithMessage(&bind.CallOpts{}, "revert")
	require.NoError(r, err)
	require.EqualValues(r, r.ERC20ZRC20Addr.Hex(), abortContext.Asset.Hex())

	// wait for the zrc20 balance to be updated
	change := utils.NewBalanceChange(true)
	utils.WaitAndVerifyZRC20BalanceChange(r, r.ERC20ZRC20, testAbortAddr, big.NewInt(0), change, r.Logger)
}
