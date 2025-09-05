package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/memo"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestBitcoinDepositInvalidMemoRevert(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	// mine blocks at normal speed
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	// CASE 1
	// make a deposit without memo output
	txHash, err := r.SendToTSSWithMemo(0.1, nil)
	require.NoError(r, err)

	// wait for the cctx to be reverted
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit without memo")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)
	require.EqualValues(r, crosschaintypes.InboundStatus_INVALID_MEMO, cctx.InboundParams.Status)

	// CASE 2
	// make a deposit with a empty memo
	txHash, err = r.SendToTSSWithMemo(0.1, []byte{})
	require.NoError(r, err)

	// wait for the cctx to be reverted
	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit empty memo")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)
	require.EqualValues(r, crosschaintypes.InboundStatus_INVALID_MEMO, cctx.InboundParams.Status)

	// CASE 3
	// make a deposit with an invalid legacy memo
	txHash, err = r.SendToTSSWithMemo(0.1, []byte("invalid legacy memo"))
	require.NoError(r, err)

	// wait for the cctx to be reverted
	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit invalid memo")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)
	require.EqualValues(r, crosschaintypes.InboundStatus_INVALID_MEMO, cctx.InboundParams.Status)

	// CASE 4
	// make a deposit with an invalid standard memo
	memo := &memo.InboundMemo{
		Header: memo.Header{
			Version:     0,
			EncodingFmt: memo.EncodingFmtCompactShort,
			OpCode:      memo.OpCodeDepositAndCall,
		},
		FieldsV0: memo.FieldsV0{
			Receiver: r.TestDAppV2ZEVMAddr,
			Payload:  []byte("payload is not allowed"),
		},
	}

	// modify the op code to 0b0000 (deposit), so payload won't be allowed
	memoBytes, err := memo.EncodeToBytes()
	require.NoError(r, err)
	memoBytes[2] = 0x00

	// deposit to TSS
	txHash, err = r.SendToTSSWithMemo(0.1, memoBytes)
	require.NoError(r, err)

	// wait for the cctx to be reverted
	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit invalid standard memo")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)
	require.EqualValues(r, crosschaintypes.InboundStatus_INVALID_MEMO, cctx.InboundParams.Status)
}
