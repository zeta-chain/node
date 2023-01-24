package types

import (
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
		MaxBondFactor:   "1.25",
		MinBondFactor:   "0.75",
		AvgBlockTime:    "6.00",
		TargetBondRatio: "67.00",
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

func validateMaxBondFactor(i interface{}) error {
	//v, ok := i.([]*string)
	//fmt.Println(v)
	//if !ok {
	//	return fmt.Errorf("invalid parameter type: %T", i)
	//}
	return nil
}

func validateMinBondFactor(i interface{}) error {
	//v, ok := i.([]*string)
	//fmt.Println(v)
	//if !ok {
	//	return fmt.Errorf("invalid parameter type: %T", i)
	//}
	return nil
}

func validateAvgBlockTime(i interface{}) error {
	//v, ok := i.([]*string)
	//fmt.Println(v)
	//if !ok {
	//	return fmt.Errorf("invalid parameter type: %T", i)
	//}
	return nil
}

func validateTargetBondRatio(i interface{}) error {
	//v, ok := i.([]*string)
	//fmt.Println(v)
	//if !ok {
	//	return fmt.Errorf("invalid parameter type: %T", i)
	//}
	return nil
}
