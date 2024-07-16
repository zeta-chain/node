package e2etests

import (
	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	solanacontract "github.com/zeta-chain/zetacore/zetaclient/chains/solana/contract"
)

func TestSolanaDeposit(r *runner.E2ERunner, _ []string) {
	// load deployer private key
	privkey := solana.MustPrivateKeyFromBase58(r.Account.RawBase58PrivateKey.String())

	// compute the gateway PDA address
	pdaComputed := r.ComputePdaAddress()
	programID := r.GatewayProgramID()

	// create 'deposit' instruction
	var inst solana.GenericInstruction
	accountSlice := []*solana.AccountMeta{}
	accountSlice = append(accountSlice, solana.Meta(privkey.PublicKey()).WRITE().SIGNER())
	accountSlice = append(accountSlice, solana.Meta(pdaComputed).WRITE())
	accountSlice = append(accountSlice, solana.Meta(solana.SystemProgramID))
	accountSlice = append(accountSlice, solana.Meta(programID))
	inst.ProgID = programID
	inst.AccountValues = accountSlice

	var err error
	inst.DataBytes, err = borsh.Serialize(solanacontract.DepositInstructionParams{
		Discriminator: solanacontract.DiscriminatorDeposit(),
		Amount:        13370000,
		Memo:          r.EVMAddress().Bytes(),
	})
	require.NoError(r, err)

	// create and sign the transaction
	signedTx := r.CreateSignedTransaction([]solana.Instruction{&inst}, privkey)

	// broadcast the transaction and wait for finalization
	sig, out := r.BroadcastTxSync(signedTx)
	r.Logger.Print("deposit receiver address: %s", r.EVMAddress().String())
	r.Logger.Print("deposit logs: %v", out.Meta.LogMessages)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)
}
