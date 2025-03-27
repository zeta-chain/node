package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	testcontract "github.com/zeta-chain/node/e2e/contracts/example"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/memo"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
)

func TestBitcoinStdMemoInscribedDepositAndCall(r *runner.E2ERunner, args []string) {
	// Given amount to send and fee rate
	require.Len(r, args, 2)
	amount := utils.ParseFloat(r, args[0])
	feeRate := utils.ParseInt(r, args[1])

	// deploy an example contract in ZEVM
	contractAddr, _, contract, err := testcontract.DeployExample(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)

	// create a standard memo > 80 bytes
	memo := &memo.InboundMemo{
		Header: memo.Header{
			Version:     0,
			EncodingFmt: memo.EncodingFmtCompactShort,
			OpCode:      memo.OpCodeDepositAndCall,
		},
		FieldsV0: memo.FieldsV0{
			Receiver: contractAddr,
			Payload:  []byte("for use case that passes a large memo > 80 bytes, inscripting the memo is the way to go"),
		},
	}
	memoBytes, err := memo.EncodeToBytes()
	require.NoError(r, err)

	// ACT
	// Send BTC to TSS address with memo
	txHash, depositAmount, commitAddress := r.InscribeToTSSWithMemo(amount, memoBytes, int64(feeRate))

	// ASSERT
	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "bitcoin_std_memo_inscribed_deposit_and_call")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// check if example contract has been called, 'bar' value should be set to correct amount
	depositFeeSats, err := common.GetSatoshis(common.DefaultDepositorFee)
	require.NoError(r, err)
	receiveAmount := depositAmount - depositFeeSats
	utils.MustHaveCalledExampleContract(
		r,
		contract,
		big.NewInt(receiveAmount),
		[]byte(commitAddress),
	)
}
