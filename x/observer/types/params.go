package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/zeta-chain/zetacore/common"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for zetaObserver module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

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
	chains := common.PrivnetChainList()
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

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyPrefix(ObserverParamsKey), &p.ObserverParams, validateVotingThresholds),
		paramtypes.NewParamSetPair(KeyPrefix(AdminPolicyParamsKey), &p.AdminPolicy, validateAdminPolicy),
		paramtypes.NewParamSetPair(KeyPrefix(BallotMaturityBlocksParamsKey), &p.BallotMaturityBlocks, validateBallotMaturityBlocks),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
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
func validateBallotMaturityBlocks(i interface{}) error {
	_, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func (p Params) GetAdminPolicyAccount(policyType Policy_Type) string {
	for _, admin := range p.AdminPolicy {
		if admin.PolicyType == policyType {
			return admin.Address
		}
	}
	return ""
}
func (p Params) GetParamsForChain(chain *common.Chain) ObserverParams {
	for _, ObserverParam := range p.GetObserverParams() {
		if ObserverParam.Chain.IsEqual(*chain) {
			return *ObserverParam
		}
	}
	return ObserverParams{}
}

// Deprecated: observer params are now stored in core params
func (p Params) GetChainFromChainName(name common.ChainName) *common.Chain {
	for _, observerParam := range p.GetObserverParams() {
		if observerParam.Chain.ChainName == name && observerParam.IsSupported {
			return observerParam.Chain
		}
	}
	return nil
}
