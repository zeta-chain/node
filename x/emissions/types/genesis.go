package types

// DefaultIndex is the default emissions global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default emissions genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}
