package e2etests

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/memo"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestBitcoinStdMemoInscribedDepositAndCall(r *runner.E2ERunner, args []string) {
	// Given amount and live network fee rate
	require.Len(r, args, 1)
	amount := utils.ParseFloat(r, args[0])
	feeRate := r.BitcoinEstimateFeeRate(1)

	oldBalance, err := r.BTCZRC20.BalanceOf(&bind.CallOpts{}, r.TestDAppV2ZEVMAddr)
	require.NoError(r, err)

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
	rawTx, depositedAmount, commitAddress := r.InscribeToTSSWithMemo(amount, memoBytes, int64(feeRate))

	// ASSERT
	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, rawTx.Txid, r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "bitcoin_std_memo_inscribed_deposit_and_call")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// sender, txOrigin should be set to the initiator of the inscription, not the commit address
	senderAddress := r.GetBtcAddress().EncodeAddress()
	require.False(r, strings.EqualFold(senderAddress, commitAddress))
	require.Equal(r, senderAddress, cctx.InboundParams.Sender)
	require.Equal(r, senderAddress, cctx.InboundParams.TxOrigin)

	// calculate received amount
	receivedAmount := r.BitcoinCalcReceivedAmount(rawTx, depositedAmount)

	// wait for the zrc20 balance to be updated
	change := utils.NewExactChange(big.NewInt(receivedAmount))
	utils.WaitAndVerifyZRC20BalanceChange(r, r.BTCZRC20, r.TestDAppV2ZEVMAddr, oldBalance, change, r.Logger)

	// check the payload was received on the contract
	r.AssertTestDAppZEVMCalled(true, payload, big.NewInt(receivedAmount))
}
