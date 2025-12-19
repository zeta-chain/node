package e2etests

import (
	"encoding/binary"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gagliardetto/solana-go"
	addresslookuptable "github.com/gagliardetto/solana-go/programs/address-lookup-table"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	solanacontract "github.com/zeta-chain/node/pkg/contracts/solana"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestSolanaDepositThroughProgramAddressLookupTable executes deposit through connected program using AddressLookupTable
// similar to TestSolanaDepositThroughProgram, but uses AddressLookupTable to provide accounts for trigger_deposit
func TestSolanaDepositThroughProgramAddressLookupTable(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	depositAmount := utils.ParseBigInt(r, args[0])

	// get ERC20 SOL balance before deposit
	balanceBefore, err := r.SOLZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SOL before deposit: %d", balanceBefore)

	// get signer for deposit
	signerPrivKey := r.GetSolanaPrivKey()
	signer := signerPrivKey.PublicKey()

	// Accounts needed for trigger_deposit:
	// 1. signer (WRITE, SIGNER)
	// 2. gateway_pda (WRITE)
	// 3. gateway_program (read-only)
	// 4. system_program (read-only)
	accounts := []solana.PublicKey{
		signer,
		r.ComputePdaAddress(),
		r.GatewayProgram,
		solana.SystemProgramID,
	}

	// Create AddressLookupTable with only the required accounts
	privkey := r.GetSolanaPrivKey()
	recentSlot, err := r.SolanaClient.GetSlot(r.Ctx, rpc.CommitmentFinalized)
	require.NoError(r, err)

	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, recentSlot)

	addressLookupTableAddress, bump, err := solana.FindProgramAddress(
		[][]byte{privkey.PublicKey().Bytes(), buf},
		solana.AddressLookupTableProgramID,
	)
	require.NoError(r, err)

	// create AddressLookupTable
	createAddressLookupTableInstruction := addresslookuptable.NewCreateAddressLookupTableInstruction(
		recentSlot,
		bump,
		addressLookupTableAddress,
		privkey.PublicKey(),
		privkey.PublicKey(),
	).Build()

	signedTx := r.CreateSignedTransaction(
		[]solana.Instruction{createAddressLookupTableInstruction},
		privkey,
		[]solana.PrivateKey{},
	)
	r.BroadcastTxSync(signedTx)

	// need to wait a bit for AddressLookupTable to be active
	time.Sleep(1 * time.Second)

	// extend AddressLookupTable with accounts
	extendAddressLookupTableInstruction := addresslookuptable.NewExtendAddressLookupTableInstruction(
		accounts,
		addressLookupTableAddress,
		privkey.PublicKey(),
		privkey.PublicKey(),
	).Build()

	signedTx = r.CreateSignedTransaction(
		[]solana.Instruction{extendAddressLookupTableInstruction},
		privkey,
		[]solana.PrivateKey{},
	)
	r.BroadcastTxSync(signedTx)

	time.Sleep(1 * time.Second)

	// Get addresses from ALT table
	altState, err := addresslookuptable.GetAddressLookupTableStateWithOpts(
		r.Ctx,
		r.SolanaClient,
		addressLookupTableAddress,
		&rpc.GetAccountInfoOpts{Commitment: rpc.CommitmentConfirmed},
	)
	require.NoError(r, err)
	require.NotNil(r, altState, "ALT state should not be nil")

	// Create trigger_deposit instruction (accounts will be resolved from ALT by TransactionBuilder)
	triggerDepositDiscriminator := [8]byte{154, 34, 24, 72, 18, 230, 27, 82}
	depositData, err := borsh.Serialize(solanacontract.DepositInstructionParams{
		Discriminator: triggerDepositDiscriminator,
		Amount:        depositAmount.Uint64(),
		Receiver:      r.EVMAddress(),
		RevertOptions: nil,
	})
	require.NoError(r, err)

	instruction := &solana.GenericInstruction{
		ProgID:    r.ConnectedProgram,
		DataBytes: depositData,
		AccountValues: []*solana.AccountMeta{
			solana.Meta(signer).WRITE().SIGNER(),
			solana.Meta(r.ComputePdaAddress()).WRITE(),
			solana.Meta(r.GatewayProgram),
			solana.Meta(solana.SystemProgramID),
		},
	}

	// Get recent blockhash
	recent, err := r.SolanaClient.GetLatestBlockhash(r.Ctx, rpc.CommitmentConfirmed)
	require.NoError(r, err)

	// Build transaction with ALT using TransactionBuilder
	builder := solana.NewTransactionBuilder().
		SetRecentBlockHash(recent.Value.Blockhash).
		SetFeePayer(signer).
		WithOpt(solana.TransactionAddressTables(map[solana.PublicKey]solana.PublicKeySlice{
			addressLookupTableAddress: altState.Addresses,
		}))

	// Add compute budget instructions
	limit := computebudget.NewSetComputeUnitLimitInstruction(500000).Build()
	feesInit := computebudget.NewSetComputeUnitPriceInstructionBuilder().
		SetMicroLamports(100000).Build()

	builder.AddInstruction(limit)
	builder.AddInstruction(feesInit)
	builder.AddInstruction(instruction)

	// Build and sign transaction
	tx, err := builder.Build()
	require.NoError(r, err)

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if signer.Equals(key) {
				return &signerPrivKey
			}
			return nil
		},
	)
	require.NoError(r, err)

	// Broadcast transaction
	sig, out := r.BroadcastTxSync(tx)
	r.Logger.Info("deposit with ALT logs: %v", out.Meta.LogMessages)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "solana_deposit_alt")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)
	require.Equal(r, cctx.GetCurrentOutboundParam().Receiver, r.EVMAddress().Hex())

	// wait for the zrc20 balance to be updated
	change := utils.NewExactChange(depositAmount)
	utils.WaitAndVerifyZRC20BalanceChange(r, r.SOLZRC20, r.EVMAddress(), balanceBefore, change, r.Logger)
}
