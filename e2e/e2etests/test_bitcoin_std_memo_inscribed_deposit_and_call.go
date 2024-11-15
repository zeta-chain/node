package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/memo"
	testcontract "github.com/zeta-chain/node/testutil/contracts"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	zetabitcoin "github.com/zeta-chain/node/zetaclient/chains/bitcoin"
)

func TestBitcoinStdMemoInscribedDepositAndCall(r *runner.E2ERunner, args []string) {
	// Start mining blocks
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	// Given amount to send and fee rate
	require.Len(r, args, 2)
	amount := parseFloat(r, args[0])
	feeRate := parseInt(r, args[1])

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
	txHash, depositAmount := r.InscribeToTSSFromDeployerWithMemo(amount, memoBytes, int64(feeRate))

	// ASSERT
	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "bitcoin_std_memo_inscribed_deposit_and_call")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// check if example contract has been called, 'bar' value should be set to correct amount
	depositFeeSats, err := zetabitcoin.GetSatoshis(zetabitcoin.DefaultDepositorFee)
	require.NoError(r, err)
	receiveAmount := depositAmount - depositFeeSats
	utils.MustHaveCalledExampleContract(r, contract, big.NewInt(receiveAmount))
}
