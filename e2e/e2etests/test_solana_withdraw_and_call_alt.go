package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	addresslookuptable "github.com/gagliardetto/solana-go/programs/address-lookup-table"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	solanacontract "github.com/zeta-chain/node/pkg/contracts/solana"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestSolanaWithdrawAndCallAddressLookupTable executes withdrawAndCall on zevm and calls connected program on solana
// similar to TestSolanaWithdrawAndCall, but uses AddressLookupTable to provide accounts for connected program
func TestSolanaWithdrawAndCallAddressLookupTable(r *runner.E2ERunner, args []string) {
	require.True(r, len(args) == 1 || len(args) == 2)

	var (
		addressLookupTableAddress solana.PublicKey
		writableIndexes           []uint8
		initAddressLookupTable    bool
		randomWallets             []solana.PublicKey
	)

	if len(args) == 2 {
		r.Logger.Info("using existing address lookup table")

		var err error
		addressLookupTableAddress, err = solana.PublicKeyFromBase58(args[1])
		require.NoError(r, err, "invalid AddressLookupTable address")

		predefinedAccountsLen := 4   // pda, connected pda, system program, sysvar instructions
		writableIndexes = []uint8{0} // only first one is mutable from predefined accounts
		alt, err := addresslookuptable.GetAddressLookupTableStateWithOpts(
			r.Ctx,
			r.SolanaClient,
			addressLookupTableAddress,
			&rpc.GetAccountInfoOpts{Commitment: rpc.CommitmentProcessed},
		)
		require.NoError(r, err)

		for i := predefinedAccountsLen; i < len(alt.Addresses); i++ {
			randomWallets = append(randomWallets, alt.Addresses[i])
			writableIndexes = append(
				writableIndexes,
				uint8(i),
			)
		}

	} else {
		initAddressLookupTable = true
	}

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

	// check balances before withdraw
	connectedPda, err := solanacontract.ComputeConnectedPdaAddress(r.ConnectedProgram)
	require.NoError(r, err)

	connectedPdaInfoBefore, err := r.SolanaClient.GetAccountInfo(r.Ctx, connectedPda)
	require.NoError(r, err)

	// in case AddressLookupTable address is not provided (eg. for local testing) use prefunded random accounts to create AddressLookupTable
	if initAddressLookupTable {
		r.Logger.Info("initializing new address lookup table")
		accounts := []solana.PublicKey{}

		accounts = append(accounts, connectedPda)
		accounts = append(accounts, r.ComputePdaAddress())
		accounts = append(accounts, solana.SystemProgramID)
		accounts = append(accounts, solana.SysVarInstructionsPubkey)
		predefinedAccountsLen := len(accounts)
		writableIndexes = []uint8{0} // only first one is mutable

		addressLookupTableAddress, randomWallets = r.SetupTestAddressLookupTableWithRandomWallets(accounts)

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

	msgDataStr := "hello"
	msg := solanacontract.ExecuteMsgAddressLookupTable{
		AddressLookupTableAddress: [32]byte(addressLookupTableAddress),
		WritableIndexes:           writableIndexes,
		Data:                      []byte(msgDataStr),
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

	pda := r.ParseConnectedPda(connectedPda)
	require.Equal(r, msgDataStr, pda.LastMessage)
	require.Equal(r, r.ZEVMAuth.From.String(), common.BytesToAddress(pda.LastSender[:]).String())

	// check if balances locally are increased in connected program
	require.Greater(r, connectedPdaInfo.Value.Lamports, connectedPdaInfoBefore.Value.Lamports)
	for i, acc := range randomWallets {
		balanceAfter, err := r.SolanaClient.GetAccountInfo(r.Ctx, acc)
		require.NoError(r, err)
		require.Greater(r, balanceAfter.Value.Lamports, randomWalletsBalanceBefore[i])
	}
}
