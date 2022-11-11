package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for zetaObserver module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{BallotThresholds: DefaultThreshold()}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{BallotThresholds: DefaultThreshold()}
}

func DefaultThreshold() []*BallotThreshold {
	chains := DefaultChainsList()
	threshold := make([]*BallotThreshold, len(chains)*2)
	for i, chain := range chains {
		threshold[i] = &BallotThreshold{
			Chain:       &Chain{ChainName: chain.ChainName, ChainId: chain.ChainId},
			Observation: ObservationType_InBoundTx,
			Threshold:   sdk.MustNewDecFromStr("0.66"),
		}
		i++
		threshold[i] = &BallotThreshold{
			Chain:       &Chain{ChainName: chain.ChainName, ChainId: chain.ChainId},
			Observation: ObservationType_OutBoundTx,
			Threshold:   sdk.MustNewDecFromStr("0.66"),
		}
	}
	return threshold
}
func DefaultChainsList() []*Chain {
	return []*Chain{
		{
			ChainName: ChainName_Eth,
			ChainId:   1,
		},
	}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyPrefix(ParamVotingThresholdsKey), &p.BallotThresholds, validateVotingThresholds),
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

func validateVotingThresholds(i interface{}) error {
	v, ok := i.([]*BallotThreshold)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	for _, threshold := range v {
		if threshold.Threshold.GT(sdk.OneDec()) {
			return ErrParamsThreshold
		}
	}
	return nil
}

func (p Params) GetVotingThreshold(chain Chain, observationType ObservationType) (BallotThreshold, bool) {
	for _, threshold := range p.GetBallotThresholds() {
		if threshold.Chain.IsEqual(chain) && threshold.Observation == observationType {
			return *threshold, true
		}
	}
	return BallotThreshold{}, false
}
