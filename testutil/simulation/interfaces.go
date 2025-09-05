package simulation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/chains"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

type ObserverKeeper interface {
	GetObserverSet(ctx sdk.Context) (val observertypes.ObserverSet, found bool)
	CheckObserverCanVote(ctx sdk.Context, address string) error
	GetSupportedChains(ctx sdk.Context) []chains.Chain
	GetNodeAccount(ctx sdk.Context, address string) (observertypes.NodeAccount, bool)
	GetAllNodeAccount(ctx sdk.Context) []observertypes.NodeAccount
}

type AuthorityKeeper interface {
	CheckAuthorization(ctx sdk.Context, msg sdk.Msg) error
	GetAdditionalChainList(ctx sdk.Context) (list []chains.Chain)
	GetPolicies(ctx sdk.Context) (val authoritytypes.Policies, found bool)
}

type FungibleKeeper interface {
	GetForeignCoins(ctx sdk.Context, zrc20Addr string) (val fungibletypes.ForeignCoins, found bool)
	GetAllForeignCoins(ctx sdk.Context) (list []fungibletypes.ForeignCoins)
}
