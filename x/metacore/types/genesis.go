package types

import (
	"fmt"
	// this line is used by starport scaffolding # ibc/genesistype/import
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # ibc/genesistype/default
		// this line is used by starport scaffolding # genesis/types/default
		ChainNoncesList:     []*ChainNonces{},
		LastBlockHeightList: []*LastBlockHeight{},
		ReceiveList:         []*Receive{},
		SendList:            []*Send{},
		NodeAccountList:     []*NodeAccount{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # ibc/genesistype/validate

	// this line is used by starport scaffolding # genesis/types/validate
	// Check for duplicated index in chainNonces
	chainNoncesIndexMap := make(map[string]bool)

	for _, elem := range gs.ChainNoncesList {
		if _, ok := chainNoncesIndexMap[elem.Index]; ok {
			return fmt.Errorf("duplicated index for chainNonces")
		}
		chainNoncesIndexMap[elem.Index] = true
	}
	// Check for duplicated index in lastBlockHeight
	lastBlockHeightIndexMap := make(map[string]bool)

	for _, elem := range gs.LastBlockHeightList {
		if _, ok := lastBlockHeightIndexMap[elem.Index]; ok {
			return fmt.Errorf("duplicated index for lastBlockHeight")
		}
		lastBlockHeightIndexMap[elem.Index] = true
	}
	// Check for duplicated index in receive
	receiveIndexMap := make(map[string]bool)

	for _, elem := range gs.ReceiveList {
		if _, ok := receiveIndexMap[elem.Index]; ok {
			return fmt.Errorf("duplicated index for receive")
		}
		receiveIndexMap[elem.Index] = true
	}
	// Check for duplicated index in send
	sendIndexMap := make(map[string]bool)

	for _, elem := range gs.SendList {
		if _, ok := sendIndexMap[elem.Index]; ok {
			return fmt.Errorf("duplicated index for send")
		}
		sendIndexMap[elem.Index] = true
	}

	// Check for duplicated index in nodeAccount
	nodeAccountIndexMap := make(map[string]bool)

	for _, elem := range gs.NodeAccountList {
		if _, ok := nodeAccountIndexMap[elem.Index]; ok {
			return fmt.Errorf("duplicated index for nodeAccount")
		}
		nodeAccountIndexMap[elem.Index] = true
	}

	return nil
}
