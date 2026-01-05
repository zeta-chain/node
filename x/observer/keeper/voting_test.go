package keeper_test

import (
	"math/rand"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

// setSupportedChain sets the supported chains for the observer module
func setSupportedChain(ctx sdk.Context, observerKeeper keeper.Keeper, chainIDs ...int64) {
	chainParamsList := make([]*types.ChainParams, len(chainIDs))
	for i, chainID := range chainIDs {
		chainParams := sample.ChainParams(chainID)
		chainParams.IsSupported = true
		chainParamsList[i] = chainParams
	}
	observerKeeper.SetChainParamsList(ctx, types.ChainParamsList{
		ChainParams: chainParamsList,
	})
}

// getValidEthChainIDWithIndex get a valid eth chain id with index
func getValidEthChainIDWithIndex(t *testing.T, index int) int64 {
	switch index {
	case 0:
		return chains.GoerliLocalnet.ChainId
	case 1:
		return chains.Goerli.ChainId
	default:
		require.Fail(t, "invalid index")
	}
	return 0
}

func TestKeeper_IsAuthorized(t *testing.T) {
	t.Run("authorized observer", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.ObserverKeeper(t)

		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		require.NoError(t, sdkk.StakingKeeper.SetValidator(ctx, validator))
		consAddress, err := validator.GetConsAddr()
		require.NoError(t, err)
		require.NoError(
			t,
			k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
				Address:             string(consAddress),
				StartHeight:         0,
				JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
				Tombstoned:          false,
				MissedBlocksCounter: 1,
			}),
		)

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)

		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String()},
		})
		require.NoError(t, k.CheckObserverCanVote(ctx, accAddressOfValidator.String()))
	})

	t.Run("not authorized if the observer address is invalid", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		invalidAddress := "invalid_address"

		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{invalidAddress},
		})
		require.ErrorIs(t, k.CheckObserverCanVote(ctx, invalidAddress), types.ErrInvalidAddress)
	})

	t.Run("not authorized for jailed observer", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.ObserverKeeper(t)

		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		validator.Jailed = true
		require.NoError(t, sdkk.StakingKeeper.SetValidator(ctx, validator))

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)

		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String()},
		})

		err = k.CheckObserverCanVote(ctx, accAddressOfValidator.String())
		require.ErrorIs(t, err, types.ErrValidatorJailed)
	})

	t.Run("not authorized for non bonded observers", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.ObserverKeeper(t)

		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		validator.Status = stakingtypes.Unbonded
		require.NoError(t, sdkk.StakingKeeper.SetValidator(ctx, validator))

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)

		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String()},
		})

		err = k.CheckObserverCanVote(ctx, accAddressOfValidator.String())
		require.ErrorIs(t, err, types.ErrValidatorStatus)
	})

	t.Run("not authorized for tombstoned observer", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.ObserverKeeper(t)

		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		require.NoError(t, sdkk.StakingKeeper.SetValidator(ctx, validator))
		consAddress, err := validator.GetConsAddr()
		require.NoError(t, err)
		err = k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             string(consAddress),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          true,
			MissedBlocksCounter: 1,
		})
		require.NoError(t, err)

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)

		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String()},
		})

		err = k.CheckObserverCanVote(ctx, accAddressOfValidator.String())
		require.ErrorIs(t, err, types.ErrValidatorTombstoned)
	})

	t.Run("not authorized for non-validator observer", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)

		consAddress, err := validator.GetConsAddr()
		require.NoError(t, err)
		err = k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             string(consAddress),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          false,
			MissedBlocksCounter: 1,
		})
		require.NoError(t, err)

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)

		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String()},
		})

		err = k.CheckObserverCanVote(ctx, accAddressOfValidator.String())
		require.ErrorIs(t, err, types.ErrNotValidator)
	})

}

func TestKeeper_CheckObserverSelfDelegation(t *testing.T) {
	t.Run("should error if accAddress invalid", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		err := k.CheckObserverSelfDelegation(ctx, "invalid")
		require.Error(t, err)
	})

	t.Run("should error if validator not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		accAddress := sample.AccAddress()
		err := k.CheckObserverSelfDelegation(ctx, accAddress)
		require.Error(t, err)
	})

	t.Run("should error if delegation not found", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.ObserverKeeper(t)

		r := rand.New(rand.NewSource(9))
		validator := sample.Validator(t, r)
		sdkk.StakingKeeper.SetValidator(ctx, validator)
		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)

		err = k.CheckObserverSelfDelegation(ctx, accAddressOfValidator.String())
		require.Error(t, err)
	})

	t.Run("should remove from observer list if tokens less than min del", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.ObserverKeeper(t)

		r := rand.New(rand.NewSource(9))
		validator := sample.Validator(t, r)
		validator.DelegatorShares = sdkmath.LegacyNewDec(100)
		sdkk.StakingKeeper.SetValidator(ctx, validator)
		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)

		sdkk.StakingKeeper.SetDelegation(ctx, stakingtypes.Delegation{
			DelegatorAddress: accAddressOfValidator.String(),
			ValidatorAddress: validator.GetOperator(),
			Shares:           sdkmath.LegacyNewDec(10),
		})

		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String()},
		})
		err = k.CheckObserverSelfDelegation(ctx, accAddressOfValidator.String())
		require.NoError(t, err)

		os, found := k.GetObserverSet(ctx)
		require.True(t, found)
		require.Empty(t, os.ObserverList)
	})

	t.Run("should not remove from observer list if tokens gte than min del", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.ObserverKeeper(t)

		r := rand.New(rand.NewSource(9))
		validator := sample.Validator(t, r)

		validator.DelegatorShares = sdkmath.LegacyNewDec(1)
		validator.Tokens = sdkmath.NewInt(1)
		sdkk.StakingKeeper.SetValidator(ctx, validator)
		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)

		minDelegation, err := types.GetMinObserverDelegationDec()
		require.NoError(t, err)
		sdkk.StakingKeeper.SetDelegation(ctx, stakingtypes.Delegation{
			DelegatorAddress: accAddressOfValidator.String(),
			ValidatorAddress: validator.GetOperator(),
			Shares:           minDelegation,
		})

		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String()},
		})
		err = k.CheckObserverSelfDelegation(ctx, accAddressOfValidator.String())
		require.NoError(t, err)

		os, found := k.GetObserverSet(ctx)
		require.True(t, found)
		require.Equal(t, 1, len(os.ObserverList))
		require.Equal(t, accAddressOfValidator.String(), os.ObserverList[0])
	})
}

func TestKeeper_IsOpeartorTombstoned(t *testing.T) {
	t.Run("should err if invalid addr", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		res, err := k.IsOperatorTombstoned(ctx, "invalid")
		require.Error(t, err)
		require.False(t, res)
	})

	t.Run("should error if validator not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		accAddress := sample.AccAddress()
		res, err := k.IsOperatorTombstoned(ctx, accAddress)
		require.Error(t, err)
		require.False(t, res)
	})

	t.Run("should not error if validator found", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.ObserverKeeper(t)

		r := rand.New(rand.NewSource(9))
		validator := sample.Validator(t, r)
		sdkk.StakingKeeper.SetValidator(ctx, validator)
		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)

		res, err := k.IsOperatorTombstoned(ctx, accAddressOfValidator.String())
		require.NoError(t, err)
		require.False(t, res)
	})
}

func TestKeeper_IsValidator(t *testing.T) {
	t.Run("should err if invalid addr", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		err := k.IsValidator(ctx, "invalid")
		require.Error(t, err)
	})

	t.Run("should error if validator not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		accAddress := sample.AccAddress()
		err := k.IsValidator(ctx, accAddress)
		require.Error(t, err)
	})

	t.Run("should err if validator not bonded", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.ObserverKeeper(t)

		r := rand.New(rand.NewSource(9))
		validator := sample.Validator(t, r)
		validator.Status = stakingtypes.Unbonded
		require.NoError(t, sdkk.StakingKeeper.SetValidator(ctx, validator))
		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)

		err = k.IsValidator(ctx, accAddressOfValidator.String())
		require.Error(t, err)
	})

	t.Run("should err if validator jailed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.ObserverKeeper(t)

		r := rand.New(rand.NewSource(9))
		validator := sample.Validator(t, r)
		validator.Status = stakingtypes.Bonded
		validator.Jailed = true
		sdkk.StakingKeeper.SetValidator(ctx, validator)
		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)

		err = k.IsValidator(ctx, accAddressOfValidator.String())
		require.Error(t, err)
	})

	t.Run("should not err if validator not jailed and bonded", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.ObserverKeeper(t)

		r := rand.New(rand.NewSource(9))
		validator := sample.Validator(t, r)
		validator.Status = stakingtypes.Bonded
		validator.Jailed = false
		sdkk.StakingKeeper.SetValidator(ctx, validator)
		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)

		err = k.IsValidator(ctx, accAddressOfValidator.String())
		require.NoError(t, err)
	})
}

func TestKeeper_FindBallot(t *testing.T) {
	t.Run("should err if chain params not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		_, _, err := k.FindBallot(ctx, "index", chains.Chain{
			ChainId: 987,
		}, types.ObservationType_InboundTx)
		require.Error(t, err)
	})
}

func TestKeeper_VoteOnBallot(t *testing.T) {
	t.Run("fails if chain is not supported", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: false,
				},
			},
		})

		chain, _ := k.GetSupportedChainFromChainID(ctx, 0)
		index := sample.ZetaIndex(t)
		_, _, _, err := k.VoteOnBallot(
			ctx,
			chain,
			index,
			types.ObservationType_InboundTx,
			sample.AccAddress(),
			types.VoteType_SuccessObservation)

		require.ErrorIs(t, err, types.ErrSupportedChains)
	})

	t.Run("fails if the user can not add a vote", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: true,
				},
			},
		})

		chain, _ := k.GetSupportedChainFromChainID(ctx, getValidEthChainIDWithIndex(t, 0))
		index := sample.ZetaIndex(t)
		_, _, _, err := k.VoteOnBallot(
			ctx,
			chain,
			index,
			types.ObservationType_InboundTx,
			sample.AccAddress(),
			types.VoteType_SuccessObservation)

		require.ErrorIs(t, err, types.ErrUnableToAddVote)
	})

	t.Run("user can create a ballot and add a vote", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: true,
				},
			},
		})

		voter := sample.AccAddress()
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{voter},
		})

		chain, _ := k.GetSupportedChainFromChainID(ctx, getValidEthChainIDWithIndex(t, 0))
		index := sample.ZetaIndex(t)
		ballot, isFinalized, isNew, err := k.VoteOnBallot(
			ctx,
			chain,
			index,
			types.ObservationType_InboundTx,
			voter,
			types.VoteType_SuccessObservation)

		require.NoError(t, err)
		require.True(t, isFinalized)
		require.True(t, isNew)
		expectedBallot, found := k.GetBallot(ctx, index)
		require.True(t, found)
		require.Equal(t, expectedBallot, ballot)
	})

	t.Run("user can create a ballot and add a vote, without finalizing ballot", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		threshold, err := sdkmath.LegacyNewDecFromStr("0.7")
		require.NoError(t, err)

		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:         getValidEthChainIDWithIndex(t, 0),
					IsSupported:     true,
					BallotThreshold: threshold,
				},
			},
		})

		voter := sample.AccAddress()
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{
				voter,
				sample.AccAddress(),
			},
		})

		chain, _ := k.GetSupportedChainFromChainID(ctx, getValidEthChainIDWithIndex(t, 0))
		index := sample.ZetaIndex(t)
		ballot, isFinalized, isNew, err := k.VoteOnBallot(
			ctx,
			chain,
			index,
			types.ObservationType_InboundTx,
			voter,
			types.VoteType_SuccessObservation)

		require.NoError(t, err)
		require.False(t, isFinalized)
		require.True(t, isNew)
		expectedBallot, found := k.GetBallot(ctx, index)
		require.True(t, found)
		require.Equal(t, expectedBallot, ballot)
	})

	t.Run("user can add a vote to an existing ballot", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: true,
				},
			},
		})

		voter := sample.AccAddress()
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{voter},
		})

		chain, _ := k.GetSupportedChainFromChainID(ctx, getValidEthChainIDWithIndex(t, 0))
		index := sample.ZetaIndex(t)
		threshold, err := sdkmath.LegacyNewDecFromStr("0.7")
		require.NoError(t, err)

		ballot := types.Ballot{
			BallotIdentifier: index,
			VoterList: []string{
				sample.AccAddress(),
				sample.AccAddress(),
				voter,
				sample.AccAddress(),
				sample.AccAddress(),
			},
			Votes:           types.CreateVotes(5),
			ObservationType: types.ObservationType_OutboundTx,
			BallotThreshold: threshold,
			BallotStatus:    types.BallotStatus_BallotInProgress,
		}
		k.SetBallot(ctx, &ballot)

		ballot, isFinalized, isNew, err := k.VoteOnBallot(
			ctx,
			chain,
			index,
			types.ObservationType_InboundTx,
			voter,
			types.VoteType_SuccessObservation)

		require.NoError(t, err)
		require.False(t, isFinalized)
		require.False(t, isNew)
		expectedBallot, found := k.GetBallot(ctx, index)
		require.True(t, found)
		require.Equal(t, expectedBallot, ballot)
	})

	t.Run("user can add a vote to an existing ballot, and finalize it", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: true,
				},
			},
		})

		voter := sample.AccAddress()
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{voter},
		})

		index := sample.ZetaIndex(t)
		threshold, err := sdkmath.LegacyNewDecFromStr("0.1")
		require.NoError(t, err)

		ballot := types.Ballot{
			BallotIdentifier: index,
			VoterList: []string{
				sample.AccAddress(),
				sample.AccAddress(),
				voter,
				sample.AccAddress(),
				sample.AccAddress(),
			},
			Votes:           types.CreateVotes(5),
			ObservationType: types.ObservationType_OutboundTx,
			BallotThreshold: threshold,
			BallotStatus:    types.BallotStatus_BallotInProgress,
		}
		k.SetBallot(ctx, &ballot)

		chain, _ := k.GetSupportedChainFromChainID(ctx, getValidEthChainIDWithIndex(t, 0))
		ballot, isFinalized, isNew, err := k.VoteOnBallot(
			ctx,
			chain,
			index,
			types.ObservationType_InboundTx,
			voter,
			types.VoteType_SuccessObservation)

		require.NoError(t, err)
		require.True(t, isFinalized)
		require.False(t, isNew)
		expectedBallot, found := k.GetBallot(ctx, index)
		require.True(t, found)
		require.Equal(t, expectedBallot, ballot)
	})
}
