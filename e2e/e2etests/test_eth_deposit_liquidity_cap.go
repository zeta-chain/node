package e2etests

import (
	"math/big"

	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/x/crosschain/types"
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
	res, err := r.ZetaTxServer.SetZRC20LiquidityCap(r.ETHZRC20Addr, liquidityCap)
	require.NoError(r, err)

	r.Logger.Info("set liquidity cap tx hash: %s", res.TxHash)
	r.Logger.Info("Depositing more than liquidity cap should make cctx reverted")

	signedTx := r.ETHDeposit(
		r.EVMAddress(),
		amountMoreThanCap,
		gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
		true,
	)

	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, signedTx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, types.CctxStatus_Reverted)

	r.Logger.Info("CCTX has been reverted")

	r.Logger.Info("Depositing less than liquidity cap should still succeed")
	initialBal, err := r.ETHZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)

	signedTx = r.ETHDeposit(
		r.EVMAddress(),
		amountLessThanCap,
		gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
		true,
	)

	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, signedTx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, types.CctxStatus_OutboundMined)

	expectedBalance := big.NewInt(0).Add(initialBal, amountLessThanCap)

	bal, err := r.ETHZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	require.Equal(r, 0, bal.Cmp(expectedBalance))

	r.Logger.Info("Deposit succeeded")

	r.Logger.Info("Removing the liquidity cap")
	res, err = r.ZetaTxServer.RemoveZRC20LiquidityCap(r.ETHZRC20Addr)
	require.NoError(r, err)
	r.Logger.Info("remove liquidity cap tx hash: %s", res.TxHash)

	initialBal, err = r.ETHZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)

	signedTx = r.ETHDeposit(
		r.EVMAddress(),
		amountMoreThanCap,
		gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
		true,
	)

	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, signedTx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, types.CctxStatus_OutboundMined)

	// wait for the zrc20 balance to be updated
	change := utils.NewExactChange(amountMoreThanCap)
	utils.WaitAndVerifyZRC20BalanceChange(r, r.ETHZRC20, r.EVMAddress(), initialBal, change, r.Logger)

	r.Logger.Info("New deposit succeeded")
}
