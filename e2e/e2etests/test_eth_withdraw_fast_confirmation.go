package e2etests

import (
	"math/big"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// TestETHWithdrawFastConfirmation tests the fast confirmation of ETH withdrawal
func TestETHWithdrawFastConfirmation(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// ARRANGE
	// parse amount
	amount := utils.ParseBigInt(r, args[0])

	// query chainID
	chainIDBig, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)
	chainID := chainIDBig.Int64()

	// enable outbound fast confirmation by updating the chain params
	reqQuery := &observertypes.QueryGetChainParamsForChainRequest{ChainId: chainID}
	resOldChainParams, err := r.ObserverClient.GetChainParamsForChain(r.Ctx, reqQuery)
	require.NoError(r, err)

	chainParams := *resOldChainParams.ChainParams
	chainParams.ConfirmationParams = &observertypes.ConfirmationParams{
		SafeInboundCount:  resOldChainParams.ChainParams.ConfirmationParams.SafeInboundCount,
		FastInboundCount:  resOldChainParams.ChainParams.ConfirmationParams.FastInboundCount,
		SafeOutboundCount: 10, // approx 10 seconds, much longer than Fast confirmation time (1 second)
		FastOutboundCount: 1,
	}
	err = r.ZetaTxServer.UpdateChainParams(&chainParams)
	require.NoError(r, err, "failed to enable outbound fast confirmation")

	// it takes 1 Zeta block time for zetaclient to pick up the new chain params
	// wait for 2 blocks to ensure the new chain params are effective
	r.WaitForBlocks(2)
	r.Logger.Info("enabled outbound fast confirmation")

	// ACT-1
	// perform the withdraw and wait for 1 Zeta block for CCTX creation
	tx := r.ETHWithdraw(r.EVMAddress(), amount, gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)})
	r.WaitForBlocks(1)

	// wait for outbound tracker submission
	// TSS keysign takes time. In order to measure confirmation time accurately, we need to wait for outbound tracker submission.
	cctx := utils.GetCCTXByInboundHash(r.Ctx, r.CctxClient, tx.Hash().Hex())
	nonce := cctx.GetCurrentOutboundParam().TssNonce
	hashes := utils.WaitOutboundTracker(r.Ctx, r.CctxClient, chainID, nonce, 1, r.Logger, 1*time.Minute)
	r.Logger.Info("outbound (fast) tracker created: %s", hashes[0])

	// ASSERT-1
	// wait for the cctx to be FAST confirmed
	timeStart := time.Now()
	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
	require.Equal(r, crosschaintypes.ConfirmationMode_FAST, cctx.GetCurrentOutboundParam().ConfirmationMode)
	fastConfirmTime := time.Since(timeStart)
	r.Logger.Info("FAST confirmed withdrawal succeeded in %f seconds", fastConfirmTime.Seconds())

	// ACT-2
	// disable outbound fast confirmation by setting FastOutboundCount to SafeOutboundCount
	chainParams.ConfirmationParams.FastOutboundCount = chainParams.ConfirmationParams.SafeOutboundCount
	err = r.ZetaTxServer.UpdateChainParams(&chainParams)
	require.NoError(r, err, "failed to disable outbound fast confirmation")

	// it takes 1 Zeta block time for zetaclient to pick up the new chain params
	// wait for 2 blocks to ensure the new chain params are effective
	r.WaitForBlocks(2)
	r.Logger.Info("disabled outbound fast confirmation")

	// perform the withdraw again and wait for 1 Zeta block for CCTX creation
	tx = r.ETHWithdraw(r.EVMAddress(), amount, gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)})
	r.WaitForBlocks(1)

	// wait for outbound tracker submission
	cctx = utils.GetCCTXByInboundHash(r.Ctx, r.CctxClient, tx.Hash().Hex())
	nonce = cctx.GetCurrentOutboundParam().TssNonce
	hashes = utils.WaitOutboundTracker(r.Ctx, r.CctxClient, chainID, nonce, 1, r.Logger, 1*time.Minute)
	r.Logger.Info("outbound (safe) tracker created: %s", hashes[0])

	// ASSERT-2
	// wait for the cctx to be SAFE confirmed
	timeStart = time.Now()
	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
	require.Equal(r, crosschaintypes.ConfirmationMode_SAFE, cctx.GetCurrentOutboundParam().ConfirmationMode)
	safeConfirmTime := time.Since(timeStart)
	r.Logger.Info("SAFE confirmed withdrawal succeeded in %f seconds", safeConfirmTime.Seconds())

	// ensure FAST confirmation is faster than SAFE confirmation
	// using 3 seconds is good enough to check the difference on local goerli network
	timeSaved := safeConfirmTime - fastConfirmTime
	r.Logger.Info("FAST confirmation saved %f seconds", timeSaved.Seconds())
	require.True(r, timeSaved > 3*time.Second)

	// TEARDOWN
	// restore old outbound confirmation params
	// note: we should NOT restore 'resOldChainParams' as it may interfere with fast confirmation tests on deposits
	chainParams.ConfirmationParams.SafeOutboundCount = resOldChainParams.ChainParams.ConfirmationParams.SafeOutboundCount
	chainParams.ConfirmationParams.FastOutboundCount = resOldChainParams.ChainParams.ConfirmationParams.FastOutboundCount
	err = r.ZetaTxServer.UpdateChainParams(&chainParams)
	require.NoError(r, err, "failed to restore chain params")
}
