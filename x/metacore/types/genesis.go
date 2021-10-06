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
		TxoutList:       []*Txout{},
		NodeAccountList: []*NodeAccount{},
		TxinVoterList:   []*TxinVoter{},
		TxinList:        []*Txin{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # ibc/genesistype/validate

	// this line is used by starport scaffolding # genesis/types/validate
	// Check for duplicated ID in txout
	txoutIdMap := make(map[uint64]bool)

	for _, elem := range gs.TxoutList {
		if _, ok := txoutIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for txout")
		}
		txoutIdMap[elem.Id] = true
	}
	// Check for duplicated index in nodeAccount
	nodeAccountIndexMap := make(map[string]bool)

	for _, elem := range gs.NodeAccountList {
		if _, ok := nodeAccountIndexMap[elem.Index]; ok {
			return fmt.Errorf("duplicated index for nodeAccount")
		}
		nodeAccountIndexMap[elem.Index] = true
	}
	// Check for duplicated index in txinVoter
	txinVoterIndexMap := make(map[string]bool)

	for _, elem := range gs.TxinVoterList {
		if _, ok := txinVoterIndexMap[elem.Index]; ok {
			return fmt.Errorf("duplicated index for txinVoter")
		}
		txinVoterIndexMap[elem.Index] = true
	}
	// Check for duplicated index in txin
	txinIndexMap := make(map[string]bool)

	for _, elem := range gs.TxinList {
		if _, ok := txinIndexMap[elem.Index]; ok {
			return fmt.Errorf("duplicated index for txin")
		}
		txinIndexMap[elem.Index] = true
	}

	return nil
}
