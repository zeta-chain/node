package types

import paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for zetaObserver module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyPrefix(ObserverParamsKey), &p.ObserverParams, validateVotingThresholds),
		paramtypes.NewParamSetPair(KeyPrefix(AdminPolicyParamsKey), &p.AdminPolicy, validateAdminPolicy),
		paramtypes.NewParamSetPair(KeyPrefix(BallotMaturityBlocksParamsKey), &p.BallotMaturityBlocks, validateBallotMaturityBlocks),
	}
}
