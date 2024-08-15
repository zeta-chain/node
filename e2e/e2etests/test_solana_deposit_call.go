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

// TestSolanaDepositAndCall tests deposit of lamports calling a example contract
func TestSolanaDepositAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse deposit amount (in lamports)
	// #nosec G115 e2e - always in range
	depositAmount := big.NewInt(int64(parseInt(r, args[0])))

	// deploy an example contract in ZEVM
	contractAddr, _, contract, err := testcontract.DeployExample(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)

	r.Logger.Info("Example contract deployed at: %s", contractAddr.String())

	// ---------------------------------------- execute the deposit transaction ----------------------------------------
	// load deployer private key
	privkey, err := solana.PrivateKeyFromBase58(r.Account.SolanaPrivateKey.String())
	require.NoError(r, err)

	// create 'deposit' instruction
	data := []byte("hello lamports")
	instruction := r.CreateDepositInstruction(privkey.PublicKey(), contractAddr, data, depositAmount.Uint64())

	// create and sign the transaction
	signedTx := r.CreateSignedTransaction([]solana.Instruction{instruction}, privkey)

	// broadcast the transaction and wait for finalization
	sig, out := r.BroadcastTxSync(signedTx)
	r.Logger.Print("deposit logs: %v", out.Meta.LogMessages)

	// ---------------------------------------- verify the cross-chain call --------------------------------------------
	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// check if example contract has been called, bar value should be set to amount
	utils.MustHaveCalledExampleContract(r, contract, depositAmount)
	r.Logger.Info("cross-chain call succeeded")
}
