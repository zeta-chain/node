package e2etests

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"

	testcontract "github.com/zeta-chain/node/e2e/contracts/example"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestSPLDepositAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)
	amount := utils.ParseBigInt(r, args[0])
	require.True(r, amount.IsUint64(), fmt.Sprintf("arg[0] is not a uint64: %s", args[0]))

	// deploy an example contract in ZEVM
	contractAddr, _, contract, err := testcontract.DeployExample(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	r.Logger.Info("Example contract deployed at: %s", contractAddr.String())

	// load deployer private key
	privKey := r.GetSolanaPrivKey()

	// get SPL balance for pda and sender atas
	pda := r.ComputePdaAddress()
	pdaAta := r.ResolveSolanaATA(privKey, pda, r.SPLAddr)

	pdaBalanceBefore, err := r.SolanaClient.GetTokenAccountBalance(r.Ctx, pdaAta, rpc.CommitmentConfirmed)
	require.NoError(r, err)

	senderAta := r.ResolveSolanaATA(privKey, privKey.PublicKey(), r.SPLAddr)
	senderBalanceBefore, err := r.SolanaClient.GetTokenAccountBalance(r.Ctx, senderAta, rpc.CommitmentConfirmed)
	require.NoError(r, err)

	// get zrc20 balance for recipient
	zrc20BalanceBefore, err := r.SPLZRC20.BalanceOf(&bind.CallOpts{}, contractAddr)
	require.NoError(r, err)

	// execute the deposit transaction
	data := []byte("hello spl tokens")
	// #nosec G115 e2eTest - always in range
	sig := r.SPLDepositAndCall(&privKey, amount.Uint64(), r.SPLAddr, contractAddr, data, nil)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "solana_deposit_spl_and_call")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)
	require.Equal(r, cctx.GetCurrentOutboundParam().Receiver, contractAddr.Hex())

	// check if example contract has been called, bar value should be set to amount
	utils.WaitAndVerifyExampleContractCallWithMsg(
		r,
		contract,
		amount,
		data,
		[]byte(r.GetSolanaPrivKey().PublicKey().String()),
	)

	// verify balances are updated
	pdaBalanceAfter, err := r.SolanaClient.GetTokenAccountBalance(r.Ctx, pdaAta, rpc.CommitmentConfirmed)
	require.NoError(r, err)

	senderBalanceAfter, err := r.SolanaClient.GetTokenAccountBalance(r.Ctx, senderAta, rpc.CommitmentConfirmed)
	require.NoError(r, err)

	// verify amount is deposited to pda ata
	require.Equal(
		r,
		new(big.Int).Add(utils.ParseBigInt(r, pdaBalanceBefore.Value.Amount), amount),
		utils.ParseBigInt(r, pdaBalanceAfter.Value.Amount),
	)

	// verify amount is subtracted from sender ata
	require.Equal(
		r,
		new(big.Int).Sub(utils.ParseBigInt(r, senderBalanceBefore.Value.Amount), amount),
		utils.ParseBigInt(r, senderBalanceAfter.Value.Amount),
	)

	// wait for the zrc20 balance to be updated
	change := utils.NewExactChange(amount)
	utils.WaitAndVerifyZRC20BalanceChange(r, r.SPLZRC20, contractAddr, zrc20BalanceBefore, change, r.Logger)
}
