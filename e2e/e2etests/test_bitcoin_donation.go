package e2etests

import (
	"time"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/constant"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	zetabitcoin "github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
)

func TestBitcoinDonation(r *runner.E2ERunner, args []string) {
	// Given amount to send
	require.Len(r, args, 1)
	amount := utils.ParseFloat(r, args[0])
	amountTotal := amount + zetabitcoin.DefaultDepositorFee

	// ACT
	// Send BTC to TSS address with donation message
	memo := []byte(constant.DonationMessage)
	txHash, err := r.SendToTSSWithMemo(amountTotal, memo)
	require.NoError(r, err)

	// ASSERT after 6 Zeta blocks
	time.Sleep(constant.ZetaBlockTime * 6)
	req := &crosschaintypes.QueryInboundHashToCctxDataRequest{InboundHash: txHash.String()}
	_, err = r.CctxClient.InboundHashToCctxData(r.Ctx, req)
	require.Error(r, err)
}
