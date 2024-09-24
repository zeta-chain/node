package e2etests

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/contracts/testbank"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

func TestPrecompilesBankThroughContract(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	// Increase the gasLimit. It's required because of the gas consumed by precompiled functions.
	previousGasLimit := r.ZEVMAuth.GasLimit
	r.ZEVMAuth.GasLimit = 10_000_000
	defer func() {
		r.ZEVMAuth.GasLimit = previousGasLimit

		// // Reset the allowance to 0; this is needed when running upgrade tests where
		// // this test runs twice.
		// tx, err := r.ERC20ZRC20.Approve(r.ZEVMAuth, bank.ContractAddress, big.NewInt(0))
		// require.NoError(r, err)
		// receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		// utils.RequireTxSuccessful(r, receipt, "Resetting allowance failed")
	}()

	spender := r.EVMAddress()
	zrc20Address := r.ERC20ZRC20Addr
	totalAmount := big.NewInt(1e3)
	fmt.Printf("DEBUG: spender %v\n", spender)
	fmt.Printf("DEBUG: zrc20Address %v\n", zrc20Address)

	// Get ERC20ZRC20.
	txHash := r.DepositERC20WithAmountAndMessage(r.EVMAddress(), totalAmount, []byte{})
	utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)

	testBankAddr, tx, testBank, err := testbank.DeployTestBank(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	fmt.Printf("DEBUG: testBankAddr %v\n", testBankAddr)

	// Initial cosmos balance should be 0.
	balance, err := testBank.BalanceOf(&bind.CallOpts{Context: r.Ctx}, zrc20Address, spender)
	require.NoError(r, err)
	require.Equal(r, int64(0), balance.Int64())
	fmt.Printf("DEBUG: initial balance %v\n", balance)
}

func checkZRC20Balance(r *runner.E2ERunner, target common.Address) *big.Int {
	bankZRC20Balance, err := r.ERC20ZRC20.BalanceOf(&bind.CallOpts{Context: r.Ctx}, target)
	require.NoError(r, err, "Call ERC20ZRC20.BalanceOf")
	return bankZRC20Balance
}

func allowZRC20fromSpender(r *runner.E2ERunner, target common.Address, amount *big.Int) {
	tx, err := r.ERC20ZRC20.Approve(r.ZEVMAuth, target, amount)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "Approve ETHZRC20 bank allowance tx failed")
}
