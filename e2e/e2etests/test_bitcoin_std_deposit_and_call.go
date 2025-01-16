package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/memo"
	testcontract "github.com/zeta-chain/node/testutil/contracts"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	zetabitcoin "github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
)

func TestBitcoinStdMemoDepositAndCall(r *runner.E2ERunner, args []string) {
	// start mining blocks if local bitcoin
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	// parse amount to deposit
	require.Len(r, args, 1)
	amount := utils.ParseFloat(r, args[0])

	// deploy an example contract in ZEVM
	contractAddr, _, contract, err := testcontract.DeployExample(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)

	// create standard memo with [receiver, payload]
	memo := &memo.InboundMemo{
		Header: memo.Header{
			Version:     0,
			EncodingFmt: memo.EncodingFmtCompactShort,
			OpCode:      memo.OpCodeDepositAndCall,
		},
		FieldsV0: memo.FieldsV0{
			Receiver: contractAddr,
			Payload:  []byte("hello satoshi"),
		},
	}

	// deposit BTC with standard memo
	txHash := r.DepositBTCWithAmount(amount, memo, true)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "bitcoin_std_memo_deposit_and_call")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// check if example contract has been called, 'bar' value should be set to amount
	amountSats, err := zetabitcoin.GetSatoshis(amount)
	require.NoError(r, err)
	utils.MustHaveCalledExampleContract(r, contract, big.NewInt(amountSats))
}
