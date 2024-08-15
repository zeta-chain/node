package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestSolanaDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// get ERC20 SOL balance before deposit
	balanceBefore, err := r.SOLZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SOL before deposit: %d", balanceBefore)

	// parse deposit amount (in lamports)
	// #nosec G115 e2e - always in range
	depositAmount := big.NewInt(int64(parseInt(r, args[0])))

	// load deployer private key
	privkey, err := solana.PrivateKeyFromBase58(r.Account.SolanaPrivateKey.String())
	require.NoError(r, err)

	// create 'deposit' instruction
	instruction := r.CreateDepositInstruction(privkey.PublicKey(), r.EVMAddress(), nil, depositAmount.Uint64())

	// create and sign the transaction
	signedTx := r.CreateSignedTransaction([]solana.Instruction{instruction}, privkey)

	// broadcast the transaction and wait for finalization
	sig, out := r.BroadcastTxSync(signedTx)
	r.Logger.Info("deposit logs: %v", out.Meta.LogMessages)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// get ERC20 SOL balance after deposit
	balanceAfter, err := r.SOLZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SOL after deposit: %d", balanceAfter)

	// the runner balance should be increased by the deposit amount
	amountIncreased := new(big.Int).Sub(balanceAfter, balanceBefore)
	require.Equal(r, depositAmount.String(), amountIncreased.String())
}
