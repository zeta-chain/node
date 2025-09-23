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
)

func TestBitcoinStdMemoDepositAndCallRevertAndAbort(r *runner.E2ERunner, args []string) {
	// Start mining blocks
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	require.Len(r, args, 0)
	amount := 0.00000001 // 1 satoshi so revert fails because of insufficient gas

	// deploy testabort contract
	testAbortAddr, _, testAbort, err := testabort.DeployTestAbort(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)

	// Create a memo to call non-existing contract
	abortMessage := "abort message"
	inboundMemo := &memo.InboundMemo{
		Header: memo.Header{
			Version:     0,
			EncodingFmt: memo.EncodingFmtCompactShort,
			OpCode:      memo.OpCodeDepositAndCall,
		},
		FieldsV0: memo.FieldsV0{
			Receiver: sample.EthAddress(), // non-existing contract
			Payload:  []byte("a payload"),
			RevertOptions: types.RevertOptions{
				CallOnRevert:  false,
				AbortAddress:  testAbortAddr.Hex(),
				RevertMessage: []byte(abortMessage),
			},
		},
	}

	// ACT
	// Deposit
	txHash := r.DepositBTCWithExactAmount(amount, inboundMemo)

	// ASSERT
	// Now we want to make sure revert TX is completed.
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "bitcoin_std_memo_deposit")
	utils.RequireCCTXStatus(r, cctx, types.CctxStatus_Aborted)

	// wait for the abort contract to receive tokens
	change := utils.NewBalanceChange(true)
	utils.WaitAndVerifyZRC20BalanceChange(r, r.BTCZRC20, testAbortAddr, big.NewInt(0), change, r.Logger)

	// check onAbort was called
	aborted, err := testAbort.IsAborted(&bind.CallOpts{})
	require.NoError(r, err)
	require.True(r, aborted)

	// check revert message was used
	abortContext, err := testAbort.GetAbortedWithMessage(&bind.CallOpts{}, abortMessage)
	require.NoError(r, err)
	require.EqualValues(r, []byte(abortMessage), abortContext.RevertMessage)
}
