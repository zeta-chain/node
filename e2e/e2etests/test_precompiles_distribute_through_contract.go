package e2etests

import (
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/e2e/contracts/testdistribute"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/precompiles/bank"
	"github.com/zeta-chain/node/precompiles/staking"
	precompiletypes "github.com/zeta-chain/node/precompiles/types"
)

func TestPrecompilesDistributeAndClaimThroughContract(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	var (
		// Addresses.
		staker               = r.EVMAddress()
		distrContractAddress = staking.ContractAddress
		lockerAddress        = bank.ContractAddress

		// Stake amount.
		stakeAmt = new(big.Int)

		// ZRC20 distribution.
		zrc20Address  = r.ERC20ZRC20Addr
		zrc20Denom    = precompiletypes.ZRC20ToCosmosDenom(zrc20Address)
		zrc20DistrAmt = big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(1e6))

		// carry is carried from the TestPrecompilesDistributeName test. It's applicable only to locker address.
		// This is needed because there's no easy way to retrieve that balance from the locker.
		carry              = big.NewInt(6210810988040846448)
		zrc20DistrAmtCarry = new(big.Int).Add(zrc20DistrAmt, carry)
		oneThousand        = big.NewInt(1e3)
		oneThousandOne     = big.NewInt(1001)
		fiveHundred        = big.NewInt(500)
		fiveHundredOne     = big.NewInt(501)
		zero               = big.NewInt(0)

		previousGasLimit = r.ZEVMAuth.GasLimit
	)

	// stakeAmt has to be as big as the validator self delegation.
	// This way the rewards will be distributed 50%.
	_, ok := stakeAmt.SetString("1000000000000000000000", 10)
	require.True(r, ok)

	// Set new gas limit to avoid out of gas errors.
	r.ZEVMAuth.GasLimit = 10_000_000

	distrContract, err := staking.NewIStaking(distrContractAddress, r.ZEVMClient)
	require.NoError(r, err, "failed to create distribute contract caller")

	// testDstrContract  is the dApp contract that uses the staking precompile under the hood.
	_, tx, testDstrContract, err := testdistribute.DeployTestDistribute(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "deployment of disitributor caller contract failed")

	// Save first validator bech32 address and ValAddress as it will be used through the test.
	validatorAddr, validatorValAddr := getValidatorAddresses(r, distrContract)

	// Reset the test after it finishes.
	defer resetDistributionTest(r, distrContract, lockerAddress, previousGasLimit, staker, validatorValAddr)

	// Get ERC20ZRC20.
	txHash := r.DepositERC20WithAmountAndMessage(staker, zrc20DistrAmt, []byte{})
	utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)

	// There is no delegation, so the response should be empty.
	dv, err := testDstrContract.GetDelegatorValidatorsThroughContract(
		&bind.CallOpts{},
		staker,
	)
	require.NoError(r, err)
	require.Empty(r, dv, "DelegatorValidators response should be empty")

	// There should be no rewards.
	rewards, err := testDstrContract.GetRewardsThroughContract(&bind.CallOpts{}, staker, validatorAddr)
	require.NoError(r, err)
	require.Empty(r, rewards, "rewards should be empty when there are no delegations")

	// Stake with spender so it's registered as a delegator.
	err = stakeThroughCosmosAPI(r, validatorValAddr, staker, stakeAmt)
	require.NoError(r, err)

	// Check initial balances.
	balanceShouldBe(r, zrc20DistrAmt, checkZRC20Balance(r, staker))
	balanceShouldBe(r, carry, checkZRC20Balance(r, lockerAddress))
	balanceShouldBe(r, zero, checkCosmosBalance(r, r.FeeCollectorAddress, zrc20Denom))

	tx, err = testDstrContract.DistributeThroughContract(r.ZEVMAuth, zrc20Address, oneThousand)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt, "distribute should fail when there's no allowance")

	// Balances shouldn't change after a failed attempt.
	balanceShouldBe(r, zrc20DistrAmt, checkZRC20Balance(r, staker))
	balanceShouldBe(r, carry, checkZRC20Balance(r, lockerAddress))
	balanceShouldBe(r, zero, checkCosmosBalance(r, r.FeeCollectorAddress, zrc20Denom))

	// Allow 500.
	approveAllowance(r, distrContractAddress, fiveHundred)

	tx, err = testDstrContract.DistributeThroughContract(r.ZEVMAuth, zrc20Address, fiveHundredOne)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt, "distribute should fail trying to distribute more than allowed")

	// Balances shouldn't change after a failed attempt.
	balanceShouldBe(r, zrc20DistrAmt, checkZRC20Balance(r, staker))
	balanceShouldBe(r, carry, checkZRC20Balance(r, lockerAddress))
	balanceShouldBe(r, zero, checkCosmosBalance(r, r.FeeCollectorAddress, zrc20Denom))

	// Raise the allowance to 1000.
	approveAllowance(r, distrContractAddress, oneThousand)

	// Shouldn't be able to distribute more than owned balance.
	tx, err = testDstrContract.DistributeThroughContract(r.ZEVMAuth, zrc20Address, oneThousandOne)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt, "distribute should fail trying to distribute more than owned balance")

	// Balances shouldn't change after a failed attempt.
	balanceShouldBe(r, zrc20DistrAmt, checkZRC20Balance(r, staker))
	balanceShouldBe(r, carry, checkZRC20Balance(r, lockerAddress))
	balanceShouldBe(r, zero, checkCosmosBalance(r, r.FeeCollectorAddress, zrc20Denom))

	// Raise the allowance to max tokens.
	approveAllowance(r, distrContractAddress, zrc20DistrAmt)

	// Should be able to distribute an amount which is within balance and allowance.
	tx, err = testDstrContract.DistributeThroughContract(r.ZEVMAuth, zrc20Address, zrc20DistrAmt)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "distribute should succeed when distributing within balance and allowance")

	balanceShouldBe(r, zero, checkZRC20Balance(r, staker))
	balanceShouldBe(r, zrc20DistrAmtCarry, checkZRC20Balance(r, lockerAddress))
	balanceShouldBe(r, zrc20DistrAmt, checkCosmosBalance(r, r.FeeCollectorAddress, zrc20Denom))

	eventDitributed, err := distrContract.ParseDistributed(*receipt.Logs[0])
	require.NoError(r, err)
	require.Equal(r, zrc20Address, eventDitributed.Zrc20Token)
	require.Equal(r, staker, eventDitributed.Zrc20Distributor)
	require.Equal(r, zrc20DistrAmt.Uint64(), eventDitributed.Amount.Uint64())

	// After one block the rewards should have been distributed and fee collector should have 0 ZRC20 balance.
	r.WaitForBlocks(1)
	balanceShouldBe(r, zero, checkCosmosBalance(r, r.FeeCollectorAddress, zrc20Denom))

	// DelegatorValidators returns the list of validator this delegator has delegated to.
	// The result should include the validator address.
	dv, err = testDstrContract.GetDelegatorValidatorsThroughContract(&bind.CallOpts{}, staker)
	require.NoError(r, err)
	require.Contains(r, dv, validatorAddr, "DelegatorValidators response should include validator address")

	// Get rewards and check it contains zrc20 tokens.
	rewards, err = testDstrContract.GetRewardsThroughContract(&bind.CallOpts{}, staker, validatorAddr)
	require.NoError(r, err)
	require.GreaterOrEqual(r, len(rewards), 2)
	found := false
	for _, coin := range rewards {
		if strings.Contains(coin.Denom, config.ZRC20DenomPrefix) {
			found = true
			break
		}
	}
	require.True(r, found, "rewards should include the ZRC20 token")

	tx, err = testDstrContract.ClaimRewardsThroughContract(r.ZEVMAuth, staker, validatorAddr)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "claim rewards should succeed")

	// Before claiming rewards the ZRC20 balance is 0. After claiming rewards the ZRC20 balance should be 14239697290875601808.
	// Which is the amount of ZRC20 distributed, divided by two validators, and subtracted the commissions.
	zrc20RewardsAmt, ok := big.NewInt(0).SetString("14239697290875601808", 10)
	require.True(r, ok)
	balanceShouldBe(r, zrc20RewardsAmt, checkZRC20Balance(r, staker))

	eventClaimed, err := distrContract.ParseClaimedRewards(*receipt.Logs[0])
	require.NoError(r, err)
	require.Equal(r, zrc20Address, eventClaimed.Zrc20Token)
	require.Equal(r, staker, eventClaimed.ClaimAddress)
	require.Equal(r, common.BytesToAddress(validatorValAddr.Bytes()), eventClaimed.Validator)
	require.Equal(r, zrc20RewardsAmt.Uint64(), eventClaimed.Amount.Uint64())

	// Locker final balance should be zrc20Distributed with carry - zrc20RewardsAmt.
	lockerFinalBalance := big.NewInt(0).Sub(zrc20DistrAmtCarry, zrc20RewardsAmt)
	balanceShouldBe(r, lockerFinalBalance, checkZRC20Balance(r, lockerAddress))

	// Staker final cosmos balance should be 0.
	balanceShouldBe(r, zero, checkCosmosBalance(r, sdk.AccAddress(staker.Bytes()), zrc20Denom))
}
