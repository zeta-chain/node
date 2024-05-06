package types

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
)

// DefaultGenesis returns the default crosschain genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		OutboundTrackerList:   []OutboundTracker{},
		InboundHashToCctxList: []InboundHashToCctx{},
		GasPriceList:          []*GasPrice{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in outTxTracker
	outboundTrackerIndexMap := make(map[string]struct{})

	for _, elem := range gs.OutboundTrackerList {
		index := string(OutboundTrackerKey(elem.Index))
		if _, ok := outboundTrackerIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for outboundTracker")
		}
		outboundTrackerIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in inTxHashToCctx
	inboundHashToCctxIndexMap := make(map[string]struct{})

	for _, elem := range gs.InboundHashToCctxList {
		index := string(InboundHashToCctxKey(elem.InboundHash))
		if _, ok := inboundHashToCctxIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for inboundHashToCctx")
		}
		inboundHashToCctxIndexMap[index] = struct{}{}
	}

	// Check for duplicated index in gasPrice
	gasPriceIndexMap := make(map[string]bool)

	for _, elem := range gs.GasPriceList {
		if _, ok := gasPriceIndexMap[elem.Index]; ok {
			return fmt.Errorf("duplicated index for gasPrice")
		}
		gasPriceIndexMap[elem.Index] = true
	}

	return gs.RateLimiterFlags.Validate()
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
