package e2etests

import (
	"math/big"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	testcontract "github.com/zeta-chain/zetacore/testutil/contracts"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// TestSolanaDepositAndCallRefund tests deposit of lamports calling a example contract
func TestSolanaDepositAndCallRefund(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse deposit amount (in lamports)
	// #nosec G115 e2e - always in range
	depositAmount := big.NewInt(int64(parseInt(r, args[0])))

	// deploy a reverter contract in ZEVM
	r.Logger.Info("Deploying reverter contract")
	reverterAddr, _, _, err := testcontract.DeployReverter(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)

	r.Logger.Info("Reverter contract deployed at: %s", reverterAddr.String())

	// ---------------------------------------- execute the deposit transaction ----------------------------------------
	// load deployer private key
	privkey, err := solana.PrivateKeyFromBase58(r.Account.SolanaPrivateKey.String())
	require.NoError(r, err)

	// create 'deposit' instruction
	data := []byte("hello reverter")
	instruction := r.CreateDepositInstruction(privkey.PublicKey(), reverterAddr, data, depositAmount.Uint64())

	// create and sign the transaction
	signedTx := r.CreateSignedTransaction([]solana.Instruction{instruction}, privkey)

	// broadcast the transaction and wait for finalization
	sig, out := r.BroadcastTxSync(signedTx)
	r.Logger.Info("deposit logs: %v", out.Meta.LogMessages)

	// ---------------------------------------- verify the cross-chain revert --------------------------------------------
	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)
	r.Logger.Info("cross-chain call reverted: %v", cctx.CctxStatus.StatusMessage)

	// check the status message contains revert error hash in case of revert
	require.Contains(r, cctx.CctxStatus.StatusMessage, utils.ErrHashRevert)
}
