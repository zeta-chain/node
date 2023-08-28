package keeper

import (
	"math/rand"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ethermint "github.com/evmos/ethermint/types"
	evmmodule "github.com/evmos/ethermint/x/evm"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/evmos/ethermint/x/evm/vm/geth"
	feemarketkeeper "github.com/evmos/ethermint/x/feemarket/keeper"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"github.com/stretchr/testify/require"
	tmdb "github.com/tendermint/tm-db"
	"github.com/zeta-chain/zetacore/testutil/sample"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	emissionstypes "github.com/zeta-chain/zetacore/x/emissions/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

// SDKKeepers is a struct containing regular SDK module keepers for test purposes
type SDKKeepers struct {
	ParamsKeeper    paramskeeper.Keeper
	AuthKeeper      authkeeper.AccountKeeper
	BankKeeper      bankkeeper.Keeper
	StakingKeeper   stakingkeeper.Keeper
	FeeMarketKeeper feemarketkeeper.Keeper
	EvmKeeper       *evmkeeper.Keeper
}

var moduleAccountPerms = map[string][]string{
	authtypes.FeeCollectorName:                      nil,
	distrtypes.ModuleName:                           nil,
	stakingtypes.BondedPoolName:                     {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName:                  {authtypes.Burner, authtypes.Staking},
	evmtypes.ModuleName:                             {authtypes.Minter, authtypes.Burner},
	crosschaintypes.ModuleName:                      {authtypes.Minter, authtypes.Burner},
	fungibletypes.ModuleName:                        {authtypes.Minter, authtypes.Burner},
	emissionstypes.ModuleName:                       nil,
	emissionstypes.UndistributedObserverRewardsPool: nil,
	emissionstypes.UndistributedTssRewardsPool:      nil,
}

// ModuleAccountAddrs returns all the app's module account addresses.
func ModuleAccountAddrs(maccPerms map[string][]string) map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// ParamsKeeper instantiates a param keeper for testing purposes
// TODO: remove https://github.com/zeta-chain/node/issues/848
func ParamsKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
) paramskeeper.Keeper {
	storeKey := sdk.NewKVStoreKey(paramstypes.StoreKey)
	tkeys := sdk.NewTransientStoreKey(paramstypes.TStoreKey)

	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	ss.MountStoreWithDB(tkeys, storetypes.StoreTypeTransient, db)

	return paramskeeper.NewKeeper(
		cdc,
		fungibletypes.Amino,
		storeKey,
		tkeys,
	)
}

// AccountKeeper instantiates an account keeper for testing purposes
func AccountKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
	paramKeeper paramskeeper.Keeper,
) authkeeper.AccountKeeper {
	storeKey := sdk.NewKVStoreKey(authtypes.StoreKey)
	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	return authkeeper.NewAccountKeeper(
		cdc,
		storeKey,
		paramKeeper.Subspace(authtypes.ModuleName),
		ethermint.ProtoAccount,
		moduleAccountPerms,
		"zeta",
	)
}

// BankKeeper instantiates a bank keeper for testing purposes
func BankKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
	paramKeeper paramskeeper.Keeper,
	authKeeper authkeeper.AccountKeeper,
) bankkeeper.Keeper {
	storeKey := sdk.NewKVStoreKey(banktypes.StoreKey)
	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	blockedAddrs := make(map[string]bool)

	return bankkeeper.NewBaseKeeper(
		cdc,
		storeKey,
		authKeeper,
		paramKeeper.Subspace(banktypes.ModuleName),
		blockedAddrs,
	)
}

// StakingKeeper instantiates a staking keeper for testing purposes
func StakingKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
	authKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	paramKeeper paramskeeper.Keeper,
) stakingkeeper.Keeper {
	storeKey := sdk.NewKVStoreKey(stakingtypes.StoreKey)
	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	return stakingkeeper.NewKeeper(
		cdc,
		storeKey,
		authKeeper,
		bankKeeper,
		paramKeeper.Subspace(stakingtypes.ModuleName),
	)
}

// DistributionKeeper instantiates a distribution keeper for testing purposes
func DistributionKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
	authKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	stakingKeeper *stakingkeeper.Keeper,
	paramKeeper paramskeeper.Keeper,
) distrkeeper.Keeper {
	storeKey := sdk.NewKVStoreKey(distrtypes.StoreKey)
	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	return distrkeeper.NewKeeper(
		cdc,
		storeKey,
		paramKeeper.Subspace(stakingtypes.ModuleName),
		authKeeper,
		bankKeeper,
		stakingKeeper,
		authtypes.FeeCollectorName,
	)
}

// ProtocolVersionSetter mock
type ProtocolVersionSetter struct{}

func (vs ProtocolVersionSetter) SetProtocolVersion(uint64) {}

// UpgradeKeeper instantiates an upgrade keeper for testing purposes
func UpgradeKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
) upgradekeeper.Keeper {
	storeKey := sdk.NewKVStoreKey(upgradetypes.StoreKey)
	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	skipUpgradeHeights := make(map[int64]bool)
	vs := ProtocolVersionSetter{}

	return upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		storeKey,
		cdc,
		"",
		vs,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
}

// FeeMarketKeeper instantiates a feemarket keeper for testing purposes
func FeeMarketKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
	paramKeeper paramskeeper.Keeper,
) feemarketkeeper.Keeper {
	storeKey := sdk.NewKVStoreKey(feemarkettypes.StoreKey)
	transientKey := sdk.NewTransientStoreKey(feemarkettypes.TransientKey)

	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	ss.MountStoreWithDB(transientKey, storetypes.StoreTypeTransient, db)

	return feemarketkeeper.NewKeeper(
		cdc,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		storeKey,
		transientKey,
		paramKeeper.Subspace(feemarkettypes.ModuleName),
	)
}

// EVMKeeper instantiates an evm keeper for testing purposes
func EVMKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
	authKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	stakingKeeper stakingkeeper.Keeper,
	feemarketKeeper feemarketkeeper.Keeper,
	paramKeeper paramskeeper.Keeper,
) *evmkeeper.Keeper {
	storeKey := sdk.NewKVStoreKey(evmtypes.StoreKey)
	transientKey := sdk.NewTransientStoreKey(evmtypes.TransientKey)

	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	ss.MountStoreWithDB(transientKey, storetypes.StoreTypeTransient, db)

	k := evmkeeper.NewKeeper(
		cdc,
		storeKey,
		transientKey,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		authKeeper,
		bankKeeper,
		stakingKeeper,
		feemarketKeeper,
		nil,
		geth.NewEVM,
		"",
		paramKeeper.Subspace(evmtypes.ModuleName),
	)

	return k
}

// NewSDKKeepers instantiates regular Cosmos SDK keeper such as staking with local storage for testing purposes
func NewSDKKeepers(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
) SDKKeepers {
	paramsKeeper := ParamsKeeper(cdc, db, ss)
	authKeeper := AccountKeeper(cdc, db, ss, paramsKeeper)
	bankKeeper := BankKeeper(cdc, db, ss, paramsKeeper, authKeeper)
	stakingKeeper := StakingKeeper(cdc, db, ss, authKeeper, bankKeeper, paramsKeeper)
	feeMarketKeeper := FeeMarketKeeper(cdc, db, ss, paramsKeeper)
	evmKeeper := EVMKeeper(cdc, db, ss, authKeeper, bankKeeper, stakingKeeper, feeMarketKeeper, paramsKeeper)

	return SDKKeepers{
		ParamsKeeper:    paramsKeeper,
		AuthKeeper:      authKeeper,
		BankKeeper:      bankKeeper,
		StakingKeeper:   stakingKeeper,
		FeeMarketKeeper: feeMarketKeeper,
		EvmKeeper:       evmKeeper,
	}
}

// InitGenesis initializes the test modules genesis state
func (sdkm SDKKeepers) InitGenesis(ctx sdk.Context) {
	sdkm.AuthKeeper.InitGenesis(ctx, *authtypes.DefaultGenesisState())
	sdkm.BankKeeper.InitGenesis(ctx, banktypes.DefaultGenesisState())
	sdkm.StakingKeeper.InitGenesis(ctx, stakingtypes.DefaultGenesisState())
	evmmodule.InitGenesis(ctx, sdkm.EvmKeeper, sdkm.AuthKeeper, *evmtypes.DefaultGenesisState())
}

// InitBlockProposer initialize the block proposer for test purposes with an associated validator
func (sdkm SDKKeepers) InitBlockProposer(t testing.TB, ctx sdk.Context) sdk.Context {
	// #nosec G404 test purpose - weak randomness is not an issue here
	r := rand.New(rand.NewSource(42))

	// Set validator in the store
	validator := sample.Validator(t, r)
	sdkm.StakingKeeper.SetValidator(ctx, validator)
	err := sdkm.StakingKeeper.SetValidatorByConsAddr(ctx, validator)
	require.NoError(t, err)

	// Validator is proposer
	consAddr, err := validator.GetConsAddr()
	require.NoError(t, err)
	return ctx.WithProposer(consAddr)
}
