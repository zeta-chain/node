package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   sdk.StoreKey
		memKey     sdk.StoreKey
		paramstore paramtypes.Subspace
		authKeeper types.AccountKeeper
		evmKeeper  evmkeeper.Keeper
		bankKeeper types.BankKeeper
		//zetacoreKeeper *zetacorekeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	authKeeper types.AccountKeeper,
	evmKeeper evmkeeper.Keeper,
	bankKeeper types.BankKeeper,
	//zetacoreKeeper *zetacorekeeper.Keeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{

		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
		authKeeper: authKeeper,
		evmKeeper:  evmKeeper,
		bankKeeper: bankKeeper,
		//zetacoreKeeper: zetacoreKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
