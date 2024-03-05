package keeper_test

import (
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/stretchr/testify/require"
	zetacommon "github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
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
		return zetacommon.GoerliLocalnetChain().ChainId
	case 1:
		return zetacommon.GoerliChain().ChainId
	default:
		require.Fail(t, "invalid index")
	}
	return 0
}

func TestKeeper_IsAuthorized(t *testing.T) {
	t.Run("authorized observer", func(t *testing.T) {
		k, ctx, _ := keepertest.ObserverKeeper(t)

		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		k.GetStakingKeeper().SetValidator(ctx, validator)
		consAddress, err := validator.GetConsAddr()
		require.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             consAddress.String(),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          false,
			MissedBlocksCounter: 1,
		})

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)

		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String()},
		})
		require.True(t, k.IsAuthorized(ctx, accAddressOfValidator.String()))

	})
	t.Run("not authorized for tombstoned observer", func(t *testing.T) {
		k, ctx, _ := keepertest.ObserverKeeper(t)

		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)
		k.GetStakingKeeper().SetValidator(ctx, validator)
		consAddress, err := validator.GetConsAddr()
		require.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             consAddress.String(),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          true,
			MissedBlocksCounter: 1,
		})

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String()},
		})

		require.False(t, k.IsAuthorized(ctx, accAddressOfValidator.String()))

	})
	t.Run("not authorized for non-validator observer", func(t *testing.T) {
		k, ctx, _ := keepertest.ObserverKeeper(t)

		r := rand.New(rand.NewSource(9))

		// Set validator in the store
		validator := sample.Validator(t, r)

		consAddress, err := validator.GetConsAddr()
		require.NoError(t, err)
		k.GetSlashingKeeper().SetValidatorSigningInfo(ctx, consAddress, slashingtypes.ValidatorSigningInfo{
			Address:             consAddress.String(),
			StartHeight:         0,
			JailedUntil:         ctx.BlockHeader().Time.Add(1000000 * time.Second),
			Tombstoned:          false,
			MissedBlocksCounter: 1,
		})

		accAddressOfValidator, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{accAddressOfValidator.String()},
		})

		require.False(t, k.IsAuthorized(ctx, accAddressOfValidator.String()))

	})
}
