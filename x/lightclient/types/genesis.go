package types

import (
	"fmt"

	"github.com/zeta-chain/node/pkg/proofs"
)

// DefaultGenesis returns the default lightclient genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		BlockHeaders:            []proofs.BlockHeader{},
		ChainStates:             []ChainState{},
		BlockHeaderVerification: BlockHeaderVerification{},
	}
}

// Validate performs basic genesis state validation returning an error upon any failure
func (gs GenesisState) Validate() error {
	// Check there is no duplicate
	blockHeaderMap := make(map[string]bool)
	for _, elem := range gs.BlockHeaders {
		if _, ok := blockHeaderMap[string(elem.Hash)]; ok {
			return fmt.Errorf("duplicated hash for block headers")
		}
		blockHeaderMap[string(elem.Hash)] = true
	}

	ChainStateMap := make(map[int64]bool)
	for _, elem := range gs.ChainStates {
		if _, ok := ChainStateMap[elem.ChainId]; ok {
			return fmt.Errorf("duplicated chain id for chain states")
		}
		ChainStateMap[elem.ChainId] = true
	}

	err := gs.BlockHeaderVerification.Validate()
	if err != nil {
		return err
	}

	return nil
}
