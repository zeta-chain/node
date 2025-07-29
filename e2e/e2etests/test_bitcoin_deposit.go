package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
)

func TestBitcoinDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// Given amount to send
	amount := utils.ParseFloat(r, args[0])
	amountSats, err := common.GetSatoshis(amount)
	require.NoError(r, err)

	oldBalance, err := r.BTCZRC20.BalanceOf(&bind.CallOpts{}, r.ZEVMAuth.From)
	require.NoError(r, err)

	// ACT
	txHash := r.DepositBTCWithAmount(amount, nil)

	// ASSERT
	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// calculate received amount
	rawTx, err := r.BtcRPCClient.GetRawTransactionVerbose(r.Ctx, txHash)
	require.NoError(r, err)
	receivedAmount := r.BitcoinCalcReceivedAmount(rawTx, amountSats)

	// wait for the zrc20 balance to be updated
	change := utils.NewExactChange(big.NewInt(receivedAmount))
	utils.WaitAndVerifyZRC20BalanceChange(r, r.BTCZRC20, r.ZEVMAuth.From, oldBalance, change, r.Logger)
}
