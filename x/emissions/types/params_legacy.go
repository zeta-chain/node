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
		paramtypes.NewParamSetPair(KeyPrefix(ParamMaxBondFactor), &p.MaxBondFactor, validateMaxBondFactor),
		paramtypes.NewParamSetPair(KeyPrefix(ParamMinBondFactor), &p.MinBondFactor, validateMinBondFactor),
		paramtypes.NewParamSetPair(KeyPrefix(ParamAvgBlockTime), &p.AvgBlockTime, validateAvgBlockTime),
		paramtypes.NewParamSetPair(KeyPrefix(ParamTargetBondRatio), &p.TargetBondRatio, validateTargetBondRatio),
		paramtypes.NewParamSetPair(KeyPrefix(ParamValidatorEmissionPercentage), &p.ValidatorEmissionPercentage, validateValidatorEmissionPercentage),
		paramtypes.NewParamSetPair(KeyPrefix(ParamObserverEmissionPercentage), &p.ObserverEmissionPercentage, validateObserverEmissionPercentage),
		paramtypes.NewParamSetPair(KeyPrefix(ParamTssSignerEmissionPercentage), &p.TssSignerEmissionPercentage, validateTssEmissionPercentage),
		paramtypes.NewParamSetPair(KeyPrefix(ParamDurationFactorConstant), &p.DurationFactorConstant, validateDurationFactorConstant),

		// TODO: enable this param
		// https://github.com/zeta-chain/node/pull/1861
		//paramtypes.NewParamSetPair(KeyPrefix(ParamObserverSlashAmount), &p.ObserverSlashAmount, validateObserverSlashAmount),
	}
}
