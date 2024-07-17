package e2etests

import (
	"github.com/gagliardetto/solana-go"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestSolanaDeposit(r *runner.E2ERunner, _ []string) {
	// load deployer private key
	privkey := solana.MustPrivateKeyFromBase58(r.Account.SolanaPrivateKey.String())

	// create 'deposit' instruction
	amount := uint64(13370000)
	instruction := r.CreateDepositInstruction(privkey.PublicKey(), r.EVMAddress(), amount)

	// create and sign the transaction
	signedTx := r.CreateSignedTransaction([]solana.Instruction{instruction}, privkey)

	// broadcast the transaction and wait for finalization
	sig, out := r.BroadcastTxSync(signedTx)
	r.Logger.Info("deposit logs: %v", out.Meta.LogMessages)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)
}
