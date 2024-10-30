package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/contracts/testdistribute"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/precompiles/bank"
	"github.com/zeta-chain/node/precompiles/staking"
	precompiletypes "github.com/zeta-chain/node/precompiles/types"
)

func TestPrecompilesDistributeThroughContract(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	var (
		spenderAddress            = r.EVMAddress()
		distributeContractAddress = staking.ContractAddress
		lockerAddress             = bank.ContractAddress

		zrc20Address = r.ERC20ZRC20Addr
		zrc20Denom   = precompiletypes.ZRC20ToCosmosDenom(zrc20Address)

		oneThousand    = big.NewInt(1e3)
		oneThousandOne = big.NewInt(1001)
		fiveHundred    = big.NewInt(500)
		fiveHundredOne = big.NewInt(501)

		previousGasLimit = r.ZEVMAuth.GasLimit
	)

	// Get ERC20ZRC20.
	txHash := r.DepositERC20WithAmountAndMessage(spenderAddress, oneThousand, []byte{})
	utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)

	dstrContract, err := staking.NewIStaking(distributeContractAddress, r.ZEVMClient)
	require.NoError(r, err, "failed to create distribute contract caller")

	_, tx, testDstrContract, err := testdistribute.DeployTestDistribute(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "deployment of disitributor caller contract failed")

	// Set new gas limit to avoid out of gas errors.
	r.ZEVMAuth.GasLimit = 10_000_000

	// Set the test to reset the state after it finishes.
	defer resetDistributionTest(r, lockerAddress, previousGasLimit, fiveHundred)

	// Check initial balances.
	balanceShouldBe(r, 1000, checkZRC20Balance(r, spenderAddress))
	balanceShouldBe(r, 500, checkZRC20Balance(r, lockerAddress)) // Carries 500 from distribute e2e.
	balanceShouldBe(r, 0, checkCosmosBalance(r, r.FeeCollectorAddress, zrc20Denom))

	receipt = distributeThroughContract(r, testDstrContract, zrc20Address, oneThousand)
	utils.RequiredTxFailed(r, receipt, "distribute should fail when there's no allowance")

	// Balances shouldn't change after a failed attempt.
	balanceShouldBe(r, 1000, checkZRC20Balance(r, spenderAddress))
	balanceShouldBe(r, 500, checkZRC20Balance(r, lockerAddress)) // Carries 500 from distribute e2e.
	balanceShouldBe(r, 0, checkCosmosBalance(r, r.FeeCollectorAddress, zrc20Denom))

	// Allow 500.
	approveAllowance(r, distributeContractAddress, fiveHundred)

	receipt = distributeThroughContract(r, testDstrContract, zrc20Address, fiveHundredOne)
	utils.RequiredTxFailed(r, receipt, "distribute should fail trying to distribute more than allowed")

	// Balances shouldn't change after a failed attempt.
	balanceShouldBe(r, 1000, checkZRC20Balance(r, spenderAddress))
	balanceShouldBe(r, 500, checkZRC20Balance(r, lockerAddress)) // Carries 500 from distribute e2e.
	balanceShouldBe(r, 0, checkCosmosBalance(r, r.FeeCollectorAddress, zrc20Denom))

	// Raise the allowance to 1000.
	approveAllowance(r, distributeContractAddress, oneThousand)

	// Shouldn't be able to distribute more than owned balance.
	receipt = distributeThroughContract(r, testDstrContract, zrc20Address, oneThousandOne)
	utils.RequiredTxFailed(r, receipt, "distribute should fail trying to distribute more than owned balance")

	// Balances shouldn't change after a failed attempt.
	balanceShouldBe(r, 1000, checkZRC20Balance(r, spenderAddress))
	balanceShouldBe(r, 500, checkZRC20Balance(r, lockerAddress)) // Carries 500 from distribute e2e.
	balanceShouldBe(r, 0, checkCosmosBalance(r, r.FeeCollectorAddress, zrc20Denom))

	// Should be able to distribute 500, which is within balance and allowance.
	receipt = distributeThroughContract(r, testDstrContract, zrc20Address, fiveHundred)
	utils.RequireTxSuccessful(r, receipt, "distribute should succeed when distributing within balance and allowance")

	balanceShouldBe(r, 500, checkZRC20Balance(r, spenderAddress))
	balanceShouldBe(r, 1000, checkZRC20Balance(r, lockerAddress)) // Carries 500 from distribute e2e.
	balanceShouldBe(r, 500, checkCosmosBalance(r, r.FeeCollectorAddress, zrc20Denom))

	eventDitributed, err := dstrContract.ParseDistributed(*receipt.Logs[0])
	require.NoError(r, err)
	require.Equal(r, zrc20Address, eventDitributed.Zrc20Token)
	require.Equal(r, spenderAddress, eventDitributed.Zrc20Distributor)
	require.Equal(r, fiveHundred.Uint64(), eventDitributed.Amount.Uint64())

	// After one block the rewards should have been distributed and fee collector should have 0 ZRC20 balance.
	r.WaitForBlocks(1)
	balanceShouldBe(r, 0, checkCosmosBalance(r, r.FeeCollectorAddress, zrc20Denom))
}

func distributeThroughContract(
	r *runner.E2ERunner,
	dstr *testdistribute.TestDistribute,
	zrc20Address common.Address,
	amount *big.Int,
) *types.Receipt {
	tx, err := dstr.DistributeThroughContract(r.ZEVMAuth, zrc20Address, amount)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	return receipt
}
