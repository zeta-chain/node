package types

// DefaultGenesis returns the default lightclient genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{}
}

// Validate performs basic genesis state validation returning an error upon any failure
func (gs GenesisState) Validate() error {
	return nil
}
