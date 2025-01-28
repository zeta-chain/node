package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/observer/types"
)

type (
	Keeper struct {
		cdc               codec.Codec
		storeKey          storetypes.StoreKey
		memKey            storetypes.StoreKey
		stakingKeeper     types.StakingKeeper
		slashingKeeper    types.SlashingKeeper
		authorityKeeper   types.AuthorityKeeper
		lightclientKeeper types.LightclientKeeper
		bankKeeper        types.BankKeeper
		authKeeper        types.AccountKeeper
		authority         string
	}
)

func NewKeeper(
	cdc codec.Codec,
	storeKey,
	memKey storetypes.StoreKey,
	stakingKeeper types.StakingKeeper,
	slashinKeeper types.SlashingKeeper,
	authorityKeeper types.AuthorityKeeper,
	lightclientKeeper types.LightclientKeeper,
	bankKeeper types.BankKeeper,
	authKeeper types.AccountKeeper,
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
		bankKeeper:        bankKeeper,
		authKeeper:        authKeeper,
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

func (k Keeper) GetBankKeeper() types.BankKeeper {
	return k.bankKeeper
}

func (k Keeper) GetAuthKeeper() types.AccountKeeper {
	return k.authKeeper
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

func (k Keeper) Codec() codec.Codec {
	return k.cdc
}

func (k Keeper) GetAuthority() string {
	return k.authority
}
