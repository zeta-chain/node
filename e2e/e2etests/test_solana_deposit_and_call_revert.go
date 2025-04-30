package e2etests

import (
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"

	testcontract "github.com/zeta-chain/node/e2e/contracts/reverter"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	solanacontracts "github.com/zeta-chain/node/pkg/contracts/solana"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestSolanaDepositAndCallRevert tests deposit of lamports calling a example contract that reverts.
func TestSolanaDepositAndCallRevert(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse deposit amount (in lamports)
	depositAmount := utils.ParseBigInt(r, args[0])

	// deploy a reverter contract in ZEVM
	// TODO: consider removing repeated deployments of reverter contract
	reverterAddr, _, _, err := testcontract.DeployReverter(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	r.Logger.Info("Reverter contract deployed at: %s", reverterAddr.String())

	// execute the deposit transaction
	data := []byte("hello reverter")
	revertAddressPrivateKey, err := solana.NewRandomPrivateKey()
	require.NoError(r, err)

	sig := r.SOLDepositAndCall(nil, reverterAddr, depositAmount, data, &solanacontracts.RevertOptions{
		RevertAddress: revertAddressPrivateKey.PublicKey(),
	})

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "solana_deposit_and_refund")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)
	require.Equal(r, cctx.GetCurrentOutboundParam().Receiver, revertAddressPrivateKey.PublicKey().String())

	require.Contains(r, cctx.CctxStatus.ErrorMessage, utils.ErrHashRevertFoo)

	// verify that revert address balance is increased
	revertAddressBalance, err := r.SolanaClient.GetBalance(
		r.Ctx,
		revertAddressPrivateKey.PublicKey(),
		rpc.CommitmentConfirmed,
	)
	require.NoError(r, err)
	require.Greater(r, revertAddressBalance.Value, uint64(0))
}
