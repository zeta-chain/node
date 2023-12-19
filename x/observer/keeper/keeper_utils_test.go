package keeper_test

import (
	"math/rand"
	"testing"
	"time"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/stretchr/testify/assert"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestKeeper_IsAuthorized(t *testing.T) {
	t.Run("authorized observer", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		chains := k.GetParams(ctx).GetSupportedChains()

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
		chains := k.GetParams(ctx).GetSupportedChains()

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
		chains := k.GetParams(ctx).GetSupportedChains()

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
