package e2etests

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestSPLDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)
	amount := utils.ParseBigInt(r, args[0])
	require.True(r, amount.IsUint64(), fmt.Sprintf("arg[0] is not a uint64: %s", args[0]))

	// load deployer private key
	privKey := r.GetSolanaPrivKey()

	// get SPL balance for pda and sender atas
	pda := r.ComputePdaAddress()
	pdaAta := r.ResolveSolanaATA(privKey, pda, r.SPLAddr)

	pdaBalanceResult, err := r.SolanaClient.GetTokenAccountBalance(r.Ctx, pdaAta, rpc.CommitmentConfirmed)
	require.NoError(r, err)
	pdaBalanceBefore := utils.ParseBigInt(r, pdaBalanceResult.Value.Amount)

	senderAta := r.ResolveSolanaATA(privKey, privKey.PublicKey(), r.SPLAddr)
	senderBalanceResult, err := r.SolanaClient.GetTokenAccountBalance(r.Ctx, senderAta, rpc.CommitmentConfirmed)
	require.NoError(r, err)
	senderBalanceBefore := utils.ParseBigInt(r, senderBalanceResult.Value.Amount)

	// get zrc20 balance for recipient
	zrc20BalanceBefore, err := r.SPLZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)

	// deposit SPL tokens
	// #nosec G115 e2eTest - always in range
	sig := r.SPLDepositAndCall(&privKey, amount.Uint64(), r.SPLAddr, r.EVMAddress(), nil, nil)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "solana_deposit_spl")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)
	require.Equal(r, cctx.GetCurrentOutboundParam().Receiver, r.EVMAddress().Hex())

	// verify amount is deposited to pda ata
	pdaChange := utils.NewExactChange(amount)
	r.WaitAndVerifySPLBalanceChange(pdaAta, pdaBalanceBefore, pdaChange)

	// verify amount is subtracted from sender ata
	senderChange := utils.NewExactChange(new(big.Int).Neg(amount))
	r.WaitAndVerifySPLBalanceChange(senderAta, senderBalanceBefore, senderChange)

	// wait for the zrc20 balance to be updated
	change := utils.NewExactChange(amount)
	utils.WaitAndVerifyZRC20BalanceChange(r, r.SPLZRC20, r.EVMAddress(), zrc20BalanceBefore, change, r.Logger)
}
