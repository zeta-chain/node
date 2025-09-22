package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// TestCriticalAdminTransactions tests critical admin transactions that are the most used on mainnet .
// The complete list is
// MsgUpdateChainParams
// MsgRefundAbortedCCTX
// MsgEnableCCTX
// MsgDisableCCTX
// MsgUpdateGasPriceIncreaseFlags
// MsgAddInboundTracker
// MsgUpdateZRC20LiquidityCap
// MsgDeploySystemContracts
// MsgWhitelistERC20
// MsgPauseZRC20
// MsgMigrateTssFunds
// MsgUpdateTssAddress
// MsgUpdateGatewayGasLimit

// However, the transactions other than `AddToInboundTracker`, `UpdateGasPriceIncreaseFlags`, and `UpdateGatewayGasLimit` have already been used in other tests.
func TestCriticalAdminTransactions(r *runner.E2ERunner, _ []string) {
	TestAddToInboundTracker(r)
	TestUpdateGasPriceIncreaseFlags(r)
	TestUpdateGatewayGasLimit(r)
}

func TestUpdateGasPriceIncreaseFlags(r *runner.E2ERunner) {
	// Set default flags on zetacore
	defaultFlags := observertypes.DefaultGasPriceIncreaseFlags
	msgGasPriceFlags := observertypes.NewMsgUpdateGasPriceIncreaseFlags(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		defaultFlags,
	)
	_, err := r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, msgGasPriceFlags)
	require.NoError(r, err)

	// create a new set of flag values by incrementing the epoch length by 1
	defaultFlagsUpdated := defaultFlags
	defaultFlagsUpdated.EpochLength = defaultFlags.EpochLength + 1

	// Update the flags on zetacore with the new values
	msgGasPriceFlags = observertypes.NewMsgUpdateGasPriceIncreaseFlags(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		defaultFlagsUpdated,
	)
	_, err = r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, msgGasPriceFlags)
	require.NoError(r, err)

	r.WaitForBlocks(1)

	// Verify that the flags have been updated
	flags, err := r.ObserverClient.CrosschainFlags(r.Ctx, &observertypes.QueryGetCrosschainFlagsRequest{})
	require.NoError(r, err)
	require.Equal(r, defaultFlagsUpdated.EpochLength, flags.CrosschainFlags.GasPriceIncreaseFlags.EpochLength)
}

func TestAddToInboundTracker(r *runner.E2ERunner) {
	chainEth := chains.GoerliLocalnet
	chainBtc := chains.BitcoinRegtest
	msgEth := crosschaintypes.NewMsgAddInboundTracker(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.EmergencyPolicyName),
		chainEth.ChainId,
		coin.CoinType_Gas,
		sample.Hash().Hex(),
	)
	_, err := r.ZetaTxServer.BroadcastTx(utils.EmergencyPolicyName, msgEth)
	require.NoError(r, err)

	msgBtc := crosschaintypes.NewMsgAddInboundTracker(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.EmergencyPolicyName),
		chainBtc.ChainId,
		coin.CoinType_Gas,
		sample.BtcHash().String(),
	)

	_, err = r.ZetaTxServer.BroadcastTx(utils.EmergencyPolicyName, msgBtc)
	require.NoError(r, err)

	r.WaitForBlocks(1)

	tracker, err := r.CctxClient.InboundTracker(r.Ctx, &crosschaintypes.QueryInboundTrackerRequest{
		ChainId: msgEth.ChainId,
		TxHash:  msgEth.TxHash,
	})
	require.NoError(r, err)
	require.NotNil(r, tracker)
	require.Equal(r, msgEth.TxHash, tracker.InboundTracker.TxHash)

	tracker, err = r.CctxClient.InboundTracker(r.Ctx, &crosschaintypes.QueryInboundTrackerRequest{
		ChainId: msgBtc.ChainId,
		TxHash:  msgBtc.TxHash,
	})
	require.NoError(r, err)
	require.NotNil(r, tracker)
	require.Equal(r, msgBtc.TxHash, tracker.InboundTracker.TxHash)
}

func TestUpdateGatewayGasLimit(r *runner.E2ERunner) {
	r.UpdateGatewayGasLimit(uint64(1_600_000))
}
