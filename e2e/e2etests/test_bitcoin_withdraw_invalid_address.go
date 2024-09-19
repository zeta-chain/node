package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

func TestBitcoinWithdrawToInvalidAddress(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	withdrawalAmount := parseFloat(r, args[0])
	amount := btcAmountFromFloat64(r, withdrawalAmount)

	r.SetBtcAddress(r.Name, false)

	withdrawToInvalidAddress(r, amount)
}

func withdrawToInvalidAddress(r *runner.E2ERunner, amount *big.Int) {
	approvalAmount := 1000000000000000000
	// approve the ZRC20 contract to spend approvalAmount BTC from the deployer address.
	// the actual amount transferred is provided as test arg BTC, but we approve more to cover withdraw fee
	tx, err := r.BTCZRC20.Approve(r.ZEVMAuth, r.BTCZRC20Addr, big.NewInt(int64(approvalAmount)))
	require.NoError(r, err)

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// mine blocks if testing on regnet
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	// withdraw amount provided as test arg BTC from ZRC20 to BTC legacy address
	// the address "1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3" is for mainnet, not regtest
	tx, err = r.BTCZRC20.Withdraw(r.ZEVMAuth, []byte("1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3"), amount)
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt)
}
