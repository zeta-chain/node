package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AuthorityKeeper interface {
	IsAuthorized(ctx sdk.Context, msg sdk.Msg) (bool, error)
}
