package wasm

import zetaCoreModuleTypes "github.com/zeta-chain/zetacore/x/zetacore/types"

type ZetaCoreQuery struct {
	WatchList *WatchlistQuery `json:"watch_list"`
}

type WatchlistQuery struct{}

type WatchlistQueryResponse struct {
	Watchlist []zetaCoreModuleTypes.OutTxTracker
}
