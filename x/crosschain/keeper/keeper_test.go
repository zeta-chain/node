package keeper

import (
	"testing"

	authkeeper2 "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper2 "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	fungibleKeeper "github.com/zeta-chain/zetacore/x/fungible/keeper"
	"github.com/zeta-chain/zetacore/x/observer/keeper"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
	"github.com/zeta-chain/zetacore/x/crosschain/types"

	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
)

func setupKeeper(t testing.TB) (*Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"ZetacoreParams",
	)
	bankkeeper := bankkeeper2.BaseKeeper{}
	authkeeper := authkeeper2.AccountKeeper{}
	zetaobserverKeeper := keeper.Keeper{}
	fungibleKeeper := fungibleKeeper.Keeper{}

	k := NewKeeper(
		codec.NewProtoCodec(registry),
		storeKey,
		memStoreKey,
		stakingkeeper.Keeper{}, // custom
		paramsSubspace,
		authkeeper,
		bankkeeper,
		zetaobserverKeeper,
		fungibleKeeper,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	return k, ctx
}
