package e2etests

import (
	"context"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/pkg/chains"
	solanacontract "github.com/zeta-chain/zetacore/zetaclient/chains/solana/contract"
)

func TestSolanaInitializeGateway(r *runner.E2ERunner, args []string) {
	// no arguments expected
	require.Len(r, args, 0, "solana gateway initialization test should have no arguments")

	// print the solana node version
	client := r.SolanaClient
	res, err := client.GetVersion(context.Background())
	require.NoError(r, err)
	r.Logger.Print("solana version: %+v", res)

	// get deployer account balance
	privkey := solana.MustPrivateKeyFromBase58(r.Account.RawBase58PrivateKey.String())
	bal, err := client.GetBalance(context.TODO(), privkey.PublicKey(), rpc.CommitmentFinalized)
	require.NoError(r, err)
	r.Logger.Print("deployer address: %s, balance: %f SOL", privkey.PublicKey().String(), float64(bal.Value)/1e9)

	// compute the gateway PDA address
	pdaComputed := r.ComputePdaAddress()
	programID := r.GatewayProgramID()

	// create 'initialize' instruction
	var inst solana.GenericInstruction
	accountSlice := []*solana.AccountMeta{}
	accountSlice = append(accountSlice, solana.Meta(privkey.PublicKey()).WRITE().SIGNER())
	accountSlice = append(accountSlice, solana.Meta(pdaComputed).WRITE())
	accountSlice = append(accountSlice, solana.Meta(solana.SystemProgramID))
	accountSlice = append(accountSlice, solana.Meta(programID))
	inst.ProgID = programID
	inst.AccountValues = accountSlice

	inst.DataBytes, err = borsh.Serialize(solanacontract.InitializeParams{
		Discriminator: solanacontract.DiscriminatorInitialize(),
		TssAddress:    r.TSSAddress,
		ChainID:       uint64(chains.SolanaLocalnet.ChainId),
	})
	require.NoError(r, err)

	// create and sign the transaction
	signedTx := r.CreateSignedTransaction([]solana.Instruction{&inst}, privkey)

	// broadcast the transaction and wait for finalization
	_, out := r.BroadcastTxSync(signedTx)
	r.Logger.Print("initialize logs: %v", out.Meta.LogMessages)

	// retrieve the PDA account info
	pdaInfo, err := client.GetAccountInfo(context.TODO(), pdaComputed)
	require.NoError(r, err)

	// deserialize the PDA info
	pda := solanacontract.PdaInfo{}
	err = borsh.Deserialize(&pda, pdaInfo.Bytes())
	require.NoError(r, err)
	tssAddress := ethcommon.BytesToAddress(pda.TssAddress[:])

	// check the TSS address
	require.Equal(r, r.TSSAddress, tssAddress, "TSS address mismatch")
}
