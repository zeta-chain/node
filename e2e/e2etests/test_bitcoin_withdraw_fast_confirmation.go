package e2etests

import (
	"time"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// TestBitcoinWithdrawFastConfirmation tests the fast confirmation of Bitcoin withdrawal
func TestBitcoinWithdrawFastConfirmation(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// ARRANGE
	// parse arguments
	chainID := r.GetBitcoinChainID()
	defaultReceiver := r.BTCDeployerAddress.EncodeAddress()
	receiver, amount := utils.ParseBitcoinWithdrawArgs(r, args, defaultReceiver, chainID)

	// enable outbound fast confirmation by updating the chain params
	reqQuery := &observertypes.QueryGetChainParamsForChainRequest{ChainId: chainID}
	resOldChainParams, err := r.ObserverClient.GetChainParamsForChain(r.Ctx, reqQuery)
	require.NoError(r, err)

	chainParams := *resOldChainParams.ChainParams
	chainParams.ConfirmationParams = &observertypes.ConfirmationParams{
		SafeInboundCount:  resOldChainParams.ChainParams.ConfirmationParams.SafeInboundCount,
		FastInboundCount:  resOldChainParams.ChainParams.ConfirmationParams.FastInboundCount,
		SafeOutboundCount: 6, // approx 36 seconds, much longer than Fast confirmation time (6 second)
		FastOutboundCount: 1,
	}
	err = r.ZetaTxServer.UpdateChainParams(&chainParams)
	require.NoError(r, err, "failed to enable outbound fast confirmation")

	// it takes 1 Zeta block time for zetaclient to pick up the new chain params
	// wait for 2 blocks to ensure the new chain params are effective
	utils.WaitForZetaBlocks(r.Ctx, r, r.ZEVMClient, 2, 20*time.Second)
	r.Logger.Info("enabled outbound fast confirmation")

	// ACT-1
	// perform the withdraw
	txResult := withdrawBTCZRC20(r, receiver, amount)

	// ASSERT-1
	// wait for the cctx to be FAST confirmed
	timeStart := time.Now()
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txResult.Txid, r.CctxClient, r.Logger, r.CctxTimeout)
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
	require.Equal(r, crosschaintypes.ConfirmationMode_FAST, cctx.InboundParams.ConfirmationMode)
	fastConfirmTime := time.Since(timeStart)

	r.Logger.Info("FAST confirmed withdrawal succeeded in %f seconds", fastConfirmTime.Seconds())

	// ACT-2
	// restore old chain params; disable outbound fast confirmation
	err = r.ZetaTxServer.UpdateChainParams(resOldChainParams.ChainParams)
	require.NoError(r, err, "failed to restore chain params")

	// perform the withdraw again
	txResult = withdrawBTCZRC20(r, receiver, amount)

	// ASSERT-2
	// wait for the cctx to be SAFE confirmed
	timeStart = time.Now()
	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, txResult.Txid, r.CctxClient, r.Logger, r.CctxTimeout)
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
	require.Equal(r, crosschaintypes.ConfirmationMode_SAFE, cctx.InboundParams.ConfirmationMode)
	safeConfirmTime := time.Since(timeStart)

	r.Logger.Info("SAFE confirmed withdrawal succeeded in %f seconds", safeConfirmTime.Seconds())

	// ensure FAST confirmation is faster than SAFE confirmation
	// using one BTC block time is good enough to check the difference
	timeSaved := safeConfirmTime - fastConfirmTime
	r.Logger.Info("FAST confirmation saved %f seconds", timeSaved.Seconds())
	require.True(r, timeSaved > runner.BTCRegnetBlockTime)
}
