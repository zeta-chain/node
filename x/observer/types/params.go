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
	return Params{
		BallotThresholds: []*BallotThreshold{
			{
				Chain:       ObserverChain_BscTestnet,
				Observation: ObservationType_InBoundTx,
				Threshold:   sdk.MustNewDecFromStr("0.66"),
			},
			{
				Chain:       ObserverChain_BscTestnet,
				Observation: ObservationType_OutBoundTx,
				Threshold:   sdk.MustNewDecFromStr("0.66"),
			},
			{
				Chain:       ObserverChain_Goerli,
				Observation: ObservationType_InBoundTx,
				Threshold:   sdk.MustNewDecFromStr("0.66"),
			},
			{
				Chain:       ObserverChain_Goerli,
				Observation: ObservationType_OutBoundTx,
				Threshold:   sdk.MustNewDecFromStr("0.66"),
			},
			{
				Chain:       ObserverChain_Mumbai,
				Observation: ObservationType_InBoundTx,
				Threshold:   sdk.MustNewDecFromStr("0.66"),
			},
			{
				Chain:       ObserverChain_Mumbai,
				Observation: ObservationType_OutBoundTx,
				Threshold:   sdk.MustNewDecFromStr("0.66"),
			},
			{
				Chain:       ObserverChain_BTCTestnet,
				Observation: ObservationType_InBoundTx,
				Threshold:   sdk.MustNewDecFromStr("0.66"),
			},
			{
				Chain:       ObserverChain_BTCTestnet,
				Observation: ObservationType_OutBoundTx,
				Threshold:   sdk.MustNewDecFromStr("0.66"),
			},
			{
				Chain:       ObserverChain_Baobab,
				Observation: ObservationType_InBoundTx,
				Threshold:   sdk.MustNewDecFromStr("0.66"),
			},
			{
				Chain:       ObserverChain_Baobab,
				Observation: ObservationType_OutBoundTx,
				Threshold:   sdk.MustNewDecFromStr("0.66"),
			},
		},
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
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

func (p Params) GetVotingThreshold(chain ObserverChain, observationType ObservationType) (BallotThreshold, bool) {
	for _, threshold := range p.GetBallotThresholds() {
		if threshold.Chain == chain && threshold.Observation == observationType {
			return *threshold, true
		}
	}
	return BallotThreshold{}, false
}
