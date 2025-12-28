package e2etests

import (
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/mod/semver"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/constant"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

const (
	startTimestampMetricName        = "zetaclient_last_start_timestamp_seconds"
	blockTimeLatencyMetricName      = "zetaclient_core_block_latency"
	blockTimeLatencySleepMetricName = "zetaclient_core_block_latency_sleep"
)

// TestZetaclientRestartHeight tests scheduling a zetaclient restart via operational flags
func TestZetaclientRestartHeight(r *runner.E2ERunner, _ []string) {
	_, err := r.Clients.Zetacore.Observer.OperationalFlags(
		r.Ctx,
		&observertypes.QueryOperationalFlagsRequest{},
	)
	require.NoError(r, err)

	currentHeight, err := r.Clients.Zetacore.GetBlockHeight(r.Ctx)
	require.NoError(r, err)

	// schedule a restart for 5 blocks in the future
	restartHeight := currentHeight + 5
	updateMsg := observertypes.NewMsgUpdateOperationalFlags(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		observertypes.OperationalFlags{
			RestartHeight: restartHeight,
		},
	)

	_, err = r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, updateMsg)
	require.NoError(r, err)

	operationalFlagsRes, err := r.Clients.Zetacore.Observer.OperationalFlags(
		r.Ctx,
		&observertypes.QueryOperationalFlagsRequest{},
	)
	require.NoError(r, err)
	require.Equal(r, restartHeight, operationalFlagsRes.OperationalFlags.RestartHeight)

	originalStartTime, err := r.Clients.ZetaclientMetrics.FetchGauge(startTimestampMetricName)
	require.NoError(r, err, "fetching zetaclient metric name")

	// wait for height above restart height
	// wait for a few extra block to account for shutdown and startup time
	require.Eventually(r, func() bool {
		height, err := r.Clients.Zetacore.GetBlockHeight(r.Ctx)
		require.NoError(r, err)
		return height > restartHeight+3
	}, time.Minute, time.Second)

	currentStartTime, err := r.Clients.ZetaclientMetrics.FetchGauge(startTimestampMetricName)
	require.NoError(r, err)

	require.Greater(r, currentStartTime, originalStartTime+1)
}

// TestZetaclientSignerOffset tests scheduling a zetaclient restart via operational flags
func TestZetaclientSignerOffset(r *runner.E2ERunner, _ []string) {
	startBlockTimeLatencySleep, err := r.Clients.ZetaclientMetrics.FetchGauge(blockTimeLatencySleepMetricName)
	require.NoError(r, err)
	require.InDelta(r, 0, startBlockTimeLatencySleep, .01, "start block time latency should be 0")

	// get starting block time latency.
	// we need to ensure it's not zero (if zetaclient just finished a restart)
	var startBlockTimeLatency float64
	require.Eventually(r, func() bool {
		startBlockTimeLatency, err = r.Clients.ZetaclientMetrics.FetchGauge(blockTimeLatencyMetricName)
		require.NoError(r, err)
		return startBlockTimeLatency > 1
	}, time.Second*15, time.Millisecond*100)

	desiredSignerBlockTimeOffset := time.Duration(startBlockTimeLatency*float64(time.Second)) + time.Millisecond*200

	updateMsg := observertypes.NewMsgUpdateOperationalFlags(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		observertypes.OperationalFlags{
			SignerBlockTimeOffset: &desiredSignerBlockTimeOffset,
		},
	)

	_, err = r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, updateMsg)
	require.NoError(r, err)

	operationalFlagsRes, err := r.Clients.Zetacore.Observer.OperationalFlags(
		r.Ctx,
		&observertypes.QueryOperationalFlagsRequest{},
	)
	require.NoError(r, err)
	require.InDelta(r, desiredSignerBlockTimeOffset, *(operationalFlagsRes.OperationalFlags.SignerBlockTimeOffset), .01)

	require.Eventually(r, func() bool {
		blockTimeLatencySleep, err := r.Clients.ZetaclientMetrics.FetchGauge(blockTimeLatencySleepMetricName)
		if err != nil {
			return false
		}
		return blockTimeLatencySleep > .05
	}, time.Second*20, time.Second*1)
}

// TestZetaclientMinimumVersion tests setting the zetaclient minimum version
func TestZetaclientMinimumVersion(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)
	minimumVersion := args[0]

	// use zetacore version by default
	if minimumVersion == "" {
		nodeInfo, err := r.Clients.Zetacore.GetNodeInfo(r.Ctx)
		require.NoError(r, err)
		minimumVersion = constant.NormalizeVersion(nodeInfo.ApplicationVersion.Version)
	}

	// validate version string
	require.True(r, semver.IsValid(minimumVersion), "invalid version: %s", minimumVersion)

	// query operational flags
	oldFlags, err := r.Clients.Zetacore.Observer.OperationalFlags(
		r.Ctx,
		&observertypes.QueryOperationalFlagsRequest{},
	)
	require.NoError(r, err)

	// update minimum version field
	newFlags := oldFlags.OperationalFlags
	newFlags.MinimumVersion = minimumVersion

	// send update tx message to zetacore
	updateMsg := observertypes.NewMsgUpdateOperationalFlags(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		newFlags,
	)
	_, err = r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, updateMsg)
	require.NoError(r, err)

	// query operational flags again
	currentFlags, err := r.Clients.Zetacore.Observer.OperationalFlags(
		r.Ctx,
		&observertypes.QueryOperationalFlagsRequest{},
	)
	require.NoError(r, err)
	require.EqualValues(r, newFlags, currentFlags.OperationalFlags)

	r.Logger.Print("set zetaclient minimum version to %s", minimumVersion)
}
