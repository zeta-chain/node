package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/memo"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestBitcoinStdMemoInscribedDepositAndCall(r *runner.E2ERunner, args []string) {
	// mine blocks at normal speed
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	// Given amount and live network fee rate
	require.Len(r, args, 1)
	amount := utils.ParseFloat(r, args[0])
	feeRate := r.BitcoinEstimateFeeRate(1)

	// ARRANGE
	// create a random payload with more than MaxDataCarrierSize (80 bytes)
	// memo size: 4b header + 20b receiver + 1b length + 100b payload == 125 bytes
	payload := randomPayloadWithSize(r, 100)

	// wrap the payload in a standard memo
	memo := &memo.InboundMemo{
		Header: memo.Header{
			Version:     0,
			EncodingFmt: memo.EncodingFmtCompactShort,
			OpCode:      memo.OpCodeDepositAndCall,
		},
		FieldsV0: memo.FieldsV0{
			Receiver: r.TestDAppV2ZEVMAddr,
			Payload:  []byte(payload),
		},
	}
	memoBytes, err := memo.EncodeToBytes()
	require.NoError(r, err)

	// ACT
	// Send BTC to TSS address with memo
	// #nosec G115 test - checked in range
	rawTx, _, _ := r.InscribeToTSSWithMemo(amount, memoBytes, int64(feeRate))

	// ASSERT
	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, rawTx.Txid, r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "bitcoin_std_memo_inscribed_deposit_and_call")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)
}
