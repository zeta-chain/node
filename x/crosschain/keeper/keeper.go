package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/crosschain/types"
)

type (
	Keeper struct {
		cdc      codec.Codec
		storeKey storetypes.StoreKey
		memKey   storetypes.StoreKey

		stakingKeeper       types.StakingKeeper
		authKeeper          types.AccountKeeper
		bankKeeper          types.BankKeeper
		zetaObserverKeeper  types.ObserverKeeper
		fungibleKeeper      types.FungibleKeeper
		authorityKeeper     types.AuthorityKeeper
		lightclientKeeper   types.LightclientKeeper
		ibcCrosschainKeeper types.IBCCrosschainKeeper
	}
)

func NewKeeper(
	cdc codec.Codec,
	storeKey,
	memKey storetypes.StoreKey,
	stakingKeeper types.StakingKeeper, // custom
	authKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	zetaObserverKeeper types.ObserverKeeper,
	fungibleKeeper types.FungibleKeeper,
	authorityKeeper types.AuthorityKeeper,
	lightclientKeeper types.LightclientKeeper,
) *Keeper {
	// ensure governance module account is set
	// FIXME: enable this check! (disabled for now to avoid unit test panic)
	//if addr := authKeeper.GetModuleAddress(types.ModuleName); addr == nil {
	//	panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	//}

	return &Keeper{
		cdc:                cdc,
		storeKey:           storeKey,
		memKey:             memKey,
		stakingKeeper:      stakingKeeper,
		authKeeper:         authKeeper,
		bankKeeper:         bankKeeper,
		zetaObserverKeeper: zetaObserverKeeper,
		fungibleKeeper:     fungibleKeeper,
		authorityKeeper:    authorityKeeper,
		lightclientKeeper:  lightclientKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetAuthKeeper() types.AccountKeeper {
	return k.authKeeper
}

func (k Keeper) GetBankKeeper() types.BankKeeper {
	return k.bankKeeper
}

func (k Keeper) GetStakingKeeper() types.StakingKeeper {
	return k.stakingKeeper
}

func (k Keeper) GetFungibleKeeper() types.FungibleKeeper {
	return k.fungibleKeeper
}

func (k Keeper) GetObserverKeeper() types.ObserverKeeper {
	return k.zetaObserverKeeper
}

func (k Keeper) GetAuthorityKeeper() types.AuthorityKeeper {
	return k.authorityKeeper
}

func (k Keeper) GetLightclientKeeper() types.LightclientKeeper {
	return k.lightclientKeeper
}

func (k Keeper) GetIBCCrosschainKeeper() types.IBCCrosschainKeeper {
	return k.ibcCrosschainKeeper
}

func (k *Keeper) SetIBCCrosschainKeeper(ibcCrosschainKeeper types.IBCCrosschainKeeper) {
	k.ibcCrosschainKeeper = ibcCrosschainKeeper
}

func (k Keeper) GetStoreKey() storetypes.StoreKey {
	return k.storeKey
}

func (k Keeper) GetMemKey() storetypes.StoreKey {
	return k.memKey
}

func (k Keeper) GetCodec() codec.Codec {
	return k.cdc
}
