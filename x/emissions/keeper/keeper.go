package keeper

import (
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/zeta-chain/zetacore/x/emissions/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type (
	Keeper struct {
		cdc              codec.BinaryCodec
		storeKey         storetypes.StoreKey
		memKey           storetypes.StoreKey
		paramstore       paramtypes.Subspace
		feeCollectorName string
		bankKeeper       types.BankKeeper
		stakingKeeper    types.StakingKeeper
		observerKeeper   types.ObserverKeeper
		authKeeper       types.AccountKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	feeCollectorName string,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	observerKeeper types.ObserverKeeper,
	authKeeper types.AccountKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:              cdc,
		storeKey:         storeKey,
		memKey:           memKey,
		paramstore:       ps,
		feeCollectorName: feeCollectorName,
		bankKeeper:       bankKeeper,
		stakingKeeper:    stakingKeeper,
		observerKeeper:   observerKeeper,
		authKeeper:       authKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetFeeCollector() string {
	return k.feeCollectorName
}

func (k Keeper) GetBankKeeper() types.BankKeeper {
	return k.bankKeeper
}

func (k Keeper) GetStakingKeeper() types.StakingKeeper {
	return k.stakingKeeper
}

func (k Keeper) GetObserverKeeper() types.ObserverKeeper {
	return k.observerKeeper
}

func (k Keeper) GetAuthKeeper() types.AccountKeeper {
	return k.authKeeper
}
