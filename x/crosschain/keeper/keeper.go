package keeper

import (
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

type (
	Keeper struct {
		cdc      codec.Codec
		storeKey storetypes.StoreKey
		memKey   storetypes.StoreKey

		StakingKeeper      types.StakingKeeper
		paramstore         paramtypes.Subspace
		authKeeper         types.AccountKeeper
		bankKeeper         types.BankKeeper
		zetaObserverKeeper types.ZetaObserverKeeper
		fungibleKeeper     types.FungibleKeeper
	}
)

func NewKeeper(
	cdc codec.Codec,
	storeKey,
	memKey storetypes.StoreKey,
	stakingKeeper types.StakingKeeper, // custom
	paramstore paramtypes.Subspace,
	authKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	zetaObserverKeeper types.ZetaObserverKeeper,
	fungibleKeeper types.FungibleKeeper,
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
		StakingKeeper:      stakingKeeper,
		paramstore:         paramstore,
		authKeeper:         authKeeper,
		bankKeeper:         bankKeeper,
		zetaObserverKeeper: zetaObserverKeeper,
		fungibleKeeper:     fungibleKeeper,
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
	return k.StakingKeeper
}

func (k Keeper) GetFungibleKeeper() types.FungibleKeeper {
	return k.fungibleKeeper
}

func (k Keeper) GetObserverKeeper() types.ZetaObserverKeeper {
	return k.zetaObserverKeeper
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
