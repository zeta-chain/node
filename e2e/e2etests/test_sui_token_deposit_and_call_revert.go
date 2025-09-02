package e2etests

import (
	"math/big"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/coin"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestSuiTokenDepositAndCallRevert(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")
	balanceBefore := r.SuiGetFungibleTokenBalance(signer.Address())
	tssBalanceBefore := r.SuiGetSUIBalance(r.SuiTSSAddress)

	// add liquidity in pool to allow revert fee to be paid
	zetaAmount := big.NewInt(1e18)
	zrc20Amount := big.NewInt(9000000000)
	r.AddLiquiditySUI(zetaAmount, zrc20Amount)
	r.AddLiquiditySuiFungibleToken(zetaAmount, zrc20Amount)

	// make the deposit transaction
	resp := r.SuiFungibleTokenDepositAndCall(r.TestDAppV2ZEVMAddr, math.NewUintFromBigInt(amount), []byte("revert"))

	r.Logger.Info("Sui deposit and call tx: %s", resp.Digest)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, resp.Digest, r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)
	require.EqualValues(r, coin.CoinType_ERC20, cctx.InboundParams.CoinType)
	require.EqualValues(r, amount.Uint64(), cctx.InboundParams.Amount.Uint64())

	// check the balance after the failed deposit is higher than balance before - amount
	// reason it's not equal is because of the gas fee for revert
	balanceAfter := r.SuiGetFungibleTokenBalance(signer.Address())
	require.Greater(r, balanceAfter, balanceBefore-amount.Uint64())

	// check the TSS balance after transaction is higher or equal to the balance before
	// reason is that the max budget is refunded to the TSS
	tssBalanceAfter := r.SuiGetSUIBalance(r.SuiTSSAddress)
	require.GreaterOrEqual(r, tssBalanceAfter, tssBalanceBefore)
}
