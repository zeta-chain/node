package e2etests

import (
	"math/big"

	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschainkeeper "github.com/zeta-chain/node/x/crosschain/keeper"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// TestEtherWithdrawRestricted tests the withdrawal to a restricted receiver address
func TestEtherWithdrawRestricted(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// ARRANGE
	// Given amount, receiver, revert address
	receiver := ethcommon.HexToAddress(args[0])
	amount := utils.ParseBigInt(r, args[1])
	revertAddress := r.EVMAddress()

	// receiver balance before
	receiverBalanceBefore, err := r.EVMClient.BalanceAt(r.Ctx, receiver, nil)
	require.NoError(r, err)

	// approve the ZRC20
	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	// ACT
	// perform the withdraw on restricted address
	tx := r.ETHWithdraw(
		receiver,
		amount,
		gatewayzevm.RevertOptions{
			RevertAddress:    revertAddress,
			OnRevertGasLimit: big.NewInt(0),
		},
	)

	r.Logger.EVMTransaction(tx, "withdraw to restricted address")

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	r.Logger.EVMReceipt(*receipt, "withdraw")
	r.Logger.ZRC20Withdrawal(r.ETHZRC20, *receipt, "withdraw")

	// revert address balance before
	revertBalanceBefore, err := r.ETHZRC20.BalanceOf(&bind.CallOpts{}, revertAddress)
	require.NoError(r, err)

	// ASSERT
	// wait for the cctx to be reverted
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, receipt.TxHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)

	// the cctx should be cancelled with zero value
	// note: the first outbound param is the cancel transaction
	r.EVMVerifyOutboundTransferAmount(cctx.OutboundParams[0].Hash, 0)

	// receiver balance should not change
	receiverBalanceAfter, err := r.EVMClient.BalanceAt(r.Ctx, receiver, nil)
	require.NoError(r, err)
	require.EqualValues(r, receiverBalanceBefore, receiverBalanceAfter)

	// revert address should receive the amount
	revertBalanceAfter, err := r.ETHZRC20.BalanceOf(&bind.CallOpts{}, revertAddress)
	require.NoError(r, err)

	userBalanceAfterUint := math.NewUintFromBigInt(revertBalanceAfter)
	userBalanceBeforeUint := math.NewUintFromBigInt(revertBalanceBefore)
	totalRevertAmount := getTotalRevertedAmount(r, cctx)

	require.EqualValues(r, userBalanceAfterUint.Sub(totalRevertAmount), userBalanceBeforeUint)
}

func getTotalRevertedAmount(r *runner.E2ERunner, cctx *crosschaintypes.CrossChainTx) math.Uint {
	outboundParams := cctx.OutboundParams[0]
	outboundTxGasUsed := math.NewUint(outboundParams.GasUsed)
	outboundTxFinalGasPrice := math.NewUintFromBigInt(outboundParams.EffectiveGasPrice.BigInt())
	outboundTxFeePaid := outboundTxGasUsed.Mul(outboundTxFinalGasPrice)
	userGasFeePaid := outboundParams.UserGasFeePaid
	totalRemainingFees := userGasFeePaid.Sub(outboundTxFeePaid)

	remainingFees := crosschainkeeper.PercentOf(totalRemainingFees, crosschaintypes.UsableRemainingFeesPercentage)
	if remainingFees.LTE(math.ZeroUint()) {
		return math.ZeroUint()
	}

	evmChainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	chainParams, err := r.ObserverClient.GetChainParamsForChain(
		r.Ctx,
		&observertypes.QueryGetChainParamsForChainRequest{
			ChainId: evmChainID.Int64(),
		},
	)
	require.NoError(r, err)

	stabilityPoolAmount := crosschainkeeper.PercentOf(remainingFees, chainParams.ChainParams.StabilityPoolPercentage)
	refundAmount := remainingFees.Sub(stabilityPoolAmount)
	return cctx.InboundParams.Amount.Add(refundAmount)
}
