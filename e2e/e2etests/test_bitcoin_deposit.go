package e2etests

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestBitcoinDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	depositAmount := utils.ParseFloat(r, args[0])
	// ZRC20 BTC amounts have 8 decimals
	depositAmountZRC20 := uint64(depositAmount * 1e8)

	startingBalance, err := r.BTCZRC20.BalanceOf(&bind.CallOpts{}, r.ZEVMAuth.From)
	require.NoError(r, err)

	txHash := r.DepositBTCWithExactAmount(depositAmount, nil)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// assert that the inbound amount is expected
	require.InDelta(r, depositAmountZRC20, cctx.InboundParams.Amount.Uint64(), 100)

	// assert that the balance increases by the expected amount
	endingBalance, err := r.BTCZRC20.BalanceOf(&bind.CallOpts{}, r.ZEVMAuth.From)
	require.NoError(r, err)
	balanceDiff := bigSub(endingBalance, startingBalance)
	require.InDelta(r, depositAmountZRC20, balanceDiff.Uint64(), 100)
}
