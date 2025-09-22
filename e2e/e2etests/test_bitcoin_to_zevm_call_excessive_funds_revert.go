package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/memo"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	zetabtc "github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
)

func TestBitcoinToZEVMCallExcessiveFundsRevert(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	// start mining blocks
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	// ARRANGE
	// wrap the payload in a standard memo
	memo := &memo.InboundMemo{
		Header: memo.Header{
			Version:     0,
			EncodingFmt: memo.EncodingFmtCompactShort,
			OpCode:      memo.OpCodeCall,
		},
		FieldsV0: memo.FieldsV0{
			Receiver: r.TestDAppV2ZEVMAddr,
			Payload:  []byte("a payload"),
		},
	}

	// Given a amount larger than maxNoAssetCallExcessAmount
	// the amount needs to be much higher than maxNoAssetCallExcessAmount (0.001 BTC) to be able to pay withdrawal fee
	amount := 0.003
	amountSats, err := zetabtc.GetSatoshis(amount)
	require.NoError(r, err)

	// ACT
	// make a NoAssetCall to ZEVM with a large amount
	txHash := r.DepositBTCWithExactAmount(amount, memo)

	// ASSERT
	// wait for the cctx to be reverted
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "call")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)

	// check cctx details
	// the revert amount should be less than amount because a withdrawal fee is charged
	require.EqualValues(r, amountSats, cctx.InboundParams.Amount.Uint64())
	require.Less(r, cctx.GetCurrentOutboundParam().Amount.BigInt().Int64(), amountSats)

	// check cctx status and error message
	require.EqualValues(r, crosschaintypes.InboundStatus_EXCESSIVE_NOASSETCALL_FUNDS, cctx.InboundParams.Status)
	require.Regexp(r, "remaining funds of [0-9]+ satoshis exceed 100000 satoshis", cctx.CctxStatus.ErrorMessage)
}
