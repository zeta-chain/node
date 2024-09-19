package e2etests

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/precompiles/staking"
)

func TestPrecompilesStaking(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	stakingContract, err := staking.NewIStaking(staking.ContractAddress, r.ZEVMClient)
	require.NoError(r, err, "Failed to create staking contract caller")

	previousGasLimit := r.ZEVMAuth.GasLimit
	r.ZEVMAuth.GasLimit = 10000000
	defer func() {
		r.ZEVMAuth.GasLimit = previousGasLimit
	}()

	validators, err := stakingContract.GetAllValidators(&bind.CallOpts{})
	require.NoError(r, err)
	require.GreaterOrEqual(r, len(validators), 2)

	CleanValidatorDelegations(r, stakingContract, validators)

	// shares are 0 for both validators at the start
	sharesBeforeVal1, err := stakingContract.GetShares(&bind.CallOpts{}, r.ZEVMAuth.From, validators[0].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, int64(0), sharesBeforeVal1.Int64())

	sharesBeforeVal2, err := stakingContract.GetShares(&bind.CallOpts{}, r.ZEVMAuth.From, validators[1].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, int64(0), sharesBeforeVal2.Int64())

	// stake 3 to validator1
	tx, err := stakingContract.Stake(r.ZEVMAuth, r.ZEVMAuth.From, validators[0].OperatorAddress, big.NewInt(3))
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// check shares are set to 3
	sharesAfterVal1, err := stakingContract.GetShares(&bind.CallOpts{}, r.ZEVMAuth.From, validators[0].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, big.NewInt(3e18).String(), sharesAfterVal1.String())

	// check delegation amount using staking keeper query client
	delegationAfterVal1, err := r.StakingClient.Delegation(r.Ctx, &types.QueryDelegationRequest{
		DelegatorAddr: sdk.AccAddress(r.ZEVMAuth.From.Bytes()).String(),
		ValidatorAddr: validators[0].OperatorAddress,
	})
	require.NoError(r, err)
	require.Equal(r, int64(3), delegationAfterVal1.DelegationResponse.Balance.Amount.Int64())

	// unstake 1 from validator1
	tx, err = stakingContract.Unstake(r.ZEVMAuth, r.ZEVMAuth.From, validators[0].OperatorAddress, big.NewInt(1))
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// check shares are set to 2
	sharesAfterVal1, err = stakingContract.GetShares(&bind.CallOpts{}, r.ZEVMAuth.From, validators[0].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, big.NewInt(2e18).String(), sharesAfterVal1.String())

	// check delegation amount using staking keeper query client
	delegationAfterVal1, err = r.StakingClient.Delegation(r.Ctx, &types.QueryDelegationRequest{
		DelegatorAddr: sdk.AccAddress(r.ZEVMAuth.From.Bytes()).String(),
		ValidatorAddr: validators[0].OperatorAddress,
	})
	require.NoError(r, err)
	require.Equal(r, int64(2), delegationAfterVal1.DelegationResponse.Balance.Amount.Int64())

	// move 1 stake from validator1 to validator2
	tx, err = stakingContract.MoveStake(
		r.ZEVMAuth,
		r.ZEVMAuth.From,
		validators[0].OperatorAddress,
		validators[1].OperatorAddress,
		big.NewInt(1),
	)
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// check shares for both validator1 and validator2 are 1
	sharesAfterVal1, err = stakingContract.GetShares(&bind.CallOpts{}, r.ZEVMAuth.From, validators[0].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, big.NewInt(1e18).String(), sharesAfterVal1.String())

	sharesAfterVal2, err := stakingContract.GetShares(&bind.CallOpts{}, r.ZEVMAuth.From, validators[1].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, big.NewInt(1e18).String(), sharesAfterVal2.String())

	// check delegation amount using staking keeper query client
	delegationAfterVal1, err = r.StakingClient.Delegation(r.Ctx, &types.QueryDelegationRequest{
		DelegatorAddr: sdk.AccAddress(r.ZEVMAuth.From.Bytes()).String(),
		ValidatorAddr: validators[0].OperatorAddress,
	})
	require.NoError(r, err)
	require.Equal(r, int64(1), delegationAfterVal1.DelegationResponse.Balance.Amount.Int64())

	delegationAfterVal2, err := r.StakingClient.Delegation(r.Ctx, &types.QueryDelegationRequest{
		DelegatorAddr: sdk.AccAddress(r.ZEVMAuth.From.Bytes()).String(),
		ValidatorAddr: validators[1].OperatorAddress,
	})
	require.NoError(r, err)
	require.Equal(r, int64(1), delegationAfterVal2.DelegationResponse.Balance.Amount.Int64())
}

// CleanValidatorDelegations unstakes all delegations from the given validators if delegations ar present
func CleanValidatorDelegations(r *runner.E2ERunner, stakingContract *staking.IStaking, validators []staking.Validator) {
	for _, validator := range validators {
		delegator := sdk.AccAddress(r.ZEVMAuth.From.Bytes()).String()
		delegation, err := r.StakingClient.Delegation(r.Ctx, &types.QueryDelegationRequest{
			DelegatorAddr: delegator,
			ValidatorAddr: validator.OperatorAddress,
		})
		if err != nil || delegation.DelegationResponse == nil {
			continue
		}

		delegationAmount := delegation.DelegationResponse.Balance.Amount.Int64()
		if delegationAmount > 0 && err == nil {
			tx, err := stakingContract.Unstake(
				r.ZEVMAuth,
				r.ZEVMAuth.From,
				validator.OperatorAddress,
				big.NewInt(delegationAmount),
			)
			require.NoError(r, err)
			utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		}
	}
}
