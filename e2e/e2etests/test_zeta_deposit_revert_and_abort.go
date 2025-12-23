package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/e2e/contracts/testabort"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestZetaDepositRevertAndAbort(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	r.ApproveZetaOnEVM(r.GatewayEVMAddr)

	// deploy testabort contract
	testAbortAddr, txDeploy, testAbort, err := testabort.DeployTestAbort(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	r.WaitForTxReceiptOnZEVM(txDeploy)

	// perform the deposit
	// Deposit (Fails as the address is non-existing)
	// Revert (Fails as the amount is too small to pay for the revert fee)
	// Zeta deposited and onAbort called
	tx := r.ZetaDepositAndCall(
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

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "zeta_deposit_revert_and_abort")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Aborted)

	if r.IsV2ZETAEnabled() {
		// V2 ZETA flows enabled: deposit should fail, revert should fail, then abort with onAbort called

		// check onAbort was called
		aborted, err := testAbort.IsAborted(&bind.CallOpts{})
		require.NoError(r, err)
		require.True(r, aborted)

		// Asset is empty as ZETA is the native gas token on ZEVM
		emptyAddress := ethcommon.Address{}
		// check abort context was passed
		abortContext, err := testAbort.GetAbortedWithMessage(&bind.CallOpts{}, "revert")
		require.NoError(r, err)
		require.EqualValues(r, emptyAddress.Hex(), abortContext.Asset.Hex())

		// check abort contract received the tokens
		balance, err := r.ZEVMClient.BalanceAt(r.Ctx, testAbortAddr, nil)
		require.NoError(r, err)
		require.True(r, balance.Uint64() > 0)
	} else {
		// V2 ZETA flows disabled: deposit should be aborted with ErrZetaThroughGateway
		require.Equal(r, cctx.CctxStatus.StatusMessage, crosschaintypes.ErrZetaThroughGateway.Error())
	}
}
