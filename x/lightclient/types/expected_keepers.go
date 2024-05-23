package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
)

type AuthorityKeeper interface {
	IsAuthorized(ctx sdk.Context, address string, policyType authoritytypes.PolicyType) bool
}
