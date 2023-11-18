package types

import (
	"fmt"
)

// DefaultGenesis returns the default fungible genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		ForeignCoinList: []ForeignCoin{},
		SystemContract:  nil,
		Params:          DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in foreignCoins
	foreignCoinsIndexMap := make(map[string]struct{})

	for _, elem := range gs.ForeignCoinList {
		index := string(ForeignCoinsKey(elem.Zrc20ContractAddress))
		if _, ok := foreignCoinsIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for foreignCoins")
		}
		foreignCoinsIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
