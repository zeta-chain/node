package e2etests

import (
	"math/big"

	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

// TestDepositEtherLiquidityCap tests depositing Ethers in a context where a liquidity cap is set
func TestDepositEtherLiquidityCap(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	liquidityCapArg := math.NewUintFromString(args[0])
	supply, err := r.ETHZRC20.TotalSupply(&bind.CallOpts{})
	require.NoError(r, err)

	liquidityCap := math.NewUintFromBigInt(supply).Add(liquidityCapArg)
	amountLessThanCap := liquidityCapArg.BigInt().Div(liquidityCapArg.BigInt(), big.NewInt(10)) // 1/10 of the cap
	amountMoreThanCap := liquidityCapArg.BigInt().Mul(liquidityCapArg.BigInt(), big.NewInt(10)) // 10 times the cap
	msg := fungibletypes.NewMsgUpdateZRC20LiquidityCap(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		r.ETHZRC20Addr.Hex(),
		liquidityCap,
	)
	res, err := r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, msg)
	require.NoError(r, err)

	r.Logger.Info("set liquidity cap tx hash: %s", res.TxHash)
	r.Logger.Info("Depositing more than liquidity cap should make cctx reverted")

	signedTx, err := r.SendEther(r.TSSAddress, amountMoreThanCap, nil)
	require.NoError(r, err)

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, signedTx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, signedTx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, types.CctxStatus_Reverted)

	r.Logger.Info("CCTX has been reverted")

	r.Logger.Info("Depositing less than liquidity cap should still succeed")
	initialBal, err := r.ETHZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)

	signedTx, err = r.SendEther(r.TSSAddress, amountLessThanCap, nil)
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, signedTx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, signedTx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, types.CctxStatus_OutboundMined)

	expectedBalance := big.NewInt(0).Add(initialBal, amountLessThanCap)

	bal, err := r.ETHZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	require.Equal(r, 0, bal.Cmp(expectedBalance))

	r.Logger.Info("Deposit succeeded")

	r.Logger.Info("Removing the liquidity cap")
	msg = fungibletypes.NewMsgUpdateZRC20LiquidityCap(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		r.ETHZRC20Addr.Hex(),
		math.ZeroUint(),
	)

	res, err = r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, msg)
	require.NoError(r, err)

	r.Logger.Info("remove liquidity cap tx hash: %s", res.TxHash)

	initialBal, err = r.ETHZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)

	signedTx, err = r.SendEther(r.TSSAddress, amountMoreThanCap, nil)
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, signedTx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	utils.WaitCctxMinedByInboundHash(r.Ctx, signedTx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	expectedBalance = big.NewInt(0).Add(initialBal, amountMoreThanCap)

	bal, err = r.ETHZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	require.Equal(r,
		0,
		bal.Cmp(expectedBalance),
		"expected balance to be %s; got %s",
		expectedBalance.String(),
		bal.String(),
	)

	r.Logger.Info("New deposit succeeded")
}
