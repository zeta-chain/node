package types

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
)

// DefaultGenesis returns the default fungible genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		ForeignCoinsList: []ForeignCoins{},
		SystemContract:   nil,
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

func GetGenesisStateFromAppState(marshaler codec.JSONCodec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		err := marshaler.UnmarshalJSON(appState[ModuleName], &genesisState)
		if err != nil {
			panic(fmt.Sprintf("Failed to get genesis state from app state: %s", err.Error()))
		}
	}
	return genesisState
}

func GetGenesisStateFromAppStateLegacy(marshaler codec.JSONCodec, appState map[string]json.RawMessage) GenesisStateLegacy {
	var genesisState GenesisStateLegacy
	if appState[ModuleName] != nil {
		err := marshaler.UnmarshalJSON(appState[ModuleName], &genesisState)
		if err != nil {
			panic(fmt.Sprintf("Failed to get genesis state from app state: %s", err.Error()))
		}
	}
	return genesisState
}
