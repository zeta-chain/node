package runner

import (
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/pkg/chains"
	solanacontract "github.com/zeta-chain/zetacore/pkg/contract/solana"
)

// SetupSolanaAccount imports the deployer's private key
func (r *E2ERunner) SetupSolanaAccount() {
	privateKey := solana.MustPrivateKeyFromBase58(r.Account.SolanaPrivateKey.String())
	r.SolanaDeployerAddress = privateKey.PublicKey()

	r.Logger.Info("SolanaDeployerAddress: %s", r.SolanaDeployerAddress)
}

// SetSolanaContracts set Solana contracts
func (r *E2ERunner) SetSolanaContracts(deployerPrivateKey string) {
	r.Logger.Print("⚙️ deploying gateway program on Solana")

	// set Solana contracts
	r.GatewayProgram = solana.MustPublicKeyFromBase58(solanacontract.SolanaGatewayProgramID)

	// get deployer account balance
	privkey := solana.MustPrivateKeyFromBase58(deployerPrivateKey)
	bal, err := r.SolanaClient.GetBalance(r.Ctx, privkey.PublicKey(), rpc.CommitmentFinalized)
	require.NoError(r, err)
	r.Logger.Info("deployer address: %s, balance: %f SOL", privkey.PublicKey().String(), float64(bal.Value)/1e9)

	// compute the gateway PDA address
	pdaComputed := r.ComputePdaAddress()

	// create 'initialize' instruction
	var inst solana.GenericInstruction
	accountSlice := []*solana.AccountMeta{}
	accountSlice = append(accountSlice, solana.Meta(privkey.PublicKey()).WRITE().SIGNER())
	accountSlice = append(accountSlice, solana.Meta(pdaComputed).WRITE())
	accountSlice = append(accountSlice, solana.Meta(solana.SystemProgramID))
	accountSlice = append(accountSlice, solana.Meta(r.GatewayProgram))
	inst.ProgID = r.GatewayProgram
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
	r.Logger.Info("initialize logs: %v", out.Meta.LogMessages)

	// retrieve the PDA account info
	pdaInfo, err := r.SolanaClient.GetAccountInfo(r.Ctx, pdaComputed)
	require.NoError(r, err)

	// deserialize the PDA info
	pda := solanacontract.PdaInfo{}
	err = borsh.Deserialize(&pda, pdaInfo.Bytes())
	require.NoError(r, err)
	fmt.Printf("pda parsed: %+v\n", pda)
	tssAddress := ethcommon.BytesToAddress(pda.TssAddress[:])

	// check the TSS address
	require.Equal(r, r.TSSAddress, tssAddress, "TSS address mismatch")

	// show the PDA balance
	balance, err := r.SolanaClient.GetBalance(r.Ctx, pdaComputed, rpc.CommitmentConfirmed)
	require.NoError(r, err)
	r.Logger.Info("initial PDA balance: %d lamports", balance.Value)
}
