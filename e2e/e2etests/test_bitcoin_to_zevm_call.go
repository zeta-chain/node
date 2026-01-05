package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/memo"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	zetabtc "github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
)

func TestBitcoinToZEVMCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	// ARRANGE
	// create a short payload less than max OP_RETURN data size (80 bytes)
	payload := randomPayloadWithSize(r, 50)
	sender := r.GetBtcAddress().EncodeAddress()
	r.AssertTestDAppZEVMCalled(false, payload, []byte(sender), big.NewInt(0))

	// wrap the payload in a standard memo
	memo := &memo.InboundMemo{
		Header: memo.Header{
			Version:     0,
			EncodingFmt: memo.EncodingFmtCompactShort,
			OpCode:      memo.OpCodeCall,
		},
		FieldsV0: memo.FieldsV0{
			Receiver: r.TestDAppV2ZEVMAddr,
			Payload:  []byte(payload),
		},
	}

	// ACT
	// make a NoAssetCall to ZEVM with tiny amount
	// the amount barely covers the depositor fee, the remaining amount is ignored
	amount := zetabtc.DefaultDepositorFee + 0.0000001
	txHash := r.DepositBTCWithAmount(amount, memo)

	// ASSERT
	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "call")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// check the payload was received on the contract
	r.AssertTestDAppZEVMCalled(true, payload, []byte(sender), big.NewInt(0))
}
