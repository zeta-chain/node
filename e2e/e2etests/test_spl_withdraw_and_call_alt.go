package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	solanacontract "github.com/zeta-chain/node/pkg/contracts/solana"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestSPLWithdrawAndCallAddressLookupTable executes spl withdrawAndCall on zevm and calls connected program on solana
// similar to TestSPLWithdrawAndCall, but uses AddressLookupTable to provide accounts for connected program
func TestSPLWithdrawAndCallAddressLookupTable(r *runner.E2ERunner, args []string) {
	require.True(r, len(args) == 1 || len(args) == 3)

	var (
		addressLookupTableAddress solana.PublicKey
		writableIndexes           []uint8
		initAddressLookupTable    bool
	)
	if len(args) == 3 {
		var err error
		addressLookupTableAddress, err = solana.PublicKeyFromBase58(args[1])
		require.NoError(r, err, "invalid AddressLookupTable address")
		writableIndexes = utils.ParseUint8Array(r, args[2])
	} else {
		initAddressLookupTable = true
	}

	withdrawAmount := utils.ParseBigInt(r, args[0])

	// get SPL ZRC20 balance before withdraw
	zrc20BalanceBefore, err := r.SPLZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SPL before withdraw: %d", zrc20BalanceBefore)

	require.Equal(r, 1, zrc20BalanceBefore.Cmp(withdrawAmount), "Insufficient balance for withdrawal")

	// parse withdraw amount (in lamports), approve amount is 1 SOL
	approvedAmount := new(big.Int).SetUint64(solana.LAMPORTS_PER_SOL)
	require.Equal(
		r,
		-1,
		withdrawAmount.Cmp(approvedAmount),
		"Withdrawal amount must be less than the %v",
		approvedAmount,
	)

	// load deployer private key
	privkey := r.GetSolanaPrivKey()

	connected := solana.MustPublicKeyFromBase58(r.ConnectedSPLProgram.String())
	connectedPda, err := solanacontract.ComputeConnectedPdaAddress(connected)
	require.NoError(r, err)

	connectedPdaAta := r.ResolveSolanaATA(privkey, connectedPda, r.SPLAddr)
	connectedPdaBalanceBefore, err := r.SolanaClient.GetTokenAccountBalance(
		r.Ctx,
		connectedPdaAta,
		rpc.CommitmentConfirmed,
	)
	require.NoError(r, err)
	r.Logger.Info("connected pda balance of SPL before withdraw: %s", connectedPdaBalanceBefore.Value.Amount)

	// in case AddressLookupTable address is not provided (eg. for local testing) use prefunded random accounts to create AddressLookupTable
	randomWallets := []solana.PublicKey{}
	if initAddressLookupTable {
		accounts := make([]solana.PublicKey, 0, 6)

		accounts = append(accounts, connectedPda)
		accounts = append(accounts, connectedPdaAta)
		accounts = append(accounts, r.SPLAddr)
		accounts = append(accounts, r.ComputePdaAddress())
		accounts = append(accounts, solana.TokenProgramID)
		accounts = append(accounts, solana.SystemProgramID)
		predefinedAccountsLen := len(accounts)
		addressLookupTableAddress, randomWallets = r.SetupTestAddressLookupTableWithRandomWalletsSPL(accounts)
		writableIndexes = make([]uint8, 0, 2+len(randomWallets))
		writableIndexes = append(writableIndexes, 0, 1) // only first 2 are mutable

		// based on example accounts from above, all random wallets are writable
		// since they will get some SPL from connected program example
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

	msgEncoded, err := msg.Encode()
	require.NoError(r, err)

	// get receivers ata balances before withdraw
	randomWalletsBalanceBefore := make([]*big.Int, 0, len(randomWallets))
	for _, acc := range randomWallets {
		receiverBalanceBefore, err := r.SolanaClient.GetTokenAccountBalance(r.Ctx, acc, rpc.CommitmentConfirmed)
		require.NoError(r, err)
		randomWalletsBalanceBefore = append(
			randomWalletsBalanceBefore,
			utils.ParseBigInt(r, receiverBalanceBefore.Value.Amount),
		)
	}

	// withdraw
	tx := r.WithdrawAndCallSPLZRC20(
		withdrawAmount,
		approvedAmount,
		msgEncoded,
		gatewayzevm.RevertOptions{
			OnRevertGasLimit: big.NewInt(0),
		},
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// get SPL ZRC20 balance after withdraw
	zrc20BalanceAfter, err := r.SPLZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SPL after withdraw: %d", zrc20BalanceAfter)

	// check pda account info of connected program
	pda := r.ParseConnectedPda(connectedPda)
	require.Equal(r, msgDataStr, pda.LastMessage)
	require.Equal(r, r.ZEVMAuth.From.String(), common.BytesToAddress(pda.LastSender[:]).String())

	// verify balances are updated
	connectedPdaBalanceAfter, err := r.SolanaClient.GetTokenAccountBalance(
		r.Ctx,
		connectedPdaAta,
		rpc.CommitmentConfirmed,
	)
	require.NoError(r, err)
	r.Logger.Info("connected pda balance of SPL after withdraw: %s", connectedPdaBalanceAfter.Value.Amount)

	// verify half of amount is added to connected pda ata
	halfWithdrawAmount := new(big.Int).Div(withdrawAmount, big.NewInt(2))
	require.EqualValues(
		r,
		new(big.Int).Add(halfWithdrawAmount, utils.ParseBigInt(r, connectedPdaBalanceBefore.Value.Amount)).String(),
		utils.ParseBigInt(r, connectedPdaBalanceAfter.Value.Amount).String(),
	)

	// verify amount is subtracted on zrc20
	require.EqualValues(r, new(big.Int).Sub(zrc20BalanceBefore, withdrawAmount).String(), zrc20BalanceAfter.String())

	// check if balances locally are increased in connected program
	for i, acc := range randomWallets {
		receiverBalanceAfter, err := r.SolanaClient.GetTokenAccountBalance(r.Ctx, acc, rpc.CommitmentConfirmed)
		require.NoError(r, err)
		require.True(r, utils.ParseBigInt(r, receiverBalanceAfter.Value.Amount).Cmp(randomWalletsBalanceBefore[i]) > 0)
	}
}
