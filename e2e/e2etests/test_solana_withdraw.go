package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
)

func TestSolanaWithdraw(r *runner.E2ERunner, args []string) {
	// require exactly 0 argument
	require.Len(r, args, 0, "TestSolanaWithdraw currently takes no argument")

	// load deployer private key
	privkey := solana.MustPrivateKeyFromBase58(r.Account.SolanaPrivateKey.String())
	r.Logger.Print("TestSolanaWithdraw...sol zrc20 %s", r.SOLZRC20Addr.String())

	solZRC20 := r.SOLZRC20
	supply, err := solZRC20.BalanceOf(&bind.CallOpts{}, r.ZEVMAuth.From)
	require.NoError(r, err)
	r.Logger.Print(" supply of %s sol zrc20: %d", r.EVMAddress(), supply)

	amount := big.NewInt(1337)
	approveAmount := big.NewInt(1e18)
	tx, err := r.SOLZRC20.Approve(r.ZEVMAuth, r.SOLZRC20Addr, approveAmount)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	tx, err = r.SOLZRC20.Withdraw(r.ZEVMAuth, []byte(privkey.PublicKey().String()), amount)
	require.NoError(r, err)
	r.Logger.EVMTransaction(*tx, "withdraw")

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)
	r.Logger.Print("Receipt txhash %s status %d", receipt.TxHash, receipt.Status)

	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
}
