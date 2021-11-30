package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

type StakingKeeper interface {
	GetAllValidators(ctx sdk.Context) (validators []types.Validator)
}
