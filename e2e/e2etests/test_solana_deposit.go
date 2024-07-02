package e2etests

import (
	"context"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
	"github.com/zeta-chain/zetacore/e2e/runner"
)

func TestSolanaInitializeGateway(r *runner.E2ERunner, args []string) {
	if len(args) != 0 {
		panic("TestSolanaIntializeGateway requires exactly zero argument for the amount.")
	}

	client := r.SolanaClient
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
		ChainId:       111111,
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
	}
	pdaInfo, err := client.GetAccountInfo(context.TODO(), pdaComputed)
	if err != nil {
		r.Logger.Print("error getting PDA info: %v", err)
		panic(err)
	}
	var pda PdaInfo
	borsh.Deserialize(&pda, pdaInfo.Bytes())

	r.Logger.Print("PDA info Tss: %v", pda.TssAddress)

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
		Amount:        1338,
		Memo:          []byte("hello this is a good memo for you to enjoy"),
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

	// wait for the cctx to be mined
	//cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	//r.Logger.CCTX(*cctx, "deposit")
	//if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
	//	panic(fmt.Sprintf(
	//		"expected mined status; got %s, message: %s",
	//		cctx.CctxStatus.Status.String(),
	//		cctx.CctxStatus.StatusMessage),
	//	)
	//}

}
