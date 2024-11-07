package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestSolanaDepositSPL(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)
	amount := parseInt(r, args[0])

	// load deployer private key
	privKey, err := solana.PrivateKeyFromBase58(r.Account.SolanaPrivateKey.String())
	require.NoError(r, err)

	// get SPL balance for pda and sender atas
	pda := r.ComputePdaAddress()
	pdaAta := r.FindOrCreateAssociatedTokenAccount(privKey, pda, r.SPLAddr)

	pdaBalanceBefore, err := r.SolanaClient.GetTokenAccountBalance(r.Ctx, pdaAta, rpc.CommitmentConfirmed)
	require.NoError(r, err)

	senderAta := r.FindOrCreateAssociatedTokenAccount(privKey, privKey.PublicKey(), r.SPLAddr)
	senderBalanceBefore, err := r.SolanaClient.GetTokenAccountBalance(r.Ctx, senderAta, rpc.CommitmentConfirmed)
	require.NoError(r, err)

	// get zrc20 balance for recepient
	zrc20BalanceBefore, err := r.SPLZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)

	// deposit SPL tokens
	sig := r.DepositSPL(&privKey, uint64(amount), r.SPLAddr, r.EVMAddress(), nil)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "solana_deposit_spl")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// verify balances are updated
	pdaBalanceAfter, err := r.SolanaClient.GetTokenAccountBalance(r.Ctx, pdaAta, rpc.CommitmentConfirmed)
	require.NoError(r, err)

	senderBalanceAfter, err := r.SolanaClient.GetTokenAccountBalance(r.Ctx, senderAta, rpc.CommitmentConfirmed)
	require.NoError(r, err)

	zrc20BalanceAfter, err := r.SPLZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)

	// verify amount is deposited to pda ata
	require.Equal(r, parseInt(r, pdaBalanceBefore.Value.Amount)+amount, parseInt(r, pdaBalanceAfter.Value.Amount))

	// verify amount is substracted from sender ata
	require.Equal(r, parseInt(r, senderBalanceBefore.Value.Amount)-amount, parseInt(r, senderBalanceAfter.Value.Amount))

	// verify amount is minted to receiver
	require.Zero(r, zrc20BalanceBefore.Add(zrc20BalanceBefore, big.NewInt(int64(amount))).Cmp(zrc20BalanceAfter))
}
