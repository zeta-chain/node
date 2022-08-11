package wasm

import (
	"encoding/json"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	zetaCoreModuleKeeper "github.com/zeta-chain/zetacore/x/zetacore/keeper"
)

func Plugins(keeper zetaCoreModuleKeeper.Keeper) *wasmkeeper.QueryPlugins {
	return &wasmkeeper.QueryPlugins{
		Custom: ZetaCoreQuerier(keeper),
	}
}

func ZetaCoreQuerier(keeper zetaCoreModuleKeeper.Keeper) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return func(context sdk.Context, request json.RawMessage) ([]byte, error) {
		var custom ZetaCoreQuery
		err := json.Unmarshal(request, &custom)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}
		switch {
		case custom.WatchList != nil:
			return PerformWatchlistQuery(keeper, context)
		}
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Zetacore query variant")
	}
}

func PerformWatchlistQuery(
	keeper zetaCoreModuleKeeper.Keeper,
	ctx sdk.Context) ([]byte, error) {

	list := keeper.GetAllOutTxTracker(ctx)

	resp := WatchlistQueryResponse{Watchlist: list}
	return json.Marshal(resp)
}
