package e2etests

import (
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	solanacontract "github.com/zeta-chain/node/pkg/contracts/solana"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestSolanaWithdrawAndCallRevertWithCall executes withdrawAndCall on zevm and calls connected program on solana
// execution is reverted in connected program on_call function and onRevert is called on ZEVM TestDapp contract
func TestSolanaWithdrawAndCallRevertWithCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	withdrawAmount := utils.ParseBigInt(r, args[0])

	// get ZRC20 SOL balance before withdraw
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

	// use a random address to get the revert amount
	revertAddress := r.TestDAppV2ZEVMAddr
	balance, err := r.SOLZRC20.BalanceOf(&bind.CallOpts{}, revertAddress)
	require.NoError(r, err)
	initialBalance := balance

	payload := randomPayload(r)
	r.AssertTestDAppEVMCalled(false, payload, withdrawAmount)

	connectedPda, err := solanacontract.ComputeConnectedPdaAddress(r.ConnectedProgram)
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
		Data: []byte("revert"),
	}

	msgEncoded, err := msg.Encode()
	require.NoError(r, err)

	// withdraw and call
	tx := r.WithdrawAndCallSOLZRC20(
		withdrawAmount,
		approvedAmount,
		msgEncoded,
		gatewayzevm.RevertOptions{
			CallOnRevert:     true,
			RevertAddress:    revertAddress,
			RevertMessage:    []byte(payload),
			OnRevertGasLimit: big.NewInt(0),
		},
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)

	// get ZRC20 SOL balance after withdraw
	balanceAfter, err := r.SOLZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SOL after withdraw: %d", balanceAfter)

	r.AssertTestDAppZEVMCalled(true, payload, nil, big.NewInt(0))

	// check expected sender was used
	senderForMsg, err := r.TestDAppV2ZEVM.SenderWithMessage(
		&bind.CallOpts{},
		[]byte(payload),
	)
	require.NoError(r, err)
	require.Equal(r, r.ZEVMAuth.From, senderForMsg)

	// check the balance of revert address is equal to withdraw amount
	finalBalance, err := r.SOLZRC20.BalanceOf(&bind.CallOpts{}, revertAddress)
	require.NoError(r, err)

	require.Equal(r, withdrawAmount.Int64(), finalBalance.Int64()-initialBalance.Int64())

	// check that failure log is attached to increment nonce instruction
	txIncNonce, err := r.SolanaClient.GetTransaction(
		r.Ctx,
		solana.MustSignatureFromBase58(cctx.OutboundParams[0].Hash),
		nil,
	)
	require.NoError(r, err)

	expectedLog := "Program log: Failure reason: Program 4xEw862A2SEwMjofPkUyd4NEekmVJKJsdHkK3UkAtDrc failed: custom program error: 0x1771"
	require.True(r, slices.Contains(txIncNonce.Meta.LogMessages, expectedLog))
}
