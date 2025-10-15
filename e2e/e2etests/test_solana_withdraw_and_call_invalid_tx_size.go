package e2etests

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	solanacontract "github.com/zeta-chain/node/pkg/contracts/solana"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestSolanaWithdrawAndCallInvalidTxSize executes withdrawAndCall, but with invalid tx size
// in that case, cctx is reverted due to "transaction is too large" error
func TestSolanaWithdrawAndCallInvalidTxSize(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// ARRANGE
	// Given withdraw amount
	withdrawAmount := utils.ParseBigInt(r, args[0])

	// ensure runner has enough balance
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

	// get connected program pda
	connectedPda, err := solanacontract.ComputeConnectedPdaAddress(r.ConnectedProgram)
	require.NoError(r, err)

	// generate a big enough data that exceeds Solana max tx size (encoded/raw 1644/1232)
	data := make([]byte, 2048)
	_, err = rand.Read(data)
	require.NoError(r, err)

	// encode msg
	msg := solanacontract.ExecuteMsg{
		Accounts: []solanacontract.AccountMeta{
			{PublicKey: [32]byte(connectedPda.Bytes()), IsWritable: true},
			{PublicKey: [32]byte(r.ComputePdaAddress().Bytes()), IsWritable: false},
			{PublicKey: [32]byte(solana.SystemProgramID.Bytes()), IsWritable: false},
			{PublicKey: [32]byte(solana.SysVarInstructionsPubkey.Bytes()), IsWritable: false},
			{PublicKey: [32]byte(r.GetSolanaPrivKey().PublicKey().Bytes()), IsWritable: true},
		},
		Data: data,
	}

	msgEncoded, err := msg.Encode()
	require.NoError(r, err)

	// ACT
	// withdraw and call
	tx := r.WithdrawAndCallSOLZRC20(
		withdrawAmount,
		approvedAmount,
		msgEncoded,
		gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)

	// refresh ERC20 SOL balance immediately after withdraw
	// the refund is now on the way, we need to remember the balance and wait for update
	balanceBefore, err = r.SOLZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SOL right after withdraw: %d", balanceBefore)

	// ASSERT
	// wait for the cctx to be reverted
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)

	// get ERC20 SOL balance after withdraw
	balanceAfter, err := r.SOLZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SOL after CCTX is reverted: %d", balanceAfter)

	// check if the balance is increased correctly
	balanceExpected := new(big.Int).Add(balanceBefore, withdrawAmount)
	require.True(r, balanceExpected.Cmp(balanceAfter) == 0, "balance is not refunded correctly")

	// check that failure log is attached to increment nonce instruction
	txIncNonce, err := r.SolanaClient.GetTransaction(
		r.Ctx,
		solana.MustSignatureFromBase58(cctx.OutboundParams[0].Hash),
		nil,
	)
	require.NoError(r, err)

	expectedLog := "Program log: Failure reason: transaction is too large"
	for _, log := range txIncNonce.Meta.LogMessages {
		if strings.Contains(log, expectedLog) {
			return
		}
	}
	require.NoError(r, errors.New("expected log not found"))
}
