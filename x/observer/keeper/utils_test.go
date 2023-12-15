package keeper_test

import (
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/stretchr/testify/assert"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// setSupportedChain sets the supported chains for the observer module
func setSupportedChain(ctx sdk.Context, observerKeeper keeper.Keeper, chainIDs ...int64) {
	coreParamsList := make([]*types.CoreParams, len(chainIDs))
	for i, chainID := range chainIDs {
		coreParams := sample.CoreParams(chainID)
		coreParams.IsSupported = true
		coreParamsList[i] = coreParams
	}
	observerKeeper.SetCoreParamsList(ctx, types.CoreParamsList{
		CoreParams: coreParamsList,
	})
}

func TestKeeper_IsAuthorized(t *testing.T) {
	t.Run("authorized observer", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		chains := k.GetSupportedChains(ctx)

		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		k.GetStakingKeeper().SetValidator(ctx, validator)
		consAddress, err := validator.GetConsAddr()
		assert.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             consAddress.String(),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          false,
			MissedBlocksCounter: 1,
		})

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		for _, chain := range chains {
			k.SetObserverMapper(ctx, &types.ObserverMapper{
				ObserverChain: chain,
				ObserverList:  []string{accAddressOfValidator.String()},
			})
		}
		for _, chain := range chains {
			assert.True(t, k.IsAuthorized(ctx, accAddressOfValidator.String(), chain))
		}
	})
	t.Run("not authorized for tombstoned observer", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		chains := k.GetSupportedChains(ctx)

		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		k.GetStakingKeeper().SetValidator(ctx, validator)
		consAddress, err := validator.GetConsAddr()
		assert.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             consAddress.String(),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          true,
			MissedBlocksCounter: 1,
		})

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		for _, chain := range chains {
			k.SetObserverMapper(ctx, &types.ObserverMapper{
				ObserverChain: chain,
				ObserverList:  []string{accAddressOfValidator.String()},
			})
		}
		for _, chain := range chains {
			assert.False(t, k.IsAuthorized(ctx, accAddressOfValidator.String(), chain))
		}
	})
	t.Run("not authorized for non-validator observer", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		chains := k.GetSupportedChains(ctx)

		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)

		consAddress, err := validator.GetConsAddr()
		assert.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             consAddress.String(),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          false,
			MissedBlocksCounter: 1,
		})

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		for _, chain := range chains {
			k.SetObserverMapper(ctx, &types.ObserverMapper{
				ObserverChain: chain,
				ObserverList:  []string{accAddressOfValidator.String()},
			})
		}
		for _, chain := range chains {
			assert.False(t, k.IsAuthorized(ctx, accAddressOfValidator.String(), chain))
		}
	})
}
