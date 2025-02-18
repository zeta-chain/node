package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// TestUpdateOperationalChainParams tests updating the operational chain params for a chain
func TestUpdateOperationalChainParams(r *runner.E2ERunner, _ []string) {
	chainIDEVM, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)
	chainID := chainIDEVM.Int64()

	// get old chain params
	res, err := r.ObserverClient.GetChainParams(r.Ctx, &observertypes.QueryGetChainParamsRequest{})
	require.NoError(r, err)

	oldChainParams := observertypes.ChainParams{}
	found := false
	for _, chainParams := range res.ChainParams.ChainParams {
		if chainParams.ChainId == chainID {
			oldChainParams = *chainParams
			found = true
			break
		}
	}
	require.True(r, found, "chain params not found")

	r.Logger.Info("Updating operational chain parameters")
	msg := observertypes.NewMsgUpdateOperationalChainParams(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		chainID,
		oldChainParams.GasPriceTicker+1,
		oldChainParams.InboundTicker+1,
		oldChainParams.OutboundTicker+1,
		oldChainParams.WatchUtxoTicker+1,
		oldChainParams.OutboundScheduleInterval+1,
		oldChainParams.OutboundScheduleLookahead+1,
		observertypes.ConfirmationParams{
			FastInboundCount:  oldChainParams.ConfirmationParams.FastInboundCount + 1,
			FastOutboundCount: oldChainParams.ConfirmationParams.FastOutboundCount + 1,
			SafeInboundCount:  oldChainParams.ConfirmationParams.SafeInboundCount + 1,
			SafeOutboundCount: oldChainParams.ConfirmationParams.SafeOutboundCount + 1,
		},
	)
	resTx, err := r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, msg)
	require.NoError(r, err)

	r.Logger.Info("Update operational flags tx hash: %s", resTx.TxHash)

	// fetch again the flags
	res, err = r.ObserverClient.GetChainParams(r.Ctx, &observertypes.QueryGetChainParamsRequest{})
	require.NoError(r, err)

	found = false
	for _, chainParams := range res.ChainParams.ChainParams {
		if chainParams.ChainId == chainID {
			require.Equal(r, oldChainParams.GasPriceTicker+1, chainParams.GasPriceTicker)
			require.Equal(r, oldChainParams.InboundTicker+1, chainParams.InboundTicker)
			require.Equal(r, oldChainParams.OutboundTicker+1, chainParams.OutboundTicker)
			require.Equal(r, oldChainParams.WatchUtxoTicker+1, chainParams.WatchUtxoTicker)
			require.Equal(r, oldChainParams.OutboundScheduleInterval+1, chainParams.OutboundScheduleInterval)
			require.Equal(r, oldChainParams.OutboundScheduleLookahead+1, chainParams.OutboundScheduleLookahead)
			require.Equal(
				r,
				oldChainParams.ConfirmationParams.FastInboundCount+1,
				chainParams.ConfirmationParams.FastInboundCount,
			)
			require.Equal(
				r,
				oldChainParams.ConfirmationParams.FastOutboundCount+1,
				chainParams.ConfirmationParams.FastOutboundCount,
			)
			require.Equal(
				r,
				oldChainParams.ConfirmationParams.SafeInboundCount+1,
				chainParams.ConfirmationParams.SafeInboundCount,
			)
			require.Equal(
				r,
				oldChainParams.ConfirmationParams.SafeOutboundCount+1,
				chainParams.ConfirmationParams.SafeOutboundCount,
			)
			found = true
			break
		}
	}
	require.True(r, found, "chain params not found")
}
