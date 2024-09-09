package keeper_test

import (
	"math/rand"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestKeeper_NotImplementedHooks(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)

	hooks := k.Hooks()
	require.Nil(t, hooks.AfterValidatorCreated(ctx, nil))
	require.Nil(t, hooks.BeforeValidatorModified(ctx, nil))
	require.Nil(t, hooks.AfterValidatorBonded(ctx, nil, nil))
	require.Nil(t, hooks.BeforeDelegationCreated(ctx, nil, nil))
	require.Nil(t, hooks.BeforeDelegationSharesModified(ctx, nil, nil))
	require.Nil(t, hooks.BeforeDelegationRemoved(ctx, nil, nil))
}

func TestKeeper_AfterValidatorRemoved(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)

	// #nosec G404 test purpose - weak randomness is not an issue here
	r := rand.New(rand.NewSource(1))
	valAddr := sample.ValAddress(r)
	accAddress, err := types.GetAccAddressFromOperatorAddress(valAddr.String())
	require.NoError(t, err)
	os := types.ObserverSet{
		ObserverList: []string{accAddress.String()},
	}
	k.SetObserverSet(ctx, os)

	hooks := k.Hooks()
	err = hooks.AfterValidatorRemoved(ctx, nil, valAddr)
	require.NoError(t, err)

	os, found := k.GetObserverSet(ctx)
	require.True(t, found)
	// observer for validator is removed from set
	require.Empty(t, os.ObserverList)
}

func TestKeeper_AfterValidatorBeginUnbonding(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.ObserverKeeper(t)

	r := rand.New(rand.NewSource(9))
	validator := sample.Validator(t, r)
	validator.DelegatorShares = sdk.NewDec(100)
	sdkk.StakingKeeper.SetValidator(ctx, validator)
	accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
	require.NoError(t, err)

	sdkk.StakingKeeper.SetDelegation(ctx, stakingtypes.Delegation{
		DelegatorAddress: accAddressOfValidator.String(),
		ValidatorAddress: validator.GetOperator().String(),
		Shares:           sdk.NewDec(10),
	})

	k.SetObserverSet(ctx, types.ObserverSet{
		ObserverList: []string{accAddressOfValidator.String()},
	})

	hooks := k.Hooks()
	err = hooks.AfterValidatorBeginUnbonding(ctx, nil, validator.GetOperator())
	require.NoError(t, err)

	os, found := k.GetObserverSet(ctx)
	require.True(t, found)
	require.Empty(t, os.ObserverList)
}

func TestKeeper_AfterDelegationModified(t *testing.T) {
	t.Run("should not clean observer if not self delegation", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.ObserverKeeper(t)

		r := rand.New(rand.NewSource(9))
		validator := sample.Validator(t, r)
		validator.DelegatorShares = sdk.NewDec(100)
		sdkk.StakingKeeper.SetValidator(ctx, validator)
		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)

		sdkk.StakingKeeper.SetDelegation(ctx, stakingtypes.Delegation{
			DelegatorAddress: accAddressOfValidator.String(),
			ValidatorAddress: validator.GetOperator().String(),
			Shares:           sdk.NewDec(10),
		})

		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String()},
		})

		hooks := k.Hooks()
		err = hooks.AfterDelegationModified(ctx, sdk.AccAddress(sample.AccAddress()), validator.GetOperator())
		require.NoError(t, err)

		os, found := k.GetObserverSet(ctx)
		require.True(t, found)
		require.Equal(t, 1, len(os.ObserverList))
		require.Equal(t, accAddressOfValidator.String(), os.ObserverList[0])
	})

	t.Run("should clean observer if self delegation", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.ObserverKeeper(t)

		r := rand.New(rand.NewSource(9))
		validator := sample.Validator(t, r)
		validator.DelegatorShares = sdk.NewDec(100)
		sdkk.StakingKeeper.SetValidator(ctx, validator)
		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)

		sdkk.StakingKeeper.SetDelegation(ctx, stakingtypes.Delegation{
			DelegatorAddress: accAddressOfValidator.String(),
			ValidatorAddress: validator.GetOperator().String(),
			Shares:           sdk.NewDec(10),
		})

		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String()},
		})

		hooks := k.Hooks()
		err = hooks.AfterDelegationModified(ctx, accAddressOfValidator, validator.GetOperator())
		require.NoError(t, err)

		os, found := k.GetObserverSet(ctx)
		require.True(t, found)
		require.Empty(t, os.ObserverList)
	})
}

func TestKeeper_BeforeValidatorSlashed(t *testing.T) {
	t.Run("should error if validator not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		r := rand.New(rand.NewSource(9))
		validator := sample.Validator(t, r)

		hooks := k.Hooks()
		err := hooks.BeforeValidatorSlashed(ctx, validator.GetOperator(), sdk.NewDec(1))
		require.Error(t, err)
	})

	t.Run("should not error if observer set not found", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.ObserverKeeper(t)

		r := rand.New(rand.NewSource(9))
		validator := sample.Validator(t, r)
		validator.DelegatorShares = sdk.NewDec(100)
		sdkk.StakingKeeper.SetValidator(ctx, validator)

		hooks := k.Hooks()
		err := hooks.BeforeValidatorSlashed(ctx, validator.GetOperator(), sdk.NewDec(1))
		require.NoError(t, err)
	})

	t.Run("should remove from observer set", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.ObserverKeeper(t)

		r := rand.New(rand.NewSource(9))
		validator := sample.Validator(t, r)
		validator.DelegatorShares = sdk.NewDec(100)
		validator.Tokens = sdk.NewInt(100)
		sdkk.StakingKeeper.SetValidator(ctx, validator)
		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)

		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String()},
		})

		hooks := k.Hooks()
		err = hooks.BeforeValidatorSlashed(ctx, validator.GetOperator(), sdk.MustNewDecFromStr("0.1"))
		require.NoError(t, err)

		os, found := k.GetObserverSet(ctx)
		require.True(t, found)
		require.Empty(t, os.ObserverList)
	})

	t.Run("should not remove from observer set", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.ObserverKeeper(t)

		r := rand.New(rand.NewSource(9))
		validator := sample.Validator(t, r)
		validator.DelegatorShares = sdk.NewDec(100)
		validator.Tokens = sdk.NewInt(1000000000000000000)
		sdkk.StakingKeeper.SetValidator(ctx, validator)
		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)

		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String()},
		})

		hooks := k.Hooks()
		err = hooks.BeforeValidatorSlashed(ctx, validator.GetOperator(), sdk.MustNewDecFromStr("0"))
		require.NoError(t, err)

		os, found := k.GetObserverSet(ctx)
		require.True(t, found)
		require.Equal(t, 1, len(os.ObserverList))
		require.Equal(t, accAddressOfValidator.String(), os.ObserverList[0])
	})
}
