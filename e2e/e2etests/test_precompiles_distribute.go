package e2etests

import (
	"math/big"

	"github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/precompiles/bank"
	"github.com/zeta-chain/node/precompiles/staking"
	ptypes "github.com/zeta-chain/node/precompiles/types"
)

func TestPrecompilesDistribute(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	var (
		spenderAddress            = r.EVMAddress()
		distributeContractAddress = staking.ContractAddress
		lockerAddress             = bank.ContractAddress

		zrc20Address = r.ERC20ZRC20Addr
		zrc20Denom   = ptypes.ZRC20ToCosmosDenom(zrc20Address)

		oneThousand    = big.NewInt(1e3)
		oneThousandOne = big.NewInt(1001)
		fiveHundred    = big.NewInt(500)
		fiveHundredOne = big.NewInt(501)

		previousGasLimit = r.ZEVMAuth.GasLimit
	)

	// Set new gas limit to avoid out of gas errors.
	r.ZEVMAuth.GasLimit = 10_000_000

	// Set the test to reset the state after it finishes.
	defer resetDistributionTest(r, lockerAddress, previousGasLimit, fiveHundred)

	// Get ERC20ZRC20.
	txHash := r.DepositERC20WithAmountAndMessage(spenderAddress, oneThousand, []byte{})
	utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)

	dstrContract, err := staking.NewIStaking(distributeContractAddress, r.ZEVMClient)
	require.NoError(r, err, "failed to create distribute contract caller")

	// DO NOT REMOVE - will be used in a subsequent PR when the ability to withdraw delegator rewards is introduced.
	// Get validators through staking contract.
	// validators, err := dstrContract.GetAllValidators(&bind.CallOpts{})
	// require.NoError(r, err)

	// Check initial balances.
	balanceShouldBe(r, 1000, checkZRC20Balance(r, spenderAddress))
	balanceShouldBe(r, 0, checkZRC20Balance(r, lockerAddress))
	balanceShouldBe(r, 0, checkCosmosBalance(r, r.FeeCollectorAddress, zrc20Denom))

	tx, err := dstrContract.Distribute(r.ZEVMAuth, zrc20Address, oneThousand)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt, "distribute should fail when there's no allowance")

	// Balances shouldn't change after a failed attempt.
	balanceShouldBe(r, 1000, checkZRC20Balance(r, spenderAddress))
	balanceShouldBe(r, 0, checkZRC20Balance(r, lockerAddress))
	balanceShouldBe(r, 0, checkCosmosBalance(r, r.FeeCollectorAddress, zrc20Denom))

	// Allow 500.
	approveAllowance(r, distributeContractAddress, fiveHundred)

	// Shouldn't be able to distribute more than allowed.
	tx, err = dstrContract.Distribute(r.ZEVMAuth, zrc20Address, fiveHundredOne)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt, "distribute should fail trying to distribute more than allowed")

	// Balances shouldn't change after a failed attempt.
	balanceShouldBe(r, 1000, checkZRC20Balance(r, spenderAddress))
	balanceShouldBe(r, 0, checkZRC20Balance(r, lockerAddress))
	balanceShouldBe(r, 0, checkCosmosBalance(r, r.FeeCollectorAddress, zrc20Denom))

	// Raise the allowance to 1000.
	approveAllowance(r, distributeContractAddress, oneThousand)

	// Shouldn't be able to distribute more than owned balance.
	tx, err = dstrContract.Distribute(r.ZEVMAuth, zrc20Address, oneThousandOne)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt, "distribute should fail trying to distribute more than owned balance")

	// Balances shouldn't change after a failed attempt.
	balanceShouldBe(r, 1000, checkZRC20Balance(r, spenderAddress))
	balanceShouldBe(r, 0, checkZRC20Balance(r, lockerAddress))
	balanceShouldBe(r, 0, checkCosmosBalance(r, r.FeeCollectorAddress, zrc20Denom))

	// Should be able to distribute 500, which is within balance and allowance.
	tx, err = dstrContract.Distribute(r.ZEVMAuth, zrc20Address, fiveHundred)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "distribute should succeed when distributing within balance and allowance")

	balanceShouldBe(r, 500, checkZRC20Balance(r, spenderAddress))
	balanceShouldBe(r, 500, checkZRC20Balance(r, lockerAddress))
	balanceShouldBe(r, 500, checkCosmosBalance(r, r.FeeCollectorAddress, zrc20Denom))

	eventDitributed, err := dstrContract.ParseDistributed(*receipt.Logs[0])
	require.NoError(r, err)
	require.Equal(r, zrc20Address, eventDitributed.Zrc20Token)
	require.Equal(r, spenderAddress, eventDitributed.Zrc20Distributor)
	require.Equal(r, fiveHundred.Uint64(), eventDitributed.Amount.Uint64())

	// After one block the rewards should have been distributed and fee collector should have 0 ZRC20 balance.
	r.WaitForBlocks(1)
	balanceShouldBe(r, 0, checkCosmosBalance(r, r.FeeCollectorAddress, zrc20Denom))

	// DO NOT REMOVE THE FOLLOWING CODE
	// This section is commented until a following PR introduces the ability to withdraw delegator rewards.
	// This validator checks will be used then to complete the whole e2e.

	// res, err := r.DistributionClient.ValidatorDistributionInfo(
	// 	r.Ctx,
	// 	&distributiontypes.QueryValidatorDistributionInfoRequest{
	// 		ValidatorAddress: validators[0].OperatorAddress,
	// 	},
	// )
	// require.NoError(r, err)
	// fmt.Printf("Validator 0 distribution info: %+v\n", res)

	// res2, err := r.DistributionClient.ValidatorOutstandingRewards(r.Ctx, &distributiontypes.QueryValidatorOutstandingRewardsRequest{
	// 	ValidatorAddress: validators[0].OperatorAddress,
	// })
	// require.NoError(r, err)
	// fmt.Printf("Validator 0 outstanding rewards: %+v\n", res2)

	// res3, err := r.DistributionClient.ValidatorCommission(r.Ctx, &distributiontypes.QueryValidatorCommissionRequest{
	// 	ValidatorAddress: validators[0].OperatorAddress,
	// })
	// require.NoError(r, err)
	// fmt.Printf("Validator 0 commission: %+v\n", res3)

	// // Validator 1
	// res, err = r.DistributionClient.ValidatorDistributionInfo(
	// 	r.Ctx,
	// 	&distributiontypes.QueryValidatorDistributionInfoRequest{
	// 		ValidatorAddress: validators[1].OperatorAddress,
	// 	},
	// )
	// require.NoError(r, err)
	// fmt.Printf("Validator 1 distribution info: %+v\n", res)

	// res2, err = r.DistributionClient.ValidatorOutstandingRewards(r.Ctx, &distributiontypes.QueryValidatorOutstandingRewardsRequest{
	// 	ValidatorAddress: validators[1].OperatorAddress,
	// })
	// require.NoError(r, err)
	// fmt.Printf("Validator 1 outstanding rewards: %+v\n", res2)

	// res3, err = r.DistributionClient.ValidatorCommission(r.Ctx, &distributiontypes.QueryValidatorCommissionRequest{
	// 	ValidatorAddress: validators[1].OperatorAddress,
	// })
	// require.NoError(r, err)
	// fmt.Printf("Validator 1 commission: %+v\n", res3)
}

func TestPrecompilesDistributeNonZRC20(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	// Increase the gasLimit. It's required because of the gas consumed by precompiled functions.
	previousGasLimit := r.ZEVMAuth.GasLimit
	r.ZEVMAuth.GasLimit = 10_000_000
	defer func() {
		r.ZEVMAuth.GasLimit = previousGasLimit
	}()

	spender, dstrAddress := r.EVMAddress(), staking.ContractAddress

	// Create a staking contract caller.
	dstrContract, err := staking.NewIStaking(dstrAddress, r.ZEVMClient)
	require.NoError(r, err, "Failed to create staking contract caller")

	// Deposit and approve 50 WZETA for the test.
	approveAmount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(50))
	r.DepositAndApproveWZeta(approveAmount)

	// Allow the staking contract to spend 25 WZeta tokens.
	tx, err := r.WZeta.Approve(r.ZEVMAuth, dstrAddress, big.NewInt(25))
	require.NoError(r, err, "Error approving allowance for staking contract")
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	require.EqualValues(r, uint64(1), receipt.Status, "approve allowance tx failed")

	// Check the allowance of the staking in WZeta tokens. Should be 25.
	allowance, err := r.WZeta.Allowance(&bind.CallOpts{Context: r.Ctx}, spender, dstrAddress)
	require.NoError(r, err, "Error retrieving staking allowance")
	require.EqualValues(r, uint64(25), allowance.Uint64(), "Error allowance for staking contract")

	// Call Distribute with 25 Non ZRC20 tokens. Should fail.
	tx, err = dstrContract.Distribute(r.ZEVMAuth, r.WZetaAddr, big.NewInt(25))
	require.NoError(r, err, "Error calling staking.distribute()")
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	require.Equal(r, uint64(0), receipt.Status, "Non ZRC20 deposit should fail")
}

// checkCosmosBalance checks the cosmos coin balance for an address. The coin is specified by its denom.
func checkCosmosBalance(r *runner.E2ERunner, address types.AccAddress, denom string) *big.Int {
	bal, err := r.BankClient.Balance(
		r.Ctx,
		&banktypes.QueryBalanceRequest{Address: address.String(), Denom: denom},
	)
	require.NoError(r, err)

	return bal.Balance.Amount.BigInt()
}

func resetDistributionTest(r *runner.E2ERunner, lockerAddress common.Address, previousGasLimit uint64, amount *big.Int) {
	r.ZEVMAuth.GasLimit = previousGasLimit

	// Reset the allowance to 0; this is needed when running upgrade tests where this test runs twice.
	tx, err := r.ERC20ZRC20.Approve(r.ZEVMAuth, lockerAddress, big.NewInt(0))
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "Resetting allowance failed")

	// Reset balance to 0 for spender; this is needed when running upgrade tests where this test runs twice.
	tx, err = r.ERC20ZRC20.Transfer(
		r.ZEVMAuth,
		common.HexToAddress("0x000000000000000000000000000000000000dEaD"),
		amount,
	)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "Resetting balance failed")
}
