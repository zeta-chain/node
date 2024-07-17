package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/testutil/sample"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestAdminTransactions(r *runner.E2ERunner, _ []string) {
	TestAddToInboundTracker(r)
	TestUpdateGasPriceIncreaseFlags(r)
}

func TestUpdateGasPriceIncreaseFlags(r *runner.E2ERunner) {
	defaultFlags := observertypes.DefaultGasPriceIncreaseFlags
	msgGasPriceFlags := observertypes.NewMsgUpdateGasPriceIncreaseFlags(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		defaultFlags,
	)
	_, err := r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, msgGasPriceFlags)
	require.NoError(r, err)

	defaultFlagsUpdated := defaultFlags
	defaultFlagsUpdated.EpochLength = defaultFlags.EpochLength + 1

	msgGasPriceFlags = observertypes.NewMsgUpdateGasPriceIncreaseFlags(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		defaultFlagsUpdated,
	)

	_, err = r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, msgGasPriceFlags)
	require.NoError(r, err)

	WaitForBlocks(r, 1)

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

	WaitForBlocks(r, 1)

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
