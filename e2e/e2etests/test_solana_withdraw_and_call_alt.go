package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	solanacontract "github.com/zeta-chain/node/pkg/contracts/solana"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestSolanaWithdrawAndCallALT executes withdrawAndCall on zevm and calls connected program on solana
// similar to TestSolanaWithdrawAndCall, but uses ALT to provide accounts for connected program
func TestSolanaWithdrawAndCallALT(r *runner.E2ERunner, args []string) {
	require.True(r, len(args) == 1 || len(args) == 3)

	withdrawAmount := utils.ParseBigInt(r, args[0])

	altAddress, err := solana.PublicKeyFromBase58(args[1])
	initALT := err != nil

	writableIndexes := utils.ParseUint8Array(r, args[2])

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

	// check balances before withdraw
	connectedPda, err := solanacontract.ComputeConnectedPdaAddress(r.ConnectedProgram)
	require.NoError(r, err)

	connectedPdaInfoBefore, err := r.SolanaClient.GetAccountInfo(r.Ctx, connectedPda)
	require.NoError(r, err)

	// in case ALT address is not provided (eg. for local testing) use prefunded random accounts to create ALT
	randomWallets := []solana.PublicKey{}
	if initALT {
		accounts := []solana.PublicKey{}

		accounts = append(accounts, connectedPda)
		accounts = append(accounts, r.ComputePdaAddress())
		accounts = append(accounts, solana.SystemProgramID)
		accounts = append(accounts, solana.SysVarInstructionsPubkey)
		predefinedAccountsLen := len(accounts)
		writableIndexes = []uint8{0} // only first one is mutable

		altAddress, randomWallets = r.SetupTestALTWithRandomWallets(accounts)

		// based on example accounts from above, all random wallets are writable
		// since they will get some lamports from connected program example
		for i := range randomWallets {
			// #nosec G115 e2eTest - always in range
			writableIndexes = append(
				writableIndexes,
				uint8(i+predefinedAccountsLen),
			)
		}
	}

	msg := solanacontract.ExecuteMsgALT{
		AltAddress:       altAddress,
		WriteableIndexes: writableIndexes,
		Data:             []byte("hello"),
	}

	encoded, err := msg.Encode()
	require.NoError(r, err)

	randomWalletsBalanceBefore := []uint64{}
	for _, acc := range randomWallets {
		balanceBefore, err := r.SolanaClient.GetAccountInfo(r.Ctx, acc)
		require.NoError(r, err)
		randomWalletsBalanceBefore = append(randomWalletsBalanceBefore, balanceBefore.Value.Lamports)
	}

	// withdraw and call
	tx := r.WithdrawAndCallSOLZRC20(
		withdrawAmount,
		approvedAmount,
		encoded,
		gatewayzevm.RevertOptions{
			OnRevertGasLimit: big.NewInt(0),
		},
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// get ERC20 SOL balance after withdraw
	balanceAfter, err := r.SOLZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SOL after withdraw: %d", balanceAfter)

	// check if the balance is reduced correctly
	amountReduced := new(big.Int).Sub(balanceBefore, balanceAfter)
	require.True(r, amountReduced.Cmp(withdrawAmount) >= 0, "balance is not reduced correctly")

	// check pda account info of connected program
	connectedPdaInfo, err := r.SolanaClient.GetAccountInfo(r.Ctx, connectedPda)
	require.NoError(r, err)

	type ConnectedPdaInfo struct {
		Discriminator     [8]byte
		LastSender        common.Address
		LastMessage       string
		LastRevertSender  solana.PublicKey
		LastRevertMessage string
	}
	pda := ConnectedPdaInfo{}
	err = borsh.Deserialize(&pda, connectedPdaInfo.Bytes())
	require.NoError(r, err)

	require.Equal(r, "hello", pda.LastMessage)
	require.Equal(r, r.ZEVMAuth.From.String(), common.BytesToAddress(pda.LastSender[:]).String())

	// check if balances locally are increased in connected program
	require.Greater(r, connectedPdaInfo.Value.Lamports, connectedPdaInfoBefore.Value.Lamports)
	for i, acc := range randomWallets {
		balanceAfter, err := r.SolanaClient.GetAccountInfo(r.Ctx, acc)
		require.NoError(r, err)
		require.Greater(r, balanceAfter.Value.Lamports, randomWalletsBalanceBefore[i])
	}
}
