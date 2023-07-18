package types

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	// this line is used by starport scaffolding # ibc/genesistype/import
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # ibc/genesistype/default
		OutTxTrackerList:   []OutTxTracker{},
		InTxHashToCctxList: []InTxHashToCctx{},
		// this line is used by starport scaffolding # genesis/types/default
		Keygen:          nil,
		GasPriceList:    []*GasPrice{},
		ChainNoncesList: []*ChainNonces{},
		//CCTX:            []*Send{},

	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # ibc/genesistype/validate

	// Check for duplicated index in outTxTracker
	outTxTrackerIndexMap := make(map[string]struct{})

	for _, elem := range gs.OutTxTrackerList {
		index := string(OutTxTrackerKey(elem.Index))
		if _, ok := outTxTrackerIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for outTxTracker")
		}
		outTxTrackerIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in inTxHashToCctx
	inTxHashToCctxIndexMap := make(map[string]struct{})

	for _, elem := range gs.InTxHashToCctxList {
		index := string(InTxHashToCctxKey(elem.InTxHash))
		if _, ok := inTxHashToCctxIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for inTxHashToCctx")
		}
		inTxHashToCctxIndexMap[index] = struct{}{}
	}
	// TODO add migrate for TSS

	// Check for duplicated index in gasPrice
	gasPriceIndexMap := make(map[string]bool)

	for _, elem := range gs.GasPriceList {
		if _, ok := gasPriceIndexMap[elem.Index]; ok {
			return fmt.Errorf("duplicated index for gasPrice")
		}
		gasPriceIndexMap[elem.Index] = true
	}
	// Check for duplicated index in chainNonces
	chainNoncesIndexMap := make(map[string]bool)

	for _, elem := range gs.ChainNoncesList {
		if _, ok := chainNoncesIndexMap[elem.Index]; ok {
			return fmt.Errorf("duplicated index for chainNonces")
		}
		chainNoncesIndexMap[elem.Index] = true
	}

	// Check for duplicated index in send
	//sendIndexMap := make(map[string]bool)

	//for _, elem := range gs.SendList {
	//	if _, ok := sendIndexMap[elem.Index]; ok {
	//		return fmt.Errorf("duplicated index for send")
	//	}
	//	sendIndexMap[elem.Index] = true
	//}

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
