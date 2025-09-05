package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/memo"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
)

func TestBitcoinStdMemoDepositAndCall(r *runner.E2ERunner, args []string) {
	// Given amount to deposit
	require.Len(r, args, 1)
	amount := utils.ParseFloat(r, args[0])
	amountSats, err := common.GetSatoshis(amount)
	require.NoError(r, err)

	oldBalance, err := r.BTCZRC20.BalanceOf(&bind.CallOpts{}, r.TestDAppV2ZEVMAddr)
	require.NoError(r, err)

	// ARRANGE
	// create a random payload exactly fit max OP_RETURN data size 80 bytes
	// memo size: 4b header + 20b receiver + 1b length + 55b payload == 80 bytes
	payload := randomPayloadWithSize(r, 55)
	r.AssertTestDAppZEVMCalled(false, payload, big.NewInt(amountSats))

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

	// ACT
	// deposit BTC with standard memo
	txHash := r.DepositBTCWithAmount(amount, memo)

	// ASSERT
	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "bitcoin_std_memo_deposit_and_call")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// calculate received amount
	rawTx, err := r.BtcRPCClient.GetRawTransactionVerbose(r.Ctx, txHash)
	require.NoError(r, err)
	receivedAmount := r.BitcoinCalcReceivedAmount(rawTx, amountSats)

	// wait for the zrc20 balance to be updated
	change := utils.NewExactChange(big.NewInt(receivedAmount))
	utils.WaitAndVerifyZRC20BalanceChange(r, r.BTCZRC20, r.TestDAppV2ZEVMAddr, oldBalance, change, r.Logger)

	// check the payload was received on the contract
	r.AssertTestDAppZEVMCalled(true, payload, big.NewInt(receivedAmount))
}
