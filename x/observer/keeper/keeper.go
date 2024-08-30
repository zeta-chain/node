package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/observer/types"
)

type (
	Keeper struct {
		cdc               codec.BinaryCodec
		storeKey          storetypes.StoreKey
		memKey            storetypes.StoreKey
		stakingKeeper     types.StakingKeeper
		slashingKeeper    types.SlashingKeeper
		authorityKeeper   types.AuthorityKeeper
		lightclientKeeper types.LightclientKeeper
		authority         string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	stakingKeeper types.StakingKeeper,
	slashinKeeper types.SlashingKeeper,
	authorityKeeper types.AuthorityKeeper,
	lightclientKeeper types.LightclientKeeper,
	authority string,
) *Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(err)
	}

	return &Keeper{
		cdc:               cdc,
		storeKey:          storeKey,
		memKey:            memKey,
		stakingKeeper:     stakingKeeper,
		slashingKeeper:    slashinKeeper,
		authorityKeeper:   authorityKeeper,
		lightclientKeeper: lightclientKeeper,
		authority:         authority,
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

func (k Keeper) GetAuthority() string {
	return k.authority
}
