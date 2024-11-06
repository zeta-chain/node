package runner

import (
	"fmt"
	"math/big"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/utils"
	solanacontract "github.com/zeta-chain/node/pkg/contracts/solana"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// ComputePdaAddress computes the PDA address for the gateway program
func (r *E2ERunner) ComputePdaAddress() solana.PublicKey {
	seed := []byte(solanacontract.PDASeed)
	GatewayProgramID := solana.MustPublicKeyFromBase58(solanacontract.SolanaGatewayProgramID)
	pdaComputed, bump, err := solana.FindProgramAddress([][]byte{seed}, GatewayProgramID)
	require.NoError(r, err)

	r.Logger.Info("computed pda: %s, bump %d\n", pdaComputed, bump)

	return pdaComputed
}

// CreateDepositInstruction creates a 'deposit' instruction
func (r *E2ERunner) CreateDepositInstruction(
	signer solana.PublicKey,
	receiver ethcommon.Address,
	data []byte,
	amount uint64,
) solana.Instruction {
	// compute the gateway PDA address
	pdaComputed := r.ComputePdaAddress()
	programID := r.GatewayProgram

	// create 'deposit' instruction
	inst := &solana.GenericInstruction{}
	accountSlice := []*solana.AccountMeta{}
	accountSlice = append(accountSlice, solana.Meta(signer).WRITE().SIGNER())
	accountSlice = append(accountSlice, solana.Meta(pdaComputed).WRITE())
	accountSlice = append(accountSlice, solana.Meta(solana.SystemProgramID))
	inst.ProgID = programID
	inst.AccountValues = accountSlice

	var err error
	inst.DataBytes, err = borsh.Serialize(solanacontract.DepositInstructionParams{
		Discriminator: solanacontract.DiscriminatorDeposit,
		Amount:        amount,
		Memo:          append(receiver.Bytes(), data...),
	})
	require.NoError(r, err)

	return inst
}

func (r *E2ERunner) CreateWhitelistSPLMintInstruction(
	signer solana.PublicKey,
	whitelistEntry solana.PublicKey,
	whitelistCandidate solana.PublicKey,
) solana.Instruction {
	// compute the gateway PDA address
	pdaComputed := r.ComputePdaAddress()
	programID := r.GatewayProgram

	// create 'whitelist_spl_mint' instruction
	inst := &solana.GenericInstruction{}
	accountSlice := []*solana.AccountMeta{}
	accountSlice = append(accountSlice, solana.Meta(whitelistEntry).WRITE())
	accountSlice = append(accountSlice, solana.Meta(whitelistCandidate))
	accountSlice = append(accountSlice, solana.Meta(pdaComputed).WRITE())
	accountSlice = append(accountSlice, solana.Meta(signer).WRITE().SIGNER())
	accountSlice = append(accountSlice, solana.Meta(solana.SystemProgramID))
	inst.ProgID = programID
	inst.AccountValues = accountSlice

	var err error
	inst.DataBytes, err = borsh.Serialize(solanacontract.WhitelistInstructionParams{
		Discriminator: solanacontract.DiscriminatorWhitelistSplMint,
		// remaining fields are empty because no tss signature is needed if signer is admin account
	})
	require.NoError(r, err)

	return inst
}

/*

 #[account(mut)]
    pub signer: Signer<'info>,

    #[account(seeds = [b"meta"], bump)]
    pub pda: Account<'info, Pda>,

    #[account(seeds=[b"whitelist", mint_account.key().as_ref()], bump)]
    pub whitelist_entry: Account<'info, WhitelistEntry>, // attach whitelist entry to show the mint_account is whitelisted

    pub mint_account: Account<'info, Mint>,

    pub token_program: Program<'info, Token>,

    #[account(mut)]
    pub from: Account<'info, TokenAccount>, // this must be owned by signer; normally the ATA of signer
    #[account(mut)]
    pub to: Account<'info, TokenAccount>, // this must be ATA of PDA
*/

func (r *E2ERunner) CreateDepositSPLInstruction(
	amount uint64,
	signer solana.PublicKey,
	whitelistEntry solana.PublicKey,
	mint solana.PublicKey,
	from solana.PublicKey,
	to solana.PublicKey,
	receiver ethcommon.Address,
	data []byte,
) solana.Instruction {
	// compute the gateway PDA address
	pdaComputed := r.ComputePdaAddress()
	programID := r.GatewayProgram

	// create 'deposit_spl' instruction
	inst := &solana.GenericInstruction{}
	accountSlice := []*solana.AccountMeta{}
	accountSlice = append(accountSlice, solana.Meta(signer).WRITE().SIGNER())
	accountSlice = append(accountSlice, solana.Meta(pdaComputed))
	accountSlice = append(accountSlice, solana.Meta(whitelistEntry))
	accountSlice = append(accountSlice, solana.Meta(mint))
	accountSlice = append(accountSlice, solana.Meta(solana.TokenProgramID))
	accountSlice = append(accountSlice, solana.Meta(from).WRITE())
	accountSlice = append(accountSlice, solana.Meta(to).WRITE())
	inst.ProgID = programID
	inst.AccountValues = accountSlice

	var err error
	inst.DataBytes, err = borsh.Serialize(solanacontract.DepositInstructionParams{
		Discriminator: solanacontract.DiscriminatorDepositSPL,
		Amount:        amount,
		Memo:          append(receiver.Bytes(), data...),
	})
	require.NoError(r, err)

	return inst
}

// CreateSignedTransaction creates a signed transaction from instructions
func (r *E2ERunner) CreateSignedTransaction(
	instructions []solana.Instruction,
	privateKey solana.PrivateKey,
	additionalPrivateKeys []solana.PrivateKey,
) *solana.Transaction {
	// get a recent blockhash
	recent, err := r.SolanaClient.GetLatestBlockhash(r.Ctx, rpc.CommitmentFinalized)
	require.NoError(r, err)

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

func (r *E2ERunner) DepositSPL(privateKey *solana.PrivateKey, tokenAccount solana.Wallet, receiver ethcommon.Address, data []byte) solana.Signature {
	// ata for pda
	pda := r.ComputePdaAddress()
	pdaAta, _, err := solana.FindAssociatedTokenAddress(pda, tokenAccount.PublicKey())
	require.NoError(r, err)

	ata, _, err := solana.FindAssociatedTokenAddress(privateKey.PublicKey(), tokenAccount.PublicKey())
	require.NoError(r, err)

	ataInstruction := associatedtokenaccount.NewCreateInstruction(privateKey.PublicKey(), pda, tokenAccount.PublicKey()).Build()
	signedTx := r.CreateSignedTransaction(
		[]solana.Instruction{ataInstruction},
		*privateKey,
		[]solana.PrivateKey{tokenAccount.PrivateKey},
	)
	// broadcast the transaction and wait for finalization
	_, out := r.BroadcastTxSync(signedTx)
	r.Logger.Info("pda ata spl logs: %v", out.Meta.LogMessages)

	_, err = r.SolanaClient.GetTokenAccountBalance(r.Ctx, pdaAta, rpc.CommitmentConfirmed)
	require.NoError(r, err)

	// deposit spl
	seed := [][]byte{[]byte("whitelist"), tokenAccount.PublicKey().Bytes()}
	whitelistEntryPDA, _, err := solana.FindProgramAddress(seed, r.GatewayProgram)
	require.NoError(r, err)

	depositSPLInstruction := r.CreateDepositSPLInstruction(uint64(500_000), privateKey.PublicKey(), whitelistEntryPDA, tokenAccount.PublicKey(), ata, pdaAta, receiver, data)
	signedTx = r.CreateSignedTransaction(
		[]solana.Instruction{depositSPLInstruction},
		*privateKey,
		[]solana.PrivateKey{},
	)
	// broadcast the transaction and wait for finalization
	sig, out := r.BroadcastTxSync(signedTx)
	r.Logger.Info("deposit spl logs: %v", out.Meta.LogMessages)

	_, err = r.SolanaClient.GetTokenAccountBalance(r.Ctx, pdaAta, rpc.CommitmentConfirmed)
	require.NoError(r, err)

	return sig
}

func (r *E2ERunner) DeploySPL(privateKey *solana.PrivateKey, whitelist bool) *solana.Wallet {
	lamport, err := r.SolanaClient.GetMinimumBalanceForRentExemption(r.Ctx, token.MINT_SIZE, rpc.CommitmentFinalized)
	require.NoError(r, err)

	// to deploy new spl token, create account instruction and initialize mint instruction have to be in the same transaction
	tokenAccount := solana.NewWallet()
	createAccountInstruction := system.NewCreateAccountInstruction(
		lamport,
		token.MINT_SIZE,
		solana.TokenProgramID,
		privateKey.PublicKey(),
		tokenAccount.PublicKey(),
	).Build()

	initializeMintInstruction := token.NewInitializeMint2Instruction(
		6,
		privateKey.PublicKey(),
		privateKey.PublicKey(),
		tokenAccount.PublicKey(),
	).Build()

	signedTx := r.CreateSignedTransaction(
		[]solana.Instruction{createAccountInstruction, initializeMintInstruction},
		*privateKey,
		[]solana.PrivateKey{tokenAccount.PrivateKey},
	)

	// broadcast the transaction and wait for finalization
	_, out := r.BroadcastTxSync(signedTx)
	r.Logger.Info("create spl logs: %v", out.Meta.LogMessages)

	if whitelist {
		seed := [][]byte{[]byte("whitelist"), tokenAccount.PublicKey().Bytes()}
		whitelistEntryPDA, _, err := solana.FindProgramAddress(seed, r.GatewayProgram)
		require.NoError(r, err)

		whitelistEntryInfo, err := r.SolanaClient.GetAccountInfo(r.Ctx, whitelistEntryPDA)
		require.Error(r, err)

		// already whitelisted
		if whitelistEntryInfo != nil {
			return tokenAccount
		}

		// create 'whitelist_spl_mint' instruction
		instruction := r.CreateWhitelistSPLMintInstruction(privateKey.PublicKey(), whitelistEntryPDA, tokenAccount.PublicKey())
		// create and sign the transaction
		signedTx := r.CreateSignedTransaction([]solana.Instruction{instruction}, *privateKey, []solana.PrivateKey{})

		// broadcast the transaction and wait for finalization
		_, out := r.BroadcastTxSync(signedTx)
		r.Logger.Info("whitelist spl mint logs: %v", out.Meta.LogMessages)

		whitelistEntryInfo, err = r.SolanaClient.GetAccountInfo(r.Ctx, whitelistEntryPDA)
		require.NoError(r, err)
		require.NotNil(r, whitelistEntryInfo)

		fmt.Println("minting tokens to deployer...")

		ata, _, err := solana.FindAssociatedTokenAddress(privateKey.PublicKey(), tokenAccount.PublicKey())
		require.NoError(r, err)

		ataInstruction := associatedtokenaccount.NewCreateInstruction(privateKey.PublicKey(), privateKey.PublicKey(), tokenAccount.PublicKey()).Build()
		signedTx = r.CreateSignedTransaction(
			[]solana.Instruction{ataInstruction},
			*privateKey,
			[]solana.PrivateKey{tokenAccount.PrivateKey},
		)
		// broadcast the transaction and wait for finalization
		_, out = r.BroadcastTxSync(signedTx)
		r.Logger.Info("ata spl logs: %v", out.Meta.LogMessages)

		_, err = r.SolanaClient.GetTokenAccountBalance(r.Ctx, ata, rpc.CommitmentConfirmed)
		require.NoError(r, err)

		amount := uint64(1_000_000)
		mintToInstruction := token.NewMintToInstruction(amount, tokenAccount.PublicKey(), ata, privateKey.PublicKey(), []solana.PublicKey{}).Build()
		signedTx = r.CreateSignedTransaction(
			[]solana.Instruction{mintToInstruction},
			*privateKey,
			[]solana.PrivateKey{},
		)

		// broadcast the transaction and wait for finalization
		_, out = r.BroadcastTxSync(signedTx)
		r.Logger.Info("mint spl logs: %v", out.Meta.LogMessages)

		_, err = r.SolanaClient.GetTokenAccountBalance(r.Ctx, ata, rpc.CommitmentConfirmed)
		require.NoError(r, err)

	}

	return tokenAccount
}

// BroadcastTxSync broadcasts a transaction and waits for it to be finalized
func (r *E2ERunner) BroadcastTxSync(tx *solana.Transaction) (solana.Signature, *rpc.GetTransactionResult) {
	// broadcast the transaction
	sig, err := r.SolanaClient.SendTransactionWithOpts(r.Ctx, tx, rpc.TransactionOpts{})
	require.NoError(r, err)
	r.Logger.Info("broadcast success! tx sig %s; waiting for confirmation...", sig)

	var (
		start   = time.Now()
		timeout = 2 * time.Minute // Solana tx expires automatically after 2 minutes
	)

	// wait for the transaction to be finalized
	var out *rpc.GetTransactionResult
	for {
		require.False(r, time.Since(start) > timeout, "waiting solana tx timeout")

		time.Sleep(1 * time.Second)
		out, err = r.SolanaClient.GetTransaction(r.Ctx, sig, &rpc.GetTransactionOpts{})
		if err == nil {
			break
		}
	}

	return sig, out
}

// SOLDepositAndCall deposits an amount of ZRC20 SOL tokens (in lamports) and calls a contract (if data is provided)
func (r *E2ERunner) SOLDepositAndCall(
	signerPrivKey *solana.PrivateKey,
	receiver ethcommon.Address,
	amount *big.Int,
	data []byte,
) solana.Signature {
	// if signer is not provided, use the runner account as default
	if signerPrivKey == nil {
		privkey, err := solana.PrivateKeyFromBase58(r.Account.SolanaPrivateKey.String())
		require.NoError(r, err)
		signerPrivKey = &privkey
	}

	// create 'deposit' instruction
	instruction := r.CreateDepositInstruction(signerPrivKey.PublicKey(), receiver, data, amount.Uint64())

	// create and sign the transaction
	signedTx := r.CreateSignedTransaction([]solana.Instruction{instruction}, *signerPrivKey, []solana.PrivateKey{})

	// broadcast the transaction and wait for finalization
	sig, out := r.BroadcastTxSync(signedTx)
	r.Logger.Info("deposit logs: %v", out.Meta.LogMessages)

	return sig
}

// WithdrawSOLZRC20 withdraws an amount of ZRC20 SOL tokens
func (r *E2ERunner) WithdrawSOLZRC20(
	to solana.PublicKey,
	amount *big.Int,
	approveAmount *big.Int,
) *crosschaintypes.CrossChainTx {
	// approve
	tx, err := r.SOLZRC20.Approve(r.ZEVMAuth, r.SOLZRC20Addr, approveAmount)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "approve")

	// withdraw
	tx, err = r.SOLZRC20.Withdraw(r.ZEVMAuth, []byte(to.String()), amount)
	require.NoError(r, err)
	r.Logger.EVMTransaction(*tx, "withdraw")

	// wait for tx receipt
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "withdraw")
	r.Logger.Info("Receipt txhash %s status %d", receipt.TxHash, receipt.Status)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	return cctx
}
