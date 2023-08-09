package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		MaxBondFactor:               "1.25",
		MinBondFactor:               "0.75",
		AvgBlockTime:                "6.00",
		TargetBondRatio:             "00.67",
		ValidatorEmissionPercentage: "00.50",
		ObserverEmissionPercentage:  "00.25",
		TssSignerEmissionPercentage: "00.25",
		DurationFactorConstant:      "0.001877876953694702",
		ObserverSlashAmount:         "100000000000000000",
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyPrefix(ParamMaxBondFactor), &p.MaxBondFactor, validateMaxBondFactor),
		paramtypes.NewParamSetPair(KeyPrefix(ParamMinBondFactor), &p.MinBondFactor, validateMinBondFactor),
		paramtypes.NewParamSetPair(KeyPrefix(ParamAvgBlockTime), &p.AvgBlockTime, validateAvgBlockTime),
		paramtypes.NewParamSetPair(KeyPrefix(ParamTargetBondRatio), &p.TargetBondRatio, validateTargetBondRatio),
		paramtypes.NewParamSetPair(KeyPrefix(ParamValidatorEmissionPercentage), &p.ValidatorEmissionPercentage, validateValidatorEmissonPercentage),
		paramtypes.NewParamSetPair(KeyPrefix(ParamObserverEmissionPercentage), &p.ObserverEmissionPercentage, validateObserverEmissonPercentage),
		paramtypes.NewParamSetPair(KeyPrefix(ParamTssSignerEmissionPercentage), &p.TssSignerEmissionPercentage, validateTssEmissonPercentage),
		paramtypes.NewParamSetPair(KeyPrefix(ParamDurationFactorConstant), &p.DurationFactorConstant, validateDurationFactorConstant),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateDurationFactorConstant(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateMaxBondFactor(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	decMaxBond := sdk.MustNewDecFromStr(v)
	if decMaxBond.GT(sdk.MustNewDecFromStr("1.25")) {
		return fmt.Errorf("max bond factor cannot be higher that 0.25")
	}
	return nil
}

func validateMinBondFactor(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	decMaxBond := sdk.MustNewDecFromStr(v)
	if decMaxBond.LT(sdk.MustNewDecFromStr("0.75")) {
		return fmt.Errorf("min bond factor cannot be lower that 0.75")
	}
	return nil
}

func validateAvgBlockTime(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateTargetBondRatio(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	decMaxBond := sdk.MustNewDecFromStr(v)
	if decMaxBond.GT(sdk.MustNewDecFromStr("100.00")) {
		return fmt.Errorf("target bond ratio cannot be more than 100 percent")
	}
	if decMaxBond.LT(sdk.ZeroDec()) {
		return fmt.Errorf("target bond ratio cannot be less than 0 percent")
	}
	return nil
}

func validateValidatorEmissonPercentage(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	dec := sdk.MustNewDecFromStr(v)
	if dec.GT(sdk.MustNewDecFromStr("100.00")) {
		return fmt.Errorf("validator emission percentage cannot be more than 100 percent")
	}
	if dec.LT(sdk.ZeroDec()) {
		return fmt.Errorf("validator emission percentage cannot be less than 0 percent")
	}
	return nil
}

func validateObserverEmissonPercentage(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	dec := sdk.MustNewDecFromStr(v)
	if dec.GT(sdk.MustNewDecFromStr("100.00")) {
		return fmt.Errorf("observer emission percentage cannot be more than 100 percent")
	}
	if dec.LT(sdk.ZeroDec()) {
		return fmt.Errorf("observer emission percentage cannot be less than 0 percent")
	}
	return nil
}

func validateTssEmissonPercentage(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	dec := sdk.MustNewDecFromStr(v)
	if dec.GT(sdk.MustNewDecFromStr("100.00")) {
		return fmt.Errorf("tss emission percentage cannot be more than 100 percent")
	}
	if dec.LT(sdk.ZeroDec()) {
		return fmt.Errorf("tss emission percentage cannot be less than 0 percent")
	}
	return nil
}
