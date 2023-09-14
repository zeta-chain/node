package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/zeta-chain/zetacore/x/fungible/types"
)

type (
	Keeper struct {
		cdc            codec.BinaryCodec
		storeKey       storetypes.StoreKey
		memKey         storetypes.StoreKey
		paramstore     paramtypes.Subspace
		authKeeper     types.AccountKeeper
		evmKeeper      types.EVMKeeper
		bankKeeper     types.BankKeeper
		observerKeeper types.ObserverKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	authKeeper types.AccountKeeper,
	evmKeeper types.EVMKeeper,
	bankKeeper types.BankKeeper,
	observerKeeper types.ObserverKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		memKey:         memKey,
		paramstore:     ps,
		authKeeper:     authKeeper,
		evmKeeper:      evmKeeper,
		bankKeeper:     bankKeeper,
		observerKeeper: observerKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetAuthKeeper() types.AccountKeeper {
	return k.authKeeper
}

func (k Keeper) GetEVMKeeper() types.EVMKeeper {
	return k.evmKeeper
}

func (k Keeper) GetBankKeeper() types.BankKeeper {
	return k.bankKeeper
}

func (k Keeper) GetObserverKeeper() types.ObserverKeeper {
	return k.observerKeeper
}
