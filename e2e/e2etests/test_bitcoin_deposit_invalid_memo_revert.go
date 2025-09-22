package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/memo"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestBitcoinDepositInvalidMemoRevert(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// mine blocks at normal speed
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	// Given amount
	amount := utils.ParseFloat(r, args[0])

	// CASE 1
	// make a deposit without memo output
	txHash, err := r.SendToTSSWithMemo(amount, nil)
	require.NoError(r, err)

	// wait for the cctx to be reverted
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit without memo")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)
	require.EqualValues(r, crosschaintypes.InboundStatus_INVALID_MEMO, cctx.InboundParams.Status)
	utils.RequireCCTXErrorMessages(r, cctx, "invalid memo: no memo found in inbound")

	// CASE 2
	// make a deposit with an empty memo
	txHash, err = r.SendToTSSWithMemo(amount, []byte("memo too short"))
	require.NoError(r, err)

	// wait for the cctx to be reverted
	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit empty memo")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)
	require.EqualValues(r, crosschaintypes.InboundStatus_INVALID_MEMO, cctx.InboundParams.Status)
	utils.RequireCCTXErrorMessages(r, cctx, "invalid memo: legacy memo length must be at least 20 bytes")

	// CASE 3
	// make a deposit with an invalid standard memo
	memoStd := &memo.InboundMemo{
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
	memoBytes, err := memoStd.EncodeToBytes()
	require.NoError(r, err)
	memoBytes[2] = 0x00

	// deposit to TSS
	txHash, err = r.SendToTSSWithMemo(amount, memoBytes)
	require.NoError(r, err)

	// wait for the cctx to be reverted
	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit invalid standard memo")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)
	require.EqualValues(r, crosschaintypes.InboundStatus_INVALID_MEMO, cctx.InboundParams.Status)
	utils.RequireCCTXErrorMessages(
		r,
		cctx,
		"invalid memo: standard memo contains improper data",
		"payload is not allowed for deposit operation",
	)

	// CASE 4
	// make a deposit with an invalid revert address
	memoStd = &memo.InboundMemo{
		Header: memo.Header{
			Version:     0,
			EncodingFmt: memo.EncodingFmtCompactShort,
			OpCode:      memo.OpCodeDepositAndCall,
		},
		FieldsV0: memo.FieldsV0{
			Receiver: r.TestDAppV2ZEVMAddr,
			Payload:  []byte("a payload"),
			RevertOptions: crosschaintypes.RevertOptions{
				// invalid revert address, not a BTC address
				RevertAddress: sample.EthAddress().Hex(),
			},
		},
	}
	memoBytes, err = memoStd.EncodeToBytes()
	require.NoError(r, err)

	// deposit to TSS
	txHash, err = r.SendToTSSWithMemo(amount, memoBytes)
	require.NoError(r, err)

	// wait for the cctx to be reverted
	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit invalid revert address")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)
	require.EqualValues(r, crosschaintypes.InboundStatus_INVALID_MEMO, cctx.InboundParams.Status)
	utils.RequireCCTXErrorMessages(
		r,
		cctx,
		"invalid memo: invalid standard memo for bitcoin",
		"invalid revert address in memo",
	)
}
