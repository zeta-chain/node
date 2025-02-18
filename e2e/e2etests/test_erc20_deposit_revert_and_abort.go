package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/e2e/contracts/testabort"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestERC20DepositRevertAndAbort(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	r.ApproveERC20OnEVM(r.GatewayEVMAddr)

	// deploy testabort contract
	testAbortAddr, _, testAbort, err := testabort.DeployTestAbort(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)

	// perform the deposit
	tx := r.ERC20DepositAndCall(
		sample.EthAddress(), // non-existing address
		amount, []byte("revert"),
		gatewayevm.RevertOptions{
			RevertAddress:    sample.EthAddress(), // non-existing address
			CallOnRevert:     true,
			RevertMessage:    []byte("revert"),
			OnRevertGasLimit: big.NewInt(200000),
			AbortAddress:     testAbortAddr,
		})

	// wait for the cctx to be reverted
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit_and_call")
	require.Equal(r, crosschaintypes.CctxStatus_Aborted, cctx.CctxStatus.Status)

	// check onAbort was called
	aborted, err := testAbort.IsAborted(&bind.CallOpts{})
	require.NoError(r, err)
	require.True(r, aborted)

	// check abort context was passed
	abortContext, err := testAbort.GetAbortedWithMessage(&bind.CallOpts{}, "revert")
	require.NoError(r, err)
	require.EqualValues(r, r.ERC20ZRC20Addr.Hex(), abortContext.Asset.Hex())

	// check abort contract received the tokens
	balance, err := r.ERC20ZRC20.BalanceOf(&bind.CallOpts{}, testAbortAddr)
	require.NoError(r, err)
	require.True(r, balance.Uint64() > 0)
}
