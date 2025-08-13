package e2etests

import (
	"math/big"

	"github.com/btcsuite/btcd/txscript"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
)

func TestBitcoinDepositAndCall(r *runner.E2ERunner, args []string) {
	// Given amount to send
	require.Len(r, args, 1)
	amount := utils.ParseFloat(r, args[0])
	amountSats, err := common.GetSatoshis(amount)
	require.NoError(r, err)

	oldBalance, err := r.BTCZRC20.BalanceOf(&bind.CallOpts{}, r.TestDAppV2ZEVMAddr)
	require.NoError(r, err)

	// ARRANGE
	// create a random payload exactly fit max OP_RETURN data size 80 bytes (20 receiver + 60 payload)
	size := txscript.MaxDataCarrierSize - ethcommon.AddressLength
	payload := randomPayloadWithSize(r, size)
	r.AssertTestDAppZEVMCalled(false, payload, big.NewInt(amountSats))

	// ACT
	// Send BTC to TSS address with a dummy memo
	memo := append(r.TestDAppV2ZEVMAddr.Bytes(), payload...)
	txHash, err := r.SendToTSSWithMemo(amount, memo)
	require.NoError(r, err)

	// ASSERT
	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "bitcoin_deposit_and_call")
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
