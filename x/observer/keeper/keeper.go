package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

type (
	Keeper struct {
		cdc               codec.BinaryCodec
		storeKey          storetypes.StoreKey
		memKey            storetypes.StoreKey
		paramstore        paramtypes.Subspace
		stakingKeeper     types.StakingKeeper
		slashingKeeper    types.SlashingKeeper
		authorityKeeper   types.AuthorityKeeper
		lightclientKeeper types.LightclientKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	stakingKeeper types.StakingKeeper,
	slashinKeeper types.SlashingKeeper,
	authorityKeeper types.AuthorityKeeper,
	lightclientKeeper types.LightclientKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:               cdc,
		storeKey:          storeKey,
		memKey:            memKey,
		paramstore:        ps,
		stakingKeeper:     stakingKeeper,
		slashingKeeper:    slashinKeeper,
		authorityKeeper:   authorityKeeper,
		lightclientKeeper: lightclientKeeper,
	}
}

func (k Keeper) GetSlashingKeeper() types.SlashingKeeper {
	return k.slashingKeeper
}

func (k Keeper) GetStakingKeeper() types.StakingKeeper {
	return k.stakingKeeper
}

func (k Keeper) GetAuthorityKeeper() types.AuthorityKeeper {
	return k.authorityKeeper
}

func (k Keeper) GetLightclientKeeper() types.LightclientKeeper {
	return k.lightclientKeeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) StoreKey() storetypes.StoreKey {
	return k.storeKey
}

func (k Keeper) Codec() codec.BinaryCodec {
	return k.cdc
}
