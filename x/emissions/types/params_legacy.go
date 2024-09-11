/*
NOTE: Usage of x/params to manage parameters is deprecated in favor of x/gov
controlled execution of MsgUpdateParams messages. These types remains solely
for migration purposes and will be removed in a future release.
*/
package types

import paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(
			KeyPrefix(ParamValidatorEmissionPercentage),
			&p.ValidatorEmissionPercentage,
			validateValidatorEmissionPercentage,
		),
		paramtypes.NewParamSetPair(
			KeyPrefix(ParamObserverEmissionPercentage),
			&p.ObserverEmissionPercentage,
			validateObserverEmissionPercentage,
		),
		paramtypes.NewParamSetPair(
			KeyPrefix(ParamTssSignerEmissionPercentage),
			&p.TssSignerEmissionPercentage,
			validateTssEmissionPercentage,
		),
	}
}
