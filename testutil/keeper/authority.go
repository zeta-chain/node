package keeper

import (
	"testing"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/authority/keeper"
	"github.com/zeta-chain/node/x/authority/types"
)

var (
	AuthorityGovAddress = sample.Bech32AccAddress()
)

func initAuthorityKeeper(
	cdc codec.Codec,
	ss store.CommitMultiStore,
) keeper.Keeper {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, nil)
	ss.MountStoreWithDB(memKey, storetypes.StoreTypeMemory, nil)

	return keeper.NewKeeper(
		cdc,
		storeKey,
		memKey,
		AuthorityGovAddress,
	)
}

// AuthorityKeeper instantiates an authority keeper for testing purposes
func AuthorityKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	// Initialize local store
	db := tmdb.NewMemDB()
	stateStore := rootmulti.NewStore(db, log.NewNopLogger())
	cdc := NewCodec()

	// Create regular keepers
	sdkKeepers := NewSDKKeepers(cdc, db, stateStore)

	// Create the observer keeper
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, nil)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := NewContext(stateStore)

	// Initialize modules genesis
	sdkKeepers.InitGenesis(ctx)

	// Add a proposer to the context
	ctx = sdkKeepers.InitBlockProposer(t, ctx)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		AuthorityGovAddress,
	)

	return &k, ctx
}

// MockCheckAuthorization mocks the CheckAuthorization method of the authority keeper.
func MockCheckAuthorization(m *mock.Mock, msg sdk.Msg, authorizationResult error) {
	m.On("CheckAuthorization", mock.Anything, msg).Return(authorizationResult).Once()
}

// MockGetChainList mocks the GetAdditionalChainList method of the authority keeper.
func MockGetChainList(m *mock.Mock, chainList []chains.Chain) {
	m.On("GetAdditionalChainList", mock.Anything).Return(chainList).Once()
}

// MockGetChainListEmpty mocks the GetAdditionalChainList method of the authority keeper.
func MockGetChainListEmpty(m *mock.Mock) {
	m.On("GetAdditionalChainList", mock.Anything).Return([]chains.Chain{})
}

func SetAdminPolicies(ctx sdk.Context, ak *keeper.Keeper) string {
	admin := sample.AccAddress()
	ak.SetPolicies(ctx, types.Policies{Items: []*types.Policy{
		{
			Address:    admin,
			PolicyType: types.PolicyType_groupAdmin,
		},
	}})
	return admin
}
