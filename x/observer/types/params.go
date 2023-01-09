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

func NewParams(observerParams []*ObserverParams) Params {
	return Params{ObserverParams: observerParams}
}
func DefaultParams() Params {
	chains := DefaultChainsList()
	observerParams := make([]*ObserverParams, len(chains)*2)
	i := 0
	for _, chain := range chains {
		observerParams[i] = &ObserverParams{
			Chain:                 chain,
			Observation:           ObservationType_InBoundTx,
			BallotThreshold:       sdk.MustNewDecFromStr("0.66"),
			MinObserverDelegation: sdk.MustNewDecFromStr("10000000000"),
		}
		i++
		observerParams[i] = &ObserverParams{
			Chain:                 chain,
			Observation:           ObservationType_OutBoundTx,
			BallotThreshold:       sdk.MustNewDecFromStr("0.66"),
			MinObserverDelegation: sdk.MustNewDecFromStr("10000000000"),
		}
		i++
	}
	return NewParams(observerParams)
}

func DefaultChainsList() []*Chain {
	return []*Chain{
		{
			ChainName: ChainName_Eth,
			ChainId:   1,
		},
		{
			ChainName: ChainName_Goerli,
			ChainId:   5,
		},
		{
			ChainName: ChainName_Ropsten,
			ChainId:   3,
		},
		{
			ChainName: ChainName_BscMainnet,
			ChainId:   56,
		},
		{
			ChainName: ChainName_BscTestnet,
			ChainId:   97,
		},
		{
			ChainName: ChainName_Baobab,
			ChainId:   1001,
		},
		{
			ChainName: ChainName_ZetaChain,
			ChainId:   2374,
		},
		{
			ChainName: ChainName_Btc,
			ChainId:   55555,
		},
		{
			ChainName: ChainName_Polygon,
			ChainId:   137,
		},
		{
			ChainName: ChainName_Mumbai,
			ChainId:   80001,
		},
	}
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

func (p Params) GetParamsForChainAndType(chain *Chain, observationType ObservationType) (ObserverParams, bool) {
	for _, ObserverParam := range p.GetObserverParams() {
		fmt.Println(ObserverParam.String())
		if ObserverParam.Chain.IsEqual(chain) && ObserverParam.Observation == observationType {
			return *ObserverParam, true
		}
	}
	return ObserverParams{}, false
}
