package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AuthorityKeeper interface {
	CheckAuthorization(ctx sdk.Context, msg sdk.Msg) error
}
