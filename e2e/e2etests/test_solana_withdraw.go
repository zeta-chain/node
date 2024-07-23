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

func TestSolanaWithdraw(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	r.Logger.Print("TestSolanaWithdraw...sol zrc20 %s", r.SOLZRC20Addr.String())

	solZRC20 := r.SOLZRC20
	supply, err := solZRC20.BalanceOf(&bind.CallOpts{}, r.ZEVMAuth.From)
	if err != nil {
		r.Logger.Error("Error getting total supply of sol zrc20: %v", err)
		panic(err)
	}
	r.Logger.Print(" from %s supply of %s sol zrc20: %d", r.ZEVMAuth.From, r.EVMAddress(), supply)

	// parse withdraw amount (in lamports), approve amount is 1 SOL
	approvedAmount := new(big.Int).SetUint64(solana.LAMPORTS_PER_SOL)
	withdrawAmount, ok := new(big.Int).SetString(args[0], 10)
	require.True(r, ok, "Invalid withdrawal amount specified for TestSolanaWithdraw.")
	require.Equal(
		r,
		-1,
		withdrawAmount.Cmp(approvedAmount),
		"Withdrawal amount must be less than the approved amount (1e9).",
	)

	// load deployer private key
	privkey := solana.MustPrivateKeyFromBase58(r.Account.SolanaPrivateKey.String())
	r.Logger.Print("TestSolanaWithdraw...sol zrc20 %s", r.SOLZRC20Addr.String())

	// approve
	tx, err := r.SOLZRC20.Approve(r.ZEVMAuth, r.SOLZRC20Addr, big.NewInt(1e18))
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// withdraw
	tx, err = r.SOLZRC20.Withdraw(r.ZEVMAuth, []byte(privkey.PublicKey().String()), withdrawAmount)
	require.NoError(r, err)
	r.Logger.EVMTransaction(*tx, "withdraw")

	// wait for tx receipt
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)
	r.Logger.Print("Receipt txhash %s status %d", receipt.TxHash, receipt.Status)

	supply, err = solZRC20.BalanceOf(&bind.CallOpts{}, r.ZEVMAuth.From)
	if err != nil {
		r.Logger.Error("Error getting total supply of sol zrc20: %v", err)
		panic(err)
	}
	r.Logger.Print(" from %s supply of %s sol zrc20 after: %d", r.ZEVMAuth.From, r.EVMAddress(), supply)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)
}
