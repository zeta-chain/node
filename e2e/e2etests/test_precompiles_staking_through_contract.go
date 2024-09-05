package e2etests

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/contracts/teststaking"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

func TestPrecompilesStakingThroughContract(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	testStakingAddr, testStakingTx, testStaking, err := teststaking.DeployTestStaking(r.ZEVMAuth, r.ZEVMClient, r.WZetaAddr)
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, testStakingTx, r.Logger, r.ReceiptTimeout)

	validators, err := testStaking.GetAllValidators(nil)
	require.NoError(r, err)
	require.GreaterOrEqual(r, len(validators), 2)

	// shares are 0 for both validators at the start
	sharesBeforeVal1, err := testStaking.GetShares(nil, testStakingAddr, validators[0].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, int64(0), sharesBeforeVal1.Int64())

	sharesBeforeVal2, err := testStaking.GetShares(nil, testStakingAddr, validators[1].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, int64(0), sharesBeforeVal2.Int64())

	// stake 3 to validator1 using user_precompile account should fail because sender is not provided staker
	_, err = testStaking.Stake(r.ZEVMAuth, r.ZEVMAuth.From, validators[0].OperatorAddress, big.NewInt(3))
	require.Error(r, err)

	// stake 3 to validator1 using testStaking smart contract should fail because it doesn't have any balance yet
	_, err = testStaking.Stake(r.ZEVMAuth, testStakingAddr, validators[0].OperatorAddress, big.NewInt(3))
	require.Error(r, err)

	r.ZEVMAuth.Value = big.NewInt(1000000000000)
	r.ZEVMAuth.GasLimit = 10000000

	// fund testStaking contract with azeta
	tx, err := testStaking.DepositWZETA(r.ZEVMAuth)
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.ZEVMAuth.Value = big.NewInt(0)

	tx, err = testStaking.WithdrawWZETA(r.ZEVMAuth, big.NewInt(100000000000))
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// stake 3 to validator1 using testStaking smart contract
	tx, err = testStaking.Stake(r.ZEVMAuth, testStakingAddr, validators[0].OperatorAddress, big.NewInt(3))
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// check shares are set to 3
	sharesAfterVal1, err := testStaking.GetShares(nil, testStakingAddr, validators[0].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, big.NewInt(3e18).String(), sharesAfterVal1.String())

	// check delegation amount using staking keeper query client
	delegationAfterVal1, err := r.StakingClient.Delegation(r.Ctx, &types.QueryDelegationRequest{
		DelegatorAddr: sdk.AccAddress(testStakingAddr.Bytes()).String(),
		ValidatorAddr: validators[0].OperatorAddress,
	})
	require.NoError(r, err)
	require.Equal(r, int64(3), delegationAfterVal1.DelegationResponse.Balance.Amount.Int64())

	// unstake 1 from validator1
	tx, err = testStaking.Unstake(r.ZEVMAuth, testStakingAddr, validators[0].OperatorAddress, big.NewInt(1))
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// check shares are set to 2
	sharesAfterVal1, err = testStaking.GetShares(nil, testStakingAddr, validators[0].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, big.NewInt(2e18).String(), sharesAfterVal1.String())

	// check delegation amount using staking keeper query client
	delegationAfterVal1, err = r.StakingClient.Delegation(r.Ctx, &types.QueryDelegationRequest{
		DelegatorAddr: sdk.AccAddress(testStakingAddr.Bytes()).String(),
		ValidatorAddr: validators[0].OperatorAddress,
	})
	require.NoError(r, err)
	require.Equal(r, int64(2), delegationAfterVal1.DelegationResponse.Balance.Amount.Int64())

	// move 1 stake from validator1 to validator2
	tx, err = testStaking.MoveStake(
		r.ZEVMAuth,
		testStakingAddr,
		validators[0].OperatorAddress,
		validators[1].OperatorAddress,
		big.NewInt(1),
	)
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// check shares for both validator1 and validator2 are 1
	sharesAfterVal1, err = testStaking.GetShares(nil, testStakingAddr, validators[0].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, big.NewInt(1e18).String(), sharesAfterVal1.String())

	// check delegation amount using staking keeper query client
	delegationAfterVal1, err = r.StakingClient.Delegation(r.Ctx, &types.QueryDelegationRequest{
		DelegatorAddr: sdk.AccAddress(testStakingAddr.Bytes()).String(),
		ValidatorAddr: validators[0].OperatorAddress,
	})
	require.NoError(r, err)
	require.Equal(r, int64(1), delegationAfterVal1.DelegationResponse.Balance.Amount.Int64())

	sharesAfterVal2, err := testStaking.GetShares(nil, testStakingAddr, validators[1].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, big.NewInt(1e18).String(), sharesAfterVal2.String())

	// check delegation amount using staking keeper query client
	delegationAfterVal2, err := r.StakingClient.Delegation(r.Ctx, &types.QueryDelegationRequest{
		DelegatorAddr: sdk.AccAddress(testStakingAddr.Bytes()).String(),
		ValidatorAddr: validators[1].OperatorAddress,
	})
	require.NoError(r, err)
	require.Equal(r, int64(1), delegationAfterVal2.DelegationResponse.Balance.Amount.Int64())
}
