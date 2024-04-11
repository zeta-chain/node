package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"gopkg.in/yaml.v2"
)

func NewParams(observerParams []*ObserverParams, adminParams []*Admin_Policy, ballotMaturityBlocks int64) Params {
	return Params{
		ObserverParams:       observerParams,
		AdminPolicy:          adminParams,
		BallotMaturityBlocks: ballotMaturityBlocks,
	}
}

// DefaultParams returns a default set of parameters.
// privnet chains are supported by default for testing purposes
// custom params must be provided in genesis for other networks
func DefaultParams() Params {
	chains := chains.PrivnetChainList()
	observerParams := make([]*ObserverParams, len(chains))
	for i, chain := range chains {
		observerParams[i] = &ObserverParams{
			IsSupported:           true,
			Chain:                 chain,
			BallotThreshold:       sdk.MustNewDecFromStr("0.66"),
			MinObserverDelegation: sdk.MustNewDecFromStr("1000000000000000000000"), // 1000 ZETA
		}
	}
	return NewParams(observerParams, DefaultAdminPolicy(), 100)
}

func DefaultAdminPolicy() []*Admin_Policy {
	return []*Admin_Policy{
		{
			PolicyType: Policy_Type_group1,
			Address:    GroupID1Address,
		},
		{
			PolicyType: Policy_Type_group2,
			Address:    GroupID1Address,
		},
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateAdminPolicy(p.AdminPolicy); err != nil {
		return err
	}
	if err := validateBallotMaturityBlocks(p.BallotMaturityBlocks); err != nil {
		return err
	}
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, err := yaml.Marshal(p)
	if err != nil {
		return ""
	}
	return string(out)
}

// Deprecated: observer params are now stored in core params
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

func validateAdminPolicy(i interface{}) error {
	_, ok := i.([]*Admin_Policy)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

// https://github.com/zeta-chain/node/issues/1983
func validateBallotMaturityBlocks(i interface{}) error {
	_, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
