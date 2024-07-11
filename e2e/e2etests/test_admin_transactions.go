package e2etests

import (
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/testutil/sample"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestAdminTransactions(r *runner.E2ERunner, args []string) {
	r.Logger.Info("Adding a inbound tracker ")
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

	time.Sleep(8 * time.Second)

	tracker, err := r.CctxClient.InboundTracker(r.Ctx, &crosschaintypes.QueryInboundTrackerRequest{
		ChainID: msgEth.ChainId,
		TxHash:  msgEth.TxHash,
	})
	require.NoError(r, err)
	require.NotNil(r, tracker)
	require.Equal(r, msgEth.TxHash, tracker.InboundTracker.TxHash)

	tracker, err = r.CctxClient.InboundTracker(r.Ctx, &crosschaintypes.QueryInboundTrackerRequest{
		ChainID: msgBtc.ChainId,
		TxHash:  msgBtc.TxHash,
	})
	require.NoError(r, err)
	require.NotNil(r, tracker)
	require.Equal(r, msgBtc.TxHash, tracker.InboundTracker.TxHash)
}
