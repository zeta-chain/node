package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestSolanaWithdrawRevertExecutableReceiver reverts if receiver is solana program
func TestSolanaWithdrawRevertExecutableReceiver(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	withdrawAmount := utils.ParseBigInt(r, args[0])

	// get ERC20 SOL balance before withdraw
	balanceBefore, err := r.SOLZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SOL before withdraw: %d", balanceBefore)

	require.Equal(r, 1, balanceBefore.Cmp(withdrawAmount), "Insufficient balance for withdrawal")

	// parse withdraw amount (in lamports), approve amount is 1 SOL
	approvedAmount := new(big.Int).SetUint64(solana.LAMPORTS_PER_SOL)
	require.Equal(
		r,
		-1,
		withdrawAmount.Cmp(approvedAmount),
		"Withdrawal amount must be less than the approved amount: %v",
		approvedAmount,
	)

	// withdraw using example connected program as receiver
	tx := r.WithdrawSOLZRC20(
		r.ConnectedProgram,
		withdrawAmount,
		approvedAmount,
		gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)

	// wait for the cctx to be reverted
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)
}
