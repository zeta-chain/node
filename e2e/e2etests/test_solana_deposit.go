package e2etests

import (
	"math/big"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestSolanaDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse deposit amount (in lamports)
	depositAmount, ok := new(big.Int).SetString(args[0], 10)
	require.True(r, ok, "Invalid deposit amount specified for TestSolanaDeposit.")

	// load deployer private key
	privkey := solana.MustPrivateKeyFromBase58(r.Account.SolanaPrivateKey.String())

	// create 'deposit' instruction
	instruction := r.CreateDepositInstruction(privkey.PublicKey(), r.EVMAddress(), depositAmount.Uint64())

	// create and sign the transaction
	signedTx := r.CreateSignedTransaction([]solana.Instruction{instruction}, privkey)

	// broadcast the transaction and wait for finalization
	sig, out := r.BroadcastTxSync(signedTx)
	r.Logger.Info("deposit logs: %v", out.Meta.LogMessages)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)
}
