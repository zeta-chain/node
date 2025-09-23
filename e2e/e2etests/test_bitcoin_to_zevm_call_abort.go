package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/contracts/testabort"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/memo"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
	zetabtc "github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
)

func TestBitcoinToZEVMCallAbort(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	// deploy testabort contract
	testAbortAddr, _, testAbort, err := testabort.DeployTestAbort(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)

	// ARRANGE
	// create a short payload less than max OP_RETURN data size (80 bytes)
	payload := randomPayloadWithSize(r, 10)
	sender := r.GetBtcAddress().EncodeAddress()
	r.AssertTestDAppZEVMCalled(false, payload, []byte(sender), big.NewInt(0))

	// wrap the payload in a standard memo
	abortMessage := "message abort"
	memo := &memo.InboundMemo{
		Header: memo.Header{
			Version:     0,
			EncodingFmt: memo.EncodingFmtCompactShort,
			OpCode:      memo.OpCodeCall, // NoAssetCall
		},
		FieldsV0: memo.FieldsV0{
			Receiver: sample.EthAddress(), // non-existing contract
			Payload:  []byte(payload),
			RevertOptions: types.RevertOptions{
				AbortAddress:  testAbortAddr.Hex(),
				RevertMessage: []byte(abortMessage),
			},
		},
	}

	// ACT
	// make a NoAssetCall to ZEVM with standard memo
	// the amount matches the exact depositor fee and should be accepted by observers
	txHash := r.DepositBTCWithAmount(zetabtc.DefaultDepositorFee, memo)

	// ASSERT
	// wait for the cctx to be aborted
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "call")
	utils.RequireCCTXStatus(r, cctx, types.CctxStatus_Aborted)

	// check onAbort was called
	aborted, err := testAbort.IsAborted(&bind.CallOpts{})
	require.NoError(r, err)
	require.True(r, aborted)

	// check revert message was used
	abortContext, err := testAbort.GetAbortedWithMessage(&bind.CallOpts{}, abortMessage)
	require.NoError(r, err)
	require.EqualValues(r, []byte(abortMessage), abortContext.RevertMessage)
}
