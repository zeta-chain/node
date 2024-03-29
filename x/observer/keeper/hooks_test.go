package keeper_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
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

	hooks.AfterValidatorRemoved(ctx, nil, valAddr)

	os, found := k.GetObserverSet(ctx)
	require.True(t, found)
	require.Empty(t, os.ObserverList)
}
