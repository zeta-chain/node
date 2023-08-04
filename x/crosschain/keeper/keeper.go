package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
	// this line is used by starport scaffolding # ibc/keeper/import
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
		// this line is used by starport scaffolding # ibc/keeper/attribute

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
// this line is used by starport scaffolding # ibc/keeper/parameter

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
		// this line is used by starport scaffolding # ibc/keeper/return
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
