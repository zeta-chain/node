package types

// DefaultGenesis returns the default authority genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Policies:          DefaultPolicies(),
		ChainInfo:         DefaultChainInfo(),
		AuthorizationList: DefaultAuthorizationsList(),
	}
}

// Validate performs basic genesis state validation returning an error upon any failure
func (gs GenesisState) Validate() error {
	if err := gs.Policies.Validate(); err != nil {
		return err
	}
	if err := gs.AuthorizationList.Validate(); err != nil {
		return err
	}
	if err := gs.ChainInfo.Validate(); err != nil {
		return err
	}
	return nil
}
