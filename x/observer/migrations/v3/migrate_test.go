package v3_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	v3 "github.com/zeta-chain/zetacore/x/observer/migrations/v3"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMigrateStore(t *testing.T) {
	k, ctx := keepertest.ObserverKeeper(t)

	// nothing if no admin policy
	params := types.DefaultParams()
	params.AdminPolicy = []*types.Admin_Policy{}
	k.SetParams(ctx, params)
	err := v3.MigrateStore(ctx, k)
	assert.NoError(t, err)
	params = k.GetParams(ctx)
	assert.Len(t, params.AdminPolicy, 0)

	// update admin policy
	admin := sample.AccAddress()
	params = types.DefaultParams()
	params.AdminPolicy = []*types.Admin_Policy{
		{
			Address:    admin,
			PolicyType: 0,
		},
		{
			Address:    sample.AccAddress(),
			PolicyType: 5,
		},
		{
			Address:    admin,
			PolicyType: 10,
		},
	}
	k.SetParams(ctx, params)
	err = v3.MigrateStore(ctx, k)
	assert.NoError(t, err)
	params = k.GetParams(ctx)
	assert.Len(t, params.AdminPolicy, 2)
	assert.Equal(t, params.AdminPolicy[0].PolicyType, types.Policy_Type_group1)
	assert.Equal(t, params.AdminPolicy[1].PolicyType, types.Policy_Type_group2)
	assert.Equal(t, params.AdminPolicy[0].Address, admin)
	assert.Equal(t, params.AdminPolicy[1].Address, admin)
}
