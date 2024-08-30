package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/chains"
)

type AuthorityKeeper interface {
	CheckAuthorization(ctx sdk.Context, msg sdk.Msg) error
	GetAdditionalChainList(ctx sdk.Context) (list []chains.Chain)
}
