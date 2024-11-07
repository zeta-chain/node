package types

import (
	"fmt"
)

// DefaultGenesis returns the default fungible genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		ForeignCoinsList: []ForeignCoins{},
		SystemContract: &SystemContract{
			SystemContract: "0x91d18e54DAf4F677cB28167158d6dd21F6aB3921",
			ConnectorZevm:  "0x239e96c8f17C85c30100AC26F635Ea15f23E9c67",
		},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in foreignCoins
	foreignCoinsIndexMap := make(map[string]struct{})

	for _, elem := range gs.ForeignCoinsList {
		index := string(ForeignCoinsKey(elem.Zrc20ContractAddress))
		if _, ok := foreignCoinsIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for foreignCoins")
		}
		foreignCoinsIndexMap[index] = struct{}{}
	}

	return nil
}
