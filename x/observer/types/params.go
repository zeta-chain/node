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
		ObserverParams: []*ObserverParams{
			{
				Chain:                 ObserverChain_Eth,
				Observation:           ObservationType_InBoundTx,
				BallotThreshold:       sdk.MustNewDecFromStr("0.66"),
				MinObserverDelegation: sdk.MustNewDecFromStr("1000000000000000000000000"),
			},
			{
				Chain:                 ObserverChain_Eth,
				Observation:           ObservationType_OutBoundTx,
				BallotThreshold:       sdk.MustNewDecFromStr("0.66"),
				MinObserverDelegation: sdk.MustNewDecFromStr("1000000000000000000000000"),
			},
			{
				Chain:                 ObserverChain_BscMainnet,
				Observation:           ObservationType_InBoundTx,
				BallotThreshold:       sdk.MustNewDecFromStr("0.66"),
				MinObserverDelegation: sdk.MustNewDecFromStr("1000000000000000000000000"),
			},
			{
				Chain:                 ObserverChain_BscMainnet,
				Observation:           ObservationType_OutBoundTx,
				BallotThreshold:       sdk.MustNewDecFromStr("0.66"),
				MinObserverDelegation: sdk.MustNewDecFromStr("1000000000000000000000000"),
			},
			{
				Chain:                 ObserverChain_Goerli,
				Observation:           ObservationType_InBoundTx,
				BallotThreshold:       sdk.MustNewDecFromStr("0.66"),
				MinObserverDelegation: sdk.MustNewDecFromStr("1000000000000000000000000"),
			},
			{
				Chain:                 ObserverChain_Goerli,
				Observation:           ObservationType_OutBoundTx,
				BallotThreshold:       sdk.MustNewDecFromStr("0.66"),
				MinObserverDelegation: sdk.MustNewDecFromStr("1000000000000000000000000"),
			},
			{
				Chain:                 ObserverChain_Ropsten,
				Observation:           ObservationType_InBoundTx,
				BallotThreshold:       sdk.MustNewDecFromStr("0.66"),
				MinObserverDelegation: sdk.MustNewDecFromStr("1000000000000000000000000"),
			},
			{
				Chain:                 ObserverChain_Ropsten,
				Observation:           ObservationType_OutBoundTx,
				BallotThreshold:       sdk.MustNewDecFromStr("0.66"),
				MinObserverDelegation: sdk.MustNewDecFromStr("1000000000000000000000000"),
			},
			{
				Chain:                 ObserverChain_BscTestnet,
				Observation:           ObservationType_InBoundTx,
				BallotThreshold:       sdk.MustNewDecFromStr("0.66"),
				MinObserverDelegation: sdk.MustNewDecFromStr("1000000000000000000000000"),
			},
			{
				Chain:                 ObserverChain_BscTestnet,
				Observation:           ObservationType_OutBoundTx,
				BallotThreshold:       sdk.MustNewDecFromStr("0.66"),
				MinObserverDelegation: sdk.MustNewDecFromStr("1000000000000000000000000"),
			},
			{
				Chain:                 ObserverChain_Mumbai,
				Observation:           ObservationType_InBoundTx,
				BallotThreshold:       sdk.MustNewDecFromStr("0.66"),
				MinObserverDelegation: sdk.MustNewDecFromStr("1000000000000000000000000"),
			},
			{
				Chain:                 ObserverChain_Mumbai,
				Observation:           ObservationType_OutBoundTx,
				BallotThreshold:       sdk.MustNewDecFromStr("0.66"),
				MinObserverDelegation: sdk.MustNewDecFromStr("1000000000000000000000000"),
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
		paramtypes.NewParamSetPair(KeyPrefix(ParamVotingThresholdsKey), &p.ObserverParams, validateVotingThresholds),
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
	v, ok := i.([]*ObserverParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	for _, threshold := range v {
		if threshold.BallotThreshold.GT(sdk.OneDec()) {
			return ErrParamsThreshold
		}
	}
	return nil
}

func (p Params) GetParamsForChainAndType(chain ObserverChain, observationType ObservationType) (ObserverParams, bool) {
	for _, threshold := range p.GetObserverParams() {
		if threshold.Chain == chain && threshold.Observation == observationType {
			return *threshold, true
		}
	}
	return ObserverParams{}, false
}
