package e2etests

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/pkg/chains"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestSolanaInitializeGateway(r *runner.E2ERunner, args []string) {
	if len(args) != 0 {
		panic("TestSolanaIntializeGateway requires exactly zero argument for the amount.")
	}

	client := r.SolanaClient
	//r.Logger.Print("solana client URL", client.)
	if client == nil {
		r.Logger.Error("Solana client is nil")
		panic("Solana client is nil")
	}
	{
		res, err := client.GetVersion(context.Background())
		if err != nil {
			r.Logger.Error("error getting solana version: %v", err)
			panic(err)
		}
		r.Logger.Print("solana RPC version: %+v", res)
	}

	// building the transaction
	recent, err := client.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		panic(err)
	}
	r.Logger.Print("recent blockhash: %s", recent.Value.Blockhash)

	programId := solana.MustPublicKeyFromBase58("94U5AHQMKkV5txNJ17QPXWoh474PheGou6cNP2FEuL1d")
	seed := []byte("meta")
	pdaComputed, bump, err := solana.FindProgramAddress([][]byte{seed}, programId)
	if err != nil {
		panic(err)
	}
	r.Logger.Print("computed pda: %s, bump %d\n", pdaComputed, bump)

	privkey := solana.MustPrivateKeyFromBase58("4yqSQxDeTBvn86BuxcN5jmZb2gaobFXrBqu8kiE9rZxNkVMe3LfXmFigRsU4sRp7vk4vVP1ZCFiejDKiXBNWvs2C")
	r.Logger.Print("user pubkey: %s", privkey.PublicKey().String())
	bal, err := client.GetBalance(context.TODO(), privkey.PublicKey(), rpc.CommitmentFinalized)
	if err != nil {
		panic(err)
	}
	r.Logger.Print("account balance in SOL %f:", float64(bal.Value)/1e9)

	var inst solana.GenericInstruction
	accountSlice := []*solana.AccountMeta{}
	accountSlice = append(accountSlice, solana.Meta(privkey.PublicKey()).WRITE().SIGNER())
	accountSlice = append(accountSlice, solana.Meta(pdaComputed).WRITE())
	accountSlice = append(accountSlice, solana.Meta(solana.SystemProgramID))
	accountSlice = append(accountSlice, solana.Meta(programId))
	inst.ProgID = programId
	inst.AccountValues = accountSlice

	type InitializeParams struct {
		Discriminator [8]byte
		TssAddress    [20]byte
		ChainId       uint64
	}
	r.Logger.Print("TSS EthAddress: %s", r.TSSAddress)

	inst.DataBytes, err = borsh.Serialize(InitializeParams{
		Discriminator: [8]byte{175, 175, 109, 31, 13, 152, 155, 237},
		TssAddress:    r.TSSAddress,
		ChainId:       uint64(chains.SolanaLocalnet.ChainId),
	})
	if err != nil {
		panic(err)
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{&inst},
		recent.Value.Blockhash,
		solana.TransactionPayer(privkey.PublicKey()),
	)
	if err != nil {
		panic(err)
	}
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if privkey.PublicKey().Equals(key) {
				return &privkey
			}
			return nil
		},
	)
	if err != nil {
		panic(fmt.Errorf("unable to sign transaction: %w", err))
	}
	sig, err := client.SendTransactionWithOpts(
		context.TODO(),
		tx,
		rpc.TransactionOpts{},
	)
	if err != nil {
		panic(err)
	}
	r.Logger.Print("broadcast success! tx sig %s; waiting for confirmation...", sig)
	time.Sleep(16 * time.Second)
	type PdaInfo struct {
		Discriminator [8]byte
		Nonce         uint64
		TssAddress    [20]byte
		Authority     [32]byte
		ChainID       uint64
	}
	pdaInfo, err := client.GetAccountInfo(context.TODO(), pdaComputed)
	if err != nil {
		r.Logger.Print("error getting PDA info: %v", err)
		panic(err)
	}
	var pda PdaInfo
	borsh.Deserialize(&pda, pdaInfo.Bytes())

	r.Logger.Print("PDA info Tss: %v, chain id %d", pda.TssAddress, pda.ChainID)

}

func TestSolanaDeposit(r *runner.E2ERunner, args []string) {
	client := r.SolanaClient

	privkey := solana.MustPrivateKeyFromBase58("4yqSQxDeTBvn86BuxcN5jmZb2gaobFXrBqu8kiE9rZxNkVMe3LfXmFigRsU4sRp7vk4vVP1ZCFiejDKiXBNWvs2C")

	// build & bcast a Depsosit tx
	bal, err := client.GetBalance(context.TODO(), privkey.PublicKey(), rpc.CommitmentFinalized)
	if err != nil {
		r.Logger.Error("Error getting balance: %v", err)
		panic(fmt.Sprintf("Error getting balance: %v", err))
	}
	r.Logger.Print("account balance in SOL %f", float64(bal.Value)/1e9)

	// building the transaction
	recent, err := client.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		r.Logger.Error("Error getting recent blockhash: %v", err)
		panic(err)
	}
	r.Logger.Print("recent blockhash:", recent.Value.Blockhash)

	programId := solana.MustPublicKeyFromBase58("94U5AHQMKkV5txNJ17QPXWoh474PheGou6cNP2FEuL1d")
	seed := []byte("meta")
	pdaComputed, bump, err := solana.FindProgramAddress([][]byte{seed}, programId)
	if err != nil {
		r.Logger.Error("Error finding program address: %v", err)
		panic(err)
	}
	r.Logger.Print("computed pda: %s, bump %d\n", pdaComputed, bump)

	//pdaAccount := solana.MustPublicKeyFromBase58("4hA43LCh2Utef8EwCyWwYmWBoSeNq6RS2HdoLkWGm5z5")
	var inst solana.GenericInstruction
	accountSlice := []*solana.AccountMeta{}
	accountSlice = append(accountSlice, solana.Meta(privkey.PublicKey()).WRITE().SIGNER())
	accountSlice = append(accountSlice, solana.Meta(pdaComputed).WRITE())
	accountSlice = append(accountSlice, solana.Meta(solana.SystemProgramID))
	accountSlice = append(accountSlice, solana.Meta(programId))
	inst.ProgID = programId
	inst.AccountValues = accountSlice

	type DepositInstructionParams struct {
		Discriminator [8]byte
		Amount        uint64
		Memo          []byte
	}

	inst.DataBytes, err = borsh.Serialize(DepositInstructionParams{
		Discriminator: [8]byte{0xf2, 0x23, 0xc6, 0x89, 0x52, 0xe1, 0xf2, 0xb6},
		Amount:        13370000,
		Memo:          r.EVMAddress().Bytes(),
	})
	if err != nil {
		r.Logger.Error("Error serializing deposit instruction: %v", err)
		panic(err)
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{&inst},
		recent.Value.Blockhash,
		solana.TransactionPayer(privkey.PublicKey()),
	)
	if err != nil {
		r.Logger.Error("Error creating transaction: %v", err)
		panic(err)
	}
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if privkey.PublicKey().Equals(key) {
				return &privkey
			}
			return nil
		},
	)
	if err != nil {
		r.Logger.Error("Error signing transaction: %v", err)
		panic(fmt.Errorf("unable to sign transaction: %w", err))
	}

	//spew.Dump(tx)

	sig, err := client.SendTransactionWithOpts(
		context.TODO(),
		tx,
		rpc.TransactionOpts{},
	)
	if err != nil {
		r.Logger.Error("Error sending transaction: %v", err)
		panic(err)
	}
	r.Logger.Print("broadcast success! tx sig %s; waiting for confirmation...", sig)
	time.Sleep(16 * time.Second)

	//spew.Dump(sig)
	out, err := client.GetTransaction(context.TODO(), sig, &rpc.GetTransactionOpts{})
	if err != nil {
		r.Logger.Error("Error getting transaction: %v", err)
		panic(err)
	}
	r.Logger.Print("transaction status: %v, %v", out.Meta.Err, out.Meta.Status)
	r.Logger.Print("transaction logs: %v", out.Meta.LogMessages)

	//wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		panic(fmt.Sprintf(
			"expected mined status; got %s, message: %s",
			cctx.CctxStatus.Status.String(),
			cctx.CctxStatus.StatusMessage),
		)
	}

}

func TestSolanaWithdraw(r *runner.E2ERunner, args []string) {
	r.Logger.Print("TestSolanaWithdraw...sol zrc20 %s", r.SOLZRC20Addr.String())
	privkey := solana.MustPrivateKeyFromBase58("4yqSQxDeTBvn86BuxcN5jmZb2gaobFXrBqu8kiE9rZxNkVMe3LfXmFigRsU4sRp7vk4vVP1ZCFiejDKiXBNWvs2C")

	solZRC20 := r.SOLZRC20
	supply, err := solZRC20.BalanceOf(&bind.CallOpts{}, r.ZEVMAuth.From)
	if err != nil {
		r.Logger.Error("Error getting total supply of sol zrc20: %v", err)
		panic(err)
	}
	r.Logger.Print(" supply of %s sol zrc20: %d", r.EVMAddress(), supply)

	amount := big.NewInt(1337)
	approveAmount := big.NewInt(1e18)
	//r.Logger.Print("Approving %s sol zrc20 to spend %d", r.ZEVMAuth.From.Hex(), approveAmount)
	tx, err := r.SOLZRC20.Approve(r.ZEVMAuth, r.SOLZRC20Addr, approveAmount)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	tx, err = r.SOLZRC20.Withdraw(r.ZEVMAuth, []byte(privkey.PublicKey().String()), amount)
	require.NoError(r, err)
	r.Logger.EVMTransaction(*tx, "withdraw")
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)
	r.Logger.Print("Receipt txhash %s status %d", receipt.TxHash, receipt.Status)

	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
}
