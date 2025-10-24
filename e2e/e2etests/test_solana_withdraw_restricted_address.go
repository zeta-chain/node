package e2etests

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestSolanaWithdrawRestricted(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// ARRANGE
	// Given amount, receiver, revert address
	receiverRestricted, err := chains.DecodeSolanaWalletAddress(args[0])
	require.NoError(r, err, fmt.Sprintf("unable to decode solana wallet address: %s", args[0]))

	// parse withdraw amount (in lamports), approve amount is 1 SOL
	approvedAmount := new(big.Int).SetUint64(solana.LAMPORTS_PER_SOL)
	withdrawAmount := utils.ParseBigInt(r, args[1])
	require.Equal(
		r,
		-1,
		withdrawAmount.Cmp(approvedAmount),
		"Withdrawal amount must be less than the approved amount (1e9).",
	)
	revertAddress := r.EVMAddress()

	// receiver balance before
	result, err := r.SolanaClient.GetBalance(r.Ctx, receiverRestricted, rpc.CommitmentFinalized)
	require.NoError(r, err)
	receiverBalanceBefore := result.Value

	// ACT
	// withdraw
	tx := r.WithdrawSOLZRC20(
		receiverRestricted,
		withdrawAmount,
		approvedAmount,
		gatewayzevm.RevertOptions{
			RevertAddress:    revertAddress,
			OnRevertGasLimit: big.NewInt(0),
		},
	)

	// wait for the withdraw tx to be mined
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// revert address balance before
	revertBalanceBefore, err := r.SOLZRC20.BalanceOf(&bind.CallOpts{}, revertAddress)
	require.NoError(r, err)

	// ASSERT
	// wait for the cctx to be reverted
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)

	// the outbound should be cancelled with zero value
	// note: the first outbound param is the cancel transaction
	r.SolanaVerifyWithdrawalAmount(cctx.OutboundParams[0].Hash, 0)

	// receiver balance should not change
	result, err = r.SolanaClient.GetBalance(r.Ctx, receiverRestricted, rpc.CommitmentFinalized)
	require.NoError(r, err)
	receiverBalanceAfter := result.Value
	require.EqualValues(r, receiverBalanceBefore, receiverBalanceAfter)

	// revert address should receive the amount
	revertBalanceAfter, err := r.SOLZRC20.BalanceOf(&bind.CallOpts{}, revertAddress)
	require.NoError(r, err)
	require.EqualValues(r, new(big.Int).Add(revertBalanceBefore, withdrawAmount), revertBalanceAfter)
}
