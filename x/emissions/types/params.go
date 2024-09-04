package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v2"
)

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		ValidatorEmissionPercentage: "00.50",
		ObserverEmissionPercentage:  "00.25",
		TssSignerEmissionPercentage: "00.25",
		ObserverSlashAmount:         ObserverSlashAmount,
		BallotMaturityBlocks:        int64(BallotMaturityBlocks),
		BlockRewardAmount:           BlockReward,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateValidatorEmissionPercentage(p.ValidatorEmissionPercentage); err != nil {
		return err
	}
	if err := validateObserverEmissionPercentage(p.ObserverEmissionPercentage); err != nil {
		return err
	}
	if err := validateTssEmissionPercentage(p.TssSignerEmissionPercentage); err != nil {
		return err
	}
	if err := validateBallotMaturityBlocks(p.BallotMaturityBlocks); err != nil {
		return err
	}
	if err := validateBlockRewardsAmount(p.BlockRewardAmount); err != nil {
		return err
	}
	return validateObserverSlashAmount(p.ObserverSlashAmount)
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, err := yaml.Marshal(p)
	if err != nil {
		return ""
	}
	return string(out)
}

func validateValidatorEmissionPercentage(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	dec := sdk.MustNewDecFromStr(v)
	if dec.GT(sdk.OneDec()) {
		return fmt.Errorf("validator emission percentage cannot be more than 100 percent")
	}
	if dec.LT(sdk.ZeroDec()) {
		return fmt.Errorf("validator emission percentage cannot be less than 0 percent")
	}
	return nil
}

func validateObserverEmissionPercentage(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	dec := sdk.MustNewDecFromStr(v)
	if dec.GT(sdk.OneDec()) {
		return fmt.Errorf("observer emission percentage cannot be more than 100 percent")
	}
	if dec.LT(sdk.ZeroDec()) {
		return fmt.Errorf("observer emission percentage cannot be less than 0 percent")
	}
	return nil
}

func validateTssEmissionPercentage(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	dec := sdk.MustNewDecFromStr(v)
	if dec.GT(sdk.OneDec()) {
		return fmt.Errorf("tss emission percentage cannot be more than 100 percent")
	}
	if dec.LT(sdk.ZeroDec()) {
		return fmt.Errorf("tss emission percentage cannot be less than 0 percent")
	}
	return nil
}

func validateObserverSlashAmount(i interface{}) error {
	v, ok := i.(sdkmath.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.LT(sdk.ZeroInt()) {
		return fmt.Errorf("slash amount cannot be less than 0")
	}
	return nil
}

func validateBallotMaturityBlocks(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("ballot maturity types must be gte 0")
	}

	return nil
}

func validateBlockRewardsAmount(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.LT(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("block reward amount cannot be less than 0")
	}
	return nil
}
