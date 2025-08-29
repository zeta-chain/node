package runner

import (
	"math/big"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/utils"
	solanacontract "github.com/zeta-chain/node/pkg/contracts/solana"
)

// solanaNodeSyncTolerance is the time tolerance for the Solana nodes behind a RPC to be synced
const solanaNodeSyncTolerance = 30 * time.Second

// Connected programs used to test sol and spl withdraw and call
var ConnectedProgramID = solana.MustPublicKeyFromBase58("4xEw862A2SEwMjofPkUyd4NEekmVJKJsdHkK3UkAtDrc")
var ConnectedSPLProgramID = solana.MustPublicKeyFromBase58("8iUjRRhUCn8BjrvsWPfj8mguTe9L81ES4oAUApiF8JFC")

// ComputePdaAddress computes the PDA address for the gateway program
func (r *E2ERunner) ComputePdaAddress() solana.PublicKey {
	seed := []byte(solanacontract.PDASeed)
	pdaComputed, bump, err := solana.FindProgramAddress([][]byte{seed}, r.GatewayProgram)
	require.NoError(r, err)

	r.Logger.Info("computed pda for gateway program: %s, bump %d\n", pdaComputed, bump)

	return pdaComputed
}

// CreateDepositInstruction creates a 'deposit' instruction
func (r *E2ERunner) CreateDepositInstruction(
	signer solana.PublicKey,
	receiver ethcommon.Address,
	data []byte,
	amount uint64,
	revertOptions *solanacontract.RevertOptions,
) solana.Instruction {
	var err error
	var depositData []byte
	if data == nil {
		depositData, err = borsh.Serialize(solanacontract.DepositInstructionParams{
			Discriminator: solanacontract.DiscriminatorDeposit,
			Amount:        amount,
			Receiver:      receiver,
			RevertOptions: revertOptions,
		})
		require.NoError(r, err)
	} else {
		depositData, err = borsh.Serialize(solanacontract.DepositAndCallInstructionParams{
			Discriminator: solanacontract.DiscriminatorDepositAndCall,
			Amount:        amount,
			Receiver:      receiver,
			Memo:          data,
			RevertOptions: revertOptions,
		})
		require.NoError(r, err)
	}

	return &solana.GenericInstruction{
		ProgID:    r.GatewayProgram,
		DataBytes: depositData,
		AccountValues: []*solana.AccountMeta{
			solana.Meta(signer).WRITE().SIGNER(),
			solana.Meta(r.ComputePdaAddress()).WRITE(),
			solana.Meta(solana.SystemProgramID),
		},
	}
}

// CreateDepositInstructionThroughProgram creates a 'trigger_deposit' instruction for connected example program
func (r *E2ERunner) CreateDepositInstructionThroughProgram(
	signer solana.PublicKey,
	receiver ethcommon.Address,
	amount uint64,
	revertOptions *solanacontract.RevertOptions,
) solana.Instruction {
	triggerDepositDiscriminator := [8]byte{154, 34, 24, 72, 18, 230, 27, 82}
	var err error
	var depositData []byte
	depositData, err = borsh.Serialize(solanacontract.DepositInstructionParams{
		Discriminator: triggerDepositDiscriminator,
		Amount:        amount,
		Receiver:      receiver,
		RevertOptions: revertOptions,
	})
	require.NoError(r, err)

	return &solana.GenericInstruction{
		ProgID:    r.ConnectedProgram,
		DataBytes: depositData,
		AccountValues: []*solana.AccountMeta{
			solana.Meta(signer).WRITE().SIGNER(),
			solana.Meta(r.ComputePdaAddress()).WRITE(),
			solana.Meta(r.GatewayProgram),
			solana.Meta(solana.SystemProgramID),
		},
	}
}

// CreateSOLCallInstruction creates a 'call' instruction
func (r *E2ERunner) CreateSOLCallInstruction(
	signer solana.PublicKey,
	receiver ethcommon.Address,
	data []byte,
	revertOptions *solanacontract.RevertOptions,
) solana.Instruction {
	callData, err := borsh.Serialize(solanacontract.CallInstructionParams{
		Discriminator: solanacontract.DiscriminatorCall,
		Receiver:      receiver,
		Memo:          data,
		RevertOptions: revertOptions,
	})
	require.NoError(r, err)

	return &solana.GenericInstruction{
		ProgID:    r.GatewayProgram,
		DataBytes: callData,
		AccountValues: []*solana.AccountMeta{
			solana.Meta(signer).WRITE().SIGNER(),
		},
	}
}

// CreateWhitelistSPLMintInstruction creates a 'whitelist_spl_mint' instruction
func (r *E2ERunner) CreateWhitelistSPLMintInstruction(
	signer, whitelistEntry, whitelistCandidate solana.PublicKey,
) solana.Instruction {
	data, err := borsh.Serialize(solanacontract.WhitelistInstructionParams{
		Discriminator: solanacontract.DiscriminatorWhitelistSplMint,
	})
	require.NoError(r, err)

	return &solana.GenericInstruction{
		ProgID:    r.GatewayProgram,
		DataBytes: data,
		AccountValues: []*solana.AccountMeta{
			solana.Meta(signer).WRITE().SIGNER(),
			solana.Meta(r.ComputePdaAddress()).WRITE(),
			solana.Meta(whitelistEntry).WRITE(),
			solana.Meta(whitelistCandidate),
			solana.Meta(solana.SystemProgramID),
		},
	}
}

// CreateDepositSPLInstruction creates a 'deposit_spl' instruction
func (r *E2ERunner) CreateDepositSPLInstruction(
	amount uint64,
	signer solana.PublicKey,
	whitelistEntry solana.PublicKey,
	mint solana.PublicKey,
	from solana.PublicKey,
	to solana.PublicKey,
	receiver ethcommon.Address,
	data []byte,
	revertOptions *solanacontract.RevertOptions,
) solana.Instruction {
	var err error
	var depositSPLData []byte
	if data == nil {
		depositSPLData, err = borsh.Serialize(solanacontract.DepositSPLInstructionParams{
			Discriminator: solanacontract.DiscriminatorDepositSPL,
			Amount:        amount,
			Receiver:      receiver,
			RevertOptions: revertOptions,
		})
		require.NoError(r, err)
	} else {
		depositSPLData, err = borsh.Serialize(solanacontract.DepositSPLAndCallInstructionParams{
			Discriminator: solanacontract.DiscriminatorDepositSPLAndCall,
			Amount:        amount,
			Receiver:      receiver,
			Memo:          data,
			RevertOptions: revertOptions,
		})
		require.NoError(r, err)
	}

	return &solana.GenericInstruction{
		ProgID:    r.GatewayProgram,
		DataBytes: depositSPLData,
		AccountValues: []*solana.AccountMeta{
			solana.Meta(signer).WRITE().SIGNER(),
			solana.Meta(r.ComputePdaAddress()).WRITE(),
			solana.Meta(whitelistEntry),
			solana.Meta(mint),
			solana.Meta(solana.TokenProgramID),
			solana.Meta(from).WRITE(),
			solana.Meta(to).WRITE(),
			solana.Meta(solana.SystemProgramID),
		},
	}
}

// CreateSignedTransaction creates a signed transaction from instructions
func (r *E2ERunner) CreateSignedTransaction(
	instructions []solana.Instruction,
	privateKey solana.PrivateKey,
	additionalPrivateKeys []solana.PrivateKey,
) *solana.Transaction {
	// get a recent blockhash
	recent, err := r.SolanaClient.GetLatestBlockhash(r.Ctx, rpc.CommitmentConfirmed)
	require.NoError(r, err)

	r.Logger.Info("Latest valid block height for tx %d", recent.Value.LastValidBlockHeight)

	// create the initialize transaction
	tx, err := solana.NewTransaction(
		instructions,
		recent.Value.Blockhash,
		solana.TransactionPayer(privateKey.PublicKey()),
	)
	require.NoError(r, err)

	// sign the initialize transaction
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if privateKey.PublicKey().Equals(key) {
				return &privateKey
			}
			for _, apk := range additionalPrivateKeys {
				if apk.PublicKey().Equals(key) {
					return &apk
				}
			}
			return nil
		},
	)
	require.NoError(r, err)

	return tx
}

// ResolveSolanaATA finds or creates SOL associated token account
func (r *E2ERunner) ResolveSolanaATA(
	payer solana.PrivateKey,
	owner solana.PublicKey,
	mintAccount solana.PublicKey,
) solana.PublicKey {
	pdaAta, _, err := solana.FindAssociatedTokenAddress(owner, mintAccount)
	require.NoError(r, err)

	info, _ := r.SolanaClient.GetAccountInfoWithOpts(
		r.Ctx,
		pdaAta,
		&rpc.GetAccountInfoOpts{Commitment: rpc.CommitmentConfirmed},
	)
	if info != nil {
		// already exists
		return pdaAta
	}
	// doesn't exist, create it
	ataInstruction := associatedtokenaccount.NewCreateInstruction(payer.PublicKey(), owner, mintAccount).Build()
	signedTx := r.CreateSignedTransaction(
		[]solana.Instruction{ataInstruction},
		payer,
		[]solana.PrivateKey{},
	)
	// broadcast the transaction and wait for finalization
	r.BroadcastTxSync(signedTx)

	return pdaAta
}

// SPLDepositAndCall deposits an amount of SPL tokens and calls a contract (if data is provided)
func (r *E2ERunner) SPLDepositAndCall(
	privateKey *solana.PrivateKey,
	amount uint64,
	mintAccount solana.PublicKey,
	receiver ethcommon.Address,
	data []byte,
	revertOptions *solanacontract.RevertOptions,
) solana.Signature {
	// ata for pda
	pda := r.ComputePdaAddress()
	pdaAta := r.ResolveSolanaATA(*privateKey, pda, mintAccount)

	// deployer ata
	ata := r.ResolveSolanaATA(*privateKey, privateKey.PublicKey(), mintAccount)

	// deposit spl
	seed := [][]byte{[]byte("whitelist"), mintAccount.Bytes()}
	whitelistEntryPDA, _, err := solana.FindProgramAddress(seed, r.GatewayProgram)
	require.NoError(r, err)

	depositSPLInstruction := r.CreateDepositSPLInstruction(
		amount,
		privateKey.PublicKey(),
		whitelistEntryPDA,
		mintAccount,
		ata,
		pdaAta,
		receiver,
		data,
		revertOptions,
	)

	limit := computebudget.NewSetComputeUnitLimitInstruction(500000).Build() // 500k compute unit limit
	feesInit := computebudget.NewSetComputeUnitPriceInstructionBuilder().
		SetMicroLamports(100000).Build() // 0.1 lamports per compute unit
	signedTx := r.CreateSignedTransaction(
		[]solana.Instruction{limit, feesInit, depositSPLInstruction},
		*privateKey,
		[]solana.PrivateKey{},
	)
	// broadcast the transaction and wait for finalization
	sig, out := r.BroadcastTxSync(signedTx)
	r.Logger.Info("deposit spl logs: %v", out.Meta.LogMessages)

	return sig
}

func (r *E2ERunner) DeploySPL(privateKey *solana.PrivateKey, whitelist bool) *solana.Wallet {
	lamport, err := r.SolanaClient.GetMinimumBalanceForRentExemption(r.Ctx, token.MINT_SIZE, rpc.CommitmentConfirmed)
	require.NoError(r, err)

	// to deploy new spl token, create account instruction and initialize mint instruction have to be in the same transaction
	mintAccount := solana.NewWallet()
	createAccountInstruction := system.NewCreateAccountInstruction(
		lamport,
		token.MINT_SIZE,
		solana.TokenProgramID,
		privateKey.PublicKey(),
		mintAccount.PublicKey(),
	).Build()

	initializeMintInstruction := token.NewInitializeMint2Instruction(
		6,
		privateKey.PublicKey(),
		privateKey.PublicKey(),
		mintAccount.PublicKey(),
	).Build()

	signedTx := r.CreateSignedTransaction(
		[]solana.Instruction{createAccountInstruction, initializeMintInstruction},
		*privateKey,
		[]solana.PrivateKey{mintAccount.PrivateKey},
	)

	// broadcast the transaction and wait for finalization
	_, out := r.BroadcastTxSync(signedTx)
	r.Logger.Info("create spl logs: %v", out.Meta.LogMessages)

	// minting some tokens to deployer for testing
	ata := r.ResolveSolanaATA(*privateKey, privateKey.PublicKey(), mintAccount.PublicKey())

	mintToInstruction := token.NewMintToInstruction(uint64(100_000_000_000_000), mintAccount.PublicKey(), ata, privateKey.PublicKey(), []solana.PublicKey{}).
		Build()
	signedTx = r.CreateSignedTransaction(
		[]solana.Instruction{mintToInstruction},
		*privateKey,
		[]solana.PrivateKey{},
	)

	// broadcast the transaction and wait for finalization
	_, out = r.BroadcastTxSync(signedTx)
	r.Logger.Info("mint spl logs: %v", out.Meta.LogMessages)

	// optionally whitelist spl token in gateway
	if whitelist {
		seed := [][]byte{[]byte("whitelist"), mintAccount.PublicKey().Bytes()}
		whitelistEntryPDA, _, err := solana.FindProgramAddress(seed, r.GatewayProgram)
		require.NoError(r, err)

		whitelistEntryInfo, err := r.SolanaClient.GetAccountInfoWithOpts(
			r.Ctx,
			whitelistEntryPDA,
			&rpc.GetAccountInfoOpts{Commitment: rpc.CommitmentConfirmed},
		)
		require.Error(r, err)

		// already whitelisted
		if whitelistEntryInfo != nil {
			return mintAccount
		}

		// create 'whitelist_spl_mint' instruction
		instruction := r.CreateWhitelistSPLMintInstruction(
			privateKey.PublicKey(),
			whitelistEntryPDA,
			mintAccount.PublicKey(),
		)
		// create and sign the transaction
		signedTx := r.CreateSignedTransaction([]solana.Instruction{instruction}, *privateKey, []solana.PrivateKey{})

		// broadcast the transaction and wait for finalization
		_, out := r.BroadcastTxSync(signedTx)
		r.Logger.Info("whitelist spl mint logs: %v", out.Meta.LogMessages)

		whitelistEntryInfo, err = r.SolanaClient.GetAccountInfoWithOpts(
			r.Ctx,
			whitelistEntryPDA,
			&rpc.GetAccountInfoOpts{
				Commitment: rpc.CommitmentConfirmed,
			},
		)
		require.NoError(r, err)
		require.NotNil(r, whitelistEntryInfo)
	}

	return mintAccount
}

// BroadcastTxSyncOnce broadcasts a transaction once and checks if it's confirmed
func (r *E2ERunner) BroadcastTxSyncOnce(tx *solana.Transaction) (solana.Signature, *rpc.GetTransactionResult, bool) {
	// broadcast the transaction
	r.Logger.Info("Broadcast once start")
	maxRetries := uint(1)
	sig, err := r.SolanaClient.SendTransactionWithOpts(r.Ctx, tx, rpc.TransactionOpts{
		SkipPreflight:       true,
		MaxRetries:          &maxRetries,
		PreflightCommitment: rpc.CommitmentConfirmed,
	})
	if err != nil { // try to fetch tx to see if error is not because it is already broadcasted, since we manually retry
		r.Logger.Info("Error sending tx %s, check if it's already broadcasted, err: %s", sig, err.Error())

		out, errGet := r.SolanaClient.GetTransaction(r.Ctx, sig, &rpc.GetTransactionOpts{
			Commitment: rpc.CommitmentConfirmed,
		})

		if errGet == nil {
			return sig, out, true
		}

		r.Logger.Info("Error getting tx %s", errGet.Error())
		require.NoError(r, err) // fail the test with send tx error
	}
	r.Logger.Info("Broadcast success! tx sig %s; waiting for confirmation...", sig)

	// wait for the transaction to be finalized
	var out *rpc.GetTransactionResult
	time.Sleep(5 * time.Second) // wait a bit and check if its confirmed
	blockHeight, err := r.SolanaClient.GetBlockHeight(r.Ctx, rpc.CommitmentConfirmed)
	require.NoError(r, err)
	r.Logger.Info("Current block height %d", blockHeight)

	out, err = r.SolanaClient.GetTransaction(r.Ctx, sig, &rpc.GetTransactionOpts{
		Commitment: rpc.CommitmentConfirmed,
	})
	if err != nil {
		r.Logger.Info("Error getting tx %s", err.Error())
	}

	isConfirmed := err == nil
	r.Logger.Info("Broadcast once finished, tx: %s, confirmed: %t", sig, isConfirmed)
	return sig, out, isConfirmed
}

// BroadcastTxSync broadcasts a transaction and waits for it to be finalized
func (r *E2ERunner) BroadcastTxSync(tx *solana.Transaction) (solana.Signature, *rpc.GetTransactionResult) {
	r.Logger.Info("Broadcast start")
	start := time.Now()
	timeout := 2 * time.Minute // Expires after 2 mins
	sig, out, isConfirmed := r.BroadcastTxSyncOnce(tx)
	for {
		require.False(r, time.Since(start) > timeout, "solana tx timeout")

		if isConfirmed {
			r.Logger.Info("Tx broadcasted and confirmed")
			require.Nil(r, out.Meta.Err, out.Meta.LogMessages)
			return sig, out
		}

		r.Logger.Info("Manually retrying tx")
		sig, out, isConfirmed = r.BroadcastTxSyncOnce(tx)
	}
}

// SOLDepositAndCall deposits an amount of ZRC20 SOL tokens (in lamports) and calls a contract (if data is provided)
func (r *E2ERunner) SOLDepositAndCall(
	signerPrivKey *solana.PrivateKey,
	receiver ethcommon.Address,
	amount *big.Int,
	data []byte,
	revertOptions *solanacontract.RevertOptions,
) solana.Signature {
	// if signer is not provided, use the runner account as default
	if signerPrivKey == nil {
		privkey := r.GetSolanaPrivKey()
		signerPrivKey = &privkey
	}

	// create 'deposit' instruction
	instruction := r.CreateDepositInstruction(signerPrivKey.PublicKey(), receiver, data, amount.Uint64(), revertOptions)

	return r.solDepositAndCall(signerPrivKey, instruction)
}

// SOLDepositAndCallThroughProgram deposits an amount of ZRC20 SOL tokens (in lamports) through program
func (r *E2ERunner) SOLDepositAndCallThroughProgram(
	signerPrivKey *solana.PrivateKey,
	receiver ethcommon.Address,
	amount *big.Int,
	revertOptions *solanacontract.RevertOptions,
) solana.Signature {
	// if signer is not provided, use the runner account as default
	if signerPrivKey == nil {
		privkey := r.GetSolanaPrivKey()
		signerPrivKey = &privkey
	}

	// create 'deposit' instruction
	instruction := r.CreateDepositInstructionThroughProgram(
		signerPrivKey.PublicKey(),
		receiver,
		amount.Uint64(),
		revertOptions,
	)

	return r.solDepositAndCall(signerPrivKey, instruction)
}

func (r *E2ERunner) solDepositAndCall(
	signerPrivKey *solana.PrivateKey,
	instruction solana.Instruction,
) solana.Signature {
	// create and sign the transaction
	limit := computebudget.NewSetComputeUnitLimitInstruction(500000).Build() // 500k compute unit limit
	feesInit := computebudget.NewSetComputeUnitPriceInstructionBuilder().
		SetMicroLamports(100000).Build() // 0.1 lamports per compute unit
	signedTx := r.CreateSignedTransaction(
		[]solana.Instruction{limit, feesInit, instruction},
		*signerPrivKey,
		[]solana.PrivateKey{},
	)

	// broadcast the transaction and wait for finalization
	sig, out := r.BroadcastTxSync(signedTx)
	r.Logger.Info("deposit logs: %v", out.Meta.LogMessages)

	return sig
}

// SOLCall calls a contract on zevm
func (r *E2ERunner) SOLCall(
	signerPrivKey *solana.PrivateKey,
	receiver ethcommon.Address,
	data []byte,
	revertOptions *solanacontract.RevertOptions,
) solana.Signature {
	// if signer is not provided, use the runner account as default
	if signerPrivKey == nil {
		privkey := r.GetSolanaPrivKey()
		signerPrivKey = &privkey
	}

	// create 'call' instruction
	instruction := r.CreateSOLCallInstruction(signerPrivKey.PublicKey(), receiver, data, revertOptions)

	// create and sign the transaction
	limit := computebudget.NewSetComputeUnitLimitInstruction(100000).Build() // 100k compute unit limit
	feesInit := computebudget.NewSetComputeUnitPriceInstructionBuilder().
		SetMicroLamports(100000).Build() // 0.1 lamports per compute unit
	signedTx := r.CreateSignedTransaction(
		[]solana.Instruction{limit, feesInit, instruction},
		*signerPrivKey,
		[]solana.PrivateKey{},
	)

	// broadcast the transaction and wait for finalization
	sig, out := r.BroadcastTxSync(signedTx)
	r.Logger.Info("call logs: %v", out.Meta.LogMessages)

	return sig
}

// WithdrawSOLZRC20 withdraws an amount of ZRC20 SOL tokens
func (r *E2ERunner) WithdrawSOLZRC20(
	to solana.PublicKey,
	amount *big.Int,
	approveAmount *big.Int,
	revertOptions gatewayzevm.RevertOptions,
) *ethtypes.Transaction {
	// approve
	tx, err := r.SOLZRC20.Approve(r.ZEVMAuth, r.GatewayZEVMAddr, approveAmount)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "approve")

	// withdraw
	tx, err = r.GatewayZEVM.Withdraw0(
		r.ZEVMAuth,
		[]byte(to.String()),
		amount,
		r.SOLZRC20Addr,
		revertOptions,
	)
	require.NoError(r, err)
	r.Logger.EVMTransaction(tx, "withdraw")

	// wait for tx receipt
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "withdraw")
	r.Logger.Info("Receipt txhash %s status %d", receipt.TxHash, receipt.Status)

	return tx
}

// WithdrawAndCallSOLZRC20 withdraws an amount of ZRC20 SOL tokens and calls program on solana
func (r *E2ERunner) WithdrawAndCallSOLZRC20(
	amount *big.Int,
	approveAmount *big.Int,
	msgEncoded []byte,
	revertOptions gatewayzevm.RevertOptions,
) *ethtypes.Transaction {
	receiver := r.ConnectedProgram.String()

	// approve
	tx, err := r.SOLZRC20.Approve(r.ZEVMAuth, r.GatewayZEVMAddr, approveAmount)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "approve")

	// withdraw
	tx, err = r.GatewayZEVM.WithdrawAndCall(
		r.ZEVMAuth,
		[]byte(receiver),
		amount,
		r.SOLZRC20Addr,
		msgEncoded,
		gatewayzevm.CallOptions{GasLimit: big.NewInt(250000)},
		revertOptions,
	)
	require.NoError(r, err)
	r.Logger.EVMTransaction(tx, "withdraw_and_call")

	// wait for tx receipt
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "withdraw_and_call")
	r.Logger.Info("Receipt txhash %s status %d", receipt.TxHash, receipt.Status)

	return tx
}

// CallSOLZRC20 calls program on solana
func (r *E2ERunner) CallSOLZRC20(
	to solana.PublicKey,
	approveAmount *big.Int,
	data []byte,
	revertOptions gatewayzevm.RevertOptions,
) *ethtypes.Transaction {
	// approve
	tx, err := r.SOLZRC20.Approve(r.ZEVMAuth, r.GatewayZEVMAddr, approveAmount)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "approve")

	// create encoded msg
	connectedPda, err := solanacontract.ComputeConnectedPdaAddress(r.ConnectedProgram)
	require.NoError(r, err)

	msg := solanacontract.ExecuteMsg{
		Accounts: []solanacontract.AccountMeta{
			{PublicKey: [32]byte(connectedPda.Bytes()), IsWritable: true},
			{PublicKey: [32]byte(r.ComputePdaAddress().Bytes()), IsWritable: false},
			{PublicKey: [32]byte(r.GetSolanaPrivKey().PublicKey().Bytes()), IsWritable: true},
			{PublicKey: [32]byte(solana.SystemProgramID.Bytes()), IsWritable: false},
			{PublicKey: [32]byte(solana.SysVarInstructionsPubkey.Bytes()), IsWritable: false},
		},
		Data: data,
	}

	msgEncoded, err := msg.Encode()
	require.NoError(r, err)

	// call
	tx, err = r.GatewayZEVM.Call(
		r.ZEVMAuth,
		[]byte(to.String()),
		r.SOLZRC20Addr,
		msgEncoded,
		gatewayzevm.CallOptions{GasLimit: big.NewInt(250000)},
		revertOptions,
	)
	require.NoError(r, err)
	utils.RequireTxSuccessful(r, receipt, "withdraw_and_call")
	r.Logger.Info("Receipt txhash %s status %d", receipt.TxHash, receipt.Status)

	return tx
}

// WithdrawSPLZRC20 withdraws an amount of ZRC20 SPL tokens
func (r *E2ERunner) WithdrawSPLZRC20(
	to solana.PublicKey,
	amount *big.Int,
	approveAmount *big.Int,
) *ethtypes.Transaction {
	// approve splzrc20 to spend gas tokens to pay gas fee
	tx, err := r.SOLZRC20.Approve(r.ZEVMAuth, r.SPLZRC20Addr, approveAmount)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "approve")

	// withdraw
	tx, err = r.SPLZRC20.Withdraw(r.ZEVMAuth, []byte(to.String()), amount)
	require.NoError(r, err)
	r.Logger.EVMTransaction(tx, "withdraw")

	// wait for tx receipt
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "withdraw")
	r.Logger.Info("Receipt txhash %s status %d", receipt.TxHash, receipt.Status)

	return tx
}

// WithdrawAndCallSPLZRC20 withdraws an amount of ZRC20 SPL tokens and calls program on solana
func (r *E2ERunner) WithdrawAndCallSPLZRC20(
	amount *big.Int,
	approveAmount *big.Int,
	data []byte,
	revertOptions gatewayzevm.RevertOptions,
) *ethtypes.Transaction {
	receiver := r.ConnectedSPLProgram.String()

	// approve
	tx, err := r.SOLZRC20.Approve(r.ZEVMAuth, r.GatewayZEVMAddr, approveAmount)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "approve")
	tx, err = r.SPLZRC20.Approve(r.ZEVMAuth, r.GatewayZEVMAddr, approveAmount)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "approve")

	// create encoded msg
	connected := solana.MustPublicKeyFromBase58(receiver)
	connectedPda, err := solanacontract.ComputeConnectedPdaAddress(connected)
	require.NoError(r, err)

	connectedPdaAta := r.ResolveSolanaATA(r.GetSolanaPrivKey(), connectedPda, r.SPLAddr)
	randomWalletAta := r.ResolveSolanaATA(r.GetSolanaPrivKey(), r.GetSolanaPrivKey().PublicKey(), r.SPLAddr)

	msg := solanacontract.ExecuteMsg{
		Accounts: []solanacontract.AccountMeta{
			{PublicKey: [32]byte(connectedPda.Bytes()), IsWritable: true},
			{PublicKey: [32]byte(connectedPdaAta.Bytes()), IsWritable: true},
			{PublicKey: [32]byte(r.SPLAddr), IsWritable: false},
			{PublicKey: [32]byte(r.ComputePdaAddress().Bytes()), IsWritable: false},
			{PublicKey: [32]byte(r.GetSolanaPrivKey().PublicKey().Bytes()), IsWritable: false},
			{PublicKey: [32]byte(randomWalletAta), IsWritable: true},
			{PublicKey: [32]byte(solana.TokenProgramID.Bytes()), IsWritable: false},
			{PublicKey: [32]byte(solana.SystemProgramID.Bytes()), IsWritable: false},
		},
		Data: data,
	}

	msgEncoded, err := msg.Encode()
	require.NoError(r, err)

	// withdraw
	tx, err = r.GatewayZEVM.WithdrawAndCall(
		r.ZEVMAuth,
		[]byte(receiver),
		amount,
		r.SPLZRC20Addr,
		msgEncoded,
		gatewayzevm.CallOptions{GasLimit: big.NewInt(250000)},
		revertOptions,
	)
	require.NoError(r, err)
	r.Logger.EVMTransaction(tx, "withdraw_and_call")

	// wait for tx receipt
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "withdraw_and_call")
	r.Logger.Info("Receipt txhash %s status %d", receipt.TxHash, receipt.Status)

	return tx
}

// WaitAndVerifySPLBalanceChange waits for the SPL balance of the given address to change by the given delta amount
// This function is to tolerate the fact that the balance update may not be synced across Solana nodes behind a RPC.
func (r *E2ERunner) WaitAndVerifySPLBalanceChange(
	ata solana.PublicKey,
	oldBalance *big.Int,
	change utils.BalanceChange,
) {
	// wait until the expected balance is reached or timeout
	startTime := time.Now()
	checkInterval := 2 * time.Second

	for {
		time.Sleep(checkInterval)
		require.False(r, time.Since(startTime) > solanaNodeSyncTolerance, "timeout waiting for SPL balance change")

		result, err := r.SolanaClient.GetTokenAccountBalance(r.Ctx, ata, rpc.CommitmentConfirmed)
		if err != nil {
			r.Logger.Info("unable to get SPL balance: %s", err.Error())
			continue
		}
		newBalance := utils.ParseBigInt(r, result.Value.Amount)

		if oldBalance.Cmp(newBalance) == 0 {
			r.Logger.Info("SPL balance has not changed yet")
			continue
		}
		r.Logger.Info("SPL balance changed from %d to %d on address %s", oldBalance, newBalance, ata.String())

		change.Verify(r, oldBalance, newBalance)

		return
	}
}
