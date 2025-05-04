package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/memo"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestBitcoinStdMemoDepositAndCallRevertOtherAddress(r *runner.E2ERunner, args []string) {
	// Start mining blocks
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	// Parse amount to send
	require.Len(r, args, 1)
	amount := utils.ParseFloat(r, args[0])

	// Create a memo to call non-existing contract
	revertAddress := "bcrt1qy9pqmk2pd9sv63g27jt8r657wy0d9uee4x2dt2"
	memo := &memo.InboundMemo{
		Header: memo.Header{
			Version:     0,
			EncodingFmt: memo.EncodingFmtCompactShort,
			OpCode:      memo.OpCodeDepositAndCall,
		},
		FieldsV0: memo.FieldsV0{
			Receiver: sample.EthAddress(), // non-existing contract
			Payload:  []byte("a payload"),
			RevertOptions: types.RevertOptions{
				RevertAddress: revertAddress,
			},
		},
	}

	// ACT
	// Deposit
	txHash := r.DepositBTCWithExactAmount(amount, memo)

	// ASSERT
	// Now we want to make sure revert TX is completed.
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "bitcoin_std_memo_deposit")
	utils.RequireCCTXStatus(r, cctx, types.CctxStatus_Reverted)

	// Make sure inbound sender and revert address are correct
	require.Equal(r, cctx.InboundParams.Sender, r.GetBtcAddress().EncodeAddress())
	require.Equal(r, cctx.GetCurrentOutboundParam().Receiver, revertAddress)

	// Check revert tx receiver address and amount
	receiver, value := r.QueryOutboundReceiverAndAmount(cctx.OutboundParams[1].Hash)
	require.Equal(r, revertAddress, receiver)
	require.True(r, value > 0)

	r.Logger.Info(
		"Sent %f BTC to TSS to call non-existing contract, got refund of %d satoshis to other address",
		amount,
		value,
	)
}
