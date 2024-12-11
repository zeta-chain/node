package e2etests

import (
	"math/big"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/precompiles/bank"
	"github.com/zeta-chain/node/precompiles/staking"
	precompiletypes "github.com/zeta-chain/node/precompiles/types"
)

func TestPrecompilesDistributeAndClaim(r *runner.E2ERunner, args []string) {
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

		// Amounts to test with.
		higherThanBalance = big.NewInt(0).Add(zrc20DistrAmt, big.NewInt(1))
		fiveHundred       = big.NewInt(500)
		fiveHundredOne    = big.NewInt(501)
		zero              = big.NewInt(0)
		stake             = "1000000000000000000000"

		previousGasLimit = r.ZEVMAuth.GasLimit
	)

	// stakeAmt has to be as big as the validator self delegation.
	// This way the rewards will be distributed 50%.
	_, ok := stakeAmt.SetString(stake, 10)
	require.True(r, ok)

	// Set new gas limit to avoid out of gas errors.
	r.ZEVMAuth.GasLimit = 10_000_000

	distrContract, err := staking.NewIStaking(distrContractAddress, r.ZEVMClient)
	require.NoError(r, err, "failed to create distribute contract caller")

	// Retrieve the list of validators.
	validators, err := distrContract.GetAllValidators(&bind.CallOpts{})
	require.NoError(r, err)
	require.GreaterOrEqual(r, len(validators), 2)

	// Save first validator bech32 address and as it will be used through the test.
	validatorAddr, validatorValAddr := getValidatorAddresses(r, distrContract)

	// Reset the test after it finishes.
	defer resetDistributionTest(r, distrContract, lockerAddress, previousGasLimit, staker, validatorValAddr)

	// Get ERC20ZRC20.
	txHash := r.LegacyDepositERC20WithAmountAndMessage(staker, zrc20DistrAmt, []byte{})
	utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)

	// There is no delegation, so the response should be empty.
	dv, err := distrContract.GetDelegatorValidators(&bind.CallOpts{}, staker)
	require.NoError(r, err)
	require.Empty(r, dv, "DelegatorValidators response should be empty")

	// Shares at this point should be 0.
	sharesBefore, err := distrContract.GetShares(&bind.CallOpts{}, r.ZEVMAuth.From, validatorAddr)
	require.NoError(r, err)
	require.Equal(r, int64(0), sharesBefore.Int64(), "shares should be 0 when there are no delegations")

	// There should be no rewards.
	rewards, err := distrContract.GetRewards(&bind.CallOpts{}, staker, validatorAddr)
	require.NoError(r, err)
	require.Empty(r, rewards, "rewards should be empty when there are no delegations")

	// Stake with spender so it's registered as a delegator.
	err = stakeThroughCosmosAPI(r, validatorValAddr, staker, stakeAmt)
	require.NoError(r, err)

	// Check initial balances.
	balanceShouldBe(r, zrc20DistrAmt, checkZRC20Balance(r, staker))
	balanceShouldBe(r, zero, checkZRC20Balance(r, lockerAddress))
	balanceShouldBe(r, zero, checkCosmosBalance(r, r.FeeCollectorAddress, zrc20Denom))

	// Failed attempt!
	tx, err := distrContract.Distribute(r.ZEVMAuth, zrc20Address, zrc20DistrAmt)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt, "distribute should fail when there's no allowance")

	// Balances shouldn't change after a failed attempt.
	balanceShouldBe(r, zrc20DistrAmt, checkZRC20Balance(r, staker))
	balanceShouldBe(r, zero, checkZRC20Balance(r, lockerAddress))
	balanceShouldBe(r, zero, checkCosmosBalance(r, r.FeeCollectorAddress, zrc20Denom))

	// Allow 500.
	approveAllowance(r, distrContractAddress, fiveHundred)

	// Failed attempt! Shouldn't be able to distribute more than allowed.
	tx, err = distrContract.Distribute(r.ZEVMAuth, zrc20Address, fiveHundredOne)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt, "distribute should fail trying to distribute more than allowed")

	// Balances shouldn't change after a failed attempt.
	balanceShouldBe(r, zrc20DistrAmt, checkZRC20Balance(r, staker))
	balanceShouldBe(r, zero, checkZRC20Balance(r, lockerAddress))
	balanceShouldBe(r, zero, checkCosmosBalance(r, r.FeeCollectorAddress, zrc20Denom))

	// Raise the allowance to the maximum ZRC20 amount.
	approveAllowance(r, distrContractAddress, zrc20DistrAmt)

	// Failed attempt! Shouldn't be able to distribute more than owned balance.
	tx, err = distrContract.Distribute(r.ZEVMAuth, zrc20Address, higherThanBalance)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt, "distribute should fail trying to distribute more than owned balance")

	// Balances shouldn't change after a failed attempt.
	balanceShouldBe(r, zrc20DistrAmt, checkZRC20Balance(r, staker))
	balanceShouldBe(r, zero, checkZRC20Balance(r, lockerAddress))
	balanceShouldBe(r, zero, checkCosmosBalance(r, r.FeeCollectorAddress, zrc20Denom))

	// Should be able to distribute an amount which is within balance and allowance.
	tx, err = distrContract.Distribute(r.ZEVMAuth, zrc20Address, zrc20DistrAmt)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "distribute should succeed when distributing within balance and allowance")

	balanceShouldBe(r, zero, checkZRC20Balance(r, staker))
	balanceShouldBe(r, zrc20DistrAmt, checkZRC20Balance(r, lockerAddress))
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
	dv, err = distrContract.GetDelegatorValidators(&bind.CallOpts{}, staker)
	require.NoError(r, err)
	require.Contains(r, dv, validatorAddr, "DelegatorValidators response should include validator address")

	// Get rewards and check it contains zrc20 tokens.
	rewards, err = distrContract.GetRewards(&bind.CallOpts{}, staker, validatorAddr)
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

	// Claim the rewards, they'll be unlocked as ZRC20 tokens.
	tx, err = distrContract.ClaimRewards(r.ZEVMAuth, staker, validatorAddr)
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

	// Locker final balance should be zrc20Disitributed - zrc20RewardsAmt.
	lockerFinalBalance := big.NewInt(0).Sub(zrc20DistrAmt, zrc20RewardsAmt)
	balanceShouldBe(r, lockerFinalBalance, checkZRC20Balance(r, lockerAddress))

	// Staker final cosmos balance should be 0.
	balanceShouldBe(r, zero, checkCosmosBalance(r, sdk.AccAddress(staker.Bytes()), zrc20Denom))
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
	r.LegacyDepositAndApproveWZeta(approveAmount)

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
func checkCosmosBalance(r *runner.E2ERunner, address sdk.AccAddress, denom string) *big.Int {
	bal, err := r.BankClient.Balance(
		r.Ctx,
		&banktypes.QueryBalanceRequest{Address: address.String(), Denom: denom},
	)
	require.NoError(r, err)

	return bal.Balance.Amount.BigInt()
}

func stakeThroughCosmosAPI(
	r *runner.E2ERunner,
	validator sdk.ValAddress,
	staker common.Address,
	amount *big.Int,
) error {
	msg := stakingtypes.NewMsgDelegate(
		sdk.AccAddress(staker.Bytes()),
		validator,
		sdk.Coin{
			Denom:  config.BaseDenom,
			Amount: math.NewIntFromBigInt(amount),
		},
	)

	_, err := r.ZetaTxServer.BroadcastTx(sdk.AccAddress(staker.Bytes()).String(), msg)
	if err != nil {
		return err
	}

	return nil
}

func resetDistributionTest(
	r *runner.E2ERunner,
	distrContract *staking.IStaking,
	lockerAddress common.Address,
	previousGasLimit uint64,
	staker common.Address,
	validator sdk.ValAddress,
) {
	validatorAddr, _ := getValidatorAddresses(r, distrContract)

	amount, err := distrContract.GetShares(&bind.CallOpts{}, r.ZEVMAuth.From, validatorAddr)
	require.NoError(r, err)

	// Restore the gas limit.
	r.ZEVMAuth.GasLimit = previousGasLimit

	// Reset the allowance to 0; this is needed when running upgrade tests where this test runs twice.
	tx, err := r.ERC20ZRC20.Approve(r.ZEVMAuth, lockerAddress, big.NewInt(0))
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "resetting allowance failed")

	// Reset balance to 0 for spender; this is needed when running upgrade tests where this test runs twice.
	balance, err := r.ERC20ZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)

	// Burn all ERC20 balance.
	tx, err = r.ERC20ZRC20.Transfer(
		r.ZEVMAuth,
		common.HexToAddress("0x000000000000000000000000000000000000dEaD"),
		balance,
	)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "Resetting balance failed")

	// Clean the delegation.
	// Delegator will always delegate on the first validator.
	msg := stakingtypes.NewMsgUndelegate(
		sdk.AccAddress(staker.Bytes()),
		validator,
		sdk.Coin{
			Denom:  config.BaseDenom,
			Amount: math.NewIntFromBigInt(amount.Div(amount, big.NewInt(1e18))),
		},
	)

	_, err = r.ZetaTxServer.BroadcastTx(sdk.AccAddress(staker.Bytes()).String(), msg)
	require.NoError(r, err)
}

func getValidatorAddresses(r *runner.E2ERunner, distrContract *staking.IStaking) (string, sdk.ValAddress) {
	// distrContract, err := staking.NewIStaking(staking.ContractAddress, r.ZEVMClient)
	// require.NoError(r, err, "failed to create distribute contract caller")

	// Retrieve the list of validators.
	validators, err := distrContract.GetAllValidators(&bind.CallOpts{})
	require.NoError(r, err)
	require.GreaterOrEqual(r, len(validators), 2)

	// Save first validators as it will be used through the test.
	validatorAddr, err := sdk.ValAddressFromBech32(validators[0].OperatorAddress)
	require.NoError(r, err)

	return validators[0].OperatorAddress, validatorAddr
}
